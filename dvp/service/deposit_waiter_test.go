package service_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp/testdata"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDepositWaiter_WaitUntilDepositIsConfirmed(t *testing.T) {
	t.Run("returns immediately when deposit is already unspent", func(t *testing.T) {
		config := service.WaitConfig{
			MaxRetries:    10,
			RetryInterval: 0, // 0 for testing, no actual sleep
		}

		deposit := testdata.NewDvpDeposit()
		deposit.Status = types.DvpDepositUnspent

		repoMock := &consolidationDepositRepositoryMock{
			GetDepositByCommitmentFunc: func(ctx context.Context, commitment *big.Int) (*types.DvpDeposit, error) {
				// Should not be called if deposit is already unspent
				t.Fatal("GetDepositByCommitment should not be called")
				return nil, nil //nolint:nilnil // intentional nil return in test mock
			},
		}

		waiter := service.NewDepositWaiter(config, repoMock)
		err := waiter.WaitUntilDepositIsConfirmed(context.Background(), deposit)

		require.NoError(t, err)
	})

	t.Run("returns nil after deposit becomes unspent after N retries", func(t *testing.T) {
		tests := []struct {
			name          string
			retryCount    int
			expectedCalls int
		}{
			{"1 retry", 1, 1},
			{"5 retries", 5, 5},
			{"9 retries (boundary)", 9, 9},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				config := service.WaitConfig{
					MaxRetries:    10,
					RetryInterval: 0,
				}

				deposit := testdata.NewDvpDeposit()
				deposit.Status = types.DvpDepositPending

				callCount := 0
				repoMock := &consolidationDepositRepositoryMock{
					GetDepositByCommitmentFunc: func(ctx context.Context, commitment *big.Int) (*types.DvpDeposit, error) {
						callCount++
						updatedDeposit := testdata.NewDvpDeposit()
						updatedDeposit.Commitment = commitment

						// Return Pending until retryCount is reached, then Unspent
						if callCount < tt.retryCount {
							updatedDeposit.Status = types.DvpDepositPending
						} else {
							updatedDeposit.Status = types.DvpDepositUnspent
						}
						return updatedDeposit, nil
					},
				}

				waiter := service.NewDepositWaiter(config, repoMock)
				err := waiter.WaitUntilDepositIsConfirmed(context.Background(), deposit)

				require.NoError(t, err)
				assert.Equal(
					t,
					tt.expectedCalls,
					callCount,
					"should call GetDepositByCommitment correct number of times",
				)
			})
		}
	})

	t.Run("returns error when max retries exceeded", func(t *testing.T) {
		config := service.WaitConfig{
			MaxRetries:    10,
			RetryInterval: 0,
		}

		deposit := testdata.NewDvpDeposit()
		deposit.Status = types.DvpDepositPending

		callCount := 0
		repoMock := &consolidationDepositRepositoryMock{
			GetDepositByCommitmentFunc: func(ctx context.Context, commitment *big.Int) (*types.DvpDeposit, error) {
				callCount++
				updatedDeposit := testdata.NewDvpDeposit()
				updatedDeposit.Status = types.DvpDepositPending // Never becomes Unspent
				updatedDeposit.Commitment = commitment
				return updatedDeposit, nil
			},
		}

		waiter := service.NewDepositWaiter(config, repoMock)
		err := waiter.WaitUntilDepositIsConfirmed(context.Background(), deposit)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "deposit not ready after 10 retries")
		assert.Equal(
			t,
			9,
			callCount,
			"should call GetDepositByCommitment 9 times before hitting max retries on 10th iteration",
		)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		tests := []struct {
			name          string
			failAtCall    int
			expectedCalls int
		}{
			{"fails immediately", 1, 1},
			{"fails on 3rd call", 3, 3},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				config := service.WaitConfig{
					MaxRetries:    10,
					RetryInterval: 0,
				}

				deposit := testdata.NewDvpDeposit()
				deposit.Status = types.DvpDepositPending

				callCount := 0
				repoMock := &consolidationDepositRepositoryMock{
					GetDepositByCommitmentFunc: func(ctx context.Context, commitment *big.Int) (*types.DvpDeposit, error) {
						callCount++
						// Fail at specified call
						if callCount == tt.failAtCall {
							return nil, errors.New("database connection lost")
						}
						updatedDeposit := testdata.NewDvpDeposit()
						updatedDeposit.Status = types.DvpDepositPending
						updatedDeposit.Commitment = commitment
						return updatedDeposit, nil
					},
				}

				waiter := service.NewDepositWaiter(config, repoMock)
				err := waiter.WaitUntilDepositIsConfirmed(context.Background(), deposit)

				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed to fetch updated deposit status")
				assert.Contains(t, err.Error(), "database connection lost")
				assert.Equal(
					t,
					tt.expectedCalls,
					callCount,
					"should call GetDepositByCommitment expected number of times",
				)
			})
		}
	})

	t.Run("handles deposit with locked status", func(t *testing.T) {
		config := service.WaitConfig{
			MaxRetries:    10,
			RetryInterval: 0,
		}

		deposit := testdata.NewDvpDeposit()
		deposit.Status = types.DvpDepositLocked

		callCount := 0
		repoMock := &consolidationDepositRepositoryMock{
			GetDepositByCommitmentFunc: func(ctx context.Context, commitment *big.Int) (*types.DvpDeposit, error) {
				callCount++
				updatedDeposit := testdata.NewDvpDeposit()
				updatedDeposit.Commitment = commitment

				// Locked -> Pending -> Unspent
				switch callCount {
				case 1:
					updatedDeposit.Status = types.DvpDepositLocked
				case 2:
					updatedDeposit.Status = types.DvpDepositPending
				default:
					updatedDeposit.Status = types.DvpDepositUnspent
				}
				return updatedDeposit, nil
			},
		}

		waiter := service.NewDepositWaiter(config, repoMock)
		err := waiter.WaitUntilDepositIsConfirmed(context.Background(), deposit)

		require.NoError(t, err)
		assert.Equal(t, 3, callCount, "should call GetDepositByCommitment 3 times for status transitions")
	})
}
