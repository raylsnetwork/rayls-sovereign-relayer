package repository

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestConnect_RetriesPingUntilContextDeadline asserts the startup ping loop
// keeps retrying until the caller's context deadline elapses, then surfaces
// an error wrap that names both the database and the elapsed retry budget so
// an operator can diagnose "DB never came up" failures without re-reading logs.
//
// Pre-2026-05-12 the function pinged exactly once and bailed; the postgres
// interrupted-shutdown recovery window (~65s observed in the outage) was
// enough to make that single ping fail, which crash-looped every dependent
// service. The retry loop added in postgres.go absorbs that window.
func TestConnect_RetriesPingUntilContextDeadline(t *testing.T) {
	// Address that nothing is listening on, so every ping attempt fails with
	// connection-refused. Bind a TCP socket to an ephemeral port, close it,
	// then point Connect at that port — the kernel won't immediately re-bind,
	// so we get a reliable connection-refused for the test duration.
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := l.Addr().(*net.TCPAddr)
	require.NoError(t, l.Close())

	connStr := "postgres://x:x@" + addr.String() + "/no_such_db?sslmode=disable"

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err = Connect(ctx, connStr)
	elapsed := time.Since(start)

	require.Error(t, err, "Connect must surface an error when ping never succeeds")

	// The error must blame the actual database (so operators reading logs
	// know which connection is broken) and include attempt count (so it's
	// distinguishable from a single-shot ping failure).
	msg := err.Error()
	require.Contains(t, msg, "no_such_db", "error must name the database; got: %s", msg)
	require.Contains(t, msg, "attempts", "error must report the attempt count; got: %s", msg)

	// We retried at least twice (initial + at least one backoff). The first
	// ping fails fast (~250ms ConnectTimeout doesn't apply to TCP RST), and
	// pingRetryInitialBackoff = 250ms means we should have observed ≥2 attempts
	// in a 1500ms budget.
	require.True(t, strings.Contains(msg, "attempts") && !strings.Contains(msg, "1 attempts"),
		"expected more than 1 attempt within a 1500ms budget; got: %s", msg)

	// Honor the context budget rather than overshoot by orders of magnitude.
	// (The last backoff might extend us slightly past the deadline because
	// the timer fires before ctx.Done is checked; allow a small grace window.)
	require.Less(t, elapsed, 4*time.Second,
		"Connect ran %s, much longer than the 1500ms context budget — backoff may not respect ctx", elapsed)
}

// TestConnect_HonorsCancelDuringBackoff covers the case where the caller
// cancels the context while we're sleeping between ping attempts. The
// previous single-shot Connect always returned on the immediate ping failure;
// the retry version must NOT block past a cancel mid-backoff.
func TestConnect_HonorsCancelDuringBackoff(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := l.Addr().(*net.TCPAddr)
	require.NoError(t, l.Close())

	connStr := "postgres://x:x@" + addr.String() + "/db_x?sslmode=disable"

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel after ~600ms — long enough for the first ping to fail and the
	// retry loop to enter a backoff sleep, but well before the natural
	// expiration of a multi-second backoff.
	go func() {
		time.Sleep(600 * time.Millisecond)
		cancel()
	}()
	defer cancel()

	start := time.Now()
	_, err = Connect(ctx, connStr)
	elapsed := time.Since(start)

	require.Error(t, err)
	require.True(t, elapsed < 2*time.Second,
		"Connect should return within ~600ms of cancel; took %s", elapsed)
	require.Contains(t, err.Error(), "cancel", "error must indicate cancellation; got: %s", err.Error())
}
