package logparser_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/logparser"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/logparser/testdata"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/logrouter"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools/fake"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools/spy"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
)

func newNoopEnygmaEventMQ() *EnygmaEventMQMock {
	return &EnygmaEventMQMock{
		PushFunc:             func(ctx context.Context, _ service.EnygmaSerializedEvent) error { return nil },
		PushBatchChunkedFunc: func(ctx context.Context, _ []service.EnygmaSerializedEvent) error { return nil },
	}
}

func newNoopDvpBatchMQ() *DvpBatchMQMock {
	return &DvpBatchMQMock{
		PushFunc: func(ctx context.Context, _ service.DvpSerializedEventBatch) error { return nil },
	}
}

func TestEnygmaParser_EnygmaEvents(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("batches EnygmaCreation events", func(t *testing.T) {
		wantResourceIDHash := common.HexToHash("0xc001babe")
		wantResourceID := wantResourceIDHash.Hex()[2:]
		wantBlockNumber := uint64(100)
		wantEventType := service.EnygmaCreationEvent

		wantInitialSupply := big.NewInt(42)

		wantEnygmaCreateEvent := service.EnygmaCreation{
			InitialSupply: wantInitialSupply,
		}

		log := testdata.NewEnygmaCreationLogWith(
			testdata.WithCreateBlockNumber(wantBlockNumber),
			testdata.WithCreateResourceID(wantResourceIDHash),
			testdata.WithCreateInitialSupply(wantInitialSupply),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   []ethTypes.Log{log},
				},
				Ack: ackSpy.Fn(),
			}),
		}
		eventMQ := &EnygmaEventMQMock{
			PushBatchChunkedFunc: func(ctx context.Context, events []service.EnygmaSerializedEvent) error {
				assert.Equal(t, 1, len(events))
				ev := events[0]
				assert.Equal(t, wantResourceID, ev.ResourceID)
				assert.Equal(t, wantBlockNumber, ev.BlockNumber)
				assert.Equal(t, wantEventType, ev.Type)

				var event service.EnygmaCreation
				err := json.Unmarshal(ev.SerializedEvent, &event)
				assert.NoError(t, err)

				assert.Equal(t, wantEnygmaCreateEvent, event)
				return nil
			},
		}

		zkDvpBatchMQ := newNoopDvpBatchMQ()

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 1, len(eventMQ.PushBatchChunkedCalls()))
	})

	t.Run("batches two enygma creation events", func(t *testing.T) {
		wantResourceIDHash := common.HexToHash("0xc001babe")
		wantResourceID := wantResourceIDHash.Hex()[2:]
		wantBlockNumber := uint64(100)
		wantEventType := service.EnygmaCreationEvent

		wantInitialSupplyFirst := big.NewInt(42)
		wantInitialSupplySecond := big.NewInt(100)

		wantEnygmaCreateEvents := []service.EnygmaCreation{
			{
				InitialSupply: wantInitialSupplyFirst,
			},
			{
				InitialSupply: wantInitialSupplySecond,
			},
		}

		logFirst := testdata.NewEnygmaCreationLogWith(
			testdata.WithCreateBlockNumber(wantBlockNumber),
			testdata.WithCreateResourceID(wantResourceIDHash),
			testdata.WithCreateInitialSupply(wantInitialSupplyFirst),
		)
		logSecond := testdata.NewEnygmaCreationLogWith(
			testdata.WithCreateBlockNumber(wantBlockNumber),
			testdata.WithCreateResourceID(wantResourceIDHash),
			testdata.WithCreateInitialSupply(wantInitialSupplySecond),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   []ethTypes.Log{logFirst, logSecond},
				},
				Ack: ackSpy.Fn(),
			}),
		}
		eventMQ := &EnygmaEventMQMock{
			PushBatchChunkedFunc: func(ctx context.Context, events []service.EnygmaSerializedEvent) error {
				assert.Equal(t, 2, len(events))

				for i, ev := range events {
					assert.Equal(t, wantResourceID, ev.ResourceID)
					assert.Equal(t, wantBlockNumber, ev.BlockNumber)
					assert.Equal(t, wantEventType, ev.Type)

					var event service.EnygmaCreation
					err := json.Unmarshal(ev.SerializedEvent, &event)
					assert.NoError(t, err)
					assert.Equal(t, wantEnygmaCreateEvents[i], event)
				}
				return nil
			},
		}

		zkDvpBatchMQ := newNoopDvpBatchMQ()

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 1, len(eventMQ.PushBatchChunkedCalls()))
	})

	t.Run("batches two enygma creation events with different resoure IDs", func(t *testing.T) {
		wantBlockNumber := uint64(100)
		wantEventType := service.EnygmaCreationEvent

		wantResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantResourceIDFirst := wantResourceIDHashFirst.Hex()[2:]
		wantResourceIDHashSecond := common.HexToHash("0xc0cac01a")
		wantResourceIDSecond := wantResourceIDHashSecond.Hex()[2:]

		wantInitialSupplyFirst := big.NewInt(42)
		wantInitialSupplySecond := big.NewInt(100)

		wantEnygmaCreateEventFirst := service.EnygmaCreation{
			InitialSupply: wantInitialSupplyFirst,
		}
		wantEnygmaCreateEventSecond := service.EnygmaCreation{
			InitialSupply: wantInitialSupplySecond,
		}

		logFirst := testdata.NewEnygmaCreationLogWith(
			testdata.WithCreateBlockNumber(wantBlockNumber),
			testdata.WithCreateResourceID(wantResourceIDHashFirst),
			testdata.WithCreateInitialSupply(wantInitialSupplyFirst),
		)
		logSecond := testdata.NewEnygmaCreationLogWith(
			testdata.WithCreateBlockNumber(wantBlockNumber),
			testdata.WithCreateResourceID(wantResourceIDHashSecond),
			testdata.WithCreateInitialSupply(wantInitialSupplySecond),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   []ethTypes.Log{logFirst, logSecond},
				},
				Ack: ackSpy.Fn(),
			}),
		}

		// Use eventVerifier to handle events in any order
		verifier := newEventVerifier(t, wantEventType)
		verifier.expectEvent(wantResourceIDFirst, service.EnygmaCreationEvent, wantEnygmaCreateEventFirst)
		verifier.expectEvent(wantResourceIDSecond, service.EnygmaCreationEvent, wantEnygmaCreateEventSecond)

		eventMQ := &EnygmaEventMQMock{
			PushBatchChunkedFunc: func(ctx context.Context, events []service.EnygmaSerializedEvent) error {
				return verifier.verifyBatch(events)
			},
		}

		zkDvpBatchMQ := newNoopDvpBatchMQ()

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 1, len(eventMQ.PushBatchChunkedCalls()))
		verifier.assertAllSeen()
	})

	t.Run("batches create and deposit event", func(t *testing.T) {
		wantResourceIDHash := common.HexToHash("0xc001babe")
		wantResourceID := wantResourceIDHash.Hex()[2:]
		wantBlockNumber := uint64(100)

		wantInitialSupply := big.NewInt(42)
		wantEnygmaCreateEvent := service.EnygmaCreation{
			InitialSupply: wantInitialSupply,
		}

		wantAmount := big.NewInt(67)
		wantFrom := common.HexToAddress("0xc0cac01a")
		wantReferenceID := common.HexToHash("0xdeadc0de")
		wantTxHash := common.HexToHash("0xabcdef03")
		wantEnygmaDepositEvent := service.EnygmaDepositToDvp{
			Amount:        wantAmount,
			From:          wantFrom,
			ReferenceId:   wantReferenceID,
			TxHash:        wantTxHash,
			TxBlockNumber: big.NewInt(int64(wantBlockNumber)),
		}

		logCreate := testdata.NewEnygmaCreationLogWith(
			testdata.WithCreateBlockNumber(wantBlockNumber),
			testdata.WithCreateResourceID(wantResourceIDHash),
			testdata.WithCreateInitialSupply(wantInitialSupply),
		)
		logDeposit := testdata.NewEnygmaDepositToDvpLogWith(
			testdata.WithDepositBlockNumber(wantBlockNumber),
			testdata.WithDepositResourceID(wantResourceIDHash),
			testdata.WithDepositFrom(wantFrom),
			testdata.WithDepositAmount(wantAmount),
			testdata.WithDepositReferenceID(wantReferenceID),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   []ethTypes.Log{logDeposit, logCreate},
				},
				Ack: ackSpy.Fn(),
			}),
		}
		eventMQ := &EnygmaEventMQMock{
			PushBatchChunkedFunc: func(ctx context.Context, events []service.EnygmaSerializedEvent) error {
				// New implementation pushes all events in one batch
				assert.Equal(t, 2, len(events))

				// Find and verify each event by type
				var foundCreate, foundDeposit bool
				for _, ev := range events {
					assert.Equal(t, wantResourceID, ev.ResourceID)
					assert.Equal(t, wantBlockNumber, ev.BlockNumber)

					if ev.Type == service.EnygmaCreationEvent {
						foundCreate = true
						var event service.EnygmaCreation
						err := json.Unmarshal(ev.SerializedEvent, &event)
						assert.NoError(t, err)
						assert.Equal(t, wantEnygmaCreateEvent, event)
					} else if ev.Type == service.EnygmaDepositEvent {
						foundDeposit = true
						var event service.EnygmaDepositToDvp
						err := json.Unmarshal(ev.SerializedEvent, &event)
						assert.NoError(t, err)
						assert.Equal(t, wantEnygmaDepositEvent, event)
					}
				}
				assert.True(t, foundCreate, "expected creation event")
				assert.True(t, foundDeposit, "expected deposit event")
				return nil
			},
		}

		zkDvpBatchMQ := newNoopDvpBatchMQ()

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 1, len(eventMQ.PushBatchChunkedCalls()))
	})

	t.Run("batches two enygma creation events with different resoure IDs", func(t *testing.T) {
		wantBlockNumber := uint64(100)
		wantEventType := service.EnygmaDepositEvent

		wantResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantResourceIDFirst := wantResourceIDHashFirst.Hex()[2:]
		wantResourceIDHashSecond := common.HexToHash("0xc0cac01a")
		wantResourceIDSecond := wantResourceIDHashSecond.Hex()[2:]

		wantAmount := big.NewInt(67)
		wantFrom := common.HexToAddress("0xc0cac01a")
		wantReferenceID := common.HexToHash("0xdeadc0de")
		wantTxHash := common.HexToHash("0xabcdef03")

		wantEnygmaDepositEventFirst := service.EnygmaDepositToDvp{
			Amount:        wantAmount,
			From:          wantFrom,
			ReferenceId:   wantReferenceID,
			TxHash:        wantTxHash,
			TxBlockNumber: big.NewInt(int64(wantBlockNumber)),
		}
		wantEnygmaDepositEventSecond := service.EnygmaDepositToDvp{
			Amount:        wantAmount,
			From:          wantFrom,
			ReferenceId:   wantReferenceID,
			TxHash:        wantTxHash,
			TxBlockNumber: big.NewInt(int64(wantBlockNumber)),
		}

		logFirst := testdata.NewEnygmaDepositToDvpLogWith(
			testdata.WithDepositBlockNumber(wantBlockNumber),
			testdata.WithDepositResourceID(wantResourceIDHashFirst),
			testdata.WithDepositFrom(wantFrom),
			testdata.WithDepositAmount(wantAmount),
			testdata.WithDepositReferenceID(wantReferenceID),
		)
		logSecond := testdata.NewEnygmaDepositToDvpLogWith(
			testdata.WithDepositBlockNumber(wantBlockNumber),
			testdata.WithDepositResourceID(wantResourceIDHashSecond),
			testdata.WithDepositFrom(wantFrom),
			testdata.WithDepositAmount(wantAmount),
			testdata.WithDepositReferenceID(wantReferenceID),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   []ethTypes.Log{logFirst, logSecond},
				},
				Ack: ackSpy.Fn(),
			}),
		}

		// Use eventVerifier to handle events in any order
		verifier := newEventVerifier(t, wantEventType)
		verifier.expectEvent(wantResourceIDFirst, service.EnygmaDepositEvent, wantEnygmaDepositEventFirst)
		verifier.expectEvent(wantResourceIDSecond, service.EnygmaDepositEvent, wantEnygmaDepositEventSecond)

		eventMQ := &EnygmaEventMQMock{
			PushBatchChunkedFunc: func(ctx context.Context, events []service.EnygmaSerializedEvent) error {
				return verifier.verifyBatch(events)
			},
		}

		zkDvpBatchMQ := newNoopDvpBatchMQ()

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 1, len(eventMQ.PushBatchChunkedCalls()))
		verifier.assertAllSeen()
	})

	t.Run("batches create and deposit events from different resource IDs", func(t *testing.T) {
		wantBlockNumber := uint64(42)

		wantResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantResourceIDFirst := wantResourceIDHashFirst.Hex()[2:]
		wantResourceIDHashSecond := common.HexToHash("0xc0cac01a")
		wantResourceIDSecond := wantResourceIDHashSecond.Hex()[2:]

		wantInitialSupply := big.NewInt(42)

		wantAmount := big.NewInt(67)
		wantFrom := common.HexToAddress("0xc0cac01a")
		wantReferenceID := common.HexToHash("0xdeadc0de")
		wantTxHash := common.HexToHash("0xabcdef03")

		wantEnygmaCreateEventFirst := service.EnygmaCreation{
			InitialSupply: wantInitialSupply,
		}
		wantEnygmaCreateEventSecond := service.EnygmaCreation{
			InitialSupply: wantInitialSupply,
		}

		wantEnygmaDepositEventFirst := service.EnygmaDepositToDvp{
			Amount:        wantAmount,
			From:          wantFrom,
			ReferenceId:   wantReferenceID,
			TxHash:        wantTxHash,
			TxBlockNumber: big.NewInt(int64(wantBlockNumber)),
		}
		wantEnygmaDepositEventSecond := service.EnygmaDepositToDvp{
			Amount:        wantAmount,
			From:          wantFrom,
			ReferenceId:   wantReferenceID,
			TxHash:        wantTxHash,
			TxBlockNumber: big.NewInt(int64(wantBlockNumber)),
		}

		// Logs must be grouped by resource ID with events in type order
		logs := []ethTypes.Log{
			buildEnygmaDepositToDvpLog(wantEnygmaDepositEventSecond, wantBlockNumber, wantResourceIDSecond),
			buildEnygmaCreationLog(wantEnygmaCreateEventFirst, wantBlockNumber, wantResourceIDFirst),
			buildEnygmaCreationLog(wantEnygmaCreateEventSecond, wantBlockNumber, wantResourceIDSecond),
			buildEnygmaDepositToDvpLog(wantEnygmaDepositEventFirst, wantBlockNumber, wantResourceIDFirst),
		}

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   logs,
				},
				Ack: ackSpy.Fn(),
			}),
		}

		// Use eventVerifier to handle events in any order
		verifier := newEventVerifier(t, service.EnygmaCreationEvent)
		verifier.expectEvent(wantResourceIDFirst, service.EnygmaCreationEvent, wantEnygmaCreateEventFirst)
		verifier.expectEvent(wantResourceIDSecond, service.EnygmaCreationEvent, wantEnygmaCreateEventSecond)
		verifier.expectEvent(wantResourceIDFirst, service.EnygmaDepositEvent, wantEnygmaDepositEventFirst)
		verifier.expectEvent(wantResourceIDSecond, service.EnygmaDepositEvent, wantEnygmaDepositEventSecond)

		eventMQ := &EnygmaEventMQMock{
			PushBatchChunkedFunc: func(ctx context.Context, events []service.EnygmaSerializedEvent) error {
				return verifier.verifyBatch(events)
			},
		}

		zkDvpBatchMQ := newNoopDvpBatchMQ()

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 1, len(eventMQ.PushBatchChunkedCalls())) // 1 batch with all 4 events
		verifier.assertAllSeen()
	})

	t.Run("batches two enygma transfer events with different resource IDs", func(t *testing.T) {
		wantBlockNumber := uint64(100)
		wantEventType := service.EnygmaTransferEvent

		wantResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantResourceIDFirst := wantResourceIDHashFirst.Hex()[2:]
		wantResourceIDHashSecond := common.HexToHash("0xc0cac01a")
		wantResourceIDSecond := wantResourceIDHashSecond.Hex()[2:]

		wantValue := []*big.Int{big.NewInt(100), big.NewInt(200)}

		wantToChainId := []*big.Int{big.NewInt(1), big.NewInt(2)}

		wantTo := []common.Address{
			common.HexToAddress("0x1111111111111111111111111111111111111111"),
			common.HexToAddress("0x2222222222222222222222222222222222222222"),
		}
		wantFrom := common.HexToAddress("0x3333333333333333333333333333333333333333")
		wantReferenceID := common.HexToHash("0xdeadc0de")

		// Each contract event with 2 destinations gets split into 2 EnygmaTransferTx
		// Expected transfers indexed by (resourceID, destination index). Each recipient carries the
		// default single-element [mintStep] from the testdata fixture, with distinct Args per
		// recipient so the parallel split's per-recipient ordering is asserted.
		wantProgramData := [][]types.EnygmaProgramData{
			{{ResourceId: [32]byte{0x0A}, ContractAddress: common.Address{}, Selector: [4]byte{0x11}, Args: []byte{0xDE, 0xAD}}},
			{{ResourceId: [32]byte{0x0B}, ContractAddress: common.Address{}, Selector: [4]byte{0x22}, Args: []byte{0xBE, 0xEF}}},
		}
		wantTransfer := func(toIndex int) service.EnygmaTransferTx {
			return service.EnygmaTransferTx{
				ReferenceId: wantReferenceID,
				FromAddress: wantFrom,
				ToChainId:   wantToChainId[toIndex],
				ToAddress:   wantTo[toIndex],
				ToAmount:    wantValue[toIndex],
				ProgramData: wantProgramData[toIndex],
			}
		}

		logFirst := testdata.NewEnygmaSendTransferPNHLogWith(
			testdata.WithTransferBlockNumber(wantBlockNumber),
			testdata.WithTransferResourceID(wantResourceIDHashFirst),
			testdata.WithTransferValue(wantValue),
			testdata.WithTransferToChainId(wantToChainId),
			testdata.WithTransferTo(wantTo),
			testdata.WithTransferFrom(wantFrom),
			testdata.WithTransferReferenceID(wantReferenceID),
		)
		logSecond := testdata.NewEnygmaSendTransferPNHLogWith(
			testdata.WithTransferBlockNumber(wantBlockNumber),
			testdata.WithTransferResourceID(wantResourceIDHashSecond),
			testdata.WithTransferValue(wantValue),
			testdata.WithTransferToChainId(wantToChainId),
			testdata.WithTransferTo(wantTo),
			testdata.WithTransferFrom(wantFrom),
			testdata.WithTransferReferenceID(wantReferenceID),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   []ethTypes.Log{logFirst, logSecond},
				},
				Ack: ackSpy.Fn(),
			}),
		}

		eventMQ := &EnygmaEventMQMock{
			PushBatchChunkedFunc: func(ctx context.Context, events []service.EnygmaSerializedEvent) error {
				// 2 logs × 2 destinations = 4 events
				assert.Equal(t, 4, len(events))

				// Track events by resource ID
				eventsByResource := make(map[string][]service.EnygmaTransferTx)
				for _, ev := range events {
					assert.Equal(t, wantEventType, ev.Type)
					assert.Equal(t, wantBlockNumber, ev.BlockNumber)

					var transfer service.EnygmaTransferTx
					err := json.Unmarshal(ev.SerializedEvent, &transfer)
					assert.NoError(t, err)
					eventsByResource[ev.ResourceID] = append(eventsByResource[ev.ResourceID], transfer)
				}

				// Verify each resource ID has 2 transfer events
				assert.Len(t, eventsByResource[wantResourceIDFirst], 2)
				assert.Len(t, eventsByResource[wantResourceIDSecond], 2)

				// Verify transfer events for each resource (ignoring MessageId which is dynamically generated)
				for _, transfers := range eventsByResource {
					for i, transfer := range transfers {
						expected := wantTransfer(i)
						assert.Equal(t, i, transfer.DestIdx)
						assert.Equal(t, expected.ReferenceId, transfer.ReferenceId)
						assert.Equal(t, expected.FromAddress, transfer.FromAddress)
						assert.Equal(t, expected.ToChainId, transfer.ToChainId)
						assert.Equal(t, expected.ToAddress, transfer.ToAddress)
						assert.Equal(t, expected.ToAmount, transfer.ToAmount)
						assert.Equal(t, expected.ProgramData, transfer.ProgramData)
					}
				}
				return nil
			},
		}

		zkDvpBatchMQ := newNoopDvpBatchMQ()

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 1, len(eventMQ.PushBatchChunkedCalls()))
	})

	t.Run("batches two enygma withdraw events with different resource IDs", func(t *testing.T) {
		wantBlockNumber := uint64(100)
		wantEventType := service.EnygmaWithdrawEvent

		wantResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantResourceIDFirst := wantResourceIDHashFirst.Hex()[2:]
		wantResourceIDHashSecond := common.HexToHash("0xc0cac01a")
		wantResourceIDSecond := wantResourceIDHashSecond.Hex()[2:]

		wantAmount := big.NewInt(150)
		wantTo := common.HexToAddress("0x4444444444444444444444444444444444444444")
		wantReferenceID := common.HexToHash("0xdeadc0de")
		wantTxHash := common.HexToHash("0xabcdef04")

		wantEnygmaWithdrawEventFirst := service.EnygmaWithdrawFromDvp{
			Amount:        wantAmount,
			To:            wantTo,
			ReferenceId:   wantReferenceID,
			TxHash:        wantTxHash,
			TxBlockNumber: big.NewInt(int64(wantBlockNumber)),
		}
		wantEnygmaWithdrawEventSecond := service.EnygmaWithdrawFromDvp{
			Amount:        wantAmount,
			To:            wantTo,
			ReferenceId:   wantReferenceID,
			TxHash:        wantTxHash,
			TxBlockNumber: big.NewInt(int64(wantBlockNumber)),
		}

		logFirst := testdata.NewEnygmaWithdrawFromDvpLogWith(
			testdata.WithWithdrawBlockNumber(wantBlockNumber),
			testdata.WithWithdrawResourceID(wantResourceIDHashFirst),
			testdata.WithWithdrawAmount(wantAmount),
			testdata.WithWithdrawTo(wantTo),
			testdata.WithWithdrawReferenceID(wantReferenceID),
		)
		logSecond := testdata.NewEnygmaWithdrawFromDvpLogWith(
			testdata.WithWithdrawBlockNumber(wantBlockNumber),
			testdata.WithWithdrawResourceID(wantResourceIDHashSecond),
			testdata.WithWithdrawAmount(wantAmount),
			testdata.WithWithdrawTo(wantTo),
			testdata.WithWithdrawReferenceID(wantReferenceID),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   []ethTypes.Log{logFirst, logSecond},
				},
				Ack: ackSpy.Fn(),
			}),
		}

		verifier := newEventVerifier(t, wantEventType)
		verifier.expectEvent(wantResourceIDFirst, service.EnygmaWithdrawEvent, wantEnygmaWithdrawEventFirst)
		verifier.expectEvent(wantResourceIDSecond, service.EnygmaWithdrawEvent, wantEnygmaWithdrawEventSecond)

		eventMQ := &EnygmaEventMQMock{
			PushBatchChunkedFunc: func(ctx context.Context, events []service.EnygmaSerializedEvent) error {
				return verifier.verifyBatch(events)
			},
		}

		zkDvpBatchMQ := newNoopDvpBatchMQ()

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 1, len(eventMQ.PushBatchChunkedCalls()))
		verifier.assertAllSeen()
	})

	t.Run("batches mint and burn events as supply update with different resource IDs", func(t *testing.T) {
		wantBlockNumber := uint64(100)
		wantEventType := service.EnygmaSupplyUpdateEvent

		wantResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantResourceIDFirst := wantResourceIDHashFirst.Hex()[2:]

		wantResourceIDHashSecond := common.HexToHash("0xc0cac01a")
		wantResourceIDSecond := wantResourceIDHashSecond.Hex()[2:]

		wantTxHashFirst := common.HexToHash("0xdeadc0de")
		wantToFirst := common.HexToAddress("0x1111111111111111111111111111111111111111")
		wantMintAmount := big.NewInt(500)

		wantTxHashSecond := common.HexToHash("0xc0fefeed")
		wantToSecond := common.HexToAddress("0x2222222222222222222222222222222222222222")
		wantBurnAmount := big.NewInt(200)

		wantSupplyUpdateEventFirst := service.EnygmaSupplyUpdate{
			TxHash: wantTxHashFirst,
			Amount: wantMintAmount,
			To:     wantToFirst,
		}
		wantSupplyUpdateEventSecond := service.EnygmaSupplyUpdate{
			TxHash: wantTxHashSecond,
			To:     wantToSecond,
			Amount: new(big.Int).Neg(wantBurnAmount),
		}

		logMint := testdata.NewEnygmaMintLogWith(
			testdata.WithMintBlockNumber(wantBlockNumber),
			testdata.WithMintResourceID(wantResourceIDHashFirst),
			testdata.WithMintTxHash(wantTxHashFirst),
			testdata.WithMintToAddress(wantToFirst),
			testdata.WithMintAmount(wantMintAmount),
		)

		logBurn := testdata.NewEnygmaBurnLogWith(
			testdata.WithBurnBlockNumber(wantBlockNumber),
			testdata.WithBurnResourceID(wantResourceIDHashSecond),
			testdata.WithBurnTxHash(wantTxHashSecond),
			testdata.WithBurnFromAddress(wantToSecond),
			testdata.WithBurnAmount(wantBurnAmount),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   []ethTypes.Log{logMint, logBurn},
				},
				Ack: ackSpy.Fn(),
			}),
		}

		verifier := newEventVerifier(t, wantEventType)
		verifier.expectEvent(wantResourceIDFirst, service.EnygmaSupplyUpdateEvent, wantSupplyUpdateEventFirst)
		verifier.expectEvent(wantResourceIDSecond, service.EnygmaSupplyUpdateEvent, wantSupplyUpdateEventSecond)

		eventMQ := &EnygmaEventMQMock{
			PushBatchChunkedFunc: func(ctx context.Context, events []service.EnygmaSerializedEvent) error {
				return verifier.verifyBatch(events)
			},
		}

		zkDvpBatchMQ := newNoopDvpBatchMQ()

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 1, len(eventMQ.PushBatchChunkedCalls()))
		verifier.assertAllSeen()
	})
}

