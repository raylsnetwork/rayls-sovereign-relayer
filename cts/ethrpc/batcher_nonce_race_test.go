// Regression tests for the CTS keyQueue nonce race.
//
// The bug:
//   Two signing paths for the SAME key (RPCBatcher.Send and AuthGen, both fed
//   by the keyQueue) resolve the nonce via an unsynchronised eth_getTransactionCount
//   (PendingNonceAt) read. Two such reads that interleave before either tx is
//   broadcast return the SAME value. PendingNonceAt is mempool-aware, so the
//   collision window is read-to-broadcast, NOT broadcast-to-mining. Both paths
//   then sign DIFFERENT calldata at the same (sender, nonce). On a chain that
//   permits 0-fee same-nonce replacement (axyl/Reth; Geth/Besu reject it as
//   "replacement underpriced"), the second tx evicts the first from the mempool,
//   so one is silently dropped, and CTS records the dropped tx's hash and polls it
//   forever; the message is permanently lost from CTS's view.
//
//   NOTE: the stubs below model the collision by returning the SAME nonce from
//   PendingNonceAt across reads: the simplest reproduction of the interleaved
//   read. The fix is independent of which sub-mechanism produced the duplicate
//   read; it removes the read-modify-write entirely.
//
// The fix: maintain a per-address local nonce cache inside RPCBatcher.
// First use of a key seeds the cache from the chain. Every Send reserves
// `len(inputs)` contiguous nonces atomically and advances the cache. Error
// paths invalidate the cache so the next Send re-queries the chain.
//
// These tests exercise the contract:
//   - No two outputs from sequential Sends with the same dequeued key share
//     a nonce, even when PendingNonceAt is stubbed to return the same value.
//   - First use seeds from chain; subsequent uses do NOT re-query.
//   - Batch size N advances cache by exactly N.
//   - Different keys maintain independent caches.
//   - Error paths invalidate the cache (next Send re-queries).
//
// Each test should FAIL against the pre-Fix-A `RPCBatcher.Send` implementation
// (which calls getNonceForPrivateKey on every Send) and PASS post-Fix.

package ethrpc_test

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/ethrpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Fixtures ──────────────────────────────────────────────────────────────

// Deterministic test keys (anvil-style). Useful for sanity checks where
// we want to verify the sender address.
const (
	testPK0 = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	testPK1 = "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"
)

func mustParsePK(t *testing.T, h string) *ecdsa.PrivateKey {
	t.Helper()
	pk, err := crypto.HexToECDSA(h)
	require.NoError(t, err)
	return pk
}

// newBatcherWithMocks wires a RPCBatcher with three mocks pre-set to the
// minimum needed for NewBatcher() to succeed. The test can then override
// PendingNonceAtFunc / BatchCallContextFunc / Dequeue/EnqueueFunc as needed.
func newBatcherWithMocks(t *testing.T, key *ecdsa.PrivateKey) (
	*ethrpc.RPCBatcher,
	*EthereumClientMock,
	*RPCClientMock,
	*KeyQueueMock,
) {
	t.Helper()
	ethMock := &EthereumClientMock{
		ChainIDFunc: func(ctx context.Context) (*big.Int, error) {
			return big.NewInt(12345), nil
		},
		SuggestGasPriceFunc: func(ctx context.Context) (*big.Int, error) {
			return big.NewInt(0), nil
		},
		EstimateGasFunc: func(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
			return 100_000, nil
		},
	}
	rpcMock := &RPCClientMock{
		BatchCallContextFunc: func(ctx context.Context, elems []rpc.BatchElem) error {
			// Default: every tx broadcasts successfully. Return a fresh
			// hash per request so the test doesn't accidentally collide.
			for i := range elems {
				h := common.BigToHash(big.NewInt(int64(i + 1))).Hex()
				p := elems[i].Result.(*string)
				*p = h
			}
			return nil
		},
	}
	keyQueueMock := &KeyQueueMock{
		DequeueFunc: func(ctx context.Context) (*ecdsa.PrivateKey, error) {
			return key, nil
		},
		EnqueueFunc: func(privateKey *ecdsa.PrivateKey) {},
	}
	b, err := ethrpc.NewBatcher(context.Background(), keyQueueMock, ethMock, rpcMock)
	require.NoError(t, err)
	return b, ethMock, rpcMock, keyQueueMock
}

