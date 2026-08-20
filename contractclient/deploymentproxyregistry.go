package contractclient

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/DeploymentProxyRegistryV1"
)

type DeploymentProxyRegistryClient struct {
	addresses map[string]common.Address
}

func NewDeploymentProxyRegistryClient(
	address common.Address,
	backend bind.ContractBackend,
) (*DeploymentProxyRegistryClient, error) {
	contract := DeploymentProxyRegistryV1.NewDeploymentProxyRegistryV1()

	calldata := contract.PackGetAllContracts()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	raw, err := backend.CallContract(ctx, ethereum.CallMsg{
		To:   &address,
		Data: calldata,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get contracts from DeploymentProxy: %w", err)
	}

	namesAndAddresses, err := contract.UnpackGetAllContracts(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack contracts from DeploymentProxy: %w", err)
	}

	addresses := make(map[string]common.Address)
	for i, nameAndAddress := range namesAndAddresses.Names {
		if i < len(namesAndAddresses.Addresses) {
			addresses[nameAndAddress] = namesAndAddresses.Addresses[i]
		}
	}
	return &DeploymentProxyRegistryClient{
		addresses: addresses,
	}, nil
}

func (d *DeploymentProxyRegistryClient) GetContractAddress(name string) (common.Address, error) {
	address, ok := d.addresses[name]
	if !ok {
		return common.Address{}, fmt.Errorf("contract '%s' not found in deployment registry", name)
	}
	if address == (common.Address{}) {
		return common.Address{}, fmt.Errorf("contract '%s' has zero address in deployment registry", name)
	}
	return address, nil
}
