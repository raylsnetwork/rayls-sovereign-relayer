package keymanager_test

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	ps "github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/ParticipantStorageV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cryptography"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/domain"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/keygen"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/keymanager"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testChainID    = big.NewInt(101)
	testVenChainID = big.NewInt(999)
	testBlock      = uint64(42)
)

func TestParticipantRegistrarRun(t *testing.T) {
	t.Run("happy_path_all_fresh", func(t *testing.T) {
		pair := mustGenerateViewPair(t)
		paymentKeys := samplePaymentSpendKeys(t)

		participants := []ps.ParticipantStructsPrivacyNodeViewData{
			makeViewData(t, big.NewInt(201)),
			makeViewData(t, big.NewInt(202)),
		}

		keysSvc := allFreshKeysService(t, pair, paymentKeys)
		keysSvc.GetRaylsViewKeyPairFunc = func(context.Context, uint64) (domain.RaylsViewKeyPair, error) {
			return domain.RaylsViewKeyPair{}, service.ErrNoRaylsViewKeysSet
		}
		keysSvc.GetPaymentSpendKeyFunc = func(context.Context) (domain.PaymentSpendKeys, error) {
			return domain.PaymentSpendKeys{}, service.ErrNoPaymentSpendKeySet
		}
		storage := &ParticipantStorageMock{
			GetChainViewDataFunc: func(context.Context, *big.Int, *big.Int) (ps.ParticipantStructsPrivacyNodeViewData, error) {
				return ps.ParticipantStructsPrivacyNodeViewData{}, contractclient.ErrNoChainInfo
			},
			SetChainViewDataFunc: func(context.Context, *big.Int, string, *big.Int) error { return nil },
			GetKeyAgreementsFunc: func(context.Context, *big.Int, *big.Int) ([]ps.ParticipantStructsKeyAgreementData, error) {
				return nil, nil
			},
			GetAllParticipantsViewDataFunc: func(context.Context, *big.Int) ([]ps.ParticipantStructsPrivacyNodeViewData, error) {
				return participants, nil
			},
			InitiateKeyAgreementFunc: func(context.Context, *big.Int, []byte, []byte, *big.Int) error { return nil },
			GetAuditInfoFunc: func(context.Context, *big.Int, *big.Int) (ps.ParticipantStructsAuditInfoData, error) {
				return ps.ParticipantStructsAuditInfoData{}, contractclient.ErrNoAuditInfo
			},
			GetVenOperatorChainInfoFunc: func(context.Context, *big.Int) (ps.ParticipantStructsPrivacyNodeViewData, error) {
				return ps.ParticipantStructsPrivacyNodeViewData{ChainId: testVenChainID}, nil
			},
			SetAuditInfoFunc: func(context.Context, *big.Int, *big.Int, string, []byte, []byte) error { return nil },
			GetPaymentSpendPublicKeyFunc: func(context.Context, *big.Int) (*big.Int, error) {
				return nil, contractclient.ErrNoPaymentSpendPublicKey
			},
			SetPaymentSpendPublicKeyFunc: func(context.Context, *big.Int, *big.Int, []common.Address) error { return nil },
		}
		encrypt := okEncryptPrivateKey()

		mgr := newTestRegistrar(t, keysSvc, encrypt, storage)
		require.NoError(t, mgr.Run(context.Background()))

		assert.Len(t, keysSvc.CreateRaylsViewKeyPairCalls(), 1)
		assert.Len(t, storage.SetChainViewDataCalls(), 1)
		assert.Len(t, keysSvc.GenerateKeyAgreementCalls(), 2)
		assert.Len(t, storage.InitiateKeyAgreementCalls(), 2)
		assert.Len(t, keysSvc.CreateKeyAgreementCalls(), 2)
		assert.Len(t, encrypt.EncryptPrivateKeyCalls(), 1)
		assert.Len(t, storage.SetAuditInfoCalls(), 1)
		assert.Len(t, keysSvc.CreatePaymentSpendKeyCalls(), 1)
		assert.Len(t, storage.SetPaymentSpendPublicKeyCalls(), 1)
	})

	t.Run("happy_path_all_registered", func(t *testing.T) {
		pair := mustGenerateViewPair(t)
		paymentKeys := samplePaymentSpendKeys(t)
		onChainHash := mustPaymentHash(t, paymentKeys.SecretKey)
		publicKeyHex := hex.EncodeToString(pair.RaylsViewPublicKey.Bytes())

		participants := []ps.ParticipantStructsPrivacyNodeViewData{
			makeViewData(t, big.NewInt(201)),
		}
		existingAgreements := []ps.ParticipantStructsKeyAgreementData{
			{ChainId: big.NewInt(201)},
		}

		keysSvc := &RegistrarKeysServiceMock{
			GetRaylsViewKeyPairFunc: func(context.Context, uint64) (domain.RaylsViewKeyPair, error) { return pair, nil },
			GetPaymentSpendKeyFunc:  func(context.Context) (domain.PaymentSpendKeys, error) { return paymentKeys, nil },
		}
		storage := &ParticipantStorageMock{
			GetChainViewDataFunc: func(context.Context, *big.Int, *big.Int) (ps.ParticipantStructsPrivacyNodeViewData, error) {
				return ps.ParticipantStructsPrivacyNodeViewData{RaylsViewPublicKey: publicKeyHex}, nil
			},
			GetKeyAgreementsFunc: func(context.Context, *big.Int, *big.Int) ([]ps.ParticipantStructsKeyAgreementData, error) {
				return existingAgreements, nil
			},
			GetAllParticipantsViewDataFunc: func(context.Context, *big.Int) ([]ps.ParticipantStructsPrivacyNodeViewData, error) {
				return participants, nil
			},
			GetAuditInfoFunc: func(context.Context, *big.Int, *big.Int) (ps.ParticipantStructsAuditInfoData, error) {
				return ps.ParticipantStructsAuditInfoData{}, nil
			},
			GetPaymentSpendPublicKeyFunc: func(context.Context, *big.Int) (*big.Int, error) { return onChainHash, nil },
		}

		mgr := newTestRegistrar(t, keysSvc, okEncryptPrivateKey(), storage)
		require.NoError(t, mgr.Run(context.Background()))

		assert.Empty(t, storage.SetChainViewDataCalls())
		assert.Empty(t, storage.InitiateKeyAgreementCalls())
		assert.Empty(t, storage.SetAuditInfoCalls())
		assert.Empty(t, storage.SetPaymentSpendPublicKeyCalls())
	})

	t.Run("view_key_mismatch", func(t *testing.T) {
		pair := mustGenerateViewPair(t)

		keysSvc := &RegistrarKeysServiceMock{
			GetRaylsViewKeyPairFunc: func(context.Context, uint64) (domain.RaylsViewKeyPair, error) { return pair, nil },
		}
		storage := &ParticipantStorageMock{
			GetChainViewDataFunc: func(context.Context, *big.Int, *big.Int) (ps.ParticipantStructsPrivacyNodeViewData, error) {
				return ps.ParticipantStructsPrivacyNodeViewData{RaylsViewPublicKey: "deadbeef"}, nil
			},
		}

		mgr := newTestRegistrar(t, keysSvc, okEncryptPrivateKey(), storage)
		err := mgr.Run(context.Background())
		require.ErrorIs(t, err, keymanager.ErrRaylsViewKeyMismatch)

		assert.Empty(t, storage.GetKeyAgreementsCalls())
		assert.Empty(t, storage.GetAuditInfoCalls())
	})

	t.Run("view_key_create_on_no_keys_set", func(t *testing.T) {
		pair := mustGenerateViewPair(t)
		paymentKeys := samplePaymentSpendKeys(t)

		keysSvc := allFreshKeysService(t, pair, paymentKeys)
		keysSvc.GetRaylsViewKeyPairFunc = func(context.Context, uint64) (domain.RaylsViewKeyPair, error) {
			return domain.RaylsViewKeyPair{}, service.ErrNoRaylsViewKeysSet
		}
		storage := allFreshStorage(t)

		mgr := newTestRegistrar(t, keysSvc, okEncryptPrivateKey(), storage)
		require.NoError(t, mgr.Run(context.Background()))

		require.Len(t, keysSvc.CreateRaylsViewKeyPairCalls(), 1)
		assert.Equal(t, testBlock, keysSvc.CreateRaylsViewKeyPairCalls()[0].InitialBlock)
	})

	t.Run("view_key_create_on_no_applicable", func(t *testing.T) {
		pair := mustGenerateViewPair(t)
		paymentKeys := samplePaymentSpendKeys(t)

		keysSvc := allFreshKeysService(t, pair, paymentKeys)
		keysSvc.GetRaylsViewKeyPairFunc = func(context.Context, uint64) (domain.RaylsViewKeyPair, error) {
			return domain.RaylsViewKeyPair{}, service.ErrNoApplicableRaylsViewKeys
		}
		storage := allFreshStorage(t)

		mgr := newTestRegistrar(t, keysSvc, okEncryptPrivateKey(), storage)
		require.NoError(t, mgr.Run(context.Background()))

		assert.Len(t, keysSvc.CreateRaylsViewKeyPairCalls(), 1)
	})

	t.Run("key_agreements_skip_existing", func(t *testing.T) {
		pair := mustGenerateViewPair(t)
		paymentKeys := samplePaymentSpendKeys(t)

		participants := []ps.ParticipantStructsPrivacyNodeViewData{
			makeViewData(t, big.NewInt(201)),
			makeViewData(t, big.NewInt(202)),
			makeViewData(t, big.NewInt(203)),
		}
		existing := []ps.ParticipantStructsKeyAgreementData{
			{ChainId: big.NewInt(202)},
		}

		keysSvc := allFreshKeysService(t, pair, paymentKeys)
		storage := allFreshStorage(t)
		storage.GetKeyAgreementsFunc = func(context.Context, *big.Int, *big.Int) ([]ps.ParticipantStructsKeyAgreementData, error) {
			return existing, nil
		}
		storage.GetAllParticipantsViewDataFunc = func(context.Context, *big.Int) ([]ps.ParticipantStructsPrivacyNodeViewData, error) {
			return participants, nil
		}

		mgr := newTestRegistrar(t, keysSvc, okEncryptPrivateKey(), storage)
		require.NoError(t, mgr.Run(context.Background()))

		assert.Len(t, keysSvc.GenerateKeyAgreementCalls(), 2)
		assert.Len(t, storage.InitiateKeyAgreementCalls(), 2)

		called := []string{}
		for _, c := range storage.InitiateKeyAgreementCalls() {
			called = append(called, c.ChainID.String())
		}
		assert.ElementsMatch(t, []string{"201", "203"}, called)
	})

	t.Run("key_agreement_outdated_skips_create", func(t *testing.T) {
		pair := mustGenerateViewPair(t)
		paymentKeys := samplePaymentSpendKeys(t)

		participants := []ps.ParticipantStructsPrivacyNodeViewData{
			makeViewData(t, big.NewInt(201)),
			makeViewData(t, big.NewInt(202)),
		}

		keysSvc := allFreshKeysService(t, pair, paymentKeys)
		storage := allFreshStorage(t)
		storage.GetAllParticipantsViewDataFunc = func(context.Context, *big.Int) ([]ps.ParticipantStructsPrivacyNodeViewData, error) {
			return participants, nil
		}
		storage.InitiateKeyAgreementFunc = func(_ context.Context, chainID *big.Int, _, _ []byte, _ *big.Int) error {
			if chainID.Cmp(big.NewInt(201)) == 0 {
				return contractclient.ErrOutdatedKeyAgreement
			}
			return nil
		}

		mgr := newTestRegistrar(t, keysSvc, okEncryptPrivateKey(), storage)
		require.NoError(t, mgr.Run(context.Background()))

		assert.Len(t, storage.InitiateKeyAgreementCalls(), 2)
		require.Len(t, keysSvc.CreateKeyAgreementCalls(), 1)
		assert.Equal(t, "202", keysSvc.CreateKeyAgreementCalls()[0].ChainID.String())
	})

	t.Run("key_agreement_initiate_error_propagates", func(t *testing.T) {
		pair := mustGenerateViewPair(t)
		paymentKeys := samplePaymentSpendKeys(t)

		participants := []ps.ParticipantStructsPrivacyNodeViewData{
			makeViewData(t, big.NewInt(201)),
			makeViewData(t, big.NewInt(202)),
		}
		wantErr := errors.New("initiate boom")

		keysSvc := allFreshKeysService(t, pair, paymentKeys)
		storage := allFreshStorage(t)
		storage.GetAllParticipantsViewDataFunc = func(context.Context, *big.Int) ([]ps.ParticipantStructsPrivacyNodeViewData, error) {
			return participants, nil
		}
		storage.InitiateKeyAgreementFunc = func(context.Context, *big.Int, []byte, []byte, *big.Int) error { return wantErr }

		mgr := newTestRegistrar(t, keysSvc, okEncryptPrivateKey(), storage)
		err := mgr.Run(context.Background())
		require.ErrorIs(t, err, wantErr)

		assert.Len(t, storage.InitiateKeyAgreementCalls(), 1)
		assert.Empty(t, storage.GetAuditInfoCalls())
	})

	t.Run("audit_info_already_set", func(t *testing.T) {
		pair := mustGenerateViewPair(t)
		paymentKeys := samplePaymentSpendKeys(t)

		keysSvc := allFreshKeysService(t, pair, paymentKeys)
		storage := allFreshStorage(t)
		storage.GetAuditInfoFunc = func(context.Context, *big.Int, *big.Int) (ps.ParticipantStructsAuditInfoData, error) {
			return ps.ParticipantStructsAuditInfoData{}, nil
		}
		encrypt := okEncryptPrivateKey()

		mgr := newTestRegistrar(t, keysSvc, encrypt, storage)
		require.NoError(t, mgr.Run(context.Background()))

		assert.Empty(t, encrypt.EncryptPrivateKeyCalls())
		assert.Empty(t, storage.SetAuditInfoCalls())
	})

	t.Run("audit_info_set_uses_ven_chain_id", func(t *testing.T) {
		pair := mustGenerateViewPair(t)
		paymentKeys := samplePaymentSpendKeys(t)

		keysSvc := allFreshKeysService(t, pair, paymentKeys)
		storage := allFreshStorage(t)
		encrypt := okEncryptPrivateKey()

		mgr := newTestRegistrar(t, keysSvc, encrypt, storage)
		require.NoError(t, mgr.Run(context.Background()))

		require.Len(t, encrypt.EncryptPrivateKeyCalls(), 1)
		assert.Equal(t, testVenChainID.Uint64(), encrypt.EncryptPrivateKeyCalls()[0].ChainID)
		assert.Equal(t, testBlock, encrypt.EncryptPrivateKeyCalls()[0].BlockNumber)

		require.Len(t, storage.SetAuditInfoCalls(), 1)
		call := storage.SetAuditInfoCalls()[0]
		assert.Equal(t, uint64(77), call.BlockNumber.Uint64()) // mirrors okEncryptPrivateKey's initialBlock
	})

	t.Run("audit_info_ven_operator_lookup_error", func(t *testing.T) {
		pair := mustGenerateViewPair(t)
		paymentKeys := samplePaymentSpendKeys(t)
		wantErr := errors.New("ven lookup boom")

		keysSvc := allFreshKeysService(t, pair, paymentKeys)
		storage := allFreshStorage(t)
		storage.GetVenOperatorChainInfoFunc = func(context.Context, *big.Int) (ps.ParticipantStructsPrivacyNodeViewData, error) {
			return ps.ParticipantStructsPrivacyNodeViewData{}, wantErr
		}
		encrypt := okEncryptPrivateKey()

		mgr := newTestRegistrar(t, keysSvc, encrypt, storage)
		err := mgr.Run(context.Background())
		require.ErrorIs(t, err, wantErr)

		assert.Empty(t, encrypt.EncryptPrivateKeyCalls())
		assert.Empty(t, storage.GetPaymentSpendPublicKeyCalls())
	})

	t.Run("payment_key_mismatch", func(t *testing.T) {
		pair := mustGenerateViewPair(t)
		paymentKeys := samplePaymentSpendKeys(t)

		keysSvc := allFreshKeysService(t, pair, paymentKeys)
		storage := allFreshStorage(t)
		storage.GetPaymentSpendPublicKeyFunc = func(context.Context, *big.Int) (*big.Int, error) {
			return big.NewInt(987654321), nil
		}

		mgr := newTestRegistrar(t, keysSvc, okEncryptPrivateKey(), storage)
		err := mgr.Run(context.Background())
		require.ErrorIs(t, err, keymanager.ErrPaymentSpendKeyMismatch)

		assert.Empty(t, storage.SetPaymentSpendPublicKeyCalls())
	})

	t.Run("payment_key_create_on_missing", func(t *testing.T) {
		pair := mustGenerateViewPair(t)
		paymentKeys := samplePaymentSpendKeys(t)

		keysSvc := allFreshKeysService(t, pair, paymentKeys)
		keysSvc.GetPaymentSpendKeyFunc = func(context.Context) (domain.PaymentSpendKeys, error) {
			return domain.PaymentSpendKeys{}, service.ErrNoPaymentSpendKeySet
		}
		storage := allFreshStorage(t)

		mgr := newTestRegistrar(t, keysSvc, okEncryptPrivateKey(), storage)
		require.NoError(t, mgr.Run(context.Background()))

		require.Len(t, keysSvc.CreatePaymentSpendKeyCalls(), 1)
		assert.Equal(t, testChainID, keysSvc.CreatePaymentSpendKeyCalls()[0].ChainID)
	})

	t.Run("payment_key_set_with_relayer_hub_addresses", func(t *testing.T) {
		pair := mustGenerateViewPair(t)
		paymentKeys := samplePaymentSpendKeys(t)

		keysSvc := allFreshKeysService(t, pair, paymentKeys)
		storage := allFreshStorage(t)
		hubKeys := mustGenerateKeys(t, 3)

		mgr := keymanager.NewParticipantRegistrar(
			testChainID, testVenChainID,
			keymanager.ChainConfig{},
			keysSvc, okEncryptPrivateKey(),
			hubKeys,
			keymanager.WithBuildParticipantStorage(func(
				*ethclient.Client, common.Address, *big.Int, *big.Int, domain.RaylsSignKeyList,
			) (keymanager.ParticipantStorage, error) {
				return storage, nil
			}),
			keymanager.WithBlockNumberFunc(okBlockNumber()),
		)
		require.NoError(t, mgr.Run(context.Background()))

		require.Len(t, storage.SetPaymentSpendPublicKeyCalls(), 1)
		gotAddrs := storage.SetPaymentSpendPublicKeyCalls()[0].PlAddresses

		wantAddrs := make([]common.Address, len(hubKeys))
		for i, k := range hubKeys {
			wantAddrs[i] = crypto.PubkeyToAddress(k.PublicKey)
		}
		assert.Equal(t, wantAddrs, gotAddrs)
	})

	t.Run("block_number_error_propagates", func(t *testing.T) {
		pair := mustGenerateViewPair(t)
		paymentKeys := samplePaymentSpendKeys(t)
		wantErr := errors.New("node down")

		keysSvc := allFreshKeysService(t, pair, paymentKeys)
		storage := allFreshStorage(t)

		mgr := keymanager.NewParticipantRegistrar(
			testChainID, testVenChainID,
			keymanager.ChainConfig{},
			keysSvc, okEncryptPrivateKey(),
			mustGenerateKeys(t, 1),
			keymanager.WithBuildParticipantStorage(func(
				*ethclient.Client, common.Address, *big.Int, *big.Int, domain.RaylsSignKeyList,
			) (keymanager.ParticipantStorage, error) {
				return storage, nil
			}),
			keymanager.WithBlockNumberFunc(func(context.Context, *ethclient.Client) (uint64, error) {
				return 0, wantErr
			}),
		)

		err := mgr.Run(context.Background())
		require.ErrorIs(t, err, wantErr)

		assert.Empty(t, storage.GetChainViewDataCalls())
	})

	t.Run("storage_factory_error_propagates", func(t *testing.T) {
		wantErr := errors.New("factory boom")

		mgr := keymanager.NewParticipantRegistrar(
			testChainID, testVenChainID,
			keymanager.ChainConfig{},
			&RegistrarKeysServiceMock{}, okEncryptPrivateKey(),
			mustGenerateKeys(t, 1),
			keymanager.WithBuildParticipantStorage(func(
				*ethclient.Client, common.Address, *big.Int, *big.Int, domain.RaylsSignKeyList,
			) (keymanager.ParticipantStorage, error) {
				return nil, wantErr
			}),
			keymanager.WithBlockNumberFunc(okBlockNumber()),
		)

		err := mgr.Run(context.Background())
		require.ErrorIs(t, err, wantErr)
	})
}

