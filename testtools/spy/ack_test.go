package spy_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools/spy"
)

func TestAckBehaviour(t *testing.T) {
	t.Run("starts not-called", func(t *testing.T) {
		a := spy.NewAck()
		assert.False(t, a.Called())
	})

	t.Run("Fn flips Called to true on first invocation", func(t *testing.T) {
		a := spy.NewAck()
		require.NoError(t, a.Fn()(context.Background()))
		assert.True(t, a.Called())
	})

	t.Run("AssertCalled passes after Fn was invoked", func(t *testing.T) {
		a := spy.NewAck()
		_ = a.Fn()(context.Background())
		a.AssertCalled(t)
	})

	t.Run("AssertNotCalled passes before Fn was invoked", func(t *testing.T) {
		a := spy.NewAck()
		a.AssertNotCalled(t, "fresh spy must not be called")
	})

	t.Run("Fn is idempotent on the synchronisation signal", func(t *testing.T) {
		a := spy.NewAck()
		fn := a.Fn()
		require.NoError(t, fn(context.Background()))
		// Second call must not panic on close(signal); sync.Once guards it.
		require.NoError(t, fn(context.Background()))
		assert.True(t, a.Called())
	})
}

func TestAck_WaitForCall(t *testing.T) {
	t.Run("returns true immediately once Fn fires", func(t *testing.T) {
		a := spy.NewAck()
		go func() {
			_ = a.Fn()(context.Background())
		}()
		require.True(t, a.WaitForCall(t, 2*time.Second))
	})

	t.Run("returns true after Fn fired earlier", func(t *testing.T) {
		a := spy.NewAck()
		_ = a.Fn()(context.Background())
		// Should not block — channel is already closed.
		require.True(t, a.WaitForCall(t, time.Millisecond))
	})

	t.Run("returns false when Fn is never called within the timeout", func(t *testing.T) {
		a := spy.NewAck()
		// Sub-tested with a tight timeout because we WANT it to time out;
		// the goroutine scheduling slack is fine in the failure direction.
		assert.False(t, a.WaitForCall(t, 10*time.Millisecond))
	})

	t.Run("multiple goroutines can WaitForCall concurrently", func(t *testing.T) {
		a := spy.NewAck()
		var wg sync.WaitGroup
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				assert.True(t, a.WaitForCall(t, 2*time.Second))
			}()
		}
		// Give the waiters a moment to subscribe to the channel.
		time.Sleep(10 * time.Millisecond)
		_ = a.Fn()(context.Background())
		wg.Wait()
	})
}
