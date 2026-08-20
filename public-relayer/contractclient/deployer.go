// Decommissioning Teleport (vanilla, atomic).

package contractclient

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/PublicChainERC1155"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/PublicChainERC20"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/PublicChainERC721"
)

type DeployExecutor interface {
	Deploy(ctx context.Context, bytecode []byte, constructor []byte) (common.Address, *types.Receipt, error)
}

type DeployerClient struct {
	deployer DeployExecutor

	client bind.ContractBackend

	erc20Contract   *PublicChainERC20.PublicChainERC20
	erc721Contract  *PublicChainERC721.PublicChainERC721
	erc1155Contract *PublicChainERC1155.PublicChainERC1155
}

func NewDeployerClient(
	deployer DeployExecutor,
) *DeployerClient {
	return &DeployerClient{
		deployer: deployer,

		erc20Contract:   PublicChainERC20.NewPublicChainERC20(),
		erc721Contract:  PublicChainERC721.NewPublicChainERC721(),
		erc1155Contract: PublicChainERC1155.NewPublicChainERC1155(),
	}
}

func (d *DeployerClient) DeployPublicChainERC20(
	ctx context.Context,
	name, symbol string,
	privateAddress, raylsNodeEndpoint common.Address,
) (common.Address, error) {
	constructor := d.erc20Contract.PackConstructor(
		name,
		symbol,
		raylsNodeEndpoint,
		big.NewInt(0),
		privateAddress,
	)
	byteCode := PublicChainERC20.PublicChainERC20MetaData.Bin

	address, receipt, err := d.deployer.Deploy(ctx, common.FromHex(byteCode), constructor)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to execute transaction: %w", err)
	}
	if receipt.Status == 0 {
		return common.Address{}, fmt.Errorf("transaction failed with receipt status: %d", receipt.Status)
	}

	return address, nil
}

func (d *DeployerClient) DeployPublicChainERC721(
	ctx context.Context,
	uri, name, symbol string,
	privateAddress, raylsNodeEndpoint common.Address,
) (common.Address, error) {
	constructor := d.erc721Contract.PackConstructor(
		uri,
		name,
		symbol,
		raylsNodeEndpoint,
		privateAddress,
	)
	byteCode := PublicChainERC721.PublicChainERC721MetaData.Bin

	address, receipt, err := d.deployer.Deploy(ctx, common.FromHex(byteCode), constructor)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to execute transaction: %w", err)
	}
	if receipt.Status == 0 {
		return common.Address{}, fmt.Errorf("transaction failed with receipt status: %d", receipt.Status)
	}

	return address, nil
}

func (d *DeployerClient) DeployPublicChainERC1155(
	ctx context.Context,
	uri, name string,
	privateAddress, raylsNodeEndpoint common.Address,
) (common.Address, error) {
	constructor := d.erc1155Contract.PackConstructor(
		uri,
		name,
		raylsNodeEndpoint,
		privateAddress,
	)
	byteCode := PublicChainERC1155.PublicChainERC1155MetaData.Bin

	address, receipt, err := d.deployer.Deploy(ctx, common.FromHex(byteCode), constructor)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to execute transaction: %w", err)
	}
	if receipt.Status == 0 {
		return common.Address{}, fmt.Errorf("transaction failed with receipt status: %d", receipt.Status)
	}

	return address, nil
}
