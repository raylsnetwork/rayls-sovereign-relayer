package run

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/configinit"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/batcher"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/config"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/crypto"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/etherror"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/ethrpc"
	ctsgrpc "github.com/raylsnetwork/rayls-privacy-relayer-api/cts/grpc"
	ctshttp "github.com/raylsnetwork/rayls-privacy-relayer-api/cts/http"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/keymanager"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/repo"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/faultinjector"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/keyqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/logger"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/mtls"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/noncecache"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/otel"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

const encryptorInitTimeout = 5 * time.Second

func RunServer(path string) error {
	conf, err := configinit.InitConfig[config.Config](path)
	if err != nil {
		slog.Error("Failed to initialize config", slog.Any("error", err))
		return fmt.Errorf("initializing CTS config: %w", err)
	}
	slog.Info("Configuration validated successful")

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

	// ********************************************************** //
	// ************** Fault Injection Initialization ************* //
	// Binaries built without the faultinjection tag get a no-op
	// NewHTTPServer and the whole block becomes a runtime no-op.
	// The server is torn down via defer so it stops even when later
	// startup steps return with an error.
	var faultServer *http.Server
	if conf.FaultInjectionEnabled {
		if conf.FaultInjectionPersistPath != "" {
			faultinjector.SetPersistPath(conf.FaultInjectionPersistPath)
		}
		faultinjector.Enable()
		faultServer = faultinjector.NewHTTPServer(":" + conf.FaultInjectionPort)
		if faultServer != nil {
			slog.Warn("Starting fault injection API server", slog.String("addr", ":"+conf.FaultInjectionPort))
			go func() {
				if srvErr := faultServer.ListenAndServe(); srvErr != nil && !errors.Is(srvErr, http.ErrServerClosed) {
					slog.Error("Fault injection server error", slog.Any("error", srvErr))
				}
			}()
			defer func() {
				const faultInjectionShutdownTimeout = 5 * time.Second
				fiCtx, fiCancel := context.WithTimeout(context.Background(), faultInjectionShutdownTimeout)
				defer fiCancel()
				if shutdownErr := faultServer.Shutdown(fiCtx); shutdownErr != nil {
					slog.Error("Error while shutting down fault injection server", slog.Any("error", shutdownErr))
				}
			}()
		} else {
			slog.Warn("Fault injection enabled in config but binary built without faultinjection tag — no-op")
		}
	}
	// ************** Fault Injection Initialization End ********** //
	// ************************************************************ //

	db, err := repo.NewDatabaseConnection(conf.DatabaseConn)
	if err != nil {
		slog.Error("Failed to connect to the database", slog.Any("error", err))
		return withstack.Wrap(fmt.Errorf("connecting to CTS database: %w", err))
	}
	defer db.Close()

	// Run migrations
	slog.Info("Running database migrations")
	if err := db.Migrate(); err != nil {
		slog.Error("Failed to run migrations", slog.Any("error", err))
		return fmt.Errorf("running CTS database migrations: %w", err)
	}

	plaintextEncryptorFactory := crypto.NewPlaintextEncryptorFactory()
	awsEncryptorFactory := crypto.NewAWSEncryptorFactory(conf.AWSProfile, conf.AWSAlias)
	gcpEncryptorFactory := crypto.NewGCPEncryptorFactory(
		conf.GCPProject,
		conf.GCPLocation,
		conf.GCPKeyRing,
		conf.GCPCryptoKey,
	)

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), encryptorInitTimeout)
	defer cancel()

	encryptorClient, cleanup, err := crypto.GetEncryptor(
		ctxWithTimeout,
		conf.EncryptorService,
		plaintextEncryptorFactory,
		awsEncryptorFactory,
		gcpEncryptorFactory,
	)
	if err != nil {
		slog.Error("Failed to create encryptor client", slog.Any("details", err))
		return fmt.Errorf("creating encryptor client: %w", err)
	}
	defer func() { _ = cleanup() }()

	raylsViewKeysRepo := db.NewRaylsViewKeysRepository(
		prefixedTableName(conf.EncryptorService, repo.RaylsViewKeysCollectionName),
	)

	ecdsaKeysRepo := db.NewRaylsSignKeysRepository(
		prefixedTableName(conf.EncryptorService, repo.RaylsSignKeysCollectionName),
	)

	paymentSpendKeysRepo := db.NewPaymentSpendKeysRepository(
		prefixedTableName(conf.EncryptorService, repo.PaymentSpendKeysCollectionName),
	)

	sharedSecretsRepo := db.NewSharedSecretsRepository(
		prefixedTableName(conf.EncryptorService, repo.SharedSecretsCollectionName),
	)

	enygmaSelfSecretsRepo := db.NewEnygmaSelfSecretsRepository(
		prefixedTableName(conf.EncryptorService, repo.EnygmaSelfSecretsCollectionName),
	)

	keysService := service.NewKeysService(
		encryptorClient,
		raylsViewKeysRepo,
		ecdsaKeysRepo,
		sharedSecretsRepo,
		enygmaSelfSecretsRepo,
		paymentSpendKeysRepo,
	)
	encryptService := service.NewEncryptService(conf.ChainID, keysService, encryptorClient)

	// Start HTTP server early so /health and /public/addresses are
	// reachable while the key manager waits for on-chain authorization.
	// The contracts deployment service polls /public/addresses to learn
	// which addresses to authorize — if this server isn't up, the
	// authorization can never happen (chicken-and-egg).
	httpAddr := ":" + conf.HTTPPort
	httpSrv, err := ctshttp.NewServer(httpAddr, keysService)
	if err != nil {
		return fmt.Errorf("creating HTTP server: %w", err)
	}
	httpErrCh := make(chan error, 1)
	go func() {
		slog.Info("Running CTS HTTP server...", slog.String("addr", httpAddr))
		if err := httpSrv.Serve(); err != nil {
			httpErrCh <- err
		}
	}()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpSrv.Shutdown(shutdownCtx)
	}()

	encryptHandler := ctsgrpc.NewEncryptHandler(&encryptService)
	keysHandler := ctsgrpc.NewKeysHandler(keysService)

	var hubClient *ethclient.Client
	if conf.PrivateHubEnabled() {
		hubClient, err = ethclient.Dial(conf.PrivateHubURL)
		if err != nil {
			return withstack.Wrap(fmt.Errorf("failed to connect to private hub ledger: %w", err))
		}
		defer hubClient.Close()
	} else {
		slog.Info("private hub disabled — skipping hub ethclient, key authorization wait, ledger probe, participant registrar, hub + dvpoperator TxOps handlers, and batch pipelines")
	}

	nodeClient, err := ethclient.Dial(conf.PrivateNodeURL)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to connect to private node ledger: %w", err))
	}
	defer nodeClient.Close()

	var publicClient *ethclient.Client
	if conf.PublicChainEnabled() {
		publicClient, err = ethclient.Dial(conf.PublicChainURL)
		if err != nil {
			return withstack.Wrap(fmt.Errorf("failed to connect to public ledger: %w", err))
		}
		defer publicClient.Close()
	} else {
		slog.Info("public chain disabled — skipping public-chain ethclient, key authorization wait, TxOps handler, and batch pipeline")
	}

	// Now run the key manager which blocks until on-chain authorization
	// completes. The HTTP and gRPC servers are already accepting
	// connections so the contracts deployer can authorize addresses and
	// relayers won't get connection-refused errors.
	signKeyMgr := keymanager.NewSignKeyManager(
		keysService,
		keymanager.ChainConfig{Backend: hubClient, DeploymentProxyRegistry: conf.PrivateHubDeploymentProxyRegistry},
		keymanager.ChainConfig{Backend: nodeClient, DeploymentProxyRegistry: conf.PrivateNodeDeploymentProxyRegistry},
		keymanager.ChainConfig{Backend: publicClient, DeploymentProxyRegistry: conf.PublicChainDeploymentProxyRegistry},
		keymanager.WithPrivateHub(conf.PrivateHubEnabled()),
		keymanager.WithPublicChain(conf.PublicChainEnabled()),
	)
	privateRelayerKeys, publicRelayerKeys, err := signKeyMgr.Run(context.Background())
	if err != nil {
		return withstack.Wrap(fmt.Errorf("ensuring and authorizing relayer sign keys: %w", err))
	}

	// Create empty key queues and start the gRPC server BEFORE the key
	// manager blocks on authorization. This way relayers that depend on
	// the CTS gRPC port can connect immediately. Any TxOps call that
	// arrives before keys are populated will block on Dequeue until keys
	// are ready — no connection-refused errors.
	nodeQueue := keyqueue.New(5)
	var hubQueue *keyqueue.RaylsSignPrivateKeyManager
	var operatorQueue *keyqueue.RaylsSignPrivateKeyManager
	if conf.PrivateHubEnabled() {
		hubQueue = keyqueue.New(5)
		operatorQueue = keyqueue.New(1)
	}
	var privateQueue *keyqueue.RaylsSignPrivateKeyManager
	var publicQueue *keyqueue.RaylsSignPrivateKeyManager
	if conf.PublicChainEnabled() {
		privateQueue = keyqueue.New(5)
		publicQueue = keyqueue.New(5)
	}

	// Populate the key queues now that authorization is complete. TxOps
	// handlers that were blocking on Dequeue will unblock.
	populateKeyQueue(nodeQueue, privateRelayerKeys.PrivateNodeKeys)
	if conf.PrivateHubEnabled() {
		populateKeyQueue(hubQueue, privateRelayerKeys.PrivateHubKeys)
		populateKeyQueue(operatorQueue, privateRelayerKeys.PrivateHubDvpOperatorKeys)
	}
	if conf.PublicChainEnabled() {
		populateKeyQueue(privateQueue, publicRelayerKeys.PrivateChainKeys)
		populateKeyQueue(publicQueue, publicRelayerKeys.PublicChainKeys)
	}

	// Verify the node's error wire strings against our constants before serving.
	probeCtx, probeCancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer probeCancel()
	type ledgerProbe struct {
		name   string
		client *ethclient.Client
		key    *ecdsa.PrivateKey
	}
	probes := []ledgerProbe{
		{"private-node", nodeClient, firstKey(privateRelayerKeys.PrivateNodeKeys)},
	}
	if conf.PrivateHubEnabled() {
		probes = append(probes, ledgerProbe{"private-hub", hubClient, firstKey(privateRelayerKeys.PrivateHubKeys)})
	}
	if conf.PublicChainEnabled() {
		probes = append(probes, ledgerProbe{"publicchain", publicClient, firstKey(publicRelayerKeys.PublicChainKeys)})
	}
	for _, p := range probes {
		if p.key == nil {
			return withstack.Wrap(fmt.Errorf("verifying %s ledger error strings: no signing key available", p.name))
		}
		if err := etherror.VerifyLedgerErrorStrings(probeCtx, p.client, p.key); err != nil {
			return withstack.Wrap(fmt.Errorf("verifying %s ledger error strings: %w", p.name, err))
		}
	}

	// buildTxOpsService touches each identity's RPC endpoint (ChainID/gasPrice in
	// both NewAuthGenWithCache and NewBatcherWithCache), now retried over the
	// flaky VPN-routed link. Give the shared construction budget headroom for
	// those retries across all identities (construction-only; the gRPC handlers
	// it builds run under their own contexts).
	txOpsCtx, txOpsCancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer txOpsCancel()

	// Shared sync-tx repo backing the idempotency path of every TxOps service.
	// cts_sync_tx is a single table keyed by the per-request idempotency id, so
	// one repo instance is shared across all signing identities.
	txSyncRepo := repo.NewCTSSyncTxRepository(db.Pool())

	// One shared nonce cache per (ethClient, keyQueue) pair. AuthGen
	// (sync gRPC path) and the RPCBatcher used by both the gRPC handler
	// AND the async BatchPipeline all read/advance through the SAME
	// cache, so they cannot race each other into the same nonce slot
	// while a prior tx is still pending.
	// Identity name only used for log scoping if we ever surface cache
	// metrics; the cache itself is keyed by address.
	nodeNonceCache := noncecache.New(nodeClient)
	var hubNonceCache, operatorNonceCache *noncecache.Cache
	if conf.PrivateHubEnabled() {
		hubNonceCache = noncecache.New(hubClient)
		operatorNonceCache = noncecache.New(hubClient)
	}
	var privateNonceCache, publicNonceCache *noncecache.Cache
	if conf.PublicChainEnabled() {
		privateNonceCache = noncecache.New(nodeClient)
		publicNonceCache = noncecache.New(publicClient)
	}

	waitMinedTimeout := time.Duration(conf.WaitMinedTimeoutSecs) * time.Second

	nodeTxOps, err := buildTxOpsService(txOpsCtx, nodeClient, nodeQueue, txSyncRepo, nodeNonceCache, waitMinedTimeout)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("building private node TxOps service: %w", err))
	}
	var hubTxOpsHandler *ctsgrpc.TxOpsHandler
	var operatorTxOpsHandler *ctsgrpc.TxOpsHandler
	if conf.PrivateHubEnabled() {
		hubTxOps, err := buildTxOpsService(txOpsCtx, hubClient, hubQueue, txSyncRepo, hubNonceCache, waitMinedTimeout)
		if err != nil {
			return withstack.Wrap(fmt.Errorf("building private hub TxOps service: %w", err))
		}
		hubTxOpsHandler = ctsgrpc.NewTxOpsHandler(hubTxOps)

		operatorTxOps, err := buildTxOpsService(txOpsCtx, hubClient, operatorQueue, txSyncRepo, operatorNonceCache, waitMinedTimeout)
		if err != nil {
			return withstack.Wrap(fmt.Errorf("building dvp operator TxOps service: %w", err))
		}
		operatorTxOpsHandler = ctsgrpc.NewTxOpsHandler(operatorTxOps)
	}
	var privateTxOpsHandler *ctsgrpc.TxOpsHandler
	var publicTxOpsHandler *ctsgrpc.TxOpsHandler
	if conf.PublicChainEnabled() {
		privateTxOps, err := buildTxOpsService(txOpsCtx, nodeClient, privateQueue, txSyncRepo, privateNonceCache, waitMinedTimeout)
		if err != nil {
			return withstack.Wrap(fmt.Errorf("building private chain TxOps service: %w", err))
		}
		privateTxOpsHandler = ctsgrpc.NewTxOpsHandler(privateTxOps)

		publicTxOps, err := buildTxOpsService(txOpsCtx, publicClient, publicQueue, txSyncRepo, publicNonceCache, waitMinedTimeout)
		if err != nil {
			return withstack.Wrap(fmt.Errorf("building public chain TxOps service: %w", err))
		}
		publicTxOpsHandler = ctsgrpc.NewTxOpsHandler(publicTxOps)
	}

	grpcAddr := ":" + conf.GRPCPort
	grpcSrv, err := ctsgrpc.NewServer(
		grpcAddr,
		ctsgrpc.TLSConfig{
			CAFile:   conf.TLSCAFile,
			CertFile: conf.TLSCertFile,
			KeyFile:  conf.TLSKeyFile,
		},
		encryptHandler,
		keysHandler,

		hubTxOpsHandler,
		ctsgrpc.NewTxOpsHandler(nodeTxOps),
		operatorTxOpsHandler,

		privateTxOpsHandler,
		publicTxOpsHandler,
	)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("creating gRPC server: %w", err))
	}

	grpcErrCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	go func() {
		slog.Info("Running CTS gRPC server...", slog.String("addr", grpcAddr))
		if err := grpcSrv.Serve(); err != nil {
			grpcErrCh <- err
		}
	}()

	if conf.PrivateHubEnabled() {
		participantRegistrar := keymanager.NewParticipantRegistrar(
			conf.ChainID,
			conf.VenChainID,
			keymanager.ChainConfig{Backend: hubClient, DeploymentProxyRegistry: conf.PrivateHubDeploymentProxyRegistry},
			keysService,
			&encryptService,
			privateRelayerKeys.PrivateHubKeys,
		)
		if err := participantRegistrar.Run(context.Background()); err != nil {
			return withstack.Wrap(fmt.Errorf("registering relayer with participant storage: %w", err))
		}
	}

	// ***** Async batch pipeline: NATS + shared repo + per-identity pipelines *****
	natsTLS, err := mtls.LoadClientConfig(conf.NATSTLSCAFile, conf.NATSTLSCertFile, conf.NATSTLSKeyFile)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("loading NATS client TLS config: %w", err))
	}
	nc, err := nats.Connect(conf.NATSUrl, nats.Secure(natsTLS))
	if err != nil {
		return withstack.Wrap(fmt.Errorf("connecting to NATS: %w", err))
	}
	defer func() { _ = nc.Drain() }()

	js, err := jetstream.New(nc)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("creating jetstream: %w", err))
	}

	managerCtx, managerCancel := context.WithTimeout(context.Background(), 30*time.Second)
	manager, err := msgqueue.NewManager(managerCtx, js, conf.ChainID.String())
	managerCancel()
	if err != nil {
		return withstack.Wrap(fmt.Errorf("creating message queue manager: %w", err))
	}

	// Shared repo — the cts_transaction table is scoped by the `identity`
	// column, not by separate tables, so one instance serves all pipelines.
	ctsTransactionRepo := repo.NewCTSTransactionRepository(db.Pool())

	// Constructing each pipeline touches its RPC endpoint (ChainID/gasPrice).
	// Over VPN-routed remote-dev links those calls can be briefly slow or
	// flaky; NewBatcherWithCache retries them, so give the shared construction
	// budget headroom for retries across all pipelines instead of failing CTS
	// startup on a single transient blip.
	pipelineCtx, pipelineCancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer pipelineCancel()

	type pipelineSpec struct {
		identity   string
		queue      *keyqueue.RaylsSignPrivateKeyManager
		ethClient  *ethclient.Client
		nonceCache *noncecache.Cache
	}
	pipelineSpecs := []pipelineSpec{
		{"privatenode", nodeQueue, nodeClient, nodeNonceCache},
	}
	if conf.PrivateHubEnabled() {
		pipelineSpecs = append(pipelineSpecs,
			pipelineSpec{"privatehub", hubQueue, hubClient, hubNonceCache},
			pipelineSpec{"dvpoperator", operatorQueue, hubClient, operatorNonceCache},
		)
	}
	if conf.PublicChainEnabled() {
		pipelineSpecs = append(pipelineSpecs,
			pipelineSpec{"privatechain", privateQueue, nodeClient, privateNonceCache},
			pipelineSpec{"publicchain", publicQueue, publicClient, publicNonceCache},
		)
	}

	pipelines := make([]*BatchPipeline, 0, len(pipelineSpecs))
	for _, spec := range pipelineSpecs {
		p, err := NewBatchPipeline(
			pipelineCtx, spec.identity, manager, spec.queue, spec.ethClient, ctsTransactionRepo, spec.nonceCache,
		)
		if err != nil {
			return withstack.Wrap(fmt.Errorf("creating batch pipeline %s: %w", spec.identity, err))
		}
		pipelines = append(pipelines, p)
	}

	for _, p := range pipelines {
		p.Run()
	}
	slog.Info("Batch pipelines started", slog.Int("count", len(pipelines)))

	// Keys authorized, queues populated, gRPC serving, participant
	// registered — mark CTS as fully ready so /ready returns 200 and
	// docker-compose dependents (relayers) can start.
	httpSrv.MarkReady()
	slog.Info("CTS is ready")

	select {
	case sig := <-sigCh:
		slog.Info("Shutdown signal received", slog.String("signal", sig.String()))
	case err := <-grpcErrCh:
		return withstack.Wrap(fmt.Errorf("gRPC server: %w", err))
	case err := <-httpErrCh:
		return withstack.Wrap(fmt.Errorf("HTTP server: %w", err))
	}

	// ***** Shutdown: stop pipelines first (they hold keys + NATS subs),
	// then the gRPC server. 10s deadline per pipeline is plenty for
	// four tick loops to notice ctx.Done and return. *****
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	for _, p := range pipelines {
		if err := p.Stop(shutdownCtx); err != nil {
			slog.Error("batch pipeline stop exceeded deadline", slog.Any("err", err))
		}
	}

	grpcSrv.GracefulStop()

	return nil
}

