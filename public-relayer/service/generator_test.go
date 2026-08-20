// Decommissioning Teleport (vanilla, atomic).

package service_test

import (
	"bytes"
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/batcher"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/RNMessageDispatcherV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/public-relayer/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func wrapMessages(msgs []service.Message) []msgqueue.Message[service.Message] {
	result := make([]msgqueue.Message[service.Message], len(msgs))
	for i, m := range msgs {
		m := m
		result[i] = msgqueue.Message[service.Message]{
			V:   m,
			Ack: func(ctx context.Context) error { return nil },
		}
	}
	return result
}

func wrapMessagesWithAckCounter(msgs []service.Message, counter *atomic.Int32) []msgqueue.Message[service.Message] {
	result := make([]msgqueue.Message[service.Message], len(msgs))
	for i, m := range msgs {
		m := m
		result[i] = msgqueue.Message[service.Message]{
			V: m,
			Ack: func(ctx context.Context) error {
				counter.Add(1)
				return nil
			},
		}
	}
	return result
}

type genFixture struct {
	consumer          *MessageConsumerMock
	generator         *TransactionGeneratorMock
	forwardBatcher    *BatcherMock
	revertBatcher     *BatcherMock
	revertRepo        *RevertSignatureRepositoryMock
	messageRecordRepo *MessageRecordRepositoryMock
	config            service.GeneratorServiceConfig
}

// newGenFixture wires a GeneratorService with mocks that all succeed and
// are empty/no-op by default. Tests override fields as needed.
func newGenFixture() *genFixture {
	f := &genFixture{
		consumer: &MessageConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[service.Message], error) {
				return nil, nil
			},
		},
		generator: &TransactionGeneratorMock{
			GenerateFunc: func(fromAddress, toAddress common.Address, message RNMessageDispatcherV1.RaylsNodeMessage, id common.Hash) ([]byte, error) {
				return []byte("default"), nil
			},
		},
		forwardBatcher: &BatcherMock{
			SendFunc: func(ctx context.Context, msgs []batcher.Message) error { return nil },
		},
		revertBatcher: &BatcherMock{
			SendFunc: func(ctx context.Context, msgs []batcher.Message) error { return nil },
		},
		revertRepo: &RevertSignatureRepositoryMock{
			BatchCreateFunc: func(ctx context.Context, sigs []service.RevertSignature) error { return nil },
			GetByIDsFunc: func(ctx context.Context, ids []string) ([]service.RevertSignature, error) {
				return nil, nil
			},
		},
		messageRecordRepo: &MessageRecordRepositoryMock{
			BatchCreateFunc: func(ctx context.Context, records []service.MessageRecord) error { return nil },
			UpdateForwardResultsFunc: func(ctx context.Context, updates []service.ForwardResultUpdate) error {
				return nil
			},
			UpdateRevertResultsFunc: func(ctx context.Context, updates []service.RevertResultUpdate) error {
				return nil
			},
		},
		config: service.GeneratorServiceConfig{
			Interval:              time.Second,
			EndpointAddress:       common.HexToAddress("0xDEAD"),
			SourceEndpointAddress: common.HexToAddress("0xBEEF"),
		},
	}
	return f
}

func (f *genFixture) newService() *service.GeneratorService {
	return service.NewGeneratorService(
		f.config,
		f.consumer,
		f.generator,
		f.forwardBatcher,
		f.revertBatcher,
		f.revertRepo,
		f.messageRecordRepo,
	)
}

