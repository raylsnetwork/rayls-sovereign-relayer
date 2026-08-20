package keymanager_test

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/domain"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/keymanager"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("happy_path_existing_keys", func(t *testing.T) {
		priv, pub := sampleBundles(t)
		keysSvc := &KeysServiceMock{
			GetPrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) { return priv, nil },
			GetPublicRelayerRaylsSignKeysFunc:  func(context.Context) (domain.PublicRelayerRaylsSignKeys, error) { return pub, nil },
		}

		waits := newWaitRecorder()
		mgr := newTestManager(keysSvc, stubFactoryAlwaysOK(waits))

		gotPriv, gotPub, err := mgr.Run(context.Background())
		require.NoError(t, err)
		require.Equal(t, priv, gotPriv)
		require.Equal(t, pub, gotPub)

		assert.Empty(t, keysSvc.CreatePrivateRelayerRaylsSignKeysCalls())
		assert.Empty(t, keysSvc.CreatePublicRelayerRaylsSignKeysCalls())

		require.Len(t, waits.records(), 4)
		assertWaitRecord(t, waits.records()[0], "private-hub", priv.PrivateHubKeys)
		assertWaitRecord(t, waits.records()[1], "privacy-node", priv.PrivateNodeKeys)
		assertWaitRecord(t, waits.records()[2], "private-chain", pub.PrivateChainKeys)
		assertWaitRecord(t, waits.records()[3], "public-chain", pub.PublicChainKeys)
	})

	t.Run("public_chain_disabled_skips_public_wait", func(t *testing.T) {
		priv, pub := sampleBundles(t)
		keysSvc := &KeysServiceMock{
			GetPrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) { return priv, nil },
			GetPublicRelayerRaylsSignKeysFunc:  func(context.Context) (domain.PublicRelayerRaylsSignKeys, error) { return pub, nil },
		}

		waits := newWaitRecorder()
		mgr := keymanager.NewSignKeyManager(
			keysSvc,
			keymanager.ChainConfig{},
			keymanager.ChainConfig{},
			keymanager.ChainConfig{},
			keymanager.WithBuildAccessManager(stubFactoryAlwaysOK(waits)),
			keymanager.WithPublicChain(false),
		)

		gotPriv, gotPub, err := mgr.Run(context.Background())
		require.NoError(t, err)
		require.Equal(t, priv, gotPriv)
		require.Equal(t, pub, gotPub)

		require.Len(t, waits.records(), 2)
		assertWaitRecord(t, waits.records()[0], "private-hub", priv.PrivateHubKeys)
		assertWaitRecord(t, waits.records()[1], "privacy-node", priv.PrivateNodeKeys)
	})

	t.Run("private_hub_disabled_skips_hub_wait", func(t *testing.T) {
		priv, pub := sampleBundles(t)
		keysSvc := &KeysServiceMock{
			GetPrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) { return priv, nil },
			GetPublicRelayerRaylsSignKeysFunc:  func(context.Context) (domain.PublicRelayerRaylsSignKeys, error) { return pub, nil },
		}

		waits := newWaitRecorder()
		mgr := keymanager.NewSignKeyManager(
			keysSvc,
			keymanager.ChainConfig{},
			keymanager.ChainConfig{},
			keymanager.ChainConfig{},
			keymanager.WithBuildAccessManager(stubFactoryAlwaysOK(waits)),
			keymanager.WithPrivateHub(false),
			keymanager.WithPublicChain(true),
		)

		gotPriv, gotPub, err := mgr.Run(context.Background())
		require.NoError(t, err)
		require.Equal(t, priv, gotPriv)
		require.Equal(t, pub, gotPub)

		// Hub wait (#1) is skipped; node + both public-chain waits remain.
		require.Len(t, waits.records(), 3)
		assertWaitRecord(t, waits.records()[0], "privacy-node", priv.PrivateNodeKeys)
		assertWaitRecord(t, waits.records()[1], "private-chain", pub.PrivateChainKeys)
		assertWaitRecord(t, waits.records()[2], "public-chain", pub.PublicChainKeys)
	})

	t.Run("private_hub_disabled_allows_empty_hub_keys", func(t *testing.T) {
		// With the hub disabled, an empty PrivateHubKeys keyset must not fail
		// validateRequiredKeysets — only the waited-on keysets are required.
		priv, pub := sampleBundles(t)
		priv.PrivateHubKeys = nil
		priv.PrivateHubDvpOperatorKeys = nil
		keysSvc := &KeysServiceMock{
			GetPrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) { return priv, nil },
			GetPublicRelayerRaylsSignKeysFunc:  func(context.Context) (domain.PublicRelayerRaylsSignKeys, error) { return pub, nil },
		}

		waits := newWaitRecorder()
		mgr := keymanager.NewSignKeyManager(
			keysSvc,
			keymanager.ChainConfig{},
			keymanager.ChainConfig{},
			keymanager.ChainConfig{},
			keymanager.WithBuildAccessManager(stubFactoryAlwaysOK(waits)),
			keymanager.WithPrivateHub(false),
			keymanager.WithPublicChain(true),
		)

		_, _, err := mgr.Run(context.Background())
		require.NoError(t, err)
		require.Len(t, waits.records(), 3)
		assertWaitRecord(t, waits.records()[0], "privacy-node", priv.PrivateNodeKeys)
	})

	t.Run("public_chain_disabled_allows_empty_public_keys", func(t *testing.T) {
		// With the public chain disabled, empty PrivateChainKeys/PublicChainKeys
		// keysets must not fail validateRequiredKeysets — only the waited-on
		// keysets are required. Mirrors the private-hub-disabled treatment.
		priv, pub := sampleBundles(t)
		pub.PrivateChainKeys = nil
		pub.PublicChainKeys = nil
		keysSvc := &KeysServiceMock{
			GetPrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) { return priv, nil },
			GetPublicRelayerRaylsSignKeysFunc:  func(context.Context) (domain.PublicRelayerRaylsSignKeys, error) { return pub, nil },
		}

		waits := newWaitRecorder()
		mgr := keymanager.NewSignKeyManager(
			keysSvc,
			keymanager.ChainConfig{},
			keymanager.ChainConfig{},
			keymanager.ChainConfig{},
			keymanager.WithBuildAccessManager(stubFactoryAlwaysOK(waits)),
			keymanager.WithPublicChain(false),
		)

		_, _, err := mgr.Run(context.Background())
		require.NoError(t, err)
		// Hub (default enabled) + privacy-node waits only; no public-chain waits.
		require.Len(t, waits.records(), 2)
		assertWaitRecord(t, waits.records()[0], "private-hub", priv.PrivateHubKeys)
		assertWaitRecord(t, waits.records()[1], "privacy-node", priv.PrivateNodeKeys)
	})

	t.Run("first_boot_both", func(t *testing.T) {
		priv, pub := sampleBundles(t)
		keysSvc := &KeysServiceMock{
			GetPrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) {
				return domain.PrivateRelayerRaylsSignKeys{}, service.ErrNoRaylsSignKeys
			},
			CreatePrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) { return priv, nil },
			GetPublicRelayerRaylsSignKeysFunc: func(context.Context) (domain.PublicRelayerRaylsSignKeys, error) {
				return domain.PublicRelayerRaylsSignKeys{}, service.ErrNoRaylsSignKeys
			},
			CreatePublicRelayerRaylsSignKeysFunc: func(context.Context) (domain.PublicRelayerRaylsSignKeys, error) { return pub, nil },
		}

		mgr := newTestManager(keysSvc, stubFactoryAlwaysOK(newWaitRecorder()))

		_, _, err := mgr.Run(context.Background())
		require.NoError(t, err)
		assert.Len(t, keysSvc.CreatePrivateRelayerRaylsSignKeysCalls(), 1)
		assert.Len(t, keysSvc.CreatePublicRelayerRaylsSignKeysCalls(), 1)
	})

	t.Run("first_boot_private_only", func(t *testing.T) {
		priv, pub := sampleBundles(t)
		keysSvc := &KeysServiceMock{
			GetPrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) {
				return domain.PrivateRelayerRaylsSignKeys{}, service.ErrNoRaylsSignKeys
			},
			CreatePrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) { return priv, nil },
			GetPublicRelayerRaylsSignKeysFunc:     func(context.Context) (domain.PublicRelayerRaylsSignKeys, error) { return pub, nil },
		}

		mgr := newTestManager(keysSvc, stubFactoryAlwaysOK(newWaitRecorder()))

		_, _, err := mgr.Run(context.Background())
		require.NoError(t, err)
		assert.Len(t, keysSvc.CreatePrivateRelayerRaylsSignKeysCalls(), 1)
		assert.Empty(t, keysSvc.CreatePublicRelayerRaylsSignKeysCalls())
	})

	t.Run("first_boot_public_only", func(t *testing.T) {
		priv, pub := sampleBundles(t)
		keysSvc := &KeysServiceMock{
			GetPrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) { return priv, nil },
			GetPublicRelayerRaylsSignKeysFunc: func(context.Context) (domain.PublicRelayerRaylsSignKeys, error) {
				return domain.PublicRelayerRaylsSignKeys{}, service.ErrNoRaylsSignKeys
			},
			CreatePublicRelayerRaylsSignKeysFunc: func(context.Context) (domain.PublicRelayerRaylsSignKeys, error) { return pub, nil },
		}

		mgr := newTestManager(keysSvc, stubFactoryAlwaysOK(newWaitRecorder()))

		_, _, err := mgr.Run(context.Background())
		require.NoError(t, err)
		assert.Empty(t, keysSvc.CreatePrivateRelayerRaylsSignKeysCalls())
		assert.Len(t, keysSvc.CreatePublicRelayerRaylsSignKeysCalls(), 1)
	})

	t.Run("private_get_unknown_error", func(t *testing.T) {
		wantErr := errors.New("db exploded")
		keysSvc := &KeysServiceMock{
			GetPrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) {
				return domain.PrivateRelayerRaylsSignKeys{}, wantErr
			},
		}

		mgr := newTestManager(keysSvc, stubFactoryPanic(t))

		_, _, err := mgr.Run(context.Background())
		require.ErrorIs(t, err, wantErr)
		assert.Empty(t, keysSvc.CreatePrivateRelayerRaylsSignKeysCalls())
	})

	t.Run("public_get_unknown_error", func(t *testing.T) {
		priv, _ := sampleBundles(t)
		wantErr := errors.New("db exploded")
		keysSvc := &KeysServiceMock{
			GetPrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) { return priv, nil },
			GetPublicRelayerRaylsSignKeysFunc: func(context.Context) (domain.PublicRelayerRaylsSignKeys, error) {
				return domain.PublicRelayerRaylsSignKeys{}, wantErr
			},
		}

		mgr := newTestManager(keysSvc, stubFactoryPanic(t))

		_, _, err := mgr.Run(context.Background())
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("private_create_error", func(t *testing.T) {
		wantErr := errors.New("encrypt failed")
		keysSvc := &KeysServiceMock{
			GetPrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) {
				return domain.PrivateRelayerRaylsSignKeys{}, service.ErrNoRaylsSignKeys
			},
			CreatePrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) {
				return domain.PrivateRelayerRaylsSignKeys{}, wantErr
			},
		}

		mgr := newTestManager(keysSvc, stubFactoryPanic(t))

		_, _, err := mgr.Run(context.Background())
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("public_create_error", func(t *testing.T) {
		priv, _ := sampleBundles(t)
		wantErr := errors.New("encrypt failed")
		keysSvc := &KeysServiceMock{
			GetPrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) { return priv, nil },
			GetPublicRelayerRaylsSignKeysFunc: func(context.Context) (domain.PublicRelayerRaylsSignKeys, error) {
				return domain.PublicRelayerRaylsSignKeys{}, service.ErrNoRaylsSignKeys
			},
			CreatePublicRelayerRaylsSignKeysFunc: func(context.Context) (domain.PublicRelayerRaylsSignKeys, error) {
				return domain.PublicRelayerRaylsSignKeys{}, wantErr
			},
		}

		mgr := newTestManager(keysSvc, stubFactoryPanic(t))

		_, _, err := mgr.Run(context.Background())
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("empty_keyset_guard", func(t *testing.T) {
		basePriv, basePub := sampleBundles(t)

		cases := []struct {
			name       string
			mutate     func(priv *domain.PrivateRelayerRaylsSignKeys, pub *domain.PublicRelayerRaylsSignKeys)
			errSnippet string
		}{
			{"private PrivateHubKeys",
				func(p *domain.PrivateRelayerRaylsSignKeys, _ *domain.PublicRelayerRaylsSignKeys) {
					p.PrivateHubKeys = nil
				},
				"private relayer PrivateHubKeys"},
			{"private PrivateNodeKeys",
				func(p *domain.PrivateRelayerRaylsSignKeys, _ *domain.PublicRelayerRaylsSignKeys) {
					p.PrivateNodeKeys = nil
				},
				"private relayer PrivateNodeKeys"},
			{"public PrivateChainKeys",
				func(_ *domain.PrivateRelayerRaylsSignKeys, p *domain.PublicRelayerRaylsSignKeys) {
					p.PrivateChainKeys = nil
				},
				"public relayer PrivateChainKeys"},
			{"public PublicChainKeys",
				func(_ *domain.PrivateRelayerRaylsSignKeys, p *domain.PublicRelayerRaylsSignKeys) {
					p.PublicChainKeys = nil
				},
				"public relayer PublicChainKeys"},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				priv, pub := basePriv, basePub
				tc.mutate(&priv, &pub)

				keysSvc := &KeysServiceMock{
					GetPrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) { return priv, nil },
					GetPublicRelayerRaylsSignKeysFunc:  func(context.Context) (domain.PublicRelayerRaylsSignKeys, error) { return pub, nil },
				}

				mgr := newTestManager(keysSvc, stubFactoryPanic(t))

				_, _, err := mgr.Run(context.Background())
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errSnippet)
			})
		}
	})

	t.Run("factory_error_identifies_chain", func(t *testing.T) {
		priv, pub := sampleBundles(t)
		keysSvc := &KeysServiceMock{
			GetPrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) { return priv, nil },
			GetPublicRelayerRaylsSignKeysFunc:  func(context.Context) (domain.PublicRelayerRaylsSignKeys, error) { return pub, nil },
		}

		wantErr := errors.New("factory kaboom")
		var calls int
		stubFactory := func(_ *ethclient.Client, _ common.Address, _ domain.RaylsSignKeyList) (keymanager.AccessManager, error) {
			calls++
			if calls == 2 {
				return nil, wantErr
			}
			return &AccessManagerMock{
				WaitForAuthorizationFunc: func(context.Context, string, []string) error { return nil },
			}, nil
		}

		mgr := newTestManager(keysSvc, stubFactory)

		_, _, err := mgr.Run(context.Background())
		require.ErrorIs(t, err, wantErr)
		assert.Contains(t, err.Error(), "privacy-node")
	})

	t.Run("wait_error_identifies_chain", func(t *testing.T) {
		priv, pub := sampleBundles(t)
		keysSvc := &KeysServiceMock{
			GetPrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) { return priv, nil },
			GetPublicRelayerRaylsSignKeysFunc:  func(context.Context) (domain.PublicRelayerRaylsSignKeys, error) { return pub, nil },
		}

		wantErr := errors.New("wait boom")
		stubFactory := func(_ *ethclient.Client, _ common.Address, _ domain.RaylsSignKeyList) (keymanager.AccessManager, error) {
			return &AccessManagerMock{
				WaitForAuthorizationFunc: func(_ context.Context, component string, _ []string) error {
					if component == "public-chain" {
						return wantErr
					}
					return nil
				},
			}, nil
		}

		mgr := newTestManager(keysSvc, stubFactory)

		_, _, err := mgr.Run(context.Background())
		require.ErrorIs(t, err, wantErr)
		assert.Contains(t, err.Error(), "public-chain")
	})

	t.Run("context_cancellation_propagates", func(t *testing.T) {
		priv, pub := sampleBundles(t)
		keysSvc := &KeysServiceMock{
			GetPrivateRelayerRaylsSignKeysFunc: func(context.Context) (domain.PrivateRelayerRaylsSignKeys, error) { return priv, nil },
			GetPublicRelayerRaylsSignKeysFunc:  func(context.Context) (domain.PublicRelayerRaylsSignKeys, error) { return pub, nil },
		}

		stubFactory := func(_ *ethclient.Client, _ common.Address, _ domain.RaylsSignKeyList) (keymanager.AccessManager, error) {
			return &AccessManagerMock{
				WaitForAuthorizationFunc: func(context.Context, string, []string) error {
					return context.Canceled
				},
			}, nil
		}

		mgr := newTestManager(keysSvc, stubFactory)

		_, _, err := mgr.Run(context.Background())
		require.ErrorIs(t, err, context.Canceled)
	})
}

