package ethretry

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestIsUnmarshalTypeError(t *testing.T) {
	t.Parallel()

	wrapped := &json.UnmarshalTypeError{Value: "number", Type: reflect.TypeOf("")}
	if !IsUnmarshalTypeError(wrapped) {
		t.Fatal("expected *json.UnmarshalTypeError to be classified as unmarshal type error")
	}
	if !IsUnmarshalTypeError(errors.Join(errors.New("wrap"), wrapped)) {
		t.Fatal("expected wrapped *json.UnmarshalTypeError to be classified as unmarshal type error")
	}
	if IsUnmarshalTypeError(errors.New("plain")) {
		t.Fatal("plain error should not classify as unmarshal type error")
	}
}

func TestIsConnectionError(t *testing.T) {
	t.Parallel()

	if !IsConnectionError(&net.OpError{Op: "dial", Err: errors.New("refused")}) {
		t.Fatal("expected *net.OpError to be classified as connection error")
	}
	if !IsConnectionError(&url.Error{Op: "Get", URL: "http://x", Err: errors.New("x")}) {
		t.Fatal("expected *url.Error to be classified as connection error")
	}
	if IsConnectionError(errors.New("plain")) {
		t.Fatal("plain error should not classify as connection error")
	}
}

func TestWithRetry_SuccessReturnsImmediately(t *testing.T) {
	t.Parallel()

	calls := 0
	err := WithRetry(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestWithRetry_NonTransientErrorReturnsImmediately(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("nope")
	calls := 0
	err := WithRetry(context.Background(), func() error {
		calls++
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call (no retry), got %d", calls)
	}
}

// TestWithRetry_ContextCancellationDuringRetry verifies that a retryable
// error followed by a cancelled context surfaces the underlying retryable
// error (not ctx.Err()) — the contract preserved across every move of this
// helper (out of batcher.go, then into the ethretry package).
func TestWithRetry_ContextCancellationDuringRetry(t *testing.T) {
	t.Parallel()

	transient := &json.UnmarshalTypeError{Value: "number", Type: reflect.TypeOf("")}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := WithRetry(ctx, func() error {
		return transient
	})
	if !errors.Is(err, transient) {
		t.Fatalf("expected to receive the transient error, got %v", err)
	}
}

func TestNextRetryDelay(t *testing.T) {
	t.Parallel()

	const incr = time.Second
	const maxDelay = 30 * time.Second

	cases := []struct {
		name          string
		current, want time.Duration
	}{
		{"first increment", 3 * time.Second, 4 * time.Second},
		{"just below cap", 29 * time.Second, 30 * time.Second},
		{"hits cap exactly", 30 * time.Second, 30 * time.Second},
		{"saturated, stays at cap", 100 * time.Second, 30 * time.Second},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := nextRetryDelay(tc.current, incr, maxDelay); got != tc.want {
				t.Errorf("nextRetryDelay(%v, %v, %v) = %v, want %v",
					tc.current, incr, maxDelay, got, tc.want)
			}
		})
	}
}
