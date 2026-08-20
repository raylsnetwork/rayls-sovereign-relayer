package txgen

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

//go:generate moq --pkg txgen_test -out ethereum_client_mock_test.go . EthereumClient

type EthereumClient interface {
	ChainID(context.Context) (*big.Int, error)
	SuggestGasPrice(context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
}

type PrivateNodeGenerator struct {
	client EthereumClient

	endpointABI *abi.ABI
}

func NewPrivateNodeGenerator(client EthereumClient) (*PrivateNodeGenerator, error) {
	endpointABI, err := EndpointV1.EndpointV1MetaData.ParseABI()
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("error parsing endpoint ABI: %w", err))
	}

	return &PrivateNodeGenerator{
		client:      client,
		endpointABI: endpointABI,
	}, nil
}

func (g *PrivateNodeGenerator) Generate(
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
