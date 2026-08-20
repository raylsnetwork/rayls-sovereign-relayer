// Package keymanager owns the CTS-side bootstrap for relayer sign keys:
// create or fetch them on boot, then block until each relevant address is
// authorized in its chain's on-chain RelayAuthorizationRegistry.
package keymanager

//go:generate moq --pkg keymanager_test -out signkeys_mock_test.go . KeysService AccessManager

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/domain"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/keyqueue"
)

type KeysService interface {
	GetPrivateRelayerRaylsSignKeys(ctx context.Context) (domain.PrivateRelayerRaylsSignKeys, error)
	CreatePrivateRelayerRaylsSignKeys(ctx context.Context) (domain.PrivateRelayerRaylsSignKeys, error)
	GetPublicRelayerRaylsSignKeys(ctx context.Context) (domain.PublicRelayerRaylsSignKeys, error)
	CreatePublicRelayerRaylsSignKeys(ctx context.Context) (domain.PublicRelayerRaylsSignKeys, error)
}

type AccessManager interface {
	WaitForAuthorization(ctx context.Context, component string, keys []string) error
}

type BuildAccessManagerFunc func(
	backend *ethclient.Client,
	deploymentProxyRegistry common.Address,
	keys domain.RaylsSignKeyList,
) (AccessManager, error)

type ChainConfig struct {
	Backend                 *ethclient.Client
	DeploymentProxyRegistry common.Address
}

const (
	componentPrivateHub   = "private-hub"
	componentPrivacyNode  = "privacy-node"
	componentPrivateChain = "private-chain"
	componentPublicChain  = "public-chain"
)

type SignKeyManager struct {
	keysService KeysService

	hub    ChainConfig
	node   ChainConfig
	public ChainConfig

	privateHubEnabled  bool
	publicChainEnabled bool

	buildAccessManager BuildAccessManagerFunc
}

type Option func(*SignKeyManager)

// WithBuildAccessManager overrides the production factory. Testing seam only.
func WithBuildAccessManager(fn BuildAccessManagerFunc) Option {
	return func(m *SignKeyManager) { m.buildAccessManager = fn }
}

// WithPublicChain configures whether the public-chain authorization waits
// run. When false, wait #3 (public-relayer PN-side) and wait #4 (public
// chain) are both skipped and the public ChainConfig passed to
// NewSignKeyManager is ignored. Use false when CTS runs without a deployed
// public chain (e.g. start_dev.sh --no-public-chain).
func WithPublicChain(enabled bool) Option {
	return func(m *SignKeyManager) { m.publicChainEnabled = enabled }
}

// WithPrivateHub configures whether the private-hub (PNH) authorization wait
// runs. When false, wait #1 (private hub) is skipped and the hub ChainConfig
// passed to NewSignKeyManager is ignored. Use false when CTS runs without a
// deployed private hub (e.g. start_dev.sh --no-hub). Defaults to true so
// existing callers keep the hub wait unless they opt out.
func WithPrivateHub(enabled bool) Option {
	return func(m *SignKeyManager) { m.privateHubEnabled = enabled }
}