func TestGeneratorService_Run(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("supports graceful shutdown", func(t *testing.T) {
		f := newGenFixture()
		f.consumer.FetchFunc = func(ctx context.Context, count int) ([]msgqueue.Message[service.Message], error) {
			return wrapMessages([]service.Message{{ID: common.Hash([32]byte{0xDE, 0xAD, 0xBE, 0xEF})}}), nil
		}

		svc := f.newService()
		assert.True(t, testtools.ShutdownFixture(t, svc.Run, time.Millisecond))
	})

	t.Run("doesn't skip context check on error", func(t *testing.T) {
		f := newGenFixture()
		f.consumer.FetchFunc = func(ctx context.Context, count int) ([]msgqueue.Message[service.Message], error) {
			return nil, errors.New("example-error")
		}

		svc := f.newService()
		assert.True(t, testtools.ShutdownFixture(t, svc.Run, time.Millisecond))
	})

	t.Run("returns to waiting on empty fetch result", func(t *testing.T) {
		f := newGenFixture()
		f.consumer.FetchFunc = func(ctx context.Context, count int) ([]msgqueue.Message[service.Message], error) {
			return nil, nil //nolint:nilnil // intentional nil return in test mock
		}
		f.generator.GenerateFunc = func(fromAddress, toAddress common.Address, message RNMessageDispatcherV1.RaylsNodeMessage, id common.Hash) ([]byte, error) {
			assert.Fail(t, "shouldn't have called generate function")
			return nil, nil //nolint:nilnil
		}
		f.forwardBatcher.SendFunc = func(ctx context.Context, msgs []batcher.Message) error {
			assert.Fail(t, "shouldn't have called forward batcher")
			return nil
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
		defer cancel()

		svc := f.newService()
		require.NoError(t, svc.Run(ctx))
		assert.Len(t, f.consumer.FetchCalls(), 1)
	})

	t.Run("generates forward + revert, persists records, publishes forward batch, acks", func(t *testing.T) {
		msgs := []service.Message{
			{
				ID: common.Hash([32]byte{0xDE, 0xAD, 0xBE, 0xEF}),
				Data: RNMessageDispatcherV1.RaylsNodeMessage{
					Payload: common.Hex2Bytes("0xc0fedEaD"),
					MessageMetadata: RNMessageDispatcherV1.RaylsNodeMessageMetadata{
						RevertPayloadData: common.Hex2Bytes("0xdEaDC0dE"),
					},
				},
			},
		}
		forwardBytes := common.Hex2Bytes("0xc0febAbE")
		revertBytes := common.Hex2Bytes("0xDeAdBeEf")

		var ackCount atomic.Int32
		f := newGenFixture()
		f.consumer.FetchFunc = func(ctx context.Context, count int) ([]msgqueue.Message[service.Message], error) {
			return wrapMessagesWithAckCounter(msgs, &ackCount), nil
		}
		f.generator.GenerateFunc = func(fromAddress, toAddress common.Address, message RNMessageDispatcherV1.RaylsNodeMessage, id common.Hash) ([]byte, error) {
			if bytes.Equal(message.Payload, msgs[0].Data.Payload) {
				return forwardBytes, nil
			} else if bytes.Equal(message.Payload, msgs[0].Data.MessageMetadata.RevertPayloadData) {
				return revertBytes, nil
			}
			return nil, errors.New("unknown payload")
		}
		f.forwardBatcher.SendFunc = func(ctx context.Context, published []batcher.Message) error {
			require.Len(t, published, 1)
			assert.Equal(t, msgs[0].ID.Hex(), published[0].ID)
			assert.Equal(t, f.config.EndpointAddress, published[0].Address)
			assert.Equal(t, forwardBytes, []byte(published[0].Calldata))
			return nil
		}
		f.revertRepo.BatchCreateFunc = func(ctx context.Context, sigs []service.RevertSignature) error {
			require.Len(t, sigs, 1)
			assert.Equal(t, msgs[0].ID.Hex(), sigs[0].ID)
			assert.Equal(t, revertBytes, sigs[0].Data)
			return nil
		}
		f.messageRecordRepo.BatchCreateFunc = func(ctx context.Context, records []service.MessageRecord) error {
			require.Len(t, records, 1)
			assert.Equal(t, msgs[0].ID.Hex(), records[0].ID)
			assert.Equal(t, service.MessageRecordStatusNew, records[0].Status)
			assert.False(t, records[0].CreatedAt.IsZero(), "CreatedAt should be set")
			assert.False(t, records[0].UpdatedAt.IsZero(), "UpdatedAt should be set")
			return nil
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
		defer cancel()

		svc := f.newService()
		require.NoError(t, svc.Run(ctx))

		assert.Len(t, f.consumer.FetchCalls(), 1)
		assert.Len(t, f.messageRecordRepo.BatchCreateCalls(), 1, "didn't create message records")
		assert.Len(t, f.revertRepo.BatchCreateCalls(), 1, "didn't persist revert signatures")
		assert.Len(t, f.forwardBatcher.SendCalls(), 1, "didn't publish forward batch")
		assert.Equal(t, int32(1), ackCount.Load(), "didn't ack message")
	})

	t.Run("does not ack when message record persist fails", func(t *testing.T) {
		msgs := []service.Message{{ID: common.Hash([32]byte{0xCA, 0xFE})}}

		var ackCount atomic.Int32
		f := newGenFixture()
		f.consumer.FetchFunc = func(ctx context.Context, count int) ([]msgqueue.Message[service.Message], error) {
			return wrapMessagesWithAckCounter(msgs, &ackCount), nil
		}
		f.messageRecordRepo.BatchCreateFunc = func(ctx context.Context, records []service.MessageRecord) error {
			return errors.New("db down")
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
		defer cancel()

		require.NoError(t, f.newService().Run(ctx))
		assert.Equal(t, int32(0), ackCount.Load(), "must not ack on persist failure")
		assert.Len(t, f.forwardBatcher.SendCalls(), 0, "must not publish on persist failure")
	})

	t.Run("does not ack when forward publish fails", func(t *testing.T) {
		msgs := []service.Message{{ID: common.Hash([32]byte{0xCA, 0xFE})}}

		var ackCount atomic.Int32
		f := newGenFixture()
		f.consumer.FetchFunc = func(ctx context.Context, count int) ([]msgqueue.Message[service.Message], error) {
			return wrapMessagesWithAckCounter(msgs, &ackCount), nil
		}
		f.forwardBatcher.SendFunc = func(ctx context.Context, published []batcher.Message) error {
			return errors.New("nats down")
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
		defer cancel()

		require.NoError(t, f.newService().Run(ctx))
		assert.Equal(t, int32(0), ackCount.Load(), "must not ack on publish failure")
	})
}

func TestGeneratorService_HandleForwardResults(t *testing.T) {
	testtools.SilenceLogger()

	okID := "0xAAA"
	failID := "0xBBB"
	okHash := common.HexToHash("0xc0febabe")
	revertData := []byte("revert-calldata")

	t.Run("success path updates records to Succeeded and publishes nothing", func(t *testing.T) {
		f := newGenFixture()
		f.messageRecordRepo.UpdateForwardResultsFunc = func(ctx context.Context, updates []service.ForwardResultUpdate) error {
			require.Len(t, updates, 1)
			assert.Equal(t, okID, updates[0].ID)
			assert.Equal(t, service.MessageRecordStatusSucceeded, updates[0].Status)
			assert.Equal(t, okHash, updates[0].Hash)
			assert.Empty(t, updates[0].Error)
			return nil
		}
		f.revertBatcher.SendFunc = func(ctx context.Context, msgs []batcher.Message) error {
			assert.Fail(t, "revert batcher must not be called on success")
			return nil
		}
		f.revertRepo.GetByIDsFunc = func(ctx context.Context, ids []string) ([]service.RevertSignature, error) {
			assert.Fail(t, "revert signatures must not be fetched on success")
			return nil, nil //nolint:nilnil
		}

		results := []types.TxResult{
			{CorrelationID: okID, Kind: types.TxResultSuccess, TxHash: okHash},
		}
		require.NoError(t, f.newService().HandleForwardResults(context.Background(), results))

		assert.Len(t, f.messageRecordRepo.UpdateForwardResultsCalls(), 1)
		assert.Len(t, f.revertBatcher.SendCalls(), 0)
	})

	t.Run("failure path updates records to Failed, fetches revert sigs, publishes via revertBatcher", func(t *testing.T) {
		f := newGenFixture()
		f.messageRecordRepo.UpdateForwardResultsFunc = func(ctx context.Context, updates []service.ForwardResultUpdate) error {
			require.Len(t, updates, 1)
			assert.Equal(t, failID, updates[0].ID)
			assert.Equal(t, service.MessageRecordStatusFailed, updates[0].Status)
			assert.Equal(t, "reverted", updates[0].Error)
			assert.Equal(t, common.Hash{}, updates[0].Hash)
			return nil
		}
		f.revertRepo.GetByIDsFunc = func(ctx context.Context, ids []string) ([]service.RevertSignature, error) {
			require.Equal(t, []string{failID}, ids)
			return []service.RevertSignature{{ID: failID, Data: revertData}}, nil
		}
		f.revertBatcher.SendFunc = func(ctx context.Context, msgs []batcher.Message) error {
			require.Len(t, msgs, 1)
			assert.Equal(t, failID, msgs[0].ID)
			assert.Equal(t, f.config.SourceEndpointAddress, msgs[0].Address)
			assert.Equal(t, revertData, []byte(msgs[0].Calldata))
			return nil
		}

		results := []types.TxResult{
			{CorrelationID: failID, Kind: types.TxResultRevert, ErrorReason: "reverted"},
		}
		require.NoError(t, f.newService().HandleForwardResults(context.Background(), results))

		assert.Len(t, f.messageRecordRepo.UpdateForwardResultsCalls(), 1)
		assert.Len(t, f.revertRepo.GetByIDsCalls(), 1)
		assert.Len(t, f.revertBatcher.SendCalls(), 1)
	})

	t.Run("mixed success and failure fans correctly", func(t *testing.T) {
		f := newGenFixture()
		f.messageRecordRepo.UpdateForwardResultsFunc = func(ctx context.Context, updates []service.ForwardResultUpdate) error {
			assert.Len(t, updates, 2)
			byID := map[string]service.ForwardResultUpdate{}
			for _, u := range updates {
				byID[u.ID] = u
			}
			assert.Equal(t, service.MessageRecordStatusSucceeded, byID[okID].Status)
			assert.Equal(t, service.MessageRecordStatusFailed, byID[failID].Status)
			return nil
		}
		f.revertRepo.GetByIDsFunc = func(ctx context.Context, ids []string) ([]service.RevertSignature, error) {
			require.Equal(t, []string{failID}, ids, "must only fetch revert sigs for failed ids")
			return []service.RevertSignature{{ID: failID, Data: revertData}}, nil
		}
		f.revertBatcher.SendFunc = func(ctx context.Context, msgs []batcher.Message) error {
			require.Len(t, msgs, 1)
			assert.Equal(t, failID, msgs[0].ID)
			return nil
		}

		results := []types.TxResult{
			{CorrelationID: okID, Kind: types.TxResultSuccess, TxHash: okHash},
			{CorrelationID: failID, Kind: types.TxResultFailed, ErrorReason: "oom"},
		}
		require.NoError(t, f.newService().HandleForwardResults(context.Background(), results))
	})

	t.Run("empty input is no-op", func(t *testing.T) {
		f := newGenFixture()
		require.NoError(t, f.newService().HandleForwardResults(context.Background(), nil))
		assert.Len(t, f.messageRecordRepo.UpdateForwardResultsCalls(), 0)
		assert.Len(t, f.revertRepo.GetByIDsCalls(), 0)
	})
}

func TestGeneratorService_HandleRevertResults(t *testing.T) {
	testtools.SilenceLogger()

	okID := "0xAAA"
	failID := "0xBBB"
	revertHash := common.HexToHash("0xfeedface")

	t.Run("success path marks row RevertSucceeded", func(t *testing.T) {
		f := newGenFixture()
		f.messageRecordRepo.UpdateRevertResultsFunc = func(ctx context.Context, updates []service.RevertResultUpdate) error {
			require.Len(t, updates, 1)
			assert.Equal(t, okID, updates[0].ID)
			assert.Equal(t, service.MessageRecordStatusRevertSucceeded, updates[0].Status)
			assert.Equal(t, revertHash, updates[0].Hash)
			return nil
		}

		results := []types.TxResult{
			{CorrelationID: okID, Kind: types.TxResultSuccess, TxHash: revertHash},
		}
		require.NoError(t, f.newService().HandleRevertResults(context.Background(), results))
		assert.Len(t, f.messageRecordRepo.UpdateRevertResultsCalls(), 1)
	})

	t.Run("failure path marks row RevertFailed with error", func(t *testing.T) {
		f := newGenFixture()
		f.messageRecordRepo.UpdateRevertResultsFunc = func(ctx context.Context, updates []service.RevertResultUpdate) error {
			require.Len(t, updates, 1)
			assert.Equal(t, failID, updates[0].ID)
			assert.Equal(t, service.MessageRecordStatusRevertFailed, updates[0].Status)
			assert.Equal(t, "boom", updates[0].Error)
			return nil
		}

		results := []types.TxResult{
			{CorrelationID: failID, Kind: types.TxResultFailed, ErrorReason: "boom"},
		}
		require.NoError(t, f.newService().HandleRevertResults(context.Background(), results))
	})

	t.Run("does not publish anything further", func(t *testing.T) {
		f := newGenFixture()
		f.forwardBatcher.SendFunc = func(ctx context.Context, msgs []batcher.Message) error {
			assert.Fail(t, "forward batcher must not fire in revert callback")
			return nil
		}
		f.revertBatcher.SendFunc = func(ctx context.Context, msgs []batcher.Message) error {
			assert.Fail(t, "revert batcher must not fire in revert callback")
			return nil
		}

		results := []types.TxResult{
			{CorrelationID: okID, Kind: types.TxResultSuccess, TxHash: revertHash},
			{CorrelationID: failID, Kind: types.TxResultRevert, ErrorReason: "bad"},
		}
		require.NoError(t, f.newService().HandleRevertResults(context.Background(), results))
	})

	t.Run("empty input is no-op", func(t *testing.T) {
		f := newGenFixture()
		require.NoError(t, f.newService().HandleRevertResults(context.Background(), nil))
		assert.Len(t, f.messageRecordRepo.UpdateRevertResultsCalls(), 0)
	})
}
