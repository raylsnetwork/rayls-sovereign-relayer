package logparser

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/backoff"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/logrouter"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/shared/service"
)

//go:generate moq --pkg logparser_test -out endpoint_mock_test.go . EndpointMQ PrivateHubMQ EndpointBackoff
type EndpointMQ interface {
	Next(ctx context.Context) (msgqueue.Message[logrouter.Block], error)
}

type PrivateHubMQ interface {
	PushBatch(context.Context, []service.PrivateHubMessage) error
}

type EndpointBackoff interface {
	Do(ctx context.Context, maxAttempts int, fn func() error) error
}

type EndpointLogParser struct {
	nodeChainID *big.Int

	endpointMQ   EndpointMQ
	privateHubMQ PrivateHubMQ

	endpointEventParser *EndpointV1.EndpointV1

	backoff EndpointBackoff
}

func NewEndpointLogParser(
	nodeChainID *big.Int,
	endpointMQ EndpointMQ,
	privateHubMQ PrivateHubMQ,
) *EndpointLogParser {
	b, _ := backoff.NewExponential(time.Second, 2, time.Minute)
	return NewEndpointLogParserWithBackoff(nodeChainID, endpointMQ, privateHubMQ, b)
}

func NewEndpointLogParserWithBackoff(
	nodeChainID *big.Int,
	endpointMQ EndpointMQ,
	privateHubMQ PrivateHubMQ,
	backoff EndpointBackoff,
) *EndpointLogParser {
	return &EndpointLogParser{
		nodeChainID: nodeChainID,

		endpointMQ:   endpointMQ,
		privateHubMQ: privateHubMQ,

		endpointEventParser: EndpointV1.NewEndpointV1(),

		backoff: backoff,
	}
}

func (p *EndpointLogParser) Run(ctx context.Context) error {
	slog.Info("EndpointLogParser started")
	for {
		slog.Debug("Fetching next endpoint block from queue")
		block, err := p.endpointMQ.Next(ctx)
		if err != nil {
			if ctx.Err() != nil {
				slog.Info("EndpointLogParser shutting down")
				return nil //nolint:nilerr // context cancellation is a clean shutdown signal
			}
			slog.Error("Failed to get next message from endpoint MQ", slog.Any("error", err))
			continue
		}

		slog.Info(
			"Processing endpoint logs",
			slog.Uint64("block_number", block.V.Number),
			slog.Int("log_count", len(block.V.Logs)),
		)

		var msgs []service.PrivateHubMessage
		for _, log := range block.V.Logs {
			if event, parseErr := p.endpointEventParser.UnpackMessageDispatchedEvent(&log); parseErr == nil {
				// CHAIN_ID_ALL_PARTICIPANTS (0) means broadcast to all participants
				isAllParticipants := event.ToChainId.Cmp(big.NewInt(0)) == 0
				isTargetedToThisNode := event.ToChainId.Cmp(p.nodeChainID) == 0
				if !isAllParticipants && !isTargetedToThisNode {
					continue
				}

				// "dst-" prefix scopes the ID to the destination direction (PNH→PN)
				// so it cannot collide with a source-side ID at the same block
				// number; see the matching note on source/logparser/endpoint.go.
				messageID := fmt.Sprintf("dst-%d-%d-%d", log.BlockNumber, log.TxIndex, log.Index)

				msgs = append(msgs, service.PrivateHubMessage{
					ID: messageID,

					MessageID: event.MessageId,
					From:      event.From,
					To:        event.To,
					Data:      event.Data,
				})
			}
		}

		if len(msgs) > 0 {
			err = p.backoff.Do(ctx, 10, func() error {
				return p.privateHubMQ.PushBatch(ctx, msgs)
			})
			if err != nil {
				slog.Error("Failed to push message batch to private hub MQ after retries, block will be redelivered",
					slog.Uint64("block_number", block.V.Number),
					slog.Int("message_count", len(msgs)),
					slog.Any("error", err))
				continue
			}
			slog.Info("Processed endpoint messages",
				slog.Uint64("block_number", block.V.Number),
				slog.Int("count", len(msgs)))
		}

		slog.Debug("Acknowledging endpoint block", slog.Uint64("block_number", block.V.Number))
		if err := p.backoff.Do(ctx, 10, func() error {
			return block.Ack(ctx)
		}); err != nil {
			slog.Error("Failed to acknowledge block after retries",
				slog.Uint64("blockNumber", block.V.Number),
				slog.Any("error", err),
			)
		}
	}
}
