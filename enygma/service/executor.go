package service

//go:generate moq --skip-ensure --pkg service_test -out executor_mock_test.go . executorTracer executorEnygmaBatcher executorProofGenerator executorEnygmaHistoryRepository executorEnygmaClient ExecutorDvpIntegrationClient executorKeysClient executorCommitmentCalculator

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	keyspb "github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/keys"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/enygma"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/faultinjector"
	telemetry "github.com/raylsnetwork/rayls-privacy-relayer-api/otel"
	repository "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/txsim"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/txutil"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

const millisPerSecond = txutil.MillisPerSecond

type ExecutorConfig struct {
	DefaultContextTimeout time.Duration
	MaxNumberOfJSDeposits int
}

type executorTracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

type executorEnygmaBatcher interface {
	CreateBatchesWithAnonimity(
		ctx context.Context,
		resourceId string,
		blockNumber *big.Int,
		txsByChainID map[string][]*types.EnygmaTransferBatchTx,
	) ([]*types.EnygmaTransferBatch, error)
}

type executorProofGenerator interface {
	GenerateTransferProof(
		ctx context.Context,
		params enygma.TransferProofParams,
	) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error)
	GenerateDepositProof(
		ctx context.Context,
		params enygma.DepositProofParams,
	) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error)
	GenerateWithdrawProof(
		ctx context.Context,
		params enygma.WithdrawProofParams,
	) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error)
}

type executorEnygmaHistoryRepository interface {
	InsertEnygmaHistory(ctx context.Context, history types.EnygmaHistory) error
}

type executorEnygmaClient interface {
	SupplyUpdate(
		ctx context.Context,
		batchID string,
		tokenAddress common.Address,
		senderChainId *big.Int,
		blockNumber *big.Int,
		update types.EnygmaSupplyUpdate,
	) error
	TransferBatch(
		ctx context.Context,
		batchID string,
		tokenAddress common.Address,
		batches []*types.EnygmaTransferBatch,
		proof *types.EnygmaProofResponse,
		blockNumber *big.Int,
	) error
	GetDvpIntegrationContractAddress(ctx context.Context, tokenAddress common.Address) (common.Address, error)
}

type ExecutorDvpIntegrationClient interface {
	Deposit(ctx context.Context, chainEventID string, batches []*types.EnygmaTransferBatch, proof *types.EnygmaProofResponse, blockNumber *big.Int, chainId *big.Int, resourceId string, amount *big.Int, from common.Address, sourceTxHash common.Hash, dvpIntegrationAddress common.Address) error
	Withdraw(ctx context.Context, chainEventID string, batches []*types.EnygmaTransferBatch, proof *types.EnygmaProofResponse, blockNumber *big.Int, jsProof *dvp.ProofReceipt, chainId *big.Int, resourceId string, amount *big.Int, from common.Address, sourceTxHash common.Hash, dvpIntegrationAddress common.Address) error
}

type executorKeysClient interface {
	GetPaymentSpendKey(ctx context.Context, in *keyspb.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error)
}

type executorCommitmentCalculator interface {
	CalculatePaymentCommitment(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error)
}

type EnygmaExecutor struct {
	conf                    ExecutorConfig
	tracer                  executorTracer
	batcher                 executorEnygmaBatcher
	proofGen                executorProofGenerator
	enygmaHistoryRepository executorEnygmaHistoryRepository
	keysClient              executorKeysClient
	enygmaClient            executorEnygmaClient
	enygmaIntegrationClient ExecutorDvpIntegrationClient
	commitmentCalculator    executorCommitmentCalculator
	plChainId               *big.Int
}

func NewEnygmaExecutor(
	conf ExecutorConfig,
	tracer executorTracer,
	batcher executorEnygmaBatcher,
	proofGen executorProofGenerator,
	enygmaHistoryRepository executorEnygmaHistoryRepository,
	keysClient executorKeysClient,
	enygmaClient executorEnygmaClient,
	enygmaIntegrationClient ExecutorDvpIntegrationClient,
	commitmentCalculator executorCommitmentCalculator,
	plChainId *big.Int,
) *EnygmaExecutor {
	return &EnygmaExecutor{
		conf:                    conf,
		tracer:                  tracer,
		batcher:                 batcher,
		proofGen:                proofGen,
		enygmaHistoryRepository: enygmaHistoryRepository,
		keysClient:              keysClient,
		enygmaClient:            enygmaClient,
		enygmaIntegrationClient: enygmaIntegrationClient,
		commitmentCalculator:    commitmentCalculator,
		plChainId:               plChainId,
	}
}

