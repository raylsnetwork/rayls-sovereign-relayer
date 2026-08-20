// Decommissioning Teleport (vanilla, atomic).

// Package app implements the legacy public-chain (RN) Teleport bridge relayer.
//
// Deprecated: Decommissioning Teleport (vanilla, atomic).
package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/public-relayer/config"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/txsim"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type PublicRelayer struct {
	// Raw dependencies (passed in by caller)
	pool          *pgxpool.Pool
	natsConn      *nats.Conn
	publicClient  *ethclient.Client
	privateClient *ethclient.Client
	ctsConn       *grpc.ClientConn

	repositories    *Repositories
	ctsClients      *CTSClients
	contractClients *ContractClients
	messageQueues   *MessageQueues
	listeners       *Listeners
	services        *Services

	healthCheckServer *http.Server

	eg       *errgroup.Group
	shutdown context.CancelFunc
}

// New creates a new PublicRelayer with external dependencies injected.
// Call Initialize() after to wire internal components.
func New(
	natsConn *nats.Conn,
	pool *pgxpool.Pool,
	publicClient *ethclient.Client,
	privateClient *ethclient.Client,
	ctsConn *grpc.ClientConn,
) *PublicRelayer {
	return &PublicRelayer{
		pool:          pool,
		natsConn:      natsConn,
		publicClient:  publicClient,
		privateClient: privateClient,
		ctsConn:       ctsConn,
	}
}

// Initialize sets up all components in dependency order.
func (p *PublicRelayer) Initialize(conf *config.Config, nc *nats.Conn) error {
	var err error

	err = txsim.PopulateErrorMap("./contracts/")
	if err != nil {
		return fmt.Errorf("failed to populate errors ABI map: %w", err)
	}

	p.initializeRepositories()
	p.initializeCTSClients()

	err = p.initializeContractClients(ContractClientsConfig{
		PublicDeploymentProxyRegistryAddress:  conf.PublicChainDeploymentProxyRegistry,
		PrivateDeploymentProxyRegistryAddress: conf.PrivateNodeDeploymentProxyRegistry,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize contracts clients: %w", err)
	}

	err = p.initializeMessageQueues(MessageQueueConfig{
		NATSConn: nc,
		ChainId:  conf.PrivateNodeChainID.String(),
	})
	if err != nil {
		return fmt.Errorf("failed to initialize message queues: %w", err)
	}

	err = p.initializeServices(ServiceConfig{
		RelayInterval: time.Second,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	err = p.initializeListeners(context.Background(), ListenersConfig{
		BatchSize: conf.ListenersBlockBatchSize,

		PublicChainStartingBlock:  conf.PublicChainStartingBlock,
		PrivateChainStartingBlock: conf.PrivateNodeStartingBlock,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize listeners: %w", err)
	}

	p.initializeHealthCheckServer(HealthCheckConfig{
		Addr: ":9000",
		Path: "/healthcheck",
	})

	return nil
}

func (p *PublicRelayer) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	p.shutdown = cancel

	eg, egCtx := errgroup.WithContext(ctx)
	// Copy the eg to the pubrelayer state so we can wait on it in the
	// shutdown routine and not close connections prematurely.
	p.eg = eg

	slog.Info("Starting Health Check")
	eg.Go(func() error {
		return p.healthCheckServer.ListenAndServe()
	})

	slog.Info("Starting Public Listener...")
	eg.Go(func() error {
		return p.listeners.public.Run(egCtx)
	})

	slog.Info("Starting Private Listener...")
	eg.Go(func() error {
		return p.listeners.private.Run(egCtx)
	})

	slog.Info("Starting Public Generator Service...")
	eg.Go(func() error {
		return p.services.publicGenerator.Run(egCtx)
	})

	slog.Info("Starting Private Generator Service...")
	eg.Go(func() error {
		return p.services.privateGenerator.Run(egCtx)
	})

	slog.Info("Starting Deployer Service...")
	eg.Go(func() error {
		return p.services.deployer.Deploy(egCtx)
	})

	slog.Info("Starting publicchain Result Router...")
	eg.Go(func() error {
		return p.services.publicChainResultRouter.Run(egCtx)
	})

	slog.Info("Starting privatechain Result Router...")
	eg.Go(func() error {
		return p.services.privateChainResultRouter.Run(egCtx)
	})

	return eg.Wait()
}

// Shutdown gracefully shuts down all services.
func (p *PublicRelayer) Shutdown(ctx context.Context) error {
	// Shut down healthcheck server
	if p.healthCheckServer != nil {
		if err := p.healthCheckServer.Shutdown(ctx); err != nil {
			slog.Error("Error while shutting down health check", slog.Any("error", err))
		}
	}

	p.shutdown()
	_ = p.eg.Wait()

	// Close message queues (drains subscriptions)
	if p.messageQueues != nil {
		p.messageQueues.Close()
	}

	// Note: We don't close pool, natsConn, publicClient, or privateClient
	// because they're owned by the caller who passed them in

	return nil
}