// --- helpers ---------------------------------------------------------------

func newTestRegistrar(
	t *testing.T,
	keysSvc keymanager.RegistrarKeysService,
	encrypt keymanager.EncryptPrivateKeyer,
	storage keymanager.ParticipantStorage,
) *keymanager.ParticipantRegistrar {
	t.Helper()
	return keymanager.NewParticipantRegistrar(
		testChainID, testVenChainID,
		keymanager.ChainConfig{},
		keysSvc, encrypt,
		mustGenerateKeys(t, 1),
		keymanager.WithBuildParticipantStorage(func(
			*ethclient.Client, common.Address, *big.Int, *big.Int, domain.RaylsSignKeyList,
		) (keymanager.ParticipantStorage, error) {
			return storage, nil
		}),
		keymanager.WithBlockNumberFunc(okBlockNumber()),
	)
}

func okBlockNumber() keymanager.GetBlockNumberFunc {
	return func(context.Context, *ethclient.Client) (uint64, error) { return testBlock, nil }
}

func okEncryptPrivateKey() *EncryptPrivateKeyerMock {
	return &EncryptPrivateKeyerMock{
		EncryptPrivateKeyFunc: func(context.Context, uint64, uint64) (string, []byte, []byte, uint64, error) {
			return "pubkeyhex", []byte("encrypted"), []byte("mac"), 77, nil
		},
	}
}

