//go:generate moq --pkg service_test --out keys_mock_test.go . Encryptor RaylsViewKeysRepository RaylsSignKeysRepository SharedSecretsRepository EnygmaSelfSecretsRepository PaymentSpendKeysRepository

package service_test

import (
	"context"
	"crypto/mlkem"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cryptography"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/domain"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/service"
	"github.com/stretchr/testify/require"
)

func TestCreateRaylsViewKeyPair(t *testing.T) {
	t.Run("generates View keys and persists them with initial PNH block number", func(t *testing.T) {
		wantInitialBlock := uint64(10)

		// Setup mocks
		encryptorMock := &EncryptorMock{
			EncryptFunc: func(bytes []byte) ([]byte, error) {
				return []byte("encrypted_" + string(bytes)), nil
			},
		}

		viewRepoMock := &RaylsViewKeysRepositoryMock{
			CreateFunc: func(ctx context.Context, key domain.EncryptedRaylsViewKeyPair) error {
				return nil
			},
		}

		ecdsaRepoMock := &RaylsSignKeysRepositoryMock{}
		sharedSecretsRepoMock := &SharedSecretsRepositoryMock{}
		enygmaSelfSecretsRepoMock := &EnygmaSelfSecretsRepositoryMock{}
		paymentSpendRepoMock := &PaymentSpendKeysRepositoryMock{}

		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		pair, err := svc.CreateRaylsViewKeyPair(context.Background(), wantInitialBlock)

		require.Nil(t, err)
		require.NotEmpty(t, pair, "didn't return View keys")
		require.Equal(t, wantInitialBlock, pair.InitialBlock, "didn't set initial block number")

		// Verify encrypt was called during Enygma rayls view key pair encryption
		require.Len(t, encryptorMock.EncryptCalls(), 1, "should encrypt private key")

		// Verify repository create was called
		require.Len(t, viewRepoMock.CreateCalls(), 1, "should persist keys in repo")
		require.Equal(t, wantInitialBlock, viewRepoMock.CreateCalls()[0].Key.InitialBlock)
	})

	t.Run("returns KeyRepositoryError on unknown error from repo on Create", func(t *testing.T) {
		wantErrMsg := "example error"
		initialBlock := uint64(10)

		// Setup mocks
		encryptorMock := &EncryptorMock{
			EncryptFunc: func(bytes []byte) ([]byte, error) {
				return []byte("encrypted_" + string(bytes)), nil
			},
		}

		viewRepoMock := &RaylsViewKeysRepositoryMock{
			CreateFunc: func(ctx context.Context, key domain.EncryptedRaylsViewKeyPair) error {
				return errors.New(wantErrMsg)
			},
		}

		ecdsaRepoMock := &RaylsSignKeysRepositoryMock{}
		sharedSecretsRepoMock := &SharedSecretsRepositoryMock{}
		enygmaSelfSecretsRepoMock := &EnygmaSelfSecretsRepositoryMock{}
		paymentSpendRepoMock := &PaymentSpendKeysRepositoryMock{}

		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		_, err := svc.CreateRaylsViewKeyPair(context.Background(), initialBlock)

		require.NotNil(t, err)

		var keyRepositoryErr service.KeyRepositoryError
		require.IsType(t, &keyRepositoryErr, err)
		require.Equal(t, wantErrMsg, err.Error())
	})
}

