package testtools_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools/spy"
)

func TestRunUntilAcked(t *testing.T) {
	t.Run("returns nil and stops fn after ack fires", func(t *testing.T) {
		ackSpy := spy.NewAck()

		// fn simulates an orchestrator-style Run loop: invoke ack once, then
		// block on ctx until cancelled.
		runStarted := make(chan struct{})
		fn := func(ctx context.Context) error {
			close(runStarted)
			// Process one message: call ack.
			_ = ackSpy.Fn()(ctx)
			<-ctx.Done()
			return nil
		}

		err := testtools.RunUntilAcked(t, fn, ackSpy, 2*time.Second)
		require.NoError(t, err)

		// Sanity: the goroutine actually started.
		select {
		case <-runStarted:
		case <-time.After(time.Second):
			t.Fatal("fn never started")
		}
	})

	t.Run("propagates fn error after ack fires", func(t *testing.T) {
		ackSpy := spy.NewAck()
		wantErr := errors.New("expected fn error")

		fn := func(ctx context.Context) error {
			_ = ackSpy.Fn()(ctx)
			<-ctx.Done()
			return wantErr
		}

		err := testtools.RunUntilAcked(t, fn, ackSpy, 2*time.Second)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("returns timeout error when ack is never called", func(t *testing.T) {
		ackSpy := spy.NewAck()

		fn := func(ctx context.Context) error {
			<-ctx.Done()
			return nil
		}

		err := testtools.RunUntilAcked(t, fn, ackSpy, 50*time.Millisecond)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "ack not called"),
			"expected ack-not-called diagnostic, got: %v", err)
	})

	t.Run("returns timeout error when fn ignores ctx cancel", func(t *testing.T) {
		ackSpy := spy.NewAck()
		done := make(chan struct{})
		t.Cleanup(func() { close(done) })

		// fn fires ack but then refuses to exit on ctx.Done(). It will exit
		// when the test's t.Cleanup fires `done`. This simulates a buggy
		// Run loop that doesn't honour cancellation.
		fn := func(ctx context.Context) error {
			_ = ackSpy.Fn()(ctx)
			<-done
			return nil
		}

		err := testtools.RunUntilAcked(t, fn, ackSpy, 50*time.Millisecond)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "did not exit"),
			"expected did-not-exit diagnostic, got: %v", err)
	})
}
