package msgqueue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

type JetStreamConsumer interface {
	FetchNoWait(batch int) (jetstream.MessageBatch, error)
	Messages(opts ...jetstream.PullMessagesOpt) (jetstream.MessagesContext, error)
}

type Message[T MessageWithID] struct {
	V   T
	Ack func(context.Context) error
}

type Consumer[T MessageWithID] struct {
	jsCons JetStreamConsumer

	iter jetstream.MessagesContext
}

func newConsumer[T MessageWithID](jsCons JetStreamConsumer) (*Consumer[T], error) {
	iter, err := jsCons.Messages()
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to create NATS messages iterator: %w", err))
	}
	return &Consumer[T]{
		jsCons: jsCons,
		iter:   iter,
	}, nil
}

func (c *Consumer[T]) Fetch(ctx context.Context, count int) ([]Message[T], error) {
	batch, err := c.jsCons.FetchNoWait(count)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to fetch batch: %w", err))
	}

	defer logBatchError(batch)

	msgs := make([]Message[T], 0, count)

	var (
		ok    bool
		jsMsg jetstream.Msg
	)
	for {
		select {
		case <-ctx.Done():
			return msgs, ctx.Err()
		case jsMsg, ok = <-batch.Messages():
			// channel is closed -> we read all the messages
			if !ok {
				return msgs, nil
			}
		}

		obj, err := unmarshalInto[T](jsMsg.Data())
		if err != nil {
			slog.Error("Failed to unmarshall message", slog.Any("error", err))
			continue
		}

		msgs = append(msgs, Message[T]{
			V:   obj,
			Ack: jsMsg.DoubleAck,
		})

		if len(msgs) >= count {
			return msgs, nil
		}
	}
}

func (c *Consumer[T]) Next(ctx context.Context) (Message[T], error) {
	for {
		jsMsg, err := c.iter.Next(
			jetstream.NextContext(ctx),
		)
		if err != nil {
			// If it's the internal MaxWait timeout, just loop again.
			if errors.Is(err, context.DeadlineExceeded) {
				// no message during this poll; try again
				continue
			}

			// If ctx is done, surface that.
			if ctx.Err() != nil {
				return Message[T]{}, ctx.Err()
			}
			return Message[T]{}, withstack.Wrap(fmt.Errorf("next message failed: %w", err))
		}

		var empty Message[T]

		obj, err := unmarshalInto[T](jsMsg.Data())
		if err != nil {
			return empty, fmt.Errorf("failed to unmarshall message: %w", err)
		}

		return Message[T]{
			V:   obj,
			Ack: jsMsg.DoubleAck,
		}, nil
	}
}

func logBatchError(batch jetstream.MessageBatch) {
	if batch.Error() != nil {
		slog.Error("Batch of messages returned an error", slog.Any("error", batch.Error()))
	}
}

func unmarshalInto[T any](data []byte) (obj T, err error) {
	err = json.Unmarshal(data, &obj)
	return
}
