package service

//go:generate moq --skip-ensure --pkg service_test -out consolidation_mock_test.go . consolidationDepositRepository consolidationTxManager consolidationCommitmentCalculator consolidationSpendKeyProvider consolidationProofService consolidationDvpClient ConsolidationEnygmaClient ConsolidationDvpIntegrationClient consolidationDepositWaiter

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cryptography"
	keyspb "github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/keys"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"google.golang.org/grpc"
)

type ConsolidationConfig struct {
	ChainID               *big.Int
	MaxNumberOfJSDeposits int
}

type consolidationDepositRepository interface {
	BatchUpdateStatusForCommitments(ctx context.Context, commitments []string, status types.DvpDepositStatus) error
	BatchUpsertNullifiers(ctx context.Context, commitmentNullifierMap map[string]string) error
	GetDepositByCommitment(ctx context.Context, commitment *big.Int) (*types.DvpDeposit, error)
	CreateDeposit(ctx context.Context, deposit *types.DvpDeposit) error
}

type consolidationTxManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type consolidationCommitmentCalculator interface {
	CalculateERC1155Commitment(spendPK, salt *big.Int, tokenAddress string, tokenID string, tokenAmount *big.Int) (*big.Int, error)
	CalculatePaymentCommitment(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error)
}

type consolidationSpendKeyProvider interface {
	GetPaymentSpendKey(ctx context.Context, in *keyspb.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error)
}

type consolidationProofService interface {
	GenerateERC1155JSProof(
		ctx context.Context,
		sourceViewPublicKey []byte,
		paymentCommitment *big.Int,
		destinationPaymentPublicKey *big.Int,
		destinationSalt *big.Int,
		paymentAmount *big.Int,
		userAddress string,
		tokenAddress string,
		tokenID string,
		deposits []*types.DvpDeposit,
	) (*dvp.ProofReceipt, error)
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
}

type consolidationDvpClient interface {
	MixFundsERC1155(ctx context.Context, id string, contractAddress common.Address, proofReceipt *dvp.ProofReceipt) error
}

// EnygmaClient interface defines only the methods needed by consolidation service
type ConsolidationEnygmaClient interface {
	GetDvpIntegrationContractAddress(ctx context.Context, tokenAddress common.Address) (common.Address, error)
}

// DvpIntegrationClient interface defines only the methods needed by consolidation service
type ConsolidationDvpIntegrationClient interface {
	ConsolidateFunds(ctx context.Context, id string, enygmaIntegrationAddress common.Address, proofReceipt *dvp.ProofReceipt) error
}

type consolidationDepositWaiter interface {
	WaitUntilDepositIsConfirmed(ctx context.Context, deposit *types.DvpDeposit) error
}

// Compile-time assertions to verify concrete types implement the interfaces
var (
	_ consolidationDepositRepository    = (*repository.DvpDepositRepository)(nil)
	_ consolidationCommitmentCalculator = (*CommitmentCalculator)(nil)
	_ consolidationProofService         = (*ProofService)(nil)
	_ consolidationDvpClient            = (*contractclient.DvpClient)(nil)
)

var (
	_ ConsolidationEnygmaClient         = (*contractclient.EnygmaClient)(nil)
	_ ConsolidationDvpIntegrationClient = (*contractclient.DvpIntegrationClient)(nil)
)

var _ consolidationDepositWaiter = (*DepositWaiter)(nil)
var _ consolidationTxManager = (*repository.TransactionManager)(nil)

type ConsolidationService struct {
	conf                    ConsolidationConfig
	depositRepository       consolidationDepositRepository
	spendKeyProvider        consolidationSpendKeyProvider
	commitmentCalculator    consolidationCommitmentCalculator
	proofService            consolidationProofService
	dvpClient               consolidationDvpClient
	enygmaClient            ConsolidationEnygmaClient
	enygmaIntegrationClient ConsolidationDvpIntegrationClient
	depositWaiter           consolidationDepositWaiter
	txManager               consolidationTxManager
}

func NewConsolidationService(
	conf ConsolidationConfig,
	depositRepository consolidationDepositRepository,
	spendKeyProvider consolidationSpendKeyProvider,
	commitmentCalculator consolidationCommitmentCalculator,
	proofService consolidationProofService,
	dvpClient consolidationDvpClient,
	enygmaClient ConsolidationEnygmaClient,
	enygmaIntegrationClient ConsolidationDvpIntegrationClient,
	depositWaiter consolidationDepositWaiter,
	txManager consolidationTxManager,
) *ConsolidationService {
	return &ConsolidationService{
		conf:                    conf,
		depositRepository:       depositRepository,
		spendKeyProvider:        spendKeyProvider,
		commitmentCalculator:    commitmentCalculator,
		proofService:            proofService,
		dvpClient:               dvpClient,
		enygmaClient:            enygmaClient,
		enygmaIntegrationClient: enygmaIntegrationClient,
		depositWaiter:           depositWaiter,
		txManager:               txManager,
	}
}

