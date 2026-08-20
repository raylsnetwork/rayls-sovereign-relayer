package spy

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Ack is a test helper that tracks whether an acknowledgment function was
// called. It is safe to consult from the test goroutine while the function
// it returned from Fn is being invoked from another goroutine (e.g. the
// service-under-test's Run loop).
//
// Two ways to observe the call:
//
//   - Polling: AssertCalled / AssertNotCalled / Called check a boolean flag
//     and are intended for use AFTER the service-under-test is known to have
//     finished processing.
//   - Synchronisation: WaitForCall blocks until Fn is invoked or the timeout
//     elapses. Prefer this for tests that drive a Run-loop, so the test
//     doesn't have to guess how long the loop needs.
//
// Usage:
//
//	ackSpy := spy.NewAck()
//	msg := msgqueue.Message[T]{
//	    V:   value,
//	    Ack: ackSpy.Fn(),
//	}
//
//	// Synchronisation pattern (preferred for Run-loop tests):
//	go svc.Run(ctx)
//	require.True(t, ackSpy.WaitForCall(t, 2*time.Second))
//	cancel()
//
//	// Polling pattern (after svc.Run already returned):
//	ackSpy.AssertCalled(t)
type Ack struct {
	mu     sync.Mutex
	called bool
	signal chan struct{}
	once   sync.Once
}

// NewAck creates a new Ack spy for tracking acknowledgment calls.
func NewAck() *Ack {
	return &Ack{signal: make(chan struct{})}
}

// Fn returns a function that can be used as an Ack callback.
// On first call it flips the internal called flag and closes the
// synchronisation channel observed by WaitForCall. Subsequent calls
// are no-ops on the channel (idempotent via sync.Once).
func (a *Ack) Fn() func(context.Context) error {
	return func(context.Context) error {
		a.once.Do(func() {
			a.mu.Lock()
			a.called = true
			a.mu.Unlock()
			close(a.signal)
		})
		return nil
	}
}

// AssertCalled asserts that the ack function was called.
func (a *Ack) AssertCalled(t testing.TB) {
	t.Helper()
	a.mu.Lock()
	defer a.mu.Unlock()
	assert.True(t, a.called, "expected ack to be called")
}

// AssertNotCalled asserts that the ack function was not called.
func (a *Ack) AssertNotCalled(t testing.TB, msg string) {
	t.Helper()
	a.mu.Lock()
	defer a.mu.Unlock()
	assert.False(t, a.called, msg)
}

// Called returns whether the ack function was called.
func (a *Ack) Called() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.called
}

// WaitForCall blocks until the ack function has been called or the timeout
// elapses. Returns true if ack was called within the timeout. The timeout
// is a safety bound, not a target — under normal conditions the call returns
// as soon as Fn fires. Tests should pass a generous value (e.g. several
// seconds) so they don't flake on a slow runner.
func (a *Ack) WaitForCall(t testing.TB, timeout time.Duration) bool {
	t.Helper()
	select {
	case <-a.signal:
		return true
	case <-time.After(timeout):
		return false
	}
}
