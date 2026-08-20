package service_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	"github.com/raylsnetwork/rayls-sovereign-relayer/enygma/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/enygma/testutils"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockRetryBlockWaiter struct {
	WaitForBlockFunc func(ctx context.Context, targetBlockNumber uint64) (uint64, error)

	WaitForBlockCalls struct {
		sync.RWMutex
		Calls []struct {
			Ctx               context.Context
			TargetBlockNumber uint64
		}
	}
}

func (m *MockRetryBlockWaiter) WaitForBlock(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
	if m.WaitForBlockFunc == nil {
		panic("MockRetryBlockWaiter.WaitForBlockFunc: method is nil but WaitForBlock was just called")
	}

	m.WaitForBlockCalls.Lock()
	m.WaitForBlockCalls.Calls = append(m.WaitForBlockCalls.Calls, struct {
		Ctx               context.Context
		TargetBlockNumber uint64
	}{Ctx: ctx, TargetBlockNumber: targetBlockNumber})
	m.WaitForBlockCalls.Unlock()

	return m.WaitForBlockFunc(ctx, targetBlockNumber)
}

type MockExecuteOperation struct {
	ExecuteFunc func(ctx context.Context, blockNumber uint64) error

	ExecuteCalls struct {
		sync.RWMutex
		Calls []struct {
			Ctx         context.Context
			BlockNumber uint64
		}
	}
}

func (m *MockExecuteOperation) Execute(ctx context.Context, blockNumber uint64) error {
	if m.ExecuteFunc == nil {
		panic("MockExecuteOperation.ExecuteFunc: method is nil but Execute was just called")
	}

	m.ExecuteCalls.Lock()
	m.ExecuteCalls.Calls = append(m.ExecuteCalls.Calls, struct {
		Ctx         context.Context
		BlockNumber uint64
	}{Ctx: ctx, BlockNumber: blockNumber})
	m.ExecuteCalls.Unlock()

	return m.ExecuteFunc(ctx, blockNumber)
}

