package service_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/msgqueue"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools/fake"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools/spy"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSingleMessageMQ(msg msgqueue.Message[service.EnygmaDestMessage]) *EnygmaDestMQMock {
	return &EnygmaDestMQMock{
		NextFunc: fake.NextMQ(msg),
	}
}

func newDefaultReceiverMock(t *testing.T) *EnygmaReceiverMock {
	return &EnygmaReceiverMock{
		HandleEnygmaCrossTransferFunc: func(ctx context.Context, batch *types.EnygmaTransferBatch) error {
			assert.Fail(t, "shouldn't call receiver handler")
			return nil
		},
	}
}

func newDefaultCheckpointRepoMock(t *testing.T) *EnygmaCheckpointRepositoryMock {
	return &EnygmaCheckpointRepositoryMock{
		GetLatestCheckpointByFiltersFunc: func(ctx context.Context, resourceId string, status *types.EnygmaCheckpointStatus, finalizedBlockNumber *big.Int, pendingBlockNumber *big.Int) (*types.EnygmaCheckpoint, error) {
			return nil, nil //nolint:nilnil // intentional nil return in test mock
		},
		CreateEnygmaCheckpointFunc: func(ctx context.Context, checkpoint types.EnygmaCheckpoint) error {
			return nil
		},
	}
}

func newDefaultBackoffMock() *BackoffMock {
	return &BackoffMock{
		DoFunc: func(ctx context.Context, maxAttempts int, fn func() error) error {
			return fn()
		},
	}
}

