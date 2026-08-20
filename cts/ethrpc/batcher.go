package ethrpc

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/raylsnetwork/rayls-sovereign-relayer/ethretry"
	"github.com/raylsnetwork/rayls-sovereign-relayer/noncecache"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

const (
	batchRetryTimeout  = 5 * time.Minute
	gasEstimateTimeout = 30 * time.Second
	rpcCallTimeout     = time.Minute
	// batcherInitCallTimeout bounds each ChainID/SuggestGasPrice attempt made
	// while constructing a batcher, so a stalled RPC is aborted and retried
	// (via WithRetry) rather than blocking until the caller's whole
	// construction budget runs out.
	batcherInitCallTimeout = 10 * time.Second
)

// TransactionInput is one leg of a batch: the calldata to broadcast and
// the destination address. ID is opaque to the batcher and only used
// for logging / error correlation.
type TransactionInput struct {
	ID      string
	Data    []byte
	Address common.Address
}

// Revert carries raw ABI-encoded revert bytes. Callers decode against
// their own binding error unpackers — same contract as cts/service/txops.go.
type Revert struct {
	Data []byte
}

// SendResult is the per-transaction outcome of Send. Exactly one arm is set:
//
//	Hash   — broadcast accepted by the node, tx is in the mempool
//	Revert — gas estimation caught a pre-flight revert; retrying the
//	         same calldata will produce the same revert
//	Err    — retryable infra failure specific to this tx (RPC, signing)
//
// Whole-batch retryable failures (BatchCallContext failed after all
// retries) are returned as Send's function-level error, not folded here.
type SendResult struct {
	Hash   common.Hash
	Revert *Revert
	Err    error
}

// ReceiptResult is the per-hash outcome of GetReceipts. Exactly one arm
// is set:
//
//	Receipt — status == 1 (mined successfully)
//	Revert  — status == 0 (mined and reverted); Data is the raw revert
//	          bytes extracted by replaying the tx at the mined block
//	Pending — no receipt yet; caller decides when to re-poll
//	Err     — retryable failure fetching or decoding this receipt
type ReceiptResult struct {
	Receipt *types.Receipt
	Revert  *Revert
	Pending bool
	Err     error
}

//go:generate moq --pkg ethrpc_test -out batcher_mock_test.go . EthereumClient RPCClient KeyQueue
type EthereumClient interface {
	ChainID(context.Context) (*big.Int, error)
	SuggestGasPrice(context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	PendingNonceAt(context.Context, common.Address) (uint64, error)

	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)
}

type RPCClient interface {
	BatchCallContext(context.Context, []rpc.BatchElem) error
}

type KeyQueue interface {
	Enqueue(*ecdsa.PrivateKey)
	Dequeue(context.Context) (*ecdsa.PrivateKey, error)
}

type RPCBatcher struct {
	keyQueue  KeyQueue
	ethClient EthereumClient
	rpcClient RPCClient

	chainID  *big.Int
	gasPrice *big.Int

	// Shared per-address local nonce counter.
	//
	// Two signing paths for the SAME key (this RPCBatcher and AuthGen, both
	// fed by one `keyQueue`) used to resolve the nonce via an unsynchronised
	// `eth_getTransactionCount` read. Two such reads that interleave before
	// either tx is broadcast return the same value. `PendingNonceAt` is
	// mempool-aware, so the collision window is read-to-broadcast, NOT
	// broadcast-to-mining. Both paths then sign DIFFERENT calldata at the
	// same (sender, nonce). On a chain that permits 0-fee same-nonce
	// replacement (axyl/Reth; Geth/Besu reject it as "replacement
	// underpriced"), the second tx evicts the first from the mempool, so one
	// is silently dropped, and CTS records the dropped tx's hash and polls it
	// forever.
	//
	// The cache itself is the noncecache.Cache (see that package for the
	// invariants); we hold it through a pointer so multiple components
	// signing for the same identity (AuthGen, RPCBatcher, reaper) can
	// SHARE a single counter. Without that sharing the within-RPCBatcher
	// race is closed but the cross-component race (e.g. AuthGen's
	// SignAndSend and the async pipeline's RPCBatcher both reading
	// PendingNonceAt for the same key while a prior tx is pending) stays
	// open. NewBatcher creates an own cache for legacy callers; the
	// production wiring in cts/cmd/run/app.go uses NewBatcherWithCache
	// to pass in the per-identity instance.
	nonceCache *noncecache.Cache
}

