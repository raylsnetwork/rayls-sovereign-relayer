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

//go:generate moq --pkg service_test -out dvp_mock_test.go . DvpDestMQ DvpReceiver DvpBackoff

type DvpBackoff interface {
	Do(ctx context.Context, maxAttempts int, fn func() error) error
}

type DvpDestMQ interface {
	Next(ctx context.Context) (msgqueue.Message[DvpDestMessage], error)
}

type DvpReceiver interface {
	HandleCommitments(ctx context.Context, data *types.DvpCommitmentsData) error
	HandleNullifiers(ctx context.Context, tokenAddress string, nullifiers []*big.Int) error
	HandleSwapInitiated(ctx context.Context, blockNum uint64, data *types.DvpSwapInitiatedData) error
	HandleSwapCompleted(ctx context.Context, sharedID string) error
	HandleSwapRevert(ctx context.Context, sharedId string, status types.DvpSwapStatus) error
}

// defaultHandleMessageTimeout bounds a single handler invocation. A handler that does not return within
// this window (e.g. an on-chain notify whose TX never mines, so bind.WaitMined polls forever on the
// long-lived run ctx) is treated as a transient failure: its ctx is cancelled, the chain wait returns
// "context deadline exceeded", the handler returns shouldAck=false (the same no-ack/redeliver path a
// handler error already takes), and the single-threaded loop continues with the next message instead of
// stalling indefinitely on one stuck handler. Generous so normal mining latency never trips it.
const defaultHandleMessageTimeout = 2 * time.Minute

type DvpOrchestrator struct {
	dvpMQ DvpDestMQ

	receiver DvpReceiver

	backoff DvpBackoff

	handleMessageTimeout time.Duration
}

func NewDvpOrchestrator(
	dvpMQ DvpDestMQ,
	receiver DvpReceiver,
) *DvpOrchestrator {
	backoff, _ := backoff.NewExponential(time.Second, 2, time.Minute)
	return NewDvpOrchestratorWithBackoff(dvpMQ, receiver, backoff)
}

func NewDvpOrchestratorWithBackoff(
	dvpMQ DvpDestMQ,
	receiver DvpReceiver,
	backoff DvpBackoff,
) *DvpOrchestrator {
	return NewDvpOrchestratorWithBackoffAndTimeout(dvpMQ, receiver, backoff, defaultHandleMessageTimeout)
}

// NewDvpOrchestratorWithBackoffAndTimeout is like NewDvpOrchestratorWithBackoff but lets callers (and
// tests) set the per-message handler timeout. A non-positive timeout falls back to the default.
func NewDvpOrchestratorWithBackoffAndTimeout(
	dvpMQ DvpDestMQ,
	receiver DvpReceiver,
	backoff DvpBackoff,
	handleMessageTimeout time.Duration,
) *DvpOrchestrator {
	if handleMessageTimeout <= 0 {
		handleMessageTimeout = defaultHandleMessageTimeout
	}
	return &DvpOrchestrator{
		dvpMQ:                dvpMQ,
		receiver:             receiver,
		backoff:              backoff,
		handleMessageTimeout: handleMessageTimeout,
	}
}

