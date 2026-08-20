package handler

//go:generate moq --skip-ensure --pkg handler_test -out receiver_mock_test.go . receiverSwapRepository receiverDepositRepository receiverMerkleClient receiverEndpointClient receiverERC721HandlerClient receiverERC1155HandlerClient receiverEnygmaHandlerClient receiverDepositFinder receiverDepositConsolidator receiverCommitmentCalculator receiverProofGenerator receiverTxPersister receiverRetryService receiverTxManager

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cryptography"
	"github.com/raylsnetwork/rayls-sovereign-relayer/dvp"
	dvpService "github.com/raylsnetwork/rayls-sovereign-relayer/dvp/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

type ReceiverConfig struct {
	ChainID       *big.Int
	ViewPublicKey []byte
}

type receiverSwapRepository interface {
	CreateSwap(ctx context.Context, swap *types.DvpSwap) error
	GetSwapBySharedID(ctx context.Context, sharedID string) (*types.DvpSwap, error)
	UpdateSwapStatus(ctx context.Context, sharedID string, status types.DvpSwapStatus) error
	CancelSwap(ctx context.Context, sharedID string, status types.DvpSwapStatus) error
}

type receiverDepositRepository interface {
	CreateDeposit(ctx context.Context, deposit *types.DvpDeposit) error
	UpdateDepositStatus(ctx context.Context, commitment *big.Int, status types.DvpDepositStatus) error
	BatchUpdateStatusForCommitments(ctx context.Context, commitments []string, status types.DvpDepositStatus) error
	UpsertDepositNullifier(ctx context.Context, commitment *big.Int, nullifier *big.Int) error
	BatchUpsertNullifiers(ctx context.Context, commitmentNullifierMap map[string]string) error
	GetNonFungibleDeposit(
		ctx context.Context,
		tokenID string,
		tokenAddress string,
		userAddress string,
		tokenType types.DvpTokenType,
		status types.DvpDepositStatus,
	) (*types.DvpDeposit, error)
	GetDepositByCommitmentAndStatus(ctx context.Context, commitment *big.Int, status types.DvpDepositStatus) (*types.DvpDeposit, error)
	ConfirmDeposit(ctx context.Context, commitment *big.Int, treeNumber *big.Int) error
	UpdateDepositStatusByNullifier(ctx context.Context, tokenAddress string, nullifier *big.Int, status types.DvpDepositStatus) error
}

type receiverMerkleClient interface {
	PopulateMerkleDbTree(
		ctx context.Context,
		tokenAddress string,
		bigIntTokenType *big.Int,
		bigIntTreeNumber *big.Int,
		bigIntLeaves []*big.Int,
	) error
	GetNonFungibleDeposit(
		ctx context.Context,
		tokenID string,
		tokenAddress string,
		userAddress string,
		tokenType types.DvpTokenType,
		status types.DvpDepositStatus,
	) (*types.DvpDeposit, error)
}

type receiverEndpointClient interface {
	GetResourceAddress(ctx context.Context, resourceId string) (common.Address, error)
}

type receiverERC721HandlerClient interface {
	NotifySenderWithPNCommunicator(
		ctx context.Context,
		tokenAddress common.Address,
		sharedId string,
		status types.DvpCommunicatorStatus,
		message string,
	) error
	MarkSwapCompleted(
		ctx context.Context,
		tokenAddress common.Address,
		to common.Address,
		sourceChainId *big.Int,
		destChainId *big.Int,
		resourceId string,
		tokenId string,
		amount *big.Int,
		sharedId string,
	) error
}

type receiverERC1155HandlerClient interface {
	NotifySenderWithPNCommunicator(
		ctx context.Context,
		tokenAddress common.Address,
		sharedId string,
		status types.DvpCommunicatorStatus,
		message string,
	) error
	MarkSwapCompleted(
		ctx context.Context,
		tokenAddress common.Address,
		from common.Address,
		to common.Address,
		sourceChainId *big.Int,
		destChainId *big.Int,
		resourceId string,
		tokenId string,
		amount *big.Int,
		data []byte,
		sharedId string,
	) error
}

type receiverEnygmaHandlerClient interface {
	NotifySenderWithPNCommunicator(
		ctx context.Context,
		tokenAddress common.Address,
		sharedId string,
		status types.DvpCommunicatorStatus,
		message string,
	) error
	MarkSwapCompleted(ctx context.Context, tokenAddress common.Address, chainID *big.Int, sharedId string) error
}

type receiverDepositFinder interface {
	FindEnygmaDeposits(ctx context.Context, from string, tokenInAddress string, tokenInAmount *big.Int) ([]*types.DvpDeposit, error)
	FindERC1155DepositsForJSProof(
		ctx context.Context,
		userAddress string,
		tokenAddress string,
		tokenId string,
		paymentAmount *big.Int,
	) ([]*types.DvpDeposit, error)
}

