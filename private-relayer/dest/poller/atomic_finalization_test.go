// Decommissioning Teleport (vanilla, atomic).

package poller_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/poller"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StubFinalizationSignatureService struct {
	spyCalledDestinationExecuted  bool
	spyDestinationUnlockSharedIDs []string

	spyCalledDestinationReverted bool
	spyrevertSharedIDs           []string
}

func (s *StubFinalizationSignatureService) HandleDestinationExecuted(ctx context.Context, sharedIDs []string) error {
	s.spyCalledDestinationExecuted = true
	s.spyDestinationUnlockSharedIDs = sharedIDs
	return nil
}

func (s *StubFinalizationSignatureService) HandleDestinationReverted(ctx context.Context, sharedIDs []string) error {
	s.spyCalledDestinationReverted = true
	s.spyrevertSharedIDs = sharedIDs
	return nil
}

type StubFinalizationTransactionRepo struct {
	txs []types.Transaction

	spyState    types.TransactionState
	spyIsAtomic bool
}

func (r *StubFinalizationTransactionRepo) GetByStateAndAtomicity(
	ctx context.Context,
	state types.TransactionState,
	isAtomic bool,
	opts ...repository.Option,
) ([]types.Transaction, error) {
	r.spyState = state
	r.spyIsAtomic = isAtomic

	return r.txs, nil
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
	testtools.SilenceLogger()

	batchSize := 100

	executeSharedID := "destination-execute-shared-id"
	revertSharedID := "destination-revert-shared-id"

	executedTx := types.Transaction{
		SharedID: executeSharedID,
		Outcome:  types.OutcomeSuccess,
	}
	failedTx := types.Transaction{
		SharedID: revertSharedID,
		Outcome:  types.OutcomeSuccess,
	}

	executedSUM := types.AtomicStatusUpdateMessage{
		SharedID: executeSharedID,
		Status:   types.AtomicExecutedStatus,
	}
	failedSUM := types.AtomicStatusUpdateMessage{
		SharedID: revertSharedID,
		Status:   types.AtomicRevertedStatus,
	}

	t.Run("exits gracefully on canceled context", func(t *testing.T) {
		sumsRepo := &StubFinalizationAtomicStatusRepo{}
		txRepo := &StubFinalizationTransactionRepo{}
		svc := &StubFinalizationSignatureService{}
		p := poller.NewFinalizationPoller(batchSize, svc, txRepo, sumsRepo)

		hasGracefulShutdown := testtools.ShutdownFixture(t, p.Run, time.Millisecond)

		require.True(t, hasGracefulShutdown, "poller doesn't support graceful shutdown")
	})

	t.Run("routes executed transactions to HandleDestinationExecuted", func(t *testing.T) {
		txs := []types.Transaction{
			executedTx,
			failedTx,
		}
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
		svc := &StubFinalizationSignatureService{}
		p := poller.NewFinalizationPoller(batchSize, svc, txRepo, sumsRepo)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = p.Run(ctx)

		assert.Equal(t, []string{executeSharedID}, svc.spyDestinationUnlockSharedIDs)
	})

	t.Run("routes revert transactions to HandleDestinationReverted", func(t *testing.T) {
		txs := []types.Transaction{
			executedTx,
			failedTx,
		}
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
		svc := &StubFinalizationSignatureService{}
		p := poller.NewFinalizationPoller(batchSize, svc, txRepo, sumsRepo)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = p.Run(ctx)

		assert.Equal(t, []string{revertSharedID}, svc.spyrevertSharedIDs)
	})

	t.Run("doesn't call execute on no transactions to process", func(t *testing.T) {
		txs := []types.Transaction{
			failedTx,
		}
		sums := []types.AtomicStatusUpdateMessage{
			failedSUM,
		}

		txRepo := &StubFinalizationTransactionRepo{
			txs: txs,
		}
		sumsRepo := &StubFinalizationAtomicStatusRepo{
			sums: sums,
		}
		svc := &StubFinalizationSignatureService{}
		p := poller.NewFinalizationPoller(batchSize, svc, txRepo, sumsRepo)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := p.Run(ctx)
		require.Nil(t, err)

		assert.False(t, svc.spyCalledDestinationExecuted, "called destination execute when shouldn't have")
	})

	t.Run("doesn't call destination revert on no transactions to process", func(t *testing.T) {
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
		svc := &StubFinalizationSignatureService{}
		p := poller.NewFinalizationPoller(batchSize, svc, txRepo, sumsRepo)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := p.Run(ctx)
		require.Nil(t, err)

		assert.False(t, svc.spyCalledDestinationReverted, "called destination revert when shouldn't have")
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
		svc := &StubFinalizationSignatureService{}
		p := poller.NewFinalizationPoller(batchSize, svc, txRepo, sumsRepo)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := p.Run(ctx)
		require.Nil(t, err)

		assert.False(t, sumsRepo.spyCalledGetUnprocessed, "queried sums repository but shouldn't have")
	})

	t.Run("routes early-race rows (HubNotifiedExec+reverted) directly to HandleDestinationReverted", func(t *testing.T) {
		earlyRaceSharedID := "early-race-shared-id"
		earlyRaceTx := types.Transaction{
			SharedID: earlyRaceSharedID,
			Outcome:  types.OutcomeReverted,
		}

		txRepo := &StubFinalizationTransactionRepo{
			txs: []types.Transaction{earlyRaceTx},
		}
		sumsRepo := &StubFinalizationAtomicStatusRepo{}
		svc := &StubFinalizationSignatureService{}
		p := poller.NewFinalizationPoller(batchSize, svc, txRepo, sumsRepo)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = p.Run(ctx)

		assert.False(t, sumsRepo.spyCalledGetUnprocessed,
			"early-race rows should not consult atomic_status; outcome already records PNH's verdict")
		assert.Equal(t, []string{earlyRaceSharedID}, svc.spyrevertSharedIDs,
			"early-race row should route directly to HandleDestinationReverted")
	})

	t.Run("merges early-race and late-race rows into a single HandleDestinationReverted call", func(t *testing.T) {
		earlyRaceSharedID := "early-race-shared-id"
		earlyRaceTx := types.Transaction{
			SharedID: earlyRaceSharedID,
			Outcome:  types.OutcomeReverted,
		}

		txRepo := &StubFinalizationTransactionRepo{
			txs: []types.Transaction{earlyRaceTx, failedTx},
		}
		sumsRepo := &StubFinalizationAtomicStatusRepo{
			sums: []types.AtomicStatusUpdateMessage{failedSUM},
		}
		svc := &StubFinalizationSignatureService{}
		p := poller.NewFinalizationPoller(batchSize, svc, txRepo, sumsRepo)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = p.Run(ctx)

		assert.ElementsMatch(t, []string{revertSharedID, earlyRaceSharedID}, svc.spyrevertSharedIDs,
			"both early- and late-race rows should be reverted in one call")
	})
}
