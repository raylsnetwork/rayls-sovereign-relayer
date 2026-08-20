package handler

//go:generate moq --skip-ensure --pkg handler_test -out initiator_mock_test.go . initiatorSwapRepository initiatorParticipantStorageClient initiatorDepositRepository initiatorKeysClient initiatorEndpointClient InitiatorERC721Client InitiatorERC721HandlerClient InitiatorERC1155Client InitiatorERC1155HandlerClient initiatorDvpClient initiatorEncryptor initiatorDepositFinder initiatorDepositConsolidator initiatorCommitmentCalculator initiatorProofGenerator initiatorSwapWaiter initiatorSwapAgreement initiatorTxPersister initiatorEthClient initiatorTxManager

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	ps "github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/ParticipantStorageV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cryptography"
	keyspb "github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/keys"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp"
	dvpService "github.com/raylsnetwork/rayls-privacy-relayer-api/dvp/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"google.golang.org/grpc"
)

type InitiatorConfig struct {
	ChainID            *big.Int
	DvpContractAddress common.Address
	DvpOperatorAddress common.Address
}

type initiatorSwapRepository interface {
	GetSwapBySharedID(ctx context.Context, sharedId string) (*types.DvpSwap, error)
	CreateSwap(ctx context.Context, swap *types.DvpSwap) error
	UpdateSwapStatus(ctx context.Context, sharedID string, status types.DvpSwapStatus) error
	UpdateSwapFrom(ctx context.Context, sharedID string, from string) error
}

type initiatorParticipantStorageClient interface {
	GetPaymentSpendPublicKey(ctx context.Context, chainID *big.Int) (*big.Int, error)
	GetChainViewData(ctx context.Context, chainID *big.Int, blockNumber *big.Int) (ps.ParticipantStructsPrivacyNodeViewData, error)
}

type initiatorDepositRepository interface {
	CreateDeposit(ctx context.Context, deposit *types.DvpDeposit) error
	UpdateDepositStatus(ctx context.Context, commitment *big.Int, status types.DvpDepositStatus) error
	BatchUpdateStatusForCommitments(ctx context.Context, commitments []string, status types.DvpDepositStatus) error
	BatchUpsertNullifiers(ctx context.Context, commitmentNullifierMap map[string]string) error
	UpsertDepositNullifier(ctx context.Context, commitment *big.Int, nullifier *big.Int) error
	GetNonFungibleDeposit(ctx context.Context, tokenId string, tokenAddress string, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) (*types.DvpDeposit, error)
}

type initiatorKeysClient interface {
	GetPaymentSpendKey(ctx context.Context, in *keyspb.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error)
	GetViewPublicKey(ctx context.Context, in *keyspb.GetViewPublicKeyRequest, opts ...grpc.CallOption) (*keyspb.GetViewPublicKeyResponse, error)
}

type initiatorEndpointClient interface {
	GetResourceAddress(ctx context.Context, resourceId string) (common.Address, error)
}

type InitiatorERC721Client interface {
	Approve(ctx context.Context, chainEventID string, tokenAddress common.Address, to common.Address, nftId *big.Int) error
	Burn(ctx context.Context, chainEventID string, tokenAddress common.Address, nftId *big.Int) error
	UpdateExtraData(ctx context.Context, chainEventID string, tokenAddress common.Address, nftId *big.Int, chainId *big.Int, extraDataBytes []byte, newOwner common.Address) error
	MintBatch(ctx context.Context, mintDatas []*dvp.DvpERC721MintData) (map[string]contractclient.BatchResult, error)
}

type InitiatorERC721HandlerClient interface {
	GetTotalSupply(ctx context.Context, tokenAddress common.Address) ([]*big.Int, error)
	GetExtraData(ctx context.Context, tokenAddress common.Address, nftId *big.Int) ([]byte, error)
	Unlock(ctx context.Context, chainEventID string, tokenAddress common.Address, nftId *big.Int) error
}

type InitiatorERC1155Client interface {
	Burn(ctx context.Context, chainEventID string, tokenAddress common.Address, tokenOwner common.Address, tokenId *big.Int, tokenAmount *big.Int) error
	Approve(ctx context.Context, chainEventID string, tokenAddress common.Address, to common.Address) error
	UpdateExtraData(ctx context.Context, chainEventID string, tokenAddress common.Address, tokenId *big.Int, tokenAmount *big.Int, chainId *big.Int, extraDataBytes []byte, newOwner common.Address) error
	MintBatch(ctx context.Context, mintDatas []*dvp.DvpERC1155MintData) (map[string]contractclient.BatchResult, error)
}

type InitiatorERC1155HandlerClient interface {
	GetAllTokenIdsWithSupply(ctx context.Context, tokenAddress common.Address) ([]dvp.DvpERC1155Supply, error)
	Unlock(ctx context.Context, chainEventID string, tokenAddress common.Address, tokenId *big.Int, tokenAmount *big.Int, to common.Address) error
	GetTokenExtraData(ctx context.Context, tokenAddress common.Address, tokenId *big.Int) ([]byte, error)
}

type initiatorDvpClient interface {
	DepositERC721(ctx context.Context, chainEventID string, contractAddress common.Address, nftId *big.Int, publicKey *big.Int, salt *big.Int, encryptedBalanceUpdate []byte) error
	WithdrawERC721(ctx context.Context, chainEventID string, contractAddress common.Address, nftId *big.Int, to common.Address, salt *big.Int, proof *dvp.ProofReceipt, encryptedBalanceUpdate []byte) error
	// SignWithdrawERC721(ctx context.Context, contractAddress common.Address, nftId *big.Int, to common.Address, salt *big.Int, proof *dvp.ProofReceipt, encryptedBalanceUpdate []byte) (*ethTypes.Transaction, error)
	DepositERC1155(ctx context.Context, chainEventID string, contractAddress common.Address, tokenId *big.Int, tokenAmount *big.Int, tokenData []byte, publicKey *big.Int, salt *big.Int, encryptedBalanceUpdate []byte) error
	WithdrawERC1155(ctx context.Context, chainEventID string, contractAddress common.Address, tokenId *big.Int, tokenAmount *big.Int, to common.Address, salt *big.Int, proof *dvp.ProofReceipt, encryptedBalanceUpdate []byte) error
	// WithdrawERC1155(ctx context.Context, contractAddress common.Address, tokenId *big.Int, tokenAmount *big.Int, to common.Address, salt *big.Int, proof *dvp.ProofReceipt, encryptedBalanceUpdate []byte) (*ethTypes.Transaction, error)

	InitiateSwap(ctx context.Context, salt *big.Int, ciphertext []byte, msg *types.DvpSwapMessage, proof *dvp.ProofReceipt, validityTime uint64, passphrase *big.Int) error
	CompleteSwap(ctx context.Context, salt *big.Int, swap *types.DvpSwapMessage, proof *dvp.ProofReceipt) error
	CancelSwap(ctx context.Context, sharedId string, preimage *big.Int) error
}

type initiatorEncryptor interface {
	EncryptDvpBalanceUpdated(ctx context.Context, message types.DvpBalanceUpdated) ([]byte, error)
}

type initiatorDepositFinder interface {
	FindEnygmaDeposits(ctx context.Context, from string, tokenInAddress string, tokenInAmount *big.Int) ([]*types.DvpDeposit, error)
	FindERC1155DepositsForJSProof(ctx context.Context, userAddress string, tokenAddress string, tokenId string, paymentAmount *big.Int) ([]*types.DvpDeposit, error)
	FindERC721Deposit(ctx context.Context, userAddress string, tokenAddress string, tokenId string) (*types.DvpDeposit, error)
}

type initiatorDepositConsolidator interface {
	PrepareDepositsForJSProof(ctx context.Context, id string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error)
	ConsolidateERC1155Deposits(ctx context.Context, id string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit, consolidationAmount *big.Int) ([]*types.DvpDeposit, error)
}

type initiatorCommitmentCalculator interface {
	CalculateNFTCommitment(spendPK, salt *big.Int, nftID string, nftAddress string) (*big.Int, error)
	CalculateERC1155Commitment(spendPK, salt *big.Int, tokenAddress string, tokenID string, tokenAmount *big.Int) (*big.Int, error)
	CalculatePaymentCommitment(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error)
	CalculateNullifier(spendSK, leafIndex *big.Int) (*big.Int, error)
}

type initiatorProofGenerator interface {
	GenerateEnygmaToERC721SwapProof(ctx context.Context, swap *types.DvpSwap, deposits []*types.DvpDeposit, sourceViewPublicKey []byte, selfSalt *big.Int, destSalt *big.Int, destSpendPubKey *big.Int) (*dvp.ProofReceipt, error)
	GenerateEnygmaToERC1155SwapProof(ctx context.Context, swap *types.DvpSwap, deposits []*types.DvpDeposit, sourceViewPublicKey []byte, selfSalt *big.Int, destSalt *big.Int, destSpendPubKey *big.Int) (*dvp.ProofReceipt, error)
	GenerateERC721ToEnygmaSwapProof(ctx context.Context, swap *types.DvpSwap, deposit *types.DvpDeposit, sourceViewPublicKey []byte, selfSalt *big.Int, destSalt *big.Int, destSpendPubKey *big.Int) (*dvp.ProofReceipt, error)
	GenerateERC1155ToEnygmaSwapProof(ctx context.Context, swap *types.DvpSwap, deposits []*types.DvpDeposit, sourceViewPublicKey []byte, selfSalt *big.Int, destSalt *big.Int, destSpendPubKey *big.Int) (*dvp.ProofReceipt, error)
	GenerateERC721WithdrawProof(ctx context.Context, sourceViewPublicKey []byte, destSalt *big.Int, operatorPublicKey *big.Int, deposit *types.DvpDeposit) (*dvp.ProofReceipt, error)
	GenerateERC1155WithdrawProof(ctx context.Context, sourceViewPublicKey []byte, destSalt *big.Int, operatorPublicKey *big.Int, userAddress string, tokenAddress string, tokenID string, tokenAmount *big.Int, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error)
}
type initiatorSwapWaiter interface {
	WaitForSwapInitiation(ctx context.Context, sharedId string) (*types.DvpSwap, error)
}

