package batcher

import (
	"context"
	"log/slog"
	"time"
)

type IngestionConfig struct {
	Identity  string
	BatchSize int
	Interval  time.Duration
}

type IngestionService struct {
	identity  string
	batchSize int
	interval  time.Duration

	consumer MessageConsumer
	repo     Repository
}

// cts.send.private_hub
// cts.send.private_node
// cts.send.public_chain
// cts.send. -> [id, message_type, calldata, address], [id, calldata, address], [id, calldata, address], [id, calldata, address]
// cts.result. -> [id, receipt, revert_data] [id, receipt, revert_data] [id, receipt, revert_data]
// cts.result.public_chain -> [id, receipt, public.foward] [id, receipt, public.revert] [id, receipt, public.foward] [id, receipt, public.foward]
// cts.result.private_hub

func NewIngestionService(conf IngestionConfig, consumer MessageConsumer, repo Repository) *IngestionService {
	return &IngestionService{
		identity:  conf.Identity,
		batchSize: conf.BatchSize,
		interval:  conf.Interval,
		consumer:  consumer,
		repo:      repo,
	}
}

func (s *IngestionService) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	initial := make(chan struct{}, 1)
	initial <- struct{}{}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		case <-initial:
		}
		s.ingest(ctx)
	}
}

func (s *IngestionService) ingest(ctx context.Context) {
	msgs, err := s.consumer.Fetch(ctx, s.batchSize)
	if err != nil {
		slog.Error("ingestion fetch failed", slog.String("identity", s.identity), slog.Any("err", err))
		return
	}
	if len(msgs) == 0 {
		return
	}

	slog.Info("Ingesting messages into batching pipeline", slog.String("service", s.identity), slog.Int("count", len(msgs)))
	for _, m := range msgs {
		// Identity on the wire is redundant — the cts.send.<identity>
		// subject already implies it. Stamp it from the IngestionService
		// config so the row's `identity` column matches what the
		// downstream BatcherService.ClaimPending filters on. Without
		// this the rows land with identity='' and the sender finds
		// nothing.
		req := m.V
		req.Identity = s.identity
		if err := s.repo.Insert(ctx, req); err != nil {
			slog.Error("ingestion insert failed",
				slog.String("identity", s.identity),
				slog.String("correlation_id", m.V.CorrelationID),
				slog.Any("err", err),
			)
			continue
		}
		if err := m.Ack(ctx); err != nil {
			slog.Error("ingestion ack failed",
				slog.String("identity", s.identity),
				slog.String("correlation_id", m.V.CorrelationID),
				slog.Any("err", err),
			)
		}
	}
}
