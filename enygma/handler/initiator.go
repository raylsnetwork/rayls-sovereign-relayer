package handler

//go:generate moq --skip-ensure --pkg handler_test -out initiator_mock_test.go . initiatorEnygmaHistoryRepository initiatorDvpDepositRepository initiatorTxManager initiatorKeysClient initiatorEndpointClient initiatorEnygmaHandlerClient initiatorCommitmentCalculator initiatorDepositFinder initiatorDepositConsolidator initiatorDvpProofGenerator initiatorTracer initiatorRetryService initiatorEnygmaBatcher initiatorEnygmaExecutor initiatorEnygmaFinalizationService initiatorEnygmaCreationService

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cryptography"
	keyspb "github.com/raylsnetwork/rayls-sovereign-relayer/cts/gen/proto/keys"
	"github.com/raylsnetwork/rayls-sovereign-relayer/dvp"
	dvpService "github.com/raylsnetwork/rayls-sovereign-relayer/dvp/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/enygma/service"
	telemetry "github.com/raylsnetwork/rayls-sovereign-relayer/otel"
	repository2 "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

var ErrNoSupplyChanges = errors.New("no supply changes")

type initiatorEnygmaHistoryRepository interface {
	InsertEnygmaHistory(ctx context.Context, history types.EnygmaHistory) error
}

type initiatorDvpDepositRepository interface {
	CreateDeposit(ctx context.Context, deposit *types.DvpDeposit) error
	BatchUpdateStatusForCommitments(ctx context.Context, commitments []string, status types.DvpDepositStatus) error
	BatchUpsertNullifiers(ctx context.Context, commitmentNullifierMap map[string]string) error
}

type initiatorTxManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type initiatorKeysClient interface {
	GetPaymentSpendKey(ctx context.Context, in *keyspb.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error)
	GetViewPublicKey(ctx context.Context, in *keyspb.GetViewPublicKeyRequest, opts ...grpc.CallOption) (*keyspb.GetViewPublicKeyResponse, error)
}

type initiatorEndpointClient interface {
	GetResourceAddress(ctx context.Context, resourceId string) (common.Address, error)
	ReceivePayload(
		ctx context.Context,
		fromChainId *big.Int,
		from common.Address,
		to common.Address,
		data EndpointV1.RaylsMessage,
		messageId [32]byte,
	) (common.Hash, error)
}

type initiatorEnygmaHandlerClient interface {
	ReceiveWithdraw(ctx context.Context, chainEventID string, tokenAddress common.Address, to common.Address, value *big.Int, referenceId [32]byte) error
}

// DVP ports
type initiatorCommitmentCalculator interface {
	CalculatePaymentCommitment(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error)
}
type initiatorDepositFinder interface {
	FindEnygmaDeposits(ctx context.Context, userAddress string, tokenAddress string, paymentAmount *big.Int) ([]*types.DvpDeposit, error)
}

type initiatorDepositConsolidator interface {
	PrepareDepositsForJSProof(ctx context.Context, id string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error)
	ConsolidateEnygmaDeposits(
		ctx context.Context,
		id string,
		sourceViewPublicKey []byte,
		depositsToConsolidate []*types.DvpDeposit,
		consolidationAmount *big.Int,
	) ([]*types.DvpDeposit, error)
}