func TestGetRaylsViewKeyPair(t *testing.T) {
	t.Run("returns Enygma rayls view key pair applicable for block number", func(t *testing.T) {
		// Generate a valid ML-KEM key pair for testing
		dk, err := mlkem.GenerateKey768()
		require.Nil(t, err)
		privateKeyBytes := dk.Bytes()
		publicKeyBytes := dk.EncapsulationKey().Bytes()

		// Setup mocks
		encryptorMock := &EncryptorMock{
			EncryptFunc: func(bytes []byte) ([]byte, error) {
				return []byte("encrypted_" + string(bytes)), nil
			},
			DecryptFunc: func(bytes []byte) ([]byte, error) {
				// Return the actual ML-KEM private key bytes
				return privateKeyBytes, nil
			},
		}

		// Mock repository to return encrypted Enygma rayls view key pair for block 200
		viewRepoMock := &RaylsViewKeysRepositoryMock{
			CreateFunc: func(ctx context.Context, key domain.EncryptedRaylsViewKeyPair) error {
				return nil
			},
			GetForBlockNumberFunc: func(ctx context.Context, blockNumber uint64) (domain.EncryptedRaylsViewKeyPair, error) {
				return domain.EncryptedRaylsViewKeyPair{
					InitialBlock:                 200,
					EncryptedRaylsViewPrivateKey: []byte("encrypted_private"),
					RaylsViewPublicKey:           publicKeyBytes,
				}, nil
			},
		}

		ecdsaRepoMock := &RaylsSignKeysRepositoryMock{}
		sharedSecretsRepoMock := &SharedSecretsRepositoryMock{}
		enygmaSelfSecretsRepoMock := &EnygmaSelfSecretsRepositoryMock{}
		paymentSpendRepoMock := &PaymentSpendKeysRepositoryMock{}

		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		blockNumber := uint64(250)
		gotPair, err := svc.GetRaylsViewKeyPair(context.Background(), blockNumber)

		require.Nil(t, err)
		require.Equal(t, uint64(200), gotPair.InitialBlock, "should return pair for block 200")

		// Verify decrypt was called
		require.Len(t, encryptorMock.DecryptCalls(), 1, "should decrypt private key")
		require.Len(t, viewRepoMock.GetForBlockNumberCalls(), 1, "should query repository")
	})

	t.Run("returns ErrNoKeys on empty keys repository", func(t *testing.T) {
		wantErr := service.ErrNoRaylsViewKeysSet

		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		viewRepoMock.GetForBlockNumberFunc = func(ctx context.Context, blockNumber uint64) (domain.EncryptedRaylsViewKeyPair, error) {
			return domain.EncryptedRaylsViewKeyPair{}, service.ErrNoRaylsViewKeysSet
		}
		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		blockNumber := uint64(250)
		_, gotErr := svc.GetRaylsViewKeyPair(context.Background(), blockNumber)

		require.NotNil(t, gotErr, "didn't return error")
		require.ErrorIs(t, gotErr, wantErr, "returned wrong sentinel error")
	})

	t.Run("returns KeyRepositoryError on unknown error from GetAll", func(t *testing.T) {
		wantErrMsg := "example error"
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		viewRepoMock.GetForBlockNumberFunc = func(ctx context.Context, blockNumber uint64) (domain.EncryptedRaylsViewKeyPair, error) {
			return domain.EncryptedRaylsViewKeyPair{}, errors.New(wantErrMsg)
		}
		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		_, err := svc.GetRaylsViewKeyPair(context.Background(), 100)

		require.NotNil(t, err, "didn't return error")

		var keyRepositoryErr service.KeyRepositoryError
		require.IsType(t, &keyRepositoryErr, err, "returned wrong error type")
		require.Equal(t, wantErrMsg, err.Error(), "didn't set error message")
	})
}

// Helper function to create default mocks
func createDefaultMocks() (*EncryptorMock, *RaylsViewKeysRepositoryMock, *RaylsSignKeysRepositoryMock, *SharedSecretsRepositoryMock, *EnygmaSelfSecretsRepositoryMock, *PaymentSpendKeysRepositoryMock) {
	encryptorMock := &EncryptorMock{
		EncryptFunc: func(bytes []byte) ([]byte, error) {
			return []byte("encrypted_" + string(bytes)), nil
		},
		DecryptFunc: func(bytes []byte) ([]byte, error) {
			return bytes[10:], nil // remove "encrypted_" prefix
		},
	}

	viewRepoMock := &RaylsViewKeysRepositoryMock{
		CreateFunc: func(ctx context.Context, key domain.EncryptedRaylsViewKeyPair) error {
			return nil
		},
		GetForBlockNumberFunc: func(ctx context.Context, blockNumber uint64) (domain.EncryptedRaylsViewKeyPair, error) {
			return domain.EncryptedRaylsViewKeyPair{}, service.ErrNoRaylsViewKeysSet
		},
	}

	ecdsaRepoMock := &RaylsSignKeysRepositoryMock{
		CreatePublicRelayerRaylsSignKeysFunc:  func(ctx context.Context, keys domain.EncryptedPublicRelayerRaylsSignKeys) error { return nil },
		CreatePrivateRelayerRaylsSignKeysFunc: func(ctx context.Context, keys domain.EncryptedPrivateRelayerRaylsSignKeys) error { return nil },
		CreateAtomicServiceRaylsSignKeysFunc:  func(ctx context.Context, keys domain.EncryptedAtomicServiceRaylsSignKeys) error { return nil },
		GetPublicRelayerRaylsSignKeysFunc: func(ctx context.Context) (domain.EncryptedPublicRelayerRaylsSignKeys, error) {
			return domain.EncryptedPublicRelayerRaylsSignKeys{}, service.ErrNoRaylsSignKeys
		},
		GetPrivateRelayerRaylsSignKeysFunc: func(ctx context.Context) (domain.EncryptedPrivateRelayerRaylsSignKeys, error) {
			return domain.EncryptedPrivateRelayerRaylsSignKeys{}, service.ErrNoRaylsSignKeys
		},
		GetAtomicServiceRaylsSignKeysFunc: func(ctx context.Context) (domain.EncryptedAtomicServiceRaylsSignKeys, error) {
			return domain.EncryptedAtomicServiceRaylsSignKeys{}, service.ErrNoRaylsSignKeys
		},
	}

	sharedSecretsRepoMock := &SharedSecretsRepositoryMock{
		CreateFunc: func(ctx context.Context, secret domain.EncryptedSharedSecret) error { return nil },
		GetAllFunc: func(ctx context.Context, blockNumber uint64) ([]domain.EncryptedSharedSecret, error) {
			return []domain.EncryptedSharedSecret{}, nil
		},
		GetByChainIdsFunc: func(ctx context.Context, chainIds []*big.Int, blockNumber uint64) ([]domain.EncryptedSharedSecret, error) {
			return []domain.EncryptedSharedSecret{}, nil
		},
		GetByChainIdFunc: func(ctx context.Context, chainId string, blockNumber uint64) (domain.EncryptedSharedSecret, error) {
			return domain.EncryptedSharedSecret{}, nil
		},
	}

	enygmaSelfSecretsRepoMock := &EnygmaSelfSecretsRepositoryMock{
		CreateFunc: func(ctx context.Context, secret domain.EncryptedEnygmaSelfSecret) error { return nil },
		GetByBlockNumberAndResourceFunc: func(ctx context.Context, blockNumber uint64, resourceID []byte) (domain.EncryptedEnygmaSelfSecret, error) {
			return domain.EncryptedEnygmaSelfSecret{}, service.ErrNoApplicableEnygmaSelfSecret
		},
	}

	paymentSpendRepoMock := &PaymentSpendKeysRepositoryMock{
		CreateFunc: func(ctx context.Context, paymentSpendKeys domain.EncryptedPaymentSpendKeys) error { return nil },
		GetFunc: func(ctx context.Context) (domain.EncryptedPaymentSpendKeys, error) {
			return domain.EncryptedPaymentSpendKeys{}, service.ErrNoPaymentSpendKeySet
		},
	}

	return encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock
}

