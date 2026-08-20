package service_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/dvp/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDepositFinder_FindERC1155DepositsForJSProof(t *testing.T) {
	t.Run("returns all deposits when all are needed to meet requirement", func(t *testing.T) {
		deposits := []types.DvpDeposit{
			createTestDeposit("user1", "0x1234", "100"),
			createTestDeposit("user1", "0x1234", "200"),
		}

		repo := &finderDepositRepositoryMock{
			GetDepositsByTokenFunc: func(ctx context.Context, tokenAddress string, tokenId string, tokenType types.DvpTokenType, userAddress string, status types.DvpDepositStatus) ([]types.DvpDeposit, error) {
				return deposits, nil
			},
		}

		finder := service.NewDepositFinder(repo)
		result, err := finder.FindERC1155DepositsForJSProof(context.Background(), "user1", "0x1234", "42", big.NewInt(150))

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, "100", result[0].TokenAmount.String())
		assert.Equal(t, "200", result[1].TokenAmount.String())
	})

	t.Run("returns only deposits needed and stops before processing excess", func(t *testing.T) {
		deposits := []types.DvpDeposit{
			createTestDeposit("user1", "0x1234", "100"),
			createTestDeposit("user1", "0x1234", "200"),
			createTestDeposit("user1", "0x1234", "300"),
			createTestDeposit("user1", "0x1234", "400"),
		}

		repo := &finderDepositRepositoryMock{
			GetDepositsByTokenFunc: func(ctx context.Context, tokenAddress string, tokenId string, tokenType types.DvpTokenType, userAddress string, status types.DvpDepositStatus) ([]types.DvpDeposit, error) {
				return deposits, nil
			},
		}

		finder := service.NewDepositFinder(repo)
		result, err := finder.FindERC1155DepositsForJSProof(context.Background(), "user1", "0x1234", "42", big.NewInt(250))

		require.NoError(t, err)
		// Should return only 100 + 200 = 300 >= 250, not the 300 and 400
		assert.Len(t, result, 2)
		assert.Equal(t, "100", result[0].TokenAmount.String())
		assert.Equal(t, "200", result[1].TokenAmount.String())
	})

	t.Run("returns error when insufficient deposits", func(t *testing.T) {
		deposits := []types.DvpDeposit{
			createTestDeposit("user1", "0x1234", "100"),
		}

		repo := &finderDepositRepositoryMock{
			GetDepositsByTokenFunc: func(ctx context.Context, tokenAddress string, tokenId string, tokenType types.DvpTokenType, userAddress string, status types.DvpDepositStatus) ([]types.DvpDeposit, error) {
				return deposits, nil
			},
		}

		finder := service.NewDepositFinder(repo)
		result, err := finder.FindERC1155DepositsForJSProof(context.Background(), "user1", "0x1234", "42", big.NewInt(500))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not enough deposits")
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repo := &finderDepositRepositoryMock{
			GetDepositsByTokenFunc: func(ctx context.Context, tokenAddress string, tokenId string, tokenType types.DvpTokenType, userAddress string, status types.DvpDepositStatus) ([]types.DvpDeposit, error) {
				return nil, errors.New("database error")
			},
		}

		finder := service.NewDepositFinder(repo)
		result, err := finder.FindERC1155DepositsForJSProof(context.Background(), "user1", "0x1234", "42", big.NewInt(100))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "database error")
	})

	t.Run("returns error when no deposits available", func(t *testing.T) {
		repo := &finderDepositRepositoryMock{
			GetDepositsByTokenFunc: func(ctx context.Context, tokenAddress string, tokenId string, tokenType types.DvpTokenType, userAddress string, status types.DvpDepositStatus) ([]types.DvpDeposit, error) {
				return []types.DvpDeposit{}, nil
			},
		}

		finder := service.NewDepositFinder(repo)
		result, err := finder.FindERC1155DepositsForJSProof(context.Background(), "user1", "0x1234", "42", big.NewInt(100))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not enough deposits")
	})

	t.Run("verifies correct parameters passed to repository", func(t *testing.T) {
		var capturedParams struct {
			tokenAddress string
			tokenId      string
			tokenType    types.DvpTokenType
			userAddress  string
			status       types.DvpDepositStatus
		}

		repo := &finderDepositRepositoryMock{
			GetDepositsByTokenFunc: func(ctx context.Context, tokenAddress string, tokenId string, tokenType types.DvpTokenType, userAddress string, status types.DvpDepositStatus) ([]types.DvpDeposit, error) {
				capturedParams.tokenAddress = tokenAddress
				capturedParams.tokenId = tokenId
				capturedParams.tokenType = tokenType
				capturedParams.userAddress = userAddress
				capturedParams.status = status
				return []types.DvpDeposit{
					createTestDeposit("user1", "0x1234", "100"),
				}, nil
			},
		}

		finder := service.NewDepositFinder(repo)
		_, _ = finder.FindERC1155DepositsForJSProof(context.Background(), "user1", "0x1234", "42", big.NewInt(100))

		assert.Equal(t, "0x1234", capturedParams.tokenAddress)
		assert.Equal(t, "42", capturedParams.tokenId)
		assert.Equal(t, types.DvpERC1155, capturedParams.tokenType)
		assert.Equal(t, "user1", capturedParams.userAddress)
		assert.Equal(t, types.DvpDepositUnspent, capturedParams.status)
	})

	t.Run("handles large payment amounts", func(t *testing.T) {
		largeAmount := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
		deposits := []types.DvpDeposit{
			createTestDeposit("user1", "0x1234", largeAmount.String()),
		}

		repo := &finderDepositRepositoryMock{
			GetDepositsByTokenFunc: func(ctx context.Context, tokenAddress string, tokenId string, tokenType types.DvpTokenType, userAddress string, status types.DvpDepositStatus) ([]types.DvpDeposit, error) {
				return deposits, nil
			},
		}

		finder := service.NewDepositFinder(repo)
		result, err := finder.FindERC1155DepositsForJSProof(context.Background(), "user1", "0x1234", "42", largeAmount)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result, 1)
	})

	t.Run("handles zero payment amount", func(t *testing.T) {
		deposits := []types.DvpDeposit{
			createTestDeposit("user1", "0x1234", "100"),
		}

		repo := &finderDepositRepositoryMock{
			GetDepositsByTokenFunc: func(ctx context.Context, tokenAddress string, tokenId string, tokenType types.DvpTokenType, userAddress string, status types.DvpDepositStatus) ([]types.DvpDeposit, error) {
				return deposits, nil
			},
		}

		finder := service.NewDepositFinder(repo)
		result, err := finder.FindERC1155DepositsForJSProof(context.Background(), "user1", "0x1234", "42", big.NewInt(0))

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result)
	})
}