func prefixedTableName(encrName, tableName string) string {
	return encrName + "_" + tableName
}

// buildTxOpsService wires a sync gRPC TxOps handler against a pre-built key
// queue and an externally-owned nonce cache. The queue is supplied by the
// caller so it can be shared with the async BatchPipeline for the same
// identity; the cache is similarly shared so AuthGen (sync gRPC) and the
// RPCBatcher driven by both this service AND the async pipeline coordinate
// nonce assignment through a single counter — closing the cross-path
// nonce-collision race.
func buildTxOpsService(
	ctx context.Context,
	ethClient *ethclient.Client,
	queue *keyqueue.RaylsSignPrivateKeyManager,
	txSyncRepo *repo.CTSSyncTxRepository,
	nonceCache *noncecache.Cache,
	waitMinedTimeout time.Duration,
) (*service.TxOpsService, error) {
	auth, err := contractclient.NewAuthGenWithCache(ctx, ethClient, nonceCache)
	if err != nil {
		return nil, fmt.Errorf("creating auth generator: %w", err)
	}
	ethBatcher, err := ethrpc.NewBatcherWithCache(ctx, queue, ethClient, ethClient.Client(), nonceCache)
	if err != nil {
		return nil, fmt.Errorf("creating ethereum batcher: %w", err)
	}
	return service.NewTxOpsServiceWithWaitMinedTimeout(auth, queue, ethClient, ethBatcher, txSyncRepo, waitMinedTimeout), nil
}

