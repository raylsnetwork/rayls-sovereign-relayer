// Decommissioning Teleport (vanilla, atomic).

package poller_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/poller"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StubFinalizationAtomicService struct {
	spyCalledSourceExecuted bool
	spySourceSharedIDs      []string

	spyCalledSourceReverted  bool
	spySourceRevertSharedIDs []string
}

func (s *StubFinalizationAtomicService) HandleSourceExecuted(ctx context.Context, sharedIDs []string) error {
	s.spyCalledSourceExecuted = true
	s.spySourceSharedIDs = sharedIDs
	return nil
}

func (s *StubFinalizationAtomicService) HandleSourceReverted(ctx context.Context, sharedIDs []string) error {
	s.spyCalledSourceReverted = true
	s.spySourceRevertSharedIDs = sharedIDs
	return nil
}

type StubFinalizationTransactionRepo struct {
	txs        []types.Transaction
	failedTxs  []types.Transaction
	expiredTxs []types.Transaction

	spyState    []types.TransactionState
	spyIsAtomic bool
}

func (r *StubFinalizationTransactionRepo) GetByStateOutcomeAndAtomicity(
	ctx context.Context,
	state types.TransactionState,
	outcome types.TransactionOutcome,
	isAtomic bool,
	opts ...repository.Option,
) ([]types.Transaction, error) {
	r.spyState = append(r.spyState, state)
	r.spyIsAtomic = isAtomic

	if state == types.SourceTimeoutRevert {
		return r.expiredTxs, nil
	}
	return r.txs, nil
}

func (r *StubFinalizationTransactionRepo) GetByStateOutcomesAndAtomicity(
	ctx context.Context,
	state types.TransactionState,
	outcomes []types.TransactionOutcome,
	isAtomic bool,
	opts ...repository.Option,
) ([]types.Transaction, error) {
	r.spyState = append(r.spyState, state)
	r.spyIsAtomic = isAtomic

	return r.failedTxs, nil
}

type StubFinalizationAtomicStatusRepo struct {
	spyCalledGetUnprocessed bool
	spySharedIDs            []string

	sums []types.AtomicStatusUpdateMessage
}

func (r *StubFinalizationAtomicStatusRepo) GetUnprocessedBySharedIDs(
	ctx context.Context,
	sharedIDs []string,
	opts ...repository.Option,
) ([]types.AtomicStatusUpdateMessage, error) {
	r.spyCalledGetUnprocessed = true
	r.spySharedIDs = sharedIDs

	var sums []types.AtomicStatusUpdateMessage
	for _, sum := range r.sums {
		if slices.Contains(sharedIDs, sum.SharedID) {
			sums = append(sums, sum)
		}
	}

	return sums, nil
}

