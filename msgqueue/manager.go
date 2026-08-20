package msgqueue

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

type Manager struct {
	streamName  string
	subjectRoot string
	chainId     string

	js jetstream.JetStream
}

func NewManager(ctx context.Context, js jetstream.JetStream, chainId string) (*Manager, error) {
	var (
		streamName  = "EVENTS"
		subjectRoot = "events." + chainId
	)
	cfg := jetstream.StreamConfig{
		Name:       streamName,
		Retention:  jetstream.WorkQueuePolicy,
		Subjects:   []string{"events.>"},
		Duplicates: time.Hour,
		// Enable atomic batch publishing (NATS 2.12+)
		AllowDirect: true,
	}

	_, err := js.CreateOrUpdateStream(ctx, cfg)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to create EVENTS stream: %w", err))
	}

	return &Manager{
		streamName:  streamName,
		subjectRoot: subjectRoot,
		chainId:     chainId,

		js: js,
	}, nil
}

func NewPublisher[T MessageWithID](m *Manager, subject string) *Publisher[T] {
	return newPublisher[T](constructSubject(m.subjectRoot, subject), m.js)
}

func NewConsumer[T MessageWithID](ctx context.Context, m *Manager, group string, subject string) (*Consumer[T], error) {
	// Prefix consumer group with chainId to ensure uniqueness across relayers
	durableName := m.chainId + "_" + group
	cons, err := m.js.CreateConsumer(ctx, m.streamName, jetstream.ConsumerConfig{
		Durable:       durableName,
		MaxDeliver:    10,
		FilterSubject: constructSubject(m.subjectRoot, subject),
	})
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to create jetstream consumer: %w", err))
	}
	return newConsumer[T](cons)
}

func constructSubject(stream string, subject string) string {
	return stream + "." + subject
}
