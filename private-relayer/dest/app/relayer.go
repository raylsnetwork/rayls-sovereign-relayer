package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/config"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/resultrouter"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// CTSClient composes the gRPC key, encrypt, and TxOps service clients from
// the Cryptography Trust Suite. It satisfies the consumer-defined
// interfaces throughout the relayer (each consumer picks the methods it
// needs). The TxOps embedded client drives the signing migration — private
// hub and private node signing both go through CTS instead of a local
// key queue.
type DestPrivateRelayer struct {
	// Raw dependencies (passed in)
	pool            *pgxpool.Pool
	natsConn        *nats.Conn
	hubClient       *ethclient.Client
	nodeClient      *ethclient.Client
	ctsConn         *grpc.ClientConn
	pnResultRouter  *resultrouter.Router
	pnhResultRouter *resultrouter.Router

	// Derived components
	repositories    *Repositories
	contractClients *ContractClients
	msgqueues       *MessageQueue
	infrastructure  *Infrastructure
	ctsClient       *CTSClient
	txGen           *TransactionGenerator
	batchers        *Batchers
	receivers       *Receivers
	services        *Services
	logRouter       *LogRouter
	logParsers      *LogParsers
	listeners       *Listeners

	atomicServices *AtomicServices
	atomicPollers  *AtomicPollers

	eg       *errgroup.Group
	shutdown context.CancelFunc
}

// New creates a new destination private relayer with external dependencies injected
func New(
	natsConn *nats.Conn,
	pool *pgxpool.Pool,
	hubClient *ethclient.Client,
	nodeClient *ethclient.Client,
	ctsConn *grpc.ClientConn,
	pnResultRouter *resultrouter.Router,
	pnhResultRouter *resultrouter.Router,
) *DestPrivateRelayer {
	return &DestPrivateRelayer{
		pool:       pool,
		natsConn:   natsConn,
		hubClient:  hubClient,
		nodeClient: nodeClient,
		ctsConn:    ctsConn,

		pnResultRouter:  pnResultRouter,
		pnhResultRouter: pnhResultRouter,
	}
}

