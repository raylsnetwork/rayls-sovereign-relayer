package msgqueue

import (
	"context"
	"errors"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testMsg is a minimal MessageWithID for exercising the publisher.
type testMsg struct{ id string }

func (m testMsg) GetID() string { return m.id }

// fakeJSPublisher records every PublishMsg call and can be made to fail on a
// specific (1-based) call index. It counts Nats-Batch-Commit headers, which the
// publisher sets once per atomic sub-batch — i.e. once per chunk.
type fakeJSPublisher struct {
	calls   int
	commits int
	failOn  int // 1-based call index to fail on; 0 = never fail
}

func (f *fakeJSPublisher) PublishMsg(_ context.Context, msg *nats.Msg, _ ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	f.calls++
	if msg.Header.Get("Nats-Batch-Commit") == "1" {
		f.commits++
	}
	if f.failOn != 0 && f.calls == f.failOn {
		return nil, errors.New("boom")
	}
	return &jetstream.PubAck{}, nil
}

func makeMsgs(n int) []testMsg {
	msgs := make([]testMsg, n)
	for i := range msgs {
		msgs[i] = testMsg{id: string(rune('a')) + string(rune(i%26))}
	}
	return msgs
}

func TestPushBatchChunked(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		n           int
		wantPublish int // expected PublishMsg calls
		wantChunks  int // expected atomic sub-batches (commit headers)
	}{
		{name: "empty is a no-op", n: 0, wantPublish: 0, wantChunks: 0},
		{name: "well under the cap", n: 10, wantPublish: 10, wantChunks: 1},
		{name: "exactly at the cap is one chunk", n: maxBatchSize, wantPublish: maxBatchSize, wantChunks: 1},
		{name: "one over the cap splits into two", n: maxBatchSize + 1, wantPublish: maxBatchSize + 1, wantChunks: 2},
		{name: "multiple full chunks plus remainder", n: 2*maxBatchSize + 500, wantPublish: 2*maxBatchSize + 500, wantChunks: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fake := &fakeJSPublisher{}
			pub := newPublisher[testMsg]("events.test", fake)

			err := pub.PushBatchChunked(context.Background(), makeMsgs(tt.n))

			require.NoError(t, err)
			assert.Equal(t, tt.wantPublish, fake.calls, "PublishMsg call count")
			assert.Equal(t, tt.wantChunks, fake.commits, "number of atomic sub-batches")
		})
	}
}

// The regression this fix targets: a >maxBatchSize batch that PushBatch rejects
// outright is published successfully when chunked.
func TestPushBatchChunked_RegressionOverCap(t *testing.T) {
	t.Parallel()
	msgs := makeMsgs(maxBatchSize + 1)

	rejected := newPublisher[testMsg]("events.test", &fakeJSPublisher{})
	require.ErrorIs(t, rejected.PushBatch(context.Background(), msgs), ErrTooManyMessages,
		"sanity: plain PushBatch must still reject an over-cap batch")

	chunked := newPublisher[testMsg]("events.test", &fakeJSPublisher{})
	require.NoError(t, chunked.PushBatchChunked(context.Background(), msgs),
		"PushBatchChunked must accept an over-cap batch")
}

// On a mid-stream failure, PushBatchChunked stops and reports the offending chunk
// range; earlier chunks have already been committed (acceptable because consumers
// dedup on redelivery).
func TestPushBatchChunked_FailsMidStream(t *testing.T) {
	t.Parallel()
	// Fail inside the second chunk (call 1200 of a 1500-message, 2-chunk push).
	fake := &fakeJSPublisher{failOn: maxBatchSize + 200}
	pub := newPublisher[testMsg]("events.test", fake)

	err := pub.PushBatchChunked(context.Background(), makeMsgs(maxBatchSize+500))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "chunk [1000:1500)", "error should name the failing chunk range")
	assert.Equal(t, 1, fake.commits, "only the first chunk should have committed")
}
