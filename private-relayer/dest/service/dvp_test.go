package service_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools/fake"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools/spy"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDvpSingleMessageMQ(msg msgqueue.Message[service.DvpDestMessage]) *DvpDestMQMock {
	return &DvpDestMQMock{
		NextFunc: fake.NextMQ(msg),
	}
}

func newDefaultDvpReceiverMock(t *testing.T) *DvpReceiverMock {
	return &DvpReceiverMock{
		HandleCommitmentsFunc: func(ctx context.Context, data *types.DvpCommitmentsData) error {
			assert.Fail(t, "shouldn't call HandleCommitments")
			return nil
		},
		HandleNullifiersFunc: func(ctx context.Context, tokenAddress string, nullifiers []*big.Int) error {
			assert.Fail(t, "shouldn't call HandleNullifiers")
			return nil
		},
		HandleSwapInitiatedFunc: func(ctx context.Context, blockNum uint64, data *types.DvpSwapInitiatedData) error {
			assert.Fail(t, "shouldn't call HandleSwapInitiated")
			return nil
		},
		HandleSwapCompletedFunc: func(ctx context.Context, sharedID string) error {
			assert.Fail(t, "shouldn't call HandleSwapCompleted")
			return nil
		},
		HandleSwapRevertFunc: func(ctx context.Context, sharedId string, status types.DvpSwapStatus) error {
			assert.Fail(t, "shouldn't call HandleSwapRevert")
			return nil
		},
	}
}

func newDefaultDvpOrchestratorBackoffMock() *DvpBackoffMock {
	return &DvpBackoffMock{
		DoFunc: func(ctx context.Context, maxAttempts int, fn func() error) error {
			return fn()
		},
	}
}

