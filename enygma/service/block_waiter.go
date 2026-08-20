package service

import (
	"context"
	"log/slog"
	"time"

	telemetry "github.com/raylsnetwork/rayls-sovereign-relayer/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	blockWaiterInterval = 1 * time.Second
)

type blockWaiterTracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}
type blockWaiterEthereumClient interface {
	BlockNumber(context.Context) (uint64, error)
}

type BlockWaiterService struct {
	tracer    blockWaiterTracer
	ethClient blockWaiterEthereumClient
}

func NewBlockWaiterService(tracer blockWaiterTracer, ethClient blockWaiterEthereumClient) *BlockWaiterService {
	return &BlockWaiterService{
		tracer:    tracer,
		ethClient: ethClient,
	}
}

// WaitForBlock blocks until the blockchain reaches or exceeds targetBlockNumber
// Returns the actual block number reached or an error if timeout occurs
func (bws *BlockWaiterService) WaitForBlock(ctx context.Context, targetBlockNumber uint64) (uint64, error) {
	ctx, span := bws.tracer.Start(ctx, telemetry.SPAN_WAIT_FOR_NEXT_BLOCK)
	defer span.End()

	ticker := time.NewTicker(blockWaiterInterval)

	// Wait until the latest block number is greater than the target block number.
	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-ticker.C:
		}

		timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 30*time.Second)
		latestBlockNumber, err := bws.ethClient.BlockNumber(timeoutCtx)
		timeoutCancel() // Clean up immediately to prevent leak

		if err != nil {
			slog.Warn("Transient error getting latest block number, will retry",
				slog.String("error", err.Error()),
				slog.Uint64("targetBlockNumber", targetBlockNumber),
			)
			continue
		}

		if latestBlockNumber > targetBlockNumber {
			return latestBlockNumber, nil
		}

		slog.Debug(
			"Waiting for the latest block number to be greater than the target block number",
			slog.Uint64("targetBlockNumber", targetBlockNumber),
			slog.Uint64("latestBlockNumber", latestBlockNumber),
		)
	}
}
