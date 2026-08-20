// Decommissioning Teleport (vanilla, atomic).

package service_test

import (
	"context"
	"errors"
	"testing"

	sharedservice "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StubReceiptTransactionRepo struct {
	spySharedIDs []string
	spyState     types.TransactionState
	spyOutcome   types.TransactionOutcome
	err          error
}

func (r *StubReceiptTransactionRepo) BatchSetStateAndOutcome(
	ctx context.Context,
	sharedIDs []string,
	state types.TransactionState,
	outcome types.TransactionOutcome,
) error {
	r.spySharedIDs = sharedIDs
	r.spyState = state
	r.spyOutcome = outcome
	return r.err
}

func TestReceiptService_HandleSuccessfullyMined(t *testing.T) {
	t.Run("updates state to SourcePublish + OutcomeSuccess", func(t *testing.T) {
		wantSharedIDs := []string{"example-shared-id"}
		wantState := types.SourcePublish
		wantOutcome := types.OutcomeSuccess

		txRepo := &StubReceiptTransactionRepo{}
		svc := service.NewReceiptService(txRepo)

		err := svc.HandleSuccessfullyMined(context.Background(), wantSharedIDs)
		require.Nil(t, err)

		assert.Equal(t, wantSharedIDs, txRepo.spySharedIDs)
		assert.Equal(t, wantState, txRepo.spyState)
		assert.Equal(t, wantOutcome, txRepo.spyOutcome)
	})

	t.Run("wraps errors in AtomicServiceError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &sharedservice.AtomicServiceError{}

		txRepo := &StubReceiptTransactionRepo{err: wantError}
		svc := service.NewReceiptService(txRepo)

		err := svc.HandleSuccessfullyMined(context.Background(), []string{"shared-id"})

		require.ErrorAs(t, err, &wantErrorType)
		require.ErrorIs(t, err, wantError)
	})
}

func TestReceiptService_HandleFailedMined(t *testing.T) {
	t.Run("updates state to SourcePublish + OutcomeFailed", func(t *testing.T) {
		wantSharedIDs := []string{"example-shared-id"}
		wantState := types.SourcePublish
		wantOutcome := types.OutcomeFailed

		txRepo := &StubReceiptTransactionRepo{}
		svc := service.NewReceiptService(txRepo)

		err := svc.HandleFailedMined(context.Background(), wantSharedIDs)
		require.Nil(t, err)

		assert.Equal(t, wantSharedIDs, txRepo.spySharedIDs)
		assert.Equal(t, wantState, txRepo.spyState)
		assert.Equal(t, wantOutcome, txRepo.spyOutcome)
	})

	t.Run("wraps errors in AtomicServiceError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &sharedservice.AtomicServiceError{}

		txRepo := &StubReceiptTransactionRepo{err: wantError}
		svc := service.NewReceiptService(txRepo)

		err := svc.HandleFailedMined(context.Background(), []string{"shared-id"})

		require.ErrorAs(t, err, &wantErrorType)
		require.ErrorIs(t, err, wantError)
	})
}
