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

// Regression for the DVP-dest head-of-line stall.
//
// The orchestrator processes messages in ONE sequential loop, and a handler that touches chain calls
// bind.WaitMined with the long-lived run ctx and (before the fix) no per-call bound. If a notify TX
// never mines, WaitMined blocks forever, and because the loop is single-threaded that one stuck handler
// permanently stalls ALL subsequent DVP-dest processing.
//
// The fix bounds each handler with a per-message timeout: a handler that overruns observes a cancelled
// ctx (its chain wait returns "context deadline exceeded"), returns shouldAck=false (the same no-ack/
// redeliver path a handler error already takes), and the loop continues. This test proves the loop is
// NOT head-of-line stalled by a handler that blocks until its ctx is cancelled.
func TestDvpOrchestrator_HeadOfLineStall(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("a handler that blocks past its deadline does not stall the loop and does not ack", func(t *testing.T) {
		const slowSharedID = "shared-slow"
		const fastSharedID = "shared-fast"

		var firstAcked atomic.Bool
		var secondAcked atomic.Bool

		first := msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:       "msg-slow",
				Type:     service.DvpSwapCancelledMessage,
				SharedID: slowSharedID,
			},
			Ack: func(ctx context.Context) error { firstAcked.Store(true); return nil },
		}
		second := msgqueue.Message[service.DvpDestMessage]{
			V: service.DvpDestMessage{
				ID:       "msg-fast",
				Type:     service.DvpSwapCompletedMessage,
				SharedID: fastSharedID,
			},
			Ack: func(ctx context.Context) error { secondAcked.Store(true); return nil },
		}

		// Deliver the slow message once, then keep delivering the fast one until the run ctx is cancelled.
		var firstOut atomic.Bool
		dvpMQ := &DvpDestMQMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[service.DvpDestMessage], error) {
				if firstOut.CompareAndSwap(false, true) {
					return first, nil
				}
				select {
				case <-ctx.Done():
					return msgqueue.Message[service.DvpDestMessage]{}, ctx.Err()
				default:
					return second, nil
				}
			},
		}

		var fastHandled atomic.Bool
		receiver := newDefaultDvpReceiverMock(t)

		// Slow handler: simulates bind.WaitMined on a never-mining TX — it returns only when its handler
		// ctx is cancelled (by the per-message deadline).
		receiver.HandleSwapRevertFunc = func(ctx context.Context, sharedID string, status types.DvpSwapStatus) error {
			assert.Equal(t, slowSharedID, sharedID)
			assert.Equal(t, types.DvpSwapCancelled, status)
			<-ctx.Done()
			return ctx.Err() // "context deadline exceeded" -> shouldAck=false
		}

		// Fast handler proves the loop advanced to the next message.
		receiver.HandleSwapCompletedFunc = func(ctx context.Context, sharedID string) error {
			assert.Equal(t, fastSharedID, sharedID)
			fastHandled.Store(true)
			return nil
		}

		// Small per-message timeout so the slow handler times out quickly in this unit test.
		svc := service.NewDvpOrchestratorWithBackoffAndTimeout(
			dvpMQ, receiver, newDefaultDvpOrchestratorBackoffMock(), 50*time.Millisecond,
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		done := make(chan error, 1)
		go func() { done <- svc.Run(ctx) }()

		// The fast (second) message must be handled even though the first handler blocked — i.e. the loop
		// was NOT head-of-line stalled.
		require.Eventually(t, fastHandled.Load, 4*time.Second, 10*time.Millisecond,
			"loop is head-of-line stalled: second message never handled")

		cancel()
		select {
		case err := <-done:
			require.NoError(t, err)
		case <-time.After(time.Second):
			t.Fatal("Run did not return after ctx cancel")
		}

		// The slow/cancelled handler timed out -> message NOT acked (so it is redelivered).
		assert.False(t, firstAcked.Load(), "timed-out handler must not ack its message")
		// The fast handler succeeded -> its message acked.
		assert.True(t, secondAcked.Load(), "successful handler must ack its message")
		assert.GreaterOrEqual(t, len(receiver.HandleSwapRevertCalls()), 1)
		assert.GreaterOrEqual(t, len(receiver.HandleSwapCompletedCalls()), 1)
	})
}