func TestFinalizationPoller(t *testing.T) {
	batchSize := 100

	executeSharedID := "executed-shared-id"
	revertSharedID := "destination-reverted-shared-id"
	expiredSharedID := "source-reverted-shared-id"

	executedTx := types.Transaction{
		SharedID: executeSharedID,
	}
	failedTx := types.Transaction{
		SharedID: revertSharedID,
	}
	expiredTx := types.Transaction{
		SharedID: expiredSharedID,
	}

	executedSUM := types.AtomicStatusUpdateMessage{
		SharedID: executeSharedID,
		Status:   types.AtomicExecutedStatus,
	}
	failedSUM := types.AtomicStatusUpdateMessage{
		SharedID: revertSharedID,
		Status:   types.AtomicRevertedStatus,
	}
	expiredSUM := types.AtomicStatusUpdateMessage{
		SharedID: expiredSharedID,
		Status:   types.AtomicRevertedStatus,
	}

	t.Run("exits gracefully on canceled context", func(t *testing.T) {
		sumsRepo := &StubFinalizationAtomicStatusRepo{}
		txRepo := &StubFinalizationTransactionRepo{}
		svc := &StubFinalizationAtomicService{}
		pol := poller.NewFinalizationPoller(batchSize, svc, txRepo, sumsRepo)

		hasGracefulShutdown := testtools.ShutdownFixture(t, pol.Run, time.Millisecond)

		require.True(t, hasGracefulShutdown, "poller doesn't support graceful shutdown")
	})

	t.Run(
		"polls for transactions in SourcePublish (success/reverted/failed) and SourceTimeoutRevert states",
		func(t *testing.T) {
			txs := []types.Transaction{
				executedTx,
			}
			failedTxs := []types.Transaction{
				failedTx,
			}
			sums := []types.AtomicStatusUpdateMessage{
				executedSUM,
			}

			txRepo := &StubFinalizationTransactionRepo{
				txs:       txs,
				failedTxs: failedTxs,
			}
			sumsRepo := &StubFinalizationAtomicStatusRepo{
				sums: sums,
			}
			svc := &StubFinalizationAtomicService{}
			pol := poller.NewFinalizationPoller(batchSize, svc, txRepo, sumsRepo)

			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()

			_ = pol.Run(ctx)

			assert.Contains(t, txRepo.spyState, types.SourcePublish)
			assert.Contains(t, txRepo.spyState, types.SourceTimeoutRevert)
			assert.Equal(t, true, txRepo.spyIsAtomic)
		},
	)

	t.Run("routes executed transactions to HandleSourceExecuted", func(t *testing.T) {
		txs := []types.Transaction{
			executedTx,
		}
		failedTxs := []types.Transaction{
			failedTx,
		}
		sums := []types.AtomicStatusUpdateMessage{
			executedSUM,
			failedSUM,
		}

		txRepo := &StubFinalizationTransactionRepo{
			txs:       txs,
			failedTxs: failedTxs,
		}
		sumsRepo := &StubFinalizationAtomicStatusRepo{
			sums: sums,
		}
		svc := &StubFinalizationAtomicService{}
		pol := poller.NewFinalizationPoller(batchSize, svc, txRepo, sumsRepo)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = pol.Run(ctx)

		assert.Equal(t, []string{executeSharedID}, svc.spySourceSharedIDs)
	})

	t.Run("routes destination reverted transactions to HandleSourceReverted", func(t *testing.T) {
		txs := []types.Transaction{
			executedTx,
			failedTx,
		}
		failedTxs := []types.Transaction{
			failedTx,
		}
		sums := []types.AtomicStatusUpdateMessage{
			executedSUM,
			failedSUM,
		}

		txRepo := &StubFinalizationTransactionRepo{
			txs:       txs,
			failedTxs: failedTxs,
		}
		sumsRepo := &StubFinalizationAtomicStatusRepo{
			sums: sums,
		}
		svc := &StubFinalizationAtomicService{}
		pol := poller.NewFinalizationPoller(batchSize, svc, txRepo, sumsRepo)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := pol.Run(ctx)
		require.Nil(t, err)

		assert.Equal(t, []string{revertSharedID}, svc.spySourceRevertSharedIDs)
	})

	t.Run("routes source reverted (expired) transactions to HandleSourceReverted", func(t *testing.T) {
		expiredTxs := []types.Transaction{
			expiredTx,
		}
		sums := []types.AtomicStatusUpdateMessage{
			expiredSUM,
		}

		txRepo := &StubFinalizationTransactionRepo{
			expiredTxs: expiredTxs,
		}
		sumsRepo := &StubFinalizationAtomicStatusRepo{
			sums: sums,
		}
		svc := &StubFinalizationAtomicService{}
		pol := poller.NewFinalizationPoller(batchSize, svc, txRepo, sumsRepo)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := pol.Run(ctx)
		require.Nil(t, err)

		assert.Equal(t, []string{expiredSharedID}, svc.spySourceRevertSharedIDs)
	})

	t.Run("doesn't call execute on no transactions to process", func(t *testing.T) {
		failedTxs := []types.Transaction{
			failedTx,
		}
		expiredTxs := []types.Transaction{
			expiredTx,
		}
		sums := []types.AtomicStatusUpdateMessage{
			failedSUM,
			executedSUM,
		}

		txRepo := &StubFinalizationTransactionRepo{
			failedTxs:  failedTxs,
			expiredTxs: expiredTxs,
		}
		sumsRepo := &StubFinalizationAtomicStatusRepo{
			sums: sums,
		}
		svc := &StubFinalizationAtomicService{}
		pol := poller.NewFinalizationPoller(batchSize, svc, txRepo, sumsRepo)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := pol.Run(ctx)
		require.Nil(t, err)

		assert.False(t, svc.spyCalledSourceExecuted, "called source execute when shouldn't have")
	})

	t.Run("doesn't call revert on no transactions to process", func(t *testing.T) {
		txs := []types.Transaction{
			executedTx,
		}
		sums := []types.AtomicStatusUpdateMessage{
			executedSUM,
		}

		txRepo := &StubFinalizationTransactionRepo{
			txs: txs,
		}
		sumsRepo := &StubFinalizationAtomicStatusRepo{
			sums: sums,
		}
		svc := &StubFinalizationAtomicService{}
		pol := poller.NewFinalizationPoller(batchSize, svc, txRepo, sumsRepo)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := pol.Run(ctx)
		require.Nil(t, err)

		assert.False(t, svc.spyCalledSourceReverted, "called source revert when shouldn't have")
	})

	t.Run("doesn't query sums repository on no transactions to process", func(t *testing.T) {
		txs := []types.Transaction{}
		sums := []types.AtomicStatusUpdateMessage{
			executedSUM,
			failedSUM,
		}

		txRepo := &StubFinalizationTransactionRepo{
			txs: txs,
		}
		sumsRepo := &StubFinalizationAtomicStatusRepo{
			sums: sums,
		}
		svc := &StubFinalizationAtomicService{}
		pol := poller.NewFinalizationPoller(batchSize, svc, txRepo, sumsRepo)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := pol.Run(ctx)
		require.Nil(t, err)

		assert.False(t, sumsRepo.spyCalledGetUnprocessed, "queried sums repository but shouldn't have")
	})
}
