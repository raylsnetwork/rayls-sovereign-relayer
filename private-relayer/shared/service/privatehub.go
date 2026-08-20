package service

import (
	"context"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/batcher"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/msgqueue"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

//go:generate moq --pkg service_test -out privatehub_mock_test.go . TransactionGenerator PrivateHubConsumer Batcher
type TransactionGenerator interface {
	Generate(
		fromChainID *big.Int,
		fromAddress, toAddress common.Address,
		message EndpointV1.RaylsMessage,
		id common.Hash,
	) ([]byte, error)
}

type PrivateHubConsumer interface {
	Fetch(ctx context.Context, count int) ([]msgqueue.Message[PrivateHubMessage], error)
}

// Batcher is the fire-and-forget publisher the PrivateHubService hands
// its generated messages to. A concrete *batcher.Batcher wired against
// the cross-chain NATS subject satisfies it.
type Batcher interface {
	Send(ctx context.Context, msgs []batcher.Message) error
}

type PrivateHubService struct {
	ticker              *time.Ticker
	srcChainID          *big.Int
	destEndpointAddress common.Address

	msgCons PrivateHubConsumer
	txGen   TransactionGenerator
	batcher Batcher
}

func NewPrivateHubService(
	tickerPeriod time.Duration,
	srcChainID *big.Int,
	destEndpointAddress common.Address,
	msgCons PrivateHubConsumer,
	txGen TransactionGenerator,
	b Batcher,
) *PrivateHubService {
	return &PrivateHubService{
		ticker:              time.NewTicker(tickerPeriod),
		srcChainID:          srcChainID,
		destEndpointAddress: destEndpointAddress,

		msgCons: msgCons,
		txGen:   txGen,
		batcher: b,
	}
}

func (s *PrivateHubService) Run(ctx context.Context) error {
	slog.Info("PrivateHubService started", slog.String("endpoint", s.destEndpointAddress.Hex()))
	initialRun := make(chan struct{}, 1)
	initialRun <- struct{}{}

	for {
		select {
		case <-ctx.Done():
			slog.Info("PrivateHubService shutting down")
			return nil
		case <-s.ticker.C:
		case <-initialRun:
		}

		slog.Debug("PrivateHubService tick")

		slog.Debug("Fetching messages", slog.Int("count", 100))
		// TODO: move fetch batch size as config parameter for service
		msgs, err := s.msgCons.Fetch(ctx, 100)
		if err != nil {
			slog.Error("Failed to read messages from consumer", slog.Any("error", err))
			continue
		}

		if len(msgs) == 0 {
			slog.Debug("No messages available")
			continue
		}

		slog.Debug("Fetched messages", slog.Int("count", len(msgs)))

		var batchMsgs []batcher.Message
		for _, msg := range msgs {
			slog.Debug("Generating calldata", slog.String("msg_id", msg.V.ID))
			calldata, genErr := s.txGen.Generate(
				s.srcChainID,
				msg.V.From,
				msg.V.To,
				msg.V.Data,
				msg.V.MessageID,
			)
			if genErr != nil {
				slog.Warn("Failed to generate transaction calldata, skipping message",
					slog.Any("error", genErr),
					slog.String("msg_id", msg.V.ID),
				)
				continue
			}

			slog.Debug("Generated calldata", slog.String("msg_id", msg.V.ID))

			batchMsgs = append(batchMsgs, batcher.Message{
				ID:       msg.V.ID,
				Address:  s.destEndpointAddress,
				Calldata: calldata,
			})
		}
		slog.Debug("Publishing batch", slog.Int("count", len(batchMsgs)))
		err = s.batcher.Send(ctx, batchMsgs)
		if err == nil {
			slog.Info("Batch published", slog.Int("count", len(batchMsgs)))
			slog.Debug("Acknowledging messages", slog.Int("count", len(msgs)))
			for _, msg := range msgs {
				if ackErr := msg.Ack(ctx); ackErr != nil {
					slog.Warn("Failed to acknowledge message",
						slog.Any("error", ackErr), slog.String("msg_id", msg.V.ID))
				}
			}
		} else {
			slog.Warn(
				"Failed to publish batch, skipping ack",
				slog.Any("error", err),
				slog.Int("count", len(batchMsgs)),
			)
		}

	}
}

func (s *PrivateHubService) HandleResults(ctx context.Context, results []types.TxResult) error {
	if len(results) == 0 {
		return nil
	}

	var (
		successfulCount int
		failedCount     int
	)
	for _, res := range results {
		switch res.Kind {
		case types.TxResultSuccess:
			successfulCount++
		default:
			slog.Warn("Failed tx", slog.String("hash", res.TxHash.Hex()), slog.String("revert data", string(res.RevertData)), slog.String("revert reason", res.ErrorReason))
			failedCount++
		}
	}
	slog.Info("Received private hub receipts", slog.Any("successful", successfulCount), slog.Any("failed", failedCount))
	return nil
}