// recordSoftFinalityLatency calculates and records the soft finality latency on the given span.
func (e *EnygmaExecutor) recordSoftFinalityLatency(span trace.Span, batches []*types.EnygmaTransferBatch) {
	if len(batches) == 0 || batches[0].SoftFinalityStartTimestamp <= 0 {
		return
	}

	softFinalityEndTime := time.Now().UnixMilli()
	softFinalityLatencyMillis := softFinalityEndTime - batches[0].SoftFinalityStartTimestamp
	latencySeconds := float64(softFinalityLatencyMillis) / millisPerSecond

	span.SetAttributes(
		attribute.Float64(telemetry.ATTR_SOFT_FINALITY_LATENCY, latencySeconds),
	)
}

func (e *EnygmaExecutor) ExecuteEnygmaSupplyUpdate(
	ctx context.Context,
	batchID string,
	resourceId string,
	blockNumber uint64,
	supplyUpdate types.EnygmaSupplyUpdate,
	enygmaAddress common.Address,
) error {
	ctx, span := e.tracer.Start(ctx, telemetry.SPAN_INITIATE_ENYGMA_SUPPLY_UPDATE)
	defer span.End()

	span.SetAttributes(
		attribute.String(telemetry.ATTR_RESOURCE_ID, resourceId),
		attribute.String(telemetry.ATTR_ENYGMA_SUPPLY_UPDATE_TYPE, supplyUpdate.Type.String()),
		attribute.String(telemetry.ATTR_AMOUNT, supplyUpdate.Amount.String()),
		//nolint:gosec // block numbers are within int range
		attribute.Int(telemetry.ATTR_BLOCK_NUMBER, int(blockNumber)),
	)

	if supplyUpdate.Type != types.EnygmaMint && supplyUpdate.Type != types.EnygmaBurn {
		return fmt.Errorf("invalid supply update type: must be mint or burn, got %d", supplyUpdate.Type)
	}

	blockNumberPrivateHub := new(big.Int).SetUint64(blockNumber)

	txCtx, cancel := context.WithTimeout(ctx, e.conf.DefaultContextTimeout)
	defer cancel()

	err := e.enygmaClient.SupplyUpdate(txCtx, batchID, enygmaAddress, e.plChainId, blockNumberPrivateHub, supplyUpdate)
	if err != nil {
		slog.Error("Error signing Enygma supply update", slog.String("resourceId", resourceId), slog.Any("err", err))
		return fmt.Errorf("signing enygma supply update: %w", err)
	}

	balanceChange := new(big.Int).Set(supplyUpdate.Amount)
	if supplyUpdate.Type == types.EnygmaBurn {
		balanceChange = balanceChange.Neg(balanceChange)
	}

	history := types.EnygmaHistory{
		ResourceId:            resourceId,
		FromChainId:           e.plChainId,
		RFactor:               big.NewInt(0),
		BlockNumberPrivateHub: blockNumberPrivateHub,
		BalanceChange:         balanceChange,
		EventType:             supplyUpdate.Type,
	}

	if err := e.enygmaHistoryRepository.InsertEnygmaHistory(ctx, history); err != nil {
		if !errors.Is(err, repository.ErrAlreadyProcessed) {
			return fmt.Errorf("inserting enygma history: %w", err)
		}
	}

	return nil
}

