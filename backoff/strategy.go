package backoff

import (
	"context"
	"time"
)

// Strategy defines the interface for backoff strategies.
// Implementations should be stateless and calculate delays based on the attempt number.
type Strategy interface {
	// Next returns the delay duration for the given attempt number (1-indexed).
	// Attempt 1 is the first retry, attempt 2 is the second retry, etc.
	Next(attempt int) time.Duration

	// Do executes the given function with retry logic using this backoff strategy.
	// It retries up to maxAttempts times, using the Next() method to calculate delays.
	// Returns nil if the function succeeds, or the last error if all attempts fail.
	// Respects context cancellation between retry attempts.
	Do(ctx context.Context, maxAttempts int, fn func() error) error
}
