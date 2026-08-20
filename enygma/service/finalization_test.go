package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/enygma/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type finalizationTracerMock struct {
	StartFunc  func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
	startCalls int
	lastSpan   string
}

func (m *finalizationTracerMock) Start(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	m.startCalls++
	m.lastSpan = spanName
	if m.StartFunc != nil {
		return m.StartFunc(ctx, spanName, opts...)
	}
	// Default: return noop span from OpenTelemetry
	tracer := otel.GetTracerProvider().Tracer("test")
	return tracer.Start(ctx, spanName, opts...)
}

type finalizationBlockWaiterServiceMock struct {
	WaitForBlockFunc  func(ctx context.Context, targetBlockNumber uint64) (uint64, error)
	waitForBlockCalls int
}

func (m *finalizationBlockWaiterServiceMock) WaitForBlock(
	ctx context.Context,
	targetBlockNumber uint64,
) (uint64, error) {
	m.waitForBlockCalls++
	if m.WaitForBlockFunc == nil {
		return targetBlockNumber, nil
	}
	return m.WaitForBlockFunc(ctx, targetBlockNumber)
}

type finalizationEndpointClientMock struct {
	GetResourceAddressFunc  func(ctx context.Context, resourceId string) (common.Address, error)
	getResourceAddressCalls int
}

func (m *finalizationEndpointClientMock) GetResourceAddress(ctx context.Context, resourceId string) (common.Address, error) {
	m.getResourceAddressCalls++
	if m.GetResourceAddressFunc == nil {
		return common.Address{}, nil
	}
	return m.GetResourceAddressFunc(ctx, resourceId)
}

type finalizationRetryServiceMock struct {
	RetryOperationFunc  func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error)
	retryOperationCalls int
	lastOperationName   string
	lastMaxRetries      int
	lastBlockNumber     uint64
}

func (m *finalizationRetryServiceMock) RetryOperation(
	ctx context.Context,
	operationName string,
	maxRetries int,
	blockNumber uint64,
	executeOperation func(ctx context.Context, nextBlockNumber uint64) error,
) (uint64, error) {
	m.retryOperationCalls++
	m.lastOperationName = operationName
	m.lastMaxRetries = maxRetries
	m.lastBlockNumber = blockNumber
	if m.RetryOperationFunc == nil {
		// Execute the operation once to simulate success
		err := executeOperation(ctx, blockNumber)
		return blockNumber, err
	}
	return m.RetryOperationFunc(ctx, operationName, maxRetries, blockNumber, executeOperation)
}

type finalizationExecutorMock struct {
	ExecuteEnygmaCrossTransferFunc  func(parentCtx context.Context, batchID string, blockNumber uint64, resourceId string, txsByChainID map[string][]*types.EnygmaTransferBatchTx, enygmaAddress common.Address) error
	executeEnygmaCrossTransferCalls int
}

func (m *finalizationExecutorMock) ExecuteEnygmaCrossTransfer(
	parentCtx context.Context,
	batchID string,
	blockNumber uint64,
	resourceId string,
	txsByChainID map[string][]*types.EnygmaTransferBatchTx,
	enygmaAddress common.Address,
) error {
	m.executeEnygmaCrossTransferCalls++
	if m.ExecuteEnygmaCrossTransferFunc == nil {
		return nil
	}
	return m.ExecuteEnygmaCrossTransferFunc(parentCtx, batchID, blockNumber, resourceId, txsByChainID, enygmaAddress)
}

// Test data helpers

func finalizationTestResourceId() string {
	return "test-resource-123"
}

func finalizationTestBlockNumber() uint64 {
	return 100
}

func finalizationTestEnygmaAddress() common.Address {
	return common.HexToAddress("0x1234567890123456789012345678901234567890")
}