func TestCreateRaylsSignKeys(t *testing.T) {
	t.Run("creates rayls sign keys and returns them", func(t *testing.T) {
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		ecdsaRepoMock.CreatePrivateRelayerRaylsSignKeysFunc = func(ctx context.Context, keys domain.EncryptedPrivateRelayerRaylsSignKeys) error { return nil }
		ecdsaRepoMock.GetPrivateRelayerRaylsSignKeysFunc = func(ctx context.Context) (domain.EncryptedPrivateRelayerRaylsSignKeys, error) {
			// Return some encrypted keys
			return domain.EncryptedPrivateRelayerRaylsSignKeys{}, nil
		}
		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		keys, err := svc.CreatePrivateRelayerRaylsSignKeys(context.Background())

		require.Nil(t, err)
		require.NotEmpty(t, keys, "didn't return keys")
	})
}

func TestGetRaylsSignKeys(t *testing.T) {
	t.Run("returns rayls sign keys", func(t *testing.T) {
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		ecdsaRepoMock.CreatePrivateRelayerRaylsSignKeysFunc = func(ctx context.Context, keys domain.EncryptedPrivateRelayerRaylsSignKeys) error { return nil }
		ecdsaRepoMock.GetPrivateRelayerRaylsSignKeysFunc = func(ctx context.Context) (domain.EncryptedPrivateRelayerRaylsSignKeys, error) {
			// Return some encrypted keys
			return domain.EncryptedPrivateRelayerRaylsSignKeys{}, nil
		}
		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		_, err := svc.CreatePrivateRelayerRaylsSignKeys(context.Background())
		require.Nil(t, err)

		keys, err := svc.GetPrivateRelayerRaylsSignKeys(context.Background())
		require.Nil(t, err)

		require.NotEmpty(t, keys, "didn't return keys")
	})

	t.Run("returns ErrNoRaylsSignKeys", func(t *testing.T) {
		wantErr := service.ErrNoRaylsSignKeys

		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		ecdsaRepoMock.GetPrivateRelayerRaylsSignKeysFunc = func(ctx context.Context) (domain.EncryptedPrivateRelayerRaylsSignKeys, error) {
			return domain.EncryptedPrivateRelayerRaylsSignKeys{}, wantErr
		}
		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		keys, err := svc.GetPrivateRelayerRaylsSignKeys(context.Background())

		require.Empty(t, keys)
		require.NotNil(t, err)
		require.ErrorIs(t, err, wantErr)
	})
}

