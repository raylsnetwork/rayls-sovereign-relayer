package backoff

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLinear(t *testing.T) {
	t.Run("creates valid linear backoff", func(t *testing.T) {
		lin, err := NewLinear(3*time.Second, time.Second, 0)
		require.NoError(t, err)
		assert.Equal(t, 3*time.Second, lin.InitialDelay)
		assert.Equal(t, time.Second, lin.Increment)
		assert.Equal(t, time.Duration(0), lin.MaxDelay)
	})

	t.Run("creates valid linear with max delay", func(t *testing.T) {
		lin, err := NewLinear(2*time.Second, time.Second, 10*time.Second)
		require.NoError(t, err)
		assert.Equal(t, 10*time.Second, lin.MaxDelay)
	})

	t.Run("rejects zero initial delay", func(t *testing.T) {
		_, err := NewLinear(0, time.Second, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "initialDelay must be positive")
	})

	t.Run("rejects negative initial delay", func(t *testing.T) {
		_, err := NewLinear(-time.Second, time.Second, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "initialDelay must be positive")
	})

	t.Run("rejects negative increment", func(t *testing.T) {
		_, err := NewLinear(time.Second, -time.Second, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "increment cannot be negative")
	})

	t.Run("accepts zero increment for constant delay", func(t *testing.T) {
		lin, err := NewLinear(5*time.Second, 0, 0)
		require.NoError(t, err)
		assert.Equal(t, time.Duration(0), lin.Increment)
	})

	t.Run("rejects negative max delay", func(t *testing.T) {
		_, err := NewLinear(time.Second, time.Second, -time.Second)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "maxDelay cannot be negative")
	})

	t.Run("rejects max delay less than initial delay", func(t *testing.T) {
		_, err := NewLinear(5*time.Second, time.Second, 3*time.Second)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "maxDelay")
		assert.Contains(t, err.Error(), "cannot be less than initialDelay")
	})
}

func TestLinear_Next(t *testing.T) {
	t.Run("returns 0 for attempt 0 or negative", func(t *testing.T) {
		lin, _ := NewLinear(time.Second, time.Second, 0)
		assert.Equal(t, time.Duration(0), lin.Next(0))
		assert.Equal(t, time.Duration(0), lin.Next(-1))
	})

	t.Run("calculates linear growth", func(t *testing.T) {
		lin, _ := NewLinear(3*time.Second, time.Second, 0)

		// 3s, 4s, 5s, 6s, 7s
		assert.Equal(t, 3*time.Second, lin.Next(1))
		assert.Equal(t, 4*time.Second, lin.Next(2))
		assert.Equal(t, 5*time.Second, lin.Next(3))
		assert.Equal(t, 6*time.Second, lin.Next(4))
		assert.Equal(t, 7*time.Second, lin.Next(5))
	})

	t.Run("calculates linear growth with different increment", func(t *testing.T) {
		lin, _ := NewLinear(2*time.Second, 500*time.Millisecond, 0)

		// 2s, 2.5s, 3s, 3.5s
		assert.Equal(t, 2*time.Second, lin.Next(1))
		assert.Equal(t, 2500*time.Millisecond, lin.Next(2))
		assert.Equal(t, 3*time.Second, lin.Next(3))
		assert.Equal(t, 3500*time.Millisecond, lin.Next(4))
	})

	t.Run("respects max delay cap", func(t *testing.T) {
		lin, _ := NewLinear(3*time.Second, time.Second, 6*time.Second)

		// 3s, 4s, 5s, then capped at 6s
		assert.Equal(t, 3*time.Second, lin.Next(1))
		assert.Equal(t, 4*time.Second, lin.Next(2))
		assert.Equal(t, 5*time.Second, lin.Next(3))
		assert.Equal(t, 6*time.Second, lin.Next(4))  // Would be 6s, at cap
		assert.Equal(t, 6*time.Second, lin.Next(5))  // Would be 7s, capped at 6s
		assert.Equal(t, 6*time.Second, lin.Next(10)) // Would be 12s, capped at 6s
	})

	t.Run("handles constant delay when increment is 0", func(t *testing.T) {
		lin, _ := NewLinear(5*time.Second, 0, 0)

		// Should remain constant at 5s
		assert.Equal(t, 5*time.Second, lin.Next(1))
		assert.Equal(t, 5*time.Second, lin.Next(2))
		assert.Equal(t, 5*time.Second, lin.Next(3))
		assert.Equal(t, 5*time.Second, lin.Next(10))
	})

	t.Run("matches evm-client behavior", func(t *testing.T) {
		// EVM client uses: InitialDelay=3s, Increment=1s
		lin, _ := NewLinear(3*time.Second, time.Second, 0)

		// Should match the cumulative duration pattern
		assert.Equal(t, 3*time.Second, lin.Next(1))
		assert.Equal(t, 4*time.Second, lin.Next(2))
		assert.Equal(t, 5*time.Second, lin.Next(3))
	})

	t.Run("matches httpclient behavior with 2s backoff", func(t *testing.T) {
		// HTTP client uses: backoffDuration=2s with cumulative growth
		lin, _ := NewLinear(2*time.Second, 2*time.Second, 0)

		// 2s, 4s, 6s, 8s
		assert.Equal(t, 2*time.Second, lin.Next(1))
		assert.Equal(t, 4*time.Second, lin.Next(2))
		assert.Equal(t, 6*time.Second, lin.Next(3))
		assert.Equal(t, 8*time.Second, lin.Next(4))
	})
}