func TestDepositFinder_FindEnygmaDeposits(t *testing.T) {
	t.Run("returns all deposits when all are needed to meet requirement", func(t *testing.T) {
		deposits := []types.DvpDeposit{
			createTestEnygmaDeposit("user1", "0x5678", "100"),
			createTestEnygmaDeposit("user1", "0x5678", "200"),
		}

		repo := &finderDepositRepositoryMock{
			GetFungibleDepositsFunc: func(ctx context.Context, tokenAddress string, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) ([]types.DvpDeposit, error) {
				return deposits, nil
			},
		}

		finder := service.NewDepositFinder(repo)
		result, err := finder.FindEnygmaDeposits(context.Background(), "user1", "0x5678", big.NewInt(150))

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, "100", result[0].TokenAmount.String())
		assert.Equal(t, "200", result[1].TokenAmount.String())
	})

	t.Run("returns only deposits needed and stops before processing excess", func(t *testing.T) {
		deposits := []types.DvpDeposit{
			createTestEnygmaDeposit("user1", "0x5678", "100"),
			createTestEnygmaDeposit("user1", "0x5678", "200"),
			createTestEnygmaDeposit("user1", "0x5678", "300"),
			createTestEnygmaDeposit("user1", "0x5678", "400"),
		}

		repo := &finderDepositRepositoryMock{
			GetFungibleDepositsFunc: func(ctx context.Context, tokenAddress string, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) ([]types.DvpDeposit, error) {
				return deposits, nil
			},
		}

		finder := service.NewDepositFinder(repo)
		result, err := finder.FindEnygmaDeposits(context.Background(), "user1", "0x5678", big.NewInt(250))

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "100", result[0].TokenAmount.String())
		assert.Equal(t, "200", result[1].TokenAmount.String())
	})

	t.Run("returns error when insufficient deposits", func(t *testing.T) {
		deposits := []types.DvpDeposit{
			createTestEnygmaDeposit("user1", "0x5678", "100"),
		}

		repo := &finderDepositRepositoryMock{
			GetFungibleDepositsFunc: func(ctx context.Context, tokenAddress string, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) ([]types.DvpDeposit, error) {
				return deposits, nil
			},
		}

		finder := service.NewDepositFinder(repo)
		result, err := finder.FindEnygmaDeposits(context.Background(), "user1", "0x5678", big.NewInt(500))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not enough deposits")
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repo := &finderDepositRepositoryMock{
			GetFungibleDepositsFunc: func(ctx context.Context, tokenAddress string, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) ([]types.DvpDeposit, error) {
				return nil, errors.New("database error")
			},
		}

		finder := service.NewDepositFinder(repo)
		result, err := finder.FindEnygmaDeposits(context.Background(), "user1", "0x5678", big.NewInt(100))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "database error")
	})

	t.Run("returns error when no deposits available", func(t *testing.T) {
		repo := &finderDepositRepositoryMock{
			GetFungibleDepositsFunc: func(ctx context.Context, tokenAddress string, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) ([]types.DvpDeposit, error) {
				return []types.DvpDeposit{}, nil
			},
		}

		finder := service.NewDepositFinder(repo)
		result, err := finder.FindEnygmaDeposits(context.Background(), "user1", "0x5678", big.NewInt(100))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not enough deposits")
	})

	t.Run("verifies correct parameters passed to repository", func(t *testing.T) {
		var capturedParams struct {
			tokenAddress string
			userAddress  string
			tokenType    types.DvpTokenType
			status       types.DvpDepositStatus
		}

		repo := &finderDepositRepositoryMock{
			GetFungibleDepositsFunc: func(ctx context.Context, tokenAddress string, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) ([]types.DvpDeposit, error) {
				capturedParams.tokenAddress = tokenAddress
				capturedParams.userAddress = userAddress
				capturedParams.tokenType = tokenType
				capturedParams.status = status
				return []types.DvpDeposit{
					createTestEnygmaDeposit("user1", "0x5678", "100"),
				}, nil
			},
		}

		finder := service.NewDepositFinder(repo)
		_, _ = finder.FindEnygmaDeposits(context.Background(), "user1", "0x5678", big.NewInt(100))

		assert.Equal(t, "0x5678", capturedParams.tokenAddress)
		assert.Equal(t, "user1", capturedParams.userAddress)
		assert.Equal(t, types.DvpEnygma, capturedParams.tokenType)
		assert.Equal(t, types.DvpDepositUnspent, capturedParams.status)
	})

	t.Run("handles large payment amounts", func(t *testing.T) {
		largeAmount := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
		deposits := []types.DvpDeposit{
			createTestEnygmaDeposit("user1", "0x5678", largeAmount.String()),
		}

		repo := &finderDepositRepositoryMock{
			GetFungibleDepositsFunc: func(ctx context.Context, tokenAddress string, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) ([]types.DvpDeposit, error) {
				return deposits, nil
			},
		}

		finder := service.NewDepositFinder(repo)
		result, err := finder.FindEnygmaDeposits(context.Background(), "user1", "0x5678", largeAmount)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result, 1)
	})
}

