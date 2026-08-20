package testtools

import (
	"context"
	"testing"
	"time"
)

func ShutdownFixture(t testing.TB, fn func(ctx context.Context) error, gracePeriod time.Duration) bool {
	t.Helper()

	shutdownSignal := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = fn(ctx)
		shutdownSignal <- struct{}{}
	}()

	// Check whether the function exits before context is canceled
	select {
	case <-shutdownSignal:
		return false
	case <-time.NewTimer(gracePeriod).C:
	}

	cancel()

	// Check whether the function exits gracefully when context is caneled
	select {
	case <-time.NewTimer(gracePeriod).C:
		return false
	case <-shutdownSignal:
		return true
	}
}
