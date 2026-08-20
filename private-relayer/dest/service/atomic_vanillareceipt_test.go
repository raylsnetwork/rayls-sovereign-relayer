package service_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/EndpointV1"
	service "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/service"
	sharedservice "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StubVanillaReceiptTeleportClient struct {
	err error

	spySharedIDs      []string
	spyAdditionalData []types.DispatchedMessageToPrivateHub
}

func (s *StubVanillaReceiptTeleportClient) StoreEncryptedDataBatch(
	_ context.Context,
	sharedIDs []string,
	additionalData []types.DispatchedMessageToPrivateHub,
	chainID *big.Int,
) (common.Hash, error) {
	s.spySharedIDs = sharedIDs
	s.spyAdditionalData = additionalData
	return common.Hash{}, s.err
}

type StubVanillaReceiptTransactionRepository struct {
	txs       []types.Transaction
	getErr    error
	updateErr error

	spyGetSharedIDs []string

	spyUpdateCalled    bool
	spyUpdateSharedIDs []string
	spyState           types.TransactionState
	spyOutcome         types.TransactionOutcome
}

func (r *StubVanillaReceiptTransactionRepository) GetBySharedIDs(
	ctx context.Context,
	sharedIDs []string,
) ([]types.Transaction, error) {
	r.spyGetSharedIDs = sharedIDs
	return r.txs, r.getErr
}

func (r *StubVanillaReceiptTransactionRepository) BatchSetStateAndOutcome(
	ctx context.Context,
	sharedIDs []string,
	state types.TransactionState,
	outcome types.TransactionOutcome,
) error {
	r.spyUpdateCalled = true

	r.spyUpdateSharedIDs = sharedIDs
	r.spyState = state
	r.spyOutcome = outcome
	return r.updateErr
}

type StubVanillaReceiptEthereumClient struct {
	receipt *ethTypes.Receipt
	block   *ethTypes.Block
	err     error

	spyTxHash    common.Hash
	spyBlockHash common.Hash
}

func (c *StubVanillaReceiptEthereumClient) TransactionReceipt(
	ctx context.Context,
	txHash common.Hash,
) (*ethTypes.Receipt, error) {
	c.spyTxHash = txHash
	return c.receipt, c.err
}

func (c *StubVanillaReceiptEthereumClient) BlockByHash(
	ctx context.Context,
	blockHash common.Hash,
) (*ethTypes.Block, error) {
	c.spyBlockHash = blockHash
	return c.block, c.err
}

type StubVanillaReceiptProofGenerator struct {
	proof []byte
	err   error

	spyHash common.Hash
}

func (g *StubVanillaReceiptProofGenerator) Generate(_ context.Context, hash common.Hash) ([]byte, error) {
	g.spyHash = hash
	return g.proof, g.err
}

