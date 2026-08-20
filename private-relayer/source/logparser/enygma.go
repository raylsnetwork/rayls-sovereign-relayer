package logparser

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EnygmaPNEvents"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/logrouter"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/service"
)

//go:generate moq --pkg logparser_test -out enygma_mock_test.go . EnygmaBlockMQ EnygmaEventMQ DvpBatchMQ
type EnygmaBlockMQ interface {
	Next(ctx context.Context) (msgqueue.Message[logrouter.Block], error)
}

type EnygmaEventMQ interface {
	Push(context.Context, service.EnygmaSerializedEvent) error
	PushBatchChunked(context.Context, []service.EnygmaSerializedEvent) error
}

type DvpBatchMQ interface {
	Push(context.Context, service.DvpSerializedEventBatch) error
}

type EnygmaLogParser struct {
	enygmaMQ     EnygmaEventMQ
	zkDvpBatchMQ DvpBatchMQ
	blockMQ      EnygmaBlockMQ

	filterer *EnygmaPNEvents.EnygmaPNEvents
}

func NewEnygmaLogParser(blockMQ EnygmaBlockMQ, enygmaMQ EnygmaEventMQ, zkDvpBatchMQ DvpBatchMQ) *EnygmaLogParser {
	filterer := EnygmaPNEvents.NewEnygmaPNEvents()
	return &EnygmaLogParser{
		enygmaMQ:     enygmaMQ,
		zkDvpBatchMQ: zkDvpBatchMQ,
		blockMQ:      blockMQ,

		filterer: filterer,
	}
}

