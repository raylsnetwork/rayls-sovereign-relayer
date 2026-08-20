// Decommissioning Teleport (vanilla, atomic).

package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	service "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/service"
	sharedservice "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/txsim"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StubReceiptTeleportClient struct {
	err error

	spyExecuteSharedIDs      []string
	spyExecuteAdditionalData []types.AtomicTeleportAdditionalData

	spyRevertSharedIDs      []string
	spyRevertAdditionalData []types.AtomicTeleportAdditionalData
}

func (c *StubReceiptTeleportClient) ExecuteAtomicMessageBatch(
	_ context.Context,
	sharedIDs []string,
	additionalData []types.AtomicTeleportAdditionalData,
) error {
	c.spyExecuteSharedIDs = sharedIDs
	c.spyExecuteAdditionalData = additionalData
	return c.err
}

func (c *StubReceiptTeleportClient) RevertAtomicMessageBatch(
	_ context.Context,
	sharedIDs []string,
	additionalData []types.AtomicTeleportAdditionalData,
) error {
	c.spyRevertSharedIDs = sharedIDs
	c.spyRevertAdditionalData = additionalData
	return c.err
}

type StubReceiptTransactionRepository struct {
	txs []types.Transaction

	spyGetSharedIDs    []string
	spyUpdateSharedIDs []string
	spyState           types.TransactionState
	spyOutcome         types.TransactionOutcome
}

func (r *StubReceiptTransactionRepository) GetBySharedIDs(
	ctx context.Context,
	sharedIDs []string,
) ([]types.Transaction, error) {
	r.spyGetSharedIDs = sharedIDs
	return r.txs, nil
}

func (r *StubReceiptTransactionRepository) BatchSetStateAndOutcome(
	ctx context.Context,
	sharedIDs []string,
	state types.TransactionState,
	outcome types.TransactionOutcome,
) error {
	r.spyUpdateSharedIDs = sharedIDs
	r.spyState = state
	r.spyOutcome = outcome
	return nil
}

type StubReceiptEthereumClient struct {
	header          *ethTypes.Header
	headerByHashErr error

	receipt               *ethTypes.Receipt
	transactionReceiptErr error

	spyTxHash    common.Hash
	spyBlockHash common.Hash
}

func (c *StubReceiptEthereumClient) TransactionReceipt(
	ctx context.Context,
	txHash common.Hash,
) (*ethTypes.Receipt, error) {
	c.spyTxHash = txHash
	return c.receipt, c.transactionReceiptErr
}

func (c *StubReceiptEthereumClient) HeaderByHash(ctx context.Context, blockHash common.Hash) (*ethTypes.Header, error) {
	c.spyBlockHash = blockHash
	return c.header, c.headerByHashErr
}

type StubReceiptTransactionSimulator struct {
	revertReason txsim.ContractError
	err          error
}

func (s *StubReceiptTransactionSimulator) GetRevertReason(
	ctx context.Context,
	hash common.Hash,
) (txsim.ContractError, error) {
	return s.revertReason, s.err
}

