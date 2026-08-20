package app

import (
	"context"
	"math/big"
	"time"

	dvpService "github.com/raylsnetwork/rayls-sovereign-relayer/dvp/service"
	enygmaAdapters "github.com/raylsnetwork/rayls-sovereign-relayer/enygma/adapters"
	enygmaHandler "github.com/raylsnetwork/rayls-sovereign-relayer/enygma/handler"
	enygmaService "github.com/raylsnetwork/rayls-sovereign-relayer/enygma/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/service"
)

const (
	enygmaDepositWaiterMaxRetries    = 30
	enygmaDepositWaiterRetryInterval = 3 * time.Second
)

// EnygmaServiceConfig holds configuration for enygma services
type EnygmaServiceConfig struct {
	TickerPeriod                   time.Duration
	PrivateHubChainID              *big.Int
	MyChainID                      *big.Int
	EnygmaBatchSize                int
	EnygmaMaxConcurrentResourceIDs int
	NumberOfJSParamsIn             int
	DefaultContextTimeout          time.Duration
}

// initializeEnygmaServices initializes all enygma-related services
func (r *SourcePrivateRelayer) initializeEnygmaServices(
	ctx context.Context,
	config EnygmaServiceConfig,
) (*service.EnygmaOrchestrator, error) {
	_ = ctx // reserved for Phase 2 idempotency follow-up (txRecoveryRepository, persistAndBroadcast wiring)
	tracer := enygmaAdapters.NewOTelTracer("enygma")

	// Initialize Enygma BlockWaiterService
	blockWaiter := enygmaService.NewBlockWaiterService(
		enygmaAdapters.NewOTelTracer("enygma-block-waiter"),
		r.hubClient, // Use private hub for block waiting
	)

	// Initialize Retry Service
	retryService := enygmaService.NewRetryService(
		enygmaAdapters.NewOTelTracer("enygma-retry"),
		blockWaiter,
	)

	// Initialize Enygma Batcher
	enygmaBatcherConfig := &enygmaService.EnygmaBatcherConfig{
		ChainID:        config.MyChainID,
		MaxTxsPerBatch: config.EnygmaBatchSize,
	}
	enygmaBatcher := enygmaService.NewEnygmaBatcher(
		enygmaBatcherConfig,
		r.contractClients.participantStorage,
		enygmaAdapters.NewOTelTracer("enygma-batcher"),
	)

	// Use pre-created endpoint clients
	privateHubEndpointClient := r.contractClients.privateHubEndpoint
	privateNodeEndpointClient := r.contractClients.privateNodeEndpoint

	// Initialize dvp services for enygma (use shared commitment calculator)
	dvpDepositFinder := dvpService.NewDepositFinder(r.repositories.dvpDeposit)

	// Initialize dvp proof service
	dvpProofService := dvpService.NewProofService(
		dvpService.ProofServiceConfig{
			ChainID:            config.MyChainID,
			MerkleTreeDepth:    r.infrastructure.dvpMerkleTreeDepth,
			NumberOfJSParamsIn: config.NumberOfJSParamsIn,
		},
		r.infrastructure.merkleService,
		r.ctsClient,
		r.infrastructure.proofAPIClient,
		r.infrastructure.commitmentCalculator,
		r.repositories.dvpDeposit,
		r.infrastructure.transactionManager,
	)

	// Initialize deposit waiter for consolidation
	depositWaiter := dvpService.NewDepositWaiter(
		dvpService.WaitConfig{
			MaxRetries:    enygmaDepositWaiterMaxRetries,
			RetryInterval: enygmaDepositWaiterRetryInterval,
		},
		r.repositories.dvpDeposit,
	)

	// Use pre-created dvp client for consolidation
	dvpClient := r.contractClients.dvpClient

	// Initialize dvp consolidation service
	dvpConsolidationService := dvpService.NewConsolidationService(
		dvpService.ConsolidationConfig{
			ChainID:               config.MyChainID,
			MaxNumberOfJSDeposits: config.NumberOfJSParamsIn,
		},
		r.repositories.dvpDeposit,
		r.ctsClient,
		r.infrastructure.commitmentCalculator,
		dvpProofService,
		dvpClient,
		r.contractClients.enygmaClient,
		r.contractClients.enygmaIntegrationClient,
		depositWaiter,
		r.infrastructure.transactionManager,
	)

	// Initialize Enygma Proof Service
	enygmaProofService := enygmaService.NewEnygmaProofService(
		config.MyChainID,
		r.ctsClient,
		r.infrastructure.proofAPIClient,
		r.repositories.enygma,
		r.contractClients.enygmaClient,
		enygmaAdapters.NewOTelTracer("enygma-proof"),
	)

	// Initialize Enygma Executor
	enygmaExecutor := enygmaService.NewEnygmaExecutor(
		enygmaService.ExecutorConfig{
			DefaultContextTimeout: config.DefaultContextTimeout,
			MaxNumberOfJSDeposits: config.NumberOfJSParamsIn,
		},
		enygmaAdapters.NewOTelTracer("enygma-executor"),
		enygmaBatcher,
		enygmaProofService,
		r.repositories.enygmaHistory,
		r.ctsClient,
		r.contractClients.enygmaClient,
		r.contractClients.enygmaIntegrationClient,
		r.infrastructure.commitmentCalculator,
		config.MyChainID,
	)

	// Initialize Enygma Finalization Service
	enygmaFinalization := enygmaService.NewEnygmaFinalizationService(
		enygmaAdapters.NewOTelTracer("enygma-finalization"),
		blockWaiter,
		privateHubEndpointClient,
		retryService,
		enygmaExecutor,
	)

	// Initialize Enygma Creation Service
	enygmaCreationService := enygmaService.NewEnygmaCreationService(
		enygmaAdapters.NewOTelTracer("enygma-creation"),
		r.infrastructure.transactionManager,
		r.repositories.enygma,
		r.repositories.enygmaHistory,
	)

	// Initialize Enygma Initiator
	enygmaInitiator := enygmaHandler.NewInitiator(
		enygmaHandler.InitiatorConfig{
			DefaultContextTimeout: config.DefaultContextTimeout,
		},
		r.ctsClient,
		r.contractClients.enygmaHandlerClient,
		privateHubEndpointClient,
		privateNodeEndpointClient,
		r.repositories.enygmaHistory,
		r.repositories.dvpDeposit,
		config.MyChainID,
		dvpDepositFinder,
		r.infrastructure.commitmentCalculator,
		dvpConsolidationService,
		dvpProofService,
		tracer,
		retryService,
		enygmaExecutor,
		enygmaFinalization,
		enygmaCreationService,
		r.infrastructure.transactionManager,
	)

	// Initialize Enygma Revert Service
	enygmaRevertService := enygmaService.NewEnygmaRevertService(
		enygmaAdapters.NewOTelTracer("enygma-revert"),
		privateNodeEndpointClient,
		r.contractClients.enygmaHandlerClient,
	)

	// Initialize Enygma orchestrator
	enygmaOrchestrator := service.NewEnygmaOrchestrator(
		config.TickerPeriod,
		r.msgqueues.enygmaBatchConsumer,
		blockWaiter,
		enygmaBatcher,
		enygmaInitiator,
		enygmaRevertService,
		enygmaFinalization,
		config.EnygmaMaxConcurrentResourceIDs,
	)

	return enygmaOrchestrator, nil
}
