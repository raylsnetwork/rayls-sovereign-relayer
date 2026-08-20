// Decommissioning Teleport (vanilla, atomic).

package poller

import (
	"context"
	"log/slog"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

type earlyRevertTransactionRepository interface {
	GetByStateOutcomesAndAtomicity(
		ctx context.Context,
		state types.TransactionState,
		outcomes []types.TransactionOutcome,
		isAtomic bool,
		opts ...repository.Option,
	) ([]types.Transaction, error)
}

type earlyRevertService interface {
	HandleEarlyRevert(context.Context, []string) error
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type EarlyRevertPoller struct {
	batchSize int

	svc    earlyRevertService
	txRepo earlyRevertTransactionRepository

	ticker     *time.Ticker
	initialRun chan struct{}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func NewEarlyRevertPoller(
	batchSize int,
	svc earlyRevertService,
	txRepo earlyRevertTransactionRepository,
) *EarlyRevertPoller {
	initialRun := make(chan struct{}, 1)
	initialRun <- struct{}{}

	return &EarlyRevertPoller{
		batchSize: batchSize,

		svc:    svc,
		txRepo: txRepo,

		ticker:     time.NewTicker(time.Second),
		initialRun: initialRun,
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (p *EarlyRevertPoller) Run(ctx context.Context) error {
	defer p.ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-p.ticker.C:
		case <-p.initialRun:
		}

		var sharedIDs []string

		// Both reverted (mined-but-reverted on PNH) and failed (never mined) outcomes
		// trigger early-revert: in either case the source-side escrow needs releasing.
		txs, err := p.txRepo.GetByStateOutcomesAndAtomicity(
			ctx,
			types.SourcePublish,
			[]types.TransactionOutcome{types.OutcomeReverted, types.OutcomeFailed},
			true,
			repository.WithLimit(p.batchSize),
		)
		if err != nil {
			slog.Error(
				"EarlyRevert: error while trying to get transactions by state and atomicity",
				slog.Any("error", err),
			)
			continue
		}
		if len(txs) == 0 {
			continue
		}

		for _, tx := range txs {
			sharedIDs = append(sharedIDs, tx.SharedID)
		}

		slog.Info("Reverting failed transactions", slog.Int("COUNT", len(sharedIDs)))
		err = p.svc.HandleEarlyRevert(ctx, sharedIDs)
		if err != nil {
			slog.Error("EarlyRevert: error while trying to revert failed transactions", slog.Any("error", err))
			continue
		}
	}
}