// Consolidates a batch of deposits
// Creates a new deposit with the consolidation amount
// Creates a new deposit for the remainder(if any)
// Returns up to 2 new deposits. Minimum 1.
func (s *ConsolidationService) ConsolidateERC1155Deposits(
	ctx context.Context,
	id string,
	sourceViewPublicKey []byte,
	depositsToConsolidate []*types.DvpDeposit,
	consolidationAmount *big.Int,
) ([]*types.DvpDeposit, error) {
	totalAmountOfDeposits := CalculateTotalAmountOfDeposits(depositsToConsolidate)

	if err := s.validateDeposits(depositsToConsolidate); err != nil {
		return nil, fmt.Errorf("validate ERC1155 deposits: %w", err)
	}

	if consolidationAmount.Cmp(totalAmountOfDeposits) > 0 {
		return nil, fmt.Errorf("consolidation amount is greater than the total amount of deposits")
	}

	userAddress := depositsToConsolidate[0].UserAddress
	tokenAddress := depositsToConsolidate[0].TokenAddress
	tokenID := depositsToConsolidate[0].TokenID

	spendKey, err := s.spendKeyProvider.GetPaymentSpendKey(ctx, &keyspb.GetPaymentSpendKeyRequest{})
	if err != nil {
		return nil, fmt.Errorf("get payment spend key for consolidation: %w", err)
	}
	spendPublicKey := new(big.Int).SetBytes(spendKey.GetPublicKey())

	consolidationSalt, _, err := cryptography.GenerateSalt(sourceViewPublicKey)
	if err != nil {
		return nil, fmt.Errorf("generate salt for ERC1155 consolidation: %w", err)
	}

	consolidationCommitment, err := s.commitmentCalculator.CalculateERC1155Commitment(
		spendPublicKey,
		consolidationSalt,
		tokenAddress,
		tokenID,
		consolidationAmount,
	)
	if err != nil {
		return nil, fmt.Errorf("calculate ERC1155 consolidation commitment: %w", err)
	}

	// Generate proof BEFORE locking deposits — proof generation is pure computation
	// and doesn't need deposits to be locked.
	proof, err := s.proofService.GenerateERC1155JSProof(
		ctx,
		sourceViewPublicKey,
		consolidationCommitment,
		spendPublicKey,
		consolidationSalt,
		consolidationAmount,
		userAddress,
		tokenAddress,
		tokenID,
		depositsToConsolidate,
	)
	if err != nil {
		return nil, err
	}

	depositCommitments := make([]string, 0)
	commitmentNullifierMap := make(map[string]string)
	for i, deposit := range depositsToConsolidate {
		depositCommitments = append(depositCommitments, deposit.Commitment.String())
		commitmentNullifierMap[deposit.Commitment.String()] = proof.Nullifiers[i].String()
	}

	// Lock deposits, set nullifiers, and create the consolidated deposit atomically
	// right before the blockchain call. If the relayer crashes before this transaction,
	// no deposits are locked. If it crashes after, the consolidation can be retried.
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if txErr := s.depositRepository.BatchUpdateStatusForCommitments(txCtx, depositCommitments, types.DvpDepositLocked); txErr != nil {
			return txErr
		}
		if txErr := s.depositRepository.BatchUpsertNullifiers(txCtx, commitmentNullifierMap); txErr != nil {
			return txErr
		}
		return s.depositRepository.CreateDeposit(
			txCtx,
			&types.DvpDeposit{
				UserAddress:  userAddress,
				Salt:         consolidationSalt,
				Commitment:   consolidationCommitment,
				TokenAmount:  consolidationAmount,
				TokenAddress: tokenAddress,
				TokenID:      tokenID,
				TokenType:    types.DvpERC1155,
				Status:       types.DvpDepositPending,
			},
		)
	})
	if err != nil {
		return nil, fmt.Errorf("create consolidated deposit: %w", err)
	}

	err = s.dvpClient.MixFundsERC1155(ctx, id, common.HexToAddress(tokenAddress), proof)
	if err != nil {
		if cErr := s.depositRepository.BatchUpdateStatusForCommitments(ctx, depositCommitments, types.DvpDepositUnspent); cErr != nil {
			return nil, cErr
		}
		return nil, fmt.Errorf("failed to consolidate ERC1155 funds: %w", err)
	}

	newDeposits, err := s.getConsolidatedDeposits(ctx, proof.Commitments)
	if err != nil {
		return nil, fmt.Errorf("get consolidated ERC1155 deposits: %w", err)
	}

	return newDeposits, nil
}

