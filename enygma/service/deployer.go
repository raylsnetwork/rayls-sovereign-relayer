package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/RaylsEnygmaHandler"
	telemetry "github.com/raylsnetwork/rayls-sovereign-relayer/otel"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type deployerTracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

// var _ deployerTracer = (*adapters.OTelTracer)(nil)

type deployerEndpointClient interface {
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

type deployerResourceRegistryClient interface {
	GetResourceById(resourceId [32]byte) (standard uint8, bytecode []byte, initializerParams []byte, err error)
}

type EnygmaDeployer struct {
	endpointClient         deployerEndpointClient
	resourceRegistryClient deployerResourceRegistryClient
	tracer                 deployerTracer
}

func NewEnygmaDeployer(
	endpointClient deployerEndpointClient,
	resourceRegistryClient deployerResourceRegistryClient,
	tracer deployerTracer,
) *EnygmaDeployer {
	return &EnygmaDeployer{
		endpointClient:         endpointClient,
		resourceRegistryClient: resourceRegistryClient,
		tracer:                 tracer,
	}
}

func (d *EnygmaDeployer) Deploy(
	ctx context.Context,
	resourceId [32]byte,
	initiatorChainId *big.Int,
) (common.Address, error) {
	ctx, span := d.tracer.Start(ctx, telemetry.SPAN_DEPLOY_ENYGMA_CONTRACT)
	defer span.End()

	enygmaPLABI, err := RaylsEnygmaHandler.RaylsEnygmaHandlerMetaData.ParseABI()
	if err != nil {
		return common.Address{}, fmt.Errorf("getting enygma handler ABI: %w", err)
	}

	payload, err := enygmaPLABI.Pack("crossTransferCheck")
	if err != nil {
		return common.Address{}, fmt.Errorf("packing crossTransferCheck payload: %w", err)
	}

	// Unique msg ID per token deployment.
	msgId := crypto.Keccak256Hash(payload, resourceId[:])
	message := types.DispatchedMessageToPrivateHub{}
	message.FromChainId = initiatorChainId
	message.MessageId = msgId
	message.Data.Payload = payload
	message.Data.MessageMetadata.Valid = true
	message.Data.MessageMetadata.Nonce = big.NewInt(0)
	message.Data.MessageMetadata.IgnoresNonce = true
	message.Data.MessageMetadata.ResourceId = resourceId
	message.Data.MessageMetadata.TransferMetadata.Id = big.NewInt(0)
	message.Data.MessageMetadata.TransferMetadata.Amount = big.NewInt(0)

	standard, _, initParams, err := d.resourceRegistryClient.GetResourceById(resourceId)
	if err != nil {
		return common.Address{}, fmt.Errorf("getting resource by ID %s: %w", hex.EncodeToString(resourceId[:]), err)
	}

	// An Enygma resource is always deployed via the ENYGMA factory template (see below), so a
	// registry standard other than ENYGMA signals a misconfiguration. We don't fail — the deploy
	// proceeds as ENYGMA regardless — but we surface it for observability.
	if standard != uint8(types.ENYGMA) {
		slog.Warn("Enygma resource has unexpected registry standard; deploying as ENYGMA anyway",
			slog.String("resourceId", hex.EncodeToString(resourceId[:])),
			slog.Uint64("registryStandard", uint64(standard)),
			slog.Uint64("expectedStandard", uint64(types.ENYGMA)))
	}

	// Deploy the Enygma handler through the PN factory (FACTORY mode) rather than
	// from issuer-supplied bytecode (BYTECODE mode). Factory deploys via CREATE2 + InitCodeStub
	// produce identical runtime bytecode on every PN. An Enygma resource is always deployed via
	// the ENYGMA factory template regardless of the standard recorded in the registry (the
	// registry standard is only checked above for observability). The issuer's `bytecode` field
	// is discarded; only initializerParams are forwarded.

	// add resource in the message metadata
	message.Data.MessageMetadata.NewResourceMetadata.Valid = true
	message.Data.MessageMetadata.NewResourceMetadata.ResourceDeployType = uint8(types.ResourceDeployTypeFactory)
	message.Data.MessageMetadata.NewResourceMetadata.FactoryTemplate = uint8(types.ENYGMA)
	message.Data.MessageMetadata.NewResourceMetadata.InitializerParams = initParams

	span.SetAttributes(
		attribute.String(telemetry.ATTR_RESOURCE_ID, hex.EncodeToString(resourceId[:])),
		attribute.String(telemetry.ATTR_SENDER_CHAIN_ID, initiatorChainId.String()),
	)

	_, err = d.endpointClient.ReceivePayload(
		ctx,
		message.FromChainId,
		message.From,
		message.To,
		message.Data,
		message.MessageId,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, telemetry.STATUS_ERROR_FAILED_TO_DEPLOY_ENYGMA_CONTRACT)
		return common.Address{}, fmt.Errorf("receiving enygma deploy payload for resource %s: %w", hex.EncodeToString(resourceId[:]), err)
	}

	deployedContractAddress, err := d.endpointClient.GetResourceAddress(ctx, hex.EncodeToString(resourceId[:]))
	if err != nil {
		return common.Address{}, fmt.Errorf(
			"error getting contract address from resource ID %s",
			hex.EncodeToString(resourceId[:]),
		)
	}

	return deployedContractAddress, nil
}
