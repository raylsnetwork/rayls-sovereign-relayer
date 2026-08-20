package service_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/batcher"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools/spy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrivateHubService(t *testing.T) {
	testtools.SilenceLogger()

	calldata := "example transaction calldata"
	myChainID := big.NewInt(1337)

	phEndpointAddress := common.HexToAddress("0xdeadbeef")
	msg := service.PrivateHubMessage{
		ID: "hello-world",

		MessageID: common.HexToHash("0xc0febabe"),
		From:      common.HexToAddress("0xc0fedead"),
		To:        common.HexToAddress("0xdeadc0de"),
		Data: EndpointV1.RaylsMessage{
			Payload: common.Hex2Bytes("0xFF00FF00"),
		},
	}

	t.Run("supports graceful shutdown", func(t *testing.T) {
		emptyEndpointAddress := common.Address{}
		cons := &PrivateHubConsumerMock{
			// make fetch function return an empty slice so that the service
			// will return to waiting on the ticker and won't call any of the
			// other mocks, thus we won't have to implement them.
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[service.PrivateHubMessage], error) {
				return []msgqueue.Message[service.PrivateHubMessage]{}, nil
			},
		}
		executorSvc := &BatcherMock{}
		gen := &TransactionGeneratorMock{}

		svc := service.NewPrivateHubService(time.Second, myChainID, emptyEndpointAddress, cons, gen, executorSvc)
		hasGracefulShutdown := testtools.ShutdownFixture(t, svc.Run, 10*time.Millisecond)

		assert.True(t, hasGracefulShutdown)
	})

	t.Run("doesn't skip context check on error", func(t *testing.T) {
		emptyEndpointAddress := common.Address{}
		cons := &PrivateHubConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[service.PrivateHubMessage], error) {
				return []msgqueue.Message[service.PrivateHubMessage]{}, errors.New("example-error")
			},
		}
		executorSvc := &BatcherMock{}
		gen := &TransactionGeneratorMock{}

		svc := service.NewPrivateHubService(time.Second, myChainID, emptyEndpointAddress, cons, gen, executorSvc)
		respectsContextOnError := testtools.ShutdownFixture(t, svc.Run, time.Millisecond)

		assert.True(t, respectsContextOnError)
	})

	t.Run("continues to ticker on no messages to process", func(t *testing.T) {
		emptyEndpointAddress := common.Address{}
		cons := &PrivateHubConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[service.PrivateHubMessage], error) {
				return []msgqueue.Message[service.PrivateHubMessage]{}, nil
			},
		}
		executorSvc := &BatcherMock{}
		gen := &TransactionGeneratorMock{}
		svc := service.NewPrivateHubService(time.Second, myChainID, emptyEndpointAddress, cons, gen, executorSvc)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.Nil(t, err)

		assert.Len(t, cons.FetchCalls(), 1)
		assert.Len(t, executorSvc.SendCalls(), 0)
		assert.Len(t, gen.GenerateCalls(), 0)
	})
	t.Run("generates transaction calldata for message", func(t *testing.T) {
		ackSpy := spy.NewAck()

		cons := &PrivateHubConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[service.PrivateHubMessage], error) {
				return []msgqueue.Message[service.PrivateHubMessage]{
					{
						V:   msg,
						Ack: ackSpy.Fn(),
					},
				}, nil
			},
		}
		executorSvc := &BatcherMock{
			SendFunc: func(ctx context.Context, msgs []batcher.Message) error {
				// verify that we do not ack the message
				// before publishing to the batcher
				ackSpy.AssertNotCalled(t, "should acknowledge messages after publish")

				// verify we get a message per fetched tx
				assert.Len(t, msgs, 1)

				// verify batcher.Message fields
				assert.Equal(t, msg.ID, msgs[0].ID)
				assert.Equal(t, phEndpointAddress, msgs[0].Address)
				assert.Equal(t, calldata, string(msgs[0].Calldata))

				return nil
			},
		}
		gen := &TransactionGeneratorMock{
			GenerateFunc: func(fromChainID *big.Int, fromAddress, toAddress common.Address, data EndpointV1.RaylsMessage, id common.Hash) ([]byte, error) {
				// verify that the correct parameters are passed
				assert.Equal(t, myChainID, fromChainID)
				assert.Equal(t, msg.From, fromAddress)
				assert.Equal(t, msg.To, toAddress)
				assert.Equal(t, msg.Data, data)
				assert.Equal(t, msg.MessageID, id)

				return []byte(calldata), nil
			},
		}
		svc := service.NewPrivateHubService(time.Second, myChainID, phEndpointAddress, cons, gen, executorSvc)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.Nil(t, err)

		// check that we fetch messages
		assert.Len(t, cons.FetchCalls(), 1, "didn't fetch messages")
		// check that the generate function was called
		assert.Len(t, gen.GenerateCalls(), 1, "didn't call generate function")
		// check that we acknowledged the message
		ackSpy.AssertCalled(t)
		// check that the transaction was published via the batcher
		assert.Len(t, executorSvc.SendCalls(), 1, "didn't call Send")
	})

	t.Run("doesn't acknowledge messages on error from executor", func(t *testing.T) {
		ackSpy := spy.NewAck()

		cons := &PrivateHubConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[service.PrivateHubMessage], error) {
				return []msgqueue.Message[service.PrivateHubMessage]{
					{
						V:   msg,
						Ack: ackSpy.Fn(),
					},
				}, nil
			},
		}
		executorSvc := &BatcherMock{
			SendFunc: func(ctx context.Context, msgs []batcher.Message) error {
				return errors.New("example-error")
			},
		}
		gen := &TransactionGeneratorMock{
			GenerateFunc: func(fromChainID *big.Int, fromAddress, toAddress common.Address, data EndpointV1.RaylsMessage, id common.Hash) ([]byte, error) {
				return []byte(calldata), nil
			},
		}
		svc := service.NewPrivateHubService(time.Second, myChainID, phEndpointAddress, cons, gen, executorSvc)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.Nil(t, err)

		// check that we don't ack the message on executor error
		ackSpy.AssertNotCalled(t, "should not ack on executor error")
	})
}
