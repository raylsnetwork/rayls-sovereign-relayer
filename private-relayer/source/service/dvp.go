package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"

	"github.com/ethereum/go-ethereum/common"
)

//go:generate moq --pkg service_test -out dvp_mock_test.go . DvpBatchMQ DvpInitiator
type DvpBatchMQ interface {
	Next(context.Context) (msgqueue.Message[DvpSerializedEventBatch], error)
}

type DvpInitiator interface {
	// Enygma Swap handlers
	HandleEnygmaSwapERC721(ctx context.Context, sharedId string, toChainId *big.Int, from common.Address, enygmaResourceId string, enygmaAmount *big.Int, nftResourceId string, nftId string, txHash string, txBlockNumber *big.Int, validityTime uint64) error
	HandleEnygmaSwapERC1155(ctx context.Context, sharedId string, toChainId *big.Int, from common.Address, enygmaResourceId string, enygmaAmount *big.Int, nftResourceId string, nftId string, nftAmount *big.Int, txHash string, txBlockNumber *big.Int, validityTime uint64) error

	// ERC721 handlers
	HandleERC721Creation(ctx context.Context, chainEventID string, resourceId string) error
	HandleERC721Mint(ctx context.Context, chainEventID string, resourceId string, nftId *big.Int) error
	HandleERC721Burn(ctx context.Context, chainEventID string, resourceId string, nftId *big.Int) error
	HandleERC721Deposit(ctx context.Context, chainEventID string, resourceId string, nftId *big.Int, from common.Address, txHash string, txBlockNumber *big.Int) error
	HandleERC721Withdraw(ctx context.Context, chainEventID string, resourceId string, nftId *big.Int, from common.Address, txHash string, txBlockNumber *big.Int) error
	HandleERC721SwapEnygma(ctx context.Context, sharedId string, toChainId *big.Int, from common.Address, nftResourceId string, nftId string, enygmaResourceId string, enygmaAmount *big.Int, txHash string, txBlockNumber *big.Int, validityTime uint64) error

	// ERC1155 handlers
	HandleERC1155Creation(ctx context.Context, chainEventID string, resourceId string) error
	HandleERC1155Mint(ctx context.Context, chainEventID string, resourceId string, tokenId *big.Int, tokenAmount *big.Int, tokenData []byte) error
	HandleERC1155Burn(ctx context.Context, chainEventID string, resourceId string, tokenId *big.Int, tokenAmount *big.Int) error
	HandleERC1155Deposit(ctx context.Context, chainEventID string, resourceId string, tokenId *big.Int, tokenAmount *big.Int, tokenData []byte, from common.Address, txHash string, txBlockNumber *big.Int) error
	HandleERC1155Withdraw(ctx context.Context, chainEventID string, resourceId string, tokenId *big.Int, tokenAmount *big.Int, from common.Address, txHash string, txBlockNumber *big.Int) error
	HandleERC1155SwapEnygma(ctx context.Context, sharedId string, toChainId *big.Int, from common.Address, erc1155ResourceId string, erc1155Amount *big.Int, erc1155Id string, enygmaResourceId string, enygmaAmount *big.Int, txHash string, txBlockNumber *big.Int, validityTime uint64) error

	// Swap cancellation handler
	HandleSwapCancellation(ctx context.Context, sharedId string, toChainId *big.Int, tokenInResourceId string, tokenInAmount *big.Int, tokenInId string, tokenInType types.DvpTokenType, tokenOutResourceId string, tokenOutAmount *big.Int, tokenOutId string, tokenOutType types.DvpTokenType) error
}

// defaultHandleMessageTimeout bounds processing of a single dvp message batch so one stuck handler
// (e.g. an on-chain call whose TX never mines, leaving bind.WaitMined polling forever on the long-lived
// run ctx) cannot head-of-line stall the single-threaded loop. On timeout the per-message ctx is
// cancelled, the handler's chain wait returns "context deadline exceeded", hasProcessingError is set
// (the same no-ack/redeliver path a handler error already takes), and the loop continues with the next
// message. Generous so normal mining latency never trips it.
const defaultHandleMessageTimeout = 2 * time.Minute

