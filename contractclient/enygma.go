package contractclient

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EnygmaV1"
	telemetry "github.com/raylsnetwork/rayls-privacy-relayer-api/otel"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type enygmaEncryptor interface {
	EncryptEnygmaTransferBatches(ctx context.Context, batches []*types.EnygmaTransferBatch, blockNumber *big.Int) ([][]byte, error)
}

type EnygmaClient struct {
	contract *EnygmaV1.EnygmaV1
	encr     enygmaEncryptor
	executor Executor
	tracer   trace.Tracer
}

func NewEnygmaClient(executor Executor, encr enygmaEncryptor) *EnygmaClient {
	return &EnygmaClient{
		contract: EnygmaV1.NewEnygmaV1(),
		encr:     encr,
		executor: executor,
		tracer:   otel.Tracer("enygma-client"),
	}
}

func (c *EnygmaClient) GetDvpIntegrationContractAddress(ctx context.Context, tokenAddress common.Address) (common.Address, error) {
	calldata := c.contract.PackGetDvpIntegrationContractAddress()

	raw, err := c.executor.Call(ctx, tokenAddress, calldata)
	if err != nil {
		return common.Address{}, WrapInEnygmaClientError("failed to get dvp integration contract address", withstack.Wrap(err))
	}

	address, err := c.contract.UnpackGetDvpIntegrationContractAddress(raw)
	if err != nil {
		return common.Address{}, WrapInEnygmaClientError("failed to unpack dvp integration contract address", withstack.Wrap(err))
	}
	return address, nil
}

func (c *EnygmaClient) GetPublicValuesFinalised(ctx context.Context, tokenAddress common.Address) (*types.EnygmaPublicValues, error) {
	_, span := c.tracer.Start(ctx, telemetry.SPAN_GET_PUBLIC_VALUES_FINALISED)
	defer span.End()

	calldata := c.contract.PackGetPublicValuesFinalised()

	raw, err := c.executor.Call(ctx, tokenAddress, calldata)
	if err != nil {
		return nil, WrapInEnygmaClientError("failed to get public values finalised", withstack.Wrap(err))
	}

	result, err := c.contract.UnpackGetPublicValuesFinalised(raw)
	if err != nil {
		return nil, WrapInEnygmaClientError("failed to unpack public values finalised", withstack.Wrap(err))
	}

	span.SetStatus(codes.Ok, "")
	return &types.EnygmaPublicValues{
		Commitments: convertCommitmentsByChainId(result.Arg0),
		PublicKeys:  convertPublicKeysByChainId(result.Arg1),
	}, nil
}

// SignSupplyUpdate signs a supply update transaction without broadcasting it.
// The caller is responsible for broadcasting and waiting for the receipt.
func (c *EnygmaClient) SupplyUpdate(
	ctx context.Context,
	chainEventID string,
	tokenAddress common.Address,
	senderChainId *big.Int,
	blockNumber *big.Int,
	update types.EnygmaSupplyUpdate,
) error {
	_, span := c.tracer.Start(ctx, telemetry.SPAN_UPDATE_ENYGMA_SUPPLY)
	defer span.End()

	supplyUpdate := EnygmaV1.IEnygmaV1SupplyUpdateTx{
		Amount: update.Amount,
		TxType: uint8(update.Type),
	}

	calldata := c.contract.PackUpdateSupply(senderChainId, blockNumber, supplyUpdate)

	_, err := c.executor.Execute(ctx, IDFor("enygmaclient.SupplyUpdate", chainEventID), calldata, tokenAddress)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, telemetry.STATUS_ERROR_FAILED_TO_SEND_SUPPLY_UPDATE)
		return WrapInEnygmaClientError("failed to execute enygma supply update", withstack.Wrap(err))
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// SignTransferBatch signs a transfer batch transaction without broadcasting it.
// The caller is responsible for broadcasting and waiting for the receipt.
func (c *EnygmaClient) TransferBatch(
	ctx context.Context,
	batchID string,
	tokenAddress common.Address,
	batches []*types.EnygmaTransferBatch,
	proof *types.EnygmaProofResponse,
	blockNumber *big.Int,
) error {
	ctx, span := c.tracer.Start(ctx, telemetry.SPAN_SEND_ENYGMA_TRANSFER_BATCH)
	defer span.End()

	encryptCtx, encryptSpan := c.tracer.Start(ctx, telemetry.SPAN_ENCRYPT_ENYGMA_TRANSFER_BATCH)
	encrBatches, err := c.encr.EncryptEnygmaTransferBatches(encryptCtx, batches, blockNumber)
	encryptSpan.End()
	if err != nil {
		return WrapInEnygmaClientError("failed to encrypt enygma transfer batches", err)
	}

	contractProof := EnygmaV1.IEnygmaV1TransferProof{
		PiA:          [2]*big.Int{proof.PiA[0], proof.PiA[1]},
		PiB:          [2][2]*big.Int{{proof.PiB[0][0], proof.PiB[0][1]}, {proof.PiB[1][0], proof.PiB[1][1]}},
		PiC:          [2]*big.Int{proof.PiC[0], proof.PiC[1]},
		PublicSignal: proof.PublicSignal,
	}

	calldata := c.contract.PackTransferBatch(contractProof, encrBatches)

	_, transferBatchSpan := c.tracer.Start(ctx, telemetry.SPAN_TRANSFER_BATCH)
	_, err = c.executor.Execute(ctx, IDFor("enygmaclient.TransferBatch", batchID), calldata, tokenAddress)
	transferBatchSpan.End()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, telemetry.STATUS_ERROR_FAILED_TO_EXECUTE_TRANSFER_BATCH)
		return WrapInEnygmaClientError("failed to sign enygma transfer batch", withstack.Wrap(err))
	}
	transferBatchSpan.SetStatus(codes.Ok, "")

	span.SetStatus(codes.Ok, telemetry.STATUS_SUCCESS_TRANSFER_BATCH_SENT)
	return nil
}

func convertCommitmentsByChainId(values []EnygmaV1.IEnygmaV1EnygmaPointWithChainId) map[string]*types.Point {
	result := make(map[string]*types.Point)
	for _, value := range values {
		result[value.ChainId.String()] = &types.Point{
			X: value.C1,
			Y: value.C2,
		}
	}
	return result
}

func convertPublicKeysByChainId(values []EnygmaV1.IEnygmaV1EnygmaPublicKeyWithChainId) map[string]*big.Int {
	result := make(map[string]*big.Int)
	for _, value := range values {
		result[value.ChainId.String()] = value.PublicKey
	}
	return result
}
