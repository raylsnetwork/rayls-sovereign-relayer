package service

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EndpointV1"
	sharedservice "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

const receiptTimeout = 30 * time.Second

type vanillaTeleportClient interface {
	StoreEncryptedDataBatch(context.Context, []string, []types.DispatchedMessageToPrivateHub, *big.Int) (common.Hash, error)
}

type vanillaEthereumClient interface {
	TransactionReceipt(context.Context, common.Hash) (*ethTypes.Receipt, error)
	BlockByHash(context.Context, common.Hash) (*ethTypes.Block, error)
}

type vanillaTransactionRepository interface {
	GetBySharedIDs(context.Context, []string) ([]types.Transaction, error)
	BatchSetStateAndOutcome(ctx context.Context, sharedIDs []string, state types.TransactionState, outcome types.TransactionOutcome) error
}

type proofGenerator interface {
	Generate(context.Context, common.Hash) ([]byte, error)
}

type VanillaReceiptService struct {
	teleportCli     vanillaTeleportClient
	ethClient       vanillaEthereumClient
	proofGen        proofGenerator
	transactionRepo vanillaTransactionRepository

	myChainID *big.Int
}

func NewVanillaReceiptService(
	myChainID *big.Int,
	teleportCli vanillaTeleportClient,
	ethClient vanillaEthereumClient,
	proofGen proofGenerator,
	transactionRepo vanillaTransactionRepository,
) *VanillaReceiptService {
	return &VanillaReceiptService{
		teleportCli:     teleportCli,
		ethClient:       ethClient,
		proofGen:        proofGen,
		transactionRepo: transactionRepo,

		myChainID: myChainID,
	}
}

func (s *VanillaReceiptService) HandleSuccessfullyMined(ctx context.Context, sharedIDs []string) error {
	var additionalDataSlice []types.DispatchedMessageToPrivateHub

	txs, err := s.transactionRepo.GetBySharedIDs(ctx, sharedIDs)
	if err != nil {
		return sharedservice.WrapInAtomicServiceError("got error while querying transactions", err)
	}

	for _, tx := range txs {
		ctxReceipt, cancelReceipt := context.WithTimeout(ctx, receiptTimeout)
		var receipt *ethTypes.Receipt
		receipt, err = s.ethClient.TransactionReceipt(ctxReceipt, tx.TxHashDestination)
		cancelReceipt()
		if err != nil {
			return sharedservice.WrapInAtomicServiceError("got error while attempting to get transaction receipt", err)
		}

		ctxBlock, cancelBlock := context.WithTimeout(ctx, receiptTimeout)
		var block *ethTypes.Block
		block, err = s.ethClient.BlockByHash(ctxBlock, receipt.BlockHash)
		cancelBlock()
		if err != nil {
			return sharedservice.WrapInAtomicServiceError("got error while attempting to get block by number", err)
		}
		var proof []byte
		proof, err = s.proofGen.Generate(ctx, tx.TxHashDestination)
		if err != nil {
			return sharedservice.WrapInAtomicServiceError("got error while attempting to generate proofs", err)
		}

		additionalData := types.DispatchedMessageToPrivateHub{
			MessageId:                  tx.MsgID,
			From:                       common.Address{}, // TODO: This is not used for this type of tx but we still need to link tx in the db with more info such as fromAddress and so on
			FromChainId:                s.myChainID,
			ToChainId:                  s.myChainID,
			Proofs:                     proof,
			TxLocation:                 tx.LogIndex,
			TxTrieProof:                block.TxHash(),
			TxHashDestination:          tx.TxHashDestination,
			TxHashDestinationTimestamp: block.Time(),
			TxHashDestinationStatus:    1,
			TransactionType:            types.Proof,
			Data: EndpointV1.RaylsMessage{
				MessageMetadata: EndpointV1.RaylsMessageMetadata{
					TransferMetadata: EndpointV1.BridgedTransferMetadata{
						AssetType: uint8(types.CUSTOM),
					},
				},
			},
			SharedId: tx.SharedID,
		}
		additionalDataSlice = append(additionalDataSlice, additionalData)
	}
	_, err = s.teleportCli.StoreEncryptedDataBatch(ctx, sharedIDs, additionalDataSlice, s.myChainID)
	if err != nil {
		return sharedservice.WrapInAtomicServiceError("got error while attempting to store encrypted data batch", err)
	}

	err = s.transactionRepo.BatchSetStateAndOutcome(ctx, sharedIDs, types.DestinationDispatch, types.OutcomeSuccess)
	if err != nil {
		return sharedservice.WrapInAtomicServiceError("got error while attempting to update state", err)
	}
	return nil
}

func (s *VanillaReceiptService) HandleFailedMined(ctx context.Context, sharedIDs []string) error {
	err := s.transactionRepo.BatchSetStateAndOutcome(ctx, sharedIDs, types.DestinationDispatch, types.OutcomeReverted)
	if err != nil {
		return sharedservice.WrapInAtomicServiceError("got error while attempting to update state", err)
	}
	return nil
}

// HandleLostMined finalises a dispatch whose tx never made it on-chain
// (CTS reaper dead-lettered it after MaxResendAttempts re-broadcasts).
// We persist OutcomeFailed — distinguishable from OutcomeReverted ("mined
// and reverted on-chain") so operators / dashboards can tell a chain-loss
// apart from a real revert, and so any future retry path can key off the
// distinct outcome without ambiguity.
func (s *VanillaReceiptService) HandleLostMined(ctx context.Context, sharedIDs []string) error {
	err := s.transactionRepo.BatchSetStateAndOutcome(ctx, sharedIDs, types.DestinationDispatch, types.OutcomeFailed)
	if err != nil {
		return sharedservice.WrapInAtomicServiceError("got error while attempting to update state", err)
	}
	return nil
}
