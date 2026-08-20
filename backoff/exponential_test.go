package backoff

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewExponential(t *testing.T) {
	t.Run("creates valid exponential with default multiplier", func(t *testing.T) {
		exp, err := NewExponential(time.Second, 2.0, 0)
		require.NoError(t, err)
		assert.Equal(t, time.Second, exp.InitialDelay)
		assert.Equal(t, 2.0, exp.Multiplier)
		assert.Equal(t, time.Duration(0), exp.MaxDelay)
	})

	t.Run("creates valid exponential with max delay", func(t *testing.T) {
		exp, err := NewExponential(time.Second, 2.0, 10*time.Second)
		require.NoError(t, err)
		assert.Equal(t, 10*time.Second, exp.MaxDelay)
	})

	t.Run("rejects zero initial delay", func(t *testing.T) {
		_, err := NewExponential(0, 2.0, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "initialDelay must be positive")
	})

	t.Run("rejects negative initial delay", func(t *testing.T) {
		_, err := NewExponential(-time.Second, 2.0, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "initialDelay must be positive")
	})

	t.Run("rejects multiplier less than 1", func(t *testing.T) {
		_, err := NewExponential(time.Second, 0.5, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "multiplier must be >= 1.0")
	})

	t.Run("accepts multiplier equal to 1", func(t *testing.T) {
		exp, err := NewExponential(time.Second, 1.0, 0)
		require.NoError(t, err)
		assert.Equal(t, 1.0, exp.Multiplier)
	})

	t.Run("rejects negative max delay", func(t *testing.T) {
		_, err := NewExponential(time.Second, 2.0, -time.Second)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "maxDelay cannot be negative")
	})

	t.Run("rejects max delay less than initial delay", func(t *testing.T) {
		_, err := NewExponential(5*time.Second, 2.0, 3*time.Second)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "maxDelay")
		assert.Contains(t, err.Error(), "cannot be less than initialDelay")
	})
}

func TestExponential_Next(t *testing.T) {
	t.Run("returns 0 for attempt 0 or negative", func(t *testing.T) {
		exp, _ := NewExponential(time.Second, 2.0, 0)
		assert.Equal(t, time.Duration(0), exp.Next(0))
		assert.Equal(t, time.Duration(0), exp.Next(-1))
	})

	t.Run("calculates exponential growth with multiplier 2.0", func(t *testing.T) {
		exp, _ := NewExponential(time.Second, 2.0, 0)

		// 1s, 2s, 4s, 8s, 16s
		assert.Equal(t, 1*time.Second, exp.Next(1))
		assert.Equal(t, 2*time.Second, exp.Next(2))
		assert.Equal(t, 4*time.Second, exp.Next(3))
		assert.Equal(t, 8*time.Second, exp.Next(4))
		assert.Equal(t, 16*time.Second, exp.Next(5))
	})

	t.Run("calculates exponential growth with multiplier 1.5", func(t *testing.T) {
		exp, _ := NewExponential(2*time.Second, 1.5, 0)

		// 2s, 3s, 4.5s, 6.75s
		assert.Equal(t, 2*time.Second, exp.Next(1))
		assert.Equal(t, 3*time.Second, exp.Next(2))
		assert.Equal(t, 4500*time.Millisecond, exp.Next(3))
		assert.Equal(t, 6750*time.Millisecond, exp.Next(4))
	})

	t.Run("respects max delay cap", func(t *testing.T) {
		exp, _ := NewExponential(time.Second, 2.0, 10*time.Second)

		// 1s, 2s, 4s, 8s, then capped at 10s
		assert.Equal(t, 1*time.Second, exp.Next(1))
		assert.Equal(t, 2*time.Second, exp.Next(2))
		assert.Equal(t, 4*time.Second, exp.Next(3))
		assert.Equal(t, 8*time.Second, exp.Next(4))
		assert.Equal(t, 10*time.Second, exp.Next(5))  // Would be 16s, capped at 10s
		assert.Equal(t, 10*time.Second, exp.Next(6))  // Would be 32s, capped at 10s
		assert.Equal(t, 10*time.Second, exp.Next(10)) // Still capped
	})

	t.Run("handles constant delay when multiplier is 1.0", func(t *testing.T) {
		exp, _ := NewExponential(5*time.Second, 1.0, 0)

		// Should remain constant at 5s
		assert.Equal(t, 5*time.Second, exp.Next(1))
		assert.Equal(t, 5*time.Second, exp.Next(2))
		assert.Equal(t, 5*time.Second, exp.Next(3))
		assert.Equal(t, 5*time.Second, exp.Next(10))
	})
}

