package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/config"
	"github.com/raylsnetwork/rayls-sovereign-relayer/resultrouter"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type SourcePrivateRelayer struct {
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
	msgqueues       *MessageQueue
	ctsClient       *CTSClient
	contractClients *ContractClients
	listeners       *Listeners
	logParsers      *LogParsers
	txGen           *TransactionGenerator
	services        *Services

	infrastructure *Infrastructure

	// Decommissioning Teleport (vanilla, atomic); declarations in source/app/atomic.go.
	atomicServices *AtomicServices
	atomicPollers  *AtomicPollers

	eg       *errgroup.Group
	shutdown context.CancelFunc
}

func New(
	natsConn *nats.Conn,
	pool *pgxpool.Pool,
	hubClient *ethclient.Client,
	nodeClient *ethclient.Client,
	ctsConn *grpc.ClientConn,
	pnResultRouter *resultrouter.Router,
	pnhResultRouter *resultrouter.Router,
) *SourcePrivateRelayer {
	return &SourcePrivateRelayer{
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
func (r *SourcePrivateRelayer) Initialize(ctx context.Context, conf config.Config) error {
	var err error

	// Initialize repositories (depends on pool)
	r.initializeRepositories()

	// Initial;ize CTS grpc clients
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
	//nolint:contextcheck // initialization code
	err = r.initializeMessageQueues(
		conf.PrivateNodeChainID.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize message queues: %w", err)
	}

	// Initialize transaction generator
	err = r.initializeTransactionGenerator()
	if err != nil {
		return fmt.Errorf("failed to initialize transaction generator: %w", err)
	}

	// Initialize infrastructure
	r.initializeInfrastructure(InfrastructureConfig{
		ProofAPIURL:        conf.ProofAPIURL,
		DvpMerkleTreeDepth: conf.PrivateHubDvpMerkleTreeDepth,
	})

	// Initialize listeners
	//nolint:contextcheck // initialization code
	err = r.initializeListeners(
		ListenersConfig{
			BatchSize:                conf.ListenersBlockBatchSize,
			PrivateNodeStartingBlock: conf.PrivateNodeStartingBlock,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to initialize listeners: %w", err)
	}

	// Initialize atomic services and pollers BEFORE services because
	// the result router constructed there registers
	// SignatureService.HandleSourceRevertedCallback as a handler.
	// Decommissioning Teleport (vanilla, atomic).
	err = r.initializeAtomic(AtomicConfig{ //nolint:contextcheck // initialization code, no parent context available
		BatchSize:      conf.ExecutorBatchMessages,
		ExpirationTime: conf.ExpirationTime,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize atomic services: %w", err)
	}

	// Initialize services
	err = r.initializeServices(ctx, ServiceConfig{
		PrivateHubChainID:      conf.PrivateHubChainID,
		MyChainID:              conf.PrivateNodeChainID,
		CrossChainTickerPeriod: time.Second,
		PrivateHubTickerPeriod: time.Second,
		EnygmaTickerPeriod:     time.Second,

		EnygmaBatchSize:                conf.EnygmaBatchSize,
		EnygmaMaxConcurrentResourceIDs: conf.EnygmaMaxConcurrentResourceIDs,
		NumberOfJSParamsIn:             conf.NumberOfJSParamsIn,
		DefaultContextTimeout:          conf.DefaultContextTimeout,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	// Initialize log parsers (after services — log parsers depend on r.services.keysService)
	err = r.initializeLogParsers(LogpParsersConfig{
		PrivateHubChainID: conf.PrivateHubChainID,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize log parsers: %w", err)
	}

	return nil
}

func (r *SourcePrivateRelayer) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	r.shutdown = cancel

	eg, egCtx := errgroup.WithContext(ctx)
	// Copy the eg to the relayer state so we can wait on it in the
	// shutdown routine and not close connections prematurely.
	r.eg = eg

	slog.Info("Starting Private Node Listener...")
	eg.Go(func() error {
		return r.listeners.privateNode.Run(egCtx)
	})

	slog.Info("Starting Endpoint Log Parser...")
	eg.Go(func() error {
		return r.logParsers.endpointParser.Fetch(egCtx)
	})

	slog.Info("Starting Cross Chain Service...")
	eg.Go(func() error {
		return r.services.crossChain.Run(egCtx)
	})

	slog.Info("Starting Private Hub Service...")
	eg.Go(func() error {
		return r.services.privateHub.Run(egCtx)
	})

	slog.Info("Starting Enygma Log Parser...")
	eg.Go(func() error {
		return r.logParsers.enygmaParser.Fetch(egCtx)
	})

	slog.Info("Starting Enygma Orchestrator...")
	eg.Go(func() error {
		return r.services.enygmaOrchestrator.Run(egCtx)
	})

	slog.Info("Starting Dvp Orchestrator...")
	eg.Go(func() error {
		return r.services.dvpOrchestrator.Run(egCtx)
	})

	slog.Info("Starting Atomic EarlyRevert Poller...")
	eg.Go(func() error {
		return r.atomicPollers.earlyRevert.Run(egCtx)
	})

	slog.Info("Starting Atomic Expired Poller...")
	eg.Go(func() error {
		return r.atomicPollers.expired.Run(egCtx)
	})

	slog.Info("Starting Atomic Source Finalization Poller...")
	eg.Go(func() error {
		return r.atomicPollers.srcFinalization.Run(egCtx)
	})

	return eg.Wait()
}

func (r *SourcePrivateRelayer) Shutdown(ctx context.Context) error {
	// Signal all goroutines to stop
	r.shutdown()

	// Wait for all goroutines to finish
	_ = r.eg.Wait()

	// Close message queue (drains connections)
	if r.msgqueues != nil {
		r.msgqueues.Close()
	}

	return nil
}
