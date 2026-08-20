// Decommissioning Teleport (vanilla, atomic).

// Package txgen implements the legacy public-chain (RN) Teleport bridge relayer.
//
// Deprecated: Decommissioning Teleport (vanilla, atomic).
package txgen

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/PublicRNEndpointV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/RNMessageDispatcherV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

//go:generate moq --pkg txgen_test -out ethereum_client_mock_test.go . ERC1564EthereumClient

type ERC1564EthereumClient interface {
	ChainID(context.Context) (*big.Int, error)
	SuggestGasPrice(context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
}

type ERC1564Generator struct {
	client ERC1564EthereumClient

	endpointABI *abi.ABI
	chainID     *big.Int
}

func NewEIP5164Generators(client ERC1564EthereumClient) (*ERC1564Generator, error) {
	endpointABI, err := PublicRNEndpointV1.PublicRNEndpointV1MetaData.ParseABI()
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("error parsing endpoint ABI: %w", err))
	}

	chainID, err := client.ChainID(context.TODO())
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("error getting chain ID: %w", err))
	}

	return &ERC1564Generator{
		client: client,

		endpointABI: endpointABI,
		chainID:     chainID,
	}, nil
}

func (g *ERC1564Generator) Generate(
	fromAddress, toAddress common.Address,
	message RNMessageDispatcherV1.RaylsNodeMessage,
	id common.Hash,
) ([]byte, error) {
	data, err := g.endpointABI.Pack("receivePayload", g.chainID, fromAddress, toAddress, message, id)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to pack function parameters: %w", err))
	}

	return data, nil
}
