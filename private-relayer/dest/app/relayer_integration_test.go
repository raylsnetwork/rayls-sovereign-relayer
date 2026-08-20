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

	"github.com/raylsnetwork/rayls-privacy-relayer-api/keyqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/config"
	destapp "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/app"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/keymanager"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
)

var (
	hubRegistryAddr  = common.HexToAddress("0x1000000000000000000000000000000000000001")
	nodeRegistryAddr = common.HexToAddress("0x2000000000000000000000000000000000000001")

	// Hub contracts (Private Hub)
	participantStorageAddr = common.HexToAddress("0xA001")
	resourceRegistryAddr   = common.HexToAddress("0xA002")
	endpointAddr           = common.HexToAddress("0xA003")
	teleportAddr           = common.HexToAddress("0xA004")
	enygmaTeleportAddr     = common.HexToAddress("0xA005")
	dvpTeleportAddr        = common.HexToAddress("0xA006")
	dvpAddr                = common.HexToAddress("0xA007")
	auditManagerAddr       = common.HexToAddress("0xA008")

	// Node contracts (Private Ledger)
	nodeEndpointAddr = common.HexToAddress("0xB001")
)

var dbConf = testtools.DBConfig{
	User:           "test",
	Pass:           "test",
	Database:       "testdb",
	MigrationsPath: "file://../repository/migrations",
}

func TestDestPrivateRelayer_Integration(t *testing.T) {
	// Set working directory to repo root so txsim.PopulateErrorMap("./contracts/") works.
	t.Chdir("../../..")

	testtools.SilenceLogger()

	// Start infrastructure containers
	pool, pgCleanup := testtools.SetupPostgres(t, dbConf)
	defer pgCleanup()

	natsConn, natsCleanup := testtools.SetupNATS(t)
	defer natsCleanup()

	// Set up fake Ethereum JSON-RPC servers
	hubChainID := big.NewInt(100)
	nodeChainID := big.NewInt(200)
	venChainID := big.NewInt(300)

	hubServer := testtools.NewFakeEthServer(t, testtools.FakeEthServerConfig{
		ChainID:  hubChainID,
		GasPrice: big.NewInt(1_000_000_000),
		DeploymentRegistries: map[common.Address]testtools.DeploymentRegistry{
			hubRegistryAddr: {
				Names: []string{
					"ParticipantStorage", "ResourceRegistry", "Endpoint",
					"Teleport", "EnygmaTeleport", "DvpTeleport", "Dvp", "AuditManager",
				},
				Addresses: []common.Address{
					participantStorageAddr, resourceRegistryAddr, endpointAddr,
					teleportAddr, enygmaTeleportAddr, dvpTeleportAddr, dvpAddr, auditManagerAddr,
				},
			},
		},
	})

	nodeServer := testtools.NewFakeEthServer(t, testtools.FakeEthServerConfig{
		ChainID:  nodeChainID,
		GasPrice: big.NewInt(1_000_000_000),
		DeploymentRegistries: map[common.Address]testtools.DeploymentRegistry{
			nodeRegistryAddr: {
				Names:     []string{"Endpoint"},
				Addresses: []common.Address{nodeEndpointAddr},
			},
		},
	})

	// Create ethclient.Client instances pointing at fake servers
	hubRPC, err := rpc.Dial(hubServer.URL)
	require.NoError(t, err)
	t.Cleanup(func() { hubRPC.Close() })
	hubClient := ethclient.NewClient(hubRPC)

	nodeRPC, err := rpc.Dial(nodeServer.URL)
	require.NoError(t, err)
	t.Cleanup(func() { nodeRPC.Close() })
	nodeClient := ethclient.NewClient(nodeRPC)

	// Key queues with buffer capacity (no real keys needed for init)
	keyQueues := &keymanager.KeyQueues{
		PrivateHub:  keyqueue.New(3),
		PrivateNode: keyqueue.New(3),
		DvpOperator: keyqueue.New(3),
	}

	conf := config.Config{
		PrivateHubChainID:  hubChainID,
		PrivateNodeChainID: nodeChainID,
		VENChainID:         venChainID,

		PrivateHubDeploymentProxyRegistry:  hubRegistryAddr,
		PrivateNodeDeploymentProxyRegistry: nodeRegistryAddr,

		CTSRootURL: "http://localhost:9999",
		CTSAPIKey:  "test-key",
		CTSSecret:  "test-secret",
		ProofAPIURL:  "http://localhost:9998",

		ListenersBlockBatchSize:  100,
		PrivateHubStartingBlock:  big.NewInt(0),
		PrivateNodeStartingBlock: big.NewInt(0),

		ExecutorBatchMessages: 10,
		ExpirationTime:        5 * time.Minute,
		DefaultContextTimeout: 2 * time.Minute,

		NumberOfJSParamsIn:           2,
		PrivateHubDvpMerkleTreeDepth: 8,
	}

	t.Run("initialize wires all components without error", func(t *testing.T) {
		relayer := destapp.New(natsConn, pool, hubClient, nodeClient, keyQueues, nil)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		err := relayer.Initialize(ctx, conf)
		require.NoError(t, err)
	})

	t.Run("run starts and shutdown completes gracefully", func(t *testing.T) {
		relayer := destapp.New(natsConn, pool, hubClient, nodeClient, keyQueues, nil)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		err := relayer.Initialize(ctx, conf)
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
			require.NoError(t, err, "Run should exit cleanly after Shutdown")
		case <-time.After(15 * time.Second):
			t.Fatal("Run did not exit within 15 seconds after Shutdown")
		}
	})
}