type receiverDepositConsolidator interface {
	PrepareDepositsForJSProof(ctx context.Context, chainEventID string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error)
}

type receiverCommitmentCalculator interface {
	CalculateNFTCommitment(spendPK, salt *big.Int, nftID string, nftAddress string) (*big.Int, error)
	CalculateERC1155Commitment(spendPK, salt *big.Int, tokenAddress string, tokenID string, tokenAmount *big.Int) (*big.Int, error)
	CalculatePaymentCommitment(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error)
	CalculateNullifier(spendSK, leafIndex *big.Int) (*big.Int, error)
}

type receiverProofGenerator interface {
	GenerateEnygmaJSProof(
		ctx context.Context,
		sourceViewPublicKey []byte,
		nftCommitment *big.Int,
		destinationPaymentPublicKey *big.Int,
		destinationSalt *big.Int,
		paymentAmount *big.Int,
		tokenAddress string,
		deposits []*types.DvpDeposit,
	) (*dvp.ProofReceipt, error)
	GenerateERC1155JSProof(
		ctx context.Context,
		sourceViewPublicKey []byte,
		paymentCommitment *big.Int,
		destinationPaymentPublicKey *big.Int,
		destinationSalt *big.Int,
		paymentAmount *big.Int,
		userAddress string,
		tokenAddress string,
		tokenId string,
		deposits []*types.DvpDeposit,
	) (*dvp.ProofReceipt, error)
	GenerateOwnershipProof(
		ctx context.Context,
		sourceViewPublicKey []byte,
		paymentCommitment *big.Int,
		destinationPaymentPublicKey *big.Int,
		destinationSalt *big.Int,
		deposit *types.DvpDeposit,
	) (*dvp.ProofReceipt, error)
}

type receiverTxPersister interface {
	PersistAndBroadcast(ctx context.Context, tx *ethTypes.Transaction, recoveryData types.TxRecoveryData) error
	CheckPendingRecovery(ctx context.Context, resourceID string, blockNumber uint64, fromChainID string, eventType types.EnygmaEventType) (*types.TxRecoveryData, error)
	ResumePendingTx(ctx context.Context, pending *types.TxRecoveryData) error
}

var _ receiverTxPersister = (*dvpService.TxPersister)(nil)

type receiverRetryService interface {
	RetryOperation(
		ctx context.Context,
		operationName string,
		maxRetries int,
		executeOperation func(ctx context.Context) error,
	) error
}

type receiverTxManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type DvpReceiver struct {
	conf                 ReceiverConfig
	swapRepository       receiverSwapRepository
	depositRepository    receiverDepositRepository
	plEndpointClient     receiverEndpointClient
	enygmaHandlerClient  receiverEnygmaHandlerClient
	erc721HandlerClient  receiverERC721HandlerClient
	erc1155HandlerClient receiverERC1155HandlerClient
	depositFinder        receiverDepositFinder
	depositConsolidator  receiverDepositConsolidator
	commitmentCalculator receiverCommitmentCalculator
	proofGen             receiverProofGenerator
	merkleClient         receiverMerkleClient
	retryService         receiverRetryService
	txManager            receiverTxManager
}

func NewDvpReceiver(
	conf ReceiverConfig,
	swapRepository receiverSwapRepository,
	depositRepository receiverDepositRepository,
	plEndpointClient receiverEndpointClient,
	enygmaHandlerClient receiverEnygmaHandlerClient,
	erc721HandlerClient receiverERC721HandlerClient,
	erc1155HandlerClient receiverERC1155HandlerClient,
	depositFinder receiverDepositFinder,
	depositConsolidator receiverDepositConsolidator,
	commitmentCalculator receiverCommitmentCalculator,
	proofGenerator receiverProofGenerator,
	merkleClient receiverMerkleClient,
	retryService receiverRetryService,
	txManager receiverTxManager,
) *DvpReceiver {
	return &DvpReceiver{
		conf:                 conf,
		swapRepository:       swapRepository,
		depositRepository:    depositRepository,
		plEndpointClient:     plEndpointClient,
		enygmaHandlerClient:  enygmaHandlerClient,
		erc721HandlerClient:  erc721HandlerClient,
		erc1155HandlerClient: erc1155HandlerClient,
		depositFinder:        depositFinder,
		depositConsolidator:  depositConsolidator,
		commitmentCalculator: commitmentCalculator,
		proofGen:             proofGenerator,
		merkleClient:         merkleClient,
		retryService:         retryService,
		txManager:            txManager,
	}
}

