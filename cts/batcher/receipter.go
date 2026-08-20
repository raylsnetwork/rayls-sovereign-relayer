package batcher

import (
	"context"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/faultinjector"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

// ReceipterConfig configures a ReceipterService.
type ReceipterConfig struct {
	Identity  string
	BatchSize int
	Interval  time.Duration
}

// ReceipterService polls sent rows for one identity, fetches their
// receipts in batches, and publishes terminal results.
//
// Publish-before-mark ordering: a crash or publish failure between
// publish and MarkFinished leaves the row 'sent'; the next tick re-polls,
// re-publishes with the same msg_id, and JetStream dedup absorbs the
// duplicate. The relayer-side handler must be idempotent (required for
// NATS at-least-once semantics regardless).
type ReceipterService struct {
	identity  string
	batchSize int
	interval  time.Duration

	repo      Repository
	batcher   Batcher
	publisher ResultPublisher
}

func NewReceipterService(
	conf ReceipterConfig, repo Repository, b Batcher, pub ResultPublisher,
) *ReceipterService {
	return &ReceipterService{
		identity:  conf.Identity,
		batchSize: conf.BatchSize,
		interval:  conf.Interval,
		repo:      repo,
		batcher:   b,
		publisher: pub,
	}
}

func (s *ReceipterService) Run(ctx context.Context) error {
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
		s.receipt(ctx)
	}
}

// receipt polls receipts for a batch of sent rows. Per-row outcome:
//   - Receipt  → publish success, then mark finished
//   - Revert   → publish revert, then mark finished
//   - Err      → transient; log and leave 'sent'
func (s *ReceipterService) receipt(ctx context.Context) {
	// Fault point: simulate the chain "losing" a broadcast tx so the receipter
	// never observes its receipt. Arming `error` here makes the receipter skip
	// this poll cycle (early return); arming `error` with a `count` > 1 keeps
	// skipping for N consecutive cycles, letting the reaper's StuckThreshold
	// elapse and dead-letter the row. Combined with the e2e regression test
	// at test/e2e/security/resilience/Vanilla_StuckTx_PermanentLoss.ts, this
	// reproduces the "axyl dropped a tx during consensus stall" failure mode
	// without requiring real chain instability.
	if faultErr := faultinjector.Check("cts.batcher.ReceipterService.receipt.before_poll"); faultErr != nil {
		slog.Warn("receipt: fault injection — skipping poll cycle",
			slog.String("identity", s.identity), slog.Any("err", faultErr))
		return
	}

	rows, err := s.repo.ClaimSent(ctx, s.identity, s.batchSize)
	if err != nil {
		slog.Error("receipt: claim sent failed", slog.String("identity", s.identity), slog.Any("err", err))
		return
	}
	if len(rows) == 0 {
		return
	}

	hashes := make([]common.Hash, len(rows))
	keys := make([]TxKey, len(rows))
	for i, r := range rows {
		hashes[i] = r.TxHash
		keys[i] = KeyOf(r)
	}

	if err := s.repo.IncrementReceiptAttempts(ctx, keys); err != nil {
		slog.Error("receipt: bump attempts failed", slog.String("identity", s.identity), slog.Any("err", err))
	}

	slog.Info("Querying receipts from ledger", slog.String("service", s.identity), slog.Int("count", len(hashes)))
	results, err := s.batcher.GetReceipts(ctx, hashes)
	if err != nil {
		slog.Error("receipt: batch receipts failed", slog.String("identity", s.identity), slog.Any("err", err))
		return
	}

	slog.Info("Publishing results to consumers", slog.String("service", s.identity), slog.Int("count", len(results)))
	for i, res := range results {
		row := rows[i]
		switch {
		case res.Pending:
			continue
		case res.Err != nil:
			slog.Warn("receipt: per-tx err",
				slog.String("identity", s.identity),
				slog.String("correlation_id", row.CorrelationID),
				slog.Any("err", res.Err),
			)
		case res.Receipt != nil:
			result := types.TxResult{
				CorrelationID: row.CorrelationID,
				MessageType:   row.MessageType,
				Identity:      s.identity,
				Kind:          types.TxResultSuccess,
				TxHash:        row.TxHash,
				Receipt:       res.Receipt,
			}
			if err := s.publisher.Push(ctx, result); err != nil {
				slog.Error("receipt: publish success failed",
					slog.String("correlation_id", row.CorrelationID), slog.Any("err", err))
				continue
			}
			if err := s.repo.MarkFinishedSuccess(ctx, KeyOf(row), 1); err != nil {
				slog.Error("receipt: mark finished success failed",
					slog.String("correlation_id", row.CorrelationID), slog.Any("err", err))
			}
		case res.Revert != nil:
			result := types.TxResult{
				CorrelationID: row.CorrelationID,
				MessageType:   row.MessageType,
				Identity:      s.identity,
				Kind:          types.TxResultRevert,
				TxHash:        row.TxHash,
				RevertData:    res.Revert.Data,
			}
			if err := s.publisher.Push(ctx, result); err != nil {
				slog.Error("receipt: publish revert failed",
					slog.String("correlation_id", row.CorrelationID), slog.Any("err", err))
				continue
			}
			if err := s.repo.MarkFinishedRevert(ctx, KeyOf(row), res.Revert.Data); err != nil {
				slog.Error("receipt: mark finished revert failed",
					slog.String("correlation_id", row.CorrelationID), slog.Any("err", err))
			}
		}
	}
}
