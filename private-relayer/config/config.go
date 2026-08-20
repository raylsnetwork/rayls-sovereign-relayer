package config

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Config represents the private relayer application configuration.
type Config struct {
	// Database configuration
	DatabaseConnectionString string `mapstructure:"PRIVATE_RELAYER_DATABASE_CONNECTIONSTRING" validate:"required,url"`

	NATSUrl string `mapstructure:"NATS_URL" validate:"required"`

	// NATS mTLS — reuses this binary's client cert (also used for the
	// CTS gRPC channel) since the relayer presents the same identity
	// on both connections.
	NATSTLSCAFile   string `mapstructure:"NATS_TLS_CA_FILE" validate:"required" default:"/app/private-relayer/certs/ca.crt"`
	NATSTLSCertFile string `mapstructure:"NATS_TLS_CERT_FILE" validate:"required" default:"/app/private-relayer/certs/private-relayer.crt"`
	NATSTLSKeyFile  string `mapstructure:"NATS_TLS_KEY_FILE" validate:"required" default:"/app/private-relayer/certs/private-relayer.key"`

	// Logger
	LogLevel   string `mapstructure:"LOG_LEVEL" validate:"required" default:"INFO"`
	LogHandler string `mapstructure:"LOG_HANDLER" validate:"required" default:"Text"`

	// OpenTelemetry
	OtelSDKDisabled bool   `mapstructure:"OTEL_SDK_DISABLED" default:"true"`
	OtelServiceName string `mapstructure:"OTEL_SERVICE_NAME"`

	// Private Hub configuration
	PrivateHubChainID  *big.Int `mapstructure:"PNH_CHAIN_ID" validate:"required"`
	PrivateHubURL      string   `mapstructure:"PNH_RPC_URL" validate:"required,url"`
	PrivateNodeChainID *big.Int `mapstructure:"PRIVACY_NODE_CHAIN_ID" validate:"required"`
	PrivateNodeURL     string   `mapstructure:"PRIVACY_NODE_RPC_URL" validate:"required,url"`
	VENChainID         *big.Int `mapstructure:"PNH_OPERATOR_CHAIN_ID" validate:"required"`

	// DeploymentProxyRegistry addresses - used to discover contract addresses
	PrivateHubDeploymentProxyRegistry  common.Address `mapstructure:"PNH_DEPLOYMENT_PROXY_REGISTRY" validate:"required"`
	PrivateNodeDeploymentProxyRegistry common.Address `mapstructure:"PRIVACY_NODE_DEPLOYMENT_PROXY_REGISTRY" validate:"required"`

	// Atomic service configuration.
	// Decommissioning Teleport (vanilla, atomic): atomic timeout-revert config.
	ExpirationTime        time.Duration `mapstructure:"PNH_EXPIRATION_REVERT_TIME_IN_MINUTES" validate:"required"`
	ExecutorBatchMessages int           `mapstructure:"PRIVACY_NODE_EXECUTOR_BATCH_MESSAGES" validate:"required" default:"800"`

	// Listeners configuration
	ListenersBlockBatchSize  int      `mapstructure:"BLOCKCHAIN_LISTENER_BATCH_BLOCKS" validate:"min=1" default:"100"`
	PrivateHubStartingBlock  *big.Int `mapstructure:"PNH_CHAIN_STARTING_BLOCK"`
	PrivateNodeStartingBlock *big.Int `mapstructure:"PRIVACY_NODE_STARTING_BLOCK"`

	// CTS configuration
	CTSRootURL string `mapstructure:"CTS_GRPC_URL" validate:"required"`
	CTSAPIKey  string `mapstructure:"CTS_API_KEY" validate:"required"`
	CTSSecret  string `mapstructure:"CTS_SECRET" validate:"required"`

	// mTLS for the CTS gRPC channel. All three files are required —
	// the channel always runs mTLS.
	CTSTLSCAFile   string `mapstructure:"CTS_CLIENT_TLS_CA_FILE" validate:"required" default:"/app/private-relayer/certs/ca.crt"`
	CTSTLSCertFile string `mapstructure:"CTS_CLIENT_TLS_CERT_FILE" validate:"required" default:"/app/private-relayer/certs/private-relayer.crt"`
	CTSTLSKeyFile  string `mapstructure:"CTS_CLIENT_TLS_KEY_FILE" validate:"required" default:"/app/private-relayer/certs/private-relayer.key"`

	// Proof API configuration
	ProofAPIURL string `mapstructure:"PRIVACY_NODE_ENYGMA_PROOF_API_ADDRESS" validate:"required,url"`

	// Dvp configuration
	PrivateHubDvpMerkleTreeDepth int `mapstructure:"DVP_MERKLE_TREE_DEPTH" validate:"min=1"  default:"8"`
	NumberOfJSParamsIn           int `mapstructure:"NUMBER_OF_JS_PARAMS_IN" validate:"required"`

	// Enygma configuration
	EnygmaBatchSize                int           `mapstructure:"PRIVACY_NODE_EXECUTOR_ENYGMA_BATCH_MESSAGES" validate:"required" default:"1000"`
	EnygmaMaxConcurrentResourceIDs int           `mapstructure:"PRIVACY_NODE_EXECUTOR_ENYGMA_MAX_CONCURRENT_RESOURCE_IDS" validate:"required,min=1" default:"2"`
	DefaultContextTimeout          time.Duration `default:"2m"`

	HealthcheckPort string `mapstructure:"PRIVATE_RELAYER_HEALTHCHECK_PORT"`

	// Fault injection (testing only — must never be enabled in production)
	FaultInjectionEnabled     bool   `mapstructure:"FAULT_INJECTION_ENABLED" default:"false"`
	FaultInjectionPort        string `mapstructure:"FAULT_INJECTION_PORT" default:"9999"`
	FaultInjectionPersistPath string `mapstructure:"FAULT_INJECTION_PERSIST_PATH"`
}
