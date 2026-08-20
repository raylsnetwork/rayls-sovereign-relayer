// Decommissioning Teleport (vanilla, atomic).

// Package run implements the legacy public-chain (RN) Teleport bridge relayer.
//
// Deprecated: Decommissioning Teleport (vanilla, atomic).
package run

import (
	"errors"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/nats-io/nats.go"
	"github.com/raylsnetwork/rayls-sovereign-relayer/client"
	"github.com/raylsnetwork/rayls-sovereign-relayer/configinit"
	"github.com/raylsnetwork/rayls-sovereign-relayer/faultinjector"
	"github.com/raylsnetwork/rayls-sovereign-relayer/logger"
	"github.com/raylsnetwork/rayls-sovereign-relayer/mtls"
	"github.com/raylsnetwork/rayls-sovereign-relayer/otel"
	"github.com/raylsnetwork/rayls-sovereign-relayer/public-relayer/app"
	"github.com/raylsnetwork/rayls-sovereign-relayer/public-relayer/config"
	"github.com/raylsnetwork/rayls-sovereign-relayer/public-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	shutdownTimeout         = 30 * time.Second
	publicRelayerRetryCount = 5
	publicRelayerBackoff    = 5 * time.Second
)

func Run(path string) error {
	var wg sync.WaitGroup

	// ********************************************************* //
	// ******** Configuration and Logger Initialization ******** //
	conf, err := configinit.InitConfig[config.Config](path)
	if err != nil {
		slog.Error("Error parsing config file", slog.Any("error", err))
		return fmt.Errorf("initializing public relayer config: %w", err)
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

	// ******************************************************** //
	// ************** Connections Initialization ************** //

	// PostgreSQL
	dbCtx, dbCancel := context.WithTimeout(context.Background(), time.Minute)
	defer dbCancel()

	dbConn, err := repository.Connect(dbCtx, conf.DatabaseConnectionString)
	if err != nil {
		slog.Error("Failed to connect to the database", slog.Any("error", err))
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer dbConn.Close()

	slog.Info("Running database migrations")
	if err := dbConn.Migrate(); err != nil {
		slog.Error("Failed to run database migrations", slog.Any("error", err))
		return fmt.Errorf("running database migrations: %w", err)
	}

	// NATS (mTLS)
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
	defer func() { _ = nc.Drain() }()

	// Ethereum RPC clients with retry transport
	retryTransport := &client.RetryTransport{
		RoundTripper:   http.DefaultTransport,
		InitialDelay:   time.Second,
		DelayIncrement: time.Second,
	}
	httpClient := &http.Client{
		Transport: retryTransport,
	}

	rpcDialCtx, rpcDialCancel := context.WithTimeout(context.Background(), time.Minute)
	defer rpcDialCancel()

	publicRPCClient, err := rpc.DialOptions(rpcDialCtx, conf.PublicChainURL, rpc.WithHTTPClient(httpClient))
	if err != nil {
		slog.Error("Failed to connect to public chain RPC", slog.Any("error", err))
		return withstack.Wrap(fmt.Errorf("creating public chain RPC client: %w", err))
	}
	defer publicRPCClient.Close()
	publicClient := ethclient.NewClient(publicRPCClient)

	privateRPCClient, err := rpc.DialOptions(rpcDialCtx, conf.PrivateNodeURL, rpc.WithHTTPClient(httpClient))
	if err != nil {
		slog.Error("Failed to connect to private node RPC", slog.Any("error", err))
		return withstack.Wrap(fmt.Errorf("creating private node RPC client: %w", err))
	}
	defer privateRPCClient.Close()
	privateClient := ethclient.NewClient(privateRPCClient)

	// CTS gRPC connection (mTLS)
	ctsTLS, err := mtls.LoadClientConfig(conf.CTSTLSCAFile, conf.CTSTLSCertFile, conf.CTSTLSKeyFile)
	if err != nil {
		return fmt.Errorf("loading CTS client TLS config: %w", err)
	}
	ctsConn, err := grpc.NewClient(conf.CTSRootURL,
		grpc.WithTransportCredentials(credentials.NewTLS(ctsTLS)),
	)
	if err != nil {
		return fmt.Errorf("connecting to CTS gRPC server: %w", err)
	}
	defer ctsConn.Close()

	// ************** Connections Initialization End ************** //
	// ************************************************************ //

	relayer := app.New(nc, dbConn.Pool(), publicClient, privateClient, ctsConn)

	err = relayer.Initialize(conf, nc)
	if err != nil {
		slog.Error("Failed to initialize public relayer", slog.Any("error", err))
		return fmt.Errorf("initializing public relayer: %w", err)
	}

	// Graceful shutdown handling
	go func() {
		slog.Info("Public relayer service starting...")
		if runErr := relayer.Run(); runErr != nil {
			slog.Error("Encountered fatal error. Shutting down.", slog.Any("error", runErr))
		} else {
			slog.Info("Stopped relaying new messages.")
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	slog.Info("Shutting down public relayer service...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := relayer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Shutdown error", slog.Any("error", err))
	}

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

	wg.Wait()

	slog.Info("Shutdown complete.")

	return nil
}

func logOnError(err error, msg string) {
	if err != nil {
		slog.Error(msg, slog.Any("error", err))
	}
}