// NewBatcher snapshots chainID and gasPrice at construction. gasPrice
// is stable on the gasless network this batcher targets, so no refresh
// loop is needed.
//
// The returned batcher owns its own NonceCache. For the production
// wiring that shares the cache between AuthGen + RPCBatcher (and any
// other component signing for the same identity) use NewBatcherWithCache.
func NewBatcher(
	ctx context.Context,
	keyQueue KeyQueue,
	ethClient EthereumClient,
	rpcClient RPCClient,
) (*RPCBatcher, error) {
	return NewBatcherWithCache(ctx, keyQueue, ethClient, rpcClient, noncecache.New(ethClient))
}

// NewBatcherWithCache is like NewBatcher but threads in an externally-owned
// nonce cache so multiple components signing for the same identity share
// a single counter. Required to close the cross-component nonce race —
// see noncecache.Cache for the invariants.
func NewBatcherWithCache(
	ctx context.Context,
	keyQueue KeyQueue,
	ethClient EthereumClient,
	rpcClient RPCClient,
	cache *noncecache.Cache,
) (*RPCBatcher, error) {
	if cache == nil {
		return nil, fmt.Errorf("nonce cache must not be nil")
	}
	// chainID and gasPrice are fetched once here, over the RPC link (VPN-routed
	// and occasionally flaky in remote dev). Retry on transient transport
	// errors, bounding each attempt with its own timeout, so a single slow or
	// stalled call can't consume the whole construction budget and abort CTS
	// startup — the send path already uses WithRetry for the same reason.
	var chainID *big.Int
	if err := ethretry.WithRetry(ctx, func() error {
		callCtx, cancel := context.WithTimeout(ctx, batcherInitCallTimeout)
		defer cancel()
		var err error
		chainID, err = ethClient.ChainID(callCtx)
		return err
	}); err != nil {
		return nil, withstack.Wrap(fmt.Errorf("fetching chain id: %w", err))
	}

	var gasPrice *big.Int
	if err := ethretry.WithRetry(ctx, func() error {
		callCtx, cancel := context.WithTimeout(ctx, batcherInitCallTimeout)
		defer cancel()
		var err error
		gasPrice, err = ethClient.SuggestGasPrice(callCtx)
		return err
	}); err != nil {
		return nil, withstack.Wrap(fmt.Errorf("fetching gas price: %w", err))
	}

	return &RPCBatcher{
		keyQueue:   keyQueue,
		ethClient:  ethClient,
		rpcClient:  rpcClient,
		chainID:    chainID,
		gasPrice:   gasPrice,
		nonceCache: cache,
	}, nil
}

// NonceCache returns the underlying shared cache. Used by callers that
// need to invalidate on error paths outside the Send hot loop (e.g. the
// reaper, which broadcasts a stuck row through Send but invalidates
// first so the resend uses the chain's current pending nonce rather
// than a possibly-strided cache value).
func (s *RPCBatcher) NonceCache() *noncecache.Cache {
	return s.nonceCache
}

// reserveNonces returns the starting nonce for a batch of `count` txs from
// `privateKey`'s address, advancing the shared nonce cache by `count`. See
// noncecache.Cache.Reserve for the invariants. Returns the starting nonce
// on success; the caller signs with `nonce, nonce+1, …, nonce+count-1`.
func (s *RPCBatcher) reserveNonces(
	ctx context.Context, privateKey *ecdsa.PrivateKey, count int,
) (uint64, error) {
	addr := crypto.PubkeyToAddress(privateKey.PublicKey)
	return s.nonceCache.Reserve(ctx, addr, count)
}