func TestGetPrivateRelayerRaylsSignPublicAddresses(t *testing.T) {
	// Generate distinct key triples so each address bucket is verifiable.
	hubKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	dvpOperatorKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	chainKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	wantHubAddr := crypto.PubkeyToAddress(hubKey.PublicKey)
	wantDvpOperatorAddr := crypto.PubkeyToAddress(dvpOperatorKey.PublicKey)
	wantChainAddr := crypto.PubkeyToAddress(chainKey.PublicKey)

	t.Run("returns addresses derived from each key bucket", func(t *testing.T) {
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()

		// Encrypt the three key buckets through the mock encryptor so that
		// service.GetPrivateRelayerRaylsSignPublicAddresses can decrypt them
		// via the encryptor mirror.
		plain := domain.PrivateRelayerRaylsSignKeys{
			PrivateHubKeys:            domain.RaylsSignKeyList{hubKey},
			PrivateHubDvpOperatorKeys: domain.RaylsSignKeyList{dvpOperatorKey},
			PrivateNodeKeys:           domain.RaylsSignKeyList{chainKey},
		}
		encrypted, encryptErr := plain.Encrypt(encryptorMock)
		require.NoError(t, encryptErr)

		ecdsaRepoMock.GetPrivateRelayerRaylsSignKeysFunc = func(ctx context.Context) (domain.EncryptedPrivateRelayerRaylsSignKeys, error) {
			return encrypted, nil
		}

		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		got, err := svc.GetPrivateRelayerRaylsSignPublicAddresses(context.Background())
		require.NoError(t, err)

		require.Equal(t, domain.AddressList{wantHubAddr}, got.PrivateHubAddresses)
		require.Equal(t, domain.AddressList{wantDvpOperatorAddr}, got.PrivateHubDvpOperatorAddresses)
		require.Equal(t, domain.AddressList{wantChainAddr}, got.PrivateChainAddresses)

		// Sanity: addresses must not collide across buckets.
		require.NotEqual(t, got.PrivateHubAddresses[0], got.PrivateHubDvpOperatorAddresses[0])
		require.NotEqual(t, got.PrivateHubAddresses[0], got.PrivateChainAddresses[0])
		require.NotEqual(t, got.PrivateHubDvpOperatorAddresses[0], got.PrivateChainAddresses[0])
	})

	t.Run("returns ErrNoRaylsSignKeys when the repo has no keys", func(t *testing.T) {
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		ecdsaRepoMock.GetPrivateRelayerRaylsSignKeysFunc = func(ctx context.Context) (domain.EncryptedPrivateRelayerRaylsSignKeys, error) {
			return domain.EncryptedPrivateRelayerRaylsSignKeys{}, service.ErrNoRaylsSignKeys
		}

		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		got, err := svc.GetPrivateRelayerRaylsSignPublicAddresses(context.Background())
		require.ErrorIs(t, err, service.ErrNoRaylsSignKeys)
		require.Equal(t, domain.PrivateRelayerRaylsSignPublicAddresses{}, got)
	})

	t.Run("wraps non-ErrNoRaylsSignKeys repo errors as KeyRepositoryError", func(t *testing.T) {
		// Covers the post-flatten branch where repo errors that are not
		// ErrNoRaylsSignKeys are wrapped with NewKeyRepositoryError.
		wantMsg := "db connection lost"
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		ecdsaRepoMock.GetPrivateRelayerRaylsSignKeysFunc = func(ctx context.Context) (domain.EncryptedPrivateRelayerRaylsSignKeys, error) {
			return domain.EncryptedPrivateRelayerRaylsSignKeys{}, errors.New(wantMsg)
		}

		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		got, err := svc.GetPrivateRelayerRaylsSignPublicAddresses(context.Background())
		require.Error(t, err)
		require.NotErrorIs(t, err, service.ErrNoRaylsSignKeys)
		var keyRepoErr service.KeyRepositoryError
		require.IsType(t, &keyRepoErr, err)
		require.Equal(t, wantMsg, err.Error())
		require.Equal(t, domain.PrivateRelayerRaylsSignPublicAddresses{}, got)
	})

	t.Run("returns empty (non-nil) AddressList for each empty key bucket", func(t *testing.T) {
		// privateKeyListToAddressList must never return nil for an empty
		// input — defensive against callers that expect a usable slice.
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		empty, encryptErr := domain.PrivateRelayerRaylsSignKeys{}.Encrypt(encryptorMock)
		require.NoError(t, encryptErr)
		ecdsaRepoMock.GetPrivateRelayerRaylsSignKeysFunc = func(ctx context.Context) (domain.EncryptedPrivateRelayerRaylsSignKeys, error) {
			return empty, nil
		}

		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		got, err := svc.GetPrivateRelayerRaylsSignPublicAddresses(context.Background())
		require.NoError(t, err)

		require.NotNil(t, got.PrivateHubAddresses, "PrivateHubAddresses should be empty slice, not nil")
		require.NotNil(t, got.PrivateHubDvpOperatorAddresses, "PrivateHubDvpOperatorAddresses should be empty slice, not nil")
		require.NotNil(t, got.PrivateChainAddresses, "PrivateChainAddresses should be empty slice, not nil")
		require.Empty(t, got.PrivateHubAddresses)
		require.Empty(t, got.PrivateHubDvpOperatorAddresses)
		require.Empty(t, got.PrivateChainAddresses)
	})
}

