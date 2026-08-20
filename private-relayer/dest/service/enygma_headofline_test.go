package service_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Regression for the Enygma-dest head-of-line stall (companion to the DVP-dest fix).
//
// The orchestrator processes messages in ONE sequential loop, and a handler that touches chain calls
// bind.WaitMined with the long-lived run ctx and (before the fix) no per-call bound. A handler that
// never returns permanently stalls ALL subsequent processing. The fix bounds each handler with a
// per-message timeout: a stuck handler observes a cancelled ctx, returns shouldAck=false (the same
// no-ack/redeliver path a handler error already takes), and the loop continues.
func TestEnygmaOrchestrator_HeadOfLineStall(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("a handler that blocks past its deadline does not stall the loop and does not ack", func(t *testing.T) {
		const slowResourceID = "resource-slow"
		const fastResourceID = "resource-fast"

		var firstAcked atomic.Bool
		var secondAcked atomic.Bool

		first := msgqueue.Message[service.EnygmaDestMessage]{
			V: service.EnygmaDestMessage{
				ID:            "msg-slow",
				Type:          service.EnygmaTransferBatchMessage,
				TransferBatch: &types.EnygmaTransferBatch{ResourceId: slowResourceID},
			},
			Ack: func(ctx context.Context) error { firstAcked.Store(true); return nil },
		}
		second := msgqueue.Message[service.EnygmaDestMessage]{
			V: service.EnygmaDestMessage{
				ID:            "msg-fast",
				Type:          service.EnygmaTransferBatchMessage,
				TransferBatch: &types.EnygmaTransferBatch{ResourceId: fastResourceID},
			},
			Ack: func(ctx context.Context) error { secondAcked.Store(true); return nil },
		}

		// Deliver the slow message once, then keep delivering the fast one until the run ctx is cancelled.
		var firstOut atomic.Bool
		enygmaMQ := &EnygmaDestMQMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[service.EnygmaDestMessage], error) {
				if firstOut.CompareAndSwap(false, true) {
					return first, nil
				}
				select {
				case <-ctx.Done():
					return msgqueue.Message[service.EnygmaDestMessage]{}, ctx.Err()
				default:
					return second, nil
				}
			},
		}

		var fastHandled atomic.Bool
		receiver := newDefaultReceiverMock(t)
		receiver.HandleEnygmaCrossTransferFunc = func(ctx context.Context, batch *types.EnygmaTransferBatch) error {
			switch batch.ResourceId {
			case slowResourceID:
				// Simulates bind.WaitMined on a never-mining TX — returns only when the handler ctx is
				// cancelled (by the per-message deadline).
				<-ctx.Done()
				return ctx.Err()
			case fastResourceID:
				fastHandled.Store(true)
				return nil
			default:
				assert.Fail(t, "unexpected resourceId", batch.ResourceId)
				return nil
			}
		}

		// Small per-message timeout so the slow handler times out quickly in this unit test.
		svc := service.NewEnygmaOrchestratorWithBackoffAndTimeout(
			enygmaMQ, receiver, newDefaultCheckpointRepoMock(t), newDefaultBackoffMock(), 50*time.Millisecond,
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		done := make(chan error, 1)
		go func() { done <- svc.Run(ctx) }()

		require.Eventually(t, fastHandled.Load, 4*time.Second, 10*time.Millisecond,
			"loop is head-of-line stalled: second message never handled")

		cancel()
		select {
		case err := <-done:
			require.NoError(t, err)
		case <-time.After(time.Second):
			t.Fatal("Run did not return after ctx cancel")
		}

		assert.False(t, firstAcked.Load(), "timed-out handler must not ack its message")
		assert.True(t, secondAcked.Load(), "successful handler must ack its message")
	})
}
