package contractclient

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/conv"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

type EndpointContract interface {
	GetAddressByResourceId(opts *bind.CallOpts, _resourceId [32]byte) (common.Address, error)
	ReceivePayload(
		opts *bind.TransactOpts,
		fromChainId *big.Int,
		from common.Address,
		to common.Address,
		data EndpointV1.RaylsMessage,
		messageId [32]byte,
	) (*ethTypes.Transaction, error)
	PackReceivePayload(srcChainId *big.Int, srcAddress common.Address, dstAddress common.Address, raylsMessage EndpointV1.RaylsMessage, messageId [32]byte) []byte
}

type EndpointClient struct {
	address  common.Address
	contract *EndpointV1.EndpointV1
	executor Executor
}

func NewEndpointClient(address common.Address, executor Executor) *EndpointClient {
	return &EndpointClient{
		contract: EndpointV1.NewEndpointV1(),
		address:  address,
		executor: executor,
	}
}

func (c *EndpointClient) GetResourceAddress(ctx context.Context, resourceId string) (common.Address, error) {
	resourceIdBytes, err := conv.StringToBytes32(resourceId)
	if err != nil {
		return common.Address{}, WrapInEndpointClientError("failed to convert resource id to bytes", err)
	}

	calldata := c.contract.PackGetAddressByResourceId(resourceIdBytes)

	raw, err := c.executor.Call(ctx, c.address, calldata)
	if err != nil {
		return common.Address{}, WrapInEndpointClientError("failed to call GetAddressByResourceId", withstack.Wrap(err))
	}

	addr, err := c.contract.UnpackGetAddressByResourceId(raw)
	if err != nil {
		return common.Address{}, WrapInEndpointClientError("failed to unpack GetAddressByResourceId result", withstack.Wrap(err))
	}
	return addr, nil
}

func (c *EndpointClient) ReceivePayload(
	ctx context.Context,
	fromChainId *big.Int,
	from common.Address,
	to common.Address,
	data EndpointV1.RaylsMessage,
	messageId [32]byte,
) (common.Hash, error) {

	calldata := c.contract.PackReceivePayload(fromChainId, from, to, data, messageId)
	receipt, err := c.executor.Execute(ctx, IDFor("endpoint.ReceivePayload", common.Hash(messageId).Hex()), calldata, c.address)

	if err != nil {
		return common.Hash{}, WrapInEndpointClientError("failed to get receive payload receipt", err)
	}

	return receipt.TxHash, nil
}