type initiatorSwapAgreement interface {
	Verify(ctx context.Context, swap *types.DvpSwap, destChainId *big.Int, tokenInResourceId string, tokenInAmount *big.Int, tokenInId string, tokenInType types.DvpTokenType, tokenOutResourceId string, tokenOutAmount *big.Int, tokenOutId string, tokenOutType types.DvpTokenType) (string, error)
	HandleSwapDisagreement(ctx context.Context, sharedId string, tokenInResourceId string, tokenInType types.DvpTokenType, message string) error
}

type initiatorTxPersister interface {
	PersistAndBroadcast(ctx context.Context, tx *ethTypes.Transaction, recoveryData types.TxRecoveryData) error
	CheckPendingRecovery(ctx context.Context, resourceID string, blockNumber uint64, fromChainID string, eventType types.EnygmaEventType) (*types.TxRecoveryData, error)
	ResumePendingTx(ctx context.Context, pending *types.TxRecoveryData) error
}

type initiatorEthClient interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*ethTypes.Block, error)
}

type initiatorTxManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

var _ initiatorSwapRepository = (*repository.DvpSwapRepository)(nil)
var _ initiatorDepositRepository = (*repository.DvpDepositRepository)(nil)
var _ initiatorEndpointClient = (*contractclient.EndpointClient)(nil)
var _ initiatorDvpClient = (*contractclient.DvpClient)(nil)
var _ initiatorDepositFinder = (*dvpService.DepositFinder)(nil)
var _ initiatorDepositConsolidator = (*dvpService.ConsolidationService)(nil)
var _ initiatorCommitmentCalculator = (*dvpService.CommitmentCalculator)(nil)
var _ initiatorProofGenerator = (*dvpService.ProofService)(nil)
var _ initiatorSwapWaiter = (*dvpService.SwapWaiter)(nil)
var _ initiatorSwapAgreement = (*dvpService.SwapAgreement)(nil)
var _ initiatorEncryptor = (*contractclient.Encryptor)(nil)
var _ initiatorTxManager = (*repository.TransactionManager)(nil)
var _ initiatorParticipantStorageClient = (*contractclient.ParticipantStorageClient)(nil)

type DvpInitiator struct {
	conf                 InitiatorConfig
	swapRepository       initiatorSwapRepository
	depositRepository    initiatorDepositRepository
	psClient             initiatorParticipantStorageClient
	keysClient           initiatorKeysClient
	ccEndpointClient     initiatorEndpointClient
	plEndpointClient     initiatorEndpointClient
	erc721Client         InitiatorERC721Client
	erc721HandlerClient  InitiatorERC721HandlerClient
	erc1155Client        InitiatorERC1155Client
	erc1155HandlerClient InitiatorERC1155HandlerClient
	dvpClient            initiatorDvpClient
	encryptor            initiatorEncryptor
	depositFinder        initiatorDepositFinder
	depositConsolidator  initiatorDepositConsolidator
	commitmentCalculator initiatorCommitmentCalculator
	proofGenerator       initiatorProofGenerator
	swapWaiter           initiatorSwapWaiter
	swapAgreement        initiatorSwapAgreement
	nodeClient           initiatorEthClient
	hubClient            initiatorEthClient
	txManager            initiatorTxManager
}

func NewDvpInitiator(
	conf InitiatorConfig,
	swapRepository initiatorSwapRepository,
	depositRepository initiatorDepositRepository,
	psClient initiatorParticipantStorageClient,
	keysClient initiatorKeysClient,
	ccEndpointClient initiatorEndpointClient,
	plEndpointClient initiatorEndpointClient,
	erc721Client InitiatorERC721Client,
	erc721HandlerClient InitiatorERC721HandlerClient,
	erc1155Client InitiatorERC1155Client,
	erc1155HandlerClient InitiatorERC1155HandlerClient,
	dvpClient initiatorDvpClient,
	encryptor initiatorEncryptor,
	depositFinder initiatorDepositFinder,
	depositConsolidator initiatorDepositConsolidator,
	commitmentCalculator initiatorCommitmentCalculator,
	proofGenerator initiatorProofGenerator,
	swapWaiter initiatorSwapWaiter,
	swapAgreement initiatorSwapAgreement,
	nodeClient initiatorEthClient,
	hubClient initiatorEthClient,
	txManager initiatorTxManager,
) *DvpInitiator {
	return &DvpInitiator{
		conf:                 conf,
		swapRepository:       swapRepository,
		depositRepository:    depositRepository,
		psClient:             psClient,
		keysClient:           keysClient,
		ccEndpointClient:     ccEndpointClient,
		plEndpointClient:     plEndpointClient,
		erc721Client:         erc721Client,
		erc721HandlerClient:  erc721HandlerClient,
		erc1155Client:        erc1155Client,
		erc1155HandlerClient: erc1155HandlerClient,
		dvpClient:            dvpClient,
		encryptor:            encryptor,
		depositFinder:        depositFinder,
		depositConsolidator:  depositConsolidator,
		commitmentCalculator: commitmentCalculator,
		proofGenerator:       proofGenerator,
		swapWaiter:           swapWaiter,
		swapAgreement:        swapAgreement,
		nodeClient:           nodeClient,
		hubClient:            hubClient,
		txManager:            txManager,
	}
}

/* DVP ENYGMA SWAPs */
func (s *DvpInitiator) HandleEnygmaSwapERC721(
	ctx context.Context,
	sharedId string,
	toChainId *big.Int,
	from common.Address,
	enygmaResourceId string,
	enygmaAmount *big.Int,
	nftResourceId string,
	nftId string,
	txHash string,
	txBlockNumber *big.Int,
	validityTime uint64,
) error {
	slog.Info("Handling Enygma -> ERC721 swap", slog.String("sharedId", sharedId))

	swap, err := s.swapRepository.GetSwapBySharedID(ctx, sharedId)
	if err != nil {
		return fmt.Errorf("getting swap by shared ID: %w", err)
	}

	// Restore block timestamp to update Privacy Network Hub
	block, err := s.nodeClient.BlockByNumber(ctx, txBlockNumber)
	if err != nil {
		return fmt.Errorf("getting block by number: %w", err)
	}
	txTimestamp := time.Unix(int64(block.Time()), 0).UTC()

	// We are the initiator of the swap
	if swap == nil {
		enygmaId := ""
		nftAmount := big.NewInt(1)

		swap, err = s.createSwap(
			ctx,
			sharedId,
			toChainId,
			from,
			enygmaResourceId,
			enygmaAmount,
			enygmaId,
			types.DvpEnygma,
			nftResourceId,
			nftAmount,
			nftId,
			types.DvpERC721,
		)

		if err != nil {
			return fmt.Errorf("creating Enygma -> ERC721 swap: %w", err)
		}
	}

	if swap.Status == types.DvpSwapCreated {
		slog.Debug("Initiating new Enygma -> ERC721 swap", slog.String("sharedId", sharedId))

		params, err := s.prepareSwapProofParams(ctx, swap.DestChainID)
		if err != nil {
			return fmt.Errorf("preparing Enygma -> ERC721 proof params: %w", err)
		}

		userDeposits, err := s.depositFinder.FindEnygmaDeposits(ctx, swap.From, swap.TokenInAddress, swap.TokenInAmount)
		if err != nil {
			return fmt.Errorf("finding enygma deposits: %w", err)
		}

		jsDeposits, err := s.depositConsolidator.PrepareDepositsForJSProof(ctx, sharedId, params.selfViewPubKey, userDeposits)
		if err != nil {
			return fmt.Errorf("preparing deposits for JS proof: %w", err)
		}

		proof, err := s.proofGenerator.GenerateEnygmaToERC721SwapProof(ctx, swap, jsDeposits, params.selfViewPubKey, params.selfSalt, params.destSalt, params.destSpendPubKey)
		if err != nil {
			return fmt.Errorf("generating Enygma -> ERC721 proof: %w", err)
		}

		err = s.initiateSwap(ctx, jsDeposits, swap, validityTime, params, proof, txHash, txTimestamp)
		if err != nil {
			return fmt.Errorf("initiating Enygma -> ERC721 swap: %w", err)
		}

		return nil
	}

	// The swap was already initiated by the other side, so we need to confirm and complete it.
	if swap.Status == types.DvpSwapWaitingConfirmation {
		slog.Debug("Confrming new Enygma -> ERC721 swap", slog.String("sharedId", sharedId))

		swap.From = from.Hex()
		if err := s.swapRepository.UpdateSwapFrom(ctx, sharedId, swap.From); err != nil {
			return fmt.Errorf("updating swap from_address for shared_id %s: %w", sharedId, err)
		}

		// TODO: This is redundant since we'll match commitments instead
		// Ensure both parties agree on the same swap information
		nftAmount := big.NewInt(1)
		enygmaId := ""
		if reason, agreementErr := s.swapAgreement.Verify(ctx, swap, toChainId, enygmaResourceId, enygmaAmount, enygmaId, types.DvpEnygma, nftResourceId, nftAmount, nftId, types.DvpERC721); agreementErr != nil {
			slog.Error("Swap Enygma -> ERC721 does not match the expected swap information", slog.String("sharedId", sharedId), slog.String("err", agreementErr.Error()), slog.String("reason", reason))

			err := s.swapAgreement.HandleSwapDisagreement(ctx, sharedId, enygmaResourceId, types.DvpEnygma, agreementErr.Error())
			if err != nil {
				return fmt.Errorf("handling swap disagreement: %w", err)
			}

			return nil
		}

		params, err := s.prepareSwapProofParams(ctx, swap.DestChainID)
		if err != nil {
			return fmt.Errorf("preparing Enygma -> ERC721 proof params: %w", err)
		}

		userDeposits, err := s.depositFinder.FindEnygmaDeposits(ctx, swap.From, swap.TokenInAddress, swap.TokenInAmount)
		if err != nil {
			return fmt.Errorf("finding enygma deposits: %w", err)
		}

		jsDeposits, err := s.depositConsolidator.PrepareDepositsForJSProof(ctx, sharedId, params.selfViewPubKey, userDeposits)
		if err != nil {
			return fmt.Errorf("preparing deposits for JS proof: %w", err)
		}

		proof, err := s.proofGenerator.GenerateEnygmaToERC721SwapProof(ctx, swap, jsDeposits, params.selfViewPubKey, swap.SelfSalt, swap.DestSalt, params.destSpendPubKey)
		if err != nil {
			return fmt.Errorf("generating Enygma -> ERC721 proof: %w", err)
		}

		slog.Info("Swap Enygma -> ERC721 sending confirmation", slog.String("sharedId", sharedId))

		msg := &types.DvpSwapMessage{
			SharedId:       sharedId,
			To:             swap.From,
			PNTxHash:       txHash,
			PNTxTimestamp:  txTimestamp,
			ChainId:        s.conf.ChainID,
			TokenInType:    swap.TokenInType,
			TokenInAddress: swap.TokenInAddress,
		}

		err = s.dvpClient.CompleteSwap(ctx, swap.DestSalt, msg, proof)
		if err != nil {
			return fmt.Errorf("confirming Enygma -> ERC721 swap: %w", err)
		}

		slog.Info("Swap Enygma -> ERC721 completed successfully", slog.String("sharedId", sharedId))
		return nil
	}

	return nil
}

