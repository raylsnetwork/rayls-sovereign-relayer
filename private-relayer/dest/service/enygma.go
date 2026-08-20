package service

import (
	"context"
	"log/slog"
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/backoff"
	"github.com/raylsnetwork/rayls-sovereign-relayer/msgqueue"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

//go:generate moq --pkg service_test -out enygma_mock_test.go . EnygmaDestMQ EnygmaReceiver EnygmaCheckpointRepository Backoff

type Backoff interface {
	Do(ctx context.Context, maxAttempts int, fn func() error) error
}

type EnygmaDestMQ interface {
	Next(ctx context.Context) (msgqueue.Message[EnygmaDestMessage], error)
}

type EnygmaReceiver interface {
	HandleEnygmaCrossTransfer(ctx context.Context, batch *types.EnygmaTransferBatch) error
}

type EnygmaCheckpointRepository interface {
	GetLatestCheckpointByFilters(
		ctx context.Context,
		resourceId string,
		status *types.EnygmaCheckpointStatus,
		finalizedBlockNumber *big.Int,
		pendingBlockNumber *big.Int,
	) (*types.EnygmaCheckpoint, error)
	CreateEnygmaCheckpoint(ctx context.Context, checkpoint types.EnygmaCheckpoint) error
}

// defaultHandleMessageTimeout (declared in dvp.go, shared across this package) bounds a single handler
// invocation so one stuck handler — e.g. an on-chain call whose TX never mines, leaving bind.WaitMined
// polling forever on the long-lived run ctx — cannot head-of-line stall the single-threaded loop.

type EnygmaOrchestrator struct {
	enygmaMQ EnygmaDestMQ

	receiver             EnygmaReceiver
	checkpointRepository EnygmaCheckpointRepository

	backoff Backoff

	handleMessageTimeout time.Duration
}

func NewEnygmaOrchestrator(
	enygmaMQ EnygmaDestMQ,
	receiver EnygmaReceiver,
	checkpointRepository EnygmaCheckpointRepository,
) *EnygmaOrchestrator {
	backoff, _ := backoff.NewExponential(time.Second, 2, time.Minute)
	return NewEnygmaOrchestratorWithBackoff(enygmaMQ, receiver, checkpointRepository, backoff)
}

func NewEnygmaOrchestratorWithBackoff(
	enygmaMQ EnygmaDestMQ,
	receiver EnygmaReceiver,
	checkpointRepository EnygmaCheckpointRepository,
	backoff Backoff,
) *EnygmaOrchestrator {
	return NewEnygmaOrchestratorWithBackoffAndTimeout(enygmaMQ, receiver, checkpointRepository, backoff, defaultHandleMessageTimeout)
}

// NewEnygmaOrchestratorWithBackoffAndTimeout is like NewEnygmaOrchestratorWithBackoff but lets callers
// (and tests) set the per-message handler timeout. A non-positive timeout falls back to the default.
func NewEnygmaOrchestratorWithBackoffAndTimeout(
	enygmaMQ EnygmaDestMQ,
	receiver EnygmaReceiver,
	checkpointRepository EnygmaCheckpointRepository,
	backoff Backoff,
	handleMessageTimeout time.Duration,
) *EnygmaOrchestrator {
	if handleMessageTimeout <= 0 {
		handleMessageTimeout = defaultHandleMessageTimeout
	}
	return &EnygmaOrchestrator{
		enygmaMQ:             enygmaMQ,
		receiver:             receiver,
		checkpointRepository: checkpointRepository,
		backoff:              backoff,
		handleMessageTimeout: handleMessageTimeout,
	}
}

func (o *EnygmaOrchestrator) Run(ctx context.Context) error {
	slog.Info("EnygmaOrchestrator started")
	for {
		slog.Debug("Fetching next enygma message from queue")
		msg, err := o.enygmaMQ.Next(ctx)
		if err != nil {
			if ctx.Err() != nil {
				slog.Info("EnygmaOrchestrator shutting down")
				return nil //nolint:nilerr // context cancellation is a clean shutdown signal
			}
			slog.Error("Failed to get next message from enygma MQ", slog.Any("error", err))
			continue
		}

		slog.Debug("Processing enygma dest message",
			slog.String("messageId", msg.V.ID),
			slog.Int("type", int(msg.V.Type)),
		)

		shouldAck := o.dispatch(ctx, msg.V)

		if shouldAck {
			slog.Debug("Acknowledging enygma message", slog.String("messageId", msg.V.ID))
			if err := o.backoff.Do(ctx, 10, func() error {
				return msg.Ack(ctx)
			}); err != nil {
				slog.Error("Failed to acknowledge message after retries",
					slog.String("messageId", msg.V.ID),
					slog.Any("error", err),
				)
			} else {
				slog.Debug("Message acknowledged successfully", slog.String("messageId", msg.V.ID))
			}
		}
	}
}

// dispatch runs the per-type handler under a bounded context so one stuck handler cannot head-of-line
// stall the single-threaded loop. The handler ctx is cancelled when the parent ctx is cancelled (clean
// shutdown) OR when handleMessageTimeout elapses. A handler that overruns observes a cancelled ctx — its
// in-flight chain wait (bind.WaitMined) returns "context deadline exceeded" — and returns shouldAck=false,
// so the message is redelivered (the same path a handler error already takes) rather than the loop
// stalling. Redelivery here is the pre-existing on-error behaviour, not a new semantic.
func (o *EnygmaOrchestrator) dispatch(ctx context.Context, v EnygmaDestMessage) bool {
	hctx, cancel := context.WithTimeout(ctx, o.handleMessageTimeout)
	defer cancel()

	switch v.Type {
	case EnygmaTransferBatchMessage:
		return o.handleEnygmaTransferBatch(hctx, v)
	case EnygmaFinalizedBalanceMessage:
		return o.handleEnygmaFinalizedBalance(hctx, v)
	default:
		slog.Warn("Unknown enygma message type", slog.Int("type", int(v.Type)))
		return true
	}
}

func (o *EnygmaOrchestrator) handleEnygmaTransferBatch(ctx context.Context, msg EnygmaDestMessage) bool {
	if msg.TransferBatch == nil {
		slog.Error("EnygmaTransferBatch message has nil batch", slog.String("messageId", msg.ID))
		return true
	}

	err := o.receiver.HandleEnygmaCrossTransfer(ctx, msg.TransferBatch)
	if err != nil {
		slog.Error("Failed to handle EnygmaCrossTransfer",
			slog.Any("error", err),
			slog.String("messageId", msg.ID),
			slog.String("resourceId", msg.TransferBatch.ResourceId),
		)
		return false
	}

	slog.Info("Successfully handled EnygmaTransferBatch",
		slog.String("messageId", msg.ID),
		slog.String("resourceId", msg.TransferBatch.ResourceId),
	)

	return true
}

func (o *EnygmaOrchestrator) handleEnygmaFinalizedBalance(ctx context.Context, msg EnygmaDestMessage) bool {
	if msg.FinalizedBalance == nil {
		slog.Error("EnygmaFinalizedBalance message has nil balance", slog.String("messageId", msg.ID))
		return true
	}

	balance := msg.FinalizedBalance

	existingCheckpoint, err := o.checkpointRepository.GetLatestCheckpointByFilters(
		ctx,
		balance.ResourceId,
		nil,
		balance.FinalizedBlockNumber,
		balance.PendingBlockNumber,
	)
	if err != nil {
		slog.Error("Failed to get existing checkpoint",
			slog.Any("error", err),
			slog.String("resourceId", balance.ResourceId),
		)
		return false
	}

	if existingCheckpoint != nil {
		slog.Debug("Checkpoint already exists, skipping",
			slog.String("resourceId", balance.ResourceId),
			slog.String("checkpointId", existingCheckpoint.ID),
		)
		return true
	}

	if balance.BalanceX.Cmp(big.NewInt(0)) == 0 && balance.BalanceY.Cmp(big.NewInt(1)) == 0 {
		slog.Debug("Skipping zero balance checkpoint",
			slog.String("resourceId", balance.ResourceId),
		)
		return true
	}

	// Check if balance changed since last checkpoint (like the original ProcessEnygmaFinalizedBalances)
	latestCheckpoint, err := o.checkpointRepository.GetLatestCheckpointByFilters(ctx, balance.ResourceId, nil, nil, nil)
	if err != nil {
		slog.Error("Failed to get latest checkpoint for resourceId",
			slog.String("resourceId", balance.ResourceId),
			slog.Any("error", err),
		)
		return false
	}

	// The public balance of our PL was not changed since the last checkpoint.
	// This means we were not involved in any transfer since then.
	// So, we don't need to create a new checkpoint.
	if latestCheckpoint != nil &&
		latestCheckpoint.FinalizedPublicBalanceX.Cmp(balance.BalanceX) == 0 &&
		latestCheckpoint.FinalizedPublicBalanceY.Cmp(balance.BalanceY) == 0 {
		slog.Info("Public balance unchanged since last checkpoint, skipping",
			slog.String("resourceId", balance.ResourceId),
			slog.String("blockNumber", balance.FinalizedBlockNumber.String()),
		)
		return true
	}

	checkpoint := types.EnygmaCheckpoint{
		ResourceId:              balance.ResourceId,
		FinalizedPublicBalanceX: balance.BalanceX,
		FinalizedPublicBalanceY: balance.BalanceY,
		FinalizedBlockNumber:    balance.FinalizedBlockNumber,
		PendingBlockNumber:      balance.PendingBlockNumber,
		Status:                  types.EnygmaCheckpointStatusTentative,
	}

	err = o.checkpointRepository.CreateEnygmaCheckpoint(ctx, checkpoint)
	if err != nil {
		slog.Error("Failed to create checkpoint",
			slog.Any("error", err),
			slog.String("resourceId", balance.ResourceId),
		)
		return false
	}

	slog.Info("Successfully created EnygmaCheckpoint",
		slog.String("messageId", msg.ID),
		slog.String("resourceId", balance.ResourceId),
	)

	return true
}
