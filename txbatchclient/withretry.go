package txbatchclient

import (
	"context"
	"fmt"
	"log/slog"
	"time"

)

const (
	retryDelay     = 3 * time.Second
	retryDelayIncr = time.Second
)

func withRetry(ctx context.Context, fn func() error) error {
	var (
		err   error
		delay = retryDelay
	)

	for {
		err = fn()
		if err == nil {
			return nil
		}
		if !isConnectionError(err) && !isUnmarshalTypeError(err) {
			return fmt.Errorf("failed to execute with retry: %w", err)
		}

		delay += retryDelayIncr
		slog.Error(
			"failed operation..retrying after delay",
			slog.String("delay", delay.String()),
			slog.Any("err", err),
		)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
}