// Initialize sets up all components in dependency order
func (r *DestPrivateRelayer) Initialize(ctx context.Context, conf config.Config) error {
	var err error

	// Initialize repositories (depends on pool)
	r.initializeRepositories(ctx)

	// Initialize CTS Clients
	r.initializeCTSClients()

	// Initialize contract clients
	err = r.initializeContractClients(ContractClientConfig{ //nolint:contextcheck // initialization code
		PrivateHubChainID:  conf.PrivateHubChainID,
		PrivateNodeChainID: conf.PrivateNodeChainID,
		VENChainID:         conf.VENChainID,

		PrivateHubDeploymentProxyRegistry:  conf.PrivateHubDeploymentProxyRegistry,
		PrivateNodeDeploymentProxyRegistry: conf.PrivateNodeDeploymentProxyRegistry,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize contract clients: %w", err)
	}

	// Initialize message queues (depends on natsConn)
	//nolint:contextcheck // initialization code, no parent context available
	err = r.initializeMessageQueues(conf.PrivateNodeChainID.String())
	if err != nil {
		return fmt.Errorf("failed to initialize message queues: %w", err)
	}

	// Initialize infrastructure
	r.initializeInfrastructure(InfrastructureConfig{
		ProofAPIURL:        conf.ProofAPIURL,
		DvpMerkleTreeDepth: conf.PrivateHubDvpMerkleTreeDepth,
	})

	// Initialize transaction generator
	err = r.initializeTransactionGenerator()
	if err != nil {
		return fmt.Errorf("failed to initialize transaction generator: %w", err)
	}

	// Initialize batchers
	r.initializeBatchers(BatcherConfig{
		DefaultContextTimeout: conf.DefaultContextTimeout,
		PrivateNodeChainID:    conf.PrivateNodeChainID,
	})

	// Initialize receivers
	err = r.initializeReceivers(ReceiverConfig{ //nolint:contextcheck // initialization code
		DefaultContextTimeout: conf.DefaultContextTimeout,
		NumberOfJSParamsIn:    conf.NumberOfJSParamsIn,
		PrivateNodeChainID:    conf.PrivateNodeChainID,
		PrivateHubChainID:     conf.PrivateHubChainID,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize receivers: %w", err)
	}

	// Initialize log parsers
	err = r.initializeLogParsers(LogParsersConfig{
		PrivateNodeChainID: conf.PrivateNodeChainID,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize log parsers: %w", err)
	}

	// Initialize atomic services BEFORE the regular services because
	// CrossChainService now takes the receipt + vanilla-receipt
	// services directly (instead of going through the old executor +
	// poller indirection).
	err = r.initializeAtomic(AtomicConfig{ //nolint:contextcheck // initialization code
		BatchSize: conf.ExecutorBatchMessages,
		MyChainID: conf.PrivateNodeChainID,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize atomic services: %w", err)
	}

	// Initialize services
	err = r.initializeServices(ServiceConfig{
		MyChainID:              conf.PrivateNodeChainID,
		PrivateHubChainID:      conf.PrivateHubChainID,
		DefaultContextTimeout:  conf.DefaultContextTimeout,
		PrivateHubTickerPeriod: time.Second, // Default 1s ticker
		EnygmaSyncMaxRetries:   10,          // Retry checkpoint validation 10 times before triggering resync
	})
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	// Initialize listeners (handles address resolution and log router internally)
	//nolint:contextcheck // initialization code, no parent context available
	err = r.initializeListeners(ListenersConfig{
		BatchSize:               conf.ListenersBlockBatchSize,
		PrivateHubStartingBlock: conf.PrivateHubStartingBlock,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize listeners: %w", err)
	}

	return nil
}

// Run starts all services using errgroup pattern
func (r *DestPrivateRelayer) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	r.shutdown = cancel

	eg, egCtx := errgroup.WithContext(ctx)
	// Copy the eg to the relayer state so we can wait on it in the
	// shutdown routine and not close connections prematurely.
	r.eg = eg

	slog.Info("Starting Private Hub Listener (listening to source messages)...")
	eg.Go(func() error {
		return r.listeners.privateHub.Run(egCtx)
	})

	// LogRouter is invoked by the listener via Handle method - no separate goroutine needed

	slog.Info("Starting Endpoint Log Parser...")
	eg.Go(func() error {
		return r.logParsers.endpointParser.Run(egCtx)
	})

	slog.Info("Starting Enygma Teleport Log Parser...")
	eg.Go(func() error {
		return r.logParsers.enygmaTeleportParser.Run(egCtx)
	})

	slog.Info("Starting Teleport Log Parser (Cross-Chain)...")
	eg.Go(func() error {
		return r.logParsers.teleportParser.Run(egCtx)
	})

	slog.Info("Starting Dvp Teleport Log Parser...")
	eg.Go(func() error {
		return r.logParsers.dvpTeleportParser.Run(egCtx)
	})

	slog.Info("Starting Audit Manager Log Parser...")
	eg.Go(func() error {
		return r.logParsers.auditManagerParser.Run(egCtx)
	})

	slog.Info("Starting PrivateHub Service...")
	eg.Go(func() error {
		return r.services.privateHub.Run(egCtx)
	})

	slog.Info("Starting Cross Chain Service...")
	eg.Go(func() error {
		return r.services.crossChain.Run(egCtx)
	})

	slog.Info("Starting Enygma Orchestrator...")
	eg.Go(func() error {
		return r.services.enygmaOrchestrator.Run(egCtx)
	})

	slog.Info("Starting Enygma Sync Service...")
	eg.Go(func() error {
		return r.runEnygmaSyncService(egCtx)
	})

	slog.Info("Starting Dvp Orchestrator...")
	eg.Go(func() error {
		return r.services.dvpOrchestrator.Run(egCtx)
	})

	slog.Info("Starting Atomic Dest Finalization Poller...")
	eg.Go(func() error {
		return r.atomicPollers.destFinalization.Run(egCtx)
	})

	return eg.Wait()
}

// runEnygmaSyncService runs the enygma sync service on a ticker
func (r *DestPrivateRelayer) runEnygmaSyncService(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Debug("Enygma Sync Service shutting down")
			return nil
		case <-ticker.C:
			if err := r.services.enygmaSyncService.Run(ctx); err != nil {
				slog.Error("Enygma Sync: Could not validate checkpoints",
					slog.Any("error", err))
			}
		}
	}
}

// Shutdown gracefully shuts down all services
func (r *DestPrivateRelayer) Shutdown(ctx context.Context) error {
	// Signal all goroutines to stop
	r.shutdown()

	// Wait for all goroutines to finish
	_ = r.eg.Wait()

	// Close message queue (drains connections)
	if r.msgqueues != nil {
		r.msgqueues.Close()
	}

	// Note: We don't close pool, natsConn, hubWrapper, or nodeWrapper
	// because they're owned by the caller who passed them in

	return nil
}
