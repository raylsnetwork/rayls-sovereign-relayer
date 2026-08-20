package logparser

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	shrservice "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/logrouter"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/service"
)

//go:generate moq --pkg logparser_test -out endpoint_mock_test.go . EndpointMQ CrossChainMQ PrivateHubMQ KeysService
type EndpointMQ interface {
	Next(ctx context.Context) (msgqueue.Message[logrouter.Block], error)
}

type CrossChainMQ interface {
	Push(context.Context, service.CrossChainMessage) error
}

type PrivateHubMQ interface {
	Push(context.Context, shrservice.PrivateHubMessage) error
}

type KeysService interface {
	UpdateRaylsViewKeys(context.Context, *big.Int) error
}

// It parses messageDispatched and BatchMessageDispatched events and splits them
// in two queue -> one for messages for the private hub, and one for cross chain messages
type EndpointLogParser struct {
	privateHubChainID *big.Int

	endpointMQ   EndpointMQ
	crossChainMQ CrossChainMQ
	privateHubMQ PrivateHubMQ
	keysService  KeysService

	endpointEventParser *EndpointV1.EndpointV1
}

func NewEndpointLogParser(
	privateHubChainID *big.Int,
	endpointMQ EndpointMQ,
	crossChainMQ CrossChainMQ,
	privateHubMQ PrivateHubMQ,
	keysService KeysService,
) *EndpointLogParser {
	endpointEventParser := EndpointV1.NewEndpointV1()

	return &EndpointLogParser{
		privateHubChainID: privateHubChainID,

		endpointMQ:   endpointMQ,
		crossChainMQ: crossChainMQ,
		privateHubMQ: privateHubMQ,
		keysService:  keysService,

		endpointEventParser: endpointEventParser,
	}
}

func (p *EndpointLogParser) Fetch( //nolint:gocognit // log parser with event type dispatching
	ctx context.Context,
) error {
	for {
		msg, err := p.endpointMQ.Next(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil //nolint:nilerr // context cancellation is a clean shutdown signal
			}
			slog.Error("Failed to get next message from endpoint MQ", slog.Any("error", err))
			continue
		}

		var hasProcessingError bool
		for _, log := range msg.V.Logs {
			if event, err := p.endpointEventParser.UnpackMessageDispatchedEvent(&log); err == nil {
				// "src-" prefix scopes the ID to the source direction (PN→PNH) so it
				// can't collide with a destination-side (PNH→PN) ID for the same
				// block number — both directions read MessageDispatched events on
				// their respective chains and the block numbers run in parallel.
				// CTS stores both directions in the same cts_transaction table with
				// (correlation_id, message_type) as the PK; without the prefix a
				// collision silently DO-NOTHINGs the second insert and the message
				// is lost.
				messageID := fmt.Sprintf("src-%d-%d-%d", log.BlockNumber, log.TxIndex, log.Index)

				if event.ToChainId.Cmp(p.privateHubChainID) == 0 {
					if err := p.privateHubMQ.Push(ctx, eventToPrivateHubMessage(messageID, event)); err != nil {
						slog.Error("Failed to push message to private hub MQ",
							slog.String("message_id", messageID),
							slog.Any("error", err))
						hasProcessingError = true
					}
				} else {
					if err := p.crossChainMQ.Push(ctx, eventToCrossChainMessage(messageID, event)); err != nil {
						slog.Error("Failed to push message to cross chain MQ",
							slog.String("message_id", messageID),
							slog.Any("error", err))
						hasProcessingError = true
					}
				}
			}

			if event, err := p.endpointEventParser.UnpackMessageBatchDispatchedEvent(&log); err == nil {
				for i, msg := range event.Messages {
					messageID := fmt.Sprintf("src-%d-%d-%d(%d)", log.BlockNumber, log.TxIndex, log.Index, i)

					if msg.ToChainId.Cmp(p.privateHubChainID) == 0 {
						if err := p.privateHubMQ.Push(
							ctx,
							batchEventToPrivateHubMessage(messageID, event.From, msg),
						); err != nil {
							slog.Error("Failed to push batch message to private hub MQ",
								slog.String("message_id", messageID),
								slog.Any("error", err))
							hasProcessingError = true
						}
					} else {
						if err := p.crossChainMQ.Push(
							ctx,
							batchEventToCrossChainMessage(messageID, msg, event),
						); err != nil {
							slog.Error("Failed to push batch message to cross chain MQ",
								slog.String("message_id", messageID),
								slog.Any("error", err))
							hasProcessingError = true
						}
					}
				}
			}

			if event, err := p.endpointEventParser.UnpackUpdateRaylsViewKeysRequestEvent(&log); err == nil {
				slog.Info("Updating View keys", slog.Any("block_number", event.BlockNumber))
				err = p.keysService.UpdateRaylsViewKeys(ctx, event.BlockNumber)
				if err == nil {
					slog.Info("Successfully updated View keys", slog.Any("block_number", event.BlockNumber))
				} else {
					slog.Error(
						"Failed to update DH keys",
						slog.Any("block_number", event.BlockNumber),
						slog.Any("error", err),
					)
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

func eventToPrivateHubMessage(
	messageID string,
	event *EndpointV1.EndpointV1MessageDispatched,
) shrservice.PrivateHubMessage {
	return shrservice.PrivateHubMessage{
		ID:        messageID,
		MessageID: event.MessageId,
		From:      event.From,
		To:        event.To,
		Data:      event.Data,
	}
}

func eventToCrossChainMessage(
	messageID string,
	event *EndpointV1.EndpointV1MessageDispatched,
) service.CrossChainMessage {
	return service.CrossChainMessage{
		ID:        messageID,
		MessageID: event.MessageId,
		From:      event.From,
		ToChainID: event.ToChainId,
		To:        event.To,
		Data:      event.Data,

		BlockHash: event.Raw.BlockHash,
		TxHash:    event.Raw.TxHash,

		BlockNumber: event.Raw.BlockNumber,
		TxIdx:       event.Raw.TxIndex,
		LogIdx:      event.Raw.Index,

		TokenAddress: event.Data.MessageMetadata.TransferMetadata.TokenAddress,
	}
}

func batchEventToPrivateHubMessage(
	messageID string,
	from common.Address,
	msg EndpointV1.BatchMessage,
) shrservice.PrivateHubMessage {
	return shrservice.PrivateHubMessage{
		ID:        messageID,
		MessageID: msg.MessageId,
		From:      from,
		To:        msg.To,
		Data:      msg.Data,
	}
}

func batchEventToCrossChainMessage(
	messageID string,
	msg EndpointV1.BatchMessage,
	event *EndpointV1.EndpointV1MessageBatchDispatched,
) service.CrossChainMessage {
	return service.CrossChainMessage{
		ID:        messageID,
		MessageID: msg.MessageId,
		From:      event.From,
		ToChainID: msg.ToChainId,
		To:        msg.To,
		Data:      msg.Data,

		BlockHash: event.Raw.BlockHash,
		TxHash:    event.Raw.TxHash,

		BlockNumber: event.Raw.BlockNumber,
		TxIdx:       event.Raw.TxIndex,
		LogIdx:      event.Raw.Index,

		TokenAddress: msg.Data.MessageMetadata.TransferMetadata.TokenAddress,
	}
}
