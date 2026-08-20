package run

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/client"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/configinit"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/txops"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/faultinjector"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/logger"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/mtls"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/otel"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/config"
	destapp "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/app"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/healthcheck"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/proof"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/app"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/resultrouter"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	retryHTTPClientRetryCount         = 100
	retryHTTPClientBackoffDuration    = time.Second * 2
	retryHTTPClientRetryCountProofAPI = 3

	routerBatchSize = 100
	routerInterval  = time.Second
)

func Run(path string) error {
	var wg sync.WaitGroup

	ctx, destCancel := context.WithCancel(context.Background())
	defer destCancel()

	// ********************************************************* //
	// ******** Configuration and Logger Initialization ******** //
	conf, err := configinit.InitConfig[config.Config](path)
	if err != nil {
		slog.Error("Error parsing config file", slog.Any("error", err))
		return fmt.Errorf("initializing relayer config: %w", err)
	}

	otelShutdown, err := otel.SetupOTelSDK(context.Background(), conf.OtelSDKDisabled)
	if err != nil {
		slog.Error("Failed to initlaize OTeL", slog.Any("error", err))
		return fmt.Errorf("initializing OpenTelemetry: %w", err)
	}
	defer func() { _ = otelShutdown(context.Background()) }()

	loggerShutdown, err := logger.InitLogger(conf.LogHandler, conf.LogLevel, conf.OtelSDKDisabled)
	if err != nil {
		slog.Error("Failed to initialize logger", slog.Any("error", err))
		return fmt.Errorf("initializing logger: %w", err)
	}
	defer func() { _ = loggerShutdown(context.Background()) }()
	// ******** Configuration and Logger Initialization End ******** //
	// ************************************************************* //

	// ******************************************************** //
	// ************** Connections Initialization ************** //
	dbCtx, dbCancel := context.WithTimeout(ctx, time.Minute)
	defer dbCancel()

	dbConn, err := repository.Connect(dbCtx, conf.DatabaseConnectionString)
	if err != nil {
		slog.Error(
			"Failed to connect to the database",
			slog.Any("error", err),
		)
		return fmt.Errorf("failed to connect to the database: %w", err)
	}
	defer dbConn.Close()

	// Run migrations
	slog.Info("Running database migrations")
	if err := dbConn.Migrate(); err != nil {
		slog.Error("Failed to run migrations", slog.Any("error", err))
		return fmt.Errorf("failed to run database migrations: %w", err)
	}

	natsTLS, err := mtls.LoadClientConfig(conf.NATSTLSCAFile, conf.NATSTLSCertFile, conf.NATSTLSKeyFile)
	if err != nil {
		slog.Error("Failed to load NATS TLS config", slog.Any("error", err))
		return withstack.Wrap(fmt.Errorf("loading NATS client TLS config: %w", err))
	}
	nc, err := nats.Connect(conf.NATSUrl, nats.Secure(natsTLS))
	if err != nil {
		slog.Error("Failed to connect to NATS", slog.Any("error", err))
		return withstack.Wrap(fmt.Errorf("connecting to NATS: %w", err))
	}

	retryTransport := &client.RetryTransport{
		RoundTripper:   http.DefaultTransport,
		InitialDelay:   time.Second,
		DelayIncrement: time.Second,
	}
	httpClient := &http.Client{
		Transport: retryTransport,
	}

	nodeRPCClient, err := rpc.DialOptions(context.TODO(), conf.PrivateNodeURL, rpc.WithHTTPClient(httpClient))
	if err != nil {
		slog.Error("Error creating private ledger client", slog.Any("error", err))
		return withstack.Wrap(fmt.Errorf("creating private node RPC client: %w", err))
	}
	defer nodeRPCClient.Close()
	nodeClient := ethclient.NewClient(nodeRPCClient)

	hubRPCClient, err := rpc.DialOptions(context.TODO(), conf.PrivateHubURL, rpc.WithHTTPClient(httpClient))
	if err != nil {
		slog.Error("Error creating private hub client", slog.Any("error", err))
		return withstack.Wrap(fmt.Errorf("creating private hub RPC client: %w", err))
	}
	defer hubRPCClient.Close()
	hubClient := ethclient.NewClient(hubRPCClient)
	// ************** Connections Initialization End ************** //
	// ************************************************************ //

	// ******************************************************** //
	// ************** Healthcheck Initialization ************** //
	healthcheckConfig := healthcheck.Config{
		Path: "/healthcheck",
		Addr: conf.HealthcheckPort,
	}
	healthcheck := healthcheck.New(healthcheckConfig, nc, dbConn.Pool(), hubClient, nodeClient)

	slog.Info("Starting healthcheck")
	wg.Add(1)
	go func() {
		defer wg.Done()
		if srvErr := healthcheck.ListenAndServe(); srvErr != nil && !errors.Is(srvErr, http.ErrServerClosed) {
			slog.Error("Failed to shutdown healthcheck server", slog.Any("error", srvErr))
		}
	}()
	// ************** Healthcheck Initialization End ************** //
	// ************************************************************ //

	// ********************************************************** //
	// ************** Fault Injection Initialization ************* //
	var faultServer *http.Server
	if conf.FaultInjectionEnabled {
		if conf.FaultInjectionPersistPath != "" {
			faultinjector.SetPersistPath(conf.FaultInjectionPersistPath)
		}
		faultinjector.Enable()
		faultServer = faultinjector.NewHTTPServer(":" + conf.FaultInjectionPort)
		if faultServer != nil {
			slog.Warn("Starting fault injection API server", slog.String("addr", ":"+conf.FaultInjectionPort))
			wg.Add(1)
			go func() {
				defer wg.Done()
				if srvErr := faultServer.ListenAndServe(); srvErr != nil && !errors.Is(srvErr, http.ErrServerClosed) {
					slog.Error("Fault injection server error", slog.Any("error", srvErr))
				}
			}()
		} else {
			slog.Warn("Fault injection enabled in config but binary built without faultinjection tag — no-op")
		}
	}
	// ************** Fault Injection Initialization End ********** //
	// ************************************************************ //

	// ************************************************************************** //
	// ************** Key Creation, Registration and Authorization ************** //
	ctsTLS, err := mtls.LoadClientConfig(conf.CTSTLSCAFile, conf.CTSTLSCertFile, conf.CTSTLSKeyFile)
	if err != nil {
		return fmt.Errorf("loading CTS client TLS config: %w", err)
	}
	ctsConn, err := grpc.NewClient(conf.CTSRootURL,
		grpc.WithTransportCredentials(credentials.NewTLS(ctsTLS)),
	)
	if err != nil {
		return fmt.Errorf("connect to grpc server: %w", err)
	}
	defer ctsConn.Close()

	// ********************************************************** //
	// ************** Header Proof Service Init ***************** //

	// Step 4: Set up the contract client infrastructure used for registering audit info
	hubRegistry, err := contractclient.NewDeploymentProxyRegistryClient(conf.PrivateHubDeploymentProxyRegistry, hubClient)
	if err != nil {
		return fmt.Errorf("failed to create private hub registry: %w", err)
	}

	if err != nil {
		return fmt.Errorf("failed to create private hub auth generator;: %w", err)
	}

	hubExecutor := contractclient.NewCTSExecutor(contractclient.NewDefaultRetryingTxOpsClient(txops.NewPrivateHubTxOpsServiceClient(ctsConn)))

	proofsAddress, err := hubRegistry.GetContractAddress("Proofs")
	if err != nil {
		return fmt.Errorf("failed to get proofs address: %w", err)
	}
	proofsClient := contractclient.NewProofsClient(proofsAddress, hubExecutor)

	headerProofService, err := proof.NewHeaderProofService(
		proof.DefaultHeaderProofConfig(conf.PrivateNodeChainID),
		nodeClient,
		proofsClient,
	)
	if err != nil {
		return fmt.Errorf("failed to create header proof service: %w", err)
	}

	slog.Info("Starting header proof service...")
	wg.Add(1)
	go func() {
		defer wg.Done()
		if proofErr := headerProofService.Run(ctx); proofErr != nil {
			slog.Error("Header proof service exited with error", slog.Any("error", proofErr))
		}
	}()
	// ************** Header Proof Service Init End ************* //
	// ********************************************************** //

	// ****************************************************************** //
	// **************** Async Message Result Router Init **************** //
	js, err := jetstream.New(nc)
	if err != nil {
		return fmt.Errorf("failed to create jetstream connection: %w", err)
	}

	// Scoped timeout for manager / consumer creation only — the routers'
	// lifetime is the outer process ctx so they keep running for the
	// entire process, not just 30 seconds.
	mqInitCtx, mqInitCancel := context.WithTimeout(ctx, 30*time.Second)
	defer mqInitCancel()

	manager, err := msgqueue.NewManager(mqInitCtx, js, conf.PrivateNodeChainID.String())
	if err != nil {
		return fmt.Errorf("failed to create message queue manager: %w", err)
	}

	privateHubResultConsumer, err := msgqueue.NewConsumer[types.TxResult](
		mqInitCtx,
		manager,
		"source_privatehub_results",
		"cts.result.privatehub",
	)
	if err != nil {
		return fmt.Errorf("failed to initialize privatehub result consumer: %w", err)
	}

	privateNodeResultConsumer, err := msgqueue.NewConsumer[types.TxResult](
		mqInitCtx,
		manager,
		"source_privatenode_results",
		"cts.result.privatenode",
	)
	if err != nil {
		return fmt.Errorf("failed to initialize privatenode result consumer: %w", err)
	}
	pnResultRouter := resultrouter.New(
		resultrouter.Config{
			Identity:  "privatenode",
			BatchSize: routerBatchSize,
			Interval:  routerInterval,
		},
		privateNodeResultConsumer,
	)
	pnhResultRouter := resultrouter.New(
		resultrouter.Config{
			Identity:  "privatehub",
			BatchSize: routerBatchSize,
			Interval:  routerInterval,
		},
		privateHubResultConsumer,
	)

	// **************** Async Message Result Router Init End **************** //
	// ********************************************************************** //

	// ************************************************************ //
	// ******** Private Relayer Source Side Initialization ******** //
	sourceRelayer := app.New(nc, dbConn.Pool(), hubClient, nodeClient, ctsConn, pnResultRouter, pnhResultRouter)
	if initErr := sourceRelayer.Initialize(ctx, *conf); initErr != nil {
		slog.Error("Failed to initialize source private relayer", slog.Any("error", initErr))
		return initErr
	}

	// Graceful shutdown handling
	go func() {
		slog.Info("Source private relayer service starting...")
		if err = sourceRelayer.Run(); err == nil {
			slog.Info("Stopped relaying new messages.")
		} else {
			slog.Error("Encountered fatal error. Shutting down.", slog.Any("error", err))
		}
	}()
	// ******** Private Relayer Source Side Initialization End ******** //
	// **************************************************************** //

	// ************************************************************ //
	// ******** Private Relayer Dest Side Initialization ********** //
	destRelayer := destapp.New(nc, dbConn.Pool(), hubClient, nodeClient, ctsConn, pnResultRouter, pnhResultRouter)
	if initErr := destRelayer.Initialize(ctx, *conf); initErr != nil {
		slog.Error("Failed to initialize dest private relayer", slog.Any("error", initErr))
		return initErr
	}

	// Graceful shutdown handling
	go func() {
		slog.Info("Dest private relayer service starting...")
		if err = destRelayer.Run(); err == nil {
			slog.Info("Stopped relaying new messages (dest).")
		} else {
			slog.Error("Encountered fatal error (dest). Shutting down.", slog.Any("error", err))
		}
	}()
	// ******** Private Relayer Dest Side Initialization End ******** //
	// ************************************************************** //

	slog.Info("Starting Private Node Result Router...")
	wg.Add(1)
	go func() {
		defer wg.Done()
		pnResultRouter.Run(ctx)
	}()

	slog.Info("Starting Private Hub Result Router...")
	wg.Add(1)
	go func() {
		defer wg.Done()
		pnhResultRouter.Run(ctx)
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	// **************************************************//
	// ************* Healthckeck  Shutdown ************* //
	const healthcheckShutdownTimeout = 5 * time.Second
	hcCtx, hcCancel := context.WithTimeout(context.Background(), healthcheckShutdownTimeout)
	defer hcCancel()

	err = healthcheck.Shutdown(hcCtx)
	if err != nil {
		slog.Error("Error while shutting down health check", slog.Any("error", err))
	}
	// ************* Healthckeck  Shutdown End ************* //
	// ******************************************************//

	// ******************************************************** //
	// ************* Fault Injection Shutdown ***************** //
	if faultServer != nil {
		const faultInjectionShutdownTimeout = 5 * time.Second
		fiCtx, fiCancel := context.WithTimeout(context.Background(), faultInjectionShutdownTimeout)
		if shutdownErr := faultServer.Shutdown(fiCtx); shutdownErr != nil {
			slog.Error("Error while shutting down fault injection server", slog.Any("error", shutdownErr))
		}
		fiCancel()
	}
	// ************* Fault Injection Shutdown End ************* //
	// ********************************************************* //

	destCancel()

	// *******************************************************//
	// ******** Private Relayer Source Side Shutdown ******** //
	slog.Info("Shutting down source private relayer service...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := sourceRelayer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Source Shutdown error", slog.Any("error", err))
	}

	slog.Info("Shutting down dest private relayer service...")
	if err := destRelayer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Dest Shutdown error", slog.Any("error", err))
	}
	slog.Info("Shutdown complete.")
	// ******** Private Relayer Source Side Shutdown End ******** //
	// ********************************************************** //

	wg.Wait()

	slog.Info("Program exited gracefully")

	return nil
}

func logOnError(err error, msg string) {
	if err != nil {
		slog.Error(msg, slog.Any("error", err))
	}
}

type txWaiterAdapter struct {
	client bind.DeployBackend
}

func (a *txWaiterAdapter) WaitMined(ctx context.Context, tx *ethTypes.Transaction) (*ethTypes.Receipt, error) {
	return bind.WaitMined(ctx, a.client, tx)
}