func (e *EnygmaExecutor) ExecuteEnygmaCrossTransfer(
	parentCtx context.Context,
	batchID string,
	blockNumber uint64,
	resourceId string,
	txsByChainID map[string][]*types.EnygmaTransferBatchTx,
	enygmaAddress common.Address,
) error {
	ctx, span := e.tracer.Start(parentCtx, telemetry.SPAN_INITIATE_ENYGMA_CROSS_TRANSFER)
	defer span.End()

	// Fault injection: drives the orchestrator's revert path under test by making
	// every executor call fail. With ActionError, the retryService exhausts its
	// retries and surfaces the error up, which moves the flow into the revert
	// branch (crossRevertMint on source PL). Used by the resilience suite —
	// Enygma_RevertCrash_DoubleRevertMint.ts in particular requires this hook.
	if err := faultinjector.Check("enygma.service.EnygmaExecutor.ExecuteEnygmaCrossTransfer.before_execute"); err != nil {
		return fmt.Errorf("fault injection at before_execute: %w", err)
	}

	blockNumberPrivateHub := new(big.Int).SetUint64(blockNumber)

	batches, err := e.batcher.CreateBatchesWithAnonimity(ctx, resourceId, blockNumberPrivateHub, txsByChainID)
	if err != nil {
		return fmt.Errorf("creating batches for cross transfer: %w", err)
	}

	senderPLTotalAmount := big.NewInt(0)
	for _, batch := range batches {
		for _, tx := range batch.Transactions {
			senderPLTotalAmount.Add(senderPLTotalAmount, tx.ToAmount)
		}
	}

	anonymityIndex := len(batches)
	proof, _, rValues, _, err := e.proofGen.GenerateTransferProof(
		ctx,
		enygma.TransferProofParams{
			ResourceId:     resourceId,
			AnonymityIndex: anonymityIndex,
			SenderAmount:   senderPLTotalAmount,
			BlockNumber:    blockNumberPrivateHub,
			Batches:        batches,
			TokenAddress:   enygmaAddress,
		},
	)
	if err != nil {
		slog.Error(
			"Error generating Enygma transfer proof",
			slog.String("resourceId", resourceId),
			slog.String("err", err.Error()),
		)
		return fmt.Errorf("generating enygma transfer proof: %w", err)
	}

	senderRFactor := big.NewInt(0)
	for i, batch := range batches {
		batch.ToRValueToAdd = rValues[i]
		if batch.ToChainID.Cmp(e.plChainId) == 0 {
			senderRFactor = rValues[i]
		}
	}

	txCtx, cancel := context.WithTimeout(ctx, e.conf.DefaultContextTimeout)
	defer cancel()

	err = e.enygmaClient.TransferBatch(txCtx, batchID, enygmaAddress, batches, proof, blockNumberPrivateHub)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, telemetry.STATUS_ERROR_FAILED_TO_SEND_TRANSFER_BATCH)
		// Mirror the destination-side programmability log: when the batch reverted on-chain, log the
		// raw revert (0x… selector + args) and a decoded `revert reason` NAME instead of dumping the
		// wrapped stack trace. Keeps programmability reverts readable and consistent across both sides.
		logAttrs := []any{slog.String("resourceId", resourceId), slog.Any("err", err)}
		var revertErr *contractclient.ErrorWithRevertData
		if errors.As(err, &revertErr) {
			logAttrs[1] = slog.String("err", revertErr.Error())
			if decoded, decErr := txsim.DecodeRevertBytes(revertErr.GetRevertData()); decErr == nil && !decoded.IsEmpty() {
				logAttrs = append(logAttrs, slog.String("revert reason", decoded.Reason()))
			}
		}
		slog.Error("Error signing Enygma transfer batch", logAttrs...)
		return fmt.Errorf("signing enygma transfer batch: %w", err)
	}

	history := types.EnygmaHistory{
		ResourceId:            resourceId,
		FromChainId:           e.plChainId,
		RFactor:               senderRFactor,
		BalanceChange:         new(big.Int).Neg(senderPLTotalAmount),
		BlockNumberPrivateHub: blockNumberPrivateHub,
		EventType:             types.EnygmaTransfer,
	}

	if err := e.enygmaHistoryRepository.InsertEnygmaHistory(ctx, history); err != nil {
		if !errors.Is(err, repository.ErrAlreadyProcessed) {
			return fmt.Errorf("inserting enygma history: %w", err)
		}
	}

	e.recordSoftFinalityLatency(span, batches)

	span.SetStatus(codes.Ok, telemetry.STATUS_SUCCESS_ENYGMA_CROSS_TRANSFER_INITIATED)
	return nil
}