type DvpOrchestrator struct {
	dvpMQ     DvpBatchMQ
	initiator DvpInitiator

	handleMessageTimeout time.Duration
}

func NewDvpOrchestrator() *DvpOrchestrator {
	return &DvpOrchestrator{handleMessageTimeout: defaultHandleMessageTimeout}
}

func NewDvpOrchestratorWithDeps(dvpMQ DvpBatchMQ, initiator DvpInitiator) *DvpOrchestrator {
	return NewDvpOrchestratorWithDepsAndTimeout(dvpMQ, initiator, defaultHandleMessageTimeout)
}

// NewDvpOrchestratorWithDepsAndTimeout is like NewDvpOrchestratorWithDeps but lets callers (and tests)
// set the per-message handler timeout. A non-positive timeout falls back to the default.
func NewDvpOrchestratorWithDepsAndTimeout(dvpMQ DvpBatchMQ, initiator DvpInitiator, handleMessageTimeout time.Duration) *DvpOrchestrator {
	if handleMessageTimeout <= 0 {
		handleMessageTimeout = defaultHandleMessageTimeout
	}
	return &DvpOrchestrator{
		dvpMQ:                dvpMQ,
		initiator:            initiator,
		handleMessageTimeout: handleMessageTimeout,
	}
}

func (z *DvpOrchestrator) Run(ctx context.Context) error {
	slog.Info("DvpOrchestrator started")

	for {
		select {
		case <-ctx.Done():
			slog.Info("DvpOrchestrator shutting down")
			return nil
		default:
		}

		slog.Debug("Fetching next dvp message from queue")
		msg, err := z.dvpMQ.Next(ctx)
		if err != nil {
			slog.Warn("Failed to get next message from queue", slog.Any("error", err))
			continue
		}

		slog.Debug("Processing dvp event", slog.Int("event_type", int(msg.V.Type)))

		z.handleMessage(ctx, msg)
	}
}