func (o *DvpOrchestrator) Run(ctx context.Context) error {
	slog.Info("DvpOrchestrator started")
	for {
		slog.Debug("Fetching next dvp message from queue")
		msg, err := o.dvpMQ.Next(ctx)
		if err != nil {
			if ctx.Err() != nil {
				slog.Info("DvpOrchestrator shutting down")
				return nil //nolint:nilerr // context cancellation is a clean shutdown signal
			}
			slog.Error("Failed to get next message from dvp MQ", slog.Any("error", err))
			continue
		}

		slog.Debug("Processing dvp dest message",
			slog.String("messageId", msg.V.ID),
			slog.Int("type", int(msg.V.Type)),
		)

		shouldAck := o.dispatch(ctx, msg.V)

		if shouldAck {
			slog.Debug("Acknowledging dvp message", slog.String("messageId", msg.V.ID))
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
// shutdown) OR when handleMessageTimeout elapses. A handler that overruns the deadline observes a
// cancelled ctx — its in-flight chain wait (bind.WaitMined) returns "context deadline exceeded" — and
// returns shouldAck=false, so the message is redelivered (the same path a handler error already takes)
// rather than the loop stalling. Redelivery here is the pre-existing on-error behaviour, not a new
// semantic: it does not change executor send/retry idempotency.
func (o *DvpOrchestrator) dispatch(ctx context.Context, v DvpDestMessage) bool {
	hctx, cancel := context.WithTimeout(ctx, o.handleMessageTimeout)
	defer cancel()

	switch v.Type {
	case DvpCommitmentsMessage:
		return o.handleCommitments(hctx, v)
	case DvpNullifierMessage:
		return o.handleNullifiers(hctx, v)
	case DvpSwapInitiatedMessage:
		return o.handleSwapInitiated(hctx, v)
	case DvpSwapCompletedMessage:
		return o.handleSwapCompleted(hctx, v)
	case DvpSwapCancelledMessage:
		return o.handleSwapCancelled(hctx, v)
	case DvpSwapTimedOutMessage:
		return o.handleSwapTimedOut(hctx, v)
	default:
		slog.Warn("Unknown dvp message type", slog.Int("type", int(v.Type)))
		return true
	}
}

func (o *DvpOrchestrator) handleCommitments(ctx context.Context, msg DvpDestMessage) bool {
	if msg.Commitments == nil {
		slog.Error("DvpCommitments message has nil commitments data", slog.String("messageId", msg.ID))
		return true
	}

	err := o.receiver.HandleCommitments(ctx, msg.Commitments)
	if err != nil {
		slog.Error("Failed to handle DvpCommitments",
			slog.Any("error", err),
			slog.String("messageId", msg.ID),
		)
		return false
	}

	slog.Info("Successfully handled DvpCommitments",
		slog.String("messageId", msg.ID),
	)

	return true
}

func (o *DvpOrchestrator) handleNullifiers(ctx context.Context, msg DvpDestMessage) bool {
	if msg.Nullifiers == nil {
		slog.Error("DvpNullifiers message has nil nullifier data", slog.String("messageId", msg.ID))
		return true
	}

	err := o.receiver.HandleNullifiers(ctx, msg.Nullifiers.TokenAddress, msg.Nullifiers.Nullifiers)
	if err != nil {
		slog.Error("Failed to handle DvpNullifiers",
			slog.Any("error", err),
			slog.String("messageId", msg.ID),
		)
		return false
	}

	slog.Info("Successfully handled DvpNullifiers",
		slog.String("messageId", msg.ID),
	)

	return true
}

func (o *DvpOrchestrator) handleSwapCompleted(ctx context.Context, msg DvpDestMessage) bool {
	if msg.SharedID == "" {
		slog.Error("DvpSwapCompleted message has empty sharedId", slog.String("messageId", msg.ID))
		return true
	}

	if err := o.receiver.HandleSwapCompleted(ctx, msg.SharedID); err != nil {
		slog.Error("Failed to handle DvpSwapCompleted",
			slog.Any("error", err),
			slog.String("messageId", msg.ID),
			slog.String("sharedId", msg.SharedID),
		)
		return false
	}

	slog.Info("Successfully handled DvpSwapCompleted",
		slog.String("messageId", msg.ID),
		slog.String("sharedId", msg.SharedID),
	)
	return true
}

func (o *DvpOrchestrator) handleSwapInitiated(ctx context.Context, msg DvpDestMessage) bool {
	if msg.SwapInitiated == nil || msg.SwapInitiated.Message == nil {
		slog.Error("DvpSwapInitiated message has nil data", slog.String("messageId", msg.ID))
		return true
	}

	err := o.receiver.HandleSwapInitiated(ctx, msg.BlockNumber, msg.SwapInitiated)
	if err != nil {
		slog.Error("Failed to handle DvpSwapInitiated",
			slog.Any("error", err),
			slog.String("messageId", msg.ID),
			slog.String("sharedId", msg.SwapInitiated.Message.SharedId),
		)
		return false
	}

	slog.Info("Successfully handled DvpSwapInitiated",
		slog.String("messageId", msg.ID),
		slog.String("sharedId", msg.SwapInitiated.Message.SharedId),
	)

	return true
}

func (o *DvpOrchestrator) handleSwapCancelled(ctx context.Context, msg DvpDestMessage) bool {
	if msg.SharedID == "" {
		slog.Error("DvpSwapCancelled message has empty sharedId", slog.String("messageId", msg.ID))
		return true
	}

	if err := o.receiver.HandleSwapRevert(ctx, msg.SharedID, types.DvpSwapCancelled); err != nil {
		slog.Error("Failed to handle DvpSwapCancelled",
			slog.Any("error", err),
			slog.String("messageId", msg.ID),
			slog.String("sharedId", msg.SharedID),
		)
		return false
	}

	slog.Info("Successfully handled DvpSwapCancelled",
		slog.String("messageId", msg.ID),
		slog.String("sharedId", msg.SharedID),
	)
	return true
}

func (o *DvpOrchestrator) handleSwapTimedOut(ctx context.Context, msg DvpDestMessage) bool {
	if msg.SharedID == "" {
		slog.Error("DvpSwapTimedOut message has empty sharedId", slog.String("messageId", msg.ID))
		return true
	}

	if err := o.receiver.HandleSwapRevert(ctx, msg.SharedID, types.DvpSwapTimedOut); err != nil {
		slog.Error("Failed to handle DvpSwapTimedOut",
			slog.Any("error", err),
			slog.String("messageId", msg.ID),
			slog.String("sharedId", msg.SharedID),
		)
		return false
	}

	slog.Info("Successfully handled DvpSwapTimedOut",
		slog.String("messageId", msg.ID),
		slog.String("sharedId", msg.SharedID),
	)
	return true
}