func (s *DvpInitiator) initiateSwap(
	ctx context.Context,
	deposits []*types.DvpDeposit,
	swap *types.DvpSwap,
	swapValidity uint64,
	proofParams *swapProofParams,
	proof *dvp.ProofReceipt,
	txHash string,
	txTimestamp time.Time,
) error {
	slog.Debug("Initiating new swap", slog.String("sharedId", swap.SharedID))

	msg := types.DvpSwapMessage{
		SharedId:           swap.SharedID,
		To:                 swap.From,
		ChainId:            swap.SourceChainID,
		TokenInResourceID:  swap.TokenInResourceID,
		TokenInAddress:     swap.TokenInAddress,
		TokenInAmount:      swap.TokenInAmount,
		TokenInType:        swap.TokenInType,
		TokenInID:          swap.TokenInID,
		TokenOutResourceID: swap.TokenOutResourceID,
		TokenOutAddress:    swap.TokenOutAddress,
		TokenOutAmount:     swap.TokenOutAmount,
		TokenOutType:       swap.TokenOutType,
		TokenOutID:         swap.TokenOutID,
		PNTxHash:           txHash,
		PNTxTimestamp:      txTimestamp,

		InitiatorSelfSalt: proofParams.selfSalt,
	}

	cancelPreimage, err := cryptography.GetPoseidonHash([]*big.Int{proofParams.destSalt})
	if err != nil {
		return fmt.Errorf("computing cancel preimage: %w", err)
	}
	passphrase, err := cryptography.GetPoseidonHash([]*big.Int{cancelPreimage, cancelPreimage})
	if err != nil {
		return fmt.Errorf("computing cancel passphrase: %w", err)
	}

	err = s.dvpClient.InitiateSwap(ctx, proofParams.destSalt, proofParams.destCiphertext, &msg, proof, swapValidity, passphrase)
	if err != nil {
		// Race condition: the other side initiated the swap first. Wait for the swap
		// to be received and persisted in our DB, then fall through to the confirmation flow.
		if errors.Is(err, dvp.ErrSwapAlreadyInitiated) {
			slog.Warn("Swap was already initiated by the other side", slog.String("sharedId", swap.SharedID))
			return fmt.Errorf("initiating swap: %w", err)
		}

		slog.Error("Swap failed to initiate.", slog.String("sharedId", swap.SharedID), slog.String("err", err.Error()))
		return fmt.Errorf("initiating swap: %w", err)
	}

	depositCommitments := make([]string, 0)
	for _, deposit := range deposits {
		depositCommitments = append(depositCommitments, deposit.Commitment.String())
	}
	err = s.depositRepository.BatchUpdateStatusForCommitments(ctx, depositCommitments, types.DvpDepositLocked)
	if err != nil {
		return fmt.Errorf("batch updating deposit status to locked: %w", err)
	}

	swap.Status = types.DvpSwapInitiated
	swap.SelfSalt = proofParams.selfSalt
	swap.DestSalt = proofParams.destSalt
	swap.CancelPreimage = cancelPreimage

	err = s.swapRepository.CreateSwap(ctx, swap)
	if err != nil {
		return fmt.Errorf("creating swap: %w", err)
	}

	slog.Info("Swap initiated successfully", slog.String("sharedId", swap.SharedID))
	return nil

}

func (s *DvpInitiator) HandleEnygmaSwapERC1155(
	ctx context.Context,
	sharedId string,
	toChainId *big.Int,
	from common.Address,
	enygmaResourceId string,
	enygmaAmount *big.Int,
	nftResourceId string,
	nftId string,
	nftAmount *big.Int,
	txHash string,
	txBlockNumber *big.Int,
	validityTime uint64,
) error {
	slog.Info("Handling Enygma -> ERC1155 swap", slog.String("sharedId", sharedId))

	swap, err := s.swapRepository.GetSwapBySharedID(ctx, sharedId)
	if err != nil {
		return fmt.Errorf("getting swap by shared ID: %w", err)
	}

	// Restore block timestamp to update Privacy Network Hub
	block, err := s.nodeClient.BlockByNumber(ctx, txBlockNumber)
	if err != nil {
		return fmt.Errorf("getting block by number: %w", err)
	}
	txTimestamp := time.Unix(int64(block.Time()), 0).UTC()

	// We are the initiator of the swap
	if swap == nil {
		slog.Info("Initiating new Enygma -> ERC1155 swap", slog.String("sharedId", sharedId))

		enygmaId := ""
		swap, err = s.createSwap(
			ctx,
			sharedId,
			toChainId,
			from,
			enygmaResourceId,
			enygmaAmount,
			enygmaId,
			types.DvpEnygma,
			nftResourceId,
			nftAmount,
			nftId,
			types.DvpERC1155,
		)
		if err != nil {
			return fmt.Errorf("creating Enygma -> ERC1155 swap: %w", err)
		}
	}

	if swap.Status == types.DvpSwapCreated {
		params, err := s.prepareSwapProofParams(ctx, swap.DestChainID)
		if err != nil {
			return fmt.Errorf("preparing Enygma -> ERC1155 proof params: %w", err)
		}

		userDeposits, err := s.depositFinder.FindEnygmaDeposits(ctx, swap.From, swap.TokenInAddress, swap.TokenInAmount)
		if err != nil {
			return fmt.Errorf("finding enygma deposits: %w", err)
		}

		jsDeposits, err := s.depositConsolidator.PrepareDepositsForJSProof(ctx, sharedId, params.selfViewPubKey, userDeposits)
		if err != nil {
			return fmt.Errorf("preparing deposits for JS proof: %w", err)
		}

		proof, err := s.proofGenerator.GenerateEnygmaToERC1155SwapProof(ctx, swap, jsDeposits, params.selfViewPubKey, params.selfSalt, params.destSalt, params.destSpendPubKey)
		if err != nil {
			return fmt.Errorf("generating Enygma -> ERC1155 proof: %w", err)
		}

		err = s.initiateSwap(ctx, jsDeposits, swap, validityTime, params, proof, txHash, txTimestamp)
		if err != nil {
			return fmt.Errorf("initiating Enygma -> ERC1155 swap: %w", err)
		}

		return nil
	}

	// The swap was already initiated by the other side, so we need to confirm and complete it.
	if swap.Status == types.DvpSwapWaitingConfirmation {
		slog.Info("Confirming Enygma -> ERC1155 swap", slog.String("sharedId", sharedId))

		swap.From = from.Hex()
		if err := s.swapRepository.UpdateSwapFrom(ctx, sharedId, swap.From); err != nil {
			return fmt.Errorf("updating swap from_address for shared_id %s: %w", sharedId, err)
		}

		// Ensure both parties agree on the same swap information
		enygmaId := ""
		if reason, agreementErr := s.swapAgreement.Verify(ctx, swap, toChainId, enygmaResourceId, enygmaAmount, enygmaId, types.DvpEnygma, nftResourceId, nftAmount, nftId, types.DvpERC1155); agreementErr != nil {
			slog.Error("Swap Enygma -> ERC1155 does not match the expected swap information", slog.String("sharedId", sharedId), slog.String("err", agreementErr.Error()), slog.String("reason", reason))

			err := s.swapAgreement.HandleSwapDisagreement(ctx, sharedId, enygmaResourceId, types.DvpEnygma, agreementErr.Error())
			if err != nil {
				return fmt.Errorf("handling swap disagreement: %w", err)
			}

			return nil
		}

		params, err := s.prepareSwapProofParams(ctx, swap.DestChainID)
		if err != nil {
			return fmt.Errorf("preparing Enygma -> ERC1155 proof params: %w", err)
		}

		userDeposits, err := s.depositFinder.FindEnygmaDeposits(ctx, swap.From, swap.TokenInAddress, swap.TokenInAmount)
		if err != nil {
			return fmt.Errorf("finding enygma deposits: %w", err)
		}

		jsDeposits, err := s.depositConsolidator.PrepareDepositsForJSProof(ctx, sharedId, params.selfViewPubKey, userDeposits)
		if err != nil {
			return fmt.Errorf("preparing deposits for JS proof: %w", err)
		}

		proof, err := s.proofGenerator.GenerateEnygmaToERC1155SwapProof(ctx, swap, jsDeposits, params.selfViewPubKey, swap.SelfSalt, swap.DestSalt, params.destSpendPubKey)
		if err != nil {
			return fmt.Errorf("generating Enygma -> ERC1155 proof: %w", err)
		}

		slog.Info("Swap Enygma -> ERC1155 sending confirmation", slog.String("sharedId", sharedId))

		msg := &types.DvpSwapMessage{
			SharedId:       sharedId,
			To:             swap.From,
			PNTxHash:       txHash,
			PNTxTimestamp:  txTimestamp,
			ChainId:        s.conf.ChainID,
			TokenInType:    swap.TokenInType,
			TokenInAddress: swap.TokenInAddress,
		}

		err = s.dvpClient.CompleteSwap(ctx, swap.DestSalt, msg, proof)
		if err != nil {
			return fmt.Errorf("confirming Enygma -> ERC1155 swap: %w", err)
		}

		slog.Info("Swap Enygma -> ERC1155 completed successfully", slog.String("sharedId", sharedId))
		return nil
	}

	return nil
}