func TestVanillaReceiptService_HandleSuccessfullyMined(t *testing.T) {
	t.Run("marks messages as executed and sends encrypted storage proofs", func(t *testing.T) {
		wantSharedID := "example-shared-id"
		wantTxHashDestination := common.Hash([32]byte{0xDE, 0xAD, 0xC0, 0xDE})
		wantAdditionalData := types.DispatchedMessageToPrivateHub{
			MessageId:                  [32]byte{0xC0, 0xFE, 0xBA, 0xBE},
			FromChainId:                new(big.Int).SetUint64(1337),
			ToChainId:                  new(big.Int).SetUint64(1337),
			Proofs:                     []byte{0xDE, 0xAD, 0xBE, 0xEF},
			TxTrieProof:                [32]byte{0xD0, 0x0D, 0x2B, 0xAD},
			TxHashDestination:          wantTxHashDestination,
			TxHashDestinationTimestamp: 1001,
			TxHashDestinationStatus:    1,
			TransactionType:            types.Proof,
			Data: EndpointV1.RaylsMessage{
				MessageMetadata: EndpointV1.RaylsMessageMetadata{
					TransferMetadata: EndpointV1.BridgedTransferMetadata{
						AssetType: uint8(types.CUSTOM),
					},
				},
			},
			SharedId: wantSharedID,
		}
		wantBlockHash := common.Hash([32]byte{0xC0, 0x01, 0x0F, 0xFF})

		chainID := wantAdditionalData.FromChainId
		tx := types.Transaction{
			SharedID:          wantSharedID,
			MsgID:             wantAdditionalData.MessageId,
			TxHashDestination: wantAdditionalData.TxHashDestination,
		}
		receipt := &ethTypes.Receipt{
			BlockHash: wantBlockHash,
		}
		block := ethTypes.NewBlockWithHeader(&ethTypes.Header{
			TxHash: wantAdditionalData.TxTrieProof,
			Time:   wantAdditionalData.TxHashDestinationTimestamp,
		})

		teleportCli := &StubVanillaReceiptTeleportClient{}
		ethereumCli := &StubVanillaReceiptEthereumClient{
			receipt: receipt,
			block:   block,
		}
		proofGen := &StubVanillaReceiptProofGenerator{
			proof: wantAdditionalData.Proofs,
		}
		transactionRepo := &StubVanillaReceiptTransactionRepository{
			txs: []types.Transaction{tx},
		}
		svc := service.NewVanillaReceiptService(chainID, teleportCli, ethereumCli, proofGen, transactionRepo)

		err := svc.HandleSuccessfullyMined(context.Background(), []string{wantSharedID})
		require.Nil(t, err)

		assert.Equal(
			t,
			[]string{wantSharedID},
			transactionRepo.spyGetSharedIDs,
			"didn't request transactions for correct shared IDs",
		)

		assert.Equal(t, wantTxHashDestination, ethereumCli.spyTxHash)
		assert.Equal(t, wantBlockHash, ethereumCli.spyBlockHash)

		assert.Equal(
			t,
			wantAdditionalData.TxHashDestination,
			proofGen.spyHash,
			"didn't generate proof for correct transaction",
		)

		assert.Equal(
			t,
			[]types.DispatchedMessageToPrivateHub{wantAdditionalData},
			teleportCli.spyAdditionalData,
			"didn't send correct additional data to teleport contract",
		)

		assert.Equal(
			t,
			[]string{wantSharedID},
			transactionRepo.spyUpdateSharedIDs,
			"didn't update state for correct shared IDs",
		)
		assert.Equal(
			t,
			types.DestinationDispatch,
			transactionRepo.spyState,
			"didn't set correct state for shared IDs",
		)
		assert.Equal(
			t,
			types.OutcomeSuccess,
			transactionRepo.spyOutcome,
			"didn't set correct outcome for shared IDs",
		)
	})

	t.Run("wraps get transactions repository errors in AtomicServiceError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &sharedservice.AtomicServiceError{}

		sharedID := "example-shared-id"
		chainID := new(big.Int).SetUint64(1337)
		tx := types.Transaction{
			SharedID:          sharedID,
			MsgID:             [32]byte{0xC0, 0xFE, 0xBA, 0xBE},
			TxHashDestination: [32]byte{0xDE, 0xAD, 0xC0, 0xDE},
		}
		receipt := &ethTypes.Receipt{
			BlockHash: common.Hash([32]byte{0xC0, 0x01, 0x0F, 0xFF}),
		}
		block := ethTypes.NewBlockWithHeader(&ethTypes.Header{
			TxHash: [32]byte{0xD0, 0x0D, 0x2B, 0xAD},
			Time:   1001,
		})

		teleportCli := &StubVanillaReceiptTeleportClient{}
		ethereumCli := &StubVanillaReceiptEthereumClient{
			receipt: receipt,
			block:   block,
		}
		proofGen := &StubVanillaReceiptProofGenerator{
			proof: []byte{0xDE, 0xAD, 0xBE, 0xEF},
		}
		transactionRepo := &StubVanillaReceiptTransactionRepository{
			txs:    []types.Transaction{tx},
			getErr: wantError,
		}
		svc := service.NewVanillaReceiptService(chainID, teleportCli, ethereumCli, proofGen, transactionRepo)

		err := svc.HandleSuccessfullyMined(context.Background(), []string{sharedID})

		assert.ErrorAs(t, err, &wantErrorType)
		assert.ErrorIs(t, err, wantError)

		assert.False(t, transactionRepo.spyUpdateCalled, "called update state but shouldn't have")
	})

	t.Run("wraps generate proof errors in AtomicServiceError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &sharedservice.AtomicServiceError{}

		sharedID := "example-shared-id"
		chainID := new(big.Int).SetUint64(1337)
		tx := types.Transaction{
			SharedID:          sharedID,
			MsgID:             [32]byte{0xC0, 0xFE, 0xBA, 0xBE},
			TxHashDestination: [32]byte{0xDE, 0xAD, 0xC0, 0xDE},
		}
		receipt := &ethTypes.Receipt{
			BlockHash: common.Hash([32]byte{0xC0, 0x01, 0x0F, 0xFF}),
		}
		block := ethTypes.NewBlockWithHeader(&ethTypes.Header{
			TxHash: [32]byte{0xD0, 0x0D, 0x2B, 0xAD},
			Time:   1001,
		})

		teleportCli := &StubVanillaReceiptTeleportClient{}
		ethereumCli := &StubVanillaReceiptEthereumClient{
			receipt: receipt,
			block:   block,
		}
		proofGen := &StubVanillaReceiptProofGenerator{
			proof: []byte{0xDE, 0xAD, 0xBE, 0xEF},
			err:   wantError,
		}
		transactionRepo := &StubVanillaReceiptTransactionRepository{
			txs: []types.Transaction{tx},
		}
		svc := service.NewVanillaReceiptService(chainID, teleportCli, ethereumCli, proofGen, transactionRepo)

		err := svc.HandleSuccessfullyMined(context.Background(), []string{sharedID})

		assert.ErrorAs(t, err, &wantErrorType)
		assert.ErrorIs(t, err, wantError)

		assert.False(t, transactionRepo.spyUpdateCalled, "called update state but shouldn't have")
	})

	t.Run("wraps teleport client errors in AtomicServiceError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &sharedservice.AtomicServiceError{}

		sharedID := "example-shared-id"
		chainID := new(big.Int).SetUint64(1337)
		tx := types.Transaction{
			SharedID:          sharedID,
			MsgID:             [32]byte{0xC0, 0xFE, 0xBA, 0xBE},
			TxHashDestination: [32]byte{0xDE, 0xAD, 0xC0, 0xDE},
		}
		receipt := &ethTypes.Receipt{
			BlockHash: common.Hash([32]byte{0xC0, 0x01, 0x0F, 0xFF}),
		}
		block := ethTypes.NewBlockWithHeader(&ethTypes.Header{
			TxHash: [32]byte{0xD0, 0x0D, 0x2B, 0xAD},
			Time:   1001,
		})

		teleportCli := &StubVanillaReceiptTeleportClient{
			err: wantError,
		}
		ethereumCli := &StubVanillaReceiptEthereumClient{
			receipt: receipt,
			block:   block,
		}
		proofGen := &StubVanillaReceiptProofGenerator{
			proof: []byte{0xDE, 0xAD, 0xBE, 0xEF},
		}
		transactionRepo := &StubVanillaReceiptTransactionRepository{
			txs: []types.Transaction{tx},
		}
		svc := service.NewVanillaReceiptService(chainID, teleportCli, ethereumCli, proofGen, transactionRepo)

		err := svc.HandleSuccessfullyMined(context.Background(), []string{sharedID})

		assert.ErrorAs(t, err, &wantErrorType)
		assert.ErrorIs(t, err, wantError)

		assert.False(t, transactionRepo.spyUpdateCalled, "called update state but shouldn't have")
	})

	t.Run("wraps ethereum client errors in AtomicServiceError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &sharedservice.AtomicServiceError{}

		sharedID := "example-shared-id"
		chainID := new(big.Int).SetUint64(1337)
		tx := types.Transaction{
			SharedID:          sharedID,
			MsgID:             [32]byte{0xC0, 0xFE, 0xBA, 0xBE},
			TxHashDestination: [32]byte{0xDE, 0xAD, 0xC0, 0xDE},
		}
		block := ethTypes.NewBlockWithHeader(&ethTypes.Header{
			TxHash: [32]byte{0xD0, 0x0D, 0x2B, 0xAD},
			Time:   1001,
		})

		teleportCli := &StubVanillaReceiptTeleportClient{}
		ethereumCli := &StubVanillaReceiptEthereumClient{
			block: block,
			err:   wantError,
		}
		proofGen := &StubVanillaReceiptProofGenerator{
			proof: []byte{0xDE, 0xAD, 0xBE, 0xEF},
		}
		transactionRepo := &StubVanillaReceiptTransactionRepository{
			txs: []types.Transaction{tx},
		}
		svc := service.NewVanillaReceiptService(chainID, teleportCli, ethereumCli, proofGen, transactionRepo)

		err := svc.HandleSuccessfullyMined(context.Background(), []string{sharedID})

		assert.ErrorAs(t, err, &wantErrorType)
		assert.ErrorIs(t, err, wantError)

		assert.False(t, transactionRepo.spyUpdateCalled, "called update state but shouldn't have")
	})

	t.Run("wraps update state repository errors in AtomicServiceError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &sharedservice.AtomicServiceError{}

		sharedID := "example-shared-id"
		chainID := new(big.Int).SetUint64(1337)
		tx := types.Transaction{
			SharedID:          sharedID,
			MsgID:             [32]byte{0xC0, 0xFE, 0xBA, 0xBE},
			TxHashDestination: [32]byte{0xDE, 0xAD, 0xC0, 0xDE},
		}
		receipt := &ethTypes.Receipt{
			BlockHash: common.Hash([32]byte{0xC0, 0x01, 0x0F, 0xFF}),
		}
		block := ethTypes.NewBlockWithHeader(&ethTypes.Header{
			TxHash: [32]byte{0xD0, 0x0D, 0x2B, 0xAD},
			Time:   1001,
		})

		teleportCli := &StubVanillaReceiptTeleportClient{}
		ethereumCli := &StubVanillaReceiptEthereumClient{
			receipt: receipt,
			block:   block,
		}
		proofGen := &StubVanillaReceiptProofGenerator{
			proof: []byte{0xDE, 0xAD, 0xBE, 0xEF},
		}
		transactionRepo := &StubVanillaReceiptTransactionRepository{
			txs:       []types.Transaction{tx},
			updateErr: wantError,
		}
		svc := service.NewVanillaReceiptService(chainID, teleportCli, ethereumCli, proofGen, transactionRepo)

		err := svc.HandleSuccessfullyMined(context.Background(), []string{sharedID})

		assert.ErrorAs(t, err, &wantErrorType)
		assert.ErrorIs(t, err, wantError)
	})
}