// Consolidates a batch of deposits
// Creates a new deposit with the consolidation amount
// Creates a new deposit for the remainder(if any)
// Returns up to 2 new deposits. Minimum 1.
func (s *ConsolidationService) ConsolidateEnygmaDeposits(
	ctx context.Context,
	chainEventID string,
	sourceViewPublicKey []byte,
	depositsToConsolidate []*types.DvpDeposit,
	consolidationAmount *big.Int,
) ([]*types.DvpDeposit, error) {
	totalAmountOfDeposits := CalculateTotalAmountOfDeposits(depositsToConsolidate)

	if err := s.validateDeposits(depositsToConsolidate); err != nil {
		return nil, fmt.Errorf("validate enygma deposits: %w", err)
	}

	if consolidationAmount.Cmp(totalAmountOfDeposits) > 0 {
		return nil, fmt.Errorf("consolidation amount is greater than the total amount of deposits")
	}

	userAddress := depositsToConsolidate[0].UserAddress
	tokenAddress := depositsToConsolidate[0].TokenAddress
	tokenID := depositsToConsolidate[0].TokenID

	spendKey, err := s.spendKeyProvider.GetPaymentSpendKey(ctx, &keyspb.GetPaymentSpendKeyRequest{})
	if err != nil {
		return nil, fmt.Errorf("get payment spend key for consolidation: %w", err)
	}
	spendPublicKey := new(big.Int).SetBytes(spendKey.GetPublicKey())

	nftCommitment := big.NewInt(0)

	consolidationSalt, _, err := cryptography.GenerateSalt(sourceViewPublicKey)
	if err != nil {
		return nil, fmt.Errorf("generate salt for Enygma consolidation: %w", err)
	}

	// Generate proof BEFORE locking deposits — proof generation is pure computation.
	proof, err := s.proofService.GenerateEnygmaJSProof(
		ctx,
		sourceViewPublicKey,
		nftCommitment,
		spendPublicKey,
		consolidationSalt,
		consolidationAmount,
		tokenAddress,
		depositsToConsolidate,
	)
	if err != nil {
		return nil, err
	}

	depositCommitments := make([]string, 0)
	commitmentNullifierMap := make(map[string]string)
	for i, deposit := range depositsToConsolidate {
		depositCommitments = append(depositCommitments, deposit.Commitment.String())
		commitmentNullifierMap[deposit.Commitment.String()] = proof.Nullifiers[i].String()
	}

	// Create the new consolidated deposit with the exact amount
	consolidationCommitment, err := s.commitmentCalculator.CalculatePaymentCommitment(
		spendPublicKey,
		consolidationSalt,
		consolidationAmount,
		tokenAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("calculate Enygma consolidation commitment: %w", err)
	}

	enygmaIntegrationAddress, err := s.enygmaClient.GetDvpIntegrationContractAddress(ctx, common.HexToAddress(tokenAddress))
	if err != nil {
		return nil, fmt.Errorf("get dvp integration contract address: %w", err)
	}
	// Lock deposits, set nullifiers, and create the consolidated deposit atomically
	// right before the blockchain call. If the relayer crashes before this transaction,
	// no deposits are locked. If it crashes after, the consolidation can be retried.
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if txErr := s.depositRepository.BatchUpdateStatusForCommitments(txCtx, depositCommitments, types.DvpDepositLocked); txErr != nil {
			return txErr
		}
		if txErr := s.depositRepository.BatchUpsertNullifiers(txCtx, commitmentNullifierMap); txErr != nil {
			return txErr
		}
		return s.depositRepository.CreateDeposit(
			txCtx,
			&types.DvpDeposit{
				UserAddress:  userAddress,
				Salt:         consolidationSalt,
				Commitment:   consolidationCommitment,
				TokenAmount:  consolidationAmount,
				TokenAddress: tokenAddress,
				TokenID:      tokenID,
				TokenType:    types.DvpEnygma,
				Status:       types.DvpDepositPending,
			},
		)
	})
	if err != nil {
		return nil, fmt.Errorf("create consolidated deposit: %w", err)
	}

	err = s.enygmaIntegrationClient.ConsolidateFunds(ctx, chainEventID, enygmaIntegrationAddress, proof)
	if err != nil {
		if cErr := s.depositRepository.BatchUpdateStatusForCommitments(ctx, depositCommitments, types.DvpDepositUnspent); cErr != nil {
			return nil, cErr
		}
		return nil, fmt.Errorf("failed to consolidate Enygma funds: %w", err)
	}

	newDeposits, err := s.getConsolidatedDeposits(ctx, proof.Commitments)
	if err != nil {
		return nil, fmt.Errorf("get consolidated enygma deposits: %w", err)
	}

	return newDeposits, nil
}