// invalidateNonceCache forces the next Send for `addr` to re-query the
// chain. Called on error paths that may have left the local nonce ahead
// of what was actually broadcast (creating a gap on chain).
func (s *RPCBatcher) invalidateNonceCache(addr common.Address) {
	s.nonceCache.Invalidate(addr)
}

// Send signs each input and broadcasts the batch in a single RPC call.
// Returns one SendResult per input, in input order.
//
// A function-level error is a whole-batch retryable failure: caller
// should retry the same inputs. Per-tx outcomes live in the slice.
func (s *RPCBatcher) Send(ctx context.Context, inputs []TransactionInput) ([]SendResult, error) {
	results := make([]SendResult, len(inputs))
	if len(inputs) == 0 {
		return results, nil
	}

	privateKey, err := s.keyQueue.Dequeue(ctx)
	if err != nil {
		return nil, fmt.Errorf("dequeue key: %w", err)
	}
	defer s.keyQueue.Enqueue(privateKey)
	addr := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Reserve a contiguous block of nonces from the local cache (seeded
	// from chain on first use). This prevents two back-to-back Sends
	// that grab the same key from both querying the chain and getting
	// the same nonce — see the docstring on the shared nonce cache for
	// the full race.
	nonce, err := s.reserveNonces(ctx, privateKey, len(inputs))
	if err != nil {
		return nil, fmt.Errorf("reserving nonces: %w", err)
	}

	requests, reqToInput := s.signBatch(ctx, inputs, privateKey, nonce, results)
	if len(requests) < len(inputs) {
		// At least one input was dropped at the sign step (pre-flight
		// revert or sign error). signBatch reuses the dropped nonce for
		// the next successful input, so the broadcast set lives in nonces
		// [nonce, nonce+len(requests)), but the cache was advanced past
		// the END of the reservation block (nonce+len(inputs)). That's
		// a gap of (len(inputs) - len(requests)) nonces that will never
		// be broadcast, stalling future mining for this key.
		//
		// The dropped txs were never signed, so nothing was broadcast for
		// them and their nonces are genuinely free. Invalidate so the next
		// Send re-queries the chain and reclaims the gap. (This is the
		// UNAMBIGUOUS case: the drop happens before callRPCBatch below.)
		//
		// (The len(requests) == 0 case is just the all-failed extreme;
		// the inequality covers both partial and total failures.)
		s.invalidateNonceCache(addr)
		if len(requests) == 0 {
			return results, nil
		}
	}

	if err := s.callRPCBatch(ctx, requests); err != nil {
		// Whole-batch broadcast failure. On a single authoritative RPC
		// backend (the assumed topology) this means nothing was accepted, so
		// invalidating lets the next Send reclaim the reserved nonces.
		//
		// CAVEAT (multi-node / load-balanced backends only): an error that
		// surfaces AFTER an upstream node already accepted the batch, e.g. a
		// 5xx from a proxy/LB post-accept, is indistinguishable here from a
		// true reject. The txs are then live, but we invalidate and the next
		// Reserve re-seeds from a pending count a lagging replica may not yet
		// reflect, which can re-open the cross-Send collision the cache exists
		// to prevent. Accepted as a documented residual (follow-up: pin a
		// sticky backend per chain or make the reseed monotonic); it needs a
		// replicated RPC backend to manifest and is backstopped by the reaper
		// re-broadcast + the on-chain executed[messageId] guard.
		s.invalidateNonceCache(addr)
		return nil, fmt.Errorf("broadcasting raw tx batch: %w", err)
	}

	needsResync := false
	for reqIdx, req := range requests {
		inIdx := reqToInput[reqIdx]
		if req.Error != nil {
			// eth_sendRawTransaction never executes the tx, so per-req
			// errors here are transport / format / sig / nonce — not reverts.
			results[inIdx].Err = fmt.Errorf("broadcast tx %s: %w", inputs[inIdx].ID, req.Error)
			// "already known" means the node already holds THIS exact tx
			// (same hash) in its mempool, e.g. an idempotent re-broadcast of
			// identical calldata at the same nonce. The slot is occupied by
			// OUR OWN tx, so the chain's pending count already reflects it and
			// the cache is correctly positioned; re-querying could only regress
			// it below a nonce we've already handed out (the residual
			// multi-node race) without reclaiming anything. Skip the resync for
			// that case only. Every OTHER per-tx error still warrants a resync:
			// "nonce too low" → chain is AHEAD and we must realign UP;
			// format/sig/transport rejects → the slot may be free to reclaim.
			if !isAlreadyKnownErr(req.Error) {
				needsResync = true
			}
			continue
		}
		results[inIdx].Hash = common.HexToHash(*req.Result.(*string))
	}

	if needsResync {
		// At least one tx wasn't accepted into the slot we expected.
		// Conservative: resync from chain on the next call to avoid striding
		// past a freed nonce slot or sitting below a chain that moved ahead.
		s.invalidateNonceCache(addr)
	}
	return results, nil
}

// signBatch signs each input and assembles the RPC requests. Signing
// outcomes are written directly into results; successful signs consume
// a nonce and get queued for broadcast.
func (s *RPCBatcher) signBatch(
	ctx context.Context,
	inputs []TransactionInput,
	privateKey *ecdsa.PrivateKey,
	nonce uint64,
	results []SendResult,
) ([]rpc.BatchElem, []int) {
	requests := make([]rpc.BatchElem, 0, len(inputs))
	reqToInput := make([]int, 0, len(inputs))

	for i, in := range inputs {
		rlpBytes, revert, signErr := s.signAndEncodeTransaction(ctx, in, privateKey, nonce)
		switch {
		case signErr != nil:
			results[i].Err = fmt.Errorf("sign tx %s: %w", in.ID, signErr)
			continue
		case revert != nil:
			results[i].Revert = revert
			continue
		}
		nonce++

		requests = append(requests, rpc.BatchElem{
			Method: "eth_sendRawTransaction",
			Args:   []interface{}{"0x" + hex.EncodeToString(rlpBytes)},
			Result: new(string),
		})
		reqToInput = append(reqToInput, i)
	}
	return requests, reqToInput
}

// GetReceipts polls receipts for the given hashes. Returns one result
// per hash, in input order. Not-yet-mined txs come back with Pending=true;
// mined-but-reverted txs come back with Revert carrying raw bytes
// extracted by replaying the tx at the mined block (same mechanism as
// cts/service/txops.go).
//
// A function-level error is a whole-batch retryable failure.
func (s *RPCBatcher) GetReceipts(ctx context.Context, hashes []common.Hash) ([]ReceiptResult, error) {
	results := make([]ReceiptResult, len(hashes))
	if len(hashes) == 0 {
		return results, nil
	}

	requests := make([]rpc.BatchElem, len(hashes))
	for i, h := range hashes {
		requests[i] = rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{h},
			Result: new(types.Receipt),
		}
	}

	if err := s.callRPCBatch(ctx, requests); err != nil {
		return nil, fmt.Errorf("fetching receipts batch: %w", err)
	}

	for i, req := range requests {
		if req.Error != nil {
			// Some nodes return a partially-populated receipt for a tx
			// that has just been included but not finalized — missing
			// required consensus fields like cumulativeGasUsed. Treat
			// it as pending; the next tick will see a complete receipt.
			if isPartialReceiptErr(req.Error) {
				results[i].Pending = true
				continue
			}
			results[i].Err = fmt.Errorf("fetch receipt %s: %w", hashes[i].Hex(), req.Error)
			continue
		}
		receipt := req.Result.(*types.Receipt)
		// A null JSON response unmarshals into a zero-valued receipt.
		if receipt.TxHash == (common.Hash{}) {
			results[i].Pending = true
			continue
		}
		if receipt.Status == types.ReceiptStatusSuccessful {
			results[i].Receipt = receipt
			continue
		}
		revert, err := s.recoverRevertReason(ctx, hashes[i], receipt.BlockNumber)
		if err != nil {
			results[i].Err = fmt.Errorf("recover revert for %s: %w", hashes[i].Hex(), err)
			continue
		}
		results[i].Revert = revert
	}
	return results, nil
}