// allFreshKeysService returns a RegistrarKeysServiceMock where every Get*
// returns a populated value; individual tests override fields as needed.
func allFreshKeysService(
	t *testing.T,
	pair domain.RaylsViewKeyPair,
	paymentKeys domain.PaymentSpendKeys,
) *RegistrarKeysServiceMock {
	t.Helper()
	return &RegistrarKeysServiceMock{
		GetRaylsViewKeyPairFunc: func(context.Context, uint64) (domain.RaylsViewKeyPair, error) { return pair, nil },
		CreateRaylsViewKeyPairFunc: func(_ context.Context, blockNum uint64) (domain.RaylsViewKeyPair, error) {
			return pair, nil
		},
		GenerateKeyAgreementFunc: func(*big.Int, []byte) ([]byte, []byte, []byte, error) {
			return []byte("ct"), []byte("ss"), []byte("dg"), nil
		},
		CreateKeyAgreementFunc: func(context.Context, *big.Int, []byte, uint64) error { return nil },
		GetPaymentSpendKeyFunc:  func(context.Context) (domain.PaymentSpendKeys, error) { return paymentKeys, nil },
		CreatePaymentSpendKeyFunc: func(_ context.Context, chainId *big.Int) (domain.PaymentSpendKeys, error) {
			return paymentKeys, nil
		},
	}
}

