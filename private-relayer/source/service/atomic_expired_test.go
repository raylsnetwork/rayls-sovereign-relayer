// Decommissioning Teleport (vanilla, atomic).

package service_test

import (
	"context"
	"errors"
	"testing"

	sharedservice "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StubExpiredTeleportClient struct {
	statuses  []types.AtomicStatusUpdateMessage
	getErr    error
	revertErr error

	spyGetSharedIDs []string

	spyCalledRevert    bool
	spyRevertSharedIDs []string
	spyAdditionalData  []types.AtomicTeleportAdditionalData
}

func (c *StubExpiredTeleportClient) GetAtomicMessageStatuses(
	_ context.Context,
	sharedIDs []string,
) ([]types.AtomicStatusUpdateMessage, error) {
	c.spyGetSharedIDs = sharedIDs
	return c.statuses, c.getErr
}

func (c *StubExpiredTeleportClient) RevertAtomicMessageBatch(
	_ context.Context,
	sharedIDs []string,
	additionalData []types.AtomicTeleportAdditionalData,
) error {
	c.spyCalledRevert = true

	c.spyRevertSharedIDs = sharedIDs
	c.spyAdditionalData = additionalData

	return c.revertErr
}

type StubExpiredTransactionRepo struct {
	spyCalledUpdate bool
	spySharedIDs    []string
	spyState        types.TransactionState
	spyOutcome      types.TransactionOutcome

	err error
}

func (r *StubExpiredTransactionRepo) BatchSetStateAndOutcome(
	ctx context.Context,
	sharedIDs []string,
	state types.TransactionState,
	outcome types.TransactionOutcome,
) error {
	r.spyCalledUpdate = true

	r.spySharedIDs = sharedIDs
	r.spyState = state
	r.spyOutcome = outcome

	return r.err
}

func TestExpiredService_RevertExpired(t *testing.T) {
	t.Run("reverts expired transactions that are still pending in the teleport contract", func(t *testing.T) {
		wantSharedID := "example-shared-id"
		wantState := types.SourceTimeoutRevert
		wantOutcome := types.OutcomeSuccess

		sharedIDs := []string{wantSharedID, "another-shared id", "third-shared-id"}

		teleportCli := &StubExpiredTeleportClient{
			statuses: []types.AtomicStatusUpdateMessage{
				{
					SharedID: wantSharedID,
					Status:   types.AtomicPendingStatus,
				},
				{
					SharedID: sharedIDs[1],
					Status:   types.AtomicExecutedStatus,
				},
				{
					SharedID: sharedIDs[2],
					Status:   types.AtomicRevertedStatus,
				},
			},
		}
		repo := &StubExpiredTransactionRepo{}
		svc := service.NewExpiredService(teleportCli, repo)

		err := svc.Handle(context.Background(), sharedIDs)
		require.Nil(t, err)

		assert.Equal(t, sharedIDs, teleportCli.spyGetSharedIDs)

		assert.Empty(t, teleportCli.spyAdditionalData)
		assert.Equal(t, []string{wantSharedID}, teleportCli.spyRevertSharedIDs)

		assert.Equal(t, wantState, repo.spyState)
		// The timeout-revert Hub tx is synchronous-confirmed, so the row must be recorded
		// as +success (NOT +pending) — otherwise the finalization poller never picks it up
		// and the sender refund (SourceRevertSigs → revertTeleportMint) never runs.
		assert.Equal(t, wantOutcome, repo.spyOutcome)
		assert.Equal(t, []string{wantSharedID}, repo.spySharedIDs)
	})

	t.Run("returns on no transactions in pending state in the teleport contract", func(t *testing.T) {
		sharedIDs := []string{"another-shared id", "third-shared-id"}

		teleportCli := &StubExpiredTeleportClient{
			statuses: []types.AtomicStatusUpdateMessage{
				{
					SharedID: sharedIDs[0],
					Status:   types.AtomicExecutedStatus,
				},
				{
					SharedID: sharedIDs[1],
					Status:   types.AtomicRevertedStatus,
				},
			},
		}
		repo := &StubExpiredTransactionRepo{}
		svc := service.NewExpiredService(teleportCli, repo)

		err := svc.Handle(context.Background(), sharedIDs)
		require.Nil(t, err)

		assert.Equal(t, sharedIDs, teleportCli.spyGetSharedIDs)

		// Check we didn't call teleport contract
		assert.False(t, teleportCli.spyCalledRevert)

		// Check we didn't update state for transactions
		assert.False(t, repo.spyCalledUpdate)
	})

	t.Run("wraps get status errors from teleport contract in AtomicServiceError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &sharedservice.AtomicServiceError{}

		sharedID := "example-shared-id"

		teleportCli := &StubExpiredTeleportClient{
			statuses: []types.AtomicStatusUpdateMessage{
				{
					SharedID: sharedID,
					Status:   types.AtomicPendingStatus,
				},
			},
			getErr: wantError,
		}
		repo := &StubExpiredTransactionRepo{}
		svc := service.NewExpiredService(teleportCli, repo)

		err := svc.Handle(context.Background(), []string{sharedID})

		assert.ErrorAs(t, err, &wantErrorType)
		assert.ErrorIs(t, err, wantError)

		assert.False(t, teleportCli.spyCalledRevert, "called revert messages but shouldn't have")
	})

	t.Run("doesn't update state on ErrAlreadyExecuted from teleport client", func(t *testing.T) {
		wantError := sharedservice.ErrAlreadyExecuted

		sharedID := "example-shared-id"

		teleportCli := &StubExpiredTeleportClient{
			statuses: []types.AtomicStatusUpdateMessage{
				{
					SharedID: sharedID,
					Status:   types.AtomicPendingStatus,
				},
			},
			revertErr: wantError,
		}
		repo := &StubExpiredTransactionRepo{}
		svc := service.NewExpiredService(teleportCli, repo)

		err := svc.Handle(context.Background(), []string{sharedID})
		require.Nil(t, err)

		require.Empty(t, repo.spySharedIDs, "updated transactions state but shouldn't had")
	})

	t.Run("wraps get status errors from teleport client in AtomicServiceError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &sharedservice.AtomicServiceError{}

		sharedID := "example-shared-id"

		teleportCli := &StubExpiredTeleportClient{
			statuses: []types.AtomicStatusUpdateMessage{
				{
					SharedID: sharedID,
					Status:   types.AtomicPendingStatus,
				},
			},
			revertErr: wantError,
		}
		repo := &StubExpiredTransactionRepo{}
		svc := service.NewExpiredService(teleportCli, repo)

		err := svc.Handle(context.Background(), []string{sharedID})

		assert.ErrorAs(t, err, &wantErrorType)
		assert.ErrorIs(t, err, wantError)
	})

	t.Run("wraps errors from transaction repository in AtomicServiceError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &sharedservice.AtomicServiceError{}

		sharedID := "example-shared-id"

		teleportCli := &StubExpiredTeleportClient{
			statuses: []types.AtomicStatusUpdateMessage{
				{
					SharedID: sharedID,
					Status:   types.AtomicPendingStatus,
				},
			},
		}
		repo := &StubExpiredTransactionRepo{
			err: wantError,
		}
		svc := service.NewExpiredService(teleportCli, repo)

		err := svc.Handle(context.Background(), []string{sharedID})

		assert.ErrorAs(t, err, &wantErrorType)
		assert.ErrorIs(t, err, wantError)
	})
}