// isPartialReceiptErr reports whether err is go-ethereum's
// "missing required field '<x>' for Receipt" — emitted by
// types.Receipt.UnmarshalJSON when a node returns a receipt that
// has been included but is not yet fully populated. The string is
// stable across go-ethereum versions (generated by gencodec) and
// there is no exported sentinel to compare against.
func isPartialReceiptErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "missing required field") && strings.Contains(msg, "for Receipt")
}

// callRPCBatch wraps BatchCallContext with the batch retry policy and
// a per-attempt deadline. Returned error is a whole-batch infra failure.
func (s *RPCBatcher) callRPCBatch(ctx context.Context, reqs []rpc.BatchElem) error {
	retryCtx, cancelRetry := context.WithTimeout(ctx, batchRetryTimeout)
	defer cancelRetry()

	return ethretry.WithRetry(retryCtx, func() error {
		opCtx, cancelOp := context.WithTimeout(retryCtx, rpcCallTimeout)
		defer cancelOp()
		return s.rpcClient.BatchCallContext(opCtx, reqs)
	})
}

// signAndEncodeTransaction returns exactly one of:
//
//	(rlp, nil, nil)    — signed; caller queues for broadcast
//	(nil, revert, nil) — gas estimation caught a pre-flight revert
//	(nil, nil, err)    — retryable infra failure
func (s *RPCBatcher) signAndEncodeTransaction(
	ctx context.Context,
	input TransactionInput,
	privateKey *ecdsa.PrivateKey,
	nonce uint64,
) ([]byte, *Revert, error) {
	estimateCtx, cancelEstimate := context.WithTimeout(ctx, gasEstimateTimeout)
	defer cancelEstimate()
	estimatedGas, err := s.ethClient.EstimateGas(estimateCtx, ethereum.CallMsg{
		From:     crypto.PubkeyToAddress(privateKey.PublicKey),
		To:       &input.Address,
		GasPrice: s.gasPrice,
		Data:     input.Data,
	})
	if err != nil {
		// Gas-estimation reverts arrive as rpc.DataError carrying the
		// raw revert blob. Surface as data, not as infra error.
		if revertBytes, ok := extractRevertData(err); ok {
			return nil, &Revert{Data: revertBytes}, nil
		}
		return nil, nil, withstack.Wrap(fmt.Errorf("estimating gas: %w", err))
	}

	// 30% gas buffer to absorb state drift between estimation and mining.
	// Without this, txs that touch storage slots whose cost depends on state
	// (e.g. SSTORE 5k vs 20k for fresh-zero slot) can OOG when concurrent
	// txs change the underlying state. The atomic-unlock dispatch hit this
	// reliably when a user race-teleported their balance to zero between
	// the unlock-signature estimation and broadcast — `_balances[user]`
	// went from non-zero to zero, and the unlock's SSTORE to refill it
	// cost an extra ~15k gas the original estimate didn't carry.
	gasLimit := (estimatedGas * 130) / 100

	ethTx := types.NewTransaction(nonce, input.Address, big.NewInt(0), gasLimit, s.gasPrice, input.Data)
	signedTx, err := types.SignTx(ethTx, types.NewEIP155Signer(s.chainID), privateKey)
	if err != nil {
		return nil, nil, withstack.Wrap(fmt.Errorf("signing tx: %w", err))
	}

	rlpBytes, err := signedTx.MarshalBinary()
	if err != nil {
		return nil, nil, fmt.Errorf("rlp-encoding signed tx: %w", err)
	}
	// Per-tx signing trace. Useful when correlating an on-chain hash
	// back to which input produced it (e.g., comparing an original
	// send vs a reaper resend for the same correlation_id). Kept at
	// Debug because it fires for every signed tx — a 60-message batch
	// would otherwise produce 60+ Info lines and drown operational
	// signal in production.
	slog.Debug("signed_tx",
		slog.String("event", "signed_tx"),
		slog.String("correlation_id", input.ID),
		slog.String("sender", crypto.PubkeyToAddress(privateKey.PublicKey).Hex()),
		slog.Uint64("nonce", nonce),
		slog.Uint64("estimated_gas", estimatedGas),
		slog.Uint64("gas_limit", gasLimit),
		slog.String("gas_price", s.gasPrice.String()),
		slog.String("hash", signedTx.Hash().Hex()),
	)
	return rlpBytes, nil, nil
}

