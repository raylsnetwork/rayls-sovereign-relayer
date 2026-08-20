package backoff

import (
	"context"
	"fmt"
	"time"
)

// Linear implements a linear backoff strategy.
// The delay increases by a fixed increment with each attempt.
// Example: with InitialDelay=3s and Increment=1s, delays are: 3s, 4s, 5s, 6s, 7s...
type Linear struct {
	// InitialDelay is the delay for the first retry attempt.
	InitialDelay time.Duration

	// Increment is the fixed amount added to the delay for each subsequent attempt.
	// Must be >= 0. Set to 0 for constant delay (no increase).
	Increment time.Duration

	// MaxDelay is the maximum delay to cap the growth.
	// Set to 0 for no maximum (unlimited growth).
	MaxDelay time.Duration
}

// NewLinear creates a new Linear backoff strategy with validation.
// Returns an error if the configuration is invalid.
func NewLinear(initialDelay time.Duration, increment time.Duration, maxDelay time.Duration) (*Linear, error) {
	if initialDelay <= 0 {
		return nil, fmt.Errorf("initialDelay must be positive, got %v", initialDelay)
	}
	if increment < 0 {
		return nil, fmt.Errorf("increment cannot be negative, got %v", increment)
	}
	if maxDelay < 0 {
		return nil, fmt.Errorf("maxDelay cannot be negative, got %v", maxDelay)
	}
	if maxDelay > 0 && maxDelay < initialDelay {
		return nil, fmt.Errorf("maxDelay (%v) cannot be less than initialDelay (%v)", maxDelay, initialDelay)
	}

	return &Linear{
		InitialDelay: initialDelay,
		Increment:    increment,
		MaxDelay:     maxDelay,
	}, nil
}

// Next calculates the delay for the given attempt number.
// Attempt numbers start at 1 (first retry).
func (l *Linear) Next(attempt int) time.Duration {
	if attempt <= 0 {
		return 0
	}

	// Calculate: InitialDelay + (Increment * (attempt - 1))
	delay := l.InitialDelay + (l.Increment * time.Duration(attempt-1))

	// Apply max delay cap if set
	if l.MaxDelay > 0 && delay > l.MaxDelay {
		return l.MaxDelay
	}

	return delay
}

// Do executes the given function with linear backoff retry logic.
// It attempts the function up to maxAttempts times, with linearly increasing delays between retries.
// Returns nil if the function succeeds, or the last error if all attempts fail.
func (l *Linear) Do(ctx context.Context, maxAttempts int, fn func() error) error {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Execute the function
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't sleep after the last attempt
		if attempt >= maxAttempts {
			break
		}

		// Calculate delay for next attempt
		delay := l.Next(attempt)

		// Wait with context cancellation support
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return lastErr
}
