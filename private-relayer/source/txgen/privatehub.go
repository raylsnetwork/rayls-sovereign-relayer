package txgen

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

//go:generate moq --pkg txgen_test -out ethereum_client_mock_test.go . EthereumClient

type EthereumClient interface {
	ChainID(context.Context) (*big.Int, error)
	SuggestGasPrice(context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
}

type PrivateHubGenerator struct {
	client EthereumClient

	endpointABI *abi.ABI
}

func NewPrivateHubGenerator(client EthereumClient) (*PrivateHubGenerator, error) {
	endpointABI, err := EndpointV1.EndpointV1MetaData.ParseABI()
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("error parsing endpoint ABI: %w", err))
	}

	return &PrivateHubGenerator{
		client: client,

		endpointABI: endpointABI,
	}, nil
}

func (g *PrivateHubGenerator) Generate(
	fromChainID *big.Int,
	fromAddress, toAddress common.Address,
	message EndpointV1.RaylsMessage,
	id common.Hash,
) ([]byte, error) {
	data, err := g.endpointABI.Pack("receivePayload", fromChainID, fromAddress, toAddress, message, id)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to pack function parameters: %w", err))
	}

	return data, nil
}