/* DVP ERC721 Operations */
func (s *DvpInitiator) HandleERC721Creation(ctx context.Context, chainEventID string, resourceId string) error {
	slog.Info("Handling ERC721 creation", slog.String("resourceId", resourceId))

	nftHandlerAddress, err := s.plEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("getting PL resource address: %w", err)
	}

	nftHandlerInitialSupply, err := s.erc721HandlerClient.GetTotalSupply(ctx, nftHandlerAddress)
	if err != nil {
		return fmt.Errorf("getting ERC721 total supply: %w", err)
	}

	if len(nftHandlerInitialSupply) == 0 {
		slog.Info("ERC721 has no initial supply. Skipping...", slog.String("resourceId", resourceId))
		return nil
	}

	slog.Debug("ERC721 has initial supply. Syncing NFTs to the PNH contract", slog.String("resourceId", resourceId))

	// We have nft supply before the token was approved on CC
	// We need to sync the local NFTs to the PNH contract
	nftAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("getting CC resource address: %w", err)
	}

	mintDatas := make([]*dvp.DvpERC721MintData, 0)
	for _, nftId := range nftHandlerInitialSupply {
		extraData, err := s.erc721HandlerClient.GetExtraData(ctx, nftHandlerAddress, nftId)
		if err != nil {
			return fmt.Errorf("getting ERC721 extra data: %w", err)
		}

		mintData := &dvp.DvpERC721MintData{
			ChainID:      s.conf.ChainID,
			To:           s.conf.DvpOperatorAddress,
			TokenID:      nftId,
			TokenAddress: nftAddress,
			ExtraData:    extraData,
		}

		mintDatas = append(mintDatas, mintData)
	}

	mintResults, err := s.erc721Client.MintBatch(ctx, mintDatas)
	if err != nil {
		return fmt.Errorf("sending ERC721 mint batch: %w", err)
	}

	for tokenID, mintResult := range mintResults {
		if mintResult.Err != nil {
			slog.Error("Failed to mint ERC721 on CC", slog.String("tokenId", tokenID), slog.Any("error", mintResult.Err))
		} else if mintResult.Receipt.Status == 1 {
			slog.Debug("ERC721 minted on CC", slog.String("tokenId", tokenID))
		} else {
			slog.Error("Failed to mint ERC721 on CC", slog.String("tokenId", tokenID), slog.String("error", "transaction failed"))
		}
	}

	slog.Info("ERC721 initial supply synced successfully", slog.String("resourceId", resourceId))
	return nil
}

func (s *DvpInitiator) HandleERC721Mint(ctx context.Context, chainEventID string, resourceId string, nftId *big.Int) error {
	slog.Info("Handling ERC721 mint", slog.String("resourceId", resourceId))

	nftHandlerAddress, err := s.plEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("getting PL resource address: %w", err)
	}

	nftAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("getting CC resource address: %w", err)
	}

	extraData, err := s.erc721HandlerClient.GetExtraData(ctx, nftHandlerAddress, nftId)
	if err != nil {
		return fmt.Errorf("getting ERC721 extra data: %w", err)
	}

	mintData := &dvp.DvpERC721MintData{
		ChainID:      s.conf.ChainID,
		To:           s.conf.DvpOperatorAddress,
		TokenID:      nftId,
		TokenAddress: nftAddress,
		ExtraData:    extraData,
	}

	mintResults, err := s.erc721Client.MintBatch(ctx, []*dvp.DvpERC721MintData{mintData})
	if err != nil {
		return fmt.Errorf("sending ERC721 mint batch: %w", err)
	}

	mintResult := mintResults[mintData.GetID()]

	if mintResult.Err != nil {
		return fmt.Errorf("minting ERC721 tokenId %s: %w", nftId.String(), mintResult.Err)
	}

	slog.Info("ERC721 minted successfully", slog.String("resourceId", resourceId), slog.String("tokenId", nftId.String()))
	return nil
}

func (s *DvpInitiator) HandleERC721Burn(ctx context.Context, chainEventID string, resourceId string, nftId *big.Int) error {
	slog.Info("Handling ERC721 burn", slog.String("resourceId", resourceId))

	nftAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("getting CC resource address: %w", err)
	}

	err = s.erc721Client.Burn(ctx, chainEventID, nftAddress, nftId)
	if err != nil {
		return fmt.Errorf("burning ERC721: %w", err)
	}

	slog.Info("ERC721 burned successfully", slog.String("resourceId", resourceId))
	return nil
}

func (s *DvpInitiator) HandleERC721Deposit(ctx context.Context, chainEventID string, resourceId string, nftId *big.Int, from common.Address, txHash string, txBlockNumber *big.Int) error {
	slog.Info("Handling ERC721 deposit", slog.String("resourceId", resourceId))

	nftAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("getting CC resource address: %w", err)
	}

	err = s.erc721Client.Approve(ctx, chainEventID, nftAddress, s.conf.DvpContractAddress, nftId)
	if err != nil {
		return fmt.Errorf("approving ERC721: %w", err)
	}

	spendKey, err := s.keysClient.GetPaymentSpendKey(ctx, &keyspb.GetPaymentSpendKeyRequest{})
	if err != nil {
		return fmt.Errorf("getting payment spend key: %w", err)
	}
	spendPublicKey := new(big.Int).SetBytes(spendKey.GetPublicKey())

	viewKey, err := s.keysClient.GetViewPublicKey(ctx, &keyspb.GetViewPublicKeyRequest{BlockNumber: txBlockNumber.Uint64()})
	if err != nil {
		return fmt.Errorf("get rayls view key pair: %w", err)
	}

	viewPubKey, err := hex.DecodeString(viewKey.GetPublicKey())
	if err != nil {
		return fmt.Errorf("decoding rayls view public key: %w", err)
	}

	salt, _, err := cryptography.GenerateSalt(viewPubKey)
	if err != nil {
		return fmt.Errorf("generating salt: %w", err)
	}

	commitment, err := s.commitmentCalculator.CalculateNFTCommitment(spendPublicKey, salt, nftId.String(), nftAddress.Hex())
	if err != nil {
		return fmt.Errorf("calculating NFT commitment: %w", err)
	}

	err = s.depositRepository.CreateDeposit(ctx, &types.DvpDeposit{
		UserAddress:  from.Hex(),
		Salt:         salt,
		Commitment:   commitment,
		TokenAmount:  big.NewInt(1),
		TokenAddress: nftAddress.Hex(),
		TokenID:      nftId.String(),
		TokenType:    types.DvpERC721,
		Status:       types.DvpDepositPending,
	})
	if err != nil {
		return fmt.Errorf("creating ERC721 deposit: %w", err)
	}

	// Restore block timestamp to update Privacy Network Hub
	block, err := s.nodeClient.BlockByNumber(ctx, txBlockNumber)
	if err != nil {
		return fmt.Errorf("getting block by number: %w", err)
	}
	txTimestamp := time.Unix(int64(block.Time()), 0).UTC()

	burnUpdate := types.DvpBalanceUpdated{
		ErcId:             nftId,
		TokenType:         uint8(types.DVPERC721),
		ResourceId:        resourceId,
		From:              from.Hex(),
		To:                s.conf.DvpContractAddress.Hex(),
		SourceChainId:     s.conf.ChainID,
		SourceTxHash:      txHash,
		SourceTxTimestamp: txTimestamp,
		Amount:            big.NewInt(1),
		UpdateType:        types.Burn,
	}

	encryptedBurnUpdate, err := s.encryptor.EncryptDvpBalanceUpdated(ctx, burnUpdate)
	if err != nil {
		return fmt.Errorf("encrypting dvp balance update: %w", err)
	}

	err = s.dvpClient.DepositERC721(ctx, chainEventID, nftAddress, nftId, spendPublicKey, salt, encryptedBurnUpdate)
	if err != nil {
		return fmt.Errorf("depositing ERC721 to dvp contract: %w", err)
	}

	slog.Info("ERC721 deposited successfully", slog.String("resourceId", resourceId))
	return nil
}