func TestVanillaReceiptService_HandleFailedMined(t *testing.T) {
	t.Run("marks transactions as failed", func(t *testing.T) {
		wantSharedID := "example-shared-id"

		chainID := new(big.Int).SetUint64(1337)

		teleportCli := &StubVanillaReceiptTeleportClient{}
		ethereumCli := &StubVanillaReceiptEthereumClient{}
		proofGen := &StubVanillaReceiptProofGenerator{}
		transactionRepo := &StubVanillaReceiptTransactionRepository{}
		svc := service.NewVanillaReceiptService(chainID, teleportCli, ethereumCli, proofGen, transactionRepo)

		err := svc.HandleFailedMined(context.Background(), []string{wantSharedID})
		require.Nil(t, err)

		assert.Equal(
			t,
			[]string{wantSharedID},
			transactionRepo.spyUpdateSharedIDs,
			"didn't update state for correct shared IDs",
		)
		assert.Equal(
			t,
			types.DestinationDispatch,
			transactionRepo.spyState,
			"didn't set correct state for shared IDs",
		)
		assert.Equal(
			t,
			types.OutcomeReverted,
			transactionRepo.spyOutcome,
			"didn't set correct outcome for shared IDs",
		)
	})

	t.Run("wraps update state repository errors in AtomicServiceError", func(t *testing.T) {
		wantError := errors.New("example error")
		wantErrorType := &sharedservice.AtomicServiceError{}

		sharedID := "example-shared-id"
		chainID := new(big.Int).SetUint64(1337)

		teleportCli := &StubVanillaReceiptTeleportClient{}
		ethereumCli := &StubVanillaReceiptEthereumClient{}
		proofGen := &StubVanillaReceiptProofGenerator{}
		transactionRepo := &StubVanillaReceiptTransactionRepository{
			updateErr: wantError,
		}
		svc := service.NewVanillaReceiptService(chainID, teleportCli, ethereumCli, proofGen, transactionRepo)

		err := svc.HandleFailedMined(context.Background(), []string{sharedID})

		assert.ErrorAs(t, err, &wantErrorType)
		assert.ErrorIs(t, err, wantError)
	})
}