func TestEnygmaParser_DvpEvents(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("batches Dvp721Creation events", func(t *testing.T) {
		wantResourceIDHash := common.HexToHash("0xc001babe")
		wantResourceID := wantResourceIDHash.Hex()[2:]
		wantBlockNumber := uint64(100)
		wantEventType := service.Dvp721CreationEvent

		wantDvp721CreationEvents := []service.Dvp721Creation{
			{
				ChainEventID: "100-6-5",
				ResourceId:   wantResourceID,
			},
		}

		log := testdata.NewDvp721CreationLogWith(
			testdata.WithDvp721CreationBlockNumber(wantBlockNumber),
			testdata.WithDvp721CreationResourceID(wantResourceIDHash),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   []ethTypes.Log{log},
				},
				Ack: ackSpy.Fn(),
			}),
		}
		eventMQ := newNoopEnygmaEventMQ()
		zkDvpBatchMQ := &DvpBatchMQMock{
			PushFunc: func(ctx context.Context, batch service.DvpSerializedEventBatch) error {
				assert.Equal(t, wantBlockNumber, batch.BlockNumber)
				assert.Equal(t, wantEventType, batch.Type)

				events, err := deserializeEvents[service.Dvp721Creation](batch.SerializedEvents)
				assert.NoError(t, err)

				assert.Equal(t, wantDvp721CreationEvents, events)
				return nil
			},
		}

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 0, len(eventMQ.PushBatchChunkedCalls()))
		assert.Equal(t, 1, len(zkDvpBatchMQ.PushCalls()))
	})

	t.Run("batches two Dvp721Creation events", func(t *testing.T) {
		wantResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantResourceIDFirst := wantResourceIDHashFirst.Hex()[2:]
		wantResourceIDHashSecond := common.HexToHash("0xdeadbeef")
		wantResourceIDSecond := wantResourceIDHashSecond.Hex()[2:]
		wantBlockNumber := uint64(100)
		wantEventType := service.Dvp721CreationEvent

		wantDvp721CreationEvents := []service.Dvp721Creation{
			{
				ChainEventID: "100-6-5",
				ResourceId:   wantResourceIDFirst,
			},
			{
				ChainEventID: "100-6-5",
				ResourceId:   wantResourceIDSecond,
			},
		}

		logFirst := testdata.NewDvp721CreationLogWith(
			testdata.WithDvp721CreationBlockNumber(wantBlockNumber),
			testdata.WithDvp721CreationResourceID(wantResourceIDHashFirst),
		)
		logSecond := testdata.NewDvp721CreationLogWith(
			testdata.WithDvp721CreationBlockNumber(wantBlockNumber),
			testdata.WithDvp721CreationResourceID(wantResourceIDHashSecond),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   []ethTypes.Log{logFirst, logSecond},
				},
				Ack: ackSpy.Fn(),
			}),
		}
		eventMQ := newNoopEnygmaEventMQ()
		zkDvpBatchMQ := &DvpBatchMQMock{
			PushFunc: func(ctx context.Context, batch service.DvpSerializedEventBatch) error {
				assert.Equal(t, wantBlockNumber, batch.BlockNumber)
				assert.Equal(t, wantEventType, batch.Type)

				events, err := deserializeEvents[service.Dvp721Creation](batch.SerializedEvents)
				assert.NoError(t, err)

				assert.Equal(t, wantDvp721CreationEvents, events)
				return nil
			},
		}

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 0, len(eventMQ.PushBatchChunkedCalls()))
		assert.Equal(t, 1, len(zkDvpBatchMQ.PushCalls()))
	})

	t.Run("batches two Dvp721Mint events", func(t *testing.T) {
		wantResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantResourceIDFirst := wantResourceIDHashFirst.Hex()[2:]
		wantResourceIDHashSecond := common.HexToHash("0xdeadbeef")
		wantResourceIDSecond := wantResourceIDHashSecond.Hex()[2:]
		wantBlockNumber := uint64(100)
		wantEventType := service.Dvp721MintEvent

		wantNftIdFirst := big.NewInt(42)
		wantNftIdSecond := big.NewInt(99)

		wantDvp721MintEvents := []service.Dvp721Mint{
			{
				ChainEventID: "100-7-6",
				ResourceId:   wantResourceIDFirst,
				NftId:        wantNftIdFirst,
			},
			{
				ChainEventID: "100-7-6",
				ResourceId:   wantResourceIDSecond,
				NftId:        wantNftIdSecond,
			},
		}

		logFirst := testdata.NewDvp721MintLogWith(
			testdata.WithDvp721MintBlockNumber(wantBlockNumber),
			testdata.WithDvp721MintResourceID(wantResourceIDHashFirst),
			testdata.WithDvp721MintNftId(wantNftIdFirst),
		)
		logSecond := testdata.NewDvp721MintLogWith(
			testdata.WithDvp721MintBlockNumber(wantBlockNumber),
			testdata.WithDvp721MintResourceID(wantResourceIDHashSecond),
			testdata.WithDvp721MintNftId(wantNftIdSecond),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   []ethTypes.Log{logFirst, logSecond},
				},
				Ack: ackSpy.Fn(),
			}),
		}
		eventMQ := newNoopEnygmaEventMQ()
		zkDvpBatchMQ := &DvpBatchMQMock{
			PushFunc: func(ctx context.Context, batch service.DvpSerializedEventBatch) error {
				assert.Equal(t, wantBlockNumber, batch.BlockNumber)
				assert.Equal(t, wantEventType, batch.Type)

				events, err := deserializeEvents[service.Dvp721Mint](batch.SerializedEvents)
				assert.NoError(t, err)

				assert.Equal(t, wantDvp721MintEvents, events)
				return nil
			},
		}

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 0, len(eventMQ.PushBatchChunkedCalls()))
		assert.Equal(t, 1, len(zkDvpBatchMQ.PushCalls()))
	})

	t.Run("batches two Dvp721Burn events", func(t *testing.T) {
		wantResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantResourceIDFirst := wantResourceIDHashFirst.Hex()[2:]
		wantResourceIDHashSecond := common.HexToHash("0xdeadbeef")
		wantResourceIDSecond := wantResourceIDHashSecond.Hex()[2:]
		wantBlockNumber := uint64(100)
		wantEventType := service.Dvp721BurnEvent

		wantNftIdFirst := big.NewInt(42)
		wantNftIdSecond := big.NewInt(99)

		wantDvp721BurnEvents := []service.Dvp721Burn{
			{
				ChainEventID: "100-8-7",
				ResourceId:   wantResourceIDFirst,
				NftId:        wantNftIdFirst,
			},
			{
				ChainEventID: "100-8-7",
				ResourceId:   wantResourceIDSecond,
				NftId:        wantNftIdSecond,
			},
		}

		logFirst := testdata.NewDvp721BurnLogWith(
			testdata.WithDvp721BurnBlockNumber(wantBlockNumber),
			testdata.WithDvp721BurnResourceID(wantResourceIDHashFirst),
			testdata.WithDvp721BurnNftId(wantNftIdFirst),
		)
		logSecond := testdata.NewDvp721BurnLogWith(
			testdata.WithDvp721BurnBlockNumber(wantBlockNumber),
			testdata.WithDvp721BurnResourceID(wantResourceIDHashSecond),
			testdata.WithDvp721BurnNftId(wantNftIdSecond),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   []ethTypes.Log{logFirst, logSecond},
				},
				Ack: ackSpy.Fn(),
			}),
		}
		eventMQ := newNoopEnygmaEventMQ()
		zkDvpBatchMQ := &DvpBatchMQMock{
			PushFunc: func(ctx context.Context, batch service.DvpSerializedEventBatch) error {
				assert.Equal(t, wantBlockNumber, batch.BlockNumber)
				assert.Equal(t, wantEventType, batch.Type)

				events, err := deserializeEvents[service.Dvp721Burn](batch.SerializedEvents)
				assert.NoError(t, err)

				assert.Equal(t, wantDvp721BurnEvents, events)
				return nil
			},
		}

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 0, len(eventMQ.PushBatchChunkedCalls()))
		assert.Equal(t, 1, len(zkDvpBatchMQ.PushCalls()))
	})

	t.Run("batches two Dvp721DepositIntoDvp events", func(t *testing.T) {
		wantResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantResourceIDFirst := wantResourceIDHashFirst.Hex()[2:]
		wantResourceIDHashSecond := common.HexToHash("0xdeadbeef")
		wantResourceIDSecond := wantResourceIDHashSecond.Hex()[2:]
		wantBlockNumber := uint64(100)
		wantEventType := service.Dvp721DepositIntoDvpEvent

		wantNftIdFirst := big.NewInt(42)
		wantNftIdSecond := big.NewInt(99)
		wantFromFirst := common.HexToAddress("0x1111111111111111111111111111111111111111")
		wantFromSecond := common.HexToAddress("0x2222222222222222222222222222222222222222")

		wantDvp721DepositIntoDvpEvents := []service.Dvp721DepositIntoDvp{
			{
				ChainEventID:  "100-9-8",
				ResourceId:    wantResourceIDFirst,
				NftId:         wantNftIdFirst,
				From:          wantFromFirst,
				TxHash:        common.HexToHash("0xabcdef08").Hex(),
				TxBlockNumber: big.NewInt(int64(wantBlockNumber)),
			},
			{
				ChainEventID:  "100-9-8",
				ResourceId:    wantResourceIDSecond,
				NftId:         wantNftIdSecond,
				From:          wantFromSecond,
				TxHash:        common.HexToHash("0xabcdef08").Hex(),
				TxBlockNumber: big.NewInt(int64(wantBlockNumber)),
			},
		}

		logFirst := testdata.NewDvp721DepositIntoDvpLogWith(
			testdata.WithDvp721DepositIntoDvpBlockNumber(wantBlockNumber),
			testdata.WithDvp721DepositIntoDvpResourceID(wantResourceIDHashFirst),
			testdata.WithDvp721DepositIntoDvpNftId(wantNftIdFirst),
			testdata.WithDvp721DepositIntoDvpFrom(wantFromFirst),
		)
		logSecond := testdata.NewDvp721DepositIntoDvpLogWith(
			testdata.WithDvp721DepositIntoDvpBlockNumber(wantBlockNumber),
			testdata.WithDvp721DepositIntoDvpResourceID(wantResourceIDHashSecond),
			testdata.WithDvp721DepositIntoDvpNftId(wantNftIdSecond),
			testdata.WithDvp721DepositIntoDvpFrom(wantFromSecond),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   []ethTypes.Log{logFirst, logSecond},
				},
				Ack: ackSpy.Fn(),
			}),
		}
		eventMQ := newNoopEnygmaEventMQ()
		zkDvpBatchMQ := &DvpBatchMQMock{
			PushFunc: func(ctx context.Context, batch service.DvpSerializedEventBatch) error {
				assert.Equal(t, wantBlockNumber, batch.BlockNumber)
				assert.Equal(t, wantEventType, batch.Type)

				events, err := deserializeEvents[service.Dvp721DepositIntoDvp](batch.SerializedEvents)
				assert.NoError(t, err)

				assert.Equal(t, wantDvp721DepositIntoDvpEvents, events)
				return nil
			},
		}

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 0, len(eventMQ.PushBatchChunkedCalls()))
		assert.Equal(t, 1, len(zkDvpBatchMQ.PushCalls()))
	})

	t.Run("batches two Dvp721WithdrawFromDvp events", func(t *testing.T) {
		wantResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantResourceIDFirst := wantResourceIDHashFirst.Hex()[2:]
		wantResourceIDHashSecond := common.HexToHash("0xdeadbeef")
		wantResourceIDSecond := wantResourceIDHashSecond.Hex()[2:]
		wantBlockNumber := uint64(100)
		wantEventType := service.Dvp721WithdrawFromDvpEvent

		wantNftIdFirst := big.NewInt(42)
		wantNftIdSecond := big.NewInt(99)
		wantOwnerFirst := common.HexToAddress("0x1111111111111111111111111111111111111111")
		wantOwnerSecond := common.HexToAddress("0x2222222222222222222222222222222222222222")

		wantDvp721WithdrawFromDvpEvents := []service.Dvp721WithdrawFromDvp{
			{
				ChainEventID:  "100-10-9",
				ResourceId:    wantResourceIDFirst,
				NftId:         wantNftIdFirst,
				Owner:         wantOwnerFirst,
				TxHash:        common.HexToHash("0xabcdef09").Hex(),
				TxBlockNumber: big.NewInt(int64(wantBlockNumber)),
			},
			{
				ChainEventID:  "100-10-9",
				ResourceId:    wantResourceIDSecond,
				NftId:         wantNftIdSecond,
				Owner:         wantOwnerSecond,
				TxHash:        common.HexToHash("0xabcdef09").Hex(),
				TxBlockNumber: big.NewInt(int64(wantBlockNumber)),
			},
		}

		logFirst := testdata.NewDvp721WithdrawFromDvpLogWith(
			testdata.WithDvp721WithdrawFromDvpBlockNumber(wantBlockNumber),
			testdata.WithDvp721WithdrawFromDvpResourceID(wantResourceIDHashFirst),
			testdata.WithDvp721WithdrawFromDvpNftId(wantNftIdFirst),
			testdata.WithDvp721WithdrawFromDvpOwner(wantOwnerFirst),
		)
		logSecond := testdata.NewDvp721WithdrawFromDvpLogWith(
			testdata.WithDvp721WithdrawFromDvpBlockNumber(wantBlockNumber),
			testdata.WithDvp721WithdrawFromDvpResourceID(wantResourceIDHashSecond),
			testdata.WithDvp721WithdrawFromDvpNftId(wantNftIdSecond),
			testdata.WithDvp721WithdrawFromDvpOwner(wantOwnerSecond),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   []ethTypes.Log{logFirst, logSecond},
				},
				Ack: ackSpy.Fn(),
			}),
		}
		eventMQ := newNoopEnygmaEventMQ()
		zkDvpBatchMQ := &DvpBatchMQMock{
			PushFunc: func(ctx context.Context, batch service.DvpSerializedEventBatch) error {
				assert.Equal(t, wantBlockNumber, batch.BlockNumber)
				assert.Equal(t, wantEventType, batch.Type)

				events, err := deserializeEvents[service.Dvp721WithdrawFromDvp](batch.SerializedEvents)
				assert.NoError(t, err)

				assert.Equal(t, wantDvp721WithdrawFromDvpEvents, events)
				return nil
			},
		}

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 0, len(eventMQ.PushBatchChunkedCalls()))
		assert.Equal(t, 1, len(zkDvpBatchMQ.PushCalls()))
	})

	t.Run("batches create and mint event", func(t *testing.T) {
		wantResourceIDHash := common.HexToHash("0xc001babe")
		wantResourceID := wantResourceIDHash.Hex()[2:]
		wantBlockNumber := uint64(100)

		wantDvp721CreationEvents := []service.Dvp721Creation{
			{
				ChainEventID: "100-6-5",
				ResourceId:   wantResourceID,
			},
		}

		wantNftId := big.NewInt(42)
		wantDvp721MintEvents := []service.Dvp721Mint{
			{
				ChainEventID: "100-7-6",
				ResourceId:   wantResourceID,
				NftId:        wantNftId,
			},
		}

		logCreate := testdata.NewDvp721CreationLogWith(
			testdata.WithDvp721CreationBlockNumber(wantBlockNumber),
			testdata.WithDvp721CreationResourceID(wantResourceIDHash),
		)
		logMint := testdata.NewDvp721MintLogWith(
			testdata.WithDvp721MintBlockNumber(wantBlockNumber),
			testdata.WithDvp721MintResourceID(wantResourceIDHash),
			testdata.WithDvp721MintNftId(wantNftId),
		)

		ackSpy := spy.NewAck()
		pushCounter := 0

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   []ethTypes.Log{logMint, logCreate},
				},
				Ack: ackSpy.Fn(),
			}),
		}
		eventMQ := newNoopEnygmaEventMQ()
		zkDvpBatchMQ := &DvpBatchMQMock{
			PushFunc: func(ctx context.Context, batch service.DvpSerializedEventBatch) error {
				assert.Equal(t, wantBlockNumber, batch.BlockNumber)
				switch pushCounter {
				case 0:
					assert.Equal(t, service.Dvp721CreationEvent, batch.Type)

					events, err := deserializeEvents[service.Dvp721Creation](batch.SerializedEvents)
					assert.NoError(t, err)

					assert.Equal(t, wantDvp721CreationEvents, events)
				case 1:
					assert.Equal(t, service.Dvp721MintEvent, batch.Type)

					events, err := deserializeEvents[service.Dvp721Mint](batch.SerializedEvents)
					assert.NoError(t, err)

					assert.Equal(t, wantDvp721MintEvents, events)
				default:
					assert.Fail(t, "called push too many times")
				}

				pushCounter += 1

				return nil
			},
		}

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 0, len(eventMQ.PushBatchChunkedCalls()))
		assert.Equal(t, 2, len(zkDvpBatchMQ.PushCalls()))
	})

	t.Run("batches two Dvp721SwapForEnygma events", func(t *testing.T) {
		wantNftResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantNftResourceIDFirst := wantNftResourceIDHashFirst.Hex()[2:]
		wantNftResourceIDHashSecond := common.HexToHash("0xdeadbeef")
		wantNftResourceIDSecond := wantNftResourceIDHashSecond.Hex()[2:]
		wantBlockNumber := uint64(100)

		wantNftIdFirst := big.NewInt(42)
		wantNftIdSecond := big.NewInt(99)

		wantEvents := []service.Dvp721SwapForEnygma{
			{
				SharedId:         "0102030000000000000000000000000000000000000000000000000000000000",
				DestChainId:      big.NewInt(1),
				From:             common.HexToAddress("0x1111111111111111111111111111111111111111"),
				NftResourceId:    wantNftResourceIDFirst,
				NftId:            wantNftIdFirst.String(),
				EnygmaResourceId: "0d0e0f0000000000000000000000000000000000000000000000000000000000",
				EnygmaAmount:     big.NewInt(1000),
				TxHash:           common.HexToHash("0xabcdef10").Hex(),
				//nolint:gosec // test data with known values
				TxBlockNumber: big.NewInt(int64(wantBlockNumber)),
			},
			{
				SharedId:         "0102030000000000000000000000000000000000000000000000000000000000",
				DestChainId:      big.NewInt(1),
				From:             common.HexToAddress("0x1111111111111111111111111111111111111111"),
				NftResourceId:    wantNftResourceIDSecond,
				NftId:            wantNftIdSecond.String(),
				EnygmaResourceId: "0d0e0f0000000000000000000000000000000000000000000000000000000000",
				EnygmaAmount:     big.NewInt(1000),
				TxHash:           common.HexToHash("0xabcdef10").Hex(),
				//nolint:gosec // test data with known values
				TxBlockNumber: big.NewInt(int64(wantBlockNumber)),
			},
		}

		logFirst := testdata.NewDvp721SwapForEnygmaLogWith(
			testdata.WithDvp721SwapForEnygmaBlockNumber(wantBlockNumber),
			testdata.WithDvp721SwapForEnygmaNftResourceID(wantNftResourceIDHashFirst),
			testdata.WithDvp721SwapForEnygmaNftId(wantNftIdFirst),
		)
		logSecond := testdata.NewDvp721SwapForEnygmaLogWith(
			testdata.WithDvp721SwapForEnygmaBlockNumber(wantBlockNumber),
			testdata.WithDvp721SwapForEnygmaNftResourceID(wantNftResourceIDHashSecond),
			testdata.WithDvp721SwapForEnygmaNftId(wantNftIdSecond),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V:   logrouter.Block{Number: wantBlockNumber, Logs: []ethTypes.Log{logFirst, logSecond}},
				Ack: ackSpy.Fn(),
			}),
		}
		eventMQ := newNoopEnygmaEventMQ()
		zkDvpBatchMQ := &DvpBatchMQMock{
			PushFunc: func(ctx context.Context, batch service.DvpSerializedEventBatch) error {
				assert.Equal(t, wantBlockNumber, batch.BlockNumber)
				assert.Equal(t, service.Dvp721SwapForEnygmaEvent, batch.Type)
				events, err := deserializeEvents[service.Dvp721SwapForEnygma](batch.SerializedEvents)
				assert.NoError(t, err)
				assert.Equal(t, wantEvents, events)
				return nil
			},
		}

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 0, len(eventMQ.PushBatchChunkedCalls()))
		assert.Equal(t, 1, len(zkDvpBatchMQ.PushCalls()))
	})

	t.Run("batches two Dvp1155Creation events", func(t *testing.T) {
		wantResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantResourceIDFirst := wantResourceIDHashFirst.Hex()[2:]
		wantResourceIDHashSecond := common.HexToHash("0xdeadbeef")
		wantResourceIDSecond := wantResourceIDHashSecond.Hex()[2:]
		wantBlockNumber := uint64(100)

		wantEvents := []service.Dvp1155Creation{
			{ChainEventID: "100-12-11", ResourceId: wantResourceIDFirst},
			{ChainEventID: "100-12-11", ResourceId: wantResourceIDSecond},
		}

		logFirst := testdata.NewDvp1155CreationLogWith(
			testdata.WithDvp1155CreationBlockNumber(wantBlockNumber),
			testdata.WithDvp1155CreationResourceID(wantResourceIDHashFirst),
		)
		logSecond := testdata.NewDvp1155CreationLogWith(
			testdata.WithDvp1155CreationBlockNumber(wantBlockNumber),
			testdata.WithDvp1155CreationResourceID(wantResourceIDHashSecond),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V:   logrouter.Block{Number: wantBlockNumber, Logs: []ethTypes.Log{logFirst, logSecond}},
				Ack: ackSpy.Fn(),
			}),
		}
		eventMQ := newNoopEnygmaEventMQ()
		zkDvpBatchMQ := &DvpBatchMQMock{
			PushFunc: func(ctx context.Context, batch service.DvpSerializedEventBatch) error {
				assert.Equal(t, wantBlockNumber, batch.BlockNumber)
				assert.Equal(t, service.Dvp1155CreationEvent, batch.Type)
				events, err := deserializeEvents[service.Dvp1155Creation](batch.SerializedEvents)
				assert.NoError(t, err)
				assert.Equal(t, wantEvents, events)
				return nil
			},
		}

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 0, len(eventMQ.PushBatchChunkedCalls()))
		assert.Equal(t, 1, len(zkDvpBatchMQ.PushCalls()))
	})

	t.Run("batches two Dvp1155Mint events", func(t *testing.T) {
		wantResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantResourceIDFirst := wantResourceIDHashFirst.Hex()[2:]
		wantResourceIDHashSecond := common.HexToHash("0xdeadbeef")
		wantResourceIDSecond := wantResourceIDHashSecond.Hex()[2:]
		wantBlockNumber := uint64(100)

		wantTokenIdFirst := big.NewInt(1)
		wantTokenIdSecond := big.NewInt(2)

		wantEvents := []service.Dvp1155Mint{
			{ChainEventID: "100-13-12", ResourceId: wantResourceIDFirst, TokenId: wantTokenIdFirst, Value: big.NewInt(100), Data: []byte{}},
			{ChainEventID: "100-13-12", ResourceId: wantResourceIDSecond, TokenId: wantTokenIdSecond, Value: big.NewInt(100), Data: []byte{}},
		}

		logFirst := testdata.NewDvp1155MintLogWith(
			testdata.WithDvp1155MintBlockNumber(wantBlockNumber),
			testdata.WithDvp1155MintResourceID(wantResourceIDHashFirst),
			testdata.WithDvp1155MintTokenId(wantTokenIdFirst),
		)
		logSecond := testdata.NewDvp1155MintLogWith(
			testdata.WithDvp1155MintBlockNumber(wantBlockNumber),
			testdata.WithDvp1155MintResourceID(wantResourceIDHashSecond),
			testdata.WithDvp1155MintTokenId(wantTokenIdSecond),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V:   logrouter.Block{Number: wantBlockNumber, Logs: []ethTypes.Log{logFirst, logSecond}},
				Ack: ackSpy.Fn(),
			}),
		}
		eventMQ := newNoopEnygmaEventMQ()
		zkDvpBatchMQ := &DvpBatchMQMock{
			PushFunc: func(ctx context.Context, batch service.DvpSerializedEventBatch) error {
				assert.Equal(t, wantBlockNumber, batch.BlockNumber)
				assert.Equal(t, service.Dvp1155MintEvent, batch.Type)
				events, err := deserializeEvents[service.Dvp1155Mint](batch.SerializedEvents)
				assert.NoError(t, err)
				assert.Equal(t, wantEvents, events)
				return nil
			},
		}

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 0, len(eventMQ.PushBatchChunkedCalls()))
		assert.Equal(t, 1, len(zkDvpBatchMQ.PushCalls()))
	})

	t.Run("batches two Dvp1155Burn events", func(t *testing.T) {
		wantResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantResourceIDFirst := wantResourceIDHashFirst.Hex()[2:]
		wantResourceIDHashSecond := common.HexToHash("0xdeadbeef")
		wantResourceIDSecond := wantResourceIDHashSecond.Hex()[2:]
		wantBlockNumber := uint64(100)

		wantTokenIdFirst := big.NewInt(1)
		wantTokenIdSecond := big.NewInt(2)

		wantEvents := []service.Dvp1155Burn{
			{ChainEventID: "100-14-13", ResourceId: wantResourceIDFirst, TokenId: wantTokenIdFirst, Value: big.NewInt(100)},
			{ChainEventID: "100-14-13", ResourceId: wantResourceIDSecond, TokenId: wantTokenIdSecond, Value: big.NewInt(100)},
		}

		logFirst := testdata.NewDvp1155BurnLogWith(
			testdata.WithDvp1155BurnBlockNumber(wantBlockNumber),
			testdata.WithDvp1155BurnResourceID(wantResourceIDHashFirst),
			testdata.WithDvp1155BurnTokenId(wantTokenIdFirst),
		)
		logSecond := testdata.NewDvp1155BurnLogWith(
			testdata.WithDvp1155BurnBlockNumber(wantBlockNumber),
			testdata.WithDvp1155BurnResourceID(wantResourceIDHashSecond),
			testdata.WithDvp1155BurnTokenId(wantTokenIdSecond),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V:   logrouter.Block{Number: wantBlockNumber, Logs: []ethTypes.Log{logFirst, logSecond}},
				Ack: ackSpy.Fn(),
			}),
		}
		eventMQ := newNoopEnygmaEventMQ()
		zkDvpBatchMQ := &DvpBatchMQMock{
			PushFunc: func(ctx context.Context, batch service.DvpSerializedEventBatch) error {
				assert.Equal(t, wantBlockNumber, batch.BlockNumber)
				assert.Equal(t, service.Dvp1155BurnEvent, batch.Type)
				events, err := deserializeEvents[service.Dvp1155Burn](batch.SerializedEvents)
				assert.NoError(t, err)
				assert.Equal(t, wantEvents, events)
				return nil
			},
		}

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 0, len(eventMQ.PushBatchChunkedCalls()))
		assert.Equal(t, 1, len(zkDvpBatchMQ.PushCalls()))
	})

	t.Run("batches two Dvp1155DepositIntoDvp events", func(t *testing.T) {
		wantResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantResourceIDFirst := wantResourceIDHashFirst.Hex()[2:]
		wantResourceIDHashSecond := common.HexToHash("0xdeadbeef")
		wantResourceIDSecond := wantResourceIDHashSecond.Hex()[2:]
		wantBlockNumber := uint64(100)

		wantTokenIdFirst := big.NewInt(1)
		wantTokenIdSecond := big.NewInt(2)

		wantEvents := []service.Dvp1155DepositIntoDvp{
			{ChainEventID: "100-15-14", ResourceId: wantResourceIDFirst, TokenId: wantTokenIdFirst, Value: big.NewInt(100), Data: []byte{}, From: common.HexToAddress("0x1111111111111111111111111111111111111111"), TxHash: common.HexToHash("0xabcdef14").Hex(), TxBlockNumber: big.NewInt(int64(wantBlockNumber))},
			{ChainEventID: "100-15-14", ResourceId: wantResourceIDSecond, TokenId: wantTokenIdSecond, Value: big.NewInt(100), Data: []byte{}, From: common.HexToAddress("0x1111111111111111111111111111111111111111"), TxHash: common.HexToHash("0xabcdef14").Hex(), TxBlockNumber: big.NewInt(int64(wantBlockNumber))},
		}

		logFirst := testdata.NewDvp1155DepositIntoDvpLogWith(
			testdata.WithDvp1155DepositIntoDvpBlockNumber(wantBlockNumber),
			testdata.WithDvp1155DepositIntoDvpResourceID(wantResourceIDHashFirst),
			testdata.WithDvp1155DepositIntoDvpTokenId(wantTokenIdFirst),
		)
		logSecond := testdata.NewDvp1155DepositIntoDvpLogWith(
			testdata.WithDvp1155DepositIntoDvpBlockNumber(wantBlockNumber),
			testdata.WithDvp1155DepositIntoDvpResourceID(wantResourceIDHashSecond),
			testdata.WithDvp1155DepositIntoDvpTokenId(wantTokenIdSecond),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V:   logrouter.Block{Number: wantBlockNumber, Logs: []ethTypes.Log{logFirst, logSecond}},
				Ack: ackSpy.Fn(),
			}),
		}
		eventMQ := newNoopEnygmaEventMQ()
		zkDvpBatchMQ := &DvpBatchMQMock{
			PushFunc: func(ctx context.Context, batch service.DvpSerializedEventBatch) error {
				assert.Equal(t, wantBlockNumber, batch.BlockNumber)
				assert.Equal(t, service.Dvp1155DepositIntoDvpEvent, batch.Type)
				events, err := deserializeEvents[service.Dvp1155DepositIntoDvp](batch.SerializedEvents)
				assert.NoError(t, err)
				assert.Equal(t, wantEvents, events)
				return nil
			},
		}

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 0, len(eventMQ.PushBatchChunkedCalls()))
		assert.Equal(t, 1, len(zkDvpBatchMQ.PushCalls()))
	})

	t.Run("batches two Dvp1155WithdrawFromDvp events", func(t *testing.T) {
		wantResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantResourceIDFirst := wantResourceIDHashFirst.Hex()[2:]
		wantResourceIDHashSecond := common.HexToHash("0xdeadbeef")
		wantResourceIDSecond := wantResourceIDHashSecond.Hex()[2:]
		wantBlockNumber := uint64(100)

		wantTokenIdFirst := big.NewInt(1)
		wantTokenIdSecond := big.NewInt(2)

		wantEvents := []service.Dvp1155WithdrawFromDvp{
			{ChainEventID: "100-16-15", ResourceId: wantResourceIDFirst, TokenId: wantTokenIdFirst, Value: big.NewInt(100), Owner: common.HexToAddress("0x1111111111111111111111111111111111111111"), TxHash: common.HexToHash("0xabcdef15").Hex(), TxBlockNumber: big.NewInt(int64(wantBlockNumber))},
			{ChainEventID: "100-16-15", ResourceId: wantResourceIDSecond, TokenId: wantTokenIdSecond, Value: big.NewInt(100), Owner: common.HexToAddress("0x1111111111111111111111111111111111111111"), TxHash: common.HexToHash("0xabcdef15").Hex(), TxBlockNumber: big.NewInt(int64(wantBlockNumber))},
		}

		logFirst := testdata.NewDvp1155WithdrawFromDvpLogWith(
			testdata.WithDvp1155WithdrawFromDvpBlockNumber(wantBlockNumber),
			testdata.WithDvp1155WithdrawFromDvpResourceID(wantResourceIDHashFirst),
			testdata.WithDvp1155WithdrawFromDvpTokenId(wantTokenIdFirst),
		)
		logSecond := testdata.NewDvp1155WithdrawFromDvpLogWith(
			testdata.WithDvp1155WithdrawFromDvpBlockNumber(wantBlockNumber),
			testdata.WithDvp1155WithdrawFromDvpResourceID(wantResourceIDHashSecond),
			testdata.WithDvp1155WithdrawFromDvpTokenId(wantTokenIdSecond),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V:   logrouter.Block{Number: wantBlockNumber, Logs: []ethTypes.Log{logFirst, logSecond}},
				Ack: ackSpy.Fn(),
			}),
		}
		eventMQ := newNoopEnygmaEventMQ()
		zkDvpBatchMQ := &DvpBatchMQMock{
			PushFunc: func(ctx context.Context, batch service.DvpSerializedEventBatch) error {
				assert.Equal(t, wantBlockNumber, batch.BlockNumber)
				assert.Equal(t, service.Dvp1155WithdrawFromDvpEvent, batch.Type)
				events, err := deserializeEvents[service.Dvp1155WithdrawFromDvp](batch.SerializedEvents)
				assert.NoError(t, err)
				assert.Equal(t, wantEvents, events)
				return nil
			},
		}

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 0, len(eventMQ.PushBatchChunkedCalls()))
		assert.Equal(t, 1, len(zkDvpBatchMQ.PushCalls()))
	})

	t.Run("batches two Dvp1155SwapForEnygma events", func(t *testing.T) {
		wantTokenResourceIDHashFirst := common.HexToHash("0xc001babe")
		wantTokenResourceIDFirst := wantTokenResourceIDHashFirst.Hex()[2:]
		wantTokenResourceIDHashSecond := common.HexToHash("0xdeadbeef")
		wantTokenResourceIDSecond := wantTokenResourceIDHashSecond.Hex()[2:]
		wantBlockNumber := uint64(100)

		wantTokenIdFirst := big.NewInt(1)
		wantTokenIdSecond := big.NewInt(2)

		wantEvents := []service.Dvp1155SwapForEnygma{
			{
				SharedId:         "0102030000000000000000000000000000000000000000000000000000000000",
				DestChainId:      big.NewInt(1),
				From:             common.HexToAddress("0x1111111111111111111111111111111111111111"),
				TokenResourceId:  wantTokenResourceIDFirst,
				TokenValue:       big.NewInt(100),
				TokenId:          wantTokenIdFirst.String(),
				EnygmaResourceId: "0d0e0f0000000000000000000000000000000000000000000000000000000000",
				EnygmaAmount:     big.NewInt(1000),
				TxHash:           common.HexToHash("0xabcdef16").Hex(),
				//nolint:gosec // test data with known values
				TxBlockNumber: big.NewInt(int64(wantBlockNumber)),
			},
			{
				SharedId:         "0102030000000000000000000000000000000000000000000000000000000000",
				DestChainId:      big.NewInt(1),
				From:             common.HexToAddress("0x1111111111111111111111111111111111111111"),
				TokenResourceId:  wantTokenResourceIDSecond,
				TokenValue:       big.NewInt(100),
				TokenId:          wantTokenIdSecond.String(),
				EnygmaResourceId: "0d0e0f0000000000000000000000000000000000000000000000000000000000",
				EnygmaAmount:     big.NewInt(1000),
				TxHash:           common.HexToHash("0xabcdef16").Hex(),
				//nolint:gosec // test data with known values
				TxBlockNumber: big.NewInt(int64(wantBlockNumber)),
			},
		}

		logFirst := testdata.NewDvp1155SwapForEnygmaLogWith(
			testdata.WithDvp1155SwapForEnygmaBlockNumber(wantBlockNumber),
			testdata.WithDvp1155SwapForEnygmaTokenResourceID(wantTokenResourceIDHashFirst),
			testdata.WithDvp1155SwapForEnygmaTokenId(wantTokenIdFirst),
		)
		logSecond := testdata.NewDvp1155SwapForEnygmaLogWith(
			testdata.WithDvp1155SwapForEnygmaBlockNumber(wantBlockNumber),
			testdata.WithDvp1155SwapForEnygmaTokenResourceID(wantTokenResourceIDHashSecond),
			testdata.WithDvp1155SwapForEnygmaTokenId(wantTokenIdSecond),
		)

		ackSpy := spy.NewAck()

		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V:   logrouter.Block{Number: wantBlockNumber, Logs: []ethTypes.Log{logFirst, logSecond}},
				Ack: ackSpy.Fn(),
			}),
		}
		eventMQ := newNoopEnygmaEventMQ()
		zkDvpBatchMQ := &DvpBatchMQMock{
			PushFunc: func(ctx context.Context, batch service.DvpSerializedEventBatch) error {
				assert.Equal(t, wantBlockNumber, batch.BlockNumber)
				assert.Equal(t, service.Dvp1155SwapForEnygmaEvent, batch.Type)
				events, err := deserializeEvents[service.Dvp1155SwapForEnygma](batch.SerializedEvents)
				assert.NoError(t, err)
				assert.Equal(t, wantEvents, events)
				return nil
			},
		}

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		_ = parser.Fetch(ctx)

		assert.Equal(t, 2, len(blockMQ.NextCalls()))
		assert.Equal(t, 0, len(eventMQ.PushBatchChunkedCalls()))
		assert.Equal(t, 1, len(zkDvpBatchMQ.PushCalls()))
	})
}

