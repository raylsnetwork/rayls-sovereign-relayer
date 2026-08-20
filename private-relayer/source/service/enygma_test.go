package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-sovereign-relayer/msgqueue"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/service/testdata"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools/fake"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools/spy"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

// Test Helpers

func newSingleMessageMQ(msg msgqueue.Message[service.EnygmaSerializedEvent]) *EnygmaBatchMQMock {
	return &EnygmaBatchMQMock{
		FetchFunc: fake.FetchMQ([]msgqueue.Message[service.EnygmaSerializedEvent]{msg}),
	}
}

func newMultiMessageMQ(msgs []msgqueue.Message[service.EnygmaSerializedEvent]) *EnygmaBatchMQMock {
	return &EnygmaBatchMQMock{
		FetchFunc: fake.FetchMQ(msgs),
	}
}

func newDefaultInitiatorMock(t *testing.T) *EnygmaInitiatorMock {
	return &EnygmaInitiatorMock{
		HandleEnygmaCreationFunc: func(context.Context, string, string, uint64, *big.Int) (uint64, error) {
			assert.Fail(t, "shouldn't call enygma creation handler")
			return 0, nil
		},
		HandleEnygmaSupplyUpdatesFunc: func(context.Context, string, string, uint64, types.EnygmaSupplyUpdate) (uint64, error) {
			assert.Fail(t, "shouldn't call enygma supply update handler")
			return 0, nil
		},
		HandleEnygmaDepositFunc: func(context.Context, string, uint64, string, [32]byte, common.Address, *big.Int, common.Hash) (uint64, error) {
			assert.Fail(t, "shouldn't call enygma deposit handler")
			return 0, nil
		},
		HandleEnygmaWithdrawalFunc: func(context.Context, string, uint64, string, [32]byte, common.Address, *big.Int, common.Hash) (uint64, error) {
			assert.Fail(t, "shouldn't call enygma withdrawal handler")
			return 0, nil
		},
		HandleEnygmaCrossTransferFunc: func(context.Context, string, string, uint64, map[string][]*types.EnygmaTransferBatchTx) (uint64, error) {
			assert.Fail(t, "shouldn't call enygma cross transfer handler")
			return 0, nil
		},
	}
}

func newDefaultReverterMock(t *testing.T) *EnygmaReverterMock {
	return &EnygmaReverterMock{
		RevertEnygmaSupplyUpdateFunc: func(ctx context.Context, resourceId string, supplyUpdateEvents []service.EnygmaSupplyUpdate) error {
			assert.Fail(t, "shouldn't call enygma revert supply update handler")
			return nil
		},
		RevertEnygmaDepositFunc: func(ctx context.Context, resourceId string, referenceId [32]byte, from common.Address, amount *big.Int) error {
			assert.Fail(t, "shouldn't call enygma revert deposit handler")
			return nil
		},
		RevertEnygmaTransferFunc: func(parentCtx context.Context, resourceId string, txsByChainID map[string][]*types.EnygmaTransferBatchTx) error {
			assert.Fail(t, "shouldn't call enygma revert transfer handler")
			return nil
		},
	}
}

func newDefaultBatcherMock(t *testing.T) *EnygmaBatcherMock {
	return &EnygmaBatcherMock{
		GroupTransfersByChainIDFunc: func([]service.EnygmaTransferTx) map[string][]*types.EnygmaTransferBatchTx {
			assert.Fail(t, "should not call group by chain ID")
			return nil
		},
		BatchTransfersFunc: func(_ context.Context, _ map[string][]*types.EnygmaTransferBatchTx) ([]map[string][]*types.EnygmaTransferBatchTx, error) {
			assert.Fail(t, "should not call batch transfers")
			return nil, nil //nolint:nilnil // intentional nil return in test mock
		},
	}
}

func newDefaultFinalizerMock(t *testing.T) *EnygmaFinalizationServiceMock {
	return &EnygmaFinalizationServiceMock{
		ExecuteFinalizationFunc: func(context.Context, string, uint64, string) error {
			return nil // Default: do nothing, succeed
		},
	}
}

