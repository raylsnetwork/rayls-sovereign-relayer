package logparser

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/backoff"
	am "github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/AuditManagerV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/keys"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/logrouter"
	"google.golang.org/grpc"
)

const maxRetries = 3

//go:generate moq --pkg logparser_test -out auditmanager_mock_test.go . AuditManagerMQ KOSKeyAgreementClient AuditManagerBackoff

type AuditManagerMQ interface {
	Next(ctx context.Context) (msgqueue.Message[logrouter.Block], error)
}

type KOSKeyAgreementClient interface {
	CompleteKeyAgreement(ctx context.Context, in *keys.CompleteKeyAgreementRequest, opts ...grpc.CallOption) (*keys.CompleteKeyAgreementResponse, error)
}

type AuditManagerBackoff interface {
	Do(ctx context.Context, maxAttempts int, fn func() error) error
}

type AuditManagerParser struct {
	localChainID         *big.Int
	auditManagerMQ       AuditManagerMQ
	kosClient            KOSKeyAgreementClient
	auditManagerFilterer *am.AuditManagerV1
	backoff              AuditManagerBackoff
}

func NewAuditManagerParser(
	localChainID *big.Int,
	auditManagerMQ AuditManagerMQ,
	kosClient KOSKeyAgreementClient,
) *AuditManagerParser {
	b, _ := backoff.NewExponential(time.Second, 2, time.Minute)
	return NewAuditManagerParserWithBackoff(localChainID, auditManagerMQ, kosClient, b)
}

func NewAuditManagerParserWithBackoff(
	localChainID *big.Int,
	auditManagerMQ AuditManagerMQ,
	kosClient KOSKeyAgreementClient,
	backoff AuditManagerBackoff,
) *AuditManagerParser {
	filterer := am.NewAuditManagerV1()

	return &AuditManagerParser{
		localChainID:         localChainID,
		auditManagerMQ:       auditManagerMQ,
		kosClient:            kosClient,
		auditManagerFilterer: filterer,
		backoff:              backoff,
	}
}

func (p *AuditManagerParser) Run(ctx context.Context) error {
	slog.Info("AuditManagerParser started")
	for {
		slog.Debug("Fetching next audit manager block from queue")
		block, err := p.auditManagerMQ.Next(ctx)
		if err != nil {
			if ctx.Err() != nil {
				slog.Info("AuditManagerParser shutting down")
				return nil //nolint:nilerr // context cancellation is a clean shutdown signal
			}
			slog.Error("Failed to get next message from audit manager MQ", slog.Any("error", err))
			continue
		}

		slog.Debug("Processing audit manager logs",
			slog.Uint64("block_number", block.V.Number),
			slog.Int("log_count", len(block.V.Logs)))

		for _, log := range block.V.Logs {
			if event, err := p.auditManagerFilterer.UnpackKeyAgreementInitiatedEvent(&log); err == nil {
				if err := p.handleKeyAgreementInitiated(ctx, event, block.V.Number); err != nil {
					slog.Error("Failed to handle KeyAgreementInitiated event", slog.Any("error", err))
				}
			}
		}

		slog.Debug("Acknowledging audit manager block", slog.Uint64("block_number", block.V.Number))
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

func (p *AuditManagerParser) handleKeyAgreementInitiated(
	ctx context.Context,
	event *am.AuditManagerV1KeyAgreementInitiated,
	blockNumber uint64,
) error {
	slog.Debug("Processing KeyAgreementInitiated event")

	if event.FromChainId.Cmp(p.localChainID) == 0 {
		slog.Debug("KeyAgreementInitiated event from ourselves. Skipping...")
		return nil
	}

	if event.ToChainId.Cmp(p.localChainID) != 0 {
		slog.Debug("KeyAgreementInitiated event not destined for us. Skipping...")
		return nil
	}

	slog.Info("Completing key agreement",
		slog.String("fromChainId", event.FromChainId.String()),
		slog.String("blockNumber", event.BlockNumber.String()))

	err := p.backoff.Do(ctx, maxRetries, func() error {
		_, err := p.kosClient.CompleteKeyAgreement(
			ctx,
			&keys.CompleteKeyAgreementRequest{
				ChainID:     event.FromChainId.Uint64(),
				Ciphertext:  event.Ciphertext,
				BlockNumber: event.BlockNumber.Uint64(),
			},
		)
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to complete key agreement with initiator %s: %w",
			event.FromChainId.String(), err)
	}

	slog.Info("Successfully completed key agreement",
		slog.String("fromChainId", event.FromChainId.String()),
		slog.Uint64("blockNumber", blockNumber),
	)

	return nil
}