// allFreshStorage returns a ParticipantStorageMock where every on-chain read
// returns the "not yet registered" sentinel and every write succeeds.
func allFreshStorage(t *testing.T) *ParticipantStorageMock {
	t.Helper()
	return &ParticipantStorageMock{
		GetChainViewDataFunc: func(context.Context, *big.Int, *big.Int) (ps.ParticipantStructsPrivacyNodeViewData, error) {
			return ps.ParticipantStructsPrivacyNodeViewData{}, contractclient.ErrNoChainInfo
		},
		SetChainViewDataFunc: func(context.Context, *big.Int, string, *big.Int) error { return nil },
		GetKeyAgreementsFunc: func(context.Context, *big.Int, *big.Int) ([]ps.ParticipantStructsKeyAgreementData, error) {
			return nil, nil
		},
		GetAllParticipantsViewDataFunc: func(context.Context, *big.Int) ([]ps.ParticipantStructsPrivacyNodeViewData, error) {
			return nil, nil
		},
		InitiateKeyAgreementFunc: func(context.Context, *big.Int, []byte, []byte, *big.Int) error { return nil },
		GetAuditInfoFunc: func(context.Context, *big.Int, *big.Int) (ps.ParticipantStructsAuditInfoData, error) {
			return ps.ParticipantStructsAuditInfoData{}, contractclient.ErrNoAuditInfo
		},
		GetVenOperatorChainInfoFunc: func(context.Context, *big.Int) (ps.ParticipantStructsPrivacyNodeViewData, error) {
			return ps.ParticipantStructsPrivacyNodeViewData{ChainId: testVenChainID}, nil
		},
		SetAuditInfoFunc: func(context.Context, *big.Int, *big.Int, string, []byte, []byte) error { return nil },
		GetPaymentSpendPublicKeyFunc: func(context.Context, *big.Int) (*big.Int, error) {
			return nil, contractclient.ErrNoPaymentSpendPublicKey
		},
		SetPaymentSpendPublicKeyFunc: func(context.Context, *big.Int, *big.Int, []common.Address) error { return nil },
	}
}

