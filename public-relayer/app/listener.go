// Decommissioning Teleport (vanilla, atomic).

package app

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/RNMessageDispatcherV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/PNTokenCoreV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/listener"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/public-relayer/handler"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

type ListenersConfig struct {
	BatchSize int

	PublicChainStartingBlock  *big.Int
	PrivateChainStartingBlock *big.Int
}

type Listeners struct {
	public  *listener.LogListener
	private *listener.LogListener
}

func (n *PublicRelayer) initializeListeners(ctx context.Context, conf ListenersConfig) error {
	// Get addresses from contracts struct
	// Get endpoint addresses from contracts struct
	publicDispatcherAddress, err := n.contractClients.publicDeploymentRegistry.GetContractAddress("RNMessageDispatcher")
	if err != nil {
		return fmt.Errorf("failed to get public message dispatcher address: %w", err)
	}
	privateDispatcherAddress, err := n.contractClients.privateDeploymentRegistry.GetContractAddress("RNMessageDispatcher")
	if err != nil {
		return fmt.Errorf("failed to get private message dispatcher address: %w", err)
	}
	tokenCoreAddress, err := n.contractClients.privateDeploymentRegistry.GetContractAddress("TokenCore")
	if err != nil {
		return fmt.Errorf("failed to get private token core address: %w", err)
	}

	publicHandler := handler.NewPublicRelayerPublicHandler(
		n.messageQueues.publicPublisher,
		RNMessageDispatcherV1.NewRNMessageDispatcherV1(),
	)

	publicListenerConfig := listener.LogListenerConfig{
		Component: types.LastProcessedBlockDocumentPublicChain,

		StartingBlock: conf.PublicChainStartingBlock,
		BatchSize:     conf.BatchSize,

		Addresses: []common.Address{publicDispatcherAddress},
	}
	publicListener, err := listener.NewLogListener(
		ctx,
		publicListenerConfig,
		listener.LogHandlerFunc(publicHandler.Handle),
		n.publicClient,
		n.repositories.lastProcessedBlockRepository)
	if err != nil {
		return fmt.Errorf("failed to initialize public listener: %w", err)
	}

	privateHandler := handler.NewPublicRelayerPrivateHandler(
		n.messageQueues.privatePublisher,
		RNMessageDispatcherV1.NewRNMessageDispatcherV1(),
		n.messageQueues.deploymentPublisher,
		PNTokenCoreV1.NewPNTokenCoreV1(),
		n.contractClients.tokenGovernance,
	)

	privateListenerConfig := listener.LogListenerConfig{
		Component: types.LastProcessedBlockDocumentPrivateChain,

		StartingBlock: conf.PrivateChainStartingBlock,
		BatchSize:     conf.BatchSize,

		Addresses: []common.Address{
			tokenCoreAddress,
			privateDispatcherAddress,
		},
	}
	privateListener, err := listener.NewLogListener(
		ctx,
		privateListenerConfig,
		listener.LogHandlerFunc(privateHandler.Handle),
		n.privateClient,
		n.repositories.lastProcessedBlockRepository)
	if err != nil {
		return fmt.Errorf("failed to initialize private listener: %w", err)
	}

	n.listeners = &Listeners{
		public:  publicListener,
		private: privateListener,
	}
	return nil
}