func (p *EnygmaLogParser) Fetch( //nolint:gocognit,gocyclo // log parser with event type dispatching
	ctx context.Context,
) error {
	for {
		msg, err := p.blockMQ.Next(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil //nolint:nilerr // context cancellation is a clean shutdown signal
			}
			slog.Error("Failed to get next message from endpoint MQ", slog.Any("error", err))
			continue
		}

		dvpBatches := service.DvpEventBatchBuilder{}
		enygmaEvents := make([]service.EnygmaTypedEvent[any], 0)
		var hasProcessingError bool
		for _, log := range msg.V.Logs {
			// Enygma Events
			if contractEvent, err := p.filterer.UnpackEnygmaCreationEvent(&log); err == nil {
				serviceEvent := convertEnygmaCreation(contractEvent)
				resourceID := common.Bytes2Hex(contractEvent.ResourceId[:])

				enygmaEvents = append(enygmaEvents, service.NewEnygmaTypedEvent[any](
					log.BlockNumber,
					log.Index,
					log.TxHash,
					resourceID,
					service.EnygmaCreationEvent,
					*serviceEvent,
				))
			} else if contractEvent, err := p.filterer.UnpackEnygmaSendTransferPNHEvent(&log); err == nil {
				serviceEvents := convertEnygmaSendTransferCC(contractEvent)
				resourceID := common.Bytes2Hex(contractEvent.ResourceId[:])
				for _, serviceEvent := range serviceEvents {
					enygmaEvents = append(enygmaEvents, service.NewEnygmaTypedEvent[any](
						log.BlockNumber,
						log.Index,
						log.TxHash,
						resourceID,
						service.EnygmaTransferEvent,
						serviceEvent,
					))
				}
			} else if contractEvent, err := p.filterer.UnpackEnygmaMintEvent(&log); err == nil {
				convertedEvent := convertEnygmaMint(contractEvent)
				serviceEvent := service.EnygmaSupplyUpdate{
					TxHash: convertedEvent.TxHash,
					To:     contractEvent.To,
					Amount: convertedEvent.Amount,
				}
				resourceID := common.Bytes2Hex(contractEvent.ResourceId[:])

				enygmaEvents = append(enygmaEvents, service.NewEnygmaTypedEvent[any](
					log.BlockNumber,
					log.Index,
					log.TxHash,
					resourceID,
					service.EnygmaSupplyUpdateEvent,
					serviceEvent,
				))
			} else if contractEvent, err := p.filterer.UnpackEnygmaBurnEvent(&log); err == nil {
				negativeAmount := new(big.Int).Neg(contractEvent.Amount)

				convertedEvent := convertEnygmaBurn(contractEvent)
				serviceEvent := service.EnygmaSupplyUpdate{
					TxHash: convertedEvent.TxHash,
					To:     contractEvent.From,
					Amount: negativeAmount,
				}
				resourceID := common.Bytes2Hex(contractEvent.ResourceId[:])

				enygmaEvents = append(enygmaEvents, service.NewEnygmaTypedEvent[any](
					log.BlockNumber,
					log.Index,
					log.TxHash,
					resourceID,
					service.EnygmaSupplyUpdateEvent,
					serviceEvent,
				))
			} else if contractEvent, err := p.filterer.UnpackEnygmaDepositToDvpEvent(&log); err == nil {
				serviceEvent := convertEnygmaDepositToDvp(contractEvent)
				resourceID := common.Bytes2Hex(contractEvent.ResourceId[:])

				enygmaEvents = append(enygmaEvents, service.NewEnygmaTypedEvent[any](
					log.BlockNumber,
					log.Index,
					log.TxHash,
					resourceID,
					service.EnygmaDepositEvent,
					*serviceEvent,
				))
			} else if contractEvent, err := p.filterer.UnpackEnygmaWithdrawFromDvpEvent(&log); err == nil {
				serviceEvent := convertEnygmaWithdrawFromDvp(contractEvent)
				resourceID := common.Bytes2Hex(contractEvent.ResourceId[:])

				enygmaEvents = append(enygmaEvents, service.NewEnygmaTypedEvent[any](
					log.BlockNumber,
					log.Index,
					log.TxHash,
					resourceID,
					service.EnygmaWithdrawEvent,
					*serviceEvent,
				))
			} else

			// Dvp ERC721 Events
			if contractEvent, err := p.filterer.UnpackDvp721CreationEvent(&log); err == nil {
				serviceEvent := convertDvp721Creation(contractEvent)
				tryPushDvpEventToBatch(dvpBatches, msg.V.Number, service.Dvp721CreationEvent, *serviceEvent)
			} else if contractEvent, err := p.filterer.UnpackDvp721MintEvent(&log); err == nil {
				serviceEvent := convertDvp721Mint(contractEvent)
				tryPushDvpEventToBatch(dvpBatches, msg.V.Number, service.Dvp721MintEvent, *serviceEvent)
			} else if contractEvent, err := p.filterer.UnpackDvp721BurnEvent(&log); err == nil {
				serviceEvent := convertDvp721Burn(contractEvent)
				tryPushDvpEventToBatch(dvpBatches, msg.V.Number, service.Dvp721BurnEvent, *serviceEvent)
			} else if contractEvent, err := p.filterer.UnpackDvp721DepositIntoDvpEvent(&log); err == nil {
				serviceEvent := convertDvp721DepositIntoDvp(contractEvent)
				tryPushDvpEventToBatch(dvpBatches, msg.V.Number, service.Dvp721DepositIntoDvpEvent, *serviceEvent)
			} else if contractEvent, err := p.filterer.UnpackDvp721WithdrawFromDvpEvent(&log); err == nil {
				serviceEvent := convertDvp721WithdrawFromDvp(contractEvent)
				tryPushDvpEventToBatch(dvpBatches, msg.V.Number, service.Dvp721WithdrawFromDvpEvent, *serviceEvent)
			} else if contractEvent, err := p.filterer.UnpackDvp721SwapForEnygmaEvent(&log); err == nil {
				serviceEvent := convertDvp721SwapForEnygma(contractEvent)
				tryPushDvpEventToBatch(dvpBatches, msg.V.Number, service.Dvp721SwapForEnygmaEvent, *serviceEvent)
			} else

			// Dvp ERC1155 Events
			if contractEvent, err := p.filterer.UnpackDvp1155CreationEvent(&log); err == nil {
				serviceEvent := convertDvp1155Creation(contractEvent)
				tryPushDvpEventToBatch(dvpBatches, msg.V.Number, service.Dvp1155CreationEvent, *serviceEvent)
			} else if contractEvent, err := p.filterer.UnpackDvp1155MintEvent(&log); err == nil {
				serviceEvent := convertDvp1155Mint(contractEvent)
				tryPushDvpEventToBatch(dvpBatches, msg.V.Number, service.Dvp1155MintEvent, *serviceEvent)
			} else if contractEvent, err := p.filterer.UnpackDvp1155BurnEvent(&log); err == nil {
				serviceEvent := convertDvp1155Burn(contractEvent)
				tryPushDvpEventToBatch(dvpBatches, msg.V.Number, service.Dvp1155BurnEvent, *serviceEvent)
			} else if contractEvent, err := p.filterer.UnpackDvp1155DepositIntoDvpEvent(&log); err == nil {
				serviceEvent := convertDvp1155DepositIntoDvp(contractEvent)
				tryPushDvpEventToBatch(dvpBatches, msg.V.Number, service.Dvp1155DepositIntoDvpEvent, *serviceEvent)
			} else if contractEvent, err := p.filterer.UnpackDvp1155WithdrawFromDvpEvent(&log); err == nil {
				serviceEvent := convertDvp1155WithdrawFromDvp(contractEvent)
				tryPushDvpEventToBatch(dvpBatches, msg.V.Number, service.Dvp1155WithdrawFromDvpEvent, *serviceEvent)
			} else if contractEvent, err := p.filterer.UnpackDvp1155SwapForEnygmaEvent(&log); err == nil {
				serviceEvent := convertDvp1155SwapForEnygma(contractEvent)
				tryPushDvpEventToBatch(dvpBatches, msg.V.Number, service.Dvp1155SwapForEnygmaEvent, *serviceEvent)
			} else

			// Enygma Swap Event
			if contractEvent, err := p.filterer.UnpackEnygmaSwapWithDvpForERC721Event(&log); err == nil {
				serviceEvent := convertEnygmaSwapWithDvpForERC721(contractEvent)
				tryPushDvpEventToBatch(dvpBatches, msg.V.Number, service.DvpEnygmaSwapERC721Event, *serviceEvent)
			} else if contractEvent, err := p.filterer.UnpackEnygmaSwapWithDvpForERC1155Event(&log); err == nil {
				serviceEvent := convertEnygmaSwapWithDvpForERC1155(contractEvent)
				tryPushDvpEventToBatch(dvpBatches, msg.V.Number, service.DvpEnygmaSwapERC1155Event, *serviceEvent)
			} else

			// Swap Cancelled Event
			if contractEvent, err := p.filterer.UnpackDvpSwapCancelledEvent(&log); err == nil {
				serviceEvent := convertDvpSwapCancelled(contractEvent)
				tryPushDvpEventToBatch(dvpBatches, msg.V.Number, service.DvpSwapCancelledEvent, *serviceEvent)
			}
		}

		enygmaBatch := make([]service.EnygmaSerializedEvent, 0)
		for _, event := range enygmaEvents {
			serialized, err := event.Serialize()
			if err != nil {
				slog.Error("Failed to serialize enygma event", slog.Any("error", err))
				hasProcessingError = true
				continue
			}
			enygmaBatch = append(enygmaBatch, serialized)
		}

		if len(enygmaBatch) > 0 {
			// Chunked: a single high-throughput block can carry more Enygma events
			// than the NATS atomic-batch cap (msgqueue.maxBatchSize=1000). A plain
			// PushBatch would reject the whole block with ErrTooManyMessages, the block
			// would never ack, and it would redeliver to MaxDeliver and be dropped.
			if err := p.enygmaMQ.PushBatchChunked(ctx, enygmaBatch); err != nil {
				slog.Error("Failed to push enygma batch to MQ", slog.Any("error", err))
				hasProcessingError = true
			}
		}

		for evType := range service.DvpEventTypeCount {
			if batch, ok := dvpBatches[evType]; ok {
				serializedBatch, err := batch.Serialize()
				if err != nil {
					slog.Error("Failed to serialize dvp batch",
						slog.Uint64("block_number", msg.V.Number),
						slog.Any("event_type", evType),
						slog.Any("error", err))
					hasProcessingError = true
					continue
				}
				if err := p.zkDvpBatchMQ.Push(ctx, serializedBatch); err != nil {
					slog.Error("Failed to push dvp batch to MQ",
						slog.Uint64("block_number", msg.V.Number),
						slog.Any("event_type", evType),
						slog.Any("error", err))
					hasProcessingError = true
				}
			}
		}

		if hasProcessingError {
			slog.Warn("Skipping ack due to processing errors, block will be redelivered",
				slog.Uint64("block_number", msg.V.Number))
			continue
		}

		_ = msg.Ack(ctx)
	}
}

func tryPushDvpEventToBatch[T any](
	dvpBatches service.DvpEventBatchBuilder,
	blockNumber uint64,
	eventType service.DvpEventType,
	event T,
) {
	batch, ok := dvpBatches[eventType]
	if !ok {
		batch = &service.DvpTypedEventBatch[T]{
			BlockNumber: blockNumber,
			Type:        eventType,
			Events:      []T{},
		}
		dvpBatches[eventType] = batch
	}
	err := batch.PushEvent(event)
	if err != nil {
		slog.Error(
			"Failed to push ZkDvP event to batch",
			slog.Any("Block Number", blockNumber),
			slog.Any("Event Type", eventType),
		)
	}
}
