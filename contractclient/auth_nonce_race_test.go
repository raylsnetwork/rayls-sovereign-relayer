// Regression tests for the AuthGen nonce race.
//
// Same bug as the RPCBatcher.Send nonce race — but the AuthGen path
// (used by cts/service/txops.go's SignAndSend / Deploy /
// BatchSignAndSend) is a different code path. These tests exercise
// AuthGen's per-key local nonce cache directly.
//
//   - First Get for an address seeds the cache from PendingNonceAt.
//   - Subsequent Gets advance the cache WITHOUT re-querying the chain.
//   - Multiple addresses keep independent caches.
//   - InvalidateNonceCache forces the next Get to re-query.
//
// All tests below FAIL against the pre-fix AuthGen (which left
// auth.Nonce nil so go-ethereum's bind package would resolve nonce at
// submission time, re-introducing the race).

package contractclient_test

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAuthPK0 = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	testAuthPK1 = "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"
)

func mustAuthPK(t *testing.T, hex string) *ecdsa.PrivateKey {
	t.Helper()
	pk, err := crypto.HexToECDSA(hex)
	require.NoError(t, err)
	return pk
}

func newAuthGenWithMocks(t *testing.T) (*contractclient.AuthGen, *AuthEthereumClientMock) {
	t.Helper()
	mock := &AuthEthereumClientMock{
		ChainIDFunc: func(ctx context.Context) (*big.Int, error) {
			return big.NewInt(1337), nil
		},
		SuggestGasPriceFunc: func(ctx context.Context) (*big.Int, error) {
			return big.NewInt(0), nil
		},
		PendingNonceAtFunc: func(ctx context.Context, addr common.Address) (uint64, error) {
			return 100, nil // default seed
		},
	}
	// Use the custom transactor func so we don't depend on the real
	// bind.NewKeyedTransactorWithChainID (which requires a real chain ID
	// match for proper signer setup — irrelevant for nonce-cache tests).
	fake := func(key *ecdsa.PrivateKey, _ *big.Int) (*bind.TransactOpts, error) {
		return &bind.TransactOpts{
			From: crypto.PubkeyToAddress(key.PublicKey),
		}, nil
	}
	gen, err := contractclient.NewAuthGenWithCustomKeyedTransactorFunc(context.Background(), mock, fake)
	require.NoError(t, err)
	return gen, mock
}

// TestAuthGen_NoNonceReuseAcrossGetsOnSameKey — the core regression for
// the AuthGen path. Two back-to-back Gets with the same key must
// produce strictly sequential nonces, even when PendingNonceAt is
// stubbed to return the same value both times (mirrors the "first tx
// not yet mined" condition on the live chain).
func TestAuthGen_NoNonceReuseAcrossGetsOnSameKey(t *testing.T) {
	gen, mock := newAuthGenWithMocks(t)
	mock.PendingNonceAtFunc = func(ctx context.Context, addr common.Address) (uint64, error) {
		return 100, nil
	}
	key := mustAuthPK(t, testAuthPK0)

	a1, err := gen.Get(context.Background(), key)
	require.NoError(t, err)
	a2, err := gen.Get(context.Background(), key)
	require.NoError(t, err)

	require.NotNil(t, a1.Nonce)
	require.NotNil(t, a2.Nonce)
	t.Logf("first Get → nonce=%d, second Get → nonce=%d (PendingNonceAt stub returns 100)",
		a1.Nonce.Uint64(), a2.Nonce.Uint64())
	assert.NotEqualf(t, a1.Nonce.Uint64(), a2.Nonce.Uint64(),
		"regression: two back-to-back Gets with the same key must NOT "+
			"return the same auth.Nonce. Pre-fix this test fails because auth.Nonce was "+
			"left nil and the race resolved at tx-submission time, not in AuthGen.")
	assert.Equal(t, uint64(100), a1.Nonce.Uint64(), "first Get seeds from chain (100)")
	assert.Equal(t, uint64(101), a2.Nonce.Uint64(), "second Get advances cache without re-query")
}

// TestAuthGen_FirstUseSeedsFromChain — exactly one PendingNonceAt call
// across N Gets for the same key.
func TestAuthGen_FirstUseSeedsFromChain(t *testing.T) {
	gen, mock := newAuthGenWithMocks(t)
	mock.PendingNonceAtFunc = func(ctx context.Context, addr common.Address) (uint64, error) {
		return 7, nil
	}
	key := mustAuthPK(t, testAuthPK0)

	for i := 0; i < 5; i++ {
		_, err := gen.Get(context.Background(), key)
		require.NoError(t, err, "Get %d", i)
	}
	calls := len(mock.PendingNonceAtCalls())
	t.Logf("PendingNonceAt called %d time(s) for 5 Gets — expect 1", calls)
	assert.Equalf(t, 1, calls,
		"PendingNonceAt should be called exactly once for the same key — "+
			"first use seeds the cache, subsequent uses don't re-query. "+
			"Pre-fix this would be %d (every Get queries via go-ethereum's bind).", 5)
}