func deserializeEvents[T any](serialized []byte) ([]T, error) {
	var events []T

	err := json.Unmarshal(serialized, &events)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal events")
	}
	return events, nil
}

// Log builder helpers - convert service events to ethTypes.Log

// buildEnygmaCreationLog creates a ethTypes.Log from an EnygmaCreation service event
func buildEnygmaCreationLog(event service.EnygmaCreation, blockNumber uint64, resourceID string) ethTypes.Log {
	// Convert hex string resourceID to [32]byte
	resourceIDBytes := common.HexToHash(resourceID)

	return testdata.NewEnygmaCreationLogWith(
		testdata.WithCreateResourceID(resourceIDBytes),
		testdata.WithCreateInitialSupply(event.InitialSupply),
		testdata.WithCreateBlockNumber(blockNumber),
	)
}

// buildEnygmaDepositToDvpLog creates a ethTypes.Log from an EnygmaDepositToDvp service event
func buildEnygmaDepositToDvpLog(event service.EnygmaDepositToDvp, blockNumber uint64, resourceID string) ethTypes.Log {
	// Convert hex string resourceID to [32]byte
	resourceIDBytes := common.HexToHash(resourceID)

	return testdata.NewEnygmaDepositToDvpLogWith(
		testdata.WithDepositResourceID(resourceIDBytes),
		testdata.WithDepositAmount(event.Amount),
		testdata.WithDepositFrom(event.From),
		testdata.WithDepositReferenceID(event.ReferenceId),
		testdata.WithDepositBlockNumber(blockNumber),
	)
}

