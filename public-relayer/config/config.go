// Decommissioning Teleport (vanilla, atomic).

// Package config implements the legacy public-chain (RN) Teleport bridge relayer.
//
// Deprecated: Decommissioning Teleport (vanilla, atomic).
package config

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Config represents the pubrelayer application configuration.
type Config struct {
	// Chain ID (used for NATS subject isolation)
	PrivateNodeChainID *big.Int `mapstructure:"PRIVACY_NODE_CHAIN_ID" validate:"required"`

	// Database configuration
	DatabaseConnectionString string `mapstructure:"RAYLS_NODE_DATABASE_CONNECTIONSTRING" validate:"required,url"`

	// Logger
	LogLevel   string `mapstructure:"LOG_LEVEL" validate:"required" default:"INFO"`
	LogHandler string `mapstructure:"LOG_HANDLER" validate:"required" default:"Text"`

	// OpenTelemetry
	OtelSDKDisabled bool   `mapstructure:"OTEL_SDK_DISABLED" default:"true"`
	OtelServiceName string `mapstructure:"OTEL_SERVICE_NAME"`

	// Ledger configuration
	PublicChainURL string `mapstructure:"PUBLIC_CHAIN_RPC_URL" validate:"required,url"`
	PrivateNodeURL string `mapstructure:"PRIVACY_NODE_RPC_URL" validate:"required,url"`

	// DeploymentProxyRegistry addresses - used to discover contract addresses
	PublicChainDeploymentProxyRegistry common.Address `mapstructure:"PUBLIC_CHAIN_DEPLOYMENT_PROXY_REGISTRY" validate:"required"`
	PrivateNodeDeploymentProxyRegistry common.Address `mapstructure:"PRIVACY_NODE_DEPLOYMENT_PROXY_REGISTRY" validate:"required"`

	// NATS configuration
	NATSUrl string `mapstructure:"NATS_URL" validate:"required"`

	// NATS mTLS — reuses this binary's client cert (also used for the
	// CTS gRPC channel) since the relayer presents the same identity
	// on both connections.
	NATSTLSCAFile   string `mapstructure:"NATS_TLS_CA_FILE" validate:"required" default:"/app/public-relayer/certs/ca.crt"`
	NATSTLSCertFile string `mapstructure:"NATS_TLS_CERT_FILE" validate:"required" default:"/app/public-relayer/certs/public-relayer.crt"`
	NATSTLSKeyFile  string `mapstructure:"NATS_TLS_KEY_FILE" validate:"required" default:"/app/public-relayer/certs/public-relayer.key"`

	// Listeners configuration
	ListenersBlockBatchSize  int      `mapstructure:"BLOCKCHAIN_LISTENER_BATCH_BLOCKS" validate:"min=1" default:"100"`
	PublicChainStartingBlock *big.Int `mapstructure:"PUBLIC_CHAIN_STARTING_BLOCK"` // todo: fetch latest block during deployment and copy
	PrivateNodeStartingBlock *big.Int `mapstructure:"PRIVACY_NODE_STARTING_BLOCK"`

	// CTS configuration
	CTSRootURL string `mapstructure:"CTS_GRPC_URL"`
	CTSAPIKey  string `mapstructure:"CTS_API_KEY"`
	CTSSecret  string `mapstructure:"CTS_SECRET"`

	// mTLS for the CTS gRPC channel. All three files are required —
	// the channel always runs mTLS.
	CTSTLSCAFile   string `mapstructure:"CTS_CLIENT_TLS_CA_FILE" validate:"required" default:"/app/public-relayer/certs/ca.crt"`
	CTSTLSCertFile string `mapstructure:"CTS_CLIENT_TLS_CERT_FILE" validate:"required" default:"/app/public-relayer/certs/public-relayer.crt"`
	CTSTLSKeyFile  string `mapstructure:"CTS_CLIENT_TLS_KEY_FILE" validate:"required" default:"/app/public-relayer/certs/public-relayer.key"`

	// Fault injection (testing only — must never be enabled in production).
	// Port and persist path are public-relayer-specific so CTS, private-relayer
	// and public-relayer instances sharing the same per-participant .env file
	// each get an isolated HTTP server and on-disk persist file. The shared
	// FAULT_INJECTION_ENABLED flag acts as a single master switch across all
	// three services; binaries built without the faultinjection tag treat it
	// as a no-op.
	FaultInjectionEnabled     bool   `mapstructure:"FAULT_INJECTION_ENABLED" default:"false"`
	FaultInjectionPort        string `mapstructure:"PUBLIC_RELAYER_FAULT_INJECTION_PORT" default:"9999"`
	FaultInjectionPersistPath string `mapstructure:"PUBLIC_RELAYER_FAULT_INJECTION_PERSIST_PATH"`
}