// txInputs builds N TransactionInput entries with deterministic per-input
// calldata so that signed RLPs differ across positions (mirrors how the
// real cts/batcher pipeline produces inputs).
func txInputs(n int) []ethrpc.TransactionInput {
	out := make([]ethrpc.TransactionInput, n)
	for i := 0; i < n; i++ {
		out[i] = ethrpc.TransactionInput{
			ID:      fmt.Sprintf("corr-%d", i),
			Data:    []byte{byte(i), 0xff},
			Address: common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
		}
	}
	return out
}

// decodeRawTx pulls a signed tx out of an rpc.BatchElem's args (the
// "0xRLP" string passed to eth_sendRawTransaction) so the test can read
// the on-wire nonce.
func decodeRawTx(t *testing.T, e rpc.BatchElem) *types.Transaction {
	t.Helper()
	args, ok := e.Args[0].(string)
	require.True(t, ok, "first arg should be 0x-hex rlp string")
	raw, err := hex.DecodeString(args[2:]) // strip "0x"
	require.NoError(t, err)
	tx := new(types.Transaction)
	require.NoError(t, tx.UnmarshalBinary(raw))
	return tx
}

// recordBroadcasts captures every signed tx that the batcher submits via
// BatchCallContext. Used by tests to assert on the wire-level nonces.
func recordBroadcasts(rpcMock *RPCClientMock, t *testing.T) *atomic.Value {
	var bucket atomic.Value
	bucket.Store(([]*types.Transaction)(nil))
	prev := rpcMock.BatchCallContextFunc
	rpcMock.BatchCallContextFunc = func(ctx context.Context, elems []rpc.BatchElem) error {
		current := bucket.Load().([]*types.Transaction)
		for i := range elems {
			current = append(current, decodeRawTx(t, elems[i]))
		}
		bucket.Store(current)
		return prev(ctx, elems)
	}
	return &bucket
}

func nonces(txs []*types.Transaction) []uint64 {
	out := make([]uint64, len(txs))
	for i, tx := range txs {
		out[i] = tx.Nonce()
	}
	return out
}

// ── Tests ─────────────────────────────────────────────────────────────────

// TestSend_NoNonceReuseAcrossSendsOnSameKey is the core nonce-reuse
// regression. Stub PendingNonceAt to ALWAYS return 100 (simulating the
// "first tx hasn't mined yet" condition that the real chain exhibits
// between rapid back-to-back Sends). Two single-input Sends back-to-back
// using the same dequeued key MUST produce txs with DIFFERENT nonces.
//
// Without the per-key local nonce cache: BOTH txs get nonce 100 → fails.
// With the cache: first gets 100, second gets 101 → passes.
func TestSend_NoNonceReuseAcrossSendsOnSameKey(t *testing.T) {
	key := mustParsePK(t, testPK0)
	b, ethMock, rpcMock, _ := newBatcherWithMocks(t, key)
	// The chain reports nonce 100 on both calls, simulating "first tx
	// not yet mined" between two consecutive Sends.
	ethMock.PendingNonceAtFunc = func(ctx context.Context, addr common.Address) (uint64, error) {
		return 100, nil
	}
	bucket := recordBroadcasts(rpcMock, t)

	_, err := b.Send(context.Background(), txInputs(1))
	require.NoError(t, err)
	_, err = b.Send(context.Background(), txInputs(1))
	require.NoError(t, err)

	got := nonces(bucket.Load().([]*types.Transaction))
	t.Logf("observed broadcast nonces: %v", got)
	t.Logf("PendingNonceAt was called %d time(s) — should be 1 (cache hit on 2nd Send)",
		len(ethMock.PendingNonceAtCalls()))

	require.Lenf(t, got, 2, "expected 2 broadcasts, got %d", len(got))
	assert.NotEqualf(t, got[0], got[1],
		"regression: two back-to-back Sends using the same dequeued key "+
			"MUST NOT sign at the same nonce. Got nonce=%d twice. The fix is a "+
			"per-key local nonce cache so the second Send doesn't re-query "+
			"PendingNonceAt and reuse the stale on-chain value.",
		got[0])
	assert.Equalf(t, uint64(100), got[0], "first Send should use the on-chain nonce as the seed")
	assert.Equalf(t, uint64(101), got[1], "second Send should advance from the cache (no re-query)")
}