// Recursively consolidates deposits until they fit in MaxNumberOfJSDepositsPerJoinSplitProof length
func (s *ConsolidationService) PrepareDepositsForJSProof(
	ctx context.Context,
	id string,
	sourceViewPublicKey []byte,
	deposits []*types.DvpDeposit,
) ([]*types.DvpDeposit, error) {
	// there are enough deposits to generate a join split proof. No need to consolidate.
	if len(deposits) <= s.conf.MaxNumberOfJSDeposits {
		return deposits, nil
	}

	if err := s.validateDeposits(deposits); err != nil {
		return nil, fmt.Errorf("validate deposits for JS proof: %w", err)
	}

	tokenType := deposits[0].TokenType

	slog.Debug("Consolidating deposits...", slog.String("Length", strconv.Itoa(len(deposits))))

	currentIterationBatch := deposits[0:s.conf.MaxNumberOfJSDeposits]
	nextIterationBatch := deposits[s.conf.MaxNumberOfJSDeposits:]
	currentIterationTotalAmount := CalculateTotalAmountOfDeposits(currentIterationBatch)

	var newDeposits []*types.DvpDeposit
	var err error

	switch tokenType {
	case types.DvpERC1155:
		newDeposits, err = s.ConsolidateERC1155Deposits(ctx, id, sourceViewPublicKey, currentIterationBatch, currentIterationTotalAmount)
	case types.DvpEnygma:
		newDeposits, err = s.ConsolidateEnygmaDeposits(ctx, id, sourceViewPublicKey, currentIterationBatch, currentIterationTotalAmount)
	default:
		err = fmt.Errorf("unsupported token to consolidate: %d", int(tokenType))
	}

	if err != nil {
		return nil, fmt.Errorf("consolidate deposits: %w", err)
	}

	finalDeposits := make([]*types.DvpDeposit, 0, len(newDeposits)+len(nextIterationBatch))
	finalDeposits = append(finalDeposits, newDeposits...)
	finalDeposits = append(finalDeposits, nextIterationBatch...)

	return s.PrepareDepositsForJSProof(ctx, id, sourceViewPublicKey, finalDeposits)
}

// validate that all deposits are of the same token
func (s *ConsolidationService) validateDeposits(deposits []*types.DvpDeposit) error {
	if len(deposits) == 0 {
		return fmt.Errorf("no deposits")
	}

	userAddress := deposits[0].UserAddress
	tokenAddress := deposits[0].TokenAddress
	tokenID := deposits[0].TokenID
	tokenType := deposits[0].TokenType

	for _, deposit := range deposits {
		if deposit.UserAddress != userAddress {
			return fmt.Errorf("deposits are not for the same user")
		}
		if deposit.TokenAddress != tokenAddress {
			return fmt.Errorf("deposits are not for the same token")
		}
		if deposit.TokenID != tokenID {
			return fmt.Errorf("deposits are not for the same token")
		}
		if deposit.TokenType != tokenType {
			return fmt.Errorf("deposits are not for the same token")
		}
	}

	return nil
}

func (s *ConsolidationService) getConsolidatedDeposits(ctx context.Context, commitments []*big.Int) ([]*types.DvpDeposit, error) {
	deposits := make([]*types.DvpDeposit, 0)

	for _, commitment := range commitments {
		deposit, err := s.depositRepository.GetDepositByCommitment(ctx, commitment)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch deposit: %w", err)
		}
		// the commitment is not in the database. We assume the deposit is dummy(0 change amount). Skip it...
		if deposit == nil {
			continue
		}

		err = s.depositWaiter.WaitUntilDepositIsConfirmed(ctx, deposit)
		if err != nil {
			return nil, fmt.Errorf("wait until deposit is confirmed: %w", err)
		}

		deposits = append(deposits, deposit)
	}

	return deposits, nil
}

func CalculateTotalAmountOfDeposits(deposits []*types.DvpDeposit) *big.Int {
	totalAmount := new(big.Int)
	for _, deposit := range deposits {
		totalAmount = totalAmount.Add(totalAmount, deposit.TokenAmount)
	}
	return totalAmount
}
