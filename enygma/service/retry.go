package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	telemetry "github.com/raylsnetwork/rayls-sovereign-relayer/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type retryBlockWaiter interface {
	WaitForBlock(ctx context.Context, targetBlockNumber uint64) (uint64, error)
}

type retryTracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

const (
	// After 10 retries, wait for 2 blocks instead of 1
	// This extra waiting helps in times of congestion.
	// Critical operations like transfer, deposit and withdraw should have >10 maxRetries
	// Good and tested value for maxRetries in critical operations is 30
	retryThreshold = 10
)

type RetryService struct {
	tracer      retryTracer
	blockWaiter retryBlockWaiter
}

func NewRetryService(tracer retryTracer, blockWaiter retryBlockWaiter) *RetryService {
	return &RetryService{
		tracer:      tracer,
		blockWaiter: blockWaiter,
	}
}

// Retries the Enygma operation until it succeeds or the max retries are reached.
// Operations such as:
// - SupplyUpdate
// - CrossTransfer
// - Deposit
// - Withdrawal
// First, the operation is executed with the initial block number provided.
// On failure, retries with the next block number mined
func (rs *RetryService) RetryOperation(
	ctx context.Context,
	operationName string,
	maxRetries int,
	blockNumber uint64,
	executeOperation func(ctx context.Context, nextBlockNumber uint64) error,
) (uint64, error) {
	ctx, span := rs.tracer.Start(ctx, telemetry.SPAN_RETRY_ENYGMA_OPERATION)
	defer span.End()

	var err error

	for retryCount := 1; retryCount <= maxRetries; retryCount++ {
		// Determine how many blocks to wait based on retry count
		var blocksToWait uint64 = 1
		if retryCount > retryThreshold {
			blocksToWait = 2
		}

		// Wait for the next block before each retry.
		if retryCount > 1 {
			targetBlockNumber := blockNumber + blocksToWait - 1
			blockNumber, err = rs.blockWaiter.WaitForBlock(ctx, targetBlockNumber)
			if err != nil {
				return 0, fmt.Errorf("error while waiting for the next block number: %w", err)
			}
		}

		slog.Info(
			"Executing Enygma operation",
			slog.String("operation", operationName),
			slog.Int("retryCount", retryCount),
			slog.Int("maxRetries", maxRetries),
		)
		err = executeOperation(ctx, blockNumber)
		if err == nil {
			span.SetStatus(codes.Ok, telemetry.STATUS_SUCCESS_OPERATION_SUCCEEDED)
			return blockNumber, nil
		}

		if rs.isRetryableError(err) {
			slog.Info(
				"Enygma operation failed. Retrying...",
				slog.String("operation", operationName),
				slog.Int("retryCount", retryCount),
				slog.Int("maxRetries", maxRetries),
			)
			span.AddEvent(telemetry.EVENT_RETRYABLE_ERROR_ENCOUNTERED, trace.WithAttributes(
				attribute.Int(telemetry.ATTR_RETRY_COUNT, retryCount),
				attribute.String(telemetry.ATTR_ERROR_MESSAGE, err.Error()),
			))
			continue
		}

		return blockNumber, fmt.Errorf("non-retryable error encountered: %w", err)

	}

	// Max retries exceeded
	combinedErr := fmt.Errorf("max retries exceeded: %w", err)
	span.RecordError(combinedErr)
	span.SetStatus(codes.Error, telemetry.STATUS_ERROR_MAX_RETRIES_EXCEEDED)
	span.SetAttributes(
		attribute.Int(telemetry.ATTR_FINAL_RETRY_COUNT, maxRetries),
		//nolint:gosec // block numbers are within int range
		attribute.Int(telemetry.ATTR_FINAL_BLOCK_NUMBER, int(blockNumber)),
	)
	return blockNumber, combinedErr
}

// isRetryableError checks if an error should trigger a retry based on contract revert reasons
// or other specific failures
func (rs *RetryService) isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errorMsg := err.Error()

	// Check for timeout errors
	if strings.Contains(errorMsg, "context deadline exceeded") ||
		strings.Contains(errorMsg, "transaction mining failed") {
		return true
	}

	// Check for transaction/simulation inconsistencies
	if strings.Contains(errorMsg, "simulation did not revert") ||
		strings.Contains(errorMsg, "timing/state inconsistency") ||
		strings.Contains(errorMsg, "transaction did not revert") {
		return true
	}

	// Blockchain reject the transaction for some reason
	if strings.Contains(errorMsg, "Transaction processing could not be completed due to an exception") {
		return true
	}

	// List of smart contract error messages that should trigger retries
	retryableContractErrors := []string{
		"Invalid public signal for balance",
		"BlockNumber in Proof was already finalised.",
		"Invalid BlockNumber Used in Proof.",
		// Emitted by the Supply update when provided block number is outdated.
		// TODO: unify this error message in contracts to match other ones related to outdated block number
		"Invalid BlockNumber.",
		"Contract is processing another transaction.",
		"Nullifier already used in pending transaction.",
	}

	// Check for smart contract errors
	for _, retryableErr := range retryableContractErrors {
		if strings.Contains(errorMsg, retryableErr) {
			return true
		}
	}
	// DB - PNH Sync errors
	if strings.Contains(errorMsg, "Outdated/wrong submitted balance or r") ||
		strings.Contains(errorMsg, "Failed to sync fresh state from smart contract") ||
		strings.Contains(errorMsg, "database state still doesn't match smart contract after fresh sync") ||
		strings.Contains(errorMsg, "balance commitment mismatch between SC and DB") ||
		strings.Contains(errorMsg, "The circuit computed different values than provided") {
		return true
	}

	// Nonce errors
	if strings.Contains(strings.ToLower(errorMsg), "nonce too low") {
		return true
	}

	// Transient network/connection errors
	connectionErrors := []string{
		"connection refused",
		"connection reset",
		"dial tcp",
		"i/o timeout",
		"eof",
		"broken pipe",
		"no such host",
		"tls handshake timeout",
	}
	lowerMsg := strings.ToLower(errorMsg)
	for _, connErr := range connectionErrors {
		if strings.Contains(lowerMsg, connErr) {
			return true
		}
	}

	return false
}