func TestEnygmaParser_NoAckOnErrors(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("does not ack block when enygma batch push fails", func(t *testing.T) {
		wantBlockNumber := uint64(100)

		log := testdata.NewEnygmaCreationLogWith(
			testdata.WithCreateBlockNumber(wantBlockNumber),
			testdata.WithCreateResourceID(common.HexToHash("0xc001babe")),
			testdata.WithCreateInitialSupply(big.NewInt(42)),
		)

		ackSpy := spy.NewAck()
		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   []ethTypes.Log{log},
				},
				Ack: ackSpy.Fn(),
			}),
		}

		eventMQ := &EnygmaEventMQMock{
			PushBatchChunkedFunc: func(ctx context.Context, events []service.EnygmaSerializedEvent) error {
				return fmt.Errorf("NATS unavailable")
			},
		}
		zkDvpBatchMQ := newNoopDvpBatchMQ()

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		ackSpy.AssertNotCalled(t, "should not ack block when enygma batch push fails, allowing MQ redelivery")
	})

	t.Run("does not ack block when dvp batch push fails", func(t *testing.T) {
		wantBlockNumber := uint64(100)

		log := testdata.NewDvp721CreationLogWith(
			testdata.WithDvp721CreationBlockNumber(wantBlockNumber),
		)

		ackSpy := spy.NewAck()
		blockMQ := &EnygmaBlockMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: wantBlockNumber,
					Logs:   []ethTypes.Log{log},
				},
				Ack: ackSpy.Fn(),
			}),
		}

		eventMQ := newNoopEnygmaEventMQ()
		zkDvpBatchMQ := &DvpBatchMQMock{
			PushFunc: func(ctx context.Context, batch service.DvpSerializedEventBatch) error {
				return fmt.Errorf("NATS unavailable")
			},
		}

		parser := logparser.NewEnygmaLogParser(blockMQ, eventMQ, zkDvpBatchMQ)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		ackSpy.AssertNotCalled(t, "should not ack block when dvp batch push fails, allowing MQ redelivery")
	})
}
