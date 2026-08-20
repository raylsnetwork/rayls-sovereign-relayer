// Decommissioning Teleport (vanilla, atomic).

package app

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	sharedcontractclient "github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/public-relayer/contractclient"
)

type ContractClientsConfig struct {
	PublicDeploymentProxyRegistryAddress  common.Address
	PrivateDeploymentProxyRegistryAddress common.Address
}

type ContractClients struct {
	deployer            *contractclient.DeployerClient
	tokenGovernance     *contractclient.TokenGovernanceClient
	publicAccessManager *sharedcontractclient.AccessManager

	publicDeploymentRegistry  *sharedcontractclient.DeploymentProxyRegistryClient
	privateDeploymentRegistry *sharedcontractclient.DeploymentProxyRegistryClient
}

func (p *PublicRelayer) initializeContractClients(conf ContractClientsConfig) error {
	// Initialize Public Chain DeploymentProxyRegistry
	publicDeploymentRegistry, err := sharedcontractclient.NewDeploymentProxyRegistryClient(conf.PublicDeploymentProxyRegistryAddress, p.publicClient)
	if err != nil {
		return fmt.Errorf("failed to create public deployment proxy: %w", err)
	}

	// Initialize Private Chain DeploymentProxyRegistry
	privateDeploymentRegistry, err := sharedcontractclient.NewDeploymentProxyRegistryClient(conf.PrivateDeploymentProxyRegistryAddress, p.privateClient)
	if err != nil {
		return fmt.Errorf("failed to create private deployment proxy: %w", err)
	}

	publicExecutor := sharedcontractclient.NewCTSExecutor(sharedcontractclient.NewDefaultRetryingTxOpsClient(p.ctsClients.PublicChainTxOpsServiceClient))
	privateExecutor := sharedcontractclient.NewCTSExecutor(sharedcontractclient.NewDefaultRetryingTxOpsClient(p.ctsClients.PrivateChainTxOpsServiceClient))

	publicDeployer := sharedcontractclient.NewCTSDeployer(p.ctsClients.PublicChainTxOpsServiceClient)
	deployer := contractclient.NewDeployerClient(publicDeployer)

	publicAccessManagerAddress, err := publicDeploymentRegistry.GetContractAddress("RaylsAccessManager")
	if err != nil {
		return fmt.Errorf("failed to get access manager address: %w", err)
	}
	publicAccessManager := sharedcontractclient.NewAccessManager(
		publicAccessManagerAddress,
		publicExecutor,
	)

	tokenGovernanceAddress, err := privateDeploymentRegistry.GetContractAddress("TokenRegistry")
	if err != nil {
		return fmt.Errorf("failed to get token registry address: %w", err)
	}
	tokenGovernance := contractclient.NewTokenGovernanceClient(
		tokenGovernanceAddress,
		privateExecutor,
	)

	p.contractClients = &ContractClients{
		deployer:            deployer,
		tokenGovernance:     tokenGovernance,
		publicAccessManager: publicAccessManager,

		publicDeploymentRegistry:  publicDeploymentRegistry,
		privateDeploymentRegistry: privateDeploymentRegistry,
	}
	return nil
}