type initiatorDvpProofGenerator interface {
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

type initiatorTracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

type initiatorRetryService interface {
	RetryOperation(
		ctx context.Context,
		operationName string,
		maxRetries int,
		blockNumber uint64,
		executeOperation func(ctx context.Context, nextBlockNumber uint64) error,
	) (uint64, error)
}

type initiatorEnygmaBatcher interface {
	CreateBatchesWithAnonimity(
		ctx context.Context,
		resourceId string,
		blockNumber *big.Int,
		txsByChainID map[string][]*types.EnygmaTransferBatchTx,
	) ([]*types.EnygmaTransferBatch, error)
}

// Service interfaces for executor, finalization, and creation services
type initiatorEnygmaExecutor interface {
	ExecuteEnygmaSupplyUpdate(ctx context.Context, batchID string, resourceId string, blockNumber uint64, batch types.EnygmaSupplyUpdate, enygmaAddress common.Address) error
	ExecuteEnygmaCrossTransfer(ctx context.Context, batchID string, blockNumber uint64, resourceId string, batch map[string][]*types.EnygmaTransferBatchTx, enygmaAddress common.Address) error
	ExecuteEnygmaDeposit(ctx context.Context, id string, resourceId string, amount *big.Int, blockNumber uint64, commitment *big.Int, salt *big.Int, from common.Address, txHash common.Hash, enygmaAddress common.Address) error
	ExecuteEnygmaWithdrawal(ctx context.Context, chainEventID string, resourceId string, amount *big.Int, deposits []*types.DvpDeposit, blockNumber uint64, enygmaAddress common.Address, proof *dvp.ProofReceipt, from common.Address, txHash common.Hash) error
}

type initiatorEnygmaFinalizationService interface {
	ExecuteFinalization(ctx context.Context, id string, blockNumber uint64, resourceId string) error
}

type initiatorEnygmaCreationService interface {
	CreateEnygma(ctx context.Context, resourceId string, fromChainId *big.Int, blockNumberPrivateHub *big.Int) error
}

// Compile-time interface implementation checks
var (
	_ initiatorEnygmaHistoryRepository   = (*repository2.EnygmaHistoryRepository)(nil)
	_ initiatorDvpDepositRepository      = (*repository2.DvpDepositRepository)(nil)
	_ initiatorDepositFinder             = (*dvpService.DepositFinder)(nil)
	_ initiatorCommitmentCalculator      = (*dvpService.CommitmentCalculator)(nil)
	_ initiatorDepositConsolidator       = (*dvpService.ConsolidationService)(nil)
	_ initiatorDvpProofGenerator         = (*dvpService.ProofService)(nil)
	_ initiatorRetryService              = (*service.RetryService)(nil)
	_ initiatorEnygmaBatcher             = (*service.EnygmaBatcher)(nil)
	_ initiatorEnygmaExecutor            = (*service.EnygmaExecutor)(nil)
	_ initiatorEnygmaFinalizationService = (*service.EnygmaFinalizationService)(nil)
	_ initiatorEnygmaCreationService     = (*service.EnygmaCreationService)(nil)
	_ initiatorTxManager                 = (*repository2.TransactionManager)(nil)
)

type InitiatorConfig struct {
	DefaultContextTimeout time.Duration
	ViewPublicKey         *big.Int
}

type Initiator struct {
	conf                    InitiatorConfig
	keysClient              initiatorKeysClient
	enygmaHistoryRepository initiatorEnygmaHistoryRepository
	dvpDepositRepository    initiatorDvpDepositRepository
	enygmaHandlerClient     initiatorEnygmaHandlerClient
	ccEndpointClient        initiatorEndpointClient
	plEndpointClient        initiatorEndpointClient
	plChainId               *big.Int
	depositFinder           initiatorDepositFinder
	commitmentCalculator    initiatorCommitmentCalculator
	depositConsolidator     initiatorDepositConsolidator
	dvpProofGen             initiatorDvpProofGenerator
	tracer                  initiatorTracer
	retryService            initiatorRetryService
	executor                initiatorEnygmaExecutor
	finalization            initiatorEnygmaFinalizationService
	creationService         initiatorEnygmaCreationService
	txManager               initiatorTxManager
}

func NewInitiator(
	conf InitiatorConfig,
	keysClient initiatorKeysClient,
	enygmaHandlerClient initiatorEnygmaHandlerClient,
	ccEndpointClient initiatorEndpointClient,
	plEndpointClient initiatorEndpointClient,
	enygmaHistoryRepository initiatorEnygmaHistoryRepository,
	dvpDepositRepository initiatorDvpDepositRepository,
	plChainId *big.Int,
	depositFinder initiatorDepositFinder,
	commitmentCalculator initiatorCommitmentCalculator,
	depositConsolidator initiatorDepositConsolidator,
	dvpProofGen initiatorDvpProofGenerator,
	tracer initiatorTracer,
	retryService initiatorRetryService,
	executor initiatorEnygmaExecutor,
	finalization initiatorEnygmaFinalizationService,
	creationService initiatorEnygmaCreationService,
	txManager initiatorTxManager,
) *Initiator {
	return &Initiator{
		conf:                    conf,
		keysClient:              keysClient,
		enygmaHistoryRepository: enygmaHistoryRepository,
		dvpDepositRepository:    dvpDepositRepository,
		plEndpointClient:        plEndpointClient,
		plChainId:               plChainId,
		enygmaHandlerClient:     enygmaHandlerClient,
		ccEndpointClient:        ccEndpointClient,
		depositFinder:           depositFinder,
		commitmentCalculator:    commitmentCalculator,
		depositConsolidator:     depositConsolidator,
		dvpProofGen:             dvpProofGen,
		tracer:                  tracer,
		retryService:            retryService,
		executor:                executor,
		finalization:            finalization,
		creationService:         creationService,
		txManager:               txManager,
	}
}

// HandleEnygmaCreation handles the creation of a new Enygma token.
// Returns the block number to be used for finalization (0 if no finalization is needed).
func (s *Initiator) HandleEnygmaCreation(
	ctx context.Context,
	chainEventID string,
	resourceId string,
	blockNumber uint64,
	initialSupply *big.Int,
) (uint64, error) {
	ctx, span := s.tracer.Start(ctx, telemetry.SPAN_HANDLE_ENYGMA_CREATION)
	defer span.End()

	span.SetAttributes(
		attribute.String(telemetry.ATTR_RESOURCE_ID, resourceId),
		attribute.Int(telemetry.ATTR_BLOCK_NUMBER, int(blockNumber)), //nolint:gosec // block numbers fit in int
		attribute.String(telemetry.ATTR_INITIAL_SUPPLY, initialSupply.String()),
	)

	err := s.creationService.CreateEnygma(
		ctx,
		resourceId,
		s.plChainId,
		new(big.Int).SetUint64(blockNumber),
	)
	if err != nil {
		return 0, fmt.Errorf("create enygma: %w", err)
	}

	// There is no initial supply, so we don't need to initiate a supply update on CC.
	if initialSupply.Cmp(big.NewInt(0)) == 0 {
		return 0, nil
	}

	enygmaAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return 0, fmt.Errorf("get resource address for initial mint: %w", err)
	}

	initialMint := types.EnygmaSupplyUpdate{
		Amount: initialSupply,
		Type:   types.EnygmaMint,
	}

	maxRetries := 30
	blockNumber, err = s.retryService.RetryOperation(
		ctx,
		"initial minting",
		maxRetries,
		blockNumber,
		func(ctx context.Context, nextBlockNumber uint64) error {
			return s.executor.ExecuteEnygmaSupplyUpdate(ctx, chainEventID, resourceId, nextBlockNumber, initialMint, enygmaAddress)
		},
	)
	if err != nil {
		return 0, fmt.Errorf("initial minting retry: %w", err)
	}

	return blockNumber, nil
}