func mustGenerateViewPair(t *testing.T) domain.RaylsViewKeyPair {
	t.Helper()
	pair, err := keygen.GenerateRaylsViewKeys()
	require.NoError(t, err)
	pair.InitialBlock = testBlock
	return pair
}

func samplePaymentSpendKeys(t *testing.T) domain.PaymentSpendKeys {
	t.Helper()
	sk, err := keygen.GenerateRandomModJubJubPrimeSubGroupWithChainId(testChainID)
	require.NoError(t, err)
	pk, err := keygen.GetPaymentSpendPublicKeyFromSpendSecretKey(sk)
	require.NoError(t, err)
	return domain.PaymentSpendKeys{SecretKey: sk, PublicKey: pk}
}

func mustPaymentHash(t *testing.T, sk *big.Int) *big.Int {
	t.Helper()
	hash, err := cryptography.GetPoseidonHashModNumber(
		[]*big.Int{sk, sk},
		cryptography.JubJubPrimeSubGroup,
	)
	require.NoError(t, err)
	return hash
}

func makeViewData(t *testing.T, chainID *big.Int) ps.ParticipantStructsPrivacyNodeViewData {
	t.Helper()
	dummy := mustGenerateViewPair(t)
	return ps.ParticipantStructsPrivacyNodeViewData{
		ChainId:            chainID,
		RaylsViewPublicKey: hex.EncodeToString(dummy.RaylsViewPublicKey.Bytes()),
		BlockNumber:        big.NewInt(int64(testBlock)),
	}
}