func TestEnygmaOrchestrator(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("supports graceful shutdown", func(t *testing.T) {
		enygmaMQ := &EnygmaDestMQMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[service.EnygmaDestMessage], error) {
				<-ctx.Done()
				return msgqueue.Message[service.EnygmaDestMessage]{}, context.Canceled
			},
		}
		receiver := newDefaultReceiverMock(t)
		checkpointRepo := newDefaultCheckpointRepoMock(t)

		svc := service.NewEnygmaOrchestratorWithBackoff(enygmaMQ, receiver, checkpointRepo, newDefaultBackoffMock())
		hasGracefulShutdown := testtools.ShutdownFixture(t, svc.Run, time.Second)

		assert.True(t, hasGracefulShutdown)
	})

	t.Run("calls receiver handler on transfer batch message", func(t *testing.T) {
		wantResourceID := "example-resource-id"
		wantBatch := &types.EnygmaTransferBatch{
			ResourceId:  wantResourceID,
			FromChainID: big.NewInt(999),
		}

		ackSpy := spy.NewAck()

		enygmaMQ := newSingleMessageMQ(msgqueue.Message[service.EnygmaDestMessage]{
			V: service.EnygmaDestMessage{
				ID:            "test-message-id",
				Type:          service.EnygmaTransferBatchMessage,
				BlockNumber:   100,
				TransferBatch: wantBatch,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultReceiverMock(t)
		receiver.HandleEnygmaCrossTransferFunc = func(ctx context.Context, batch *types.EnygmaTransferBatch) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantResourceID, batch.ResourceId)
			return nil
		}
		checkpointRepo := newDefaultCheckpointRepoMock(t)

		svc := service.NewEnygmaOrchestratorWithBackoff(enygmaMQ, receiver, checkpointRepo, newDefaultBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(receiver.HandleEnygmaCrossTransferCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack message on receiver error", func(t *testing.T) {
		wantBatch := &types.EnygmaTransferBatch{
			ResourceId:  "example-resource-id",
			FromChainID: big.NewInt(999),
		}

		ackSpy := spy.NewAck()

		enygmaMQ := newSingleMessageMQ(msgqueue.Message[service.EnygmaDestMessage]{
			V: service.EnygmaDestMessage{
				ID:            "test-message-id",
				Type:          service.EnygmaTransferBatchMessage,
				BlockNumber:   100,
				TransferBatch: wantBatch,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultReceiverMock(t)
		receiver.HandleEnygmaCrossTransferFunc = func(ctx context.Context, batch *types.EnygmaTransferBatch) error {
			return errors.New("receiver error")
		}
		checkpointRepo := newDefaultCheckpointRepoMock(t)

		svc := service.NewEnygmaOrchestratorWithBackoff(enygmaMQ, receiver, checkpointRepo, newDefaultBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(receiver.HandleEnygmaCrossTransferCalls()))
		ackSpy.AssertNotCalled(t, "should not ack message on receiver error")
	})

	t.Run("creates checkpoint on finalized balance message", func(t *testing.T) {
		wantResourceID := "example-resource-id"
		wantBalance := &types.EnygmaFinalizedBalance{
			ResourceId:           wantResourceID,
			FinalizedBlockNumber: big.NewInt(100),
			PendingBlockNumber:   big.NewInt(110),
			BalanceX:             big.NewInt(1000),
			BalanceY:             big.NewInt(2000),
		}

		ackSpy := spy.NewAck()

		enygmaMQ := newSingleMessageMQ(msgqueue.Message[service.EnygmaDestMessage]{
			V: service.EnygmaDestMessage{
				ID:               "test-message-id",
				Type:             service.EnygmaFinalizedBalanceMessage,
				BlockNumber:      100,
				FinalizedBalance: wantBalance,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultReceiverMock(t)
		checkpointRepo := newDefaultCheckpointRepoMock(t)
		checkpointRepo.CreateEnygmaCheckpointFunc = func(ctx context.Context, checkpoint types.EnygmaCheckpoint) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER creating checkpoint")
			assert.Equal(t, wantResourceID, checkpoint.ResourceId)
			assert.Equal(t, wantBalance.BalanceX, checkpoint.FinalizedPublicBalanceX)
			assert.Equal(t, wantBalance.BalanceY, checkpoint.FinalizedPublicBalanceY)
			assert.Equal(t, types.EnygmaCheckpointStatusTentative, checkpoint.Status)
			return nil
		}

		svc := service.NewEnygmaOrchestratorWithBackoff(enygmaMQ, receiver, checkpointRepo, newDefaultBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(checkpointRepo.CreateEnygmaCheckpointCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("skips checkpoint if already exists", func(t *testing.T) {
		wantResourceID := "example-resource-id"
		wantBalance := &types.EnygmaFinalizedBalance{
			ResourceId:           wantResourceID,
			FinalizedBlockNumber: big.NewInt(100),
			PendingBlockNumber:   big.NewInt(110),
			BalanceX:             big.NewInt(1000),
			BalanceY:             big.NewInt(2000),
		}

		ackSpy := spy.NewAck()

		enygmaMQ := newSingleMessageMQ(msgqueue.Message[service.EnygmaDestMessage]{
			V: service.EnygmaDestMessage{
				ID:               "test-message-id",
				Type:             service.EnygmaFinalizedBalanceMessage,
				BlockNumber:      100,
				FinalizedBalance: wantBalance,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultReceiverMock(t)
		checkpointRepo := newDefaultCheckpointRepoMock(t)
		checkpointRepo.GetLatestCheckpointByFiltersFunc = func(ctx context.Context, resourceId string, status *types.EnygmaCheckpointStatus, finalizedBlockNumber *big.Int, pendingBlockNumber *big.Int) (*types.EnygmaCheckpoint, error) {
			return &types.EnygmaCheckpoint{ID: "existing-checkpoint"}, nil
		}
		checkpointRepo.CreateEnygmaCheckpointFunc = func(ctx context.Context, checkpoint types.EnygmaCheckpoint) error {
			assert.Fail(t, "should not create checkpoint if already exists")
			return nil
		}

		svc := service.NewEnygmaOrchestratorWithBackoff(enygmaMQ, receiver, checkpointRepo, newDefaultBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 0, len(checkpointRepo.CreateEnygmaCheckpointCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("skips zero balance checkpoint", func(t *testing.T) {
		wantResourceID := "example-resource-id"
		wantBalance := &types.EnygmaFinalizedBalance{
			ResourceId:           wantResourceID,
			FinalizedBlockNumber: big.NewInt(100),
			PendingBlockNumber:   big.NewInt(110),
			BalanceX:             big.NewInt(0),
			BalanceY:             big.NewInt(1),
		}

		ackSpy := spy.NewAck()

		enygmaMQ := newSingleMessageMQ(msgqueue.Message[service.EnygmaDestMessage]{
			V: service.EnygmaDestMessage{
				ID:               "test-message-id",
				Type:             service.EnygmaFinalizedBalanceMessage,
				BlockNumber:      100,
				FinalizedBalance: wantBalance,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultReceiverMock(t)
		checkpointRepo := newDefaultCheckpointRepoMock(t)
		checkpointRepo.CreateEnygmaCheckpointFunc = func(ctx context.Context, checkpoint types.EnygmaCheckpoint) error {
			assert.Fail(t, "should not create checkpoint for zero balance")
			return nil
		}

		svc := service.NewEnygmaOrchestratorWithBackoff(enygmaMQ, receiver, checkpointRepo, newDefaultBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 0, len(checkpointRepo.CreateEnygmaCheckpointCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("acks nil transfer batch message", func(t *testing.T) {
		ackSpy := spy.NewAck()

		enygmaMQ := newSingleMessageMQ(msgqueue.Message[service.EnygmaDestMessage]{
			V: service.EnygmaDestMessage{
				ID:            "test-message-id",
				Type:          service.EnygmaTransferBatchMessage,
				BlockNumber:   100,
				TransferBatch: nil,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultReceiverMock(t)
		checkpointRepo := newDefaultCheckpointRepoMock(t)

		svc := service.NewEnygmaOrchestratorWithBackoff(enygmaMQ, receiver, checkpointRepo, newDefaultBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 0, len(receiver.HandleEnygmaCrossTransferCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("acks unknown message type", func(t *testing.T) {
		ackSpy := spy.NewAck()

		enygmaMQ := newSingleMessageMQ(msgqueue.Message[service.EnygmaDestMessage]{
			V: service.EnygmaDestMessage{
				ID:          "test-message-id",
				Type:        service.EnygmaDestMessageType(999),
				BlockNumber: 100,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultReceiverMock(t)
		checkpointRepo := newDefaultCheckpointRepoMock(t)

		svc := service.NewEnygmaOrchestratorWithBackoff(enygmaMQ, receiver, checkpointRepo, newDefaultBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack on checkpoint repository error", func(t *testing.T) {
		wantBalance := &types.EnygmaFinalizedBalance{
			ResourceId:           "example-resource-id",
			FinalizedBlockNumber: big.NewInt(100),
			PendingBlockNumber:   big.NewInt(110),
			BalanceX:             big.NewInt(1000),
			BalanceY:             big.NewInt(2000),
		}

		ackSpy := spy.NewAck()

		enygmaMQ := newSingleMessageMQ(msgqueue.Message[service.EnygmaDestMessage]{
			V: service.EnygmaDestMessage{
				ID:               "test-message-id",
				Type:             service.EnygmaFinalizedBalanceMessage,
				BlockNumber:      100,
				FinalizedBalance: wantBalance,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultReceiverMock(t)
		checkpointRepo := newDefaultCheckpointRepoMock(t)
		checkpointRepo.GetLatestCheckpointByFiltersFunc = func(ctx context.Context, resourceId string, status *types.EnygmaCheckpointStatus, finalizedBlockNumber *big.Int, pendingBlockNumber *big.Int) (*types.EnygmaCheckpoint, error) {
			return nil, errors.New("database error")
		}

		svc := service.NewEnygmaOrchestratorWithBackoff(enygmaMQ, receiver, checkpointRepo, newDefaultBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		ackSpy.AssertNotCalled(t, "should not ack on repository error")
	})

	t.Run("does not ack on checkpoint creation error", func(t *testing.T) {
		wantBalance := &types.EnygmaFinalizedBalance{
			ResourceId:           "example-resource-id",
			FinalizedBlockNumber: big.NewInt(100),
			PendingBlockNumber:   big.NewInt(110),
			BalanceX:             big.NewInt(1000),
			BalanceY:             big.NewInt(2000),
		}

		ackSpy := spy.NewAck()

		enygmaMQ := newSingleMessageMQ(msgqueue.Message[service.EnygmaDestMessage]{
			V: service.EnygmaDestMessage{
				ID:               "test-message-id",
				Type:             service.EnygmaFinalizedBalanceMessage,
				BlockNumber:      100,
				FinalizedBalance: wantBalance,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultReceiverMock(t)
		checkpointRepo := newDefaultCheckpointRepoMock(t)
		checkpointRepo.CreateEnygmaCheckpointFunc = func(ctx context.Context, checkpoint types.EnygmaCheckpoint) error {
			return errors.New("creation error")
		}

		svc := service.NewEnygmaOrchestratorWithBackoff(enygmaMQ, receiver, checkpointRepo, newDefaultBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		ackSpy.AssertNotCalled(t, "should not ack on checkpoint creation error")
	})

	t.Run("acks nil finalized balance message", func(t *testing.T) {
		ackSpy := spy.NewAck()

		enygmaMQ := newSingleMessageMQ(msgqueue.Message[service.EnygmaDestMessage]{
			V: service.EnygmaDestMessage{
				ID:               "test-message-id",
				Type:             service.EnygmaFinalizedBalanceMessage,
				BlockNumber:      100,
				FinalizedBalance: nil,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultReceiverMock(t)
		checkpointRepo := newDefaultCheckpointRepoMock(t)

		svc := service.NewEnygmaOrchestratorWithBackoff(enygmaMQ, receiver, checkpointRepo, newDefaultBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 0, len(checkpointRepo.CreateEnygmaCheckpointCalls()))
		ackSpy.AssertCalled(t)
	})
}
