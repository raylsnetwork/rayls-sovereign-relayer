package logparser_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/msgqueue"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/logparser"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/logparser/testdata"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/logrouter"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools/fake"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools/spy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newEndpointSingleMessageMQ(msg msgqueue.Message[logrouter.Block]) *EndpointMQMock {
	return &EndpointMQMock{
		NextFunc: fake.NextMQ(msg),
	}
}

func newDefaultEndpointBackoffMock() *EndpointBackoffMock {
	return &EndpointBackoffMock{
		DoFunc: func(ctx context.Context, maxAttempts int, fn func() error) error {
			return fn()
		},
	}
}

func TestEndpointParser(t *testing.T) {
	testtools.SilenceLogger()
	hasGracefulShutdown := false
	t.Run("supports graceful shutdown", func(t *testing.T) {
		endpointMQ := &EndpointMQMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[logrouter.Block], error) {
				<-ctx.Done()
				return msgqueue.Message[logrouter.Block]{}, context.Canceled
			},
		}
		privateHubMQ := &PrivateHubMQMock{
			PushBatchFunc: func(ctx context.Context, msgs []service.PrivateHubMessage) error {
				return nil
			},
		}
		svc := logparser.NewEndpointLogParser(big.NewInt(1), endpointMQ, privateHubMQ)

		hasGracefulShutdown = testtools.ShutdownFixture(t, svc.Run, time.Millisecond)
		assert.True(t, hasGracefulShutdown)
	})

	// Require graceful shutdown to continue with the rest of the tests.
	require.True(t, hasGracefulShutdown)
	t.Run("converts logs to service type and pushes batch to mq", func(t *testing.T) {
		wantMessageID := common.HexToHash("0xc0febabe")
		wantFrom := common.HexToAddress("0xdeadc0de")
		wantTo := common.HexToAddress("0xcaffe19")

		log := testdata.NewEndpointV1MessageDispatchedLogWith(
			testdata.WithMessageDispatchedMessageID(wantMessageID),
			testdata.WithMessageDispatchedFrom(wantFrom),
			testdata.WithMessageDispatchedTo(wantTo),
			testdata.WithMessageDispatchedToChainId(big.NewInt(1)),
		)

		block := logrouter.Block{
			Number: 1337,
			Logs:   []ethTypes.Log{log},
		}

		ackSpy := spy.NewAck()

		endpointMQ := newEndpointSingleMessageMQ(msgqueue.Message[logrouter.Block]{
			V:   block,
			Ack: ackSpy.Fn(),
		})
		privateHubMQ := &PrivateHubMQMock{
			PushBatchFunc: func(ctx context.Context, msgs []service.PrivateHubMessage) error {
				ackSpy.AssertNotCalled(t, "should push before acking")

				require.Len(t, msgs, 1)
				assert.Equal(t, wantMessageID, msgs[0].MessageID)
				assert.Equal(t, wantFrom, msgs[0].From)
				assert.Equal(t, wantTo, msgs[0].To)
				assert.NotEmpty(t, msgs[0].ID)
				return nil
			},
		}
		svc := logparser.NewEndpointLogParserWithBackoff(
			big.NewInt(1),
			endpointMQ,
			privateHubMQ,
			newDefaultEndpointBackoffMock(),
		)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(privateHubMQ.PushBatchCalls()), "should push batch to private hub MQ exactly once")
		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack block when PushBatch fails after backoff", func(t *testing.T) {
		log := testdata.NewEndpointV1MessageDispatchedLogWith(
			testdata.WithMessageDispatchedToChainId(big.NewInt(1)),
		)
		block := logrouter.Block{
			Number: 1337,
			Logs:   []ethTypes.Log{log},
		}

		ackSpy := spy.NewAck()
		endpointMQ := newEndpointSingleMessageMQ(msgqueue.Message[logrouter.Block]{
			V:   block,
			Ack: ackSpy.Fn(),
		})
		privateHubMQ := &PrivateHubMQMock{
			PushBatchFunc: func(ctx context.Context, msgs []service.PrivateHubMessage) error {
				return errors.New("push failed")
			},
		}

		backoff := &EndpointBackoffMock{
			DoFunc: func(ctx context.Context, maxAttempts int, fn func() error) error {
				return fn()
			},
		}

		svc := logparser.NewEndpointLogParserWithBackoff(big.NewInt(1), endpointMQ, privateHubMQ, backoff)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(privateHubMQ.PushBatchCalls()), "should attempt PushBatch")
		ackSpy.AssertNotCalled(t, "should not ack block when PushBatch fails, allowing MQ redelivery")
	})
}