func (r *DvpReceiver) HandleCommitments(ctx context.Context, data *types.DvpCommitmentsData) error {
	slog.Info("Handling commitments",
		slog.String("tokenAddress", data.TokenAddress),
		slog.String("treeNumber", data.TreeNumber.String()),
		slog.Int("commitmentCount", len(data.Commitments)),
	)

	err := r.merkleClient.PopulateMerkleDbTree(ctx, data.TokenAddress, data.TokenType, data.TreeNumber, data.Commitments)
	if err != nil {
		return fmt.Errorf("failed to populate merkle tree: %w", err)
	}

	for _, commitment := range data.Commitments {
		deposit, err := r.depositRepository.GetDepositByCommitmentAndStatus(ctx, commitment, types.DvpDepositPending)
		if err != nil {
			slog.Error("Failed to get deposit by commitment",
				slog.String("commitment", commitment.String()),
				slog.Any("error", err),
			)
			continue
		}

		if deposit == nil {
			continue
		}

		err = r.depositRepository.ConfirmDeposit(ctx, commitment, data.TreeNumber)
		if err != nil {
			slog.Error("Failed to confirm deposit",
				slog.String("commitment", commitment.String()),
				slog.Any("error", err),
			)
			continue
		}

		slog.Info("Deposit confirmed",
			slog.String("commitment", commitment.String()),
			slog.String("tokenAddress", data.TokenAddress),
			slog.String("treeNumber", data.TreeNumber.String()),
		)
	}

	return nil
}

func (r *DvpReceiver) HandleNullifiers(ctx context.Context, tokenAddress string, nullifiers []*big.Int) error {
	slog.Info("Handling nullifiers",
		slog.String("tokenAddress", tokenAddress),
		slog.Int("count", len(nullifiers)),
	)

	for _, nullifier := range nullifiers {
		err := r.depositRepository.UpdateDepositStatusByNullifier(ctx, tokenAddress, nullifier, types.DvpDepositSpent)
		if err != nil {
			return fmt.Errorf("failed to update deposit status by nullifier %s on token %s: %w", nullifier.String(), tokenAddress, err)
		}

		slog.Info("Deposit marked as spent",
			slog.String("tokenAddress", tokenAddress),
			slog.String("nullifier", nullifier.String()),
		)
	}

	return nil
}

func (r *DvpReceiver) HandleSwapCompleted(ctx context.Context, sharedId string) error {
	slog.Info("Handling swap completed",
		slog.String("sharedId", sharedId),
	)

	swap, err := r.swapRepository.GetSwapBySharedID(ctx, sharedId)
	if err != nil {
		return err
	}
	if swap == nil {
		return fmt.Errorf("swap for shared ID %s does not exist", sharedId)
	}
	err = r.swapRepository.UpdateSwapStatus(ctx, sharedId, types.DvpSwapCompleted)
	if err != nil {
		return fmt.Errorf("updating swap status to completed: %w", err)
	}

	to := common.HexToAddress(swap.To)
	from := common.HexToAddress(swap.From)

	tokenHandlerAddress, err := r.plEndpointClient.GetResourceAddress(ctx, swap.TokenInResourceID)
	if err != nil {
		return err
	}

	var swapErr error

	switch swap.TokenInType {
	case types.DvpEnygma:
		slog.Info("Completing Enygma swap", slog.String("SharedId", sharedId))
		swapErr = r.enygmaHandlerClient.MarkSwapCompleted(ctx, tokenHandlerAddress, swap.DestChainID, sharedId)

	case types.DvpERC721:
		slog.Info("Completing ERC721 swap", slog.String("SharedId", sharedId))

		// For ERC721, token amount is always 1
		tokenAmount := big.NewInt(1)
		swapErr = r.erc721HandlerClient.MarkSwapCompleted(ctx, tokenHandlerAddress, to, swap.SourceChainID, swap.DestChainID, swap.TokenInResourceID, swap.TokenInID, tokenAmount, sharedId)

	case types.DvpERC1155:
		slog.Info("Completing ERC1155 swap", slog.String("SharedId", sharedId))
		swapErr = r.erc1155HandlerClient.MarkSwapCompleted(ctx, tokenHandlerAddress, from, to, swap.SourceChainID, swap.DestChainID, swap.TokenInResourceID, swap.TokenInID, swap.TokenInAmount, []byte{}, sharedId)

	default:
		swapErr = fmt.Errorf("failed to complete swap. Unsupported token type: %d, SharedId=%s", swap.TokenInType, sharedId)
	}

	if swapErr != nil {
		return swapErr
	}

	slog.Info("Swap completed handled successfully", slog.String("SharedId", sharedId))

	return nil
}