func NewSignKeyManager(
	keysService KeysService,
	hub, node, public ChainConfig,
	opts ...Option,
) *SignKeyManager {
	m := &SignKeyManager{
		keysService:        keysService,
		hub:                hub,
		node:               node,
		public:             public,
		privateHubEnabled:  true,
		buildAccessManager: buildLocalAccessManager,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Run ensures both relayer sign-key bundles exist (get-or-create), then waits
// sequentially for each bundle's addresses to appear in the respective chain's
// RelayAuthorizationRegistry. Wait order:
//
//  1. Private Hub (skipped when WithPrivateHub(false))
//  2. Privacy Node
//  3. Private Chain (public-relayer PN-side keys — skipped when WithPublicChain(false))
//  4. Public Chain (skipped when WithPublicChain(false))
//
// WithPrivateHub(false) skips wait #1 because the hub is not deployed.
// WithPublicChain(false) skips both wait #3 and #4 because the deployer only
// grants RELAYER to those keys when the public-relayer stack runs.
//
// DvP operator keys are not waited on — mirrors the private-relayer bootstrap.
//
// Returns both decrypted bundles on success. Any failure (get/create,
// registry construction, wait error) is surfaced with a label identifying
// which step failed.
func (m *SignKeyManager) Run(ctx context.Context) (
	domain.PrivateRelayerRaylsSignKeys,
	domain.PublicRelayerRaylsSignKeys,
	error,
) {
	privateKeys, err := m.ensurePrivateRelayerSignKeys(ctx)
	if err != nil {
		return domain.PrivateRelayerRaylsSignKeys{}, domain.PublicRelayerRaylsSignKeys{},
			fmt.Errorf("ensuring private relayer sign keys: %w", err)
	}

	publicKeys, err := m.ensurePublicRelayerSignKeys(ctx)
	if err != nil {
		return domain.PrivateRelayerRaylsSignKeys{}, domain.PublicRelayerRaylsSignKeys{},
			fmt.Errorf("ensuring public relayer sign keys: %w", err)
	}

	if err := validateRequiredKeysets(privateKeys, publicKeys, m.privateHubEnabled, m.publicChainEnabled); err != nil {
		return domain.PrivateRelayerRaylsSignKeys{}, domain.PublicRelayerRaylsSignKeys{}, err
	}

	type waitSpec struct {
		chain     ChainConfig
		component string
		keys      domain.RaylsSignKeyList
	}
	var waits []waitSpec
	if m.privateHubEnabled {
		waits = append(waits, waitSpec{m.hub, componentPrivateHub, privateKeys.PrivateHubKeys})
	}
	waits = append(waits, waitSpec{m.node, componentPrivacyNode, privateKeys.PrivateNodeKeys})
	if m.publicChainEnabled {
		waits = append(waits,
			waitSpec{m.node, componentPrivateChain, publicKeys.PrivateChainKeys},
			waitSpec{m.public, componentPublicChain, publicKeys.PublicChainKeys},
		)
	}

	for _, w := range waits {
		registry, err := m.buildAccessManager(w.chain.Backend, w.chain.DeploymentProxyRegistry, w.keys)
		if err != nil {
			return domain.PrivateRelayerRaylsSignKeys{}, domain.PublicRelayerRaylsSignKeys{},
				fmt.Errorf("building access manager for %s: %w", w.component, err)
		}
		if err := registry.WaitForAuthorization(ctx, w.component, signKeyListToHex(w.keys)); err != nil {
			return domain.PrivateRelayerRaylsSignKeys{}, domain.PublicRelayerRaylsSignKeys{},
				fmt.Errorf("awaiting authorization for %s: %w", w.component, err)
		}
	}

	return privateKeys, publicKeys, nil
}

func (m *SignKeyManager) ensurePrivateRelayerSignKeys(ctx context.Context) (domain.PrivateRelayerRaylsSignKeys, error) {
	keys, err := m.keysService.GetPrivateRelayerRaylsSignKeys(ctx)
	if err == nil {
		return keys, nil
	}
	if !errors.Is(err, service.ErrNoRaylsSignKeys) {
		return domain.PrivateRelayerRaylsSignKeys{}, fmt.Errorf("getting keys: %w", err)
	}

	slog.Info("no private relayer sign keys found — creating initial keyset")
	keys, err = m.keysService.CreatePrivateRelayerRaylsSignKeys(ctx)
	if err != nil {
		return domain.PrivateRelayerRaylsSignKeys{}, fmt.Errorf("creating keys: %w", err)
	}
	return keys, nil
}

func (m *SignKeyManager) ensurePublicRelayerSignKeys(ctx context.Context) (domain.PublicRelayerRaylsSignKeys, error) {
	keys, err := m.keysService.GetPublicRelayerRaylsSignKeys(ctx)
	if err == nil {
		return keys, nil
	}
	if !errors.Is(err, service.ErrNoRaylsSignKeys) {
		return domain.PublicRelayerRaylsSignKeys{}, fmt.Errorf("getting keys: %w", err)
	}

	slog.Info("no public relayer sign keys found — creating initial keyset")
	keys, err = m.keysService.CreatePublicRelayerRaylsSignKeys(ctx)
	if err != nil {
		return domain.PublicRelayerRaylsSignKeys{}, fmt.Errorf("creating keys: %w", err)
	}
	return keys, nil
}

// validateRequiredKeysets fails fast if any keyset the waits depend on is
// empty. Without this, WaitForAuthorization would trivially pass on an empty
// address set — a silent boot-time footgun if a migration ever wrote an
// empty slice to the DB.
func validateRequiredKeysets(
	priv domain.PrivateRelayerRaylsSignKeys,
	pub domain.PublicRelayerRaylsSignKeys,
	privateHubEnabled bool,
	publicChainEnabled bool,
) error {
	type check struct {
		name string
		keys domain.RaylsSignKeyList
	}
	checks := []check{
		{"private relayer PrivateNodeKeys", priv.PrivateNodeKeys},
	}
	if privateHubEnabled {
		checks = append(checks, check{"private relayer PrivateHubKeys", priv.PrivateHubKeys})
	}
	if publicChainEnabled {
		checks = append(checks,
			check{"public relayer PrivateChainKeys", pub.PrivateChainKeys},
			check{"public relayer PublicChainKeys", pub.PublicChainKeys},
		)
	}
	for _, c := range checks {
		if len(c.keys) == 0 {
			return fmt.Errorf("empty keyset: %s", c.name)
		}
	}
	return nil
}

// signKeyListToHex matches the format WaitForAuthorization expects:
// lowercase, no 0x prefix (crypto.HexToECDSA trims either form).
func signKeyListToHex(keys domain.RaylsSignKeyList) []string {
	out := make([]string, len(keys))
	for i, k := range keys {
		out[i] = hex.EncodeToString(crypto.FromECDSA(k))
	}
	return out
}

// buildLocalAccessManager is the default BuildAccessManagerFunc used in
// production. It wires AuthGen + key queue + LocalExecutor +
// DeploymentProxyRegistry + AccessManager.
func buildLocalAccessManager(
	backend *ethclient.Client,
	deploymentProxyRegistry common.Address,
	keys domain.RaylsSignKeyList,
) (AccessManager, error) {
	authGen, err := contractclient.NewAuthGen(context.Background(), backend)
	if err != nil {
		return nil, fmt.Errorf("auth gen: %w", err)
	}

	executor := contractclient.NewLocalExecutor(authGen, newBootKeyQueue(keys), backend)

	proxy, err := contractclient.NewDeploymentProxyRegistryClient(deploymentProxyRegistry, backend)
	if err != nil {
		return nil, fmt.Errorf("deployment proxy registry: %w", err)
	}

	accessManagerAddr, err := proxy.GetContractAddress("RaylsAccessManager")
	if err != nil {
		return nil, fmt.Errorf("resolving RelayAuthorizationRegistry address: %w", err)
	}

	return contractclient.NewAccessManager(accessManagerAddr, executor), nil
}

func newBootKeyQueue(keys []*ecdsa.PrivateKey) *keyqueue.RaylsSignPrivateKeyManager {
	q := keyqueue.New(len(keys))
	for _, k := range keys {
		q.Enqueue(k)
	}
	return q
}