func TestDepositFinder_FindERC721Deposit(t *testing.T) {
	t.Run("successfully finds ERC721 deposit", func(t *testing.T) {
		deposit := createTestERC721Deposit("user1", "0x1234", "42")

		repo := &finderDepositRepositoryMock{
			GetNonFungibleDepositFunc: func(ctx context.Context, tokenId string, tokenAddress string, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) (*types.DvpDeposit, error) {
				return &deposit, nil
			},
		}

		finder := service.NewDepositFinder(repo)
		result, err := finder.FindERC721Deposit(context.Background(), "user1", "0x1234", "42")

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "42", result.TokenID)
		assert.Equal(t, types.DvpERC721, result.TokenType)
		assert.Equal(t, "user1", result.UserAddress)
		assert.Equal(t, "0x1234", result.TokenAddress)
		assert.Equal(t, int64(1), result.TokenAmount.Int64())
	})

	t.Run("returns error when deposit not found (nil)", func(t *testing.T) {
		repo := &finderDepositRepositoryMock{
			GetNonFungibleDepositFunc: func(ctx context.Context, tokenId string, tokenAddress string, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) (*types.DvpDeposit, error) {
				return nil, nil //nolint:nilnil // intentional nil return in test mock
			},
		}

		finder := service.NewDepositFinder(repo)
		result, err := finder.FindERC721Deposit(context.Background(), "user1", "0x1234", "42")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repo := &finderDepositRepositoryMock{
			GetNonFungibleDepositFunc: func(ctx context.Context, tokenId string, tokenAddress string, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) (*types.DvpDeposit, error) {
				return nil, errors.New("database error")
			},
		}

		finder := service.NewDepositFinder(repo)
		result, err := finder.FindERC721Deposit(context.Background(), "user1", "0x1234", "42")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "database error")
	})

	t.Run("verifies correct parameters passed to repository", func(t *testing.T) {
		var capturedParams struct {
			tokenAddress string
			tokenId      string
			userAddress  string
			tokenType    types.DvpTokenType
			status       types.DvpDepositStatus
		}

		deposit := createTestERC721Deposit("user1", "0x1234", "42")
		repo := &finderDepositRepositoryMock{
			GetNonFungibleDepositFunc: func(ctx context.Context, tokenId string, tokenAddress string, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) (*types.DvpDeposit, error) {
				capturedParams.tokenAddress = tokenAddress
				capturedParams.tokenId = tokenId
				capturedParams.userAddress = userAddress
				capturedParams.tokenType = tokenType
				capturedParams.status = status
				return &deposit, nil
			},
		}

		finder := service.NewDepositFinder(repo)
		_, _ = finder.FindERC721Deposit(context.Background(), "user1", "0x1234", "42")

		assert.Equal(t, "0x1234", capturedParams.tokenAddress)
		assert.Equal(t, "42", capturedParams.tokenId)
		assert.Equal(t, "user1", capturedParams.userAddress)
		assert.Equal(t, types.DvpERC721, capturedParams.tokenType)
		assert.Equal(t, types.DvpDepositUnspent, capturedParams.status)
	})

	t.Run("handles different token IDs correctly", func(t *testing.T) {
		testCases := []struct {
			tokenId string
		}{
			{"1"},
			{"999"},
			{"12345"},
		}

		for _, tc := range testCases {
			t.Run("tokenId="+tc.tokenId, func(t *testing.T) {
				deposit := createTestERC721Deposit("user1", "0x1234", tc.tokenId)

				repo := &finderDepositRepositoryMock{
					GetNonFungibleDepositFunc: func(ctx context.Context, tokenId string, tokenAddress string, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) (*types.DvpDeposit, error) {
						if tokenId == tc.tokenId {
							return &deposit, nil
						}
						return nil, nil //nolint:nilnil // intentional nil return in test mock
					},
				}

				finder := service.NewDepositFinder(repo)
				result, err := finder.FindERC721Deposit(context.Background(), "user1", "0x1234", tc.tokenId)

				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tc.tokenId, result.TokenID)
			})
		}
	})
}

