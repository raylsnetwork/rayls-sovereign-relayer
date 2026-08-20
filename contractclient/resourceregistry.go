package contractclient

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/ResourceRegistryV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

// ResourceRegistryBackend is the minimal interface needed for calling
// the ResourceRegistry contract. It is satisfied by *ethclient.Client.
type ResourceRegistryBackend interface {
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
}

type ResourceRegistryClient struct {
	address  common.Address
	backend  ResourceRegistryBackend
	contract *ResourceRegistryV1.ResourceRegistryV1
}

func NewResourceRegistryClient(address common.Address, backend ResourceRegistryBackend) *ResourceRegistryClient {
	return &ResourceRegistryClient{
		address:  address,
		backend:  backend,
		contract: ResourceRegistryV1.NewResourceRegistryV1(),
	}
}

func (c *ResourceRegistryClient) GetResourceById(resourceId [32]byte) (uint8, []byte, []byte, error) {
	calldata := c.contract.PackGetResourceById(resourceId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	raw, err := c.backend.CallContract(ctx, ethereum.CallMsg{
		To:   &c.address,
		Data: calldata,
	}, nil)
	if err != nil {
		return 0, nil, nil, WrapInResourceRegistryClientError("failed to get resource by id", withstack.Wrap(err))
	}

	resource, err := c.contract.UnpackGetResourceById(raw)
	if err != nil {
		return 0, nil, nil, WrapInResourceRegistryClientError("failed to unpack resource by id", withstack.Wrap(err))
	}

	return resource.Standard, resource.Bytecode, resource.InitializerParams, nil
}