func TestLinear_ImplementsStrategy(t *testing.T) {
	// Compile-time check that Linear implements Strategy
	var _ Strategy = (*Linear)(nil)

	lin, _ := NewLinear(3*time.Second, time.Second, 0)
	var strategy Strategy = lin

	// Use through interface
	delay := strategy.Next(3)
	assert.Equal(t, 5*time.Second, delay)
}

func TestLinear_Do(t *testing.T) {
	t.Run("succeeds on first attempt", func(t *testing.T) {
		lin, _ := NewLinear(time.Second, time.Second, 0)
		ctx := context.Background()

		callCount := 0
		err := lin.Do(ctx, 5, func() error {
			callCount++
			return nil
		})

		assert.NoError(t, err)
		assert.Equal(t, 1, callCount)
	})

	t.Run("succeeds on third attempt", func(t *testing.T) {
		lin, _ := NewLinear(100*time.Millisecond, 100*time.Millisecond, 0)
		ctx := context.Background()

		callCount := 0
		err := lin.Do(ctx, 5, func() error {
			callCount++
			if callCount < 3 {
				return fmt.Errorf("attempt %d failed", callCount)
			}
			return nil
		})

		assert.NoError(t, err)
		assert.Equal(t, 3, callCount)
	})

	t.Run("fails after max attempts", func(t *testing.T) {
		lin, _ := NewLinear(10*time.Millisecond, 10*time.Millisecond, 0)
		ctx := context.Background()

		callCount := 0
		err := lin.Do(ctx, 3, func() error {
			callCount++
			return fmt.Errorf("always fails")
		})

		assert.Error(t, err)
		assert.Equal(t, "always fails", err.Error())
		assert.Equal(t, 3, callCount)
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		lin, _ := NewLinear(time.Second, time.Second, 0)
		ctx, cancel := context.WithCancel(context.Background())

		callCount := 0
		// Cancel after first failed attempt
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		err := lin.Do(ctx, 10, func() error {
			callCount++
			return fmt.Errorf("always fails")
		})

		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
		assert.LessOrEqual(t, callCount, 2) // Should stop early due to cancellation
	})

	t.Run("respects context timeout", func(t *testing.T) {
		lin, _ := NewLinear(300*time.Millisecond, 300*time.Millisecond, 0)
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		callCount := 0
		err := lin.Do(ctx, 10, func() error {
			callCount++
			return fmt.Errorf("always fails")
		})

		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.LessOrEqual(t, callCount, 3) // Should timeout before 10 attempts
	})

	t.Run("uses linear delays", func(t *testing.T) {
		lin, _ := NewLinear(100*time.Millisecond, 100*time.Millisecond, 0)
		ctx := context.Background()

		callTimes := []time.Time{}
		err := lin.Do(ctx, 4, func() error {
			callTimes = append(callTimes, time.Now())
			return fmt.Errorf("fail")
		})

		assert.Error(t, err)
		assert.Len(t, callTimes, 4)

		// Check delays: ~100ms, ~200ms, ~300ms
		// Allow some tolerance for timing variation
		delay1 := callTimes[1].Sub(callTimes[0])
		delay2 := callTimes[2].Sub(callTimes[1])
		delay3 := callTimes[3].Sub(callTimes[2])

		assert.InDelta(t, 100*time.Millisecond, delay1, float64(50*time.Millisecond))
		assert.InDelta(t, 200*time.Millisecond, delay2, float64(50*time.Millisecond))
		assert.InDelta(t, 300*time.Millisecond, delay3, float64(50*time.Millisecond))
	})
}