func (s *DvpInitiator) HandleERC721Withdraw(ctx context.Context, chainEventID string, resourceId string, nftId *big.Int, from common.Address, txHash string, txBlockNumber *big.Int) error {
	slog.Info("Handling ERC721 withdrawal", slog.String("resourceId", resourceId))

	// Restore block timestamp to update Privacy Network Hub
	block, err := s.nodeClient.BlockByNumber(ctx, txBlockNumber)
	if err != nil {
		return fmt.Errorf("getting block by number: %w", err)
	}
	txTimestamp := time.Unix(int64(block.Time()), 0).UTC()

	nftAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("getting CC resource address: %w", err)
	}

	deposit, err := s.depositRepository.GetNonFungibleDeposit(ctx,
		nftId.String(),
		nftAddress.Hex(),
		from.Hex(),
		types.DvpERC721,
		types.DvpDepositUnspent,
	)
	if err != nil {
		return fmt.Errorf("getting non-fungible deposit: %w", err)
	}
	if deposit == nil {
		return fmt.Errorf("NFT deposit not found")
	}

	viewKeyPair, err := s.keysClient.GetViewPublicKey(ctx, &keyspb.GetViewPublicKeyRequest{BlockNumber: txBlockNumber.Uint64()})
	if err != nil {
		return fmt.Errorf("get rayls view key pair: %w", err)
	}

	viewPubKey, err := hex.DecodeString(viewKeyPair.GetPublicKey())
	if err != nil {
		return fmt.Errorf("decoding rayls view public key: %w", err)
	}

	destSalt, _, err := cryptography.GenerateSalt(viewPubKey)
	if err != nil {
		return fmt.Errorf("generating destination salt: %w", err)
	}

	// Generate proof BEFORE locking — proof generation is pure computation.
	proof, err := s.proofGenerator.GenerateERC721WithdrawProof(ctx, viewPubKey, destSalt, s.conf.DvpOperatorAddress.Big(), deposit)
	if err != nil {
		return fmt.Errorf("generating ERC721 withdraw proof: %w", err)
	}

	mintUpdate := types.DvpBalanceUpdated{
		ErcId:             nftId,
		TokenType:         uint8(types.DVPERC721),
		ResourceId:        resourceId,
		From:              from.Hex(),
		To:                s.conf.DvpContractAddress.Hex(),
		SourceChainId:     s.conf.ChainID,
		SourceTxHash:      txHash,
		SourceTxTimestamp: txTimestamp,
		Amount:            big.NewInt(1),
		UpdateType:        types.Mint,
	}

	encryptedMintUpdate, err := s.encryptor.EncryptDvpBalanceUpdated(ctx, mintUpdate)
	if err != nil {
		return fmt.Errorf("encrypting dvp balance update: %w", err)
	}

	// Lock deposit and set nullifier right before the external call.
	// If the relayer crashes before this point, the deposit stays Unspent (no recovery needed).
	// If it crashes after, PersistAndBroadcast's recovery mechanism handles re-broadcast.
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if txErr := s.depositRepository.UpdateDepositStatus(txCtx, deposit.Commitment, types.DvpDepositLocked); txErr != nil {
			return txErr
		}
		return s.depositRepository.UpsertDepositNullifier(txCtx, deposit.Commitment, proof.Nullifiers[0])
	})
	if err != nil {
		return err
	}

	// Sign without broadcasting
	err = s.dvpClient.WithdrawERC721(ctx, chainEventID, nftAddress, nftId, s.conf.DvpOperatorAddress, destSalt, proof, encryptedMintUpdate)
	if err != nil {
		return err
	}

	slog.Debug("ERC721 withdrawal completed. Syncing extra data to the PNH contract...", slog.String("resourceId", resourceId))

	nftHandlerAddress, err := s.plEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("getting PL resource address: %w", err)
	}

	extraData, err := s.erc721HandlerClient.GetExtraData(ctx, nftHandlerAddress, nftId)
	if err != nil {
		return fmt.Errorf("getting ERC721 extra data: %w", err)
	}

	err = s.erc721Client.UpdateExtraData(ctx, chainEventID, nftAddress, nftId, s.conf.ChainID, extraData, s.conf.DvpOperatorAddress)
	if err != nil {
		return fmt.Errorf("updating ERC721 extra data: %w", err)
	}

	slog.Debug("Extra data synced to successfully. Unlocking ERC721 on PL...", slog.String("resourceId", resourceId), slog.String("nftId", nftId.String()))

	err = s.erc721HandlerClient.Unlock(ctx, chainEventID, nftHandlerAddress, nftId)
	if err != nil {
		return fmt.Errorf("unlocking ERC721: %w", err)
	}

	slog.Info("ERC721 unlocked successfully", slog.String("resourceId", resourceId))

	return nil
}

func (s *DvpInitiator) HandleERC721SwapEnygma(
	ctx context.Context,
	sharedId string,
	toChainId *big.Int,
	from common.Address,
	nftResourceId string,
	nftId string,
	enygmaResourceId string,
	enygmaAmount *big.Int,
	txHash string,
	txBlockNumber *big.Int,
	validityTime uint64,
) error {
	slog.Info("Handling ERC721 -> Enygma swap", slog.String("sharedId", sharedId))

	swap, err := s.swapRepository.GetSwapBySharedID(ctx, sharedId)
	if err != nil {
		return fmt.Errorf("getting swap by shared ID: %w", err)
	}

	// Restore block timestamp to update Privacy Network Hub
	block, err := s.nodeClient.BlockByNumber(ctx, txBlockNumber)
	if err != nil {
		return fmt.Errorf("getting block by number: %w", err)
	}
	txTimestamp := time.Unix(int64(block.Time()), 0).UTC()

	// We are the initiator of the swap
	if swap == nil {
		slog.Debug("Initiating new ERC721 -> Enygma swap", slog.String("sharedId", sharedId))

		nftAmount := big.NewInt(1)
		enygmaId := ""
		swap, err = s.createSwap(
			ctx,
			sharedId,
			toChainId,
			from,
			nftResourceId,
			nftAmount,
			nftId,
			types.DvpERC721,
			enygmaResourceId,
			enygmaAmount,
			enygmaId,
			types.DvpEnygma,
		)

		if err != nil {
			return fmt.Errorf("creating ERC721 -> Enygma swap: %w", err)
		}
	}

	if swap.Status == types.DvpSwapCreated {
		slog.Info("Initiating ERC721 -> Enygma swap", slog.String("sharedId", sharedId))

		params, err := s.prepareSwapProofParams(ctx, swap.DestChainID)
		if err != nil {
			return fmt.Errorf("preparing ERC721 -> Enygma proof params: %w", err)
		}

		deposit, err := s.depositFinder.FindERC721Deposit(ctx, swap.From, swap.TokenInAddress, swap.TokenInID)
		if err != nil {
			return fmt.Errorf("finding ERC721 deposit: %w", err)
		}

		proof, err := s.proofGenerator.GenerateERC721ToEnygmaSwapProof(ctx, swap, deposit, params.selfViewPubKey, params.selfSalt, params.destSalt, params.destSpendPubKey)
		if err != nil {
			return fmt.Errorf("generating ERC721 -> Enygma proof: %w", err)
		}

		err = s.initiateSwap(ctx, []*types.DvpDeposit{deposit}, swap, validityTime, params, proof, txHash, txTimestamp)
		if err != nil {
			return fmt.Errorf("initiating ERC721 -> Enygma swap: %w", err)
		}

		return nil
	}
	// The swap was initiated by the other party, so we need to confirm the swap
	if swap.Status == types.DvpSwapWaitingConfirmation {
		slog.Info("Confirming ERC721 -> Enygma swap", slog.String("sharedId", sharedId))

		swap.From = from.Hex()
		if err := s.swapRepository.UpdateSwapFrom(ctx, sharedId, swap.From); err != nil {
			return fmt.Errorf("updating swap from_address for shared_id %s: %w", sharedId, err)
		}

		// Ensure both parties agree on the same swap information
		nftAmount := big.NewInt(1)
		enygmaId := ""
		if reason, agreementErr := s.swapAgreement.Verify(ctx, swap, toChainId, nftResourceId, nftAmount, nftId, types.DvpERC721, enygmaResourceId, enygmaAmount, enygmaId, types.DvpEnygma); agreementErr != nil {
			slog.Error("Swap ERC721 -> Enygma does not match the expected swap information", slog.String("sharedId", sharedId), slog.String("err", agreementErr.Error()), slog.String("reason", reason))

			err := s.swapAgreement.HandleSwapDisagreement(ctx, sharedId, nftResourceId, types.DvpERC721, agreementErr.Error())
			if err != nil {
				return fmt.Errorf("handling swap disagreement: %w", err)
			}

			return nil
		}

		params, err := s.prepareSwapProofParams(ctx, swap.DestChainID)
		if err != nil {
			return fmt.Errorf("preparing ERC721 -> Enygma proof params: %w", err)
		}

		deposit, err := s.depositFinder.FindERC721Deposit(ctx, swap.From, swap.TokenInAddress, swap.TokenInID)
		if err != nil {
			return fmt.Errorf("finding ERC721 deposit: %w", err)
		}

		proof, err := s.proofGenerator.GenerateERC721ToEnygmaSwapProof(ctx, swap, deposit, params.selfViewPubKey, swap.SelfSalt, swap.DestSalt, params.destSpendPubKey)
		if err != nil {
			return fmt.Errorf("generating ERC721 -> Enygma proof: %w", err)
		}

		slog.Debug("Swap ERC721 -> Enygma sending confirmation", slog.String("sharedId", sharedId))

		msg := &types.DvpSwapMessage{
			SharedId:       sharedId,
			To:             swap.From,
			PNTxHash:       txHash,
			PNTxTimestamp:  txTimestamp,
			ChainId:        s.conf.ChainID,
			TokenInType:    swap.TokenInType,
			TokenInAddress: swap.TokenInAddress,
		}

		err = s.dvpClient.CompleteSwap(ctx, swap.DestSalt, msg, proof)
		if err != nil {
			return fmt.Errorf("confirming Enygma -> ERC721 swap: %w", err)
		}

		slog.Info("Swap ERC721 -> Enygma completed successfully", slog.String("sharedId", sharedId))
	}

	return nil
}

