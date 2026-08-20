package backoff

import (
	"context"
	"fmt"
	"math"
	"time"
)

// Exponential implements an exponential backoff strategy.
// The delay grows by multiplying the previous delay by a multiplier factor.
// Example: with InitialDelay=1s and Multiplier=2.0, delays are: 1s, 2s, 4s, 8s, 16s...
type Exponential struct {
	// InitialDelay is the delay for the first retry attempt.
	InitialDelay time.Duration

	// Multiplier is the factor by which the delay grows each attempt.
	// Typical value is 2.0 for exponential doubling.
	// Must be >= 1.0.
	Multiplier float64

	// MaxDelay is the maximum delay to cap the growth.
	// Set to 0 for no maximum (unlimited growth).
	MaxDelay time.Duration
}

// NewExponential creates a new Exponential backoff strategy with validation.
// Returns an error if the configuration is invalid.
func NewExponential(initialDelay time.Duration, multiplier float64, maxDelay time.Duration) (*Exponential, error) {
	if initialDelay <= 0 {
		return nil, fmt.Errorf("initialDelay must be positive, got %v", initialDelay)
	}
	if multiplier < 1.0 {
		return nil, fmt.Errorf("multiplier must be >= 1.0, got %v", multiplier)
	}
	if maxDelay < 0 {
		return nil, fmt.Errorf("maxDelay cannot be negative, got %v", maxDelay)
	}
	if maxDelay > 0 && maxDelay < initialDelay {
		return nil, fmt.Errorf("maxDelay (%v) cannot be less than initialDelay (%v)", maxDelay, initialDelay)
	}

	return &Exponential{
		InitialDelay: initialDelay,
		Multiplier:   multiplier,
		MaxDelay:     maxDelay,
	}, nil
}

// Next calculates the delay for the given attempt number.
// Attempt numbers start at 1 (first retry).
func (e *Exponential) Next(attempt int) time.Duration {
	if attempt <= 0 {
		return 0
	}

	// Calculate: InitialDelay * (Multiplier ^ (attempt - 1))
	delay := float64(e.InitialDelay) * math.Pow(e.Multiplier, float64(attempt-1))

	// Apply max delay cap if set
	if e.MaxDelay > 0 && delay > float64(e.MaxDelay) {
		return e.MaxDelay
	}

	return time.Duration(delay)
}

// Do executes the given function with exponential backoff retry logic.
// It attempts the function up to maxAttempts times, with increasing delays between retries.
// Returns nil if the function succeeds, or the last error if all attempts fail.
func (e *Exponential) Do(ctx context.Context, maxAttempts int, fn func() error) error {
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
		delay := e.Next(attempt)

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
