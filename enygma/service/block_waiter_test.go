package service_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/raylsnetwork/rayls-sovereign-relayer/enygma/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/enygma/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockBlockWaiterEthereumClient struct {
	BlockNumberFunc func(ctx context.Context) (uint64, error)

	BlockNumberCalls struct {
		sync.RWMutex
		Calls []struct {
			Ctx context.Context
		}
	}
}

func (m *MockBlockWaiterEthereumClient) BlockNumber(ctx context.Context) (uint64, error) {
	if m.BlockNumberFunc == nil {
		panic("MockBlockWaiterEthereumClient.BlockNumberFunc: method is nil but BlockNumber was just called")
	}

	m.BlockNumberCalls.Lock()
	m.BlockNumberCalls.Calls = append(m.BlockNumberCalls.Calls, struct {
		Ctx context.Context
	}{Ctx: ctx})
	m.BlockNumberCalls.Unlock()

	return m.BlockNumberFunc(ctx)
}

func TestBlockWaiterService_WaitForBlock(t *testing.T) {
	t.Run("returns immediately when target block is already reached", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		mockEthClient := &MockBlockWaiterEthereumClient{
			BlockNumberFunc: func(ctx context.Context) (uint64, error) {
				return 150, nil // Already at target
			},
		}

		service := service.NewBlockWaiterService(mockTracer, mockEthClient)
		ctx := context.Background()

		blockNum, err := service.WaitForBlock(ctx, 100)

		require.NoError(t, err)
		assert.Equal(t, uint64(150), blockNum)
		assert.Equal(t, 1, len(mockEthClient.BlockNumberCalls.Calls))
		assert.Equal(t, 1, len(mockTracer.StartCalls.Calls))
	})

	t.Run("returns the target block number after polling", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		callCount := 0
		mockEthClient := &MockBlockWaiterEthereumClient{
			BlockNumberFunc: func(ctx context.Context) (uint64, error) {
				responses := []uint64{90, 95, 110} // Reaches target on 3rd call
				result := responses[callCount]
				callCount++
				return result, nil
			},
		}

		service := service.NewBlockWaiterService(mockTracer, mockEthClient)
		ctx := context.Background()

		blockNum, err := service.WaitForBlock(ctx, 100)

		require.NoError(t, err)
		assert.Equal(t, uint64(110), blockNum)
		assert.Equal(t, 3, len(mockEthClient.BlockNumberCalls.Calls))
	})

	t.Run("continues polling after transient errors then succeeds", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		callCount := 0
		mockEthClient := &MockBlockWaiterEthereumClient{
			BlockNumberFunc: func(ctx context.Context) (uint64, error) {
				callCount++
				// Fail first 3 times, succeed on 4th
				if callCount <= 3 {
					return 0, errors.New("RPC connection failed")
				}
				return 150, nil
			},
		}

		svc := service.NewBlockWaiterService(mockTracer, mockEthClient)
		ctx := context.Background()

		blockNum, err := svc.WaitForBlock(ctx, 100)

		require.NoError(t, err)
		assert.Equal(t, uint64(150), blockNum)
		assert.Equal(t, 4, len(mockEthClient.BlockNumberCalls.Calls))
	})

	t.Run("returns context error when cancelled during transient errors", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		mockEthClient := &MockBlockWaiterEthereumClient{
			BlockNumberFunc: func(ctx context.Context) (uint64, error) {
				return 0, errors.New("RPC connection failed")
			},
		}

		svc := service.NewBlockWaiterService(mockTracer, mockEthClient)
		ctx, cancel := context.WithCancel(context.Background())
		// Cancel immediately so that on next tick ctx.Done() fires
		cancel()

		blockNum, err := svc.WaitForBlock(ctx, 100)

		require.Error(t, err)
		assert.Equal(t, uint64(0), blockNum)
		assert.True(t, errors.Is(err, context.Canceled))
	})

	t.Run("returns the current block number when the target block number is 0", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		mockEthClient := &MockBlockWaiterEthereumClient{
			BlockNumberFunc: func(ctx context.Context) (uint64, error) {
				return 5, nil // Any positive number is > 0
			},
		}

		service := service.NewBlockWaiterService(mockTracer, mockEthClient)
		ctx := context.Background()

		blockNum, err := service.WaitForBlock(ctx, 0)

		require.NoError(t, err)
		assert.Equal(t, uint64(5), blockNum)
		assert.Equal(t, 1, len(mockEthClient.BlockNumberCalls.Calls))
	})
}
