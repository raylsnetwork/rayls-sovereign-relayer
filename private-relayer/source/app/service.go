// Decommissioning Teleport (vanilla, atomic): atomic message types/handlers noted below; vanilla/generic/Enygma/DVP retained.

package app

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/batcher"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/proofgen"
	sharedservice "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/resultrouter"
)

// MessageType strings for the source-side publishers that target the
// privatehub signing identity. They double as the keys the result
// router uses to fan TxResult batches into per-callback handlers.
const (
	// Decommissioning Teleport (vanilla, atomic): the two atomic.* consts below; msgTypePrivateHubExecute is retained.
	msgTypeAtomicSourceRevert = "atomic.source-revert"
	msgTypeAtomicEarlyRevert  = "atomic.early-revert"
	msgTypePrivateHubExecute  = "privatehub.execute"
)

const (
	routerBatchSize = 100
	routerInterval  = time.Second
)

type ServiceConfig struct {
	PrivateHubChainID      *big.Int
	MyChainID              *big.Int
	CrossChainTickerPeriod time.Duration
	PrivateHubTickerPeriod time.Duration
	EnygmaTickerPeriod     time.Duration

	// Enygma and Dvp configuration
	EnygmaBatchSize                int
	EnygmaMaxConcurrentResourceIDs int
	NumberOfJSParamsIn             int
	DefaultContextTimeout          time.Duration
}

type Services struct {
	crossChain  *service.CrossChainService
	privateHub  *sharedservice.PrivateHubService
	keysService *service.KeysService

	enygmaOrchestrator *service.EnygmaOrchestrator
	dvpOrchestrator    *service.DvpOrchestrator
}

func (r *SourcePrivateRelayer) initializeServices(ctx context.Context, config ServiceConfig) error {
	// Initialize proof generator for cross-chain service
	crossChainProofGen := proofgen.New(r.nodeClient)

	crossChainService := service.NewCrossChainService(
		config.CrossChainTickerPeriod,
		config.MyChainID,
		r.msgqueues.crossChainConsumer,
		r.nodeClient,
		crossChainProofGen,
		r.contractClients.teleport,
		r.repositories.transactionRepository,
		r.repositories.signatureRepository,
	)

	privateHubEndpointAddress, err := r.contractClients.hubRegistry.GetContractAddress("Endpoint")
	if err != nil {
		return fmt.Errorf("failed to get private hub endpoint address: %w", err)
	}

	privateHubBatcher := batcher.NewBatcher(msgTypePrivateHubExecute, r.msgqueues.privateHubSendPublisher)

	privateHubService := sharedservice.NewPrivateHubService(
		config.PrivateHubTickerPeriod,
		config.MyChainID,
		privateHubEndpointAddress,
		r.msgqueues.privateHubConsumer,
		r.txGen.privateHub,
		privateHubBatcher,
	)

	keysService := service.NewKeysService(config.MyChainID, r.ctsClient, r.contractClients.participantStorage)

	// Initialize Enygma services
	//nolint:contextcheck // startup initialization chain
	enygmaOrchestrator, err := r.initializeEnygmaServices(ctx, EnygmaServiceConfig{
		TickerPeriod:                   config.EnygmaTickerPeriod,
		PrivateHubChainID:              config.PrivateHubChainID,
		MyChainID:                      config.MyChainID,
		EnygmaBatchSize:                config.EnygmaBatchSize,
		EnygmaMaxConcurrentResourceIDs: config.EnygmaMaxConcurrentResourceIDs,
		NumberOfJSParamsIn:             config.NumberOfJSParamsIn,
		DefaultContextTimeout:          config.DefaultContextTimeout,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize enygma services: %w", err)
	}

	// Initialize Dvp services
	dvpOrchestrator, err := r.initializeDvpServices(ctx, DvpServiceConfig{
		PrivateHubChainID:     config.PrivateHubChainID,
		MyChainID:             config.MyChainID,
		NumberOfJSParamsIn:    config.NumberOfJSParamsIn,
		DefaultContextTimeout: config.DefaultContextTimeout,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize dvp services: %w", err)
	}

	signatureSvc := r.atomicServices.signature

	r.pnResultRouter.Register(
		msgTypePrivateHubExecute,
		resultrouter.HandlerFunc(privateHubService.HandleResults),
	)
	r.pnResultRouter.Register(
		msgTypeAtomicSourceRevert,
		resultrouter.HandlerFunc(signatureSvc.HandleSourceRevertedCallback),
	)
	r.pnhResultRouter.Register(
		msgTypeAtomicEarlyRevert,
		resultrouter.HandlerFunc(r.atomicServices.earlyRevert.HandleEarlyRevertCallback),
	)

	r.services = &Services{
		crossChain:  crossChainService,
		privateHub:  privateHubService,
		keysService: keysService,

		enygmaOrchestrator: enygmaOrchestrator,
		dvpOrchestrator:    dvpOrchestrator,
	}

	return nil
}