// HandleEnygmaSupplyUpdates handles supply updates for Enygma tokens
func (s *Initiator) HandleEnygmaSupplyUpdates(
	ctx context.Context,
	batchID string,
	resourceId string,
	blockNumber uint64,
	batch types.EnygmaSupplyUpdate,
) (uint64, error) {
	ctx, span := s.tracer.Start(ctx, telemetry.SPAN_HANDLE_ENYGMA_SUPPLY_UPDATE)
	defer span.End()

	if batch.Amount.Cmp(big.NewInt(0)) == 0 {
		return blockNumber, ErrNoSupplyChanges
	}

	enygmaAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return blockNumber, fmt.Errorf("get resource address for supply update: %w", err)
	}

	span.SetAttributes(
		attribute.String(telemetry.ATTR_RESOURCE_ID, resourceId),
		attribute.Int(telemetry.ATTR_BLOCK_NUMBER, int(blockNumber)), //nolint:gosec // block numbers fit in int
		attribute.String(telemetry.ATTR_ENYGMA_ADDRESS, enygmaAddress.Hex()),
	)

	maxRetries := 30
	blockNumber, err = s.retryService.RetryOperation(
		ctx,
		"supply update",
		maxRetries,
		blockNumber,
		func(ctx context.Context, nextBlockNumber uint64) error {
			return s.executor.ExecuteEnygmaSupplyUpdate(ctx, batchID, resourceId, nextBlockNumber, batch, enygmaAddress)
		},
	)
	if err != nil {
		return blockNumber, fmt.Errorf("supply update retry: %w", err)
	}

	return blockNumber, nil
}