func (e *EnygmaExecutor) ExecuteEnygmaDeposit(
	ctx context.Context,
	chainEventID string,
	resourceId string,
	amount *big.Int,
	blockNumber uint64,
	depositCommitment *big.Int,
	depositSalt *big.Int,
	from common.Address,
	txHash common.Hash,
	enygmaAddress common.Address,
) error {
	blockNumberPrivateHub := new(big.Int).SetUint64(blockNumber)

	batches, err := e.batcher.CreateBatchesWithAnonimity(ctx, resourceId, blockNumberPrivateHub, nil)
	if err != nil {
		return fmt.Errorf("creating batches for deposit: %w", err)
	}

	spendKeyResp, err := e.keysClient.GetPaymentSpendKey(ctx, &keyspb.GetPaymentSpendKeyRequest{})
	if err != nil {
		return fmt.Errorf("getting payment spend key for deposit: %w", err)
	}
	spendKey := types.PaymentSpendKey{
		SecretKey: new(big.Int).SetBytes(spendKeyResp.GetSecretKey()),
		PublicKey: new(big.Int).SetBytes(spendKeyResp.GetPublicKey()),
	}

	anonymityIndex := len(batches)
	proof, _, rValues, _, err := e.proofGen.GenerateDepositProof(
		ctx,
		enygma.DepositProofParams{
			ResourceId:        resourceId,
			AnonymityIndex:    anonymityIndex,
			SenderAmount:      amount,
			BlockNumber:       blockNumberPrivateHub,
			Batches:           batches,
			DepositCommitment: depositCommitment,
			DepositSalt:       depositSalt,
			DepositPublicKey:  spendKey.PublicKey,
			TokenAddress:      enygmaAddress,
		},
	)
	if err != nil {
		slog.Error(
			"Error generating Enygma deposit proof",
			slog.String("resourceId", resourceId),
			slog.String("err", err.Error()),
		)
		return fmt.Errorf("generating enygma deposit proof: %w", err)
	}

	senderRFactor := big.NewInt(0)
	for i, batch := range batches {
		batch.ToRValueToAdd = rValues[i]
		if batch.ToChainID.Cmp(e.plChainId) == 0 {
			senderRFactor = rValues[i]
		}
	}

	ctxTx, cancel := context.WithTimeout(ctx, e.conf.DefaultContextTimeout)
	defer cancel()

	enygmaIntegrationAddress, err := e.enygmaClient.GetDvpIntegrationContractAddress(ctx, enygmaAddress)
	if err != nil {
		return fmt.Errorf("getting dvp integration contract address for deposit: %w", err)
	}
	err = e.enygmaIntegrationClient.Deposit(ctxTx, chainEventID, batches, proof, blockNumberPrivateHub, e.plChainId,
		resourceId, amount, from, txHash, enygmaIntegrationAddress)
	if err != nil {
		slog.Error(
			"Error signing Enygma deposit",
			slog.String("resourceId", resourceId),
			slog.Any("err", err.Error()),
		)
		return fmt.Errorf("signing enygma deposit: %w", err)
	}

	history := types.EnygmaHistory{
		ResourceId:            resourceId,
		FromChainId:           e.plChainId,
		RFactor:               senderRFactor,
		BalanceChange:         new(big.Int).Neg(amount),
		BlockNumberPrivateHub: blockNumberPrivateHub,
		EventType:             types.EnygmaDepositToDvp,
	}

	if err := e.enygmaHistoryRepository.InsertEnygmaHistory(ctx, history); err != nil {
		if !errors.Is(err, repository.ErrAlreadyProcessed) {
			return fmt.Errorf("inserting enygma history: %w", err)
		}
	}

	slog.Debug("Enygma deposit done", slog.String("resourceId", resourceId))
	return nil
}

