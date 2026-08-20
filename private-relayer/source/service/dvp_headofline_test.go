package service_test

import (
	"context"
	"encoding/json"
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Regression for the source DVP head-of-line stall (companion to the dest fixes).
//
// The source orchestrator processes messages in ONE sequential loop, and a handler that touches chain
// calls bind.WaitMined with the long-lived run ctx and (before the fix) no per-call bound. A handler
// that never returns permanently stalls ALL subsequent processing. The fix bounds each message's
// processing with a per-message timeout: a stuck handler observes a cancelled ctx, the message is not
// acked (redelivered, the same path a handler error already takes), and the loop continues.
func TestSourceDvpOrchestrator_HeadOfLineStall(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("a handler that blocks past its deadline does not stall the loop and does not ack", func(t *testing.T) {
		const slowSharedID = "shared-slow"
		const fastSharedID = "shared-fast"

		mkMsg := func(sharedID string, ack func(context.Context) error) msgqueue.Message[service.DvpSerializedEventBatch] {
			events := []service.DvpEnygmaSwapERC721{{
				SharedId:      sharedID,
				DestChainId:   big.NewInt(100),
				From:          common.HexToAddress("0xc001babe"),
				ResourceId:    "resource-id",
				EnygmaAmount:  big.NewInt(1000),
				NftResourceId: "nft-resource-id",
				NftId:         "42",
			}}
			serialized, _ := json.Marshal(events)
			return msgqueue.Message[service.DvpSerializedEventBatch]{
				V: service.DvpSerializedEventBatch{
					BlockNumber:      100,
					Type:             service.DvpEnygmaSwapERC721Event,
					SerializedEvents: serialized,
				},
				Ack: ack,
			}
		}

		var firstAcked atomic.Bool
		var secondAcked atomic.Bool
		first := mkMsg(slowSharedID, func(ctx context.Context) error { firstAcked.Store(true); return nil })
		second := mkMsg(fastSharedID, func(ctx context.Context) error { secondAcked.Store(true); return nil })

		// Deliver the slow message once, then keep delivering the fast one until the run ctx is cancelled.
		var firstOut atomic.Bool
		dvpMQ := &DvpBatchMQMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[service.DvpSerializedEventBatch], error) {
				if firstOut.CompareAndSwap(false, true) {
					return first, nil
				}
				select {
				case <-ctx.Done():
					return msgqueue.Message[service.DvpSerializedEventBatch]{}, ctx.Err()
				default:
					return second, nil
				}
			},
		}

		var fastHandled atomic.Bool
		initiator := newDefaultDvpInitiatorMock(t)
		initiator.HandleEnygmaSwapERC721Func = func(ctx context.Context, sharedId string, toChainId *big.Int, from common.Address, enygmaResourceId string, enygmaAmount *big.Int, nftResourceId string, nftId string, txHash string, txBlockNumber *big.Int, validityTime uint64) error {
			switch sharedId {
			case slowSharedID:
				// Simulates bind.WaitMined on a never-mining TX — returns only when the per-message ctx is
				// cancelled (by the per-message deadline).
				<-ctx.Done()
				return ctx.Err()
			case fastSharedID:
				fastHandled.Store(true)
				return nil
			default:
				assert.Fail(t, "unexpected sharedId", sharedId)
				return nil
			}
		}

		svc := service.NewDvpOrchestratorWithDepsAndTimeout(dvpMQ, initiator, 50*time.Millisecond)

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