func TestReceiptService_HandleSuccessfullyMined(t *testing.T) {
	t.Run("marks messages as executed and sends additional data", func(t *testing.T) {
		wantTxHash := common.Hash([32]byte{0xDE, 0xAD, 0xBE, 0xEF})
		wantAdditionalData := types.AtomicTeleportAdditionalData{
			TxHashDestination:          wantTxHash,
			TxHashDestinationStatus:    1,
			TxHashDestinationTimestamp: 12345,
			SharedId:                   "example-shared-id",
			BatchPrivateHubHash:        common.Hash([32]byte{0xC0, 0xFE, 0xBA, 0xBE}),
		}
		wantBlockHash := common.Hash([32]byte{0xC0, 0x01, 0x0F, 0xFF})
		wantSharedIDs := []string{wantAdditionalData.SharedId}
		wantState := types.HubNotifiedExec

		receipt := &ethTypes.Receipt{
			BlockHash: wantBlockHash,
		}
		header := &ethTypes.Header{
			Time: wantAdditionalData.TxHashDestinationTimestamp,
		}
		tx := types.Transaction{
			TxHashDestination:   wantAdditionalData.TxHashDestination,
			BatchPrivateHubHash: wantAdditionalData.BatchPrivateHubHash,
			SharedID:            wantAdditionalData.SharedId,
		}

		teleportCli := &StubReceiptTeleportClient{}
		ethereumCli := &StubReceiptEthereumClient{
			header:  header,
			receipt: receipt,
		}
		transactionRepo := &StubReceiptTransactionRepository{
			txs: []types.Transaction{tx},
		}
		transactionSim := &StubReceiptTransactionSimulator{}
		svc := service.NewReceiptService(teleportCli, ethereumCli, transactionRepo, transactionSim)

		err := svc.HandleSuccessfullyMined(context.Background(), wantSharedIDs)
		require.Nil(t, err)

		assert.Equal(
			t,
			wantSharedIDs,
			transactionRepo.spyGetSharedIDs,
			"didn't get transactions for correct shared IDs",
		)

		assert.Equal(t, wantTxHash, ethereumCli.spyTxHash)
		assert.Equal(t, wantBlockHash, ethereumCli.spyBlockHash)

		assert.Equal(t, wantSharedIDs, teleportCli.spyExecuteSharedIDs)
		assert.Equal(t, []types.AtomicTeleportAdditionalData{wantAdditionalData}, teleportCli.spyExecuteAdditionalData)

		assert.Equal(t, wantSharedIDs, transactionRepo.spyUpdateSharedIDs)
		assert.Equal(t, wantState, transactionRepo.spyState)
	})

	t.Run("sets timestamp to zero on error while getting header", func(t *testing.T) {
		wantTimestamp := uint64(0)

		sharedID := "example-shared-ID"
		tx := types.Transaction{
			TxHashDestination:   common.Hash([32]byte{0xDE, 0xAD, 0xBE, 0xEF}),
			BatchPrivateHubHash: common.Hash([32]byte{0xCA, 0xFE, 0xBA, 0xBE}),
			SharedID:            sharedID,
		}
		receipt := &ethTypes.Receipt{
			BlockHash: common.Hash([32]byte{0xC0, 0x01, 0x0F, 0xFF}),
		}

		teleportCli := &StubReceiptTeleportClient{}
		ethereumCli := &StubReceiptEthereumClient{
			receipt:         receipt,
			headerByHashErr: errors.New("example-error"),
		}
		transactionRepo := &StubReceiptTransactionRepository{
			txs: []types.Transaction{tx},
		}
		transactionSim := &StubReceiptTransactionSimulator{}
		svc := service.NewReceiptService(teleportCli, ethereumCli, transactionRepo, transactionSim)

		err := svc.HandleSuccessfullyMined(context.Background(), []string{sharedID})
		require.Nil(t, err)

		gotTimestamp := teleportCli.spyExecuteAdditionalData[0].TxHashDestinationTimestamp
		assert.Equal(t, wantTimestamp, gotTimestamp)
	})

	t.Run("records outcome=reverted on ErrAlreadyReverted from teleport client", func(t *testing.T) {
		wantSharedIDs := []string{"example-shared-id"}
		wantState := types.HubNotifiedExec
		wantOutcome := types.OutcomeReverted

		tx := types.Transaction{
			TxHashDestination:   common.Hash([32]byte{0xDE, 0xAD, 0xBE, 0xEF}),
			BatchPrivateHubHash: common.Hash([32]byte{0xCA, 0xFE, 0xBA, 0xBE}),
			SharedID:            wantSharedIDs[0],
		}
		teleportCli := &StubReceiptTeleportClient{
			err: sharedservice.ErrAlreadyReverted,
		}
		ethereumCli := &StubReceiptEthereumClient{
			transactionReceiptErr: errors.New("no available receipt"),
			headerByHashErr:       errors.New("no available header"),
		}
		transactionRepo := &StubReceiptTransactionRepository{
			txs: []types.Transaction{tx},
		}
		transactionSim := &StubReceiptTransactionSimulator{}
		svc := service.NewReceiptService(teleportCli, ethereumCli, transactionRepo, transactionSim)

		err := svc.HandleSuccessfullyMined(context.Background(), wantSharedIDs)
		require.Nil(t, err)

		assert.Equal(t, wantSharedIDs, transactionRepo.spyUpdateSharedIDs)
		assert.Equal(t, wantState, transactionRepo.spyState)
		assert.Equal(t, wantOutcome, transactionRepo.spyOutcome)
	})

	t.Run(
		"doesn't update state on error from teleport client and wraps error in AtomicServiceError",
		func(t *testing.T) {
			wantError := errors.New("example error")
			wantErrorType := &sharedservice.AtomicServiceError{}

			sharedID := "example-shared-ID"
			tx := types.Transaction{
				TxHashDestination:   common.Hash([32]byte{0xDE, 0xAD, 0xBE, 0xEF}),
				BatchPrivateHubHash: common.Hash([32]byte{0xCA, 0xFE, 0xBA, 0xBE}),
				SharedID:            sharedID,
			}
			teleportCli := &StubReceiptTeleportClient{
				err: wantError,
			}
			ethereumCli := &StubReceiptEthereumClient{
				transactionReceiptErr: errors.New("no available receipt"),
				headerByHashErr:       errors.New("no available header"),
			}
			transactionRepo := &StubReceiptTransactionRepository{
				txs: []types.Transaction{tx},
			}
			transactionSim := &StubReceiptTransactionSimulator{}
			svc := service.NewReceiptService(teleportCli, ethereumCli, transactionRepo, transactionSim)

			err := svc.HandleSuccessfullyMined(context.Background(), []string{sharedID})

			assert.ErrorAs(t, err, &wantErrorType)
			assert.ErrorIs(t, err, wantError)

			assert.Empty(t, transactionRepo.spyUpdateSharedIDs)
		},
	)
}