// --- helpers ---------------------------------------------------------------

// newTestManager builds a SignKeyManager with zero-value ChainConfig.Backend
// slots. The stub BuildAccessManagerFunc never reads backend, so no real
// *ethclient.Client is needed.
func newTestManager(keysSvc keymanager.KeysService, factory keymanager.BuildAccessManagerFunc) *keymanager.SignKeyManager {
	return keymanager.NewSignKeyManager(
		keysSvc,
		keymanager.ChainConfig{},
		keymanager.ChainConfig{},
		keymanager.ChainConfig{},
		keymanager.WithBuildAccessManager(factory),
		keymanager.WithPublicChain(true),
	)
}

func sampleBundles(t *testing.T) (domain.PrivateRelayerRaylsSignKeys, domain.PublicRelayerRaylsSignKeys) {
	t.Helper()
	priv := domain.PrivateRelayerRaylsSignKeys{
		PrivateHubKeys:            mustGenerateKeys(t, 2),
		PrivateHubDvpOperatorKeys: mustGenerateKeys(t, 1),
		PrivateNodeKeys:           mustGenerateKeys(t, 2),
	}
	pub := domain.PublicRelayerRaylsSignKeys{
		PublicChainKeys:  mustGenerateKeys(t, 2),
		PrivateChainKeys: mustGenerateKeys(t, 2),
	}
	return priv, pub
}

