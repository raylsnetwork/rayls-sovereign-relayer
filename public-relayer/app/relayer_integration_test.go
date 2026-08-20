// Decommissioning Teleport (vanilla, atomic).

//go:build integration

package app_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-sovereign-relayer/keyqueue"
	"github.com/raylsnetwork/rayls-sovereign-relayer/public-relayer/app"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
)

var (
	publicRegistryAddr  = common.HexToAddress("0x1000000000000000000000000000000000000001")
	privateRegistryAddr = common.HexToAddress("0x2000000000000000000000000000000000000001")

	// Public chain contracts
	publicEndpointAddr   = common.HexToAddress("0xA001")
	publicDispatcherAddr = common.HexToAddress("0xA002")
	publicAccessMgrAddr    = common.HexToAddress("0xA003")

	// Private chain contracts
	privateEndpointAddr   = common.HexToAddress("0xB001")
	privateDispatcherAddr = common.HexToAddress("0xB002")
	privateAccessMgrAddr  = common.HexToAddress("0xB003")
	tokenRegistryAddr     = common.HexToAddress("0xB004")
	tokenCoreAddr         = common.HexToAddress("0xB005")
)

var dbConf = testtools.DBConfig{
	User:           "test",
	Pass:           "test",
	Database:       "testdb",
	MigrationsPath: "file://../repository/migrations",
}

func TestPublicRelayer_Integration(t *testing.T) {
	// Set working directory to repo root so txsim.PopulateErrorMap("./contracts/") works.
	t.Chdir("../..")

	testtools.SilenceLogger()

	// Start infrastructure containers
	pool, pgCleanup := testtools.SetupPostgres(t, dbConf)
	defer pgCleanup()

	natsConn, natsCleanup := testtools.SetupNATS(t)
	defer natsCleanup()

	// Set up fake Ethereum JSON-RPC servers
	publicChainID := big.NewInt(100)
	privateChainID := big.NewInt(200)

	publicServer := testtools.NewFakeEthServer(t, testtools.FakeEthServerConfig{
		ChainID:  publicChainID,
		GasPrice: big.NewInt(1_000_000_000),
		DeploymentRegistries: map[common.Address]testtools.DeploymentRegistry{
			publicRegistryAddr: {
				Names: []string{
					"PublicRNEndpoint", "RNMessageDispatcher", "RaylsAccessManager",
				},
				Addresses: []common.Address{
					publicEndpointAddr, publicDispatcherAddr, publicAccessMgrAddr,
				},
			},
		},
	})

	privateServer := testtools.NewFakeEthServer(t, testtools.FakeEthServerConfig{
		ChainID:  privateChainID,
		GasPrice: big.NewInt(1_000_000_000),
		DeploymentRegistries: map[common.Address]testtools.DeploymentRegistry{
			privateRegistryAddr: {
				Names: []string{
					"RNEndpoint", "RNMessageDispatcher", "RaylsAccessManager", "TokenRegistry", "TokenCore",
				},
				Addresses: []common.Address{
					privateEndpointAddr, privateDispatcherAddr, privateAccessMgrAddr, tokenRegistryAddr, tokenCoreAddr,
				},
			},
		},
	})

	// Create ethclient.Client instances pointing at fake servers
	publicRPC, err := rpc.Dial(publicServer.URL)
	require.NoError(t, err)
	t.Cleanup(func() { publicRPC.Close() })
	publicClient := ethclient.NewClient(publicRPC)

	privateRPC, err := rpc.Dial(privateServer.URL)
	require.NoError(t, err)
	t.Cleanup(func() { privateRPC.Close() })
	privateClient := ethclient.NewClient(privateRPC)

	// Key queues with buffer capacity (no real keys needed for init)
	keyQueues := &app.KeyQueues{
		PublicChainKeys:       keyqueue.New(3),
		PrivateNodeKeys:       keyqueue.New(3),
		PublicChainKeyStrings: []string{},
		PrivateNodeKeyStrings: []string{},
	}

	contractsConfig := app.ContractsConfig{
		PublicDeploymentProxyRegistryAddress:  publicRegistryAddr,
		PrivateDeploymentProxyRegistryAddress: privateRegistryAddr,
	}
	listenersConfig := app.ListenersConfig{
		BatchSize:                 100,
		PublicChainStartingBlock:  big.NewInt(0),
		PrivateChainStartingBlock: big.NewInt(0),
	}
	executorConfig := app.ExecutorConfig{
		SendInterval:    time.Second,
		ReceiptInterval: time.Second,
	}
	servicesConfig := app.ServiceConfig{
		RelayInterval: time.Second,
	}
	healthCheckConfig := app.HealthCheckConfig{
		Addr: ":0", // Use port 0 to let the OS pick an available port
		Path: "/healthcheck",
	}

	t.Run("initialize wires all components without error", func(t *testing.T) {
		relayer := app.New(natsConn, pool, publicClient, privateClient, keyQueues)

		err := relayer.Initialize(contractsConfig, listenersConfig, executorConfig, servicesConfig, healthCheckConfig)
		require.NoError(t, err)
	})

	t.Run("run starts and shutdown completes gracefully", func(t *testing.T) {
		relayer := app.New(natsConn, pool, publicClient, privateClient, keyQueues)

		err := relayer.Initialize(contractsConfig, listenersConfig, executorConfig, servicesConfig, healthCheckConfig)
		require.NoError(t, err)

		runErr := make(chan error, 1)
		go func() {
			runErr <- relayer.Run()
		}()

		// Give goroutines time to start
		time.Sleep(500 * time.Millisecond)

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		err = relayer.Shutdown(shutdownCtx)
		require.NoError(t, err)

		select {
		case err := <-runErr:
			// Run() may return context.Canceled when Shutdown() cancels the context
			// during WaitForAuthorization (expected with empty key strings in tests).
			if err != nil {
				require.ErrorIs(t, err, context.Canceled, "Run should only fail with context.Canceled during shutdown")
			}
		case <-time.After(15 * time.Second):
			t.Fatal("Run did not exit within 15 seconds after Shutdown")
		}
	})
}
