package fake

import (
	"context"
)

// NextMQ creates a mock message queue that returns a single message once via Next(),
// then blocks on context cancellation for subsequent calls.
//
// This is useful for testing services that consume messages one at a time.
//
// Usage:
//
//	nextFunc := fake.NextMQ(message)
//	mock := &SomeMQMock{
//	    NextFunc: nextFunc,
//	}
func NextMQ[T any](msg T) func(context.Context) (T, error) {
	firstCall := true
	return func(ctx context.Context) (T, error) {
		if firstCall {
			firstCall = false
			return msg, nil
		}
		<-ctx.Done()
		var zero T
		return zero, context.Canceled
	}
}

// FetchMQ creates a mock message queue that returns a slice of messages once via Fetch(),
// then blocks on context cancellation for subsequent calls.
//
// This is useful for testing services that fetch messages in batches.
//
// Usage:
//
//	fetchFunc := fake.FetchMQ([]message{msg1, msg2})
//	mock := &SomeMQMock{
//	    FetchFunc: fetchFunc,
//	}
func FetchMQ[T any](msgs []T) func(context.Context, int) ([]T, error) {
	firstCall := true
	return func(ctx context.Context, count int) ([]T, error) {
		if firstCall {
			firstCall = false
			return msgs, nil
		}
		<-ctx.Done()
		return nil, context.Canceled
	}
}