func (r *DvpReceiver) HandleSwapInitiated(ctx context.Context, blockNum uint64, data *types.DvpSwapInitiatedData) error {
	slog.Info("HandleSwapInitiated: processing",
		slog.String("sharedId", data.Message.SharedId),
	)

	cancelPreimage, err := cryptography.GetPoseidonHash([]*big.Int{data.InitiatorDestSalt})
	if err != nil {
		return fmt.Errorf("computing cancel preimage: %w", err)
	}

	newSwap := &types.DvpSwap{
		SharedID:      data.Message.SharedId,
		To:            data.Message.To,
		SourceChainID: r.conf.ChainID,
		DestChainID:   data.Message.ChainId,

		// Swap legs are mirrored on the receiver: the token the initiator wants to
		// receive (msg.TokenOut) is what we deposit, so it maps to our TokenIn.
		TokenInAmount:     data.Message.TokenOutAmount,
		TokenInAddress:    data.Message.TokenOutAddress,
		TokenInResourceID: data.Message.TokenOutResourceID,
		TokenInType:       data.Message.TokenOutType,
		TokenInID:         data.Message.TokenOutID,
		// Conversely, the token the initiator deposits (msg.TokenIn) is what we
		// receive, so it maps to our TokenOut.
		TokenOutAmount:     data.Message.TokenInAmount,
		TokenOutAddress:    data.Message.TokenInAddress,
		TokenOutID:         data.Message.TokenInID,
		TokenOutResourceID: data.Message.TokenInResourceID,
		TokenOutType:       data.Message.TokenInType,
		Status:             types.DvpSwapWaitingConfirmation,
		// InitiatorSelfSalt is the salt the initiator used for their self-destination
		// commitment — the commitment of the token the initiator will eventually
		// receive. From our perspective, that token is what we deposit (our outgoing
		// leg, our TokenIn). We store it as DestSalt because, in the same naming
		// convention applied from our side, it is the salt for the commitment of the
		// token at the destination — i.e. the counterparty's incoming leg.
		DestSalt: data.Message.InitiatorSelfSalt,
		// InitiatorDestSalt is the salt the initiator used for the destination
		// commitment — the commitment of the token the destination party (us) will
		// receive, i.e. the initiator's own deposit. From our perspective that token
		// is what we receive (our incoming leg, our TokenOut). We store it as SelfSalt
		// because it is the salt for our own self-destination commitment — the token
		// we will eventually receive.
		SelfSalt:       data.InitiatorDestSalt,
		CancelPreimage: cancelPreimage,
	}

	err = r.swapRepository.CreateSwap(ctx, newSwap)
	if err != nil {
		return fmt.Errorf("creating swap: %w", err)
	}

	slog.Info("Swap initiated to us handled successfully",
		slog.String("sharedId", data.Message.SharedId),
		slog.String("sourceChainId", newSwap.SourceChainID.String()),
		slog.String("destChainId", newSwap.DestChainID.String()),
	)

	return nil
}

func (r *DvpReceiver) HandleSwapRevert(
	ctx context.Context,
	sharedId string,
	status types.DvpSwapStatus,
) error {
	if status != types.DvpSwapCancelled && status != types.DvpSwapTimedOut {
		return fmt.Errorf("invalid revert status for swap: %v", status)
	}

	swap, err := r.swapRepository.GetSwapBySharedID(ctx, sharedId)
	if err != nil {
		return fmt.Errorf("getting swap by shared ID: %w", err)
	}

	if swap == nil {
		return fmt.Errorf("swap not found")
	}

	err = r.swapRepository.CancelSwap(ctx, sharedId, status)
	if err != nil {
		return fmt.Errorf("cancelling swap in repository: %w", err)
	}

	err = r.notifyCancellation(ctx, sharedId, swap.TokenInResourceID, swap.TokenInType, status)
	if err != nil {
		return fmt.Errorf("notifying communicator for swap cancellation: %w", err)
	}

	return nil
}

func (r *DvpReceiver) notifyCancellation(
	ctx context.Context,
	sharedId string,
	resourceId string,
	tokenType types.DvpTokenType,
	swapStatus types.DvpSwapStatus,
) error {
	tokenAddress, err := r.plEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("getting resource address: %w", err)
	}

	message := "swap is expired"
	status := types.SwapTimedOut
	if swapStatus == types.DvpSwapCancelled {
		message = "swap is cancelled"
		status = types.SwapCancelled
	}

	switch tokenType {
	case types.DvpEnygma:
		return r.enygmaHandlerClient.NotifySenderWithPNCommunicator(ctx, tokenAddress, sharedId, status, message)

	case types.DvpERC721:
		return r.erc721HandlerClient.NotifySenderWithPNCommunicator(ctx, tokenAddress, sharedId, status, message)

	case types.DvpERC1155:
		return r.erc1155HandlerClient.NotifySenderWithPNCommunicator(ctx, tokenAddress, sharedId, status, message)

	default:
		return fmt.Errorf("failed to notify sender. Unsupported token type: %d, SharedId=%s", tokenType, sharedId)
	}
}