func (e *EnygmaExecutor) ExecuteEnygmaWithdrawal(
	ctx context.Context,
	chainEventID string,
	resourceId string,
	amount *big.Int,
	deposits []*types.DvpDeposit,
	blockNumber uint64,
	enygmaAddress common.Address,
	jsProof *dvp.ProofReceipt,
	from common.Address,
	txHash common.Hash,
) error {
	blockNumberPrivateHub := new(big.Int).SetUint64(blockNumber)

	// Prepare EnygmaDvp related info: deposit private keys, commitments, amounts
	if len(deposits) > e.conf.MaxNumberOfJSDeposits {
		return fmt.Errorf(
			"number of deposits (%d) exceeds MaxNumberOfJSDeposits (%d)",
			len(deposits),
			e.conf.MaxNumberOfJSDeposits,
		)
	}

	depositCommitments := make([]*big.Int, e.conf.MaxNumberOfJSDeposits)
	depositSecretKeys := make([]*big.Int, e.conf.MaxNumberOfJSDeposits)
	depositAmounts := make([]*big.Int, e.conf.MaxNumberOfJSDeposits)
	depositSalts := make([]*big.Int, e.conf.MaxNumberOfJSDeposits)

	for i := range e.conf.MaxNumberOfJSDeposits {
		depositCommitments[i] = big.NewInt(0)
		depositSecretKeys[i] = big.NewInt(0)
		depositAmounts[i] = big.NewInt(0)
		depositSalts[i] = big.NewInt(0)
	}

	spendKeyResp, err := e.keysClient.GetPaymentSpendKey(ctx, &keyspb.GetPaymentSpendKeyRequest{})
	if err != nil {
		return fmt.Errorf("getting payment spend key for withdrawal: %w", err)
	}
	spendKey := types.PaymentSpendKey{
		SecretKey: new(big.Int).SetBytes(spendKeyResp.GetSecretKey()),
		PublicKey: new(big.Int).SetBytes(spendKeyResp.GetPublicKey()),
	}

	for i, deposit := range deposits {
		depositCommitment, err := e.commitmentCalculator.CalculatePaymentCommitment(
			spendKey.PublicKey,
			deposit.Salt,
			deposit.TokenAmount,
			deposit.TokenAddress,
		)
		if err != nil {
			return fmt.Errorf("calculating payment commitment: %w", err)
		}

		depositCommitments[i] = depositCommitment
		depositSecretKeys[i] = spendKey.SecretKey
		depositAmounts[i] = deposit.TokenAmount
		depositSalts[i] = deposit.Salt
	}

	batches, err := e.batcher.CreateBatchesWithAnonimity(ctx, resourceId, blockNumberPrivateHub, nil)
	if err != nil {
		return fmt.Errorf("creating batches for withdrawal: %w", err)
	}

	anonymityIndex := len(batches)
	proof, _, rValues, _, err := e.proofGen.GenerateWithdrawProof(
		ctx,
		enygma.WithdrawProofParams{
			ResourceId:         resourceId,
			AnonymityIndex:     anonymityIndex,
			SenderAmount:       amount,
			BlockNumber:        blockNumberPrivateHub,
			Batches:            batches,
			DepositCommitments: depositCommitments,
			DepositSecretKeys:  depositSecretKeys,
			DepositAmounts:     depositAmounts,
			DepositSalts:       depositSalts,
			TokenAddress:       enygmaAddress,
		},
	)
	if err != nil {
		slog.Error(
			"Error generating Enygma withdrawal proof",
			slog.String("resourceId", resourceId),
			slog.String("err", err.Error()),
		)
		return fmt.Errorf("generating enygma withdrawal proof: %w", err)
	}

	senderRFactor := big.NewInt(0)
	for i, batch := range batches {
		batch.ToRValueToAdd = rValues[i]
		if batch.ToChainID.Cmp(e.plChainId) == 0 {
			senderRFactor = rValues[i]
		}
	}

	ctxTx, cancel := context.WithTimeout(ctx, e.conf.DefaultContextTimeout)
	defer cancel()

	enygmaIntegrationAddress, err := e.enygmaClient.GetDvpIntegrationContractAddress(ctx, enygmaAddress)
	if err != nil {
		return fmt.Errorf("getting dvp integration contract address for withdrawal: %w", err)
	}

	err = e.enygmaIntegrationClient.Withdraw(ctxTx, chainEventID, batches, proof, blockNumberPrivateHub, jsProof, e.plChainId, resourceId, amount, from, txHash, enygmaIntegrationAddress)
	if err != nil {
		slog.Error("Error signing Enygma withdrawal", slog.String("resourceId", resourceId), slog.Any("err", err))
		return fmt.Errorf("signing enygma withdrawal: %w", err)
	}

	history := types.EnygmaHistory{
		ResourceId:            resourceId,
		FromChainId:           e.plChainId,
		RFactor:               senderRFactor,
		BalanceChange:         amount,
		BlockNumberPrivateHub: blockNumberPrivateHub,
		EventType:             types.EnygmaWithdrawFromDvp,
	}

	if err := e.enygmaHistoryRepository.InsertEnygmaHistory(ctx, history); err != nil {
		if !errors.Is(err, repository.ErrAlreadyProcessed) {
			return fmt.Errorf("inserting enygma history: %w", err)
		}
	}

	slog.Debug("Enygma withdrawal done", slog.String("resourceId", resourceId))
	return nil
}
