package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	telemetry "github.com/raylsnetwork/rayls-privacy-relayer-api/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type retryTracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

const (
	dvpRetryBaseDelay = 3 * time.Second
	dvpRetryDelayIncr = time.Second
)

type RetryService struct {
	tracer    retryTracer
	baseDelay time.Duration
	delayIncr time.Duration
}

func NewRetryService(tracer retryTracer) *RetryService {
	return &RetryService{
		tracer:    tracer,
		baseDelay: dvpRetryBaseDelay,
		delayIncr: dvpRetryDelayIncr,
	}
}

func NewRetryServiceWithDelays(tracer retryTracer, baseDelay, delayIncr time.Duration) *RetryService {
	return &RetryService{
		tracer:    tracer,
		baseDelay: baseDelay,
		delayIncr: delayIncr,
	}
}

func (rs *RetryService) RetryOperation(
	ctx context.Context,
	operationName string,
	maxRetries int,
	executeOperation func(ctx context.Context) error,
) error {
	ctx, span := rs.tracer.Start(ctx, telemetry.SPAN_RETRY_DVP_OPERATION)
	defer span.End()

	var err error
	delay := rs.baseDelay

	for retryCount := 1; retryCount <= maxRetries; retryCount++ {
		if retryCount > 1 {
			slog.Info("DVP operation failed. Retrying after delay...",
				slog.String("operation", operationName),
				slog.Int("retryCount", retryCount),
				slog.Int("maxRetries", maxRetries),
				slog.String("delay", delay.String()),
			)
			span.AddEvent(telemetry.EVENT_RETRYABLE_ERROR_ENCOUNTERED, trace.WithAttributes(
				attribute.Int(telemetry.ATTR_RETRY_COUNT, retryCount),
				attribute.String(telemetry.ATTR_ERROR_MESSAGE, err.Error()),
			))

			select {
			case <-ctx.Done():
				return fmt.Errorf("context cancelled during retry: %w", ctx.Err())
			case <-time.After(delay):
			}
			delay += rs.delayIncr
		}

		if ctx.Err() != nil {
			return fmt.Errorf("context cancelled before retry: %w", ctx.Err())
		}

		slog.Info("Executing DVP operation",
			slog.String("operation", operationName),
			slog.Int("retryCount", retryCount),
			slog.Int("maxRetries", maxRetries),
		)
		err = executeOperation(ctx)
		if err == nil {
			span.SetStatus(codes.Ok, telemetry.STATUS_SUCCESS_OPERATION_SUCCEEDED)
			return nil
		}

		if !rs.isRetryableError(err) {
			return fmt.Errorf("non-retryable error encountered: %w", err)
		}
	}

	combinedErr := fmt.Errorf("max retries exceeded: %w", err)
	span.RecordError(combinedErr)
	span.SetStatus(codes.Error, telemetry.STATUS_ERROR_MAX_RETRIES_EXCEEDED)
	span.SetAttributes(
		attribute.Int(telemetry.ATTR_FINAL_RETRY_COUNT, maxRetries),
	)
	return combinedErr
}

func (rs *RetryService) isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errorMsg := err.Error()

	retryableErrors := []string{
		"nonce too low",
		"transaction mining failed",
	}

	lowerMsg := strings.ToLower(errorMsg)
	for _, retryableErr := range retryableErrors {
		if strings.Contains(lowerMsg, retryableErr) {
			return true
		}
	}

	return false
}