// HandleEnygmaCrossTransferBatch handles cross-chain transfer batches
func (s *Initiator) HandleEnygmaCrossTransfer(
	parentCtx context.Context,
	batchID string,
	resourceId string,
	blockNumber uint64,
	batch map[string][]*types.EnygmaTransferBatchTx,
) (uint64, error) {
	ctx, span := s.tracer.Start(parentCtx, telemetry.SPAN_HANDLE_ENYGMA_CROSS_TRANSFER)
	defer span.End()

	enygmaAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return blockNumber, fmt.Errorf("get resource address for cross transfer: %w", err)
	}

	span.SetAttributes(
		attribute.String(telemetry.ATTR_RESOURCE_ID, resourceId),
		attribute.Int(telemetry.ATTR_BLOCK_NUMBER, int(blockNumber)), //nolint:gosec // block numbers fit in int
		attribute.Int(telemetry.ATTR_BATCH_LENGTH, len(batch)),
		attribute.String(telemetry.ATTR_ENYGMA_ADDRESS, enygmaAddress.Hex()),
	)

	maxRetries := 30
	// Amount of retries within which the SyncService must finalize the Enygma Checkpoint & the operation succeeds.
	blockNumber, err = s.retryService.RetryOperation(
		ctx,
		"cross transfer",
		maxRetries,
		blockNumber,
		func(ctx context.Context, nextBlockNumber uint64) error {
			return s.executor.ExecuteEnygmaCrossTransfer(ctx, batchID, nextBlockNumber, resourceId, batch, enygmaAddress)
		},
	)
	if err != nil {
		return blockNumber, fmt.Errorf("cross transfer retry: %w", err)
	}

	span.SetStatus(codes.Ok, telemetry.STATUS_SUCCESS_TRANSFER_EXECUTED)
	return blockNumber, nil
}

func (s *Initiator) HandleEnygmaDeposit(
	ctx context.Context,
	chainEventID string,
	blockNumber uint64,
	resourceId string,
	referenceId [32]byte,
	from common.Address,
	amount *big.Int,
	txHash common.Hash,
) (uint64, error) {
	var err error

	enygmaAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return blockNumber, fmt.Errorf("get resource address for deposit: %w", err)
	}

	viewKey, err := s.keysClient.GetViewPublicKey(ctx, &keyspb.GetViewPublicKeyRequest{BlockNumber: blockNumber})
	if err != nil {
		return blockNumber, fmt.Errorf("failed to get rayls view key pair for deposit: %w", err)
	}

	publicKeyBytes, err := hex.DecodeString(viewKey.GetPublicKey())
	if err != nil {
		return blockNumber, fmt.Errorf("decoding rayls view public key: %w", err)
	}
	salt, _, err := cryptography.GenerateSalt(publicKeyBytes)
	if err != nil {
		return blockNumber, fmt.Errorf("failed to generate salt for deposit: %w", err)
	}

	spendKey, err := s.keysClient.GetPaymentSpendKey(ctx, &keyspb.GetPaymentSpendKeyRequest{})
	if err != nil {
		return blockNumber, fmt.Errorf("failed to get payment spend key for deposit: %w", err)
	}
	spendPublicKey := new(big.Int).SetBytes(spendKey.GetPublicKey())

	commitment, err := s.commitmentCalculator.CalculatePaymentCommitment(
		spendPublicKey,
		salt,
		amount,
		enygmaAddress.Hex(),
	)
	if err != nil {
		return blockNumber, fmt.Errorf("calculate payment commitment for deposit: %w", err)
	}

	// Optimistically create the deposit. We must ensure the deposit is created before we submit the tx to CC.
	// It'll be used by the merkle tree service, once deposit is verified on CC.
	deposit := types.DvpDeposit{
		UserAddress:  from.Hex(),
		Salt:         salt,
		TokenAmount:  amount,
		Commitment:   commitment,
		TokenAddress: enygmaAddress.Hex(),
		TokenID:      "",
		TokenType:    types.DvpEnygma,
		Status:       types.DvpDepositPending,
	}
	err = s.dvpDepositRepository.CreateDeposit(ctx, &deposit)
	if err != nil {
		return blockNumber, fmt.Errorf("create deposit record: %w", err)
	}

	maxRetries := 30
	blockNumber, err = s.retryService.RetryOperation(
		ctx,
		"deposit",
		maxRetries,
		blockNumber,
		func(ctx context.Context, nextBlockNumber uint64) error {
			return s.executor.ExecuteEnygmaDeposit(
				ctx,
				chainEventID,
				resourceId,
				amount,
				nextBlockNumber,
				commitment,
				salt,
				from,
				txHash,
				enygmaAddress,
			)
		},
	)
	if err != nil {
		return blockNumber, fmt.Errorf("deposit retry: %w", err)
	}

	return blockNumber, nil
}

