package contractclient

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/ethretry"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/noncecache"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

// authInitCallTimeout bounds each ChainID/SuggestGasPrice attempt at AuthGen
// construction so a stalled (VPN-routed) RPC is aborted and retried rather than
// blocking until the caller's whole construction budget runs out. Mirrors
// batcherInitCallTimeout in cts/ethrpc.
const authInitCallTimeout = 10 * time.Second

//go:generate moq --pkg contractclient_test -out auth_mock_test.go . AuthEthereumClient
type AuthEthereumClient interface {
	ChainID(context.Context) (*big.Int, error)
	SuggestGasPrice(context.Context) (*big.Int, error)
	// PendingNonceAt is consulted to seed the per-key local nonce cache
	// on first use. See the docstring on AuthGen.nonceCache for the
	// race this prevents.
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
}

type newKeyedTransactorWithChainIDFunc func(key *ecdsa.PrivateKey, chainID *big.Int) (*bind.TransactOpts, error)

type AuthGen struct {
	chainID  *big.Int
	gasPrice *big.Int

	ethClient                     AuthEthereumClient
	newKeyedTransactorWithChainID newKeyedTransactorWithChainIDFunc

	// Per-address local nonce cache.
	//
	// `RPCBatcher.Send` and `AuthGen.Get` are two separate code paths
	// that ultimately sign for the same keys: the former is used by
	// `cts/batcher/sender.go` for the destination-side dispatch
	// (privatenode identity); the latter by `cts/service/txops.go`'s
	// SignAndSend / Deploy / BatchSignAndSend gRPC handlers, which build
	// a `bind.TransactOpts` and then `bound.RawTransact(auth, data)`.
	// The go-ethereum bind package resolves `auth.Nonce` via
	// `eth_getTransactionCount` at submission time when `auth.Nonce` is
	// nil — and two concurrent Sends that dequeue the same key from the
	// keyQueue can both observe the same pending nonce while the prior
	// tx is still in-flight, sign different calldata under it, and
	// collide on chain.
	//
	// The cache pre-resolves `auth.Nonce` from a per-address counter.
	// Each `Get` reserves exactly one nonce (the bind path is always
	// 1 tx at a time). Error paths from the callers invalidate via
	// `InvalidateNonceCache(addr)`.
	//
	// CRUCIALLY, the cache is held through a pointer so production
	// wiring can SHARE a single instance between AuthGen (sync gRPC) and
	// the RPCBatcher (async pipeline) signing for the same identity.
	// Without sharing, the per-component race is closed but a cross-
	// component race remains open: AuthGen's SignAndSend and the async
	// pipeline's RPCBatcher.Send can both read PendingNonceAt for the
	// same key in the window where a prior tx is broadcast but not yet
	// mined, get the same value, and broadcast colliding txs.
	//
	// The legacy `NewAuthGen` constructor still creates a per-AuthGen
	// cache for callers (e.g. keymanager) that don't share a signing
	// queue with anything else. The shared-cache wiring uses
	// `NewAuthGenWithCache`.
	nonceCache *noncecache.Cache
}

// NewAuthGen builds an AuthGen, fetching the chain ID and a baseline gas price
// from the supplied client. The supplied ctx scopes those bootstrap RPC calls;
// pass a context with a sane deadline at startup so a hung RPC doesn't wedge
// relayer launch.
//
// The returned AuthGen owns its own nonce cache. For wiring that needs to
// share a nonce counter with other components signing for the same identity
// (e.g. an RPCBatcher), use NewAuthGenWithCache.
func NewAuthGen(ctx context.Context, ethClient AuthEthereumClient) (*AuthGen, error) {
	return NewAuthGenWithCache(ctx, ethClient, noncecache.New(ethClient))
}

// NewAuthGenWithCustomKeyedTransactorFunc is the test-friendly constructor;
// behaves like NewAuthGen but lets callers swap the transactor factory.
func NewAuthGenWithCustomKeyedTransactorFunc(
	ctx context.Context,
	ethClient AuthEthereumClient,
	newKeyedTransactorWithChainID newKeyedTransactorWithChainIDFunc,
) (*AuthGen, error) {
	a, err := NewAuthGenWithCache(ctx, ethClient, noncecache.New(ethClient))
	if err != nil {
		return nil, err
	}
	a.newKeyedTransactorWithChainID = newKeyedTransactorWithChainID
	return a, nil
}