// populateKeyQueue enqueues decrypted private keys into a pre-existing
// queue. Called after signKeyMgr.Run returns — any TxOps handler that was
// blocking on Dequeue will unblock once keys appear.
// firstKey returns the first key in the list, or nil if it is empty. Used to
// pick an arbitrary signing key for the per-node startup error-string probe.
func firstKey(keys []*ecdsa.PrivateKey) *ecdsa.PrivateKey {
	if len(keys) == 0 {
		return nil
	}
	return keys[0]
}

func populateKeyQueue(q *keyqueue.RaylsSignPrivateKeyManager, keys []*ecdsa.PrivateKey) {
	for _, k := range keys {
		q.Enqueue(k)
	}
}

type BatchPipeline struct {
	ingester  *batcher.IngestionService
	sender    *batcher.BatcherService
	receipter *batcher.ReceipterService
	reaper    *batcher.ReaperService

	wg     sync.WaitGroup
	cancel context.CancelFunc
}

// NewBatchPipeline constructs a BatchPipeline for a single signing
// identity. The identity string is used both as the per-row `identity`
// column value in cts_transaction and as the NATS subject segment
// (`cts.send.<identity>` and `cts.result.<identity>`).
//
// ctsTransactionRepo is passed in so a single repo instance can be
// shared across every pipeline — scoping is handled by the identity
// column, not by separate tables.
func NewBatchPipeline(
	ctx context.Context,
	identity string,
	manager *msgqueue.Manager,
	keyQueue *keyqueue.RaylsSignPrivateKeyManager,
	ethClient *ethclient.Client,
	ctsTransactionRepo batcher.Repository,
	nonceCache *noncecache.Cache,
) (*BatchPipeline, error) {
	// Shares the nonce cache with the sync gRPC path (AuthGen + its
	// RPCBatcher built in buildTxOpsService) — see the cross-path
	// nonce-collision notes in that wiring.
	ethBatcher, err := ethrpc.NewBatcherWithCache(ctx, keyQueue, ethClient, ethClient.Client(), nonceCache)
	if err != nil {
		return nil, fmt.Errorf("creating ethereum batcher for %s: %w", identity, err)
	}

	sendSubject := "cts.send." + identity
	resultSubject := "cts.result." + identity

	sendConsumer, err := msgqueue.NewConsumer[types.TxRequest](
		ctx, manager, identity, sendSubject,
	)
	if err != nil {
		return nil, fmt.Errorf("creating send consumer for %s: %w", identity, err)
	}
	resultPublisher := msgqueue.NewPublisher[types.TxResult](manager, resultSubject)

	ingester := batcher.NewIngestionService(
		batcher.IngestionConfig{
			Identity:  identity,
			BatchSize: 100,
			Interval:  time.Second,
		},
		sendConsumer, ctsTransactionRepo,
	)
	sender := batcher.NewBatcherService(
		batcher.BatcherConfig{
			Identity:  identity,
			BatchSize: 100,
			Interval:  time.Second,
		},
		ctsTransactionRepo, ethBatcher, resultPublisher,
	)
	receipter := batcher.NewReceipterService(
		batcher.ReceipterConfig{
			Identity:  identity,
			BatchSize: 100,
			Interval:  time.Second,
		},
		ctsTransactionRepo, ethBatcher, resultPublisher,
	)
	reaper := batcher.NewReaperService(
		batcher.ReaperConfig{
			Identity:       identity,
			BatchSize:      100,
			Interval:       time.Minute,
			StuckThreshold: time.Minute,
			// Up to 5 broadcasts total (initial + 4 resends) before
			// dead-letter. With Interval=StuckThreshold=1m and sent_at
			// refreshed on each resend, that pushes the dead-letter
			// horizon to ~10 minutes from the first send, comfortably
			// beyond the worst-case worker/receipter blip we've seen
			// in practice.
			MaxResendAttempts: 5,
		},
		ctsTransactionRepo, ethBatcher, resultPublisher,
	)

	return &BatchPipeline{
		ingester:  ingester,
		sender:    sender,
		receipter: receipter,
		reaper:    reaper,
	}, nil
}

// Run starts the four pipeline services in their own goroutines and
// returns immediately. Stop is the single cancellation path — the
// internal context is not derived from any caller context.
func (p *BatchPipeline) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	p.wg.Add(4)
	go func() { defer p.wg.Done(); _ = p.ingester.Run(ctx) }()
	go func() { defer p.wg.Done(); _ = p.sender.Run(ctx) }()
	go func() { defer p.wg.Done(); _ = p.receipter.Run(ctx) }()
	go func() { defer p.wg.Done(); _ = p.reaper.Run(ctx) }()
}

// Stop cancels the pipeline's internal context and waits for every
// service goroutine to return. Accepts a shutdown context so the caller
// can bound the wait time — returns ctx.Err() if the shutdown deadline
// expires before services finish, nil on clean exit. Safe to call even
// if Run was never called.
func (p *BatchPipeline) Stop(ctx context.Context) error {
	if p.cancel == nil {
		return nil
	}
	p.cancel()

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
