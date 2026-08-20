package service

import (
	"context"
	"errors"
	"testing"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type swapWaiterRepositoryMock struct {
	GetSwapBySharedIDFunc func(ctx context.Context, sharedId string) (*types.DvpSwap, error)
	calls                 []string
}

func (m *swapWaiterRepositoryMock) GetSwapBySharedID(ctx context.Context, sharedId string) (*types.DvpSwap, error) {
	m.calls = append(m.calls, "GetSwapBySharedID")
	if m.GetSwapBySharedIDFunc != nil {
		return m.GetSwapBySharedIDFunc(ctx, sharedId)
	}
	return nil, nil //nolint:nilnil // intentional nil return in test mock
}

func TestSwapWaiter(t *testing.T) {
	sharedId := "test-swap-123"

	t.Run("successfully returns swap when found on first poll", func(t *testing.T) {
		swap := &types.DvpSwap{
			SharedID: sharedId,
		}

		repoMock := &swapWaiterRepositoryMock{
			GetSwapBySharedIDFunc: func(ctx context.Context, id string) (*types.DvpSwap, error) {
				return swap, nil
			},
		}

		waiter := NewSwapWaiter(
			WaitConfig{
				MaxRetries:    10,
				RetryInterval: 0,
			},
			repoMock,
		)

		result, err := waiter.WaitForSwapInitiation(context.Background(), sharedId)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, sharedId, result.SharedID)
		assert.Len(t, repoMock.calls, 1)
	})

	t.Run("successfully returns swap when found after multiple polls", func(t *testing.T) {
		swap := &types.DvpSwap{
			SharedID: sharedId,
		}

		callCount := 0
		repoMock := &swapWaiterRepositoryMock{
			GetSwapBySharedIDFunc: func(ctx context.Context, id string) (*types.DvpSwap, error) {
				callCount++
				if callCount < 3 {
					return nil, nil //nolint:nilnil // intentional nil return in test mock
				}
				return swap, nil
			},
		}

		waiter := NewSwapWaiter(
			WaitConfig{
				MaxRetries:    10,
				RetryInterval: 0,
			},
			repoMock,
		)

		result, err := waiter.WaitForSwapInitiation(context.Background(), sharedId)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, sharedId, result.SharedID)
		assert.Len(t, repoMock.calls, 3)
	})

	t.Run("returns error when swap is never found after max retries", func(t *testing.T) {
		repoMock := &swapWaiterRepositoryMock{
			GetSwapBySharedIDFunc: func(ctx context.Context, id string) (*types.DvpSwap, error) {
				return nil, nil //nolint:nilnil // intentional nil return in test mock
			},
		}

		waiter := NewSwapWaiter(
			WaitConfig{
				MaxRetries:    3,
				RetryInterval: 0,
			},
			repoMock,
		)

		result, err := waiter.WaitForSwapInitiation(context.Background(), sharedId)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Equal(t, "swap was initiated by the other side but not found in our DB after 3 retries", err.Error())
		assert.Len(t, repoMock.calls, 3)
	})

	t.Run("returns error from repository immediately", func(t *testing.T) {
		repoMock := &swapWaiterRepositoryMock{
			GetSwapBySharedIDFunc: func(ctx context.Context, id string) (*types.DvpSwap, error) {
				return nil, errors.New("database connection error")
			},
		}

		waiter := NewSwapWaiter(
			WaitConfig{
				MaxRetries:    10,
				RetryInterval: 0,
			},
			repoMock,
		)

		result, err := waiter.WaitForSwapInitiation(context.Background(), sharedId)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Equal(t, "failed to fetch swap from DB: database connection error", err.Error())
		assert.Len(t, repoMock.calls, 1)
	})

	t.Run("respects max retries configuration", func(t *testing.T) {
		callCount := 0
		repoMock := &swapWaiterRepositoryMock{
			GetSwapBySharedIDFunc: func(ctx context.Context, id string) (*types.DvpSwap, error) {
				callCount++
				//nolint:nilnil // returning nil,nil is intentional for "not found" pattern in test
				return nil, nil
			},
		}

		waiter := NewSwapWaiter(
			WaitConfig{
				MaxRetries:    5,
				RetryInterval: 0,
			},
			repoMock,
		)

		result, err := waiter.WaitForSwapInitiation(context.Background(), sharedId)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Len(t, repoMock.calls, 5)
		assert.Equal(t, "swap was initiated by the other side but not found in our DB after 5 retries", err.Error())
	})
}