// TestSend_FirstUseSeedsFromChain_SecondUseDoesNot verifies the cache
// behavior at the chain-RPC level. First Send queries PendingNonceAt;
// every subsequent Send for the same key uses the cache without
// re-querying — the chief lever for fixing the race.
func TestSend_FirstUseSeedsFromChain_SecondUseDoesNot(t *testing.T) {
	key := mustParsePK(t, testPK0)
	b, ethMock, _, _ := newBatcherWithMocks(t, key)
	ethMock.PendingNonceAtFunc = func(ctx context.Context, addr common.Address) (uint64, error) {
		return 7, nil
	}

	// 5 sequential Sends, single input each.
	for i := 0; i < 5; i++ {
		_, err := b.Send(context.Background(), txInputs(1))
		require.NoError(t, err, "Send %d", i)
	}

	calls := len(ethMock.PendingNonceAtCalls())
	t.Logf("PendingNonceAt calls across 5 Sends: %d", calls)
	assert.Equalf(t, 1, calls,
		"PendingNonceAt should be called exactly once across N Sends for the "+
			"same key — first use seeds the cache, subsequent uses don't re-query. "+
			"Pre-Fix-A this would be %d.", 5)
}

// TestSend_BatchAdvancesNonceCacheBySize verifies a Send with N inputs
// reserves N contiguous nonces (the same invariant the pre-fix code had
// internally within signBatch, now externalised to the cache).
func TestSend_BatchAdvancesNonceCacheBySize(t *testing.T) {
	key := mustParsePK(t, testPK0)
	b, ethMock, rpcMock, _ := newBatcherWithMocks(t, key)
	ethMock.PendingNonceAtFunc = func(ctx context.Context, addr common.Address) (uint64, error) {
		return 50, nil
	}
	bucket := recordBroadcasts(rpcMock, t)

	_, err := b.Send(context.Background(), txInputs(3)) // expect nonces 50, 51, 52
	require.NoError(t, err)
	_, err = b.Send(context.Background(), txInputs(2)) // expect 53, 54
	require.NoError(t, err)

	got := nonces(bucket.Load().([]*types.Transaction))
	want := []uint64{50, 51, 52, 53, 54}
	t.Logf("observed: %v  expected: %v", got, want)
	assert.Equalf(t, want, got,
		"batch of 3 then batch of 2 should produce 5 contiguous nonces starting at "+
			"the on-chain seed; got %v", got)
}

// TestSend_DifferentKeysIndependent verifies each key gets its own cache:
// PendingNonceAt is called once PER key, the nonces are independent.
func TestSend_DifferentKeysIndependent(t *testing.T) {
	keyA := mustParsePK(t, testPK0)
	keyB := mustParsePK(t, testPK1)
	addrA := crypto.PubkeyToAddress(keyA.PublicKey)
	addrB := crypto.PubkeyToAddress(keyB.PublicKey)

	b, ethMock, rpcMock, kq := newBatcherWithMocks(t, keyA)
	// Per-key on-chain nonce; PendingNonceAt is keyed by address.
	ethMock.PendingNonceAtFunc = func(ctx context.Context, addr common.Address) (uint64, error) {
		switch addr {
		case addrA:
			return 10, nil
		case addrB:
			return 200, nil
		default:
			return 0, fmt.Errorf("unexpected address %s", addr.Hex())
		}
	}
	// Alternate the dequeued key: A, B, A, B.
	seq := []*ecdsa.PrivateKey{keyA, keyB, keyA, keyB}
	idx := 0
	kq.DequeueFunc = func(ctx context.Context) (*ecdsa.PrivateKey, error) {
		k := seq[idx]
		idx++
		return k, nil
	}
	bucket := recordBroadcasts(rpcMock, t)

	for i := 0; i < 4; i++ {
		_, err := b.Send(context.Background(), txInputs(1))
		require.NoError(t, err)
	}

	got := nonces(bucket.Load().([]*types.Transaction))
	want := []uint64{10, 200, 11, 201}
	t.Logf("observed: %v  expected: %v", got, want)
	assert.Equal(t, want, got, "keys should have independent caches: A→10,11  B→200,201")
	assert.Equal(t, 2, len(ethMock.PendingNonceAtCalls()),
		"PendingNonceAt should be called twice: once per distinct key")
}