func TestDvpOrchestrator(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("supports graceful shutdown", func(t *testing.T) {
		dvpMQ := &DvpDestMQMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[service.DvpDestMessage], error) {
				<-ctx.Done()
				return msgqueue.Message[service.DvpDestMessage]{}, context.Canceled
			},
		}
		receiver := newDefaultDvpReceiverMock(t)

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())
		hasGracefulShutdown := testtools.ShutdownFixture(t, svc.Run, time.Second)

		assert.True(t, hasGracefulShutdown)
	})

	t.Run("calls receiver HandleCommitments on commitments message", func(t *testing.T) {
		wantCommitments := &types.DvpCommitmentsData{
			TokenAddress: "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			TokenType:    big.NewInt(4),
			TreeNumber:   big.NewInt(1),
			Commitments:  []*big.Int{big.NewInt(111), big.NewInt(222)},
		}

		ackSpy := spy.NewAck()

		dvpMQ := newDvpSingleMessageMQ(msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:          "test-message-id",
				Type:        service.DvpCommitmentsMessage,
				BlockNumber: 100,
				Commitments: wantCommitments,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultDvpReceiverMock(t)
		receiver.HandleCommitmentsFunc = func(ctx context.Context, data *types.DvpCommitmentsData) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantCommitments.TokenAddress, data.TokenAddress)
			assert.Equal(t, 0, wantCommitments.TokenType.Cmp(data.TokenType))
			assert.Equal(t, 0, wantCommitments.TreeNumber.Cmp(data.TreeNumber))
			require.Len(t, data.Commitments, 2)
			assert.Equal(t, 0, wantCommitments.Commitments[0].Cmp(data.Commitments[0]))
			assert.Equal(t, 0, wantCommitments.Commitments[1].Cmp(data.Commitments[1]))
			return nil
		}

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(receiver.HandleCommitmentsCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack message on receiver HandleCommitments error", func(t *testing.T) {
		ackSpy := spy.NewAck()

		dvpMQ := newDvpSingleMessageMQ(msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:   "test-message-id",
				Type: service.DvpCommitmentsMessage,
				Commitments: &types.DvpCommitmentsData{
					TokenAddress: "0xaaaa",
					TokenType:    big.NewInt(4),
					TreeNumber:   big.NewInt(1),
					Commitments:  []*big.Int{big.NewInt(111)},
				},
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultDvpReceiverMock(t)
		receiver.HandleCommitmentsFunc = func(ctx context.Context, data *types.DvpCommitmentsData) error {
			return errors.New("receiver error")
		}

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(receiver.HandleCommitmentsCalls()))
		ackSpy.AssertNotCalled(t, "should not ack message on receiver error")
	})

	t.Run("acks nil commitments message", func(t *testing.T) {
		ackSpy := spy.NewAck()

		dvpMQ := newDvpSingleMessageMQ(msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:          "test-message-id",
				Type:        service.DvpCommitmentsMessage,
				BlockNumber: 100,
				Commitments: nil,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultDvpReceiverMock(t)

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 0, len(receiver.HandleCommitmentsCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls receiver HandleNullifiers on nullifiers message", func(t *testing.T) {
		wantNullifiers := []*big.Int{big.NewInt(99999)}
		const wantTokenAddress = "0xToken"

		ackSpy := spy.NewAck()

		dvpMQ := newDvpSingleMessageMQ(msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:          "test-message-id",
				Type:        service.DvpNullifierMessage,
				BlockNumber: 100,
				Nullifiers: &types.DvpNullifierData{
					TokenAddress: wantTokenAddress,
					Nullifiers:   wantNullifiers,
				},
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultDvpReceiverMock(t)
		receiver.HandleNullifiersFunc = func(ctx context.Context, tokenAddress string, nullifiers []*big.Int) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantTokenAddress, tokenAddress)
			require.Len(t, nullifiers, 1)
			assert.Equal(t, 0, wantNullifiers[0].Cmp(nullifiers[0]))
			return nil
		}

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(receiver.HandleNullifiersCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack message on receiver HandleNullifiers error", func(t *testing.T) {
		ackSpy := spy.NewAck()

		dvpMQ := newDvpSingleMessageMQ(msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:          "test-message-id",
				Type:        service.DvpNullifierMessage,
				BlockNumber: 100,
				Nullifiers: &types.DvpNullifierData{
					TokenAddress: "0xToken",
					Nullifiers:   []*big.Int{big.NewInt(99999)},
				},
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultDvpReceiverMock(t)
		receiver.HandleNullifiersFunc = func(ctx context.Context, tokenAddress string, nullifiers []*big.Int) error {
			return errors.New("receiver error")
		}

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(receiver.HandleNullifiersCalls()))
		ackSpy.AssertNotCalled(t, "should not ack message on receiver error")
	})

	t.Run("acks nil nullifiers message", func(t *testing.T) {
		ackSpy := spy.NewAck()

		dvpMQ := newDvpSingleMessageMQ(msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:          "test-message-id",
				Type:        service.DvpNullifierMessage,
				BlockNumber: 100,
				Nullifiers:  nil,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultDvpReceiverMock(t)

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 0, len(receiver.HandleNullifiersCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls receiver HandleSwapInitiated on swap initiated message", func(t *testing.T) {
		wantData := &types.DvpSwapInitiatedData{
			Message:           &types.DvpSwapMessage{SharedId: "shared-1"},
			InitiatorDestSalt: big.NewInt(0xBBBB),
		}

		ackSpy := spy.NewAck()

		dvpMQ := newDvpSingleMessageMQ(msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:            "test-message-id",
				Type:          service.DvpSwapInitiatedMessage,
				BlockNumber:   42,
				SharedID:      "shared-1",
				SwapInitiated: wantData,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultDvpReceiverMock(t)
		receiver.HandleSwapInitiatedFunc = func(ctx context.Context, blockNum uint64, data *types.DvpSwapInitiatedData) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, uint64(42), blockNum)
			assert.Equal(t, wantData, data)
			return nil
		}

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(receiver.HandleSwapInitiatedCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack message on receiver HandleSwapInitiated error", func(t *testing.T) {
		ackSpy := spy.NewAck()

		dvpMQ := newDvpSingleMessageMQ(msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:   "test-message-id",
				Type: service.DvpSwapInitiatedMessage,
				SwapInitiated: &types.DvpSwapInitiatedData{
					Message: &types.DvpSwapMessage{SharedId: "shared-1"},
				},
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultDvpReceiverMock(t)
		receiver.HandleSwapInitiatedFunc = func(ctx context.Context, blockNum uint64, data *types.DvpSwapInitiatedData) error {
			return errors.New("receiver error")
		}

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(receiver.HandleSwapInitiatedCalls()))
		ackSpy.AssertNotCalled(t, "should not ack message on receiver error")
	})

	t.Run("acks nil swap initiated message", func(t *testing.T) {
		ackSpy := spy.NewAck()

		dvpMQ := newDvpSingleMessageMQ(msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:            "test-message-id",
				Type:          service.DvpSwapInitiatedMessage,
				BlockNumber:   100,
				SwapInitiated: nil,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultDvpReceiverMock(t)

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 0, len(receiver.HandleSwapInitiatedCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls receiver HandleSwapCompleted on swap completed message", func(t *testing.T) {
		wantSharedId := "shared-c"

		ackSpy := spy.NewAck()

		dvpMQ := newDvpSingleMessageMQ(msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:       "test-message-id",
				Type:     service.DvpSwapCompletedMessage,
				SharedID: wantSharedId,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultDvpReceiverMock(t)
		receiver.HandleSwapCompletedFunc = func(ctx context.Context, sharedID string) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantSharedId, sharedID)
			return nil
		}

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(receiver.HandleSwapCompletedCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack message on receiver HandleSwapCompleted error", func(t *testing.T) {
		ackSpy := spy.NewAck()

		dvpMQ := newDvpSingleMessageMQ(msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:       "test-message-id",
				Type:     service.DvpSwapCompletedMessage,
				SharedID: "shared-c",
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultDvpReceiverMock(t)
		receiver.HandleSwapCompletedFunc = func(ctx context.Context, sharedID string) error {
			return errors.New("receiver error")
		}

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(receiver.HandleSwapCompletedCalls()))
		ackSpy.AssertNotCalled(t, "should not ack message on receiver error")
	})

	t.Run("acks empty shared id for swap completed message", func(t *testing.T) {
		ackSpy := spy.NewAck()

		dvpMQ := newDvpSingleMessageMQ(msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:       "test-message-id",
				Type:     service.DvpSwapCompletedMessage,
				SharedID: "",
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultDvpReceiverMock(t)

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 0, len(receiver.HandleSwapCompletedCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls receiver HandleSwapRevert on swap cancelled message", func(t *testing.T) {
		wantSharedId := "shared-x"

		ackSpy := spy.NewAck()

		dvpMQ := newDvpSingleMessageMQ(msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:       "test-message-id",
				Type:     service.DvpSwapCancelledMessage,
				SharedID: wantSharedId,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultDvpReceiverMock(t)
		receiver.HandleSwapRevertFunc = func(ctx context.Context, sharedId string, status types.DvpSwapStatus) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantSharedId, sharedId)
			assert.Equal(t, types.DvpSwapCancelled, status)
			return nil
		}

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(receiver.HandleSwapRevertCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack message on receiver HandleSwapRevert cancelled error", func(t *testing.T) {
		ackSpy := spy.NewAck()

		dvpMQ := newDvpSingleMessageMQ(msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:       "test-message-id",
				Type:     service.DvpSwapCancelledMessage,
				SharedID: "shared-x",
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultDvpReceiverMock(t)
		receiver.HandleSwapRevertFunc = func(ctx context.Context, sharedId string, status types.DvpSwapStatus) error {
			return errors.New("receiver error")
		}

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(receiver.HandleSwapRevertCalls()))
		ackSpy.AssertNotCalled(t, "should not ack message on receiver error")
	})

	t.Run("calls receiver HandleSwapRevert on swap timed out message", func(t *testing.T) {
		wantSharedId := "shared-x"

		ackSpy := spy.NewAck()

		dvpMQ := newDvpSingleMessageMQ(msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:       "test-message-id",
				Type:     service.DvpSwapTimedOutMessage,
				SharedID: wantSharedId,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultDvpReceiverMock(t)
		receiver.HandleSwapRevertFunc = func(ctx context.Context, sharedId string, status types.DvpSwapStatus) error {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantSharedId, sharedId)
			assert.Equal(t, types.DvpSwapTimedOut, status)
			return nil
		}

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(receiver.HandleSwapRevertCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack message on receiver HandleSwapRevert timed out error", func(t *testing.T) {
		ackSpy := spy.NewAck()

		dvpMQ := newDvpSingleMessageMQ(msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:       "test-message-id",
				Type:     service.DvpSwapTimedOutMessage,
				SharedID: "shared-x",
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultDvpReceiverMock(t)
		receiver.HandleSwapRevertFunc = func(ctx context.Context, sharedId string, status types.DvpSwapStatus) error {
			return errors.New("receiver error")
		}

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(receiver.HandleSwapRevertCalls()))
		ackSpy.AssertNotCalled(t, "should not ack message on receiver error")
	})

	t.Run("acks unknown message type", func(t *testing.T) {
		ackSpy := spy.NewAck()

		dvpMQ := newDvpSingleMessageMQ(msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:          "test-message-id",
				Type:        service.DvpDestMessageType(999),
				BlockNumber: 100,
			},
			Ack: ackSpy.Fn(),
		})

		receiver := newDefaultDvpReceiverMock(t)

		svc := service.NewDvpOrchestratorWithBackoff(dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		ackSpy.AssertCalled(t)
	})
}