func TestExponential_ImplementsStrategy(t *testing.T) {
	// Compile-time check that Exponential implements Strategy
	var _ Strategy = (*Exponential)(nil)

	exp, _ := NewExponential(time.Second, 2.0, 0)
	var strategy Strategy = exp

	// Use through interface
	delay := strategy.Next(3)
	assert.Equal(t, 4*time.Second, delay)
}

func TestExponential_Do(t *testing.T) {
	t.Run("succeeds on first attempt", func(t *testing.T) {
		exp, _ := NewExponential(time.Second, 2.0, 0)
		ctx := context.Background()

		callCount := 0
		err := exp.Do(ctx, 5, func() error {
			callCount++
			return nil
		})

		assert.NoError(t, err)
		assert.Equal(t, 1, callCount)
	})

	t.Run("succeeds on third attempt", func(t *testing.T) {
		exp, _ := NewExponential(100*time.Millisecond, 2.0, 0)
		ctx := context.Background()

		callCount := 0
		err := exp.Do(ctx, 5, func() error {
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
		exp, _ := NewExponential(10*time.Millisecond, 2.0, 0)
		ctx := context.Background()

		callCount := 0
		err := exp.Do(ctx, 3, func() error {
			callCount++
			return fmt.Errorf("always fails")
		})

		assert.Error(t, err)
		assert.Equal(t, "always fails", err.Error())
		assert.Equal(t, 3, callCount)
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		exp, _ := NewExponential(time.Second, 2.0, 0)
		ctx, cancel := context.WithCancel(context.Background())

		callCount := 0
		// Cancel after first failed attempt
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		err := exp.Do(ctx, 10, func() error {
			callCount++
			return fmt.Errorf("always fails")
		})

		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
		assert.LessOrEqual(t, callCount, 2) // Should stop early due to cancellation
	})

	t.Run("respects context timeout", func(t *testing.T) {
		exp, _ := NewExponential(500*time.Millisecond, 2.0, 0)
		ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
		defer cancel()

		callCount := 0
		err := exp.Do(ctx, 10, func() error {
			callCount++
			return fmt.Errorf("always fails")
		})

		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.LessOrEqual(t, callCount, 3) // Should timeout before 10 attempts
	})

	t.Run("uses exponential delays", func(t *testing.T) {
		exp, _ := NewExponential(100*time.Millisecond, 2.0, 0)
		ctx := context.Background()

		callTimes := []time.Time{}
		err := exp.Do(ctx, 4, func() error {
			callTimes = append(callTimes, time.Now())
			return fmt.Errorf("fail")
		})

		assert.Error(t, err)
		assert.Len(t, callTimes, 4)

		// Check delays: ~100ms, ~200ms, ~400ms
		// Allow some tolerance for timing variation
		delay1 := callTimes[1].Sub(callTimes[0])
		delay2 := callTimes[2].Sub(callTimes[1])
		delay3 := callTimes[3].Sub(callTimes[2])

		assert.InDelta(t, 100*time.Millisecond, delay1, float64(50*time.Millisecond))
		assert.InDelta(t, 200*time.Millisecond, delay2, float64(50*time.Millisecond))
		assert.InDelta(t, 400*time.Millisecond, delay3, float64(50*time.Millisecond))
	})
}
