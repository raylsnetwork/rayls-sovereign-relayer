package contractclient_test

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthGen(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.Nil(t, err, "got error when generating key, cannot contiue with tests")

	var (
		chainID  = new(big.Int).SetUint64(1337)
		gasPrice = new(big.Int).SetUint64(1000)
	)

	t.Run("gets gas price from client", func(t *testing.T) {
		ethClient := &AuthEthereumClientMock{
			ChainIDFunc: func(ctx context.Context) (*big.Int, error) {
				return chainID, nil
			},
			SuggestGasPriceFunc: func(ctx context.Context) (*big.Int, error) {
				return gasPrice, nil
			},
			PendingNonceAtFunc: func(ctx context.Context, addr common.Address) (uint64, error) {
				return 0, nil
			},
		}
		gen, err := contractclient.NewAuthGen(context.Background(), ethClient)
		require.Nil(t, err)

		auth, err := gen.Get(context.Background(), privateKey)
		require.Nil(t, err)

		assert.Zero(t, gasPrice.Cmp(auth.GasPrice))
	})

	t.Run("sets value = 0", func(t *testing.T) {
		wantValue := new(big.Int).SetInt64(0)

		ethClient := &AuthEthereumClientMock{
			ChainIDFunc: func(ctx context.Context) (*big.Int, error) {
				return chainID, nil
			},
			SuggestGasPriceFunc: func(ctx context.Context) (*big.Int, error) {
				return gasPrice, nil
			},
			PendingNonceAtFunc: func(ctx context.Context, addr common.Address) (uint64, error) {
				return 0, nil
			},
		}
		gen, err := contractclient.NewAuthGen(context.Background(), ethClient)
		require.Nil(t, err)

		auth, err := gen.Get(context.Background(), privateKey)
		require.Nil(t, err)

		assert.Zero(t, wantValue.Cmp(auth.Value))
	})

	t.Run("sets gas limit = 0", func(t *testing.T) {
		wantGasLimit := uint64(0)

		ethClient := &AuthEthereumClientMock{
			ChainIDFunc: func(ctx context.Context) (*big.Int, error) {
				return chainID, nil
			},
			SuggestGasPriceFunc: func(ctx context.Context) (*big.Int, error) {
				return gasPrice, nil
			},
			PendingNonceAtFunc: func(ctx context.Context, addr common.Address) (uint64, error) {
				return 0, nil
			},
		}
		gen, err := contractclient.NewAuthGen(context.Background(), ethClient)
		require.Nil(t, err)

		auth, err := gen.Get(context.Background(), privateKey)
		require.Nil(t, err)

		assert.Equal(t, wantGasLimit, auth.GasLimit)
	})

	t.Run("calls NewKeyedTransactorWithChainID with private key and chain ID", func(t *testing.T) {
		// Get derives the sender address from the private key to seed
		// the per-key nonce cache. The empty ecdsa.PrivateKey{} this
		// test used to pass crashes on `crypto.PubkeyToAddress` (X/Y
		// nil) — use a real generated
		// key instead.
		wantPrivateKey, err := crypto.GenerateKey()
		require.Nil(t, err)
		wantChainID := new(big.Int).SetUint64(1337)

		var (
			gotPrivateKey *ecdsa.PrivateKey
			gotChainID    *big.Int
		)

		spyFunc := func(key *ecdsa.PrivateKey, chainID *big.Int) (*bind.TransactOpts, error) {
			gotPrivateKey = key
			gotChainID = chainID

			return &bind.TransactOpts{}, nil
		}

		ethClient := &AuthEthereumClientMock{
			ChainIDFunc: func(ctx context.Context) (*big.Int, error) {
				return wantChainID, nil
			},
			SuggestGasPriceFunc: func(ctx context.Context) (*big.Int, error) {
				return gasPrice, nil
			},
			PendingNonceAtFunc: func(ctx context.Context, addr common.Address) (uint64, error) {
				return 0, nil
			},
		}
		gen, err := contractclient.NewAuthGenWithCustomKeyedTransactorFunc(context.Background(), ethClient, spyFunc)
		require.Nil(t, err)

		_, err = gen.Get(context.Background(), wantPrivateKey)
		require.Nil(t, err)

		assert.Equal(t, wantPrivateKey, gotPrivateKey)
		assert.Equal(t, wantChainID, gotChainID)
	})

	// The previous version of this test asserted that `auth.Nonce` is
	// left nil so go-ethereum's bind package would resolve it via
	// `eth_getTransactionCount` at submission time. That was the exact
	// race the current behaviour closes — two back-to-back Gets in
	// quick succession would both hit the chain, both see the same
	// pre-mine nonce, and sign two different txs at the same
	// (sender, nonce) slot. The chain mines one and silently drops the
	// other. AuthGen now pre-resolves `auth.Nonce` from a local cache
	// (seeded from PendingNonceAt on first use), so two consecutive
	// Gets get strictly sequential nonces.
	t.Run("pre-resolves auth.Nonce from the local cache", func(t *testing.T) {
		ethClient := &AuthEthereumClientMock{
			ChainIDFunc: func(ctx context.Context) (*big.Int, error) {
				return chainID, nil
			},
			SuggestGasPriceFunc: func(ctx context.Context) (*big.Int, error) {
				return gasPrice, nil
			},
			PendingNonceAtFunc: func(ctx context.Context, addr common.Address) (uint64, error) {
				return 42, nil
			},
		}
		fakeFunc := func(key *ecdsa.PrivateKey, chainID *big.Int) (*bind.TransactOpts, error) {
			return &bind.TransactOpts{}, nil
		}
		gen, err := contractclient.NewAuthGenWithCustomKeyedTransactorFunc(context.Background(), ethClient, fakeFunc)
		require.Nil(t, err)

		auth, err := gen.Get(context.Background(), privateKey)
		require.Nil(t, err)

		require.NotNil(t, auth.Nonce, "AuthGen must pre-resolve nonce — leaving it nil reintroduces the race")
		assert.Equal(t, uint64(42), auth.Nonce.Uint64(), "first Get should seed from chain (42)")
	})

	t.Run("wraps errors from new transactor function in AuthGenError", func(t *testing.T) {
		wantError := errors.New("example error")

		ethClient := &AuthEthereumClientMock{
			ChainIDFunc: func(ctx context.Context) (*big.Int, error) {
				return chainID, nil
			},
			SuggestGasPriceFunc: func(ctx context.Context) (*big.Int, error) {
				return gasPrice, nil
			},
			PendingNonceAtFunc: func(ctx context.Context, addr common.Address) (uint64, error) {
				return 0, nil
			},
		}
		errFunc := func(key *ecdsa.PrivateKey, chainID *big.Int) (*bind.TransactOpts, error) {
			return &bind.TransactOpts{}, wantError
		}
		gen, err := contractclient.NewAuthGenWithCustomKeyedTransactorFunc(context.Background(), ethClient, errFunc)
		require.Nil(t, err)

		_, err = gen.Get(context.Background(), &ecdsa.PrivateKey{})

		authGenErrType := &contractclient.AuthGenError{}
		assert.ErrorAs(t, err, &authGenErrType)

		assert.ErrorIs(t, err, wantError)
	})
}