// NewAuthGenWithCache is like NewAuthGen but accepts an externally-owned
// nonce cache. Used by the production wiring to share a cache between
// AuthGen and one or more RPCBatchers signing for the same identity.
func NewAuthGenWithCache(
	ctx context.Context,
	ethClient AuthEthereumClient,
	cache *noncecache.Cache,
) (*AuthGen, error) {
	if cache == nil {
		return nil, fmt.Errorf("nonce cache must not be nil")
	}
	// chainID and gasPrice are fetched once here, over the (VPN-routed in remote
	// dev) RPC link. This runs before the batcher's own retry in
	// buildTxOpsService, so without retrying here a transient blip would crash
	// CTS startup regardless. Bound each attempt with its own timeout so a
	// stalled call is aborted and retried — mirrors NewBatcherWithCache.
	var chainID *big.Int
	if err := ethretry.WithRetry(ctx, func() error {
		callCtx, cancel := context.WithTimeout(ctx, authInitCallTimeout)
		defer cancel()
		var err error
		chainID, err = ethClient.ChainID(callCtx)
		return err
	}); err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to get chain ID: %w", err))
	}

	var gasPrice *big.Int
	if err := ethretry.WithRetry(ctx, func() error {
		callCtx, cancel := context.WithTimeout(ctx, authInitCallTimeout)
		defer cancel()
		var err error
		gasPrice, err = ethClient.SuggestGasPrice(callCtx)
		return err
	}); err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to get gas price: %w", err))
	}

	return &AuthGen{
		chainID:  chainID,
		gasPrice: gasPrice,

		ethClient:                     ethClient,
		newKeyedTransactorWithChainID: bind.NewKeyedTransactorWithChainID,
		nonceCache:                    cache,
	}, nil
}

// reserveNonce returns the next-to-use nonce for `addr`, advancing the shared
// nonce cache by 1. Delegates to noncecache.Cache.Reserve for the seed +
// reserve logic (see that package for the invariants).
func (g *AuthGen) reserveNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return g.nonceCache.Reserve(ctx, addr, 1)
}

// InvalidateNonceCache forces the next `Get` for `addr` to re-query the chain.
// Callers MUST call this on any error path that prevents the tx from being
// broadcast / mined — otherwise the local cache will stride past nonces that
// were never actually consumed, creating gaps that stall mining.
func (g *AuthGen) InvalidateNonceCache(addr common.Address) {
	g.nonceCache.Invalidate(addr)
}

// Get builds a *bind.TransactOpts for a single contract call. The supplied ctx
// is stored on auth.Context so abigen-generated methods (which call
// eth_getTransactionCount, eth_estimateGas, eth_sendRawTransaction, etc.
// internally) honor caller cancellation and deadlines instead of falling
// through to context.Background() inside go-ethereum's bind.ensureContext.
//
// Without this propagation, any RPC call inside an abigen method paired with
// the relayer's RetryTransport (which retries forever on transient errors)
// can hang the goroutine indefinitely.
func (g *AuthGen) Get(ctx context.Context, key *ecdsa.PrivateKey) (*bind.TransactOpts, error) {
	auth, err := g.newKeyedTransactorWithChainID(key, g.chainID)
	if err != nil {
		return nil, WrapInAuthGenError("failed to create keyed transactor", withstack.Wrap(err))
	}

	// Pre-resolve `auth.Nonce` from the shared local cache so the
	// go-ethereum bind package does NOT call `PendingNonceAt` at
	// submission time.
	//
	// The race: this AuthGen path and the async RPCBatcher pipeline both
	// sign for the same key (shared keyQueue) and both used to resolve the
	// nonce by an unsynchronised `PendingNonceAt` read. Two such reads that
	// interleave before either tx is broadcast return the SAME value. Note
	// `PendingNonceAt` is mempool-aware (eth_getTransactionCount "pending"),
	// so the collision window is read-to-broadcast, NOT broadcast-to-mining.
	// Both paths then sign DIFFERENT calldata at that nonce. On a chain that
	// permits 0-fee same-nonce replacement (axyl/Reth; Geth/Besu would
	// reject it as "replacement underpriced"), the second tx evicts the
	// first from the mempool, so one is silently dropped and its hash never
	// mines. Reserving from the shared cache closes the read-modify-write
	// window so the two paths can't hand out the same nonce.
	addr := crypto.PubkeyToAddress(key.PublicKey)
	nonce, err := g.reserveNonce(ctx, addr)
	if err != nil {
		return nil, WrapInAuthGenError("reserving nonce", err)
	}
	auth.Nonce = new(big.Int).SetUint64(nonce)

	auth.Value = new(big.Int).SetUint64(0)

	auth.GasLimit = 0
	auth.GasPrice = g.gasPrice

	auth.Context = ctx

	return auth, nil
}
