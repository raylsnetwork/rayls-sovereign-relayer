// Decommissioning Teleport (vanilla, atomic).

package app

import (
	"fmt"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/batcher"
	"github.com/raylsnetwork/rayls-sovereign-relayer/public-relayer/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/public-relayer/txgen"
	"github.com/raylsnetwork/rayls-sovereign-relayer/resultrouter"
)

type ServiceConfig struct {
	// LPS service configuration
	RelayInterval time.Duration
}

// MessageType strings — kept here as constants because they are also the
// keys both batchers (as the publishedmessageType) and result routers (as
// the registered handler key) use to round-trip a request through CTS.
const (
	msgTypePublicForward  = "publicrelayer.forward.public"
	msgTypePublicRevert   = "publicrelayer.revert.public"
	msgTypePrivateForward = "publicrelayer.forward.private"
	msgTypePrivateRevert  = "publicrelayer.revert.private"
)

// Result-router tunables. Both routers fetch on the same cadence as the
// CTS receipter, which is the only producer feeding cts.result.<identity>.
const (
	routerBatchSize = 100
	routerInterval  = time.Second
)

type Services struct {
	publicGenerator  *service.GeneratorService
	privateGenerator *service.GeneratorService

	deployer *service.DeployerService

	// Two routers, one per CTS signing identity. Each fans the inbound
	// TxResult batches into the right HandleForwardResults /
	// HandleRevertResults callback by MessageType.
	publicChainResultRouter  *resultrouter.Router
	privateChainResultRouter *resultrouter.Router
}

func (p *PublicRelayer) initializeServices(config ServiceConfig) error {
	// Get endpoint addresses from contracts struct
	publicEndpointAddr, err := p.contractClients.publicDeploymentRegistry.GetContractAddress("PublicRNEndpoint")
	if err != nil {
		return fmt.Errorf("failed to get public endpoint address: %w", err)
	}
	privateEndpointAddr, err := p.contractClients.privateDeploymentRegistry.GetContractAddress("RNEndpoint")
	if err != nil {
		return fmt.Errorf("failed to get public endpoint address: %w", err)
	}

	publicTxGenerator, err := txgen.NewEIP5164Generators(p.publicClient)
	if err != nil {
		return fmt.Errorf("failed to initialize public EIP5164 transaction generator: %w", err)
	}

	// publicGenerator: PN → PC (forward to private chain), revert on public chain.
	publicForwardBatcher := batcher.NewBatcher(msgTypePublicForward, p.messageQueues.privateChainSendPublisher)
	publicRevertBatcher := batcher.NewBatcher(msgTypePublicRevert, p.messageQueues.publicChainSendPublisher)
	publicGeneratorConfig := service.GeneratorServiceConfig{
		Interval:              config.RelayInterval,
		EndpointAddress:       privateEndpointAddr,
		SourceEndpointAddress: publicEndpointAddr,
	}
	publicGenerator := service.NewGeneratorService(
		publicGeneratorConfig,
		p.messageQueues.publicConsumer,
		publicTxGenerator,
		publicForwardBatcher,
		publicRevertBatcher,
		p.repositories.publicRevertSignature,
		p.repositories.publicMessageRecord,
	)

	privateTxGenerator, err := txgen.NewEIP5164Generators(p.privateClient)
	if err != nil {
		return fmt.Errorf("failed to initialize private EIP5164 transaction generator: %w", err)
	}

	// privateGenerator: PC → PN (forward to public chain), revert on private chain.
	privateForwardBatcher := batcher.NewBatcher(msgTypePrivateForward, p.messageQueues.publicChainSendPublisher)
	privateRevertBatcher := batcher.NewBatcher(msgTypePrivateRevert, p.messageQueues.privateChainSendPublisher)
	privateGeneratorConfig := service.GeneratorServiceConfig{
		Interval:              config.RelayInterval,
		EndpointAddress:       publicEndpointAddr,
		SourceEndpointAddress: privateEndpointAddr,
	}
	privateGenerator := service.NewGeneratorService(
		privateGeneratorConfig,
		p.messageQueues.privateConsumer,
		privateTxGenerator,
		privateForwardBatcher,
		privateRevertBatcher,
		p.repositories.privateRevertSignature,
		p.repositories.privateMessageRecord,
	)

	deployer := service.NewDeployerService(
		p.messageQueues.deploymentConsumer,
		publicEndpointAddr,
		p.contractClients.deployer,
		p.contractClients.tokenGovernance,
		p.contractClients.publicAccessManager,
	)

	// publicchain identity carries publicGenerator's revert + privateGenerator's forward.
	publicChainResultRouter := resultrouter.New(
		resultrouter.Config{
			Identity:  "publicchain",
			BatchSize: routerBatchSize,
			Interval:  routerInterval,
		},
		p.messageQueues.publicChainResultConsumer,
	)
	publicChainResultRouter.Register(msgTypePublicRevert, resultrouter.HandlerFunc(publicGenerator.HandleRevertResults))
	publicChainResultRouter.Register(msgTypePrivateForward, resultrouter.HandlerFunc(privateGenerator.HandleForwardResults))

	// privatechain identity carries publicGenerator's forward + privateGenerator's revert.
	privateChainResultRouter := resultrouter.New(
		resultrouter.Config{
			Identity:  "privatechain",
			BatchSize: routerBatchSize,
			Interval:  routerInterval,
		},
		p.messageQueues.privateChainResultConsumer,
	)
	privateChainResultRouter.Register(msgTypePublicForward, resultrouter.HandlerFunc(publicGenerator.HandleForwardResults))
	privateChainResultRouter.Register(msgTypePrivateRevert, resultrouter.HandlerFunc(privateGenerator.HandleRevertResults))

	p.services = &Services{
		publicGenerator:  publicGenerator,
		privateGenerator: privateGenerator,

		deployer: deployer,

		publicChainResultRouter:  publicChainResultRouter,
		privateChainResultRouter: privateChainResultRouter,
	}

	return nil
}