func TestGetPublicRelayerRaylsSignPublicAddresses(t *testing.T) {
	publicKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	chainKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	wantPublicAddr := crypto.PubkeyToAddress(publicKey.PublicKey)
	wantChainAddr := crypto.PubkeyToAddress(chainKey.PublicKey)

	t.Run("returns addresses derived from each key bucket", func(t *testing.T) {
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()

		plain := domain.PublicRelayerRaylsSignKeys{
			PublicChainKeys:  domain.RaylsSignKeyList{publicKey},
			PrivateChainKeys: domain.RaylsSignKeyList{chainKey},
		}
		encrypted, encryptErr := plain.Encrypt(encryptorMock)
		require.NoError(t, encryptErr)

		ecdsaRepoMock.GetPublicRelayerRaylsSignKeysFunc = func(ctx context.Context) (domain.EncryptedPublicRelayerRaylsSignKeys, error) {
			return encrypted, nil
		}

		svc := service.NewKeysService(encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock)

		got, err := svc.GetPublicRelayerRaylsSignPublicAddresses(context.Background())
		require.NoError(t, err)
		require.Equal(t, domain.AddressList{wantPublicAddr}, got.PublicChainAddresses)
		require.Equal(t, domain.AddressList{wantChainAddr}, got.PrivateChainAddresses)
	})

	t.Run("returns ErrNoRaylsSignKeys when the repo has no keys", func(t *testing.T) {
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		ecdsaRepoMock.GetPublicRelayerRaylsSignKeysFunc = func(ctx context.Context) (domain.EncryptedPublicRelayerRaylsSignKeys, error) {
			return domain.EncryptedPublicRelayerRaylsSignKeys{}, service.ErrNoRaylsSignKeys
		}

		svc := service.NewKeysService(encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock)

		got, err := svc.GetPublicRelayerRaylsSignPublicAddresses(context.Background())
		require.ErrorIs(t, err, service.ErrNoRaylsSignKeys)
		require.Equal(t, domain.PublicRelayerRaylsSignPublicAddresses{}, got)
	})

	t.Run("wraps non-ErrNoRaylsSignKeys repo errors as KeyRepositoryError", func(t *testing.T) {
		// Exercises the } else { branch in GetPublicRelayerRaylsSignPublicAddresses —
		// must continue to wrap as KeyRepositoryError after the flatten refactor.
		wantMsg := "db connection lost"
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		ecdsaRepoMock.GetPublicRelayerRaylsSignKeysFunc = func(ctx context.Context) (domain.EncryptedPublicRelayerRaylsSignKeys, error) {
			return domain.EncryptedPublicRelayerRaylsSignKeys{}, errors.New(wantMsg)
		}

		svc := service.NewKeysService(encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock)

		got, err := svc.GetPublicRelayerRaylsSignPublicAddresses(context.Background())
		require.Error(t, err)
		require.NotErrorIs(t, err, service.ErrNoRaylsSignKeys)
		var keyRepoErr service.KeyRepositoryError
		require.IsType(t, &keyRepoErr, err)
		require.Equal(t, wantMsg, err.Error())
		require.Equal(t, domain.PublicRelayerRaylsSignPublicAddresses{}, got)
	})
}

