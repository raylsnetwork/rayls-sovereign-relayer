// Decommissioning Teleport (vanilla, atomic).

package poller

import (
	"context"
	"log/slog"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

type expiredTransactionRepository interface {
	GetByStateOutcomeAndAtomicity(
		ctx context.Context,
		state types.TransactionState,
		outcome types.TransactionOutcome,
		isAtomic bool,
		opts ...repository.Option,
	) ([]types.Transaction, error)
}

type expiredExpiredService interface {
	Handle(context.Context, []string) error
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type ExpiredPoller struct {
	batchSize int

	svc    expiredExpiredService
	txRepo expiredTransactionRepository

	expirationPeriod time.Duration
	ticker           *time.Ticker
	initialRun       chan struct{}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func NewExpiredPoller(
	batchSize int,
	expirationPeriod time.Duration,
	svc expiredExpiredService,
	txRepo expiredTransactionRepository,
) *ExpiredPoller {
	initialRun := make(chan struct{}, 1)
	initialRun <- struct{}{}

	return &ExpiredPoller{
		batchSize: batchSize,

		svc:    svc,
		txRepo: txRepo,

		expirationPeriod: expirationPeriod,
		ticker:           time.NewTicker(time.Second),
		initialRun:       initialRun,
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (p *ExpiredPoller) Run(ctx context.Context) error {
	defer p.ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-p.ticker.C:
		case <-p.initialRun:
		}

		var sharedIDsToRevert []string

		txs, err := p.txRepo.GetByStateOutcomeAndAtomicity(
			ctx,
			types.SourcePublish,
			types.OutcomeSuccess,
			true,
			repository.WithLimit(p.batchSize),
		)
		if err != nil {
			slog.Error(
				"AtomicExpired: error while trying to get transactions by state and atomicity",
				slog.Any("error", err),
			)
			continue
		}
		if len(txs) == 0 {
			continue
		}

		for _, tx := range txs {
			if tx.UpdatedAt.Add(p.expirationPeriod).Before(time.Now()) {
				sharedIDsToRevert = append(sharedIDsToRevert, tx.SharedID)
			}
		}
		if len(sharedIDsToRevert) == 0 {
			continue
		}

		slog.Info("Reverting expired transactions", slog.Int("COUNT", len(sharedIDsToRevert)))
		err = p.svc.Handle(ctx, sharedIDsToRevert)
		if err != nil {
			slog.Error("AtomicExpired: error while trying to revert expired transactions", slog.Any("error", err))
			continue
		}
	}
}
