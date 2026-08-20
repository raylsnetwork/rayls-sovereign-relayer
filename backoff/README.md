# Backoff Package

Provides configurable backoff strategies for retry logic with exponential and linear growth options.

## Features

- **Exponential backoff**: Delay grows by multiplying (e.g., 1s → 2s → 4s → 8s...)
- **Linear backoff**: Delay grows by fixed increment (e.g., 3s → 4s → 5s → 6s...)
- **Optional max delay cap**: Prevents unbounded growth
- **Context cancellation support**: Respects `context.Context` for timeouts and cancellation
- **Two APIs**: Low-level `Next()` for custom logic, high-level `Do()` for simple retries

## Usage

### Exponential Backoff

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/raylsnetwork/rayls-sovereign-relayer/backoff"
)

func main() {
    // Create exponential backoff: 1s, 2s, 4s, 8s... (capped at 30s)
    exp, _ := backoff.NewExponential(
        time.Second,      // initial delay
        2.0,              // multiplier
        30*time.Second,   // max delay (0 for unlimited)
    )

    // High-level API: automatically retries with backoff
    ctx := context.Background()
    err := exp.Do(ctx, 5, func() error {
        // Your operation here
        return someOperation()
    })

    // Low-level API: manual control
    for attempt := 1; attempt <= 5; attempt++ {
        if err := someOperation(); err == nil {
            break
        }
        delay := exp.Next(attempt)
        time.Sleep(delay)
    }
}
```

### Linear Backoff

```go
// Create linear backoff: 3s, 4s, 5s, 6s...
lin, _ := backoff.NewLinear(
    3*time.Second,    // initial delay
    time.Second,      // increment per retry
    10*time.Second,   // max delay (0 for unlimited)
)

// Use with context timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

err := lin.Do(ctx, 10, func() error {
    return makeDatabaseCall()
})
```

## Real-World Examples

### Replace EVM Client Retry

Current pattern in `evm-client.go`:
```go
// Before: manual retry loop
delay := 3 * time.Second
for i := 0; i < 100; i++ {
    result, err := operation()
    if err == nil {
        return result, nil
    }
    delay += time.Second
    time.Sleep(delay)
}
```

With backoff package:
```go
// After: using backoff
lin, _ := backoff.NewLinear(3*time.Second, time.Second, 0)
var result interface{}
err := lin.Do(ctx, 100, func() error {
    var err error
    result, err = operation()
    return err
})
```

### Replace LogRouter Exponential

Current pattern in `logrouter.go`:
```go
// Before: manual exponential logic
duration := time.Second
for {
    err := f()
    if err == nil {
        break
    }
    time.Sleep(duration)
    duration *= 2
}
```

With backoff package:
```go
// After: using backoff
exp, _ := backoff.NewExponential(time.Second, 2.0, 0)
exp.Do(ctx, math.MaxInt, f) // Unlimited retries
```

## API Reference

### Strategy Interface

```go
type Strategy interface {
    Next(attempt int) time.Duration
    Do(ctx context.Context, maxAttempts int, fn func() error) error
}
```

### Exponential

```go
type Exponential struct {
    InitialDelay time.Duration
    Multiplier   float64
    MaxDelay     time.Duration
}

func NewExponential(initialDelay time.Duration, multiplier float64, maxDelay time.Duration) (*Exponential, error)
func (e *Exponential) Next(attempt int) time.Duration
func (e *Exponential) Do(ctx context.Context, maxAttempts int, fn func() error) error
```

**Formula**: `delay = InitialDelay * (Multiplier ^ (attempt - 1))`

### Linear

```go
type Linear struct {
    InitialDelay time.Duration
    Increment    time.Duration
    MaxDelay     time.Duration
}

func NewLinear(initialDelay time.Duration, increment time.Duration, maxDelay time.Duration) (*Linear, error)
func (l *Linear) Next(attempt int) time.Duration
func (l *Linear) Do(ctx context.Context, maxAttempts int, fn func() error) error
```

**Formula**: `delay = InitialDelay + (Increment * (attempt - 1))`

## Design Decisions

1. **Attempt numbering**: 1-indexed (attempt 1 = first retry)
2. **Stateless**: No internal state, calculations based on attempt number
3. **Context-aware**: Respects cancellation and timeouts
4. **Max delay cap**: Optional ceiling to prevent unbounded delays
5. **Constructor validation**: Ensures valid configuration upfront