type swapProofParams struct {
	selfSalt        *big.Int
	selfViewPubKey  []byte
	destSalt        *big.Int
	destCiphertext  []byte
	destSpendPubKey *big.Int
}

func (s *DvpInitiator) prepareSwapProofParams(ctx context.Context, destChainId *big.Int) (*swapProofParams, error) {
	block, err := s.hubClient.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("get hub node latest block: %w", err)
	}

	destSpendPubKey, err := s.psClient.GetPaymentSpendPublicKey(ctx, destChainId)
	if err != nil {
		return nil, fmt.Errorf("getting dest spend public key: %w", err)
	}

	destViewKey, err := s.psClient.GetChainViewData(ctx, destChainId, block.Number())
	if err != nil {
		return nil, fmt.Errorf("getting destination view key: %w", err)
	}

	destViewPubKey, err := hex.DecodeString(destViewKey.RaylsViewPublicKey)
	if err != nil {
		return nil, fmt.Errorf("decoding rayls view public key: %w", err)
	}

	destSalt, destCiphertext, err := cryptography.GenerateSalt(destViewPubKey)
	if err != nil {
		return nil, fmt.Errorf("generating destination salt: %w", err)
	}

	selfViewKeyPair, err := s.keysClient.GetViewPublicKey(ctx, &keyspb.GetViewPublicKeyRequest{BlockNumber: block.Number().Uint64()})
	if err != nil {
		return nil, fmt.Errorf("get rayls view key pair: %w", err)
	}

	selfViewPubKey, err := hex.DecodeString(selfViewKeyPair.GetPublicKey())
	if err != nil {
		return nil, fmt.Errorf("decoding rayls view public key: %w", err)
	}

	selfSalt, _, err := cryptography.GenerateSalt(selfViewPubKey)
	if err != nil {
		return nil, fmt.Errorf("generating salt: %w", err)
	}

	return &swapProofParams{
		selfSalt:        selfSalt,
		selfViewPubKey:  selfViewPubKey,
		destSalt:        destSalt,
		destCiphertext:  destCiphertext,
		destSpendPubKey: destSpendPubKey,
	}, nil
}

/* DVP ERC1155 Operations */
func (s *DvpInitiator) HandleERC1155Creation(ctx context.Context, chainEventID string, resourceId string) error {
	slog.Info("Handling ERC1155 creation", slog.String("resourceId", resourceId))

	tokenHandlerAddress, err := s.plEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("getting PL resource address: %w", err)
	}

	totalTokens, err := s.erc1155HandlerClient.GetAllTokenIdsWithSupply(ctx, tokenHandlerAddress)
	if err != nil {
		return fmt.Errorf("getting all ERC1155 token IDs with supply: %w", err)
	}

	if len(totalTokens) == 0 {
		slog.Info("ERC1155 has no initial supply. Skipping...", slog.String("resourceId", resourceId))
		return nil
	}

	slog.Debug("ERC1155 has initial supply. Syncing tokens to the PNH contract", slog.String("resourceId", resourceId))

	// We have nft supply before the token was approved on CC
	// We need to sync the local NFTs to the PNH contract
	tokenAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("getting CC resource address: %w", err)
	}

	mintDatas := make([]*dvp.DvpERC1155MintData, 0)
	for _, token := range totalTokens {
		tokenExtraData, err := s.erc1155HandlerClient.GetTokenExtraData(ctx, tokenHandlerAddress, token.Id)

		if err != nil {
			return fmt.Errorf("getting ERC1155 token extra data: %w", err)
		}

		mintData := &dvp.DvpERC1155MintData{
			ChainID:      s.conf.ChainID,
			To:           s.conf.DvpOperatorAddress,
			TokenID:      token.Id,
			TokenAmount:  token.Amount,
			TokenAddress: tokenAddress,
			Data:         []byte{},
			ExtraData:    tokenExtraData,
		}

		mintDatas = append(mintDatas, mintData)
	}

	mintResults, err := s.erc1155Client.MintBatch(ctx, mintDatas)
	if err != nil {
		return fmt.Errorf("sending ERC1155 mint batch: %w", err)
	}

	for tokenID, mintResult := range mintResults {
		if mintResult.Err != nil {
			slog.Error("Failed to mint ERC1155 on CC", slog.String("tokenId", tokenID), slog.Any("error", mintResult.Err))
		} else if mintResult.Receipt.Status == 1 {
			slog.Debug("ERC1155 minted on CC", slog.String("tokenId", tokenID))
		} else {
			slog.Error("Failed to mint ERC1155 on CC", slog.String("tokenId", tokenID), slog.String("error", "transaction failed"))
		}
	}

	slog.Info("ERC1155 initial supply synced successfully", slog.String("resourceId", resourceId))
	return nil
}

func (s *DvpInitiator) HandleERC1155Mint(ctx context.Context, chainEventID string, resourceId string, tokenId *big.Int, tokenAmount *big.Int, tokenData []byte) error {
	slog.Info("Handling ERC1155 mint", slog.String("resourceId", resourceId))

	tokenHandlerAddress, err := s.plEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("getting PL resource address: %w", err)
	}

	tokenAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("getting CC resource address: %w", err)
	}

	tokenExtraData, err := s.erc1155HandlerClient.GetTokenExtraData(ctx, tokenHandlerAddress, tokenId)
	if err != nil {
		return fmt.Errorf("getting ERC1155 token extra data: %w", err)
	}

	mintData := &dvp.DvpERC1155MintData{
		ChainID:      s.conf.ChainID,
		To:           s.conf.DvpOperatorAddress,
		TokenID:      tokenId,
		TokenAddress: tokenAddress,
		TokenAmount:  tokenAmount,
		Data:         tokenData,
		ExtraData:    tokenExtraData,
	}

	mintResults, err := s.erc1155Client.MintBatch(ctx, []*dvp.DvpERC1155MintData{mintData})
	if err != nil {
		return fmt.Errorf("sending ERC1155 mint batch: %w", err)
	}

	mintResult := mintResults[mintData.GetID()]

	if mintResult.Err != nil {
		slog.Error("Failed to mint ERC1155 on CC", slog.String("tokenId", tokenId.String()), slog.Any("error", mintResult.Err))
	} else if mintResult.Receipt.Status == 1 {
		slog.Debug("ERC1155 minted on CC", slog.String("tokenId", tokenId.String()))
	} else {
		slog.Error("Failed to mint ERC1155 on CC", slog.String("tokenId", tokenId.String()), slog.String("error", "transaction failed"))
	}

	slog.Info("ERC1155 minted successfully", slog.String("resourceId", resourceId))
	return nil
}

func (s *DvpInitiator) HandleERC1155Burn(ctx context.Context, chainEventID string, resourceId string, tokenId *big.Int, tokenAmount *big.Int) error {
	slog.Info("Handling ERC1155 burn", slog.String("resourceId", resourceId))

	tokenAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("getting CC resource address: %w", err)
	}

	err = s.erc1155Client.Burn(ctx, chainEventID, tokenAddress, s.conf.DvpOperatorAddress, tokenId, tokenAmount)
	if err != nil {
		return fmt.Errorf("burning ERC1155: %w", err)
	}

	slog.Info("ERC1155 burned successfully", slog.String("resourceId", resourceId))
	return nil
}

