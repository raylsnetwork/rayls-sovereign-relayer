package config

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Config represents the CTS (Cryptography Trust Suite) application configuration.
type Config struct {
	// Database
	DatabaseConn string `mapstructure:"CTS_DATABASE_CONNECTIONSTRING" validate:"required,url"`

	// Application
	CORSDomain string `mapstructure:"CTS_CORSDOMAIN" validate:"required"`
	ApiKey     string `mapstructure:"CTS_API_KEY" validate:"required"`
	Secret     string `mapstructure:"CTS_SECRET" validate:"required"`

	// GCP Encryptor
	GCPProject   string `mapstructure:"CTS_GCPPROJECT" validate:"required_if=EncryptorService gcp"`
	GCPLocation  string `mapstructure:"CTS_GCPLOCATION" validate:"required_if=EncryptorService gcp"`
	GCPKeyRing   string `mapstructure:"CTS_GCPKEYRING" validate:"required_if=EncryptorService gcp"`
	GCPCryptoKey string `mapstructure:"CTS_GCPCRYPTOKEY" validate:"required_if=EncryptorService gcp"`

	// AWS Encryptor
	AWSProfile string `mapstructure:"CTS_AWSPROFILE" validate:"required_if=EncryptorService aws"`
	AWSAlias   string `mapstructure:"CTS_AWSALIAS" validate:"required_if=EncryptorService aws"`

	// Blockchain
	ChainID    *big.Int `mapstructure:"BLOCKCHAIN_CHAIN_ID" validate:"required"`
	VenChainID *big.Int `mapstructure:"PNH_OPERATOR_CHAIN_ID"`

	// NATS — async batch pipeline transport between relayer and CTS.
	// Shared with the existing relayer env var name.
	NATSUrl string `mapstructure:"NATS_URL" validate:"required"`

	// NATS mTLS. The CTS process is a NATS *client*, so it presents
	// its own client cert (separate from the gRPC server cert above).
	NATSTLSCAFile   string `mapstructure:"NATS_TLS_CA_FILE" validate:"required" default:"/app/cts/certs/ca.crt"`
	NATSTLSCertFile string `mapstructure:"NATS_TLS_CERT_FILE" validate:"required" default:"/app/cts/certs/cts.crt"`
	NATSTLSKeyFile  string `mapstructure:"NATS_TLS_KEY_FILE" validate:"required" default:"/app/cts/certs/cts.key"`

	// TxOps targets — reuse the same env var names the relayer already
	// defines so operators don't have to set new variables during the
	// signing-migration phase.
	PrivateHubURL  string `mapstructure:"PNH_RPC_URL" validate:"omitempty,url"`
	PrivateNodeURL string `mapstructure:"PRIVACY_NODE_RPC_URL" validate:"required,url"`
	PublicChainURL string `mapstructure:"PUBLIC_CHAIN_RPC_URL" validate:"omitempty,url"`

	// DeploymentProxyRegistry addresses per chain — resolved at boot to look
	// up each chain's RelayAuthorizationRegistry. Env names match the
	// relayer configs so operators share one set of values across services.
	// PublicChainDeploymentProxyRegistry is optional: when empty, CTS treats
	// the public chain as not deployed (see PublicChainEnabled).
	// PrivateHubDeploymentProxyRegistry is likewise optional: when empty, CTS
	// treats the hub as not deployed (see PrivateHubEnabled).
	PrivateHubDeploymentProxyRegistry  common.Address `mapstructure:"PNH_DEPLOYMENT_PROXY_REGISTRY"`
	PrivateNodeDeploymentProxyRegistry common.Address `mapstructure:"PRIVACY_NODE_DEPLOYMENT_PROXY_REGISTRY" validate:"required"`
	PublicChainDeploymentProxyRegistry common.Address `mapstructure:"PUBLIC_CHAIN_DEPLOYMENT_PROXY_REGISTRY"`

	// Encryptor Service
	EncryptorService string `mapstructure:"CTS_ENCRYPTORSERVICE" validate:"required,oneof=plaintext aws gcp"`

	// Logger
	LogLevel   string `mapstructure:"LOG_LEVEL" validate:"required" default:"INFO"`
	LogHandler string `mapstructure:"LOG_HANDLER" validate:"required" default:"Text"`

	// Ports
	GRPCPort string `mapstructure:"CTS_GRPC_PORT" default:"8080"`
	HTTPPort string `mapstructure:"CTS_HTTP_PORT" default:"8090"`

	// WaitMinedTimeoutSecs caps how long a single tx broadcast waits for its
	// receipt (bind.WaitMined) before giving up, so a tx that never mines can't
	// wedge the service. A per-request caller deadline that is tighter still
	// wins (context.WithTimeout keeps the earliest deadline); this is the
	// last-resort ceiling for loose/absent caller deadlines. 0 uses the
	// service default.
	WaitMinedTimeoutSecs int `mapstructure:"CTS_WAIT_MINED_TIMEOUT_SECS" default:"300"`

	// gRPC mTLS — files are required. The CTS server requires client
	// certs and rejects any connection that doesn't present one signed
	// by the CA pinned via TLSCAFile.
	TLSCAFile   string `mapstructure:"CTS_TLS_CA_FILE" validate:"required" default:"/app/cts/certs/ca.crt"`
	TLSCertFile string `mapstructure:"CTS_TLS_CERT_FILE" validate:"required" default:"/app/cts/certs/server.crt"`
	TLSKeyFile  string `mapstructure:"CTS_TLS_KEY_FILE" validate:"required" default:"/app/cts/certs/server.key"`

	// OpenTelemetry
	OtelSDKDisabled bool   `mapstructure:"OTEL_SDK_DISABLED" default:"true"`
	OtelServiceName string `mapstructure:"OTEL_SERVICE_NAME"`

	// Fault injection (testing only — must never be enabled in production).
	// Port and persist path are CTS-specific so CTS, private-relayer and
	// public-relayer instances sharing the same per-participant .env file
	// each get an isolated HTTP server and on-disk persist file. The shared
	// FAULT_INJECTION_ENABLED flag acts as a single master switch across all
	// three services; binaries built without the faultinjection tag treat it
	// as a no-op. The KOS_-prefixed env var names are preserved for backward
	// compatibility with the resilience e2e suite, which still uses them.
	FaultInjectionEnabled     bool   `mapstructure:"FAULT_INJECTION_ENABLED" default:"false"`
	FaultInjectionPort        string `mapstructure:"KOS_FAULT_INJECTION_PORT" default:"9999"`
	FaultInjectionPersistPath string `mapstructure:"KOS_FAULT_INJECTION_PERSIST_PATH"`
}

// PublicChainEnabled reports whether the public chain has been deployed and
// its proxy registry written to the config. When false, CTS skips the public
// chain ethclient dial, key authorization wait, TxOps handler, and batch
// pipeline. Mirrors start_dev.sh's --no-public-chain flag: the contracts
// deployment service only appends PUBLIC_CHAIN_DEPLOYMENT_PROXY_REGISTRY to
// the participant env file when the public chain stack is deployed.
func (c *Config) PublicChainEnabled() bool {
	return c.PublicChainDeploymentProxyRegistry != (common.Address{})
}

// PrivateHubEnabled reports whether the private hub (PNH) has been deployed and
// its proxy registry written to the config. When false, CTS skips the hub
// ethclient dial, the hub key authorization wait, the hub ledger-error probe,
// the participant registrar, the hub + dvpoperator TxOps handlers, and the
// privatehub + dvpoperator batch pipelines. Mirrors PublicChainEnabled and
// start_dev.sh's --no-hub flag: the contracts deployment service only appends
// PNH_DEPLOYMENT_PROXY_REGISTRY to the participant env file when the hub stack
// is deployed.
func (c *Config) PrivateHubEnabled() bool {
	return c.PrivateHubDeploymentProxyRegistry != (common.Address{})
}
