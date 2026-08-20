package batcher

import (
	"context"
	"log/slog"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/ethrpc"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

// BatcherConfig configures a BatcherService.
type BatcherConfig struct {
	Identity  string        // scope of the DB queries (e.g. "privatehub")
	BatchSize int           // rows claimed per tick
	Interval  time.Duration // polling period
}

type BatcherService struct {
	identity  string
	batchSize int
	interval  time.Duration

	repo      Repository
	batcher   Batcher
	publisher ResultPublisher
}

func NewBatcherService(
	conf BatcherConfig, repo Repository, b Batcher, pub ResultPublisher,
) *BatcherService {
	return &BatcherService{
		identity:  conf.Identity,
		batchSize: conf.BatchSize,
		interval:  conf.Interval,
		repo:      repo,
		batcher:   b,
		publisher: pub,
	}
}

func (s *BatcherService) Run(ctx context.Context) error {
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
		s.send(ctx)
	}
}

func (s *BatcherService) send(ctx context.Context) {
	rows, err := s.repo.ClaimPending(ctx, s.identity, s.batchSize)
	if err != nil {
		slog.Error("send: claim pending failed", slog.String("identity", s.identity), slog.Any("err", err))
		return
	}
	if len(rows) == 0 {
		return
	}

	inputs := make([]ethrpc.TransactionInput, len(rows))
	keys := make([]TxKey, len(rows))
	for i, r := range rows {
		inputs[i] = ethrpc.TransactionInput{ID: r.CorrelationID, Data: r.Calldata, Address: r.Address}
		keys[i] = KeyOf(r)
	}

	// Bump attempts BEFORE broadcast so a mid-batch crash does not
	// leave attempts=0 on the retried rows.
	if err := s.repo.IncrementSendAttempts(ctx, keys); err != nil {
		slog.Error("send: bump attempts failed", slog.String("identity", s.identity), slog.Any("err", err))
	}

	slog.Info("Sending messages to ledger", slog.String("service", s.identity), slog.Int("count", len(inputs)))
	results, err := s.batcher.Send(ctx, inputs)
	if err != nil {
		// Whole-batch retryable failure — rows stay 'pending', next tick retries.
		slog.Error("send: batch send failed", slog.String("identity", s.identity), slog.Any("err", err))
		return
	}

	for i, res := range results {
		row := rows[i]
		id := row.CorrelationID
		switch {
		case res.Err != nil:
			slog.Warn("send: per-tx err",
				slog.String("identity", s.identity),
				slog.String("correlation_id", id),
				slog.Any("err", res.Err),
			)
		case res.Revert != nil:
			// Pre-flight revert: the tx reverted during eth_sendRawTransaction
			// simulation and never made it on-chain. Publish a TxResultRevert
			// before marking finished, otherwise the relayer's result router
			// never sees this terminal outcome and the originating service
			// callback never fires. Publish-before-mark mirrors the receipter
			// and reaper paths.
			result := types.TxResult{
				CorrelationID: id,
				MessageType:   row.MessageType,
				Identity:      s.identity,
				Kind:          types.TxResultRevert,
				RevertData:    res.Revert.Data,
			}
			if err := s.publisher.Push(ctx, result); err != nil {
				slog.Error("send: publish revert failed",
					slog.String("correlation_id", id), slog.Any("err", err))
				continue
			}
			if err := s.repo.MarkFinishedRevert(ctx, KeyOf(row), res.Revert.Data); err != nil {
				slog.Error("send: mark finished revert failed",
					slog.String("correlation_id", id), slog.Any("err", err))
			}
		default:
			if err := s.repo.MarkSent(ctx, KeyOf(row), res.Hash); err != nil {
				slog.Error("send: mark sent failed",
					slog.String("correlation_id", id), slog.Any("err", err))
			}
		}
	}
}