func TestEnygmaOrchestrator(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("supports graceful shutdown", func(t *testing.T) {
		enygmaMQ := &EnygmaBatchMQMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[service.EnygmaSerializedEvent], error) {
				<-ctx.Done()
				return nil, context.Canceled
			},
		}
		blockWaiter := &EnygmaBlockWaiterServiceMock{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				assert.EqualValues(t, 0, targetBlockNumber)
				return 1337, nil
			},
		}
		batcher := newDefaultBatcherMock(t)
		initiator := newDefaultInitiatorMock(t)
		reverter := newDefaultReverterMock(t)
		finalizer := newDefaultFinalizerMock(t)

		svc := service.NewEnygmaOrchestrator(
			time.Millisecond,
			enygmaMQ,
			blockWaiter,
			batcher,
			initiator,
			reverter,
			finalizer,
			2,
		)
		hasGracefulShutdown := testtools.ShutdownFixture(t, svc.Run, time.Second)

		assert.True(t, hasGracefulShutdown)
	})

	t.Run("calls initiator handler on creation event", func(t *testing.T) {
		wantResourceID := "example-resource-id"
		wantBlockNumber := uint64(1337)
		wantInitialSupply := big.NewInt(42)

		ackSpy := spy.NewAck()

		event := service.EnygmaCreation{
			InitialSupply: wantInitialSupply,
		}
		serializedEvent, _ := json.Marshal(event)

		enygmaMQ := newSingleMessageMQ(msgqueue.Message[service.EnygmaSerializedEvent]{
			V: service.EnygmaSerializedEvent{
				Id:              "test-creation-1",
				BlockNumber:     67,
				LogIndex:        0,
				TxHash:          common.Hash{},
				Type:            service.EnygmaCreationEvent,
				ResourceID:      wantResourceID,
				SerializedEvent: serializedEvent,
			},
			Ack: ackSpy.Fn(),
		})
		blockWaiter := &EnygmaBlockWaiterServiceMock{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				assert.EqualValues(t, 0, targetBlockNumber)
				return wantBlockNumber, nil
			},
		}

		batcher := newDefaultBatcherMock(t)
		initiator := newDefaultInitiatorMock(t)
		initiator.HandleEnygmaCreationFunc = func(ctx context.Context, _ string, resourceId string, blockNumber uint64, initialSupply *big.Int) (uint64, error) {
			ackSpy.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantResourceID, resourceId)
			assert.Equal(t, wantBlockNumber, blockNumber)
			assert.Equal(t, wantInitialSupply, initialSupply)
			return wantBlockNumber, nil
		}
		reverter := newDefaultReverterMock(t)
		finalizer := newDefaultFinalizerMock(t)

		svc := service.NewEnygmaOrchestrator(
			time.Millisecond,
			enygmaMQ,
			blockWaiter,
			batcher,
			initiator,
			reverter,
			finalizer,
			2,
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(initiator.HandleEnygmaCreationCalls()))
		ackSpy.AssertCalled(t)
	})

	t.Run("calls initiator handler on supply updates", func(t *testing.T) {
		wantResourceID := "example-resource-id"
		wantEventType := types.EnygmaBurn
		wantAmount := big.NewInt(10)

		wantInitiatorBlockNumber := uint64(1337)
		wantFinalizerBlockNumber := uint64(1340)

		ackSpy1 := spy.NewAck()
		ackSpy2 := spy.NewAck()

		event1 := service.EnygmaSupplyUpdate{Amount: big.NewInt(10)}
		event2 := service.EnygmaSupplyUpdate{Amount: big.NewInt(-20)}
		serializedEvent1, _ := json.Marshal(event1)
		serializedEvent2, _ := json.Marshal(event2)

		enygmaMQ := newMultiMessageMQ([]msgqueue.Message[service.EnygmaSerializedEvent]{
			{
				V: service.EnygmaSerializedEvent{
					Id:              "test-supply-1",
					BlockNumber:     67,
					LogIndex:        0,
					TxHash:          common.Hash{},
					ResourceID:      wantResourceID,
					Type:            service.EnygmaSupplyUpdateEvent,
					SerializedEvent: serializedEvent1,
				},
				Ack: ackSpy1.Fn(),
			},
			{
				V: service.EnygmaSerializedEvent{
					Id:              "test-supply-2",
					BlockNumber:     67,
					LogIndex:        1,
					TxHash:          common.Hash{},
					ResourceID:      wantResourceID,
					Type:            service.EnygmaSupplyUpdateEvent,
					SerializedEvent: serializedEvent2,
				},
				Ack: ackSpy2.Fn(),
			},
		})
		blockWaiter := &EnygmaBlockWaiterServiceMock{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				assert.EqualValues(t, 0, targetBlockNumber)
				return wantInitiatorBlockNumber, nil
			},
		}

		batcher := newDefaultBatcherMock(t)
		initiator := newDefaultInitiatorMock(t)
		initiator.HandleEnygmaSupplyUpdatesFunc = func(ctx context.Context, _ string, resourceId string, blockNumber uint64, batch types.EnygmaSupplyUpdate) (uint64, error) {
			ackSpy1.AssertNotCalled(t, "should ack message AFTER handling it")
			ackSpy2.AssertNotCalled(t, "should ack message AFTER handling it")
			assert.Equal(t, wantResourceID, resourceId)
			assert.Equal(t, wantInitiatorBlockNumber, blockNumber)
			assert.Equal(t, wantEventType, batch.Type)
			assert.Equal(t, wantAmount, batch.Amount)
			return wantFinalizerBlockNumber, nil
		}
		reverter := newDefaultReverterMock(t)

		finalizer := newDefaultFinalizerMock(t)
		finalizer.ExecuteFinalizationFunc = func(ctx context.Context, _ string, blockNumber uint64, resourceID string) error {
			// In the new implementation, messages are acked BEFORE finalization (after successful handling)
			ackSpy1.AssertCalled(t)
			ackSpy2.AssertCalled(t)
			assert.Equal(t, wantResourceID, resourceID)
			assert.Equal(t, wantFinalizerBlockNumber, blockNumber)
			return nil
		}

		svc := service.NewEnygmaOrchestrator(
			time.Millisecond,
			enygmaMQ,
			blockWaiter,
			batcher,
			initiator,
			reverter,
			finalizer,
			2,
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(initiator.HandleEnygmaSupplyUpdatesCalls()))
		assert.Equal(t, 1, len(finalizer.ExecuteFinalizationCalls()))
		ackSpy1.AssertCalled(t)
		ackSpy2.AssertCalled(t)
	})

	t.Run("calls revert handler on error from initiator and does not call finalizer", func(t *testing.T) {
		wantResourceID := "example-resource-id"
		initiatorBlockNumber := uint64(1337)

		ackSpy1 := spy.NewAck()
		ackSpy2 := spy.NewAck()

		wantEvent1 := service.EnygmaSupplyUpdate{Amount: big.NewInt(10)}
		wantEvent2 := service.EnygmaSupplyUpdate{Amount: big.NewInt(-20)}
		wantEvents := []service.EnygmaSupplyUpdate{wantEvent1, wantEvent2}
		serializedEvent1, _ := json.Marshal(wantEvent1)
		serializedEvent2, _ := json.Marshal(wantEvent2)

		enygmaMQ := newMultiMessageMQ([]msgqueue.Message[service.EnygmaSerializedEvent]{
			{
				V: service.EnygmaSerializedEvent{
					Id:              "test-supply-revert-1",
					BlockNumber:     67,
					LogIndex:        0,
					TxHash:          common.Hash{},
					ResourceID:      wantResourceID,
					Type:            service.EnygmaSupplyUpdateEvent,
					SerializedEvent: serializedEvent1,
				},
				Ack: ackSpy1.Fn(),
			},
			{
				V: service.EnygmaSerializedEvent{
					Id:              "test-supply-revert-2",
					BlockNumber:     67,
					LogIndex:        1,
					TxHash:          common.Hash{},
					ResourceID:      wantResourceID,
					Type:            service.EnygmaSupplyUpdateEvent,
					SerializedEvent: serializedEvent2,
				},
				Ack: ackSpy2.Fn(),
			},
		})
		blockWaiter := &EnygmaBlockWaiterServiceMock{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				assert.EqualValues(t, 0, targetBlockNumber)
				return initiatorBlockNumber, nil
			},
		}

		batcher := newDefaultBatcherMock(t)
		initiator := newDefaultInitiatorMock(t)
		initiator.HandleEnygmaSupplyUpdatesFunc = func(ctx context.Context, _ string, resourceId string, blockNumber uint64, batch types.EnygmaSupplyUpdate) (uint64, error) {
			return 0, errors.New("example-error")
		}
		reverter := newDefaultReverterMock(t)
		reverter.RevertEnygmaSupplyUpdateFunc = func(ctx context.Context, resourceId string, supplyUpdateEvents []service.EnygmaSupplyUpdate) error {
			assert.Equal(t, wantResourceID, resourceId)
			assert.Equal(t, wantEvents, supplyUpdateEvents)
			return nil
		}

		finalizer := newDefaultFinalizerMock(t)
		finalizer.ExecuteFinalizationFunc = func(ctx context.Context, _ string, blockNumber uint64, resourceID string) error {
			assert.Fail(t, "shouldn't call finalize on error from initiator")
			return nil
		}

		svc := service.NewEnygmaOrchestrator(
			time.Millisecond,
			enygmaMQ,
			blockWaiter,
			batcher,
			initiator,
			reverter,
			finalizer,
			2,
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, len(initiator.HandleEnygmaSupplyUpdatesCalls()))
		assert.Equal(t, 1, len(reverter.RevertEnygmaSupplyUpdateCalls()))
		ackSpy1.AssertCalled(t)
		ackSpy2.AssertCalled(t)
	})

	t.Run("calls initiator handle on deposit events", func(t *testing.T) {
		wantResourceID := "example-resource-id"

		calledWaitForBlockCount := 0
		firstLatestBlockNumber := uint64(1337)
		secondLatestBlockNumber := uint64(1341)

		calledDepositCount := 0
		firstDepositBlockNumber := uint64(1340)
		secondDepositBlockNumber := uint64(1347)

		ackSpy1 := spy.NewAck()
		ackSpy2 := spy.NewAck()

		event1 := service.EnygmaDepositToDvp{
			ReferenceId: common.HexToHash("0xc001a1dd"),
			From:        common.HexToAddress("0xc001babe"),
			Amount:      big.NewInt(100),
		}
		event2 := service.EnygmaDepositToDvp{
			ReferenceId: common.HexToHash("0xdeadbeef"),
			From:        common.HexToAddress("0xdeadc0de"),
			Amount:      big.NewInt(200),
		}
		events := []service.EnygmaDepositToDvp{event1, event2}
		serializedEvent1, _ := json.Marshal(event1)
		serializedEvent2, _ := json.Marshal(event2)

		enygmaMQ := newMultiMessageMQ([]msgqueue.Message[service.EnygmaSerializedEvent]{
			{
				V: service.EnygmaSerializedEvent{
					Id:              "test-deposit-1",
					BlockNumber:     67,
					LogIndex:        0,
					TxHash:          common.Hash{},
					ResourceID:      wantResourceID,
					Type:            service.EnygmaDepositEvent,
					SerializedEvent: serializedEvent1,
				},
				Ack: ackSpy1.Fn(),
			},
			{
				V: service.EnygmaSerializedEvent{
					Id:              "test-deposit-2",
					BlockNumber:     67,
					LogIndex:        1,
					TxHash:          common.Hash{},
					ResourceID:      wantResourceID,
					Type:            service.EnygmaDepositEvent,
					SerializedEvent: serializedEvent2,
				},
				Ack: ackSpy2.Fn(),
			},
		})
		// New implementation calls WaitForBlock 3 times:
		// 1. Initial call in Run() with targetBlockNumber=0
		// 2. In handleEnygmaDeposit for deposit 1 with targetBlockNumber=latestBlock (1337)
		// 3. In handleEnygmaDeposit for deposit 2 with targetBlockNumber=firstDepositBlockNumber (1340)
		blockWaiter := &EnygmaBlockWaiterServiceMock{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				calledWaitForBlockCount += 1

				switch calledWaitForBlockCount {
				case 1:
					assert.EqualValues(t, 0, targetBlockNumber)
					return firstLatestBlockNumber, nil
				case 2:
					assert.EqualValues(t, firstLatestBlockNumber, targetBlockNumber)
					return firstLatestBlockNumber, nil
				case 3:
					assert.EqualValues(t, firstDepositBlockNumber, targetBlockNumber)
					return secondLatestBlockNumber, nil
				default:
					assert.Fail(t, "called block waiter more times than required")
					return 0, nil
				}
			},
		}

		batcher := newDefaultBatcherMock(t)
		initiator := newDefaultInitiatorMock(t)
		initiator.HandleEnygmaDepositFunc = func(ctx context.Context, _ string, blockNumber uint64, resourceID string, referenceID [32]byte, from common.Address, amount *big.Int, txHash common.Hash) (uint64, error) {
			calledDepositCount += 1

			assert.Equal(t, wantResourceID, resourceID)
			switch calledDepositCount {
			case 1:
				assert.Equal(t, firstLatestBlockNumber, blockNumber)
				assert.Equal(t, events[0].ReferenceId, referenceID)
				assert.Equal(t, events[0].From, from)
				assert.Equal(t, events[0].Amount, amount)
				return firstDepositBlockNumber, nil
			case 2:
				assert.Equal(t, secondLatestBlockNumber, blockNumber)
				assert.Equal(t, events[1].ReferenceId, referenceID)
				assert.Equal(t, events[1].From, from)
				assert.Equal(t, events[1].Amount, amount)
				return secondDepositBlockNumber, nil
			default:
				assert.Fail(t, "called deposit too many times")
				return 0, nil
			}
		}
		reverter := newDefaultReverterMock(t)

		finalizer := newDefaultFinalizerMock(t)
		finalizer.ExecuteFinalizationFunc = func(ctx context.Context, _ string, blockNumber uint64, resourceID string) error {
			assert.Equal(t, secondDepositBlockNumber, blockNumber)
			assert.Equal(t, wantResourceID, resourceID)
			return nil
		}

		svc := service.NewEnygmaOrchestrator(
			time.Millisecond,
			enygmaMQ,
			blockWaiter,
			batcher,
			initiator,
			reverter,
			finalizer,
			2,
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 2, len(initiator.HandleEnygmaDepositCalls()))
		assert.Equal(t, 1, len(finalizer.ExecuteFinalizationCalls()))
		ackSpy1.AssertCalled(t)
		ackSpy2.AssertCalled(t)
	})

	t.Run("calls initiator handle on withdrawal events", func(t *testing.T) {
		wantResourceID := "example-resource-id"

		calledWaitForBlockCount := 0
		firstLatestBlockNumber := uint64(1337)
		secondLatestBlockNumber := uint64(1341)

		calledWithdrawalCount := 0
		firstWithdrawalBlockNumber := uint64(1340)
		secondWithdrawalBlockNumber := uint64(1347)

		ackSpy1 := spy.NewAck()
		ackSpy2 := spy.NewAck()

		event1 := service.EnygmaWithdrawFromDvp{
			ReferenceId: common.HexToHash("0xc001a1dd"),
			To:          common.HexToAddress("0xc001babe"),
			Amount:      big.NewInt(100),
		}
		event2 := service.EnygmaWithdrawFromDvp{
			ReferenceId: common.HexToHash("0xdeadbeef"),
			To:          common.HexToAddress("0xdeadc0de"),
			Amount:      big.NewInt(200),
		}
		events := []service.EnygmaWithdrawFromDvp{event1, event2}
		serializedEvent1, _ := json.Marshal(event1)
		serializedEvent2, _ := json.Marshal(event2)

		enygmaMQ := newMultiMessageMQ([]msgqueue.Message[service.EnygmaSerializedEvent]{
			{
				V: service.EnygmaSerializedEvent{
					Id:              "test-withdraw-1",
					BlockNumber:     67,
					LogIndex:        0,
					TxHash:          common.Hash{},
					ResourceID:      wantResourceID,
					Type:            service.EnygmaWithdrawEvent,
					SerializedEvent: serializedEvent1,
				},
				Ack: ackSpy1.Fn(),
			},
			{
				V: service.EnygmaSerializedEvent{
					Id:              "test-withdraw-2",
					BlockNumber:     67,
					LogIndex:        1,
					TxHash:          common.Hash{},
					ResourceID:      wantResourceID,
					Type:            service.EnygmaWithdrawEvent,
					SerializedEvent: serializedEvent2,
				},
				Ack: ackSpy2.Fn(),
			},
		})
		// New implementation calls WaitForBlock 3 times:
		// 1. Initial call in Run() with targetBlockNumber=0
		// 2. In handleEnygmaWithdrawal for withdrawal 1 with targetBlockNumber=latestBlock (1337)
		// 3. In handleEnygmaWithdrawal for withdrawal 2 with targetBlockNumber=firstWithdrawalBlockNumber (1340)
		blockWaiter := &EnygmaBlockWaiterServiceMock{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				calledWaitForBlockCount += 1

				switch calledWaitForBlockCount {
				case 1:
					assert.EqualValues(t, 0, targetBlockNumber)
					return firstLatestBlockNumber, nil
				case 2:
					assert.EqualValues(t, firstLatestBlockNumber, targetBlockNumber)
					return firstLatestBlockNumber, nil
				case 3:
					assert.EqualValues(t, firstWithdrawalBlockNumber, targetBlockNumber)
					return secondLatestBlockNumber, nil
				default:
					assert.Fail(t, "called block waiter more times than required")
					return 0, nil
				}
			},
		}

		batcher := newDefaultBatcherMock(t)
		initiator := newDefaultInitiatorMock(t)
		initiator.HandleEnygmaWithdrawalFunc = func(ctx context.Context, _ string, blockNumber uint64, resourceID string, referenceID [32]byte, to common.Address, amount *big.Int, txHash common.Hash) (uint64, error) {
			calledWithdrawalCount += 1

			assert.Equal(t, wantResourceID, resourceID)
			switch calledWithdrawalCount {
			case 1:
				assert.Equal(t, firstLatestBlockNumber, blockNumber)
				assert.Equal(t, events[0].ReferenceId, referenceID)
				assert.Equal(t, events[0].To, to)
				assert.Equal(t, events[0].Amount, amount)
				return firstWithdrawalBlockNumber, nil
			case 2:
				assert.Equal(t, secondLatestBlockNumber, blockNumber)
				assert.Equal(t, events[1].ReferenceId, referenceID)
				assert.Equal(t, events[1].To, to)
				assert.Equal(t, events[1].Amount, amount)
				return secondWithdrawalBlockNumber, nil
			default:
				assert.Fail(t, "called withdrawal too many times")
				return 0, nil
			}
		}
		reverter := newDefaultReverterMock(t)

		finalizer := newDefaultFinalizerMock(t)
		finalizer.ExecuteFinalizationFunc = func(ctx context.Context, _ string, blockNumber uint64, resourceID string) error {
			assert.Equal(t, secondWithdrawalBlockNumber, blockNumber)
			assert.Equal(t, wantResourceID, resourceID)
			return nil
		}

		svc := service.NewEnygmaOrchestrator(
			time.Millisecond,
			enygmaMQ,
			blockWaiter,
			batcher,
			initiator,
			reverter,
			finalizer,
			2,
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 2, len(initiator.HandleEnygmaWithdrawalCalls()))
		assert.Equal(t, 1, len(finalizer.ExecuteFinalizationCalls()))
		ackSpy1.AssertCalled(t)
		ackSpy2.AssertCalled(t)
	})

	t.Run("calls initiator handler on transfer events", func(t *testing.T) {
		wantResourceID := "example-resource-id"

		calledWaitForBlockCount := 0
		firstLatestBlockNumber := uint64(1337)
		secondLatestBlockNumber := uint64(1341)

		calledTransferCount := 0
		firstTransferBlockNumber := uint64(1340)
		secondTransferBlockNumber := uint64(1347)

		ackSpy1 := spy.NewAck()
		ackSpy2 := spy.NewAck()
		ackSpy3 := spy.NewAck()
		ackSpy4 := spy.NewAck()

		wantEvents := testdata.NewEnygmaTransferEventsFixture()
		dummyTransfersByChainID := testdata.NewTransfersByChainIDFixture()
		dummyBatchedTransfers := testdata.NewBatchedTransfersFixture()

		serializedEvent1, _ := json.Marshal(wantEvents[0])
		serializedEvent2, _ := json.Marshal(wantEvents[1])
		serializedEvent3, _ := json.Marshal(wantEvents[2])
		serializedEvent4, _ := json.Marshal(wantEvents[3])

		enygmaMQ := newMultiMessageMQ([]msgqueue.Message[service.EnygmaSerializedEvent]{
			{
				V: service.EnygmaSerializedEvent{
					Id:              wantEvents[0].MessageId,
					BlockNumber:     67,
					LogIndex:        0,
					TxHash:          common.Hash{},
					ResourceID:      wantResourceID,
					Type:            service.EnygmaTransferEvent,
					SerializedEvent: serializedEvent1,
				},
				Ack: ackSpy1.Fn(),
			},
			{
				V: service.EnygmaSerializedEvent{
					Id:              wantEvents[1].MessageId,
					BlockNumber:     67,
					LogIndex:        1,
					TxHash:          common.Hash{},
					ResourceID:      wantResourceID,
					Type:            service.EnygmaTransferEvent,
					SerializedEvent: serializedEvent2,
				},
				Ack: ackSpy2.Fn(),
			},
			{
				V: service.EnygmaSerializedEvent{
					Id:              wantEvents[2].MessageId,
					BlockNumber:     67,
					LogIndex:        2,
					TxHash:          common.Hash{},
					ResourceID:      wantResourceID,
					Type:            service.EnygmaTransferEvent,
					SerializedEvent: serializedEvent3,
				},
				Ack: ackSpy3.Fn(),
			},
			{
				V: service.EnygmaSerializedEvent{
					Id:              wantEvents[3].MessageId,
					BlockNumber:     67,
					LogIndex:        3,
					TxHash:          common.Hash{},
					ResourceID:      wantResourceID,
					Type:            service.EnygmaTransferEvent,
					SerializedEvent: serializedEvent4,
				},
				Ack: ackSpy4.Fn(),
			},
		})
		// New implementation calls WaitForBlock 3 times:
		// 1. Initial call in Run() with targetBlockNumber=0
		// 2. In handleEnygmaTransferBatch for batch 1 with targetBlockNumber=latestBlock (1337)
		// 3. In handleEnygmaTransferBatch for batch 2 with targetBlockNumber=firstTransferBlockNumber (1340)
		blockWaiter := &EnygmaBlockWaiterServiceMock{
			WaitForBlockFunc: func(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
				calledWaitForBlockCount += 1

				switch calledWaitForBlockCount {
				case 1:
					assert.EqualValues(t, 0, targetBlockNumber)
					return firstLatestBlockNumber, nil
				case 2:
					assert.EqualValues(t, firstLatestBlockNumber, targetBlockNumber)
					return firstLatestBlockNumber, nil
				case 3:
					assert.EqualValues(t, firstTransferBlockNumber, targetBlockNumber)
					return secondLatestBlockNumber, nil
				default:
					assert.Fail(t, "called block waiter more times than required")
					return 0, nil
				}
			},
		}
		batcher := newDefaultBatcherMock(t)
		batcher.GroupTransfersByChainIDFunc = func(txs []service.EnygmaTransferTx) map[string][]*types.EnygmaTransferBatchTx {
			assert.Equal(t, wantEvents, txs)
			return dummyTransfersByChainID
		}
		batcher.BatchTransfersFunc = func(_ context.Context, txsByChainID map[string][]*types.EnygmaTransferBatchTx) ([]map[string][]*types.EnygmaTransferBatchTx, error) {
			assert.Equal(t, dummyTransfersByChainID, txsByChainID)
			return dummyBatchedTransfers, nil
		}

		initiator := newDefaultInitiatorMock(t)
		initiator.HandleEnygmaCrossTransferFunc = func(ctx context.Context, _ string, resourceID string, blockNumber uint64, batch map[string][]*types.EnygmaTransferBatchTx) (uint64, error) {
			calledTransferCount += 1

			assert.Equal(t, wantResourceID, resourceID)
			switch calledTransferCount {
			case 1:
				assert.Equal(t, firstLatestBlockNumber, blockNumber)
				assert.Equal(t, dummyBatchedTransfers[0], batch)

				return firstTransferBlockNumber, nil
			case 2:
				assert.Equal(t, secondLatestBlockNumber, blockNumber)
				assert.Equal(t, dummyBatchedTransfers[1], batch)

				return secondTransferBlockNumber, nil
			default:
				assert.Fail(t, "called deposit too many times")
				return 0, nil
			}
		}
		reverter := newDefaultReverterMock(t)

		finalizer := newDefaultFinalizerMock(t)
		finalizer.ExecuteFinalizationFunc = func(ctx context.Context, _ string, blockNumber uint64, resourceID string) error {
			assert.Equal(t, secondTransferBlockNumber, blockNumber)
			assert.Equal(t, wantResourceID, resourceID)

			return nil
		}

		svc := service.NewEnygmaOrchestrator(
			time.Millisecond,
			enygmaMQ,
			blockWaiter,
			batcher,
			initiator,
			reverter,
			finalizer,
			2,
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, 2, len(initiator.HandleEnygmaCrossTransferCalls()))
		assert.Equal(t, 1, len(finalizer.ExecuteFinalizationCalls()))
		ackSpy1.AssertCalled(t)
		ackSpy2.AssertCalled(t)
		ackSpy3.AssertCalled(t)
		ackSpy4.AssertCalled(t)
	})

	t.Run("processes multiple resources concurrently", func(t *testing.T) {
		const (
			resourceA = "resource-A"
			resourceB = "resource-B"
		)

		ackSpyA := spy.NewAck()
		ackSpyB := spy.NewAck()

		event := service.EnygmaCreation{InitialSupply: big.NewInt(42)}
		serializedEvent, _ := json.Marshal(event)

		enygmaMQ := newMultiMessageMQ([]msgqueue.Message[service.EnygmaSerializedEvent]{
			{
				V: service.EnygmaSerializedEvent{
					Id:              "creation-A",
					BlockNumber:     10,
					LogIndex:        0,
					Type:            service.EnygmaCreationEvent,
					ResourceID:      resourceA,
					SerializedEvent: serializedEvent,
				},
				Ack: ackSpyA.Fn(),
			},
			{
				V: service.EnygmaSerializedEvent{
					Id:              "creation-B",
					BlockNumber:     11,
					LogIndex:        0,
					Type:            service.EnygmaCreationEvent,
					ResourceID:      resourceB,
					SerializedEvent: serializedEvent,
				},
				Ack: ackSpyB.Fn(),
			},
		})

		blockWaiter := &EnygmaBlockWaiterServiceMock{
			WaitForBlockFunc: func(ctx context.Context, target uint64) (uint64, error) {
				return 1337, nil
			},
		}

		var gateWG sync.WaitGroup
		gateWG.Add(2)
		releaseHandlers := make(chan struct{})

		batcher := newDefaultBatcherMock(t)
		initiator := newDefaultInitiatorMock(t)
		initiator.HandleEnygmaCreationFunc = func(ctx context.Context, chainEventID string, resourceID string, blockNumber uint64, initialSupply *big.Int) (uint64, error) {
			gateWG.Done()
			select {
			case <-releaseHandlers:
			case <-ctx.Done():
			}
			return blockNumber, nil
		}
		reverter := newDefaultReverterMock(t)

		var finalizedMu sync.Mutex
		finalizedResources := map[string]struct{}{}
		bothFinalized := make(chan struct{})
		finalizer := newDefaultFinalizerMock(t)
		finalizer.ExecuteFinalizationFunc = func(ctx context.Context, id string, blockNumber uint64, resourceID string) error {
			finalizedMu.Lock()
			finalizedResources[resourceID] = struct{}{}
			n := len(finalizedResources)
			finalizedMu.Unlock()
			if n == 2 {
				close(bothFinalized)
			}
			return nil
		}

		// Watcher: once both handlers are at the barrier, release them.
		// Records whether overlap was actually observed.
		overlapObserved := make(chan struct{})
		go func() {
			wgDone := make(chan struct{})
			go func() {
				gateWG.Wait()
				close(wgDone)
			}()
			select {
			case <-wgDone:
				close(overlapObserved)
				close(releaseHandlers)
			case <-time.After(2 * time.Second):
				close(releaseHandlers) // unblock whatever did arrive so the test terminates
			}
		}()

		svc := service.NewEnygmaOrchestrator(
			time.Millisecond,
			enygmaMQ,
			blockWaiter,
			batcher,
			initiator,
			reverter,
			finalizer,
			2,
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		runDone := make(chan error, 1)
		go func() {
			runDone <- svc.Run(ctx)
		}()

		select {
		case <-bothFinalized:
		case <-time.After(3 * time.Second):
			t.Fatal("did not observe both resources finalize within 3s")
		}

		cancel()
		require.NoError(t, <-runDone)

		select {
		case <-overlapObserved:
		default:
			t.Fatal("expected handlers to overlap within 2s — processing appears sequential")
		}

		assert.Equal(t, 2, len(initiator.HandleEnygmaCreationCalls()))
		finalizedMu.Lock()
		assert.Contains(t, finalizedResources, resourceA)
		assert.Contains(t, finalizedResources, resourceB)
		finalizedMu.Unlock()
		ackSpyA.AssertCalled(t)
		ackSpyB.AssertCalled(t)
	})

	t.Run("respects maxConcurrentResourceIDs limit", func(t *testing.T) {
		resources := []string{"resource-A", "resource-B", "resource-C"}

		event := service.EnygmaCreation{InitialSupply: big.NewInt(42)}
		serializedEvent, _ := json.Marshal(event)

		ackSpies := []*spy.Ack{spy.NewAck(), spy.NewAck(), spy.NewAck()}
		msgs := make([]msgqueue.Message[service.EnygmaSerializedEvent], len(resources))
		for i, r := range resources {
			msgs[i] = msgqueue.Message[service.EnygmaSerializedEvent]{
				V: service.EnygmaSerializedEvent{
					Id:              "creation-" + r,
					BlockNumber:     uint64(10 + i),
					LogIndex:        0,
					Type:            service.EnygmaCreationEvent,
					ResourceID:      r,
					SerializedEvent: serializedEvent,
				},
				Ack: ackSpies[i].Fn(),
			}
		}
		enygmaMQ := newMultiMessageMQ(msgs)

		blockWaiter := &EnygmaBlockWaiterServiceMock{
			WaitForBlockFunc: func(ctx context.Context, target uint64) (uint64, error) {
				return 1337, nil
			},
		}

		var inFlight, peak, arrived atomic.Int32
		firstTwoArrived := make(chan struct{})
		releaseHandlers := make(chan struct{})

		batcher := newDefaultBatcherMock(t)
		initiator := newDefaultInitiatorMock(t)
		initiator.HandleEnygmaCreationFunc = func(ctx context.Context, chainEventID string, resourceID string, blockNumber uint64, initialSupply *big.Int) (uint64, error) {
			n := inFlight.Add(1)
			for {
				old := peak.Load()
				if n <= old || peak.CompareAndSwap(old, n) {
					break
				}
			}
			if arrived.Add(1) == 2 {
				close(firstTwoArrived)
			}
			select {
			case <-releaseHandlers:
			case <-ctx.Done():
			}
			inFlight.Add(-1)
			return blockNumber, nil
		}
		reverter := newDefaultReverterMock(t)

		var finalizedMu sync.Mutex
		finalizedResources := map[string]struct{}{}
		allFinalized := make(chan struct{})
		finalizer := newDefaultFinalizerMock(t)
		finalizer.ExecuteFinalizationFunc = func(ctx context.Context, id string, blockNumber uint64, resourceID string) error {
			finalizedMu.Lock()
			finalizedResources[resourceID] = struct{}{}
			n := len(finalizedResources)
			finalizedMu.Unlock()
			if n == len(resources) {
				close(allFinalized)
			}
			return nil
		}

		// Snapshot the arrival count at the moment we release. Under correct
		// behavior errgroup.SetLimit(2) blocks the third eg.Go on a semaphore,
		// so the third handler cannot have entered before we release.
		var arrivedAtRelease int32
		go func() {
			select {
			case <-firstTwoArrived:
				arrivedAtRelease = arrived.Load()
				close(releaseHandlers)
			case <-time.After(2 * time.Second):
				close(releaseHandlers)
			}
		}()

		svc := service.NewEnygmaOrchestrator(
			time.Millisecond,
			enygmaMQ,
			blockWaiter,
			batcher,
			initiator,
			reverter,
			finalizer,
			2,
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		runDone := make(chan error, 1)
		go func() {
			runDone <- svc.Run(ctx)
		}()

		select {
		case <-allFinalized:
		case <-time.After(3 * time.Second):
			t.Fatal("did not observe all three resources finalize within 3s")
		}

		cancel()
		require.NoError(t, <-runDone)

		assert.EqualValues(t, 2, peak.Load(), "peak concurrency should equal the configured limit of 2")
		assert.EqualValues(t, 2, arrivedAtRelease, "third resource must be throttled while the first two are in-flight")
		assert.Equal(t, len(resources), len(initiator.HandleEnygmaCreationCalls()))
		finalizedMu.Lock()
		for _, r := range resources {
			assert.Contains(t, finalizedResources, r)
		}
		finalizedMu.Unlock()
		for _, a := range ackSpies {
			a.AssertCalled(t)
		}
	})
}