func TestRetryService_RetryOperation(t *testing.T) {
	t.Run("succeeds on first attempt", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		mockBlockWaiter := &MockRetryBlockWaiter{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				return 100, nil
			},
		}

		service := service.NewRetryService(mockTracer, mockBlockWaiter)
		ctx := context.Background()

		mockOperation := &MockExecuteOperation{
			ExecuteFunc: func(ctx context.Context, blockNumber uint64) error {
				return nil // Succeed on first attempt
			},
		}

		blockNum, err := service.RetryOperation(ctx, "testOp", 3, 100, mockOperation.Execute)

		require.NoError(t, err)
		assert.Equal(t, uint64(100), blockNum)
		assert.Equal(t, 1, len(mockOperation.ExecuteCalls.Calls))
		assert.Equal(t, 0, len(mockBlockWaiter.WaitForBlockCalls.Calls)) // No retries, so no block waiter calls
	})

	t.Run("succeeds after one retry", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		mockBlockWaiter := &MockRetryBlockWaiter{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				return 101, nil // Block waiter will return 101 for targetBlock 100
			},
		}

		service := service.NewRetryService(mockTracer, mockBlockWaiter)
		ctx := context.Background()

		retryableErr := errors.New("Invalid public signal for balance")
		callCount := 0
		mockOperation := &MockExecuteOperation{
			ExecuteFunc: func(ctx context.Context, blockNumber uint64) error {
				responses := []error{retryableErr, nil} // Fail once, then succeed
				if callCount < len(responses) {
					err := responses[callCount]
					callCount++
					return err
				}
				return responses[len(responses)-1]
			},
		}

		blockNum, err := service.RetryOperation(ctx, "testOp", 3, 100, mockOperation.Execute)

		require.NoError(t, err)
		assert.Equal(t, uint64(101), blockNum)
		assert.Equal(t, 2, len(mockOperation.ExecuteCalls.Calls))
		assert.Equal(t, 1, len(mockBlockWaiter.WaitForBlockCalls.Calls)) // One retry = one block wait
	})

	t.Run("succeeds after exceeding retry threshold", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		// Block waiter responses for retries 2-12
		// Retries 2-10 wait for 1 block (blockNumber + 0)
		// Retries 11-12 wait for 2 blocks (blockNumber + 1)
		blockResponses := []uint64{101, 102, 103, 104, 105, 106, 107, 108, 109, 111, 113}
		blockCallCount := 0
		mockBlockWaiter := &MockRetryBlockWaiter{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				if blockCallCount < len(blockResponses) {
					result := blockResponses[blockCallCount]
					blockCallCount++
					return result, nil
				}
				return blockResponses[len(blockResponses)-1], nil
			},
		}

		service := service.NewRetryService(mockTracer, mockBlockWaiter)
		ctx := context.Background()

		retryableErr := errors.New("Contract is processing another transaction.")
		// Fail 11 times, succeed on 12th attempt
		responses := make([]error, 12)
		for i := 0; i < 11; i++ {
			responses[i] = retryableErr
		}
		responses[11] = nil

		opCallCount := 0
		mockOperation := &MockExecuteOperation{
			ExecuteFunc: func(ctx context.Context, blockNumber uint64) error {
				if opCallCount < len(responses) {
					err := responses[opCallCount]
					opCallCount++
					return err
				}
				return responses[len(responses)-1]
			},
		}

		blockNum, err := service.RetryOperation(ctx, "testOp", 12, 100, mockOperation.Execute)

		require.NoError(t, err)
		assert.Equal(t, uint64(113), blockNum)
		assert.Equal(t, 12, len(mockOperation.ExecuteCalls.Calls))
		assert.Equal(t, 11, len(mockBlockWaiter.WaitForBlockCalls.Calls))
	})

	t.Run("returns error on non retryable error", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		mockBlockWaiter := &MockRetryBlockWaiter{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				return 0, errors.New("should not be called")
			},
		}

		service := service.NewRetryService(mockTracer, mockBlockWaiter)
		ctx := context.Background()

		nonRetryableErr := errors.New("unknown error")
		mockOperation := &MockExecuteOperation{
			ExecuteFunc: func(ctx context.Context, blockNumber uint64) error {
				return nonRetryableErr
			},
		}

		blockNum, err := service.RetryOperation(ctx, "testOp", 3, 100, mockOperation.Execute)

		require.Error(t, err)
		assert.Equal(t, uint64(100), blockNum)
		assert.Equal(t, 1, len(mockOperation.ExecuteCalls.Calls))
		assert.Equal(t, 0, len(mockBlockWaiter.WaitForBlockCalls.Calls))
	})

	t.Run("returns error when block waiter fails", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		waiterErr := errors.New("failed to get block")
		mockBlockWaiter := &MockRetryBlockWaiter{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				return 0, waiterErr // Fail on first WaitForBlock call
			},
		}

		service := service.NewRetryService(mockTracer, mockBlockWaiter)
		ctx := context.Background()

		retryableErr := errors.New("Invalid BlockNumber Used in Proof.")
		callCount := 0
		mockOperation := &MockExecuteOperation{
			ExecuteFunc: func(ctx context.Context, blockNumber uint64) error {
				responses := []error{retryableErr, retryableErr} // Trigger retry
				if callCount < len(responses) {
					err := responses[callCount]
					callCount++
					return err
				}
				return responses[len(responses)-1]
			},
		}

		blockNum, err := service.RetryOperation(ctx, "testOp", 3, 100, mockOperation.Execute)

		require.Error(t, err)
		assert.Equal(t, uint64(0), blockNum)
		// Check that error message contains the expected error wrapper
		assert.Contains(t, err.Error(), "error while waiting for the next block number")
	})

	t.Run("returns error when max retries exceeded", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		blockResponses := []uint64{101, 102, 103}
		blockCallCount := 0
		mockBlockWaiter := &MockRetryBlockWaiter{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				if blockCallCount < len(blockResponses) {
					result := blockResponses[blockCallCount]
					blockCallCount++
					return result, nil
				}
				return blockResponses[len(blockResponses)-1], nil
			},
		}

		svc := service.NewRetryService(mockTracer, mockBlockWaiter)
		ctx := context.Background()

		retryableErr := errors.New("Nullifier already used in pending transaction.")
		callCount := 0
		mockOperation := &MockExecuteOperation{
			ExecuteFunc: func(ctx context.Context, blockNumber uint64) error {
				responses := []error{retryableErr, retryableErr, retryableErr, retryableErr}
				if callCount < len(responses) {
					err := responses[callCount]
					callCount++
					return err
				}
				return responses[len(responses)-1]
			},
		}

		blockNum, err := svc.RetryOperation(ctx, "testOp", 3, 100, mockOperation.Execute)

		require.Error(t, err)
		assert.Equal(t, uint64(102), blockNum)
		assert.Equal(t, 3, len(mockOperation.ExecuteCalls.Calls))
	})
}

