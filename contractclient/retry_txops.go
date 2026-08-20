package contractclient

import (
	"context"
	"log/slog"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/backoff"
	txopspb "github.com/raylsnetwork/rayls-sovereign-relayer/cts/gen/proto/txops"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Default retry tuning for transient CTS transport blips. Small and bounded —
// enough to ride out a GC pause or a momentary connection reset, not to wait
// out an outage. Sustained outages are handled durably by JetStream
// redelivery, not by in-process retry (see issue #210).
const (
	defaultRetryInitialDelay = time.Second
	defaultRetryMultiplier   = 2.0
	defaultRetryMaxDelay     = 1 * time.Minute
	defaultRetryMaxAttempts  = 10
)

// retryableReadCodes are transient gRPC status codes safe to retry for
// read-only Call operations. A read can be repeated freely, so even an
// ambiguous DeadlineExceeded is safe.
var retryableReadCodes = map[codes.Code]struct{}{
	codes.Unavailable:       {},
	codes.DeadlineExceeded:  {},
	codes.ResourceExhausted: {},
}

// retryableWriteCodes are transient gRPC status codes safe to retry for
// state-changing SignAndSend / BatchSignAndSend operations.
//
// These now carry a per-request idempotency key (the `id` on SignAndSendRequest),
// and CTS dedups on it via the cts_sync_tx ledger: a resend of an already-broadcast
// tx lands on the existing row and recovers (re-broadcasts the same signed bytes or
// returns the stored terminal verdict) instead of signing a second transaction. That
// removes the double-submit hazard that previously forced this set to the single
// unambiguous code (Unavailable = request provably never reached CTS).
//
// With that hazard gone, writes retry the same transient set as reads:
//   - Unavailable       — transport blip / CTS restart (incl. a crash mid-call, which
//     surfaces as Unavailable "error reading from server: EOF").
//   - DeadlineExceeded  — ambiguous (tx may already be mining); safe now because a
//     redundant resend is absorbed by the idempotency key (#211).
//   - ResourceExhausted — CTS backpressure; safe to retry after backoff.
//
// Retrying in-process here lets a transient CTS death be ridden out on the SAME id,
// avoiding the durable JetStream redelivery path — which re-runs the orchestrator
// handler and recomputes a fresh blockNumber (the source of post-recovery resync churn).
var retryableWriteCodes = map[codes.Code]struct{}{
	codes.Unavailable:       {},
	codes.DeadlineExceeded:  {},
	codes.ResourceExhausted: {},
}

// retryingTxOpsClient decorates a TxOpsServiceClient with bounded,
// idempotency-aware retry on transient gRPC errors. Retry lives at this layer
// (below CTSExecutor) deliberately: CTSExecutor flattens the gRPC status into a
// string via wrapTxOpsStatusError, after which the code can no longer be
// classified with status.FromError.
type retryingTxOpsClient struct {
	inner       TxOpsServiceClient
	strategy    backoff.Strategy
	maxAttempts int
}

// NewRetryingTxOpsClient wraps inner with the given backoff strategy and
// attempt budget. Callers that just want sane defaults should use
// NewDefaultRetryingTxOpsClient.
func NewRetryingTxOpsClient(inner TxOpsServiceClient, strategy backoff.Strategy, maxAttempts int) TxOpsServiceClient {
	return &retryingTxOpsClient{inner: inner, strategy: strategy, maxAttempts: maxAttempts}
}

// NewDefaultRetryingTxOpsClient wraps inner with the package default backoff
// (1s initial, 2x, capped at 1min, up to 10 attempts), suitable for riding out
// transient transport blips and short CTS restarts.
func NewDefaultRetryingTxOpsClient(inner TxOpsServiceClient) TxOpsServiceClient {
	strategy, err := backoff.NewExponential(defaultRetryInitialDelay, defaultRetryMultiplier, defaultRetryMaxDelay)
	if err != nil {
		// The arguments above are compile-time constants known to be valid, so
		// this branch is unreachable. Degrade to no retry rather than panic in
		// library code.
		return inner
	}
	return NewRetryingTxOpsClient(inner, strategy, defaultRetryMaxAttempts)
}

func (c *retryingTxOpsClient) SignAndSend(ctx context.Context, in *txopspb.SignAndSendRequest, opts ...grpc.CallOption) (*txopspb.SignAndSendResponse, error) {
	return retryTxOpsCall(ctx, in.Id, c.strategy, c.maxAttempts, retryableWriteCodes, func() (*txopspb.SignAndSendResponse, error) {
		return c.inner.SignAndSend(ctx, in, opts...)
	})
}

func (c *retryingTxOpsClient) BatchSignAndSend(ctx context.Context, in *txopspb.BatchSignAndSendRequest, opts ...grpc.CallOption) (*txopspb.BatchSignAndSendResponse, error) {
	// TODO: add a proper id once implemented in gRPC method
	return retryTxOpsCall(ctx, "[id]", c.strategy, c.maxAttempts, retryableWriteCodes, func() (*txopspb.BatchSignAndSendResponse, error) {
		return c.inner.BatchSignAndSend(ctx, in, opts...)
	})
}

func (c *retryingTxOpsClient) Call(ctx context.Context, in *txopspb.CallRequest, opts ...grpc.CallOption) (*txopspb.CallResponse, error) {
	// TODO: add a proper id once implemented in gRPC method
	return retryTxOpsCall(ctx, "[id]", c.strategy, c.maxAttempts, retryableReadCodes, func() (*txopspb.CallResponse, error) {
		return c.inner.Call(ctx, in, opts...)
	})
}

// retryTxOpsCall runs fn up to maxAttempts times, retrying only when the
// returned error is a gRPC status whose code is in retryable. It sleeps using
// strategy between attempts and aborts immediately if ctx is cancelled.
func retryTxOpsCall[T any](
	ctx context.Context,
	id string,
	strategy backoff.Strategy,
	maxAttempts int,
	retryable map[codes.Code]struct{},
	fn func() (T, error),
) (T, error) {
	var out T
	var err error
	for attempt := 1; ; attempt++ {
		out, err = fn()
		if err == nil {
			return out, nil
		}

		st, ok := status.FromError(err)
		if !ok {
			return out, err // not a gRPC status — not classifiable, don't retry
		}
		if _, retry := retryable[st.Code()]; !retry {
			return out, err // permanent / non-retryable code
		}

		slog.Debug("encountered retriable error while processing gRPC message", slog.String("id", id))
		if attempt >= maxAttempts {
			return out, err // budget exhausted — propagate the last error
		}

		// Give cancellation priority: a select between time.After and ctx.Done
		// picks randomly when both are ready (e.g. a zero backoff delay on an
		// already-cancelled context), so check explicitly first.
		if ctxErr := ctx.Err(); ctxErr != nil {
			return out, ctxErr
		}

		// Strategy.Next is 1-indexed: Next(1) is the delay before the first retry.
		select {
		case <-time.After(strategy.Next(attempt)):
		case <-ctx.Done():
			return out, ctx.Err()
		}
	}
}