func mustGenerateKeys(t *testing.T, n int) domain.RaylsSignKeyList {
	t.Helper()
	keys := make(domain.RaylsSignKeyList, n)
	for i := 0; i < n; i++ {
		k, err := crypto.GenerateKey()
		require.NoError(t, err)
		keys[i] = k
	}
	return keys
}

// waitRecord captures one WaitForAuthorization invocation.
type waitRecord struct {
	component string
	keyHex    []string
}

type waitRecorder struct {
	mu      sync.Mutex
	entries []waitRecord
}

func newWaitRecorder() *waitRecorder { return &waitRecorder{} }

func (r *waitRecorder) record(component string, keys []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = append(r.entries, waitRecord{component: component, keyHex: append([]string(nil), keys...)})
}

func (r *waitRecorder) records() []waitRecord {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]waitRecord(nil), r.entries...)
}

func stubFactoryAlwaysOK(waits *waitRecorder) keymanager.BuildAccessManagerFunc {
	return func(_ *ethclient.Client, _ common.Address, _ domain.RaylsSignKeyList) (keymanager.AccessManager, error) {
		return &AccessManagerMock{
			WaitForAuthorizationFunc: func(_ context.Context, component string, keys []string) error {
				waits.record(component, keys)
				return nil
			},
		}, nil
	}
}

func stubFactoryPanic(t *testing.T) keymanager.BuildAccessManagerFunc {
	return func(_ *ethclient.Client, _ common.Address, _ domain.RaylsSignKeyList) (keymanager.AccessManager, error) {
		t.Fatalf("buildAccessManager should not be called")
		return nil, nil
	}
}

// assertWaitRecord verifies the recorded component label and that every hex
// key round-trips back to the expected address.
func assertWaitRecord(t *testing.T, rec waitRecord, wantComponent string, wantKeys domain.RaylsSignKeyList) {
	t.Helper()
	assert.Equal(t, wantComponent, rec.component)
	require.Len(t, rec.keyHex, len(wantKeys))
	for i, h := range rec.keyHex {
		k, err := crypto.HexToECDSA(h)
		require.NoError(t, err)
		assert.Equal(t, addressOf(wantKeys[i]), crypto.PubkeyToAddress(k.PublicKey))
	}
}

func addressOf(k *ecdsa.PrivateKey) common.Address { return crypto.PubkeyToAddress(k.PublicKey) }
