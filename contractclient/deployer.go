package contractclient

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

type DeployerEndpoint interface {
	ReceivePayload(
		ctx context.Context,
		fromChainId *big.Int,
		from, to common.Address,
		data EndpointV1.RaylsMessage,
		messageId [32]byte,
	) (common.Hash, error)
}

type deployerResourceRegistryClient interface {
	GetResourceById([32]byte) (standard uint8, bytecode []byte, initializerParams []byte, err error)
}

type DeployerClient struct {
	endpoint         DeployerEndpoint
	resourceRegistry deployerResourceRegistryClient
	contextTimeout   time.Duration
}

func NewDeployerClient(
	endpoint DeployerEndpoint,
	resourceRegistry deployerResourceRegistryClient,
	contextTimeout time.Duration,
) *DeployerClient {
	return &DeployerClient{
		endpoint:         endpoint,
		resourceRegistry: resourceRegistry,
		contextTimeout:   contextTimeout,
	}
}

// DeployResourceAndExecute deploys a contract by resource ID and executes the payload.
// It retrieves the resource metadata (bytecode and initializer params) from the resource registry,
// adds it to the message, and deploys the contract via the endpoint's ReceivePayload method.
// Returns the transaction hash of the deployment transaction.
func (c *DeployerClient) DeployResourceAndExecute(
	ctx context.Context,
	resourceId [32]byte,
	message *types.DispatchedMessageToPrivateHub,
) (common.Hash, error) {
	// Get resource metadata from the resource registry
	standard, bytecode, initializerParams, err := c.resourceRegistry.GetResourceById(resourceId)
	if err != nil {
		return common.Hash{}, WrapInDeployerClientError("failed to get resource by id", withstack.Wrap(err))
	}

	// Enrich message with resource metadata for deployment. A non-Custom standard is a
	// FACTORY-mode resource: the registry holds no raw bytecode for it (the receiver deploys
	// from the seeded template registry keyed by the standard), so route through FACTORY and
	// pass the template instead of empty bytecode. ErcStandard and RaylsBridgeableERC share
	// ordinals, so the standard byte is the factoryTemplate. Custom (0) stays the legacy
	// BYTECODE path. Without this, FACTORY tokens deploy with empty bytecode and the factory
	// reverts FactoryV1__EmptyBytecode (0x974918c1).
	meta := &message.Data.MessageMetadata.NewResourceMetadata
	meta.Valid = true
	meta.InitializerParams = initializerParams
	if types.ErcStandard(standard) == types.CUSTOM {
		meta.ResourceDeployType = uint8(types.ResourceDeployTypeBytecode)
		meta.Bytecode = bytecode
		meta.FactoryTemplate = uint8(types.CUSTOM)
	} else {
		meta.ResourceDeployType = uint8(types.ResourceDeployTypeFactory)
		meta.Bytecode = nil
		meta.FactoryTemplate = standard
	}

	ctxTx := ctx
	if c.contextTimeout > 0 {
		var cancel context.CancelFunc
		ctxTx, cancel = context.WithTimeout(ctx, c.contextTimeout)
		defer cancel()
	}

	txHash, err := c.endpoint.ReceivePayload(
		ctxTx,
		message.FromChainId,
		message.From,
		message.To,
		message.Data,
		message.MessageId,
	)
	if err != nil {
		return common.Hash{}, WrapInDeployerClientError(
			"failed to deploy resource via receive payload",
			withstack.Wrap(err),
		)
	}

	return txHash, nil
}
