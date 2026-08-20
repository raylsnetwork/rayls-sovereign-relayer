package testtools

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools/spy"
)

// RunUntilAcked runs `fn` in a goroutine, waits up to `timeout` for `ackSpy`
// to be called, then cancels `fn`'s context and returns whatever error `fn`
// returned.
//
// Use this for tests that drive a Run-loop (e.g. an orchestrator) and want
// to observe "one message was processed end-to-end (handler + ack) and the
// loop terminated cleanly" without relying on tight timing assumptions.
//
// The `timeout` is a safety bound, not a target — under normal conditions
// the function returns within microseconds of the ack callback firing.
// Tests should pass a generous value (a few seconds) so they don't flake
// on a slow runner. Failure of either gate (ack-not-fired or run-did-not-
// exit-cleanly) is reported via the testing.TB so the call site can use
// require / assert as it normally would.
//
// Usage:
//
//	svc := NewMyOrchestratorWithDeps(mq, initiator)
//	err := testtools.RunUntilAcked(t, svc.Run, ackSpy, 2*time.Second)
//	require.NoError(t, err)
//	// ... subsequent state assertions on the initiator / mocks ...
func RunUntilAcked(
	t testing.TB,
	fn func(ctx context.Context) error,
	ackSpy *spy.Ack,
	timeout time.Duration,
) error {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runErr := make(chan error, 1)
	go func() { runErr <- fn(ctx) }()

	if !ackSpy.WaitForCall(t, timeout) {
		cancel()
		// Drain the goroutine to avoid leaking it. Best-effort — the outer
		// failure message is the test's concern, not the drain.
		select {
		case <-runErr:
		case <-time.After(timeout):
		}
		return fmt.Errorf("RunUntilAcked: ack not called within %s", timeout)
	}

	cancel()
	select {
	case err := <-runErr:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("RunUntilAcked: fn did not exit within %s after ctx cancel", timeout)
	}
}
