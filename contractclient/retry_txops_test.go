package contractclient_test

import (
	"context"
	"testing"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	txopspb "github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/txops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// zeroBackoff is a backoff.Strategy that never sleeps, keeping tests instant.
type zeroBackoff struct{}

func (zeroBackoff) Next(int) time.Duration { return 0 }

func (zeroBackoff) Do(_ context.Context, _ int, _ func() error) error { return nil }

// scriptedTxOpsClient returns the next scripted error for each method call (a
// nil entry, or running past the end of the slice, means success) and counts
// how many times each method was invoked.
type scriptedTxOpsClient struct {
	signErrs []error
	callErrs []error

	signCalls int
	callCalls int
}

func (s *scriptedTxOpsClient) SignAndSend(_ context.Context, _ *txopspb.SignAndSendRequest, _ ...grpc.CallOption) (*txopspb.SignAndSendResponse, error) {
	i := s.signCalls
	s.signCalls++
	if i < len(s.signErrs) && s.signErrs[i] != nil {
		return nil, s.signErrs[i]
	}
	return &txopspb.SignAndSendResponse{}, nil
}

func (s *scriptedTxOpsClient) BatchSignAndSend(_ context.Context, _ *txopspb.BatchSignAndSendRequest, _ ...grpc.CallOption) (*txopspb.BatchSignAndSendResponse, error) {
	return &txopspb.BatchSignAndSendResponse{}, nil
}

func (s *scriptedTxOpsClient) Call(_ context.Context, _ *txopspb.CallRequest, _ ...grpc.CallOption) (*txopspb.CallResponse, error) {
	i := s.callCalls
	s.callCalls++
	if i < len(s.callErrs) && s.callErrs[i] != nil {
		return nil, s.callErrs[i]
	}
	return &txopspb.CallResponse{}, nil
}

const testMaxAttempts = 3

func newTestClient(inner contractclient.TxOpsServiceClient) contractclient.TxOpsServiceClient {
	return contractclient.NewRetryingTxOpsClient(inner, zeroBackoff{}, testMaxAttempts)
}

func TestRetryingTxOpsClient_SignAndSend_WriteSemantics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		signErrs  []error
		wantCalls int
		wantErr   bool
	}{
		{
			name:      "succeeds first attempt",
			signErrs:  nil,
			wantCalls: 1,
		},
		{
			name:      "retries Unavailable then succeeds",
			signErrs:  []error{status.Error(codes.Unavailable, "blip"), nil},
			wantCalls: 2,
		},
		{
			name:      "retries DeadlineExceeded for writes (idempotency key makes resend safe)",
			signErrs:  []error{status.Error(codes.DeadlineExceeded, "ambiguous"), nil},
			wantCalls: 2,
		},
		{
			name:      "retries ResourceExhausted for writes",
			signErrs:  []error{status.Error(codes.ResourceExhausted, "backpressure"), nil},
			wantCalls: 2,
		},
		{
			name:      "does not retry permanent code",
			signErrs:  []error{status.Error(codes.InvalidArgument, "bad")},
			wantCalls: 1,
			wantErr:   true,
		},
		{
			name:      "exhausts attempts on persistent Unavailable",
			signErrs:  []error{status.Error(codes.Unavailable, "x"), status.Error(codes.Unavailable, "x"), status.Error(codes.Unavailable, "x")},
			wantCalls: testMaxAttempts,
			wantErr:   true,
		},
		{
			name:      "does not retry non-status error",
			signErrs:  []error{assertPlainError()},
			wantCalls: 1,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fake := &scriptedTxOpsClient{signErrs: tt.signErrs}
			client := newTestClient(fake)

			_, err := client.SignAndSend(context.Background(), &txopspb.SignAndSendRequest{})

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.wantCalls, fake.signCalls)
		})
	}
}

func TestRetryingTxOpsClient_Call_ReadSemantics(t *testing.T) {
	t.Parallel()

	// Call is read-only: even an ambiguous DeadlineExceeded is safe to retry.
	fake := &scriptedTxOpsClient{
		callErrs: []error{status.Error(codes.DeadlineExceeded, "slow"), nil},
	}
	client := newTestClient(fake)

	_, err := client.Call(context.Background(), &txopspb.CallRequest{})

	require.NoError(t, err)
	assert.Equal(t, 2, fake.callCalls)
}

func TestRetryingTxOpsClient_ContextCancelled(t *testing.T) {
	t.Parallel()

	// Persistent retryable failures, but a cancelled context must abort before
	// exhausting the attempt budget.
	fake := &scriptedTxOpsClient{
		callErrs: []error{
			status.Error(codes.Unavailable, "x"),
			status.Error(codes.Unavailable, "x"),
			status.Error(codes.Unavailable, "x"),
		},
	}
	client := newTestClient(fake)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.Call(ctx, &txopspb.CallRequest{})

	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

// assertPlainError returns a non-gRPC error to verify it is not retried.
func assertPlainError() error {
	return assertErr("not a grpc status")
}

type assertErr string

func (e assertErr) Error() string { return string(e) }