// TestSend_WholeBatchErrorInvalidatesCache: when BatchCallContext returns
// an error (no tx was broadcast), the cache for that address must be
// invalidated so the next Send re-queries the chain and doesn't stride
// past unused nonces.
func TestSend_WholeBatchErrorInvalidatesCache(t *testing.T) {
	key := mustParsePK(t, testPK0)
	b, ethMock, rpcMock, _ := newBatcherWithMocks(t, key)
	ethMock.PendingNonceAtFunc = func(ctx context.Context, addr common.Address) (uint64, error) {
		return 5, nil
	}
	// First call: fail the whole-batch RPC. Second call: succeed.
	callIdx := 0
	rpcMock.BatchCallContextFunc = func(ctx context.Context, elems []rpc.BatchElem) error {
		callIdx++
		if callIdx == 1 {
			return errors.New("simulated batch-RPC failure")
		}
		for i := range elems {
			h := common.BigToHash(big.NewInt(int64(i + 1))).Hex()
			p := elems[i].Result.(*string)
			*p = h
		}
		return nil
	}
	bucket := recordBroadcasts(rpcMock, t)

	// First Send: errors out on broadcast.
	_, err := b.Send(context.Background(), txInputs(1))
	require.Error(t, err)
	// Second Send: should re-query the chain (cache invalidated) and
	// observe the same on-chain nonce 5.
	_, err = b.Send(context.Background(), txInputs(1))
	require.NoError(t, err)

	got := nonces(bucket.Load().([]*types.Transaction))
	t.Logf("post-error broadcast nonces: %v", got)
	calls := len(ethMock.PendingNonceAtCalls())
	t.Logf("PendingNonceAt calls: %d (expect 2 — one before error, one after)", calls)

	// Only the second Send produced an actual broadcast (recordBroadcasts
	// captures by reading rpc.BatchElem args, and the failing batch
	// still went through BatchCallContext — so we DO see both attempts'
	// signed bytes in the bucket).
	require.GreaterOrEqual(t, len(got), 2, "should have observed 2 signed-tx attempts")
	assert.Equalf(t, uint64(5), got[len(got)-1],
		"after whole-batch error, next Send must re-seed from chain (nonce=5), got %d",
		got[len(got)-1])
	assert.Equal(t, 2, calls, "PendingNonceAt should have been re-queried after the error")
}

// TestSend_AllSignsFailInvalidatesCache: every input's signAndEncode
// returns either signErr or revert (no broadcast happens at all). The
// cache must be invalidated so the next Send doesn't think nonces were
// consumed.
//
// We trigger "every input fails sign" by making EstimateGas return a
// revert error consistently. The pre-flight-revert path doesn't advance
// the on-chain nonce; the local cache must follow suit.
func TestSend_AllSignsFailInvalidatesCache(t *testing.T) {
	key := mustParsePK(t, testPK0)
	b, ethMock, _, _ := newBatcherWithMocks(t, key)
	pendingCalls := 0
	ethMock.PendingNonceAtFunc = func(ctx context.Context, addr common.Address) (uint64, error) {
		pendingCalls++
		return 42, nil
	}
	// Make estimation fail without a revert blob → caller treats as
	// retryable infra failure (signErr in signBatch).
	ethMock.EstimateGasFunc = func(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
		return 0, errors.New("simulated estimate failure")
	}

	res, err := b.Send(context.Background(), txInputs(2))
	require.NoError(t, err) // Send itself returns nil because per-tx errors
	// land in the result slice, not the function-level err.
	require.Len(t, res, 2)
	assert.NotNil(t, res[0].Err, "input 0 should have a per-tx error")
	assert.NotNil(t, res[1].Err, "input 1 should have a per-tx error")
	t.Logf("PendingNonceAt called %d times on the first (all-fail) Send", pendingCalls)

	// Next Send must re-query chain (cache invalidated by all-fail path).
	_, err = b.Send(context.Background(), txInputs(1))
	require.NoError(t, err)
	t.Logf("PendingNonceAt called %d times TOTAL after a successful Send post-failure", pendingCalls)
	assert.Equalf(t, 2, pendingCalls,
		"on the failing first Send all signs errored → no broadcast → cache must be "+
			"invalidated → the next Send must re-query PendingNonceAt. Saw %d calls.",
		pendingCalls)
}

