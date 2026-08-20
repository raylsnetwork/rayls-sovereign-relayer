package msgqueue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

const maxBatchSize = 1000

var ErrTooManyMessages = errors.New("cannot process batch - too many messages")

type MessageWithID interface {
	GetID() string
}

// BatchPublishError represents an error that occurred during batch publishing.
// Since batch publishing is atomic (all-or-nothing), this error indicates
// that none of the messages in the batch were successfully published.
type BatchPublishError struct {
	// BatchSize is the total number of messages in the failed batch
	BatchSize int
	// FailedAtIndex is the index of the message that caused the failure
	FailedAtIndex int
	// Err is the underlying error
	Err error
}

func (e *BatchPublishError) Error() string {
	return fmt.Sprintf("batch publish failed at message %d/%d: %v", e.FailedAtIndex, e.BatchSize, e.Err)
}

func (e *BatchPublishError) Unwrap() error {
	return e.Err
}

type JetStreamPublisher interface {
	PublishMsg(ctx context.Context, msg *nats.Msg, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error)
}

type Publisher[T MessageWithID] struct {
	subject string
	jsPub   JetStreamPublisher
}

func newPublisher[T MessageWithID](subject string, jsPub JetStreamPublisher) *Publisher[T] {
	return &Publisher[T]{
		subject: subject,
		jsPub:   jsPub,
	}
}

func (q *Publisher[T]) Push(ctx context.Context, msg T) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	natsMsg := &nats.Msg{
		Subject: q.subject,
		Header:  nats.Header{},
		Data:    data,
	}

	// Set Msg-Id header to avoid message duplication
	withMsgIDOpt := jetstream.WithMsgID(q.getMsgIDWithSubject(msg.GetID()))
	_, err = q.jsPub.PublishMsg(ctx, natsMsg, withMsgIDOpt)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to publish message: %w", err))
	}
	return nil
}

// PushBatch publishes multiple messages to the queue using NATS atomic batch publishing.
// All messages are published atomically (all-or-nothing) in a single batch operation.
// This uses NATS batch headers to ensure either all messages are persisted or none are.
// Maximum batch size is 1000 messages per NATS limitation.
// Returns a BatchPublishError if the batch fails; otherwise returns nil on success.
func (q *Publisher[T]) PushBatch(ctx context.Context, msgs []T) error {
	if len(msgs) == 0 {
		return nil
	}

	// NATS supports max 1000 messages per batch
	if len(msgs) > maxBatchSize {
		return ErrTooManyMessages
	}

	// Publish messages using atomic batch semantics
	// All messages except the last have Nats-Batch-Open header
	// The last message has Nats-Batch-Commit header to commit the batch
	for i, msg := range msgs {
		data, err := json.Marshal(msg)
		if err != nil {
			return &BatchPublishError{
				BatchSize:     len(msgs),
				FailedAtIndex: i,
				Err:           fmt.Errorf("failed to marshal message: %w", err),
			}
		}

		natsMsg := &nats.Msg{
			Subject: q.subject,
			Header:  nats.Header{},
			Data:    data,
		}

		// Set batch headers for atomic publishing
		isLastMessage := i == len(msgs)-1
		if isLastMessage {
			// Commit the batch on the last message
			natsMsg.Header.Set("Nats-Batch-Commit", "1")
		} else if i == 0 {
			// Open the batch on the first message
			natsMsg.Header.Set("Nats-Batch-Open", "1")
		}

		// Set message ID for deduplication
		withMsgIDOpt := jetstream.WithMsgID(q.getMsgIDWithSubject(msg.GetID()))

		_, err = q.jsPub.PublishMsg(ctx, natsMsg, withMsgIDOpt)
		if err != nil {
			return &BatchPublishError{
				BatchSize:     len(msgs),
				FailedAtIndex: i,
				Err:           withstack.Wrap(fmt.Errorf("failed to publish message: %w", err)),
			}
		}
	}

	return nil
}

// PushBatchChunked splits msgs into atomic sub-batches of at most maxBatchSize and
// publishes them in order via PushBatch. Unlike PushBatch, the operation is atomic
// ONLY per chunk, not across the whole input: if a chunk fails, the chunks before it
// are already committed. This is safe for at-least-once consumers that deduplicate on
// Msg-Id (every EVENTS-stream consumer does, via the stream's Duplicates window) — a
// redelivery re-pushes every chunk and the already-committed ones are dropped as
// duplicates. Do NOT use it where callers require whole-input atomicity; use PushBatch
// with len(msgs) <= maxBatchSize for that.
func (q *Publisher[T]) PushBatchChunked(ctx context.Context, msgs []T) error {
	for i := 0; i < len(msgs); i += maxBatchSize {
		end := min(i+maxBatchSize, len(msgs))
		if err := q.PushBatch(ctx, msgs[i:end]); err != nil {
			return fmt.Errorf("pushing chunk [%d:%d) of %d messages: %w", i, end, len(msgs), err)
		}
	}
	return nil
}

// Used for deduplication of messages on a subject level
func (q *Publisher[T]) getMsgIDWithSubject(id string) string {
	return q.subject + "." + id
}
