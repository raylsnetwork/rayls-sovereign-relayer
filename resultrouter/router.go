// Package resultrouter consumes terminal TxResult messages from a single
// `cts.result.<identity>` subject and dispatches them to per-MessageType
// handlers registered by the relayer.
//
// One Router per identity-subject. Handlers are typically the
// HandleXxxResults callbacks on services like CrossChainService or
// GeneratorService.
//
// Semantics:
//   - At-least-once delivery — handlers must be idempotent.
//   - Handler returns nil → all results in the group are acked.
//   - Handler returns error → none of the results are acked; NATS
//     redelivers (MaxDeliver bounds the retry count).
//   - No handler registered for a MessageType → results are logged and
//     acked, so the queue does not stall on a misconfiguration. Operators
//     see the warning and add the missing handler.
package resultrouter

import (
	"context"
	"log/slog"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

//go:generate moq --pkg resultrouter_test -out router_mock_test.go . Consumer Handler

// Handler is invoked with a non-empty batch of results sharing the same
// MessageType. Returning a non-nil error skips the batch's ack, so NATS
// will redeliver — handlers must therefore be idempotent.
type Handler interface {
	Handle(ctx context.Context, results []types.TxResult) error
}

// HandlerFunc adapts a function to the Handler interface. Service
// callbacks like (*CrossChainService).HandleAtomicResults satisfy this
// signature directly.
type HandlerFunc func(ctx context.Context, results []types.TxResult) error

func (f HandlerFunc) Handle(ctx context.Context, results []types.TxResult) error {
	return f(ctx, results)
}

// Consumer is the subset of *msgqueue.Consumer[types.TxResult] the
// router needs. Defined as an interface for testability.
type Consumer interface {
	Fetch(ctx context.Context, count int) ([]msgqueue.Message[types.TxResult], error)
}

// Config bundles the tunables for a Router. Identity is informational
// only — used in logs to tell concurrent routers apart.
type Config struct {
	Identity  string
	BatchSize int
	Interval  time.Duration
}

// Router consumes a single cts.result.<identity> subject and fans
// batches of TxResults into per-MessageType handlers.
type Router struct {
	identity  string
	batchSize int
	interval  time.Duration

	consumer Consumer
	handlers map[string]Handler
}

// New constructs a Router. Register handlers via Register before calling
// Run. Register is not safe for concurrent use with Run.
func New(conf Config, consumer Consumer) *Router {
	return &Router{
		identity:  conf.Identity,
		batchSize: conf.BatchSize,
		interval:  conf.Interval,
		consumer:  consumer,
		handlers:  map[string]Handler{},
	}
}

// Register binds a handler to a MessageType. Calling Register twice with
// the same MessageType replaces the previous handler. Must be called
// before Run.
func (r *Router) Register(messageType string, h Handler) {
	r.handlers[messageType] = h
}

// Run polls the result subject and dispatches batches until ctx is
// cancelled. Returns nil on clean shutdown.
func (r *Router) Run(ctx context.Context) error {
	slog.Info("resultrouter started", slog.String("identity", r.identity))
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	initial := make(chan struct{}, 1)
	initial <- struct{}{}

	for {
		select {
		case <-ctx.Done():
			slog.Info("resultrouter shutting down", slog.String("identity", r.identity))
			return nil
		case <-ticker.C:
		case <-initial:
		}
		r.tick(ctx)
	}
}

// tick fetches one batch and dispatches it. Per-iteration failures are
// logged and swallowed; the next tick retries.
func (r *Router) tick(ctx context.Context) {
	msgs, err := r.consumer.Fetch(ctx, r.batchSize)
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		slog.Error("resultrouter: fetch failed",
			slog.String("identity", r.identity), slog.Any("err", err))
		return
	}
	if len(msgs) == 0 {
		return
	}

	groups := map[string][]msgqueue.Message[types.TxResult]{}
	for _, m := range msgs {
		groups[m.V.MessageType] = append(groups[m.V.MessageType], m)
	}

	for messageType, group := range groups {
		r.dispatchGroup(ctx, messageType, group)
	}
}

// dispatchGroup runs the registered handler for one MessageType group.
// On success, acks every message in the group. On error, leaves the group
// un-acked for redelivery. If no handler is registered the group is acked
// to keep the queue moving (the misconfiguration is logged).
func (r *Router) dispatchGroup(
	ctx context.Context,
	messageType string,
	group []msgqueue.Message[types.TxResult],
) {
	handler, ok := r.handlers[messageType]
	if !ok {
		slog.Warn("resultrouter: no handler registered, acking to unblock queue",
			slog.String("identity", r.identity),
			slog.String("message_type", messageType),
			slog.Int("count", len(group)))
		ackAll(ctx, group)
		return
	}

	results := make([]types.TxResult, len(group))
	for i, m := range group {
		results[i] = m.V
	}

	if err := handler.Handle(ctx, results); err != nil {
		slog.Error("resultrouter: handler failed, leaving messages un-acked",
			slog.String("identity", r.identity),
			slog.String("message_type", messageType),
			slog.Int("count", len(group)),
			slog.Any("err", err))
		return
	}

	ackAll(ctx, group)
}

func ackAll(ctx context.Context, group []msgqueue.Message[types.TxResult]) {
	for _, m := range group {
		if err := m.Ack(ctx); err != nil {
			slog.Error("resultrouter: ack failed",
				slog.String("correlation_id", m.V.CorrelationID),
				slog.Any("err", err))
		}
	}
}
