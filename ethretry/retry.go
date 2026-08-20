// Package ethretry provides a small retry helper for transient Ethereum-RPC
// failures — network errors and the JSON decode errors flaky RPC endpoints
// occasionally return. It is shared by the CTS RPC paths (cts/ethrpc,
// cts/service) and the contract clients (contractclient), so it lives in a
// neutral top-level package rather than under any one of them.
package ethretry

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/url"
	"time"
)

const (
	retryDelay     = 3 * time.Second
	retryDelayIncr = time.Second
	maxRetryDelay  = 30 * time.Second
)

// WithRetry re-runs fn on transient failures (network errors, JSON
// decode errors from flaky RPC responses) until ctx is cancelled.
func WithRetry(ctx context.Context, fn func() error) error {
	delay := retryDelay
	for {
		err := fn()
		if err == nil {
			return nil
		}
		if !IsConnectionError(err) && !IsUnmarshalTypeError(err) {
			return err
		}

		delay = nextRetryDelay(delay, retryDelayIncr, maxRetryDelay)
		slog.Error(
			"rpc call failed, retrying",
			slog.String("delay", delay.String()),
			slog.Any("error", err),
		)

		select {
		case <-ctx.Done():
			return err
		case <-time.After(delay):
		}
	}
}

func nextRetryDelay(delay, incr, maxDelay time.Duration) time.Duration {
	next := delay + incr
	if next > maxDelay {
		return maxDelay
	}
	return next
}

func IsConnectionError(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	var urlErr *url.Error
	return errors.As(err, &urlErr)
}

func IsUnmarshalTypeError(err error) bool {
	var unmarshalErr *json.UnmarshalTypeError
	return errors.As(err, &unmarshalErr)
}