func TestEnygmaFinalizationService_ExecuteFinalization(t *testing.T) {
	t.Run("successfully executes finalization", func(t *testing.T) {
		ctx := context.Background()
		resourceId := finalizationTestResourceId()
		blockNumber := finalizationTestBlockNumber()
		enygmaAddress := finalizationTestEnygmaAddress()

		tracer := &finalizationTracerMock{}
		blockWaiter := &finalizationBlockWaiterServiceMock{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				return targetBlockNumber, nil
			},
		}
		endpointClient := &finalizationEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return enygmaAddress, nil
			},
		}
		retryService := &finalizationRetryServiceMock{
			RetryOperationFunc: func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
				err := executeOperation(ctx, blockNumber)
				return blockNumber, err
			},
		}
		executor := &finalizationExecutorMock{
			ExecuteEnygmaCrossTransferFunc: func(parentCtx context.Context, _ string, blockNumber uint64, resourceId string, txsByChainID map[string][]*types.EnygmaTransferBatchTx, enygmaAddress common.Address) error {
				return nil
			},
		}

		svc := service.NewEnygmaFinalizationService(tracer, blockWaiter, endpointClient, retryService, executor)

		err := svc.ExecuteFinalization(ctx, "test-finalization-id", blockNumber, resourceId)

		require.NoError(t, err)
		assert.Equal(t, 1, blockWaiter.waitForBlockCalls)
		assert.Equal(t, 1, endpointClient.getResourceAddressCalls)
		assert.Equal(t, 1, retryService.retryOperationCalls)
		assert.Equal(t, 1, executor.executeEnygmaCrossTransferCalls)
	})

	t.Run("returns error when block waiter fails", func(t *testing.T) {
		ctx := context.Background()
		resourceId := finalizationTestResourceId()
		blockNumber := finalizationTestBlockNumber()

		tracer := &finalizationTracerMock{}
		blockWaiter := &finalizationBlockWaiterServiceMock{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				return 0, errors.New("block waiter service down")
			},
		}
		endpointClient := &finalizationEndpointClientMock{}
		retryService := &finalizationRetryServiceMock{}
		executor := &finalizationExecutorMock{}

		svc := service.NewEnygmaFinalizationService(tracer, blockWaiter, endpointClient, retryService, executor)

		err := svc.ExecuteFinalization(ctx, "test-finalization-id", blockNumber, resourceId)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "block waiter service down")
		assert.Equal(t, 1, blockWaiter.waitForBlockCalls)
		// Should not proceed past block waiter
		assert.Equal(t, 0, endpointClient.getResourceAddressCalls)
	})

	t.Run("returns error when GetResourceAddress fails", func(t *testing.T) {
		ctx := context.Background()
		resourceId := finalizationTestResourceId()
		blockNumber := finalizationTestBlockNumber()

		tracer := &finalizationTracerMock{}
		blockWaiter := &finalizationBlockWaiterServiceMock{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				return targetBlockNumber, nil
			},
		}
		endpointClient := &finalizationEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return common.Address{}, errors.New("resource address lookup failed")
			},
		}
		retryService := &finalizationRetryServiceMock{}
		executor := &finalizationExecutorMock{}

		svc := service.NewEnygmaFinalizationService(tracer, blockWaiter, endpointClient, retryService, executor)

		err := svc.ExecuteFinalization(ctx, "test-finalization-id", blockNumber, resourceId)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "resource address lookup failed")
		assert.Equal(t, 1, blockWaiter.waitForBlockCalls)
		assert.Equal(t, 1, endpointClient.getResourceAddressCalls)
		// Should not proceed past endpoint client
		assert.Equal(t, 0, retryService.retryOperationCalls)
	})

	t.Run("returns error when retry operation fails", func(t *testing.T) {
		ctx := context.Background()
		resourceId := finalizationTestResourceId()
		blockNumber := finalizationTestBlockNumber()
		enygmaAddress := finalizationTestEnygmaAddress()

		tracer := &finalizationTracerMock{}
		blockWaiter := &finalizationBlockWaiterServiceMock{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				return targetBlockNumber, nil
			},
		}
		endpointClient := &finalizationEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return enygmaAddress, nil
			},
		}
		retryService := &finalizationRetryServiceMock{
			RetryOperationFunc: func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
				return 0, errors.New("max retries exceeded")
			},
		}
		executor := &finalizationExecutorMock{}

		svc := service.NewEnygmaFinalizationService(tracer, blockWaiter, endpointClient, retryService, executor)

		err := svc.ExecuteFinalization(ctx, "test-finalization-id", blockNumber, resourceId)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "max retries exceeded")
		assert.Equal(t, 1, retryService.retryOperationCalls)
		assert.Equal(t, 10, retryService.lastMaxRetries)
	})
}
