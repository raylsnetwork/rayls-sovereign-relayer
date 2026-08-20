package service_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/enygma/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/enygma/testutils"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mocks
type creationServiceRepositoryManagerMock struct {
	WithTransactionFunc  func(ctx context.Context, fn func(ctx context.Context) error) error
	withTransactionCalls int
}

func (m *creationServiceRepositoryManagerMock) WithTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	m.withTransactionCalls++
	if m.WithTransactionFunc == nil {
		// Execute the function (for simple tests)
		return fn(ctx)
	}
	return m.WithTransactionFunc(ctx, fn)
}

type creationServiceEnygmaHistoryRepositoryMock struct {
	InsertEnygmaHistoryFunc  func(ctx context.Context, history types.EnygmaHistory) error
	insertEnygmaHistoryCalls int
	lastHistory              types.EnygmaHistory
}

func (m *creationServiceEnygmaHistoryRepositoryMock) InsertEnygmaHistory(
	ctx context.Context,
	history types.EnygmaHistory,
) error {
	m.insertEnygmaHistoryCalls++
	m.lastHistory = history
	if m.InsertEnygmaHistoryFunc == nil {
		return nil
	}
	return m.InsertEnygmaHistoryFunc(ctx, history)
}

type creationServiceEnygmaRepositoryMock struct {
	CreateEnygmaFunc  func(ctx context.Context, enygma types.Enygma) error
	createEnygmaCalls int
	lastEnygma        types.Enygma
}

func (m *creationServiceEnygmaRepositoryMock) CreateEnygma(ctx context.Context, enygma types.Enygma) error {
	m.createEnygmaCalls++
	m.lastEnygma = enygma
	if m.CreateEnygmaFunc == nil {
		return nil
	}
	return m.CreateEnygmaFunc(ctx, enygma)
}

// Helper functions for test data
func creationTestResourceId() string {
	return "test-resource-id-123"
}

func creationTestChainId() *big.Int {
	return big.NewInt(1)
}

func creationTestBlockNumber() *big.Int {
	return big.NewInt(100)
}

func TestEnygmaCreationService_CreateEnygma(t *testing.T) {
	t.Run("successfully creates enygma", func(t *testing.T) {
		ctx := context.Background()
		resourceId := creationTestResourceId()
		chainId := creationTestChainId()
		blockNumber := creationTestBlockNumber()

		tracer := &testutils.MockTracer{}
		repositoryManager := &creationServiceRepositoryManagerMock{}
		historyRepository := &creationServiceEnygmaHistoryRepositoryMock{}
		enygmaRepository := &creationServiceEnygmaRepositoryMock{}

		svc := service.NewEnygmaCreationService(tracer, repositoryManager, enygmaRepository, historyRepository)

		err := svc.CreateEnygma(ctx, resourceId, chainId, blockNumber)

		require.NoError(t, err)
		assert.Equal(t, 1, repositoryManager.withTransactionCalls)
		assert.Equal(t, 1, historyRepository.insertEnygmaHistoryCalls)
		assert.Equal(t, 1, enygmaRepository.createEnygmaCalls)

		// Verify history was created with correct values
		assert.Equal(t, resourceId, historyRepository.lastHistory.ResourceId)
		assert.Equal(t, chainId, historyRepository.lastHistory.FromChainId)
		assert.Equal(t, blockNumber, historyRepository.lastHistory.BlockNumberPrivateHub)
		assert.Equal(t, types.EnygmaCreation, historyRepository.lastHistory.EventType)
		assert.Equal(t, big.NewInt(0), historyRepository.lastHistory.RFactor)
		assert.Equal(t, big.NewInt(0), historyRepository.lastHistory.BalanceChange)

		// Verify enygma was created with correct values
		assert.Equal(t, resourceId, enygmaRepository.lastEnygma.ResourceId)
		assert.Equal(t, big.NewInt(0), enygmaRepository.lastEnygma.FinalizedR)
		assert.Equal(t, big.NewInt(0), enygmaRepository.lastEnygma.FinalizedBalance)
		assert.Equal(t, blockNumber, enygmaRepository.lastEnygma.FinalizedBlockNumber)
		assert.Equal(t, blockNumber, enygmaRepository.lastEnygma.PendingBlockNumber)
	})

	t.Run("returns error when WithTransaction fails", func(t *testing.T) {
		ctx := context.Background()
		resourceId := creationTestResourceId()
		chainId := creationTestChainId()
		blockNumber := creationTestBlockNumber()

		tracer := &testutils.MockTracer{}
		repositoryManager := &creationServiceRepositoryManagerMock{
			WithTransactionFunc: func(context.Context, func(context.Context) error) error {
				return errors.New("transaction failed")
			},
		}
		historyRepository := &creationServiceEnygmaHistoryRepositoryMock{}
		enygmaRepository := &creationServiceEnygmaRepositoryMock{}

		svc := service.NewEnygmaCreationService(tracer, repositoryManager, enygmaRepository, historyRepository)

		err := svc.CreateEnygma(ctx, resourceId, chainId, blockNumber)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "transaction failed")
		// Verify repositories were never called
		assert.Equal(t, 0, historyRepository.insertEnygmaHistoryCalls)
		assert.Equal(t, 0, enygmaRepository.createEnygmaCalls)
	})

	t.Run("returns error when InsertEnygmaHistory fails", func(t *testing.T) {
		ctx := context.Background()
		resourceId := creationTestResourceId()
		chainId := creationTestChainId()
		blockNumber := creationTestBlockNumber()

		tracer := &testutils.MockTracer{}
		repositoryManager := &creationServiceRepositoryManagerMock{}
		historyRepository := &creationServiceEnygmaHistoryRepositoryMock{
			InsertEnygmaHistoryFunc: func(ctx context.Context, history types.EnygmaHistory) error {
				return errors.New("history insert failed")
			},
		}
		enygmaRepository := &creationServiceEnygmaRepositoryMock{}

		svc := service.NewEnygmaCreationService(tracer, repositoryManager, enygmaRepository, historyRepository)

		err := svc.CreateEnygma(ctx, resourceId, chainId, blockNumber)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "history insert failed")
		// Verify history was called but enygma was not
		assert.Equal(t, 1, historyRepository.insertEnygmaHistoryCalls)
		assert.Equal(t, 0, enygmaRepository.createEnygmaCalls)
	})

	t.Run("returns error when CreateEnygma fails", func(t *testing.T) {
		ctx := context.Background()
		resourceId := creationTestResourceId()
		chainId := creationTestChainId()
		blockNumber := creationTestBlockNumber()

		tracer := &testutils.MockTracer{}
		repositoryManager := &creationServiceRepositoryManagerMock{}
		historyRepository := &creationServiceEnygmaHistoryRepositoryMock{}
		enygmaRepository := &creationServiceEnygmaRepositoryMock{
			CreateEnygmaFunc: func(ctx context.Context, enygma types.Enygma) error {
				return errors.New("enygma creation failed")
			},
		}

		svc := service.NewEnygmaCreationService(tracer, repositoryManager, enygmaRepository, historyRepository)

		err := svc.CreateEnygma(ctx, resourceId, chainId, blockNumber)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "enygma creation failed")
		// Verify both history and enygma creation were attempted
		assert.Equal(t, 1, historyRepository.insertEnygmaHistoryCalls)
		assert.Equal(t, 1, enygmaRepository.createEnygmaCalls)
	})
}