// TestAuthGen_DifferentKeysIndependent — each address gets its own cache.
func TestAuthGen_DifferentKeysIndependent(t *testing.T) {
	keyA := mustAuthPK(t, testAuthPK0)
	keyB := mustAuthPK(t, testAuthPK1)
	addrA := crypto.PubkeyToAddress(keyA.PublicKey)
	addrB := crypto.PubkeyToAddress(keyB.PublicKey)

	gen, mock := newAuthGenWithMocks(t)
	mock.PendingNonceAtFunc = func(ctx context.Context, addr common.Address) (uint64, error) {
		switch addr {
		case addrA:
			return 10, nil
		case addrB:
			return 200, nil
		default:
			t.Fatalf("unexpected address %s", addr.Hex())
			return 0, nil
		}
	}

	// Alternate keys: A, B, A, B → expect nonces 10, 200, 11, 201.
	got := make([]uint64, 0, 4)
	for _, k := range []*ecdsa.PrivateKey{keyA, keyB, keyA, keyB} {
		a, err := gen.Get(context.Background(), k)
		require.NoError(t, err)
		got = append(got, a.Nonce.Uint64())
	}
	want := []uint64{10, 200, 11, 201}
	t.Logf("observed: %v  expected: %v", got, want)
	assert.Equal(t, want, got, "two addresses must maintain independent caches")
	assert.Equal(t, 2, len(mock.PendingNonceAtCalls()),
		"PendingNonceAt should be called twice — once per distinct address")
}

// TestAuthGen_InvalidateNonceCache — after invalidation, the next Get
// re-queries the chain. Mirrors the recovery path that
// TxOpsService.SignAndSend / Deploy invoke on broadcast / wait-mined
// errors.
func TestAuthGen_InvalidateNonceCache(t *testing.T) {
	gen, mock := newAuthGenWithMocks(t)
	mock.PendingNonceAtFunc = func(ctx context.Context, addr common.Address) (uint64, error) {
		return 50, nil
	}
	key := mustAuthPK(t, testAuthPK0)
	addr := crypto.PubkeyToAddress(key.PublicKey)

	_, err := gen.Get(context.Background(), key)
	require.NoError(t, err)
	_, err = gen.Get(context.Background(), key)
	require.NoError(t, err)

	preInvalidate := len(mock.PendingNonceAtCalls())
	t.Logf("PendingNonceAt calls before invalidation: %d (expect 1)", preInvalidate)
	require.Equal(t, 1, preInvalidate)

	gen.InvalidateNonceCache(addr)
	a3, err := gen.Get(context.Background(), key)
	require.NoError(t, err)
	t.Logf("post-invalidation Get nonce: %d (expect 50 — re-read from chain)", a3.Nonce.Uint64())
	assert.Equal(t, uint64(50), a3.Nonce.Uint64(),
		"after invalidation, next Get must re-query the chain")
	assert.Equalf(t, 2, len(mock.PendingNonceAtCalls()),
		"PendingNonceAt should be re-called after InvalidateNonceCache. Saw %d total calls.",
		len(mock.PendingNonceAtCalls()))
}

// TestAuthGen_AdvancesAcrossManyGets — sanity that the cache is
// monotonically increasing across a sustained burst.
func TestAuthGen_AdvancesAcrossManyGets(t *testing.T) {
	gen, mock := newAuthGenWithMocks(t)
	mock.PendingNonceAtFunc = func(ctx context.Context, addr common.Address) (uint64, error) {
		return 1000, nil
	}
	key := mustAuthPK(t, testAuthPK0)

	got := make([]uint64, 0, 60) // simulate the "60 messages" burst load
	for i := 0; i < 60; i++ {
		a, err := gen.Get(context.Background(), key)
		require.NoError(t, err)
		got = append(got, a.Nonce.Uint64())
	}
	for i := 1; i < 60; i++ {
		require.Equalf(t, got[i-1]+1, got[i],
			"nonces must be strictly sequential — got[%d]=%d  got[%d]=%d",
			i-1, got[i-1], i, got[i])
	}
	t.Logf("60 Gets produced contiguous nonces 1000..%d (1 PendingNonceAt call)", got[len(got)-1])
	assert.Equal(t, 1, len(mock.PendingNonceAtCalls()), "exactly one chain query for 60 Gets")
}