func TestGetAtomicServiceRaylsSignPublicAddresses(t *testing.T) {
	hubKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	chainKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	wantHubAddr := crypto.PubkeyToAddress(hubKey.PublicKey)
	wantChainAddr := crypto.PubkeyToAddress(chainKey.PublicKey)

	t.Run("returns addresses derived from each key bucket", func(t *testing.T) {
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()

		plain := domain.AtomicServiceRaylsSignKeys{
			PrivateHubKeys:   domain.RaylsSignKeyList{hubKey},
			PrivateChainKeys: domain.RaylsSignKeyList{chainKey},
		}
		encrypted, encryptErr := plain.Encrypt(encryptorMock)
		require.NoError(t, encryptErr)

		ecdsaRepoMock.GetAtomicServiceRaylsSignKeysFunc = func(ctx context.Context) (domain.EncryptedAtomicServiceRaylsSignKeys, error) {
			return encrypted, nil
		}

		svc := service.NewKeysService(encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock)

		got, err := svc.GetAtomicServiceRaylsSignPublicAddresses(context.Background())
		require.NoError(t, err)
		require.Equal(t, domain.AddressList{wantHubAddr}, got.PrivateHubAddresses)
		require.Equal(t, domain.AddressList{wantChainAddr}, got.PrivateChainAddresses)
	})

	t.Run("returns ErrNoRaylsSignKeys when the repo has no keys", func(t *testing.T) {
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		ecdsaRepoMock.GetAtomicServiceRaylsSignKeysFunc = func(ctx context.Context) (domain.EncryptedAtomicServiceRaylsSignKeys, error) {
			return domain.EncryptedAtomicServiceRaylsSignKeys{}, service.ErrNoRaylsSignKeys
		}

		svc := service.NewKeysService(encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock)

		got, err := svc.GetAtomicServiceRaylsSignPublicAddresses(context.Background())
		require.ErrorIs(t, err, service.ErrNoRaylsSignKeys)
		require.Equal(t, domain.AtomicServiceRaylsSignPublicAddresses{}, got)
	})

	t.Run("wraps non-ErrNoRaylsSignKeys repo errors as KeyRepositoryError", func(t *testing.T) {
		// Exercises the } else { branch in GetAtomicServiceRaylsSignPublicAddresses —
		// must continue to wrap as KeyRepositoryError after the flatten refactor.
		wantMsg := "db connection lost"
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		ecdsaRepoMock.GetAtomicServiceRaylsSignKeysFunc = func(ctx context.Context) (domain.EncryptedAtomicServiceRaylsSignKeys, error) {
			return domain.EncryptedAtomicServiceRaylsSignKeys{}, errors.New(wantMsg)
		}

		svc := service.NewKeysService(encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock)

		got, err := svc.GetAtomicServiceRaylsSignPublicAddresses(context.Background())
		require.Error(t, err)
		require.NotErrorIs(t, err, service.ErrNoRaylsSignKeys)
		var keyRepoErr service.KeyRepositoryError
		require.IsType(t, &keyRepoErr, err)
		require.Equal(t, wantMsg, err.Error())
		require.Equal(t, domain.AtomicServiceRaylsSignPublicAddresses{}, got)
	})
}

func TestCreatePaymentSpendKey(t *testing.T) {
	t.Run("creates enygma key and returns it", func(t *testing.T) {
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		ecdsaRepoMock.CreatePrivateRelayerRaylsSignKeysFunc = func(ctx context.Context, keys domain.EncryptedPrivateRelayerRaylsSignKeys) error { return nil }
		ecdsaRepoMock.GetPrivateRelayerRaylsSignKeysFunc = func(ctx context.Context) (domain.EncryptedPrivateRelayerRaylsSignKeys, error) {
			// Return some encrypted keys
			return domain.EncryptedPrivateRelayerRaylsSignKeys{}, nil
		}
		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		chainId := big.NewInt(1)
		key, err := svc.CreatePaymentSpendKey(context.Background(), chainId)

		require.Nil(t, err)
		require.NotEmpty(t, key)
	})
}

func TestGetPaymentSpendKey(t *testing.T) {
	t.Run("returns enygma key", func(t *testing.T) {
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()

		// Setup Enygma repository to return encrypted keys
		paymentSpendRepoMock.GetFunc = func(ctx context.Context) (domain.EncryptedPaymentSpendKeys, error) {
			publicKey := big.NewInt(12345)
			return domain.EncryptedPaymentSpendKeys{
				EncryptedSecretKey: []byte("encrypted_secret"),
				PublicKey:          publicKey.Bytes(),
			}, nil
		}

		// Setup encryptor to return valid secret key
		encryptorMock.DecryptFunc = func(bytes []byte) ([]byte, error) {
			validSecretKey := make([]byte, 32)
			validSecretKey[0] = 1 // Non-zero key
			return validSecretKey, nil
		}

		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		_, err := svc.CreatePaymentSpendKey(context.Background(), big.NewInt(1))
		require.Nil(t, err)

		key, err := svc.GetPaymentSpendKey(context.Background())

		require.Nil(t, err)
		require.NotEmpty(t, key)
	})

	t.Run("returns ErrNoPaymentSpendKeySet", func(t *testing.T) {
		wantErr := service.ErrNoPaymentSpendKeySet

		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		ecdsaRepoMock.CreatePrivateRelayerRaylsSignKeysFunc = func(ctx context.Context, keys domain.EncryptedPrivateRelayerRaylsSignKeys) error { return nil }
		ecdsaRepoMock.GetPrivateRelayerRaylsSignKeysFunc = func(ctx context.Context) (domain.EncryptedPrivateRelayerRaylsSignKeys, error) {
			// Return some encrypted keys
			return domain.EncryptedPrivateRelayerRaylsSignKeys{}, nil
		}
		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		key, err := svc.GetPaymentSpendKey(context.Background())

		require.Nil(t, key.SecretKey)
		require.Nil(t, key.PublicKey)
		require.NotNil(t, err)
		require.ErrorIs(t, err, wantErr)
	})
}