// TestSend_PerTxBroadcastErrInvalidatesCache: if at least one per-tx in
// the broadcast batch returns req.Error (e.g., "nonce too low"), the
// cache must be invalidated. This is the conservative recovery for a
// chain that rejected our claimed nonce.
func TestSend_PerTxBroadcastErrInvalidatesCache(t *testing.T) {
	key := mustParsePK(t, testPK0)
	b, ethMock, rpcMock, _ := newBatcherWithMocks(t, key)
	ethMock.PendingNonceAtFunc = func(ctx context.Context, addr common.Address) (uint64, error) {
		return 8, nil
	}
	// First batch: one per-tx error, one success.
	callIdx := 0
	rpcMock.BatchCallContextFunc = func(ctx context.Context, elems []rpc.BatchElem) error {
		callIdx++
		if callIdx == 1 {
			elems[0].Error = errors.New("simulated nonce-too-low")
			p := elems[1].Result.(*string)
			*p = common.BigToHash(big.NewInt(2)).Hex()
			return nil
		}
		for i := range elems {
			h := common.BigToHash(big.NewInt(int64(i + 10))).Hex()
			p := elems[i].Result.(*string)
			*p = h
		}
		return nil
	}

	_, err := b.Send(context.Background(), txInputs(2))
	require.NoError(t, err)
	_, err = b.Send(context.Background(), txInputs(1))
	require.NoError(t, err)

	calls := len(ethMock.PendingNonceAtCalls())
	t.Logf("PendingNonceAt calls: %d (expect 2 — invalidated after per-tx error)", calls)
	assert.Equalf(t, 2, calls,
		"a 'nonce too low' per-tx broadcast error must invalidate the cache so the "+
			"next Send re-queries on-chain (the chain has moved ahead). Saw %d calls.", calls)
}

// TestSend_AlreadyKnownPerTxErr_KeepsCacheForward: an "already known" per-tx
// error means the node already holds OUR identical tx in its mempool (an
// idempotent re-broadcast). The nonce slot is occupied by our own tx, so the
// cache is correctly positioned: re-querying could only regress it below a
// nonce already handed out and re-open the cross-Send collision the cache
// prevents. So unlike "nonce too low", "already known" must NOT invalidate.
func TestSend_AlreadyKnownPerTxErr_KeepsCacheForward(t *testing.T) {
	// Both wire strings mean the node already holds our identical tx: go-ethereum
	// returns "already known"; axyl/Reth returns "already known" too (it remaps its
	// internal "already imported" pool error to the geth-compatible RPC string), and
	// "already imported" is matched defensively. Neither must invalidate the cache.
	for _, errMsg := range []string{"already known", "already imported"} {
		t.Run(errMsg, func(t *testing.T) {
			key := mustParsePK(t, testPK0)
			b, ethMock, rpcMock, _ := newBatcherWithMocks(t, key)
			ethMock.PendingNonceAtFunc = func(ctx context.Context, addr common.Address) (uint64, error) {
				return 8, nil
			}
			// First batch (2 inputs): the first tx comes back as a duplicate of our own
			// tx, the second broadcasts fine. Set the custom func BEFORE recordBroadcasts
			// so the recorder wraps it.
			callIdx := 0
			rpcMock.BatchCallContextFunc = func(ctx context.Context, elems []rpc.BatchElem) error {
				callIdx++
				if callIdx == 1 {
					elems[0].Error = errors.New(errMsg)
					p := elems[1].Result.(*string)
					*p = common.BigToHash(big.NewInt(2)).Hex()
					return nil
				}
				for i := range elems {
					h := common.BigToHash(big.NewInt(int64(i + 10))).Hex()
					p := elems[i].Result.(*string)
					*p = h
				}
				return nil
			}
			bucket := recordBroadcasts(rpcMock, t)

			_, err := b.Send(context.Background(), txInputs(2)) // reserves 8,9
			require.NoError(t, err)
			_, err = b.Send(context.Background(), txInputs(1)) // must continue from cache (10), NOT re-query
			require.NoError(t, err)

			calls := len(ethMock.PendingNonceAtCalls())
			t.Logf("PendingNonceAt calls: %d (expect 1; %q must NOT invalidate)", calls, errMsg)
			assert.Equalf(t, 1, calls,
				"a %q per-tx error must NOT invalidate the cache: the slot holds our own "+
					"identical tx, and a reseed could only regress below a handed-out nonce. "+
					"Saw %d calls.", errMsg, calls)

			got := nonces(bucket.Load().([]*types.Transaction))
			t.Logf("observed broadcast nonces: %v", got)
			// Broadcasts: batch1 tx@8 (errored but signed), tx@9; batch2 tx@10.
			require.GreaterOrEqual(t, len(got), 3, "expected 3 signed-tx attempts")
			assert.Equalf(t, uint64(10), got[len(got)-1],
				"second Send must advance from the cache (nonce=10), not re-seed; got %d", got[len(got)-1])
		})
	}
}