// Helper functions
func createTestDeposit(userAddress, tokenAddress, amount string) types.DvpDeposit {
	tokenAmount := new(big.Int)
	tokenAmount.SetString(amount, 10)

	return types.DvpDeposit{
		UserAddress:  userAddress,
		TokenAddress: tokenAddress,
		TokenAmount:  tokenAmount,
		TokenType:    types.DvpERC1155,
		TokenID:      "42",
		Status:       types.DvpDepositUnspent,
		CreatedAt:    time.Now(),
		Salt:         big.NewInt(123),
		TreeNumber:   1,
		Commitment:   big.NewInt(456),
		Nullifier:    big.NewInt(0),
	}
}

func createTestEnygmaDeposit(userAddress, tokenAddress, amount string) types.DvpDeposit {
	tokenAmount := new(big.Int)
	tokenAmount.SetString(amount, 10)

	return types.DvpDeposit{
		UserAddress:  userAddress,
		TokenAddress: tokenAddress,
		TokenAmount:  tokenAmount,
		TokenType:    types.DvpEnygma,
		TokenID:      "",
		Status:       types.DvpDepositUnspent,
		CreatedAt:    time.Now(),
		Salt:         big.NewInt(789),
		TreeNumber:   1,
		Commitment:   big.NewInt(789),
		Nullifier:    big.NewInt(0),
	}
}

func createTestERC721Deposit(userAddress, tokenAddress, tokenId string) types.DvpDeposit {
	return types.DvpDeposit{
		UserAddress:  userAddress,
		TokenAddress: tokenAddress,
		TokenAmount:  big.NewInt(1),
		TokenType:    types.DvpERC721,
		TokenID:      tokenId,
		Status:       types.DvpDepositUnspent,
		CreatedAt:    time.Now(),
		Salt:         big.NewInt(999),
		TreeNumber:   1,
		Commitment:   big.NewInt(888),
		Nullifier:    big.NewInt(0),
	}
}