func (s *Initiator) HandleEnygmaWithdrawal(
	ctx context.Context,
	chainEventID string,
	blockNumber uint64,
	resourceId string,
	referenceId [32]byte,
	to common.Address,
	amount *big.Int,
	txHash common.Hash,
) (uint64, error) {
	var err error

	enygmaAddress, err := s.ccEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return blockNumber, fmt.Errorf("get resource address for withdrawal: %w", err)
	}

	viewKey, err := s.keysClient.GetViewPublicKey(ctx, &keyspb.GetViewPublicKeyRequest{BlockNumber: blockNumber})
	if err != nil {
		return blockNumber, fmt.Errorf("get rayls view key pair: %w", err)
	}

	viewPubKey, err := hex.DecodeString(viewKey.GetPublicKey())
	if err != nil {
		return blockNumber, fmt.Errorf("decoding rayls view public key: %w", err)
	}

	userDeposits, err := s.depositFinder.FindEnygmaDeposits(ctx, to.Hex(), enygmaAddress.Hex(), amount)
	if err != nil {
		return blockNumber, fmt.Errorf("find enygma deposits: %w", err)
	}
	jsDeposits, err := s.depositConsolidator.PrepareDepositsForJSProof(ctx, chainEventID, viewPubKey, userDeposits)
	if err != nil {
		return blockNumber, fmt.Errorf("prepare deposits for JS proof: %w", err)
	}
	totalDepositAmount := CalculateTotalAmountOfDeposits(jsDeposits)

	// User do not have exact withdrawal amount deposited. We need to consolidate those deposits in exact withdrawal amount
	if amount.Cmp(totalDepositAmount) < 0 {
		jsDeposits, err = s.depositConsolidator.ConsolidateEnygmaDeposits(ctx, chainEventID, viewPubKey, jsDeposits, amount)
		if err != nil {
			return blockNumber, fmt.Errorf("consolidate enygma deposits: %w", err)
		}
		// Get first deposit that match the withdrawal amount. In case we have 2 deposits with same amount, we will use the first one for the withdraw
		for _, deposit := range jsDeposits {
			if deposit.TokenAmount.Cmp(amount) == 0 {
				jsDeposits = []*types.DvpDeposit{deposit}
				break
			}
		}
	}

	if len(jsDeposits) == 0 {
		return blockNumber, errors.New("no deposits found to withdraw from")
	}

	depositCommitments := make([]string, 0)
	for _, deposit := range jsDeposits {
		depositCommitments = append(depositCommitments, deposit.Commitment.String())
	}

	spendKey, err := s.keysClient.GetPaymentSpendKey(ctx, &keyspb.GetPaymentSpendKeyRequest{})
	if err != nil {
		return blockNumber, fmt.Errorf("get payment spend key: %w", err)
	}
	spendPublicKey := new(big.Int).SetBytes(spendKey.GetPublicKey())

	nftCommitment := big.NewInt(0)

	destSalt, _, err := cryptography.GenerateSalt(viewPubKey)
	if err != nil {
		return blockNumber, fmt.Errorf("generating salt: %w", err)
	}

	jsProof, err := s.dvpProofGen.GenerateEnygmaJSProof(
		ctx,
		viewPubKey,
		nftCommitment,
		spendPublicKey,
		destSalt,
		amount,
		enygmaAddress.Hex(),
		jsDeposits,
	)
	if err != nil {
		return blockNumber, fmt.Errorf("generate enygma JS proof: %w", err)
	}

	commitmentNullifierMap := make(map[string]string)
	for i, deposit := range jsDeposits {
		commitmentNullifierMap[deposit.Commitment.String()] = jsProof.Nullifiers[i].String()
	}

	// Lock deposits and set nullifiers atomically in a single transaction.
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if txErr := s.dvpDepositRepository.BatchUpdateStatusForCommitments(txCtx, depositCommitments, types.DvpDepositLocked); txErr != nil {
			return txErr
		}
		return s.dvpDepositRepository.BatchUpsertNullifiers(txCtx, commitmentNullifierMap)
	})
	if err != nil {
		return blockNumber, fmt.Errorf("batch upsert nullifiers: %w", err)
	}

	maxRetries := 30
	blockNumber, err = s.retryService.RetryOperation(
		ctx,
		"withdrawal",
		maxRetries,
		blockNumber,
		func(ctx context.Context, nextBlockNumber uint64) error {
			return s.executor.ExecuteEnygmaWithdrawal(
				ctx,
				chainEventID,
				resourceId,
				amount,
				jsDeposits,
				nextBlockNumber,
				enygmaAddress,
				jsProof,
				to,
				txHash,
			)
		},
	)
	if err != nil {
		sErr := s.dvpDepositRepository.BatchUpdateStatusForCommitments(ctx, depositCommitments, types.DvpDepositUnspent)
		if sErr != nil {
			return blockNumber, fmt.Errorf("revert deposit status after withdrawal failure: %w", sErr)
		}

		return blockNumber, fmt.Errorf("withdrawal retry: %w", err)
	}

	// Retry withdrawEnygmaInPL because the PNH withdrawal already succeeded and is irreversible.
	// If this fails without retry, the deposit stays permanently Spent and all future attempts
	// to re-enter HandleEnygmaWithdrawal will fail at FindEnygmaDeposits with "not enough deposits".
	maxRetriesPL := 30
	blockNumber, err = s.retryService.RetryOperation(
		ctx,
		"withdrawEnygmaInPL",
		maxRetriesPL,
		blockNumber,
		func(ctx context.Context, nextBlockNumber uint64) error {
			return s.withdrawEnygmaInPL(ctx, chainEventID, resourceId, to, amount, referenceId)
		},
	)
	if err != nil {
		slog.Error(
			"Error while withdrawing Enygma in PL after retries",
			slog.String("resourceId", resourceId),
			slog.String("error", err.Error()),
		)
		return blockNumber, fmt.Errorf("withdraw enygma in PL retry: %w", err)
	}

	return blockNumber, nil
}

func (s *Initiator) withdrawEnygmaInPL(
	ctx context.Context,
	chainEventID string,
	resourceId string,
	to common.Address,
	amount *big.Int,
	referenceId [32]byte,
) error {
	enygmaAddress, err := s.plEndpointClient.GetResourceAddress(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("error getting contract address from resource ID %s: %w", resourceId, err)
	}
	txCtx, cancel := context.WithTimeout(ctx, s.conf.DefaultContextTimeout)
	defer cancel()

	err = s.enygmaHandlerClient.ReceiveWithdraw(txCtx, chainEventID, enygmaAddress, to, amount, referenceId)
	if err != nil {
		return fmt.Errorf("receive withdraw: %w", err)
	}

	return nil
}

func CalculateTotalAmountOfDeposits(deposits []*types.DvpDeposit) *big.Int {
	totalAmount := new(big.Int)
	for _, deposit := range deposits {
		totalAmount = totalAmount.Add(totalAmount, deposit.TokenAmount)
	}
	return totalAmount
}
