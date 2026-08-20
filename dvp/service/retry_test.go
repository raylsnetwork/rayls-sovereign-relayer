package service_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/dvp/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/enygma/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockDvpExecuteOperation struct {
	ExecuteFunc func(ctx context.Context) error

	ExecuteCalls struct {
		sync.RWMutex
		Calls []struct {
			Ctx context.Context
		}
	}
}

func (m *MockDvpExecuteOperation) Execute(ctx context.Context) error {
	m.ExecuteCalls.Lock()
	m.ExecuteCalls.Calls = append(m.ExecuteCalls.Calls, struct {
		Ctx context.Context
	}{Ctx: ctx})
	m.ExecuteCalls.Unlock()

	return m.ExecuteFunc(ctx)
}

func TestDvpRetryService_RetryOperation(t *testing.T) {
	t.Run("succeeds on first attempt", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		svc := service.NewRetryServiceWithDelays(mockTracer, time.Millisecond, time.Millisecond)
		ctx := context.Background()

		mockOperation := &MockDvpExecuteOperation{
			ExecuteFunc: func(ctx context.Context) error {
				return nil
			},
		}

		err := svc.RetryOperation(ctx, "testOp", 3, mockOperation.Execute)

		require.NoError(t, err)
		assert.Equal(t, 1, len(mockOperation.ExecuteCalls.Calls))
	})

	t.Run("retries on nonce too low and succeeds", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		svc := service.NewRetryServiceWithDelays(mockTracer, time.Millisecond, time.Millisecond)
		ctx := context.Background()

		callCount := 0
		mockOperation := &MockDvpExecuteOperation{
			ExecuteFunc: func(ctx context.Context) error {
				callCount++
				if callCount < 3 {
					return errors.New("nonce too low")
				}
				return nil
			},
		}

		err := svc.RetryOperation(ctx, "testOp", 5, mockOperation.Execute)

		require.NoError(t, err)
		assert.Equal(t, 3, len(mockOperation.ExecuteCalls.Calls))
	})

	t.Run("does not retry on context deadline exceeded", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		svc := service.NewRetryServiceWithDelays(mockTracer, time.Millisecond, time.Millisecond)
		ctx := context.Background()

		mockOperation := &MockDvpExecuteOperation{
			ExecuteFunc: func(ctx context.Context) error {
				return errors.New("context deadline exceeded")
			},
		}

		err := svc.RetryOperation(ctx, "testOp", 3, mockOperation.Execute)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "non-retryable error")
		assert.Equal(t, 1, len(mockOperation.ExecuteCalls.Calls))
	})

	t.Run("returns error on non-retryable error", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		svc := service.NewRetryServiceWithDelays(mockTracer, time.Millisecond, time.Millisecond)
		ctx := context.Background()

		nonRetryableErr := errors.New("unknown error")
		mockOperation := &MockDvpExecuteOperation{
			ExecuteFunc: func(ctx context.Context) error {
				return nonRetryableErr
			},
		}

		err := svc.RetryOperation(ctx, "testOp", 3, mockOperation.Execute)

		require.Error(t, err)
		assert.ErrorIs(t, err, nonRetryableErr)
		assert.Equal(t, 1, len(mockOperation.ExecuteCalls.Calls))
	})

	t.Run("returns error when max retries exceeded", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		svc := service.NewRetryServiceWithDelays(mockTracer, time.Millisecond, time.Millisecond)
		ctx := context.Background()

		nonceErr := errors.New("nonce too low")
		mockOperation := &MockDvpExecuteOperation{
			ExecuteFunc: func(ctx context.Context) error {
				return nonceErr
			},
		}

		err := svc.RetryOperation(ctx, "testOp", 3, mockOperation.Execute)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "max retries exceeded")
		assert.ErrorIs(t, err, nonceErr)
		assert.Equal(t, 3, len(mockOperation.ExecuteCalls.Calls))
	})
}