func (s *DvpInitiator) HandleERC1155Deposit(ctx context.Context, chainEventID string, resourceId string, tokenId *big.Int, tokenAmount *big.Int, tokenData []byte, from common.Address, txHash string, txBlockNumber *big.Int) error {
	slog.Info("Handling ERC1155 deposit", slog.String("resourceId", resourceId))

	tokenAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("getting CC resource address: %w", err)
	}

	err = s.erc1155Client.Approve(ctx, chainEventID, tokenAddress, s.conf.DvpContractAddress)
	if err != nil {
		return fmt.Errorf("approving ERC1155: %w", err)
	}

	spendKey, err := s.keysClient.GetPaymentSpendKey(ctx, &keyspb.GetPaymentSpendKeyRequest{})
	if err != nil {
		return fmt.Errorf("getting payment spend key: %w", err)
	}
	spendPublicKey := new(big.Int).SetBytes(spendKey.GetPublicKey())

	viewKey, err := s.keysClient.GetViewPublicKey(ctx, &keyspb.GetViewPublicKeyRequest{BlockNumber: txBlockNumber.Uint64()})
	if err != nil {
		return fmt.Errorf("get rayls view key pair: %w", err)
	}

	viewPubKey, err := hex.DecodeString(viewKey.GetPublicKey())
	if err != nil {
		return fmt.Errorf("decoding rayls view public key: %w", err)
	}

	salt, _, err := cryptography.GenerateSalt(viewPubKey)
	if err != nil {
		return fmt.Errorf("generating salt: %w", err)
	}

	commitment, err := s.commitmentCalculator.CalculateERC1155Commitment(spendPublicKey, salt, tokenAddress.Hex(), tokenId.String(), tokenAmount)
	if err != nil {
		return fmt.Errorf("calculating ERC1155 commitment: %w", err)
	}

	err = s.depositRepository.CreateDeposit(ctx, &types.DvpDeposit{
		UserAddress:  from.Hex(),
		Salt:         salt,
		TokenAmount:  tokenAmount,
		TokenAddress: tokenAddress.Hex(),
		TokenID:      tokenId.String(),
		TokenType:    types.DvpERC1155,
		Commitment:   commitment,
		Status:       types.DvpDepositPending,
	})
	if err != nil {
		return fmt.Errorf("creating ERC1155 deposit: %w", err)
	}

	// Restore block timestamp to update Privacy Network Hub
	block, err := s.nodeClient.BlockByNumber(ctx, txBlockNumber)
	if err != nil {
		return fmt.Errorf("getting block by number: %w", err)
	}
	txTimestamp := time.Unix(int64(block.Time()), 0).UTC()

	burnUpdate := types.DvpBalanceUpdated{
		ErcId:             tokenId,
		TokenType:         uint8(types.DVPERC1155),
		ResourceId:        resourceId,
		From:              from.Hex(),
		To:                s.conf.DvpContractAddress.Hex(),
		SourceChainId:     s.conf.ChainID,
		SourceTxHash:      txHash,
		SourceTxTimestamp: txTimestamp,
		Amount:            tokenAmount,
		UpdateType:        types.Burn,
	}

	encryptedBurnUpdate, err := s.encryptor.EncryptDvpBalanceUpdated(ctx, burnUpdate)
	if err != nil {
		return fmt.Errorf("encrypting dvp balance update: %w", err)
	}

	err = s.dvpClient.DepositERC1155(ctx, chainEventID, tokenAddress, tokenId, tokenAmount, tokenData, spendPublicKey, salt, encryptedBurnUpdate)
	if err != nil {
		return fmt.Errorf("depositing ERC1155 to dvp contract: %w", err)
	}

	slog.Info("ERC1155 deposit executed successfully", slog.String("resourceId", resourceId))

	return nil
}

func (s *DvpInitiator) HandleERC1155Withdraw(ctx context.Context, chainEventID string, resourceId string, tokenId *big.Int, tokenAmount *big.Int, from common.Address, txHash string, txBlockNumber *big.Int) error {
	slog.Info("Handling ERC1155 withdrawal", slog.String("resourceId", resourceId))

	tokenAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("getting CC resource address: %w", err)
	}

	viewKeyPair, err := s.keysClient.GetViewPublicKey(ctx, &keyspb.GetViewPublicKeyRequest{BlockNumber: txBlockNumber.Uint64()})
	if err != nil {
		return fmt.Errorf("get rayls view key pair: %w", err)
	}

	viewPubKey, err := hex.DecodeString(viewKeyPair.GetPublicKey())
	if err != nil {
		return fmt.Errorf("decoding rayls view public key: %w", err)
	}

	userDeposits, err := s.depositFinder.FindERC1155DepositsForJSProof(ctx, from.Hex(), tokenAddress.Hex(), tokenId.String(), tokenAmount)
	if err != nil {
		return fmt.Errorf("finding ERC1155 deposits for JS proof: %w", err)
	}

	jsDeposits, err := s.depositConsolidator.PrepareDepositsForJSProof(ctx, chainEventID, viewPubKey, userDeposits)
	if err != nil {
		return fmt.Errorf("preparing deposits for JS proof: %w", err)
	}

	totalDepositAmount := CalculateTotalAmountOfDeposits(jsDeposits)

	if tokenAmount.Cmp(totalDepositAmount) < 0 {
		jsDeposits, err = s.depositConsolidator.ConsolidateERC1155Deposits(ctx, chainEventID, viewPubKey, jsDeposits, tokenAmount)
		if err != nil {
			return fmt.Errorf("consolidating ERC1155 deposits: %w", err)
		}
		// Get first deposit that match the withdrawal amount. In case we have 2 deposits with same amount, we will use the first one for the withdraw
		for _, deposit := range jsDeposits {
			if deposit.TokenAmount.Cmp(tokenAmount) == 0 {
				jsDeposits = []*types.DvpDeposit{deposit}
				break
			}
		}
	}

	depositCommitments := make([]string, 0)
	for _, deposit := range jsDeposits {
		depositCommitments = append(depositCommitments, deposit.Commitment.String())
	}
	err = s.depositRepository.BatchUpdateStatusForCommitments(ctx, depositCommitments, types.DvpDepositLocked)
	if err != nil {
		return fmt.Errorf("batch updating deposit status to locked: %w", err)
	}

	destSalt, _, err := cryptography.GenerateSalt(viewPubKey)
	if err != nil {
		return fmt.Errorf("generating destination salt: %w", err)
	}

	jsProof, err := s.proofGenerator.GenerateERC1155WithdrawProof(ctx, viewPubKey, destSalt, s.conf.DvpOperatorAddress.Big(), from.Hex(), tokenAddress.Hex(), tokenId.String(), tokenAmount, jsDeposits)
	if err != nil {
		return fmt.Errorf("generating ERC1155 withdraw proof: %w", err)
	}

	commitmentNullifierMap := make(map[string]string)
	for i, deposit := range jsDeposits {
		commitmentNullifierMap[deposit.Commitment.String()] = jsProof.Nullifiers[i].String()
	}
	err = s.depositRepository.BatchUpsertNullifiers(ctx, commitmentNullifierMap)
	if err != nil {
		return fmt.Errorf("batch upserting nullifiers: %w", err)
	}

	// Restore block timestamp to update Privacy Network Hub
	block, err := s.nodeClient.BlockByNumber(ctx, txBlockNumber)
	if err != nil {
		return fmt.Errorf("getting block by number: %w", err)
	}
	txTimestamp := time.Unix(int64(block.Time()), 0).UTC()

	mintUpdate := types.DvpBalanceUpdated{
		ErcId:             tokenId,
		TokenType:         uint8(types.DVPERC1155),
		ResourceId:        resourceId,
		From:              from.Hex(),
		To:                s.conf.DvpContractAddress.Hex(),
		SourceChainId:     s.conf.ChainID,
		SourceTxHash:      txHash,
		SourceTxTimestamp: txTimestamp,
		Amount:            tokenAmount,
		UpdateType:        types.Mint,
	}

	encryptedMintUpdate, err := s.encryptor.EncryptDvpBalanceUpdated(ctx, mintUpdate)
	if err != nil {
		return fmt.Errorf("encrypting dvp balance update: %w", err)
	}

	err = s.dvpClient.WithdrawERC1155(ctx, chainEventID, tokenAddress, tokenId, tokenAmount, s.conf.DvpOperatorAddress, destSalt, jsProof, encryptedMintUpdate)
	if err != nil {
		return err
	}

	slog.Debug("ERC1155 withdrawal completed. Syncing extra data to the PNH contract...", slog.String("resourceId", resourceId))

	tokenHandlerAddress, err := s.plEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("getting PL resource address: %w", err)
	}

	extraData, err := s.erc1155HandlerClient.GetTokenExtraData(ctx, tokenHandlerAddress, tokenId)
	if err != nil {
		return fmt.Errorf("getting ERC1155 token extra data: %w", err)
	}

	err = s.erc1155Client.UpdateExtraData(ctx, chainEventID, tokenAddress, tokenId, tokenAmount, s.conf.ChainID, extraData, s.conf.DvpOperatorAddress)
	if err != nil {
		return fmt.Errorf("updating ERC1155 extra data: %w", err)
	}

	slog.Debug("ERC1155 extra data synced successfully. Unlocking ERC1155 on PL...", slog.String("resourceId", resourceId))

	err = s.erc1155HandlerClient.Unlock(ctx, chainEventID, tokenHandlerAddress, tokenId, tokenAmount, from)
	if err != nil {
		return fmt.Errorf("unlocking ERC1155: %w", err)
	}

	slog.Info("ERC1155 unlocked successfully", slog.String("resourceId", resourceId))

	return nil
}