func TestGenerateKeyAgreement(t *testing.T) {
	t.Run("returns ciphertext, shared secret, and digest", func(t *testing.T) {
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		// Generate a test ML-KEM key pair
		dk, err := mlkem.GenerateKey768()
		require.Nil(t, err)
		publicKeyBytes := dk.EncapsulationKey().Bytes()

		chainID := big.NewInt(1)
		ciphertext, sharedSecret, digest, err := svc.GenerateKeyAgreement(chainID, publicKeyBytes)

		require.Nil(t, err)
		require.NotEmpty(t, ciphertext, "ciphertext should not be empty")
		require.NotEmpty(t, sharedSecret, "shared secret should not be empty")
		require.NotEmpty(t, digest, "digest should not be empty")
	})

	t.Run("returns error on invalid public key", func(t *testing.T) {
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		invalidPublicKey := []byte("invalid public key")
		chainID := big.NewInt(1)

		_, _, _, err := svc.GenerateKeyAgreement(chainID, invalidPublicKey)

		require.NotNil(t, err)
	})
}

func TestCreateKeyAgreement(t *testing.T) {
	t.Run("encrypts and persists shared secret", func(t *testing.T) {
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		chainID := big.NewInt(1)
		sharedSecret := []byte("test shared secret 32 bytes long")
		blockNum := uint64(100)

		err := svc.CreateKeyAgreement(context.Background(), chainID, sharedSecret, blockNum)

		require.Nil(t, err)
		require.Len(t, encryptorMock.EncryptCalls(), 1, "should encrypt shared secret")
		require.Len(t, sharedSecretsRepoMock.CreateCalls(), 1, "should persist shared secret")
	})

	t.Run("returns error from encryptor", func(t *testing.T) {
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		encryptorMock.EncryptFunc = func(bytes []byte) ([]byte, error) {
			return nil, errors.New("encryption failed")
		}
		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		chainID := big.NewInt(1)
		sharedSecret := []byte("test shared secret")
		blockNum := uint64(100)

		err := svc.CreateKeyAgreement(context.Background(), chainID, sharedSecret, blockNum)

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "failed to encrypt shared secret")
	})

	t.Run("returns error from repository", func(t *testing.T) {
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		sharedSecretsRepoMock.CreateFunc = func(ctx context.Context, secret domain.EncryptedSharedSecret) error {
			return errors.New("repository error")
		}
		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		chainID := big.NewInt(1)
		sharedSecret := []byte("test shared secret")
		blockNum := uint64(100)

		err := svc.CreateKeyAgreement(context.Background(), chainID, sharedSecret, blockNum)

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "failed to create shared secret")
	})
}

func TestCompleteKeyAgreement(t *testing.T) {
	t.Run("decapsulates ciphertext and stores shared secret", func(t *testing.T) {
		// Generate a test ML-KEM key pair and ciphertext
		dk, err := mlkem.GenerateKey768()
		require.Nil(t, err)
		publicKeyBytes := dk.EncapsulationKey().Bytes()

		// Generate ciphertext using the public key
		ek, err := mlkem.NewEncapsulationKey768(publicKeyBytes)
		require.Nil(t, err)
		_, ciphertext := ek.Encapsulate()

		// Setup mocks
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()

		// Return the valid ML-KEM key pair from the repository
		viewRepoMock.GetForBlockNumberFunc = func(ctx context.Context, blockNumber uint64) (domain.EncryptedRaylsViewKeyPair, error) {
			return domain.EncryptedRaylsViewKeyPair{
				InitialBlock:                 100,
				EncryptedRaylsViewPrivateKey: dk.Bytes(),
				RaylsViewPublicKey:           publicKeyBytes,
			}, nil
		}
		// Return the actual key bytes for decryption
		encryptorMock.DecryptFunc = func(bytes []byte) ([]byte, error) {
			return bytes, nil // Return as-is since we're passing the actual key bytes
		}

		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		chainID := big.NewInt(1)
		blockNum := uint64(100)

		err = svc.CompleteKeyAgreement(context.Background(), chainID, ciphertext, blockNum)

		require.Nil(t, err)
		require.Len(t, viewRepoMock.GetForBlockNumberCalls(), 1, "should get view key pair")
		require.Len(t, sharedSecretsRepoMock.CreateCalls(), 1, "should persist shared secret")
	})

	t.Run("returns error when view key pair not found", func(t *testing.T) {
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		viewRepoMock.GetForBlockNumberFunc = func(ctx context.Context, blockNumber uint64) (domain.EncryptedRaylsViewKeyPair, error) {
			return domain.EncryptedRaylsViewKeyPair{}, service.ErrNoRaylsViewKeysSet
		}
		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		chainID := big.NewInt(1)
		ciphertext := []byte("dummy ciphertext")
		blockNum := uint64(100)

		err := svc.CompleteKeyAgreement(context.Background(), chainID, ciphertext, blockNum)

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "failed to get view key pair")
	})

	t.Run("returns error from repository", func(t *testing.T) {
		// Generate a test ML-KEM key pair and ciphertext
		dk, err := mlkem.GenerateKey768()
		require.Nil(t, err)
		publicKeyBytes := dk.EncapsulationKey().Bytes()

		// Generate ciphertext using the public key
		ek, err := mlkem.NewEncapsulationKey768(publicKeyBytes)
		require.Nil(t, err)
		_, ciphertext := ek.Encapsulate()

		// Setup mocks
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()

		viewRepoMock.GetForBlockNumberFunc = func(ctx context.Context, blockNumber uint64) (domain.EncryptedRaylsViewKeyPair, error) {
			return domain.EncryptedRaylsViewKeyPair{
				InitialBlock:                 100,
				EncryptedRaylsViewPrivateKey: dk.Bytes(),
				RaylsViewPublicKey:           publicKeyBytes,
			}, nil
		}
		encryptorMock.DecryptFunc = func(bytes []byte) ([]byte, error) {
			return bytes, nil
		}
		sharedSecretsRepoMock.CreateFunc = func(ctx context.Context, secret domain.EncryptedSharedSecret) error {
			return errors.New("repository error")
		}

		svc := service.NewKeysService(
			encryptorMock,
			viewRepoMock,
			ecdsaRepoMock,
			sharedSecretsRepoMock,
			enygmaSelfSecretsRepoMock,
			paymentSpendRepoMock,
		)

		chainID := big.NewInt(1)
		blockNum := uint64(100)

		err = svc.CompleteKeyAgreement(context.Background(), chainID, ciphertext, blockNum)

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "failed to create shared secret")
	})
}