func TestRetryService_IsRetryableError(t *testing.T) {
	t.Run("contract error is retryable", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		mockBlockWaiter := &MockRetryBlockWaiter{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				return targetBlockNumber + 1, nil
			},
		}
		svc := service.NewRetryService(mockTracer, mockBlockWaiter)
		ctx := context.Background()

		contractErr := errors.New("BlockNumber in Proof was already finalised.")
		mockOperation := &MockExecuteOperation{
			ExecuteFunc: func(ctx context.Context, blockNumber uint64) error {
				return contractErr
			},
		}

		blockNum, err := svc.RetryOperation(ctx, "testOp", 3, 100, mockOperation.Execute)

		// Should retry and eventually return the error after max retries
		require.Error(t, err)
		// Contract error is retryable, so it should attempt retries (3 retries + 1 initial = 4 attempts max)
		assert.Equal(t, 3, len(mockOperation.ExecuteCalls.Calls), "contract error should trigger retries")
		// blockNum should be >= 100 (it gets incremented through retries)
		assert.GreaterOrEqual(t, blockNum, uint64(100))
	})

	t.Run("hex-encoded ErrorWithRevertData with a retryable reason is retryable", func(t *testing.T) {
		// Regression for the hex-vs-ASCII misclassification: an on-chain revert surfaces as a
		// *contractclient.ErrorWithRevertData. Before the Error() decode fix, its message was raw hex
		// ("transaction reverted: 0x08c379a0...496e76616c6964...") so the retryable-reason substring
		// match never fired and a transient stale-proof race ("Invalid public signal for balance" — the
		// finalised balance moved between proof-gen and submission) was wrongly classified non-retryable
		// and the cross-transfer was dropped instead of regenerated+resubmitted. The plain
		// errors.New("Invalid public signal for balance") used by the case above does NOT exercise this
		// path because it already carries the ASCII reason. This case uses the real revert payload.
		mockTracer := &testutils.MockTracer{}
		mockBlockWaiter := &MockRetryBlockWaiter{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				return targetBlockNumber + 1, nil
			},
		}
		svc := service.NewRetryService(mockTracer, mockBlockWaiter)
		ctx := context.Background()

		// Wrap exactly as the production cross-transfer signing path does (fmt.Errorf("...: %w", err)).
		revertErr := fmt.Errorf("signing enygma transfer batch: %w",
			contractclient.NewErrorWithRevertData(testtools.ErrorStringRevertData(t, "Invalid public signal for balance")))
		mockOperation := &MockExecuteOperation{
			ExecuteFunc: func(ctx context.Context, blockNumber uint64) error {
				return revertErr
			},
		}

		_, err := svc.RetryOperation(ctx, "crossTransfer", 3, 100, mockOperation.Execute)

		require.Error(t, err)
		assert.Equal(t, 3, len(mockOperation.ExecuteCalls.Calls),
			"a hex-encoded ErrorWithRevertData carrying a retryable reason must trigger retries, not be dropped as non-retryable")
	})

	t.Run("database sync error is retryable", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		mockBlockWaiter := &MockRetryBlockWaiter{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				return targetBlockNumber + 1, nil
			},
		}
		svc := service.NewRetryService(mockTracer, mockBlockWaiter)
		ctx := context.Background()

		dbSyncErr := errors.New("Failed to sync fresh state from smart contract")
		mockOperation := &MockExecuteOperation{
			ExecuteFunc: func(ctx context.Context, blockNumber uint64) error {
				return dbSyncErr
			},
		}

		blockNum, err := svc.RetryOperation(ctx, "testOp", 3, 100, mockOperation.Execute)

		// Should retry and eventually return the error after max retries
		require.Error(t, err)
		// Database sync error is retryable, so it should attempt retries
		assert.Equal(t, 3, len(mockOperation.ExecuteCalls.Calls), "database sync error should trigger retries")
		// blockNum should be >= 100 (it gets incremented through retries)
		assert.GreaterOrEqual(t, blockNum, uint64(100))
	})

	t.Run("transaction simulation error is retryable", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		mockBlockWaiter := &MockRetryBlockWaiter{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				return targetBlockNumber + 1, nil
			},
		}
		svc := service.NewRetryService(mockTracer, mockBlockWaiter)
		ctx := context.Background()

		simErr := errors.New("simulation did not revert")
		mockOperation := &MockExecuteOperation{
			ExecuteFunc: func(ctx context.Context, blockNumber uint64) error {
				return simErr
			},
		}

		blockNum, err := svc.RetryOperation(ctx, "testOp", 3, 100, mockOperation.Execute)

		// Should retry and eventually return the error after max retries
		require.Error(t, err)
		// Simulation error is retryable, so it should attempt retries
		assert.Equal(t, 3, len(mockOperation.ExecuteCalls.Calls), "simulation error should trigger retries")
		// blockNum should be >= 100 (it gets incremented through retries)
		assert.GreaterOrEqual(t, blockNum, uint64(100))
	})

	t.Run("timeout error is retryable", func(t *testing.T) {
		mockTracer := &testutils.MockTracer{}
		mockBlockWaiter := &MockRetryBlockWaiter{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				return targetBlockNumber + 1, nil
			},
		}
		svc := service.NewRetryService(mockTracer, mockBlockWaiter)
		ctx := context.Background()

		timeoutErr := errors.New("context deadline exceeded")
		mockOperation := &MockExecuteOperation{
			ExecuteFunc: func(ctx context.Context, blockNumber uint64) error {
				return timeoutErr
			},
		}

		blockNum, err := svc.RetryOperation(ctx, "testOp", 3, 100, mockOperation.Execute)

		// Should retry and eventually return the error after max retries
		require.Error(t, err)
		// Timeout error is retryable, so it should attempt retries
		assert.Equal(t, 3, len(mockOperation.ExecuteCalls.Calls), "timeout error should trigger retries")
		// blockNum should be >= 100 (it gets incremented through retries)
		assert.GreaterOrEqual(t, blockNum, uint64(100))
	})

	connectionErrorCases := []struct {
		name     string
		errorMsg string
	}{
		{"connection refused", "dial tcp 127.0.0.1:8545: connect: connection refused"},
		{"connection reset", "read tcp 10.0.0.1:54321->10.0.0.2:8545: read: connection reset by peer"},
		{"dial tcp", "dial tcp: lookup hub.example.com: no such host"},
		{"i/o timeout", "dial tcp 10.0.0.2:8545: i/o timeout"},
		{"EOF", "Post \"http://hub:8545\": EOF"},
		{"broken pipe", "write tcp 10.0.0.1:54321->10.0.0.2:8545: write: broken pipe"},
		{"no such host", "dial tcp: lookup hub.internal: no such host"},
		{"TLS handshake timeout", "net/http: TLS handshake timeout"},
	}

	for _, tc := range connectionErrorCases {
		t.Run("connection error is retryable: "+tc.name, func(t *testing.T) {
			mockTracer := &testutils.MockTracer{}
			mockBlockWaiter := &MockRetryBlockWaiter{
				WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
					return targetBlockNumber + 1, nil
				},
			}
			svc := service.NewRetryService(mockTracer, mockBlockWaiter)
			ctx := context.Background()

			connErr := errors.New(tc.errorMsg)
			mockOperation := &MockExecuteOperation{
				ExecuteFunc: func(ctx context.Context, blockNumber uint64) error {
					return connErr
				},
			}

			_, err := svc.RetryOperation(ctx, "testOp", 3, 100, mockOperation.Execute)

			require.Error(t, err)
			assert.Equal(
				t,
				3,
				len(mockOperation.ExecuteCalls.Calls),
				"connection error '%s' should trigger retries",
				tc.name,
			)
			assert.Contains(t, err.Error(), "max retries exceeded")
		})
	}
}