func (s *DvpInitiator) HandleERC1155SwapEnygma(
	ctx context.Context,
	sharedId string,
	toChainId *big.Int,
	from common.Address,
	erc1155ResourceId string,
	erc1155Amount *big.Int,
	erc1155Id string,
	enygmaResourceId string,
	enygmaAmount *big.Int,
	txHash string,
	txBlockNumber *big.Int,
	validityTime uint64,
) error {
	slog.Info("Handling ERC1155 -> Enygma swap", slog.String("sharedId", sharedId))

	swap, err := s.swapRepository.GetSwapBySharedID(ctx, sharedId)
	if err != nil {
		return fmt.Errorf("getting swap by shared ID: %w", err)
	}

	// Restore block timestamp to update Privacy Network Hub
	block, err := s.nodeClient.BlockByNumber(ctx, txBlockNumber)
	if err != nil {
		return fmt.Errorf("getting block by number: %w", err)
	}
	txTimestamp := time.Unix(int64(block.Time()), 0).UTC()

	// We are the initiator of the swap
	if swap == nil {
		slog.Debug("Initiating new ERC1155 -> Enygma swap", slog.String("sharedId", sharedId))

		enygmaId := ""
		swap, err = s.createSwap(
			ctx,
			sharedId,
			toChainId,
			from,
			erc1155ResourceId,
			erc1155Amount,
			erc1155Id,
			types.DvpERC1155,
			enygmaResourceId,
			enygmaAmount,
			enygmaId,
			types.DvpEnygma,
		)
		if err != nil {
			return fmt.Errorf("creating ERC1155 -> Enygma swap: %w", err)
		}
	}

	if swap.Status == types.DvpSwapCreated {
		slog.Info("Initiating ERC1155 -> Enygma swap", slog.String("sharedId", sharedId))

		params, err := s.prepareSwapProofParams(ctx, swap.DestChainID)
		if err != nil {
			return fmt.Errorf("preparing ERC1155 -> Enygma proof params: %w", err)
		}

		userDeposits, err := s.depositFinder.FindERC1155DepositsForJSProof(ctx, swap.From, swap.TokenInAddress, swap.TokenInID, swap.TokenInAmount)
		if err != nil {
			return fmt.Errorf("finding ERC1155 deposits: %w", err)
		}

		jsDeposits, err := s.depositConsolidator.PrepareDepositsForJSProof(ctx, sharedId, params.selfViewPubKey, userDeposits)
		if err != nil {
			return fmt.Errorf("preparing deposits for JS proof: %w", err)
		}

		proof, err := s.proofGenerator.GenerateERC1155ToEnygmaSwapProof(ctx, swap, jsDeposits, params.selfViewPubKey, params.selfSalt, params.destSalt, params.destSpendPubKey)
		if err != nil {
			return fmt.Errorf("generating ERC1155 -> Enygma proof: %w", err)
		}

		err = s.initiateSwap(ctx, jsDeposits, swap, validityTime, params, proof, txHash, txTimestamp)
		if err != nil {
			return fmt.Errorf("initiating ERC1155 -> Enygma swap: %w", err)
		}

		return nil
	}

	// The swap was initiated by the other party, so we need to confirm the swap
	if swap.Status == types.DvpSwapWaitingConfirmation {
		slog.Info("Confirming ERC1155 -> Enygma swap", slog.String("sharedId", sharedId))

		swap.From = from.Hex()
		if err := s.swapRepository.UpdateSwapFrom(ctx, sharedId, swap.From); err != nil {
			return fmt.Errorf("updating swap from_address for shared_id %s: %w", sharedId, err)
		}

		// Ensure both parties agree on the same swap information
		enygmaId := ""
		if reason, agreementErr := s.swapAgreement.Verify(ctx, swap, toChainId, erc1155ResourceId, erc1155Amount, erc1155Id, types.DvpERC1155, enygmaResourceId, enygmaAmount, enygmaId, types.DvpEnygma); agreementErr != nil {
			slog.Error("Swap ERC1155 -> Enygma does not match the expected swap information", slog.String("sharedId", sharedId), slog.String("err", agreementErr.Error()), slog.String("reason", reason))

			err := s.swapAgreement.HandleSwapDisagreement(ctx, sharedId, erc1155ResourceId, types.DvpERC1155, agreementErr.Error())
			if err != nil {
				return fmt.Errorf("handling swap disagreement: %w", err)
			}

			return nil
		}

		params, err := s.prepareSwapProofParams(ctx, swap.DestChainID)
		if err != nil {
			return fmt.Errorf("preparing ERC1155 -> Enygma proof params: %w", err)
		}

		userDeposits, err := s.depositFinder.FindERC1155DepositsForJSProof(ctx, swap.From, swap.TokenInAddress, swap.TokenInID, swap.TokenInAmount)
		if err != nil {
			return fmt.Errorf("finding ERC1155 deposits: %w", err)
		}

		jsDeposits, err := s.depositConsolidator.PrepareDepositsForJSProof(ctx, sharedId, params.selfViewPubKey, userDeposits)
		if err != nil {
			return fmt.Errorf("preparing deposits for JS proof: %w", err)
		}

		proof, err := s.proofGenerator.GenerateERC1155ToEnygmaSwapProof(ctx, swap, jsDeposits, params.selfViewPubKey, swap.SelfSalt, swap.DestSalt, params.destSpendPubKey)
		if err != nil {
			return fmt.Errorf("generating ERC1155 -> Enygma proof: %w", err)
		}

		slog.Debug("Swap ERC1155 -> Enygma sending confirmation", slog.String("sharedId", sharedId))

		msg := &types.DvpSwapMessage{
			SharedId:       sharedId,
			To:             swap.From,
			PNTxHash:       txHash,
			PNTxTimestamp:  txTimestamp,
			ChainId:        s.conf.ChainID,
			TokenInType:    swap.TokenInType,
			TokenInAddress: swap.TokenInAddress,
		}

		err = s.dvpClient.CompleteSwap(ctx, swap.DestSalt, msg, proof)
		if err != nil {
			return fmt.Errorf("confirming ERC1155 -> Enygma swap: %w", err)
		}

		slog.Info("Swap ERC1155 -> Enygma completed successfully", slog.String("sharedId", sharedId))
	}

	return nil
}

func (s *DvpInitiator) HandleSwapCancellation(
	ctx context.Context,
	sharedId string,
	toChainId *big.Int,
	tokenInResourceId string,
	tokenInAmount *big.Int,
	tokenInID string,
	tokenInType types.DvpTokenType,
	tokenOutResourceId string,
	tokenOutAmount *big.Int,
	tokenOutID string,
	tokenOutType types.DvpTokenType,
) error {
	slog.Info("Handling swap cancellation", slog.String("sharedId", sharedId))

	swap, err := s.swapRepository.GetSwapBySharedID(ctx, sharedId)
	if err != nil {
		return fmt.Errorf("getting swap by shared ID: %w", err)
	}

	if swap == nil {
		return fmt.Errorf("swap not found")
	}

	// Enygma tokens do not have an ID, so we need to set it to an empty string.
	if tokenInType == types.DvpEnygma {
		tokenInID = ""
	} else if tokenOutType == types.DvpEnygma {
		tokenOutID = ""
	}

	reason, err := s.swapAgreement.Verify(ctx, swap, toChainId, tokenInResourceId, tokenInAmount, tokenInID, tokenInType, tokenOutResourceId, tokenOutAmount, tokenOutID, tokenOutType)
	if err != nil {
		slog.Error("Swap cancellation does not match the expected swap information", slog.String("sharedId", sharedId), slog.String("err", err.Error()), slog.String("reason", reason))

		return fmt.Errorf("verifying swap agreement for cancellation: %w", err)
	}

	if swap.CancelPreimage == nil {
		return fmt.Errorf("swap %s has no cancel preimage, cannot cancel", sharedId)
	}

	err = s.dvpClient.CancelSwap(ctx, sharedId, swap.CancelPreimage)
	if err != nil {
		return fmt.Errorf("cancelling swap: %w", err)
	}

	slog.Info("Swap cancellation initiated successfully", slog.String("sharedId", sharedId))

	return nil
}

func (s *DvpInitiator) createSwap(
	ctx context.Context,
	sharedId string,
	toChainId *big.Int,
	from common.Address,
	tokenInResourceId string,
	tokenInAmount *big.Int,
	tokenInID string,
	tokenInType types.DvpTokenType,
	tokenOutResourceId string,
	tokenOutAmount *big.Int,
	tokenOutID string,
	tokenOutType types.DvpTokenType,
) (*types.DvpSwap, error) {
	tokenInAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, tokenInResourceId)
	if err != nil {
		return nil, fmt.Errorf("getting token-in resource address: %w", err)
	}

	tokenOutAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, tokenOutResourceId)
	if err != nil {
		return nil, fmt.Errorf("getting token-out resource address: %w", err)
	}

	newSwap := &types.DvpSwap{
		SharedID:           sharedId,
		From:               from.Hex(),
		SourceChainID:      s.conf.ChainID,
		DestChainID:        toChainId,
		TokenInAmount:      tokenInAmount,
		TokenInAddress:     tokenInAddress.Hex(),
		TokenInResourceID:  tokenInResourceId,
		TokenInType:        tokenInType,
		TokenInID:          tokenInID,
		TokenOutAmount:     tokenOutAmount,
		TokenOutAddress:    tokenOutAddress.Hex(),
		TokenOutID:         tokenOutID,
		TokenOutResourceID: tokenOutResourceId,
		TokenOutType:       tokenOutType,
		Status:             types.DvpSwapCreated,
		DestSalt:           nil,
		SelfSalt:           nil,
	}

	return newSwap, nil
}

func CalculateTotalAmountOfDeposits(deposits []*types.DvpDeposit) *big.Int {
	totalAmount := big.NewInt(0)
	for _, deposit := range deposits {
		totalAmount = totalAmount.Add(totalAmount, deposit.TokenAmount)
	}
	return totalAmount
}