func TestReceiptService_HandleFailedMined(t *testing.T) {
	t.Run("marks messages as reverted and sends additional data", func(t *testing.T) {
		wantTxHash := common.Hash([32]byte{0xDE, 0xAD, 0xBE, 0xEF})
		wantAdditionalData := types.AtomicTeleportAdditionalData{
			TxHashDestination:          wantTxHash,
			TxHashDestinationStatus:    1,
			TxHashDestinationTimestamp: 12345,
			RevertReason:               "",
			SharedId:                   "example-shared-id",
			BatchPrivateHubHash:        common.Hash([32]byte{0xCA, 0xFE, 0xBA, 0xBE}),
		}
		wantBlockHash := common.Hash([32]byte{0xC0, 0x01, 0x0F, 0xFF})
		wantSharedIDs := []string{wantAdditionalData.SharedId}
		wantState := types.HubNotifiedRevert

		receipt := &ethTypes.Receipt{
			BlockHash: wantBlockHash,
		}
		header := &ethTypes.Header{
			Time: wantAdditionalData.TxHashDestinationTimestamp,
		}
		tx := types.Transaction{
			TxHashDestination:   wantAdditionalData.TxHashDestination,
			BatchPrivateHubHash: wantAdditionalData.BatchPrivateHubHash,
			SharedID:            wantAdditionalData.SharedId,
		}

		teleportCli := &StubReceiptTeleportClient{}
		ethereumCli := &StubReceiptEthereumClient{
			header:  header,
			receipt: receipt,
		}
		transactionRepo := &StubReceiptTransactionRepository{
			txs: []types.Transaction{tx},
		}
		transactionSim := &StubReceiptTransactionSimulator{}
		svc := service.NewReceiptService(teleportCli, ethereumCli, transactionRepo, transactionSim)

		err := svc.HandleFailedMined(context.Background(), wantSharedIDs)
		require.Nil(t, err)

		assert.Equal(
			t,
			wantSharedIDs,
			transactionRepo.spyGetSharedIDs,
			"didn't get transactions for correct shared IDs",
		)

		assert.Equal(t, wantTxHash, ethereumCli.spyTxHash)
		assert.Equal(t, wantBlockHash, ethereumCli.spyBlockHash)

		assert.Equal(t, wantSharedIDs, teleportCli.spyRevertSharedIDs)
		assert.Equal(t, []types.AtomicTeleportAdditionalData{wantAdditionalData}, teleportCli.spyRevertAdditionalData)

		assert.Equal(t, wantSharedIDs, transactionRepo.spyUpdateSharedIDs)
		assert.Equal(t, wantState, transactionRepo.spyState)
	})

	t.Run("leaves revert reason empty in case of error from transaction simulator", func(t *testing.T) {
		sharedID := "example-shared-ID"
		tx := types.Transaction{
			TxHashDestination:   common.Hash([32]byte{0xDE, 0xAD, 0xBE, 0xEF}),
			BatchPrivateHubHash: common.Hash([32]byte{0xCA, 0xFE, 0xBA, 0xBE}),
			SharedID:            sharedID,
		}
		teleportCli := &StubReceiptTeleportClient{}
		ethereumCli := &StubReceiptEthereumClient{
			transactionReceiptErr: errors.New("no available receipt"),
			headerByHashErr:       errors.New("no available header"),
		}
		transactionRepo := &StubReceiptTransactionRepository{
			txs: []types.Transaction{tx},
		}
		transactionSim := &StubReceiptTransactionSimulator{
			err: errors.New("example error"),
		}
		svc := service.NewReceiptService(teleportCli, ethereumCli, transactionRepo, transactionSim)

		err := svc.HandleFailedMined(context.Background(), []string{sharedID})
		require.Nil(t, err)

		require.Len(t, teleportCli.spyRevertAdditionalData, 1)

		gotRevertReason := teleportCli.spyRevertAdditionalData[0].RevertReason
		assert.Equal(t, "", gotRevertReason, "zero-value ContractError.String() returns empty string")
	})
}