// handleMessage processes one dvp message batch under a per-message bounded context, so one stuck handler
// (an on-chain call whose TX never mines → bind.WaitMined polling forever on the long-lived run ctx)
// cannot head-of-line stall the single-threaded Run loop. On timeout the handler's chain wait returns
// "context deadline exceeded" → hasProcessingError → the message is not acked and is redelivered (the same
// path a handler error already takes). Mirrors the named dispatch() helpers on the dest orchestrators.
func (z *DvpOrchestrator) handleMessage(parentCtx context.Context, msg msgqueue.Message[DvpSerializedEventBatch]) {
	ctx, cancel := context.WithTimeout(parentCtx, z.handleMessageTimeout)
	defer cancel()

	var hasProcessingError bool
	switch msg.V.Type {
	// Enygma Swap Events
	case DvpEnygmaSwapERC721Event:
		slog.Info("Processing Enygma swap ERC721 event batch")
		eventBatch, err := deserializeDvpMessage[DvpEnygmaSwapERC721](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize Enygma swap ERC721 events", slog.Any("error", err))
			return
		}
		slog.Debug("Processing Enygma swap ERC721 events", slog.Int("event_count", len(eventBatch.Events)))
		for _, event := range eventBatch.Events {
			if err := z.initiator.HandleEnygmaSwapERC721(ctx, event.SharedId, event.DestChainId, event.From, event.ResourceId, event.EnygmaAmount, event.NftResourceId, event.NftId, event.TxHash, event.TxBlockNumber, event.ValidityTime); err != nil {
				slog.Error("Error executing Enygma swap ERC721",
					slog.String("sharedId", event.SharedId),
					slog.Any("error", err))
				hasProcessingError = true
			} else {
				slog.Debug("Successfully handled Enygma swap ERC721", slog.String("sharedId", event.SharedId))
			}
		}

	case DvpEnygmaSwapERC1155Event:
		slog.Info("Processing Enygma swap ERC1155 event batch")
		eventBatch, err := deserializeDvpMessage[DvpEnygmaSwapERC1155](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize Enygma swap ERC1155 events", slog.Any("error", err))
			return
		}
		slog.Debug("Processing Enygma swap ERC1155 events", slog.Int("event_count", len(eventBatch.Events)))
		for _, event := range eventBatch.Events {
			if err := z.initiator.HandleEnygmaSwapERC1155(ctx, event.SharedId, event.DestChainId, event.From, event.ResourceId, event.EnygmaAmount, event.NftResourceId, event.NftId, event.NftAmountOrOne, event.TxHash, event.TxBlockNumber, event.ValidityTime); err != nil {
				slog.Error("Error executing Enygma swap ERC1155",
					slog.String("sharedId", event.SharedId),
					slog.Any("error", err))
				hasProcessingError = true
			} else {
				slog.Debug("Successfully handled Enygma swap ERC1155", slog.String("sharedId", event.SharedId))
			}
		}

	// ERC721 Events
	case Dvp721CreationEvent:
		slog.Info("Processing ERC721 creation event batch")
		eventBatch, err := deserializeDvpMessage[Dvp721Creation](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize ERC721 creation events", slog.Any("error", err))
			return
		}
		slog.Debug("Processing ERC721 creation events", slog.Int("event_count", len(eventBatch.Events)))
		for _, event := range eventBatch.Events {
			if err := z.initiator.HandleERC721Creation(ctx, event.ChainEventID, event.ResourceId); err != nil {
				slog.Error("Error executing ERC721 creation",
					slog.String("resourceId", event.ResourceId),
					slog.Any("error", err))
				hasProcessingError = true
			} else {
				slog.Debug("Successfully handled ERC721 creation", slog.String("resourceId", event.ResourceId))
			}
		}

	case Dvp721MintEvent:
		slog.Info("Processing ERC721 mint event batch")
		eventBatch, err := deserializeDvpMessage[Dvp721Mint](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize ERC721 mint events", slog.Any("error", err))
			return
		}
		slog.Debug("Processing ERC721 mint events", slog.Int("event_count", len(eventBatch.Events)))
		for _, event := range eventBatch.Events {
			if err := z.initiator.HandleERC721Mint(ctx, event.ChainEventID, event.ResourceId, event.NftId); err != nil {
				slog.Error("Error executing ERC721 mint",
					slog.String("resourceId", event.ResourceId),
					slog.String("nftId", event.NftId.String()),
					slog.Any("error", err))
				hasProcessingError = true
			} else {
				slog.Debug("Successfully handled ERC721 mint", slog.String("resourceId", event.ResourceId), slog.String("nftId", event.NftId.String()))
			}
		}

	case Dvp721BurnEvent:
		slog.Info("Processing ERC721 burn event batch")
		eventBatch, err := deserializeDvpMessage[Dvp721Burn](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize ERC721 burn events", slog.Any("error", err))
			return
		}
		slog.Debug("Processing ERC721 burn events", slog.Int("event_count", len(eventBatch.Events)))
		for _, event := range eventBatch.Events {
			if err := z.initiator.HandleERC721Burn(ctx, event.ChainEventID, event.ResourceId, event.NftId); err != nil {
				slog.Error("Error executing ERC721 burn",
					slog.String("resourceId", event.ResourceId),
					slog.String("nftId", event.NftId.String()),
					slog.Any("error", err))
				hasProcessingError = true
			} else {
				slog.Debug("Successfully handled ERC721 burn", slog.String("resourceId", event.ResourceId), slog.String("nftId", event.NftId.String()))
			}
		}

	case Dvp721DepositIntoDvpEvent:
		slog.Info("Processing ERC721 deposit event batch")
		eventBatch, err := deserializeDvpMessage[Dvp721DepositIntoDvp](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize ERC721 deposit events", slog.Any("error", err))
			return
		}
		slog.Debug("Processing ERC721 deposit events", slog.Int("event_count", len(eventBatch.Events)))
		for _, event := range eventBatch.Events {
			if err := z.initiator.HandleERC721Deposit(ctx, event.ChainEventID, event.ResourceId, event.NftId, event.From, event.TxHash, event.TxBlockNumber); err != nil {
				slog.Error("Error executing ERC721 deposit",
					slog.String("resourceId", event.ResourceId),
					slog.String("nftId", event.NftId.String()),
					slog.Any("error", err))
				hasProcessingError = true
			} else {
				slog.Debug("Successfully handled ERC721 deposit", slog.String("resourceId", event.ResourceId), slog.String("nftId", event.NftId.String()))
			}
		}

	case Dvp721WithdrawFromDvpEvent:
		slog.Info("Processing ERC721 withdraw event batch")
		eventBatch, err := deserializeDvpMessage[Dvp721WithdrawFromDvp](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize ERC721 withdraw events", slog.Any("error", err))
			return
		}
		slog.Debug("Processing ERC721 withdraw events", slog.Int("event_count", len(eventBatch.Events)))
		for _, event := range eventBatch.Events {
			if err := z.initiator.HandleERC721Withdraw(ctx, event.ChainEventID, event.ResourceId, event.NftId, event.Owner, event.TxHash, event.TxBlockNumber); err != nil {
				slog.Error("Error executing ERC721 withdraw",
					slog.String("resourceId", event.ResourceId),
					slog.String("nftId", event.NftId.String()),
					slog.Any("error", err))
				hasProcessingError = true
			} else {
				slog.Debug("Successfully handled ERC721 withdraw", slog.String("resourceId", event.ResourceId), slog.String("nftId", event.NftId.String()))
			}
		}

	case Dvp721SwapForEnygmaEvent:
		slog.Info("Processing ERC721 swap for Enygma event batch")
		eventBatch, err := deserializeDvpMessage[Dvp721SwapForEnygma](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize ERC721 swap for Enygma events", slog.Any("error", err))
			return
		}
		slog.Debug("Processing ERC721 swap for Enygma events", slog.Int("event_count", len(eventBatch.Events)))
		for _, event := range eventBatch.Events {
			if err := z.initiator.HandleERC721SwapEnygma(ctx, event.SharedId, event.DestChainId, event.From, event.NftResourceId, event.NftId, event.EnygmaResourceId, event.EnygmaAmount, event.TxHash, event.TxBlockNumber, event.ValidityTime); err != nil {
				slog.Error("Error executing Swap ERC721 for Enygma",
					slog.String("sharedId", event.SharedId),
					slog.Any("error", err))
				hasProcessingError = true
			} else {
				slog.Debug("Successfully handled ERC721 swap for Enygma", slog.String("sharedId", event.SharedId))
			}
		}

	// ERC1155 Events
	case Dvp1155CreationEvent:
		slog.Info("Processing ERC1155 creation event batch")
		eventBatch, err := deserializeDvpMessage[Dvp1155Creation](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize ERC1155 creation events", slog.Any("error", err))
			return
		}
		slog.Debug("Processing ERC1155 creation events", slog.Int("event_count", len(eventBatch.Events)))
		for _, event := range eventBatch.Events {
			if err := z.initiator.HandleERC1155Creation(ctx, event.ChainEventID, event.ResourceId); err != nil {
				slog.Error("Error executing ERC1155 creation",
					slog.String("resourceId", event.ResourceId),
					slog.Any("error", err))
				hasProcessingError = true
			} else {
				slog.Debug("Successfully handled ERC1155 creation", slog.String("resourceId", event.ResourceId))
			}
		}

	case Dvp1155MintEvent:
		slog.Info("Processing ERC1155 mint event batch")
		eventBatch, err := deserializeDvpMessage[Dvp1155Mint](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize ERC1155 mint events", slog.Any("error", err))
			return
		}
		slog.Debug("Processing ERC1155 mint events", slog.Int("event_count", len(eventBatch.Events)))
		for _, event := range eventBatch.Events {
			if err := z.initiator.HandleERC1155Mint(ctx, event.ChainEventID, event.ResourceId, event.TokenId, event.Value, event.Data); err != nil {
				slog.Error("Error executing ERC1155 mint",
					slog.String("resourceId", event.ResourceId),
					slog.String("tokenId", event.TokenId.String()),
					slog.Any("error", err))
				hasProcessingError = true
			} else {
				slog.Debug("Successfully handled ERC1155 mint", slog.String("resourceId", event.ResourceId), slog.String("tokenId", event.TokenId.String()))
			}
		}

	case Dvp1155BurnEvent:
		slog.Info("Processing ERC1155 burn event batch")
		eventBatch, err := deserializeDvpMessage[Dvp1155Burn](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize ERC1155 burn events", slog.Any("error", err))
			return
		}
		slog.Debug("Processing ERC1155 burn events", slog.Int("event_count", len(eventBatch.Events)))
		for _, event := range eventBatch.Events {
			if err := z.initiator.HandleERC1155Burn(ctx, event.ChainEventID, event.ResourceId, event.TokenId, event.Value); err != nil {
				slog.Error("Error executing ERC1155 burn",
					slog.String("resourceId", event.ResourceId),
					slog.String("tokenId", event.TokenId.String()),
					slog.Any("error", err))
				hasProcessingError = true
			} else {
				slog.Debug("Successfully handled ERC1155 burn", slog.String("resourceId", event.ResourceId), slog.String("tokenId", event.TokenId.String()))
			}
		}

	case Dvp1155DepositIntoDvpEvent:
		slog.Info("Processing ERC1155 deposit event batch")
		eventBatch, err := deserializeDvpMessage[Dvp1155DepositIntoDvp](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize ERC1155 deposit events", slog.Any("error", err))
			return
		}
		slog.Debug("Processing ERC1155 deposit events", slog.Int("event_count", len(eventBatch.Events)))
		for _, event := range eventBatch.Events {
			if err := z.initiator.HandleERC1155Deposit(ctx, event.ChainEventID, event.ResourceId, event.TokenId, event.Value, event.Data, event.From, event.TxHash, event.TxBlockNumber); err != nil {
				slog.Error("Error executing ERC1155 deposit",
					slog.String("resourceId", event.ResourceId),
					slog.String("tokenId", event.TokenId.String()),
					slog.Any("error", err))
				hasProcessingError = true
			} else {
				slog.Debug("Successfully handled ERC1155 deposit", slog.String("resourceId", event.ResourceId), slog.String("tokenId", event.TokenId.String()))
			}
		}

	case Dvp1155WithdrawFromDvpEvent:
		slog.Info("Processing ERC1155 withdraw event batch")
		eventBatch, err := deserializeDvpMessage[Dvp1155WithdrawFromDvp](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize ERC1155 withdraw events", slog.Any("error", err))
			return
		}
		slog.Debug("Processing ERC1155 withdraw events", slog.Int("event_count", len(eventBatch.Events)))
		for _, event := range eventBatch.Events {
			if err := z.initiator.HandleERC1155Withdraw(ctx, event.ChainEventID, event.ResourceId, event.TokenId, event.Value, event.Owner, event.TxHash, event.TxBlockNumber); err != nil {
				slog.Error("Error executing ERC1155 withdraw",
					slog.String("resourceId", event.ResourceId),
					slog.String("tokenId", event.TokenId.String()),
					slog.Any("error", err))
				hasProcessingError = true
			} else {
				slog.Debug("Successfully handled ERC1155 withdraw", slog.String("resourceId", event.ResourceId), slog.String("tokenId", event.TokenId.String()))
			}
		}

	case Dvp1155SwapForEnygmaEvent:
		slog.Info("Processing ERC1155 swap for Enygma event batch")
		eventBatch, err := deserializeDvpMessage[Dvp1155SwapForEnygma](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize ERC1155 swap for Enygma events", slog.Any("error", err))
			return
		}
		slog.Debug("Processing ERC1155 swap for Enygma events", slog.Int("event_count", len(eventBatch.Events)))
		for _, event := range eventBatch.Events {
			if err := z.initiator.HandleERC1155SwapEnygma(ctx, event.SharedId, event.DestChainId, event.From, event.TokenResourceId, event.TokenValue, event.TokenId, event.EnygmaResourceId, event.EnygmaAmount, event.TxHash, event.TxBlockNumber, event.ValidityTime); err != nil {
				slog.Error("Error executing Swap ERC1155 for Enygma",
					slog.String("sharedId", event.SharedId),
					slog.Any("error", err))
				hasProcessingError = true
			} else {
				slog.Debug("Successfully handled ERC1155 swap for Enygma", slog.String("sharedId", event.SharedId))
			}
		}

	case DvpSwapCancelledEvent:
		slog.Info("Processing swap cancelled event batch")
		eventBatch, err := deserializeDvpMessage[DvpSwapCancelled](msg.V)
		if err != nil {
			slog.Error("Failed to deserialize swap cancelled events", slog.Any("error", err))
			return
		}
		slog.Debug("Processing swap cancelled events", slog.Int("event_count", len(eventBatch.Events)))
		for _, event := range eventBatch.Events {
			tokenInType, err := ercStandardToDvpTokenType(event.TokenInStandard)
			if err != nil {
				slog.Error("Error converting ERC standard to Dvp token type", slog.Any("error", err))
				hasProcessingError = true
				return
			}
			tokenOutType, err := ercStandardToDvpTokenType(event.TokenOutStandard)
			if err != nil {
				slog.Error("Error converting ERC standard to Dvp token type", slog.Any("error", err))
				hasProcessingError = true
				return
			}
			if err := z.initiator.HandleSwapCancellation(ctx, event.SharedId, event.DestChainId, event.TokenInResourceId, event.TokenInAmount, event.TokenInId.String(), tokenInType, event.TokenOutResourceId, event.TokenOutAmount, event.TokenOutId.String(), tokenOutType); err != nil {
				slog.Error("Error executing swap cancellation",
					slog.String("sharedId", event.SharedId),
					slog.Any("error", err))
				hasProcessingError = true
			} else {
				slog.Debug("Successfully handled swap cancelled", slog.String("sharedId", event.SharedId))
			}
		}

	default:
		slog.Warn("Unknown Dvp event type", slog.Int("event_type", int(msg.V.Type)))
		return
	}

	if hasProcessingError {
		slog.Warn("Skipping ack due to processing errors, message will be redelivered",
			slog.Int("event_type", int(msg.V.Type)))
		return
	}

	// Acknowledge message after processing all events
	slog.Debug("Acknowledging dvp message")
	if err := msg.Ack(ctx); err != nil {
		slog.Error("Failed to acknowledge message", slog.Any("error", err))
		return
	}
	slog.Debug("Message acknowledged successfully")
}

func deserializeDvpMessage[T any](serialized DvpSerializedEventBatch) (DvpTypedEventBatch[T], error) {
	var zero DvpTypedEventBatch[T]

	typed := DvpTypedEventBatch[T]{
		BlockNumber: serialized.BlockNumber,
		Type:        serialized.Type,
	}
	err := json.Unmarshal(serialized.SerializedEvents, &typed.Events)
	if err != nil {
		return zero, fmt.Errorf("failed to unmarshal events")
	}
	return typed, nil
}

func ercStandardToDvpTokenType(ercStandard uint8) (types.DvpTokenType, error) {
	standard, err := types.Uint8ToErcStandard(ercStandard)
	if err != nil {
		return 0, fmt.Errorf("failed to convert uint8 %d to ERC standard: %w", ercStandard, err)
	}
	switch standard {
	case types.ENYGMA:
		return types.DvpEnygma, nil
	case types.DVPERC721:
		return types.DvpERC721, nil
	case types.DVPERC1155:
		return types.DvpERC1155, nil
	default:
		return 0, fmt.Errorf("unsupported ERC standard: %d", standard)
	}
}