// recoverRevertReason replays a mined failed tx as an eth_call at its
// mined block to extract the raw revert bytes. Mirrors
// cts/service/txops.go:recoverRevertReason.
func (s *RPCBatcher) recoverRevertReason(
	ctx context.Context,
	hash common.Hash,
	blockNumber *big.Int,
) (*Revert, error) {
	tx, _, err := s.ethClient.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("fetching tx by hash: %w", err))
	}
	from, err := types.Sender(types.LatestSignerForChainID(s.chainID), tx)
	if err != nil {
		return nil, fmt.Errorf("recovering sender: %w", err)
	}
	_, callErr := s.ethClient.CallContract(ctx, ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}, blockNumber)
	if callErr == nil {
		// State drift between mine-time and simulate-time. The tx
		// reverted on-chain but replay succeeded. Classify as revert
		// with empty bytes so the caller still sees the final state.
		return &Revert{}, nil
	}
	if revertBytes, ok := extractRevertData(callErr); ok {
		return &Revert{Data: revertBytes}, nil
	}
	return nil, withstack.Wrap(fmt.Errorf("replay call: %w", callErr))
}

// extractRevertData pulls raw revert bytes from a go-ethereum RPC error.
// Mirrors cts/service/txops.go:extractRevertData.
func extractRevertData(err error) ([]byte, bool) {
	var de rpc.DataError
	if !errors.As(err, &de) {
		return nil, false
	}
	raw, ok := de.ErrorData().(string)
	if !ok {
		return nil, false
	}
	data, decodeErr := hex.DecodeString(strings.TrimPrefix(raw, "0x"))
	if decodeErr != nil {
		return nil, false
	}
	return data, true
}

// isAlreadyKnownErr reports whether a per-tx broadcast error means the node
// already holds THIS exact transaction (same hash) in its mempool: an idempotent
// re-broadcast of identical signed bytes. There is no exported sentinel, so match
// the stable wire strings (case-insensitive):
//
//   - "already known": go-ethereum, and also axyl/Reth. Reth maps its internal
//     PoolError::AlreadyImported to the geth-compatible RpcPoolError::AlreadyKnown,
//     whose JSON-RPC message is "already known" (verified against Reth v1.11.0,
//     the axyl runtime). So the wire string is the same on both backends.
//   - "already imported": Reth's raw pool Display string, matched defensively in
//     case a path surfaces it before the RPC remap.
//
// Distinct from "replacement transaction underpriced" (a DIFFERENT tx contends for
// the nonce) and "nonce too low" (the chain has moved past it); both of those still
// warrant a resync, so they are intentionally NOT matched here.
func isAlreadyKnownErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "already known") || strings.Contains(msg, "already imported")
}