func TestRecoverViewSalt(t *testing.T) {
	t.Run("recovers the same salt the sender derived via GenerateSalt", func(t *testing.T) {
		// Sender (relayer) generated a salt by encapsulating against this view pubkey.
		dk, err := mlkem.GenerateKey768()
		require.Nil(t, err)
		privateKeyBytes := dk.Bytes()
		publicKeyBytes := dk.EncapsulationKey().Bytes()

		wantSalt, ctxt, err := cryptography.GenerateSalt(publicKeyBytes)
		require.Nil(t, err)

		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		encryptorMock.DecryptFunc = func([]byte) ([]byte, error) {
			return privateKeyBytes, nil
		}
		viewRepoMock.GetForBlockNumberFunc = func(_ context.Context, _ uint64) (domain.EncryptedRaylsViewKeyPair, error) {
			return domain.EncryptedRaylsViewKeyPair{
				InitialBlock:                 100,
				EncryptedRaylsViewPrivateKey: []byte("encrypted_private"),
				RaylsViewPublicKey:           publicKeyBytes,
			}, nil
		}

		svc := service.NewKeysService(encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock)

		gotSalt, err := svc.RecoverViewSalt(context.Background(), 150, ctxt)

		require.Nil(t, err)
		require.Equal(t, 0, wantSalt.Cmp(gotSalt), "recovered salt must equal the salt the sender generated")
	})

	t.Run("returns repository error when no view keys exist for block", func(t *testing.T) {
		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		// createDefaultMocks already returns ErrNoRaylsViewKeysSet from the view repo.

		svc := service.NewKeysService(encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock)

		_, err := svc.RecoverViewSalt(context.Background(), 150, []byte("anything"))

		require.NotNil(t, err)
		require.ErrorIs(t, err, service.ErrNoRaylsViewKeysSet)
	})

	t.Run("returns error when ctxt cannot be decapsulated", func(t *testing.T) {
		dk, err := mlkem.GenerateKey768()
		require.Nil(t, err)
		privateKeyBytes := dk.Bytes()
		publicKeyBytes := dk.EncapsulationKey().Bytes()

		encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock := createDefaultMocks()
		encryptorMock.DecryptFunc = func([]byte) ([]byte, error) { return privateKeyBytes, nil }
		viewRepoMock.GetForBlockNumberFunc = func(_ context.Context, _ uint64) (domain.EncryptedRaylsViewKeyPair, error) {
			return domain.EncryptedRaylsViewKeyPair{
				InitialBlock:                 100,
				EncryptedRaylsViewPrivateKey: []byte("encrypted_private"),
				RaylsViewPublicKey:           publicKeyBytes,
			}, nil
		}

		svc := service.NewKeysService(encryptorMock, viewRepoMock, ecdsaRepoMock, sharedSecretsRepoMock, enygmaSelfSecretsRepoMock, paymentSpendRepoMock)

		_, err = svc.RecoverViewSalt(context.Background(), 150, []byte("not-a-valid-mlkem-ciphertext"))

		require.NotNil(t, err)
		require.Contains(t, err.Error(), "ml-kem decapsulation")
	})
}
