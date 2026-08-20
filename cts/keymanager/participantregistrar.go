package keymanager

//go:generate moq --pkg keymanager_test -out participantregistrar_mock_test.go . RegistrarKeysService EncryptPrivateKeyer ParticipantStorage

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"slices"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	ps "github.com/raylsnetwork/rayls-sovereign-relayer/contracts/ParticipantStorageV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cryptography"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/domain"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/service"
)

// RegistrarKeysService is the narrow slice of *service.KeysService the
// registrar depends on — only the methods needed by the four ensure-and-register
// flows.
type RegistrarKeysService interface {
	GetRaylsViewKeyPair(ctx context.Context, blockNum uint64) (domain.RaylsViewKeyPair, error)
	CreateRaylsViewKeyPair(ctx context.Context, initialBlock uint64) (domain.RaylsViewKeyPair, error)

	GenerateKeyAgreement(chainID *big.Int, publicKey []byte) (ciphertext, sharedSecret, digest []byte, err error)
	CreateKeyAgreement(ctx context.Context, chainID *big.Int, sharedSecret []byte, blockNum uint64) error

	GetPaymentSpendKey(ctx context.Context) (domain.PaymentSpendKeys, error)
	CreatePaymentSpendKey(ctx context.Context, chainID *big.Int) (domain.PaymentSpendKeys, error)
}

// EncryptPrivateKeyer is the one method from *service.EncryptService needed
// for audit-info registration.
type EncryptPrivateKeyer interface {
	EncryptPrivateKey(ctx context.Context, chainID uint64, blockNumber uint64) (publicKeyHex string, encrypted []byte, mac []byte, initialBlock uint64, err error)
}

// ParticipantStorage is the narrow slice of *contractclient.ParticipantStorageClient
// the registrar uses.
type ParticipantStorage interface {
	GetChainViewData(ctx context.Context, chainID, blockNumber *big.Int) (ps.ParticipantStructsPrivacyNodeViewData, error)
	SetChainViewData(ctx context.Context, chainID *big.Int, raylsViewPublicKey string, blockNumber *big.Int) error

	GetAllParticipantsViewData(ctx context.Context, blockNumber *big.Int) ([]ps.ParticipantStructsPrivacyNodeViewData, error)
	GetKeyAgreements(ctx context.Context, chainID, blockNumber *big.Int) ([]ps.ParticipantStructsKeyAgreementData, error)
	InitiateKeyAgreement(ctx context.Context, chainID *big.Int, ciphertext, digest []byte, blockNumber *big.Int) error

	GetAuditInfo(ctx context.Context, chainID, blockNumber *big.Int) (ps.ParticipantStructsAuditInfoData, error)
	SetAuditInfo(ctx context.Context, chainID, blockNumber *big.Int, publicKey string, encrPrivateKey, mac []byte) error
	GetVenOperatorChainInfo(ctx context.Context, blockNumber *big.Int) (ps.ParticipantStructsPrivacyNodeViewData, error)

	GetPaymentSpendPublicKey(ctx context.Context, chainID *big.Int) (*big.Int, error)
	SetPaymentSpendPublicKey(ctx context.Context, chainID *big.Int, paymentSpendPublicKey *big.Int, plAddresses []common.Address) error
}

// GetBlockNumberFunc reads the current block number from the given backend.
// Production uses getHubBlockNumber; tests substitute a stub.
type GetBlockNumberFunc func(ctx context.Context, hub *ethclient.Client) (uint64, error)

// BuildParticipantStorageFunc is the injected factory used by the registrar.
// Production uses buildHubParticipantStorage; tests substitute a stub.
type BuildParticipantStorageFunc func(
	hub *ethclient.Client,
	deploymentProxyRegistry common.Address,
	plChainID, venChainID *big.Int,
	hubKeys domain.RaylsSignKeyList,
) (ParticipantStorage, error)

// ErrRaylsViewKeyMismatch indicates the view key stored locally does not
// match the one registered on-chain — a fatal boot-time condition.
var ErrRaylsViewKeyMismatch = errors.New("local rayls view public key does not match on-chain value")

// ErrPaymentSpendKeyMismatch indicates the Poseidon hash of the local
// payment-spend secret key does not match the value registered on-chain — a
// fatal boot-time condition.
var ErrPaymentSpendKeyMismatch = errors.New("local payment-spend public key hash does not match on-chain value")

type ParticipantRegistrar struct {
	chainID    *big.Int
	venChainID *big.Int

	hub ChainConfig

	keysService RegistrarKeysService
	encrypt     EncryptPrivateKeyer

	hubKeys             domain.RaylsSignKeyList
	relayerHubAddresses []common.Address

	buildStorage   BuildParticipantStorageFunc
	getBlockNumber GetBlockNumberFunc
}

type ParticipantRegistrarOption func(*ParticipantRegistrar)

// WithBuildParticipantStorage overrides the production factory. Testing seam only.
func WithBuildParticipantStorage(fn BuildParticipantStorageFunc) ParticipantRegistrarOption {
	return func(r *ParticipantRegistrar) { r.buildStorage = fn }
}

// WithBlockNumberFunc overrides the production block-number source
// (getHubBlockNumber). Testing seam only.
func WithBlockNumberFunc(fn GetBlockNumberFunc) ParticipantRegistrarOption {
	return func(r *ParticipantRegistrar) { r.getBlockNumber = fn }
}

func NewParticipantRegistrar(
	chainID, venChainID *big.Int,
	hub ChainConfig,
	keysService RegistrarKeysService,
	encrypt EncryptPrivateKeyer,
	hubKeys domain.RaylsSignKeyList,
	opts ...ParticipantRegistrarOption,
) *ParticipantRegistrar {
	addresses := make([]common.Address, len(hubKeys))
	for i, k := range hubKeys {
		addresses[i] = crypto.PubkeyToAddress(k.PublicKey)
	}

	r := &ParticipantRegistrar{
		chainID:             chainID,
		venChainID:          venChainID,
		hub:                 hub,
		keysService:         keysService,
		encrypt:             encrypt,
		hubKeys:             hubKeys,
		relayerHubAddresses: addresses,
		buildStorage:        buildHubParticipantStorage,
		getBlockNumber:      getHubBlockNumber,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Run performs the four on-hub registrations sequentially. Each step is
// idempotent against on-chain state, so restarting the service is safe.
//
// Prerequisite: SignKeyManager must have completed — in particular the Hub
// RelayAuthorizationRegistry wait — because every write below goes through
// ParticipantStorage which enforces onlyAuthorizedCaller.
func (r *ParticipantRegistrar) Run(ctx context.Context) error {
	storage, err := r.buildStorage(r.hub.Backend, r.hub.DeploymentProxyRegistry, r.chainID, r.venChainID, r.hubKeys)
	if err != nil {
		return fmt.Errorf("building participant storage client: %w", err)
	}

	blockNumber, err := r.getBlockNumber(ctx, r.hub.Backend)
	if err != nil {
		return fmt.Errorf("reading node block number: %w", err)
	}

	if err := r.ensureViewKeys(ctx, blockNumber, storage); err != nil {
		return fmt.Errorf("registering view keys: %w", err)
	}
	if err := r.ensureKeyAgreements(ctx, blockNumber, storage); err != nil {
		return fmt.Errorf("registering key agreements: %w", err)
	}
	if err := r.ensureAuditInfo(ctx, blockNumber, storage); err != nil {
		return fmt.Errorf("registering audit info: %w", err)
	}
	if err := r.ensurePaymentSpendKey(ctx, storage); err != nil {
		return fmt.Errorf("registering payment spend key: %w", err)
	}
	return nil
}

func (r *ParticipantRegistrar) ensureViewKeys(ctx context.Context, blockNumber uint64, storage ParticipantStorage) error {
	pair, err := r.getOrCreateViewKeyPair(ctx, blockNumber)
	if err != nil {
		return err
	}
	localPublicKeyHex := hex.EncodeToString(pair.RaylsViewPublicKey.Bytes())

	blockBig := new(big.Int).SetUint64(blockNumber)
	onChain, err := storage.GetChainViewData(ctx, r.chainID, blockBig)
	switch {
	case errors.Is(err, contractclient.ErrNoChainInfo):
		if err := storage.SetChainViewData(ctx, r.chainID, localPublicKeyHex, blockBig); err != nil {
			return fmt.Errorf("setting chain view data: %w", err)
		}
		return nil
	case err == nil:
		if onChain.RaylsViewPublicKey != localPublicKeyHex {
			return ErrRaylsViewKeyMismatch
		}
		return nil
	default:
		return fmt.Errorf("getting chain view data: %w", err)
	}
}

func (r *ParticipantRegistrar) getOrCreateViewKeyPair(ctx context.Context, blockNumber uint64) (domain.RaylsViewKeyPair, error) {
	pair, err := r.keysService.GetRaylsViewKeyPair(ctx, blockNumber)
	if err == nil {
		return pair, nil
	}
	if !errors.Is(err, service.ErrNoRaylsViewKeysSet) && !errors.Is(err, service.ErrNoApplicableRaylsViewKeys) {
		return domain.RaylsViewKeyPair{}, fmt.Errorf("getting rayls view key pair: %w", err)
	}

	slog.Info("no rayls view key pair found — creating initial pair", slog.Uint64("block", blockNumber))
	pair, err = r.keysService.CreateRaylsViewKeyPair(ctx, blockNumber)
	if err != nil {
		return domain.RaylsViewKeyPair{}, fmt.Errorf("creating rayls view key pair: %w", err)
	}
	return pair, nil
}

func (r *ParticipantRegistrar) ensureKeyAgreements(ctx context.Context, blockNumber uint64, storage ParticipantStorage) error {
	blockBig := new(big.Int).SetUint64(blockNumber)

	existing, err := storage.GetKeyAgreements(ctx, r.chainID, blockBig)
	if err != nil {
		return fmt.Errorf("getting existing key agreements: %w", err)
	}

	// VALIDATION (race-condition theory): GetAllParticipantsViewData can race when
	// participants register near-simultaneously — a peer's view data may not have
	// propagated to the storage we read yet, so the returned set is missing it and
	// the pairwise key agreement below is silently skipped (there is no retry). Poll
	// with naive exponential backoff until a peer chain (any chain that is neither
	// this node nor the VEN operator) appears, logging each attempt so we can observe
	// whether the set starts incomplete and fills in over time.
	const maxViewDataAttempts = 8
	viewDataBackoff := 500 * time.Millisecond
	const maxViewDataBackoff = 8 * time.Second

	var all []ps.ParticipantStructsPrivacyNodeViewData
	for attempt := 1; ; attempt++ {
		all, err = storage.GetAllParticipantsViewData(ctx, blockBig)
		if err != nil {
			return fmt.Errorf("getting all participants view data: %w", err)
		}

		hasPeer := slices.ContainsFunc(all, func(p ps.ParticipantStructsPrivacyNodeViewData) bool {
			return p.ChainId.Cmp(r.chainID) != 0 && p.ChainId.Cmp(r.venChainID) != 0
		})

		slog.Info("ensureKeyAgreements: participants view data poll",
			slog.Int("attempt", attempt),
			slog.Int("count", len(all)),
			slog.Bool("has_peer", hasPeer),
			slog.Uint64("block", blockNumber),
		)

		if hasPeer || attempt >= maxViewDataAttempts {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(viewDataBackoff):
		}
		if viewDataBackoff *= 2; viewDataBackoff > maxViewDataBackoff {
			viewDataBackoff = maxViewDataBackoff
		}
	}

	for _, participant := range all {
		if slices.ContainsFunc(existing, func(ka ps.ParticipantStructsKeyAgreementData) bool {
			return ka.ChainId.Cmp(participant.ChainId) == 0
		}) {
			continue
		}

		if err := r.initiateKeyAgreementFor(ctx, participant, blockNumber, storage); err != nil {
			return err
		}
	}
	return nil
}

func (r *ParticipantRegistrar) initiateKeyAgreementFor(
	ctx context.Context,
	participant ps.ParticipantStructsPrivacyNodeViewData,
	blockNumber uint64,
	storage ParticipantStorage,
) error {
	slog.Info("creating key agreement for participant", slog.String("chain_id", participant.ChainId.String()))

	viewKeyBytes, err := hex.DecodeString(participant.RaylsViewPublicKey)
	if err != nil {
		return fmt.Errorf("decoding view public key for chain %s: %w", participant.ChainId, err)
	}

	ciphertext, sharedSecret, digest, err := r.keysService.GenerateKeyAgreement(participant.ChainId, viewKeyBytes)
	if err != nil {
		return fmt.Errorf("generating key agreement for chain %s: %w", participant.ChainId, err)
	}

	blockBig := new(big.Int).SetUint64(blockNumber)
	err = storage.InitiateKeyAgreement(ctx, participant.ChainId, ciphertext, digest, blockBig)
	switch {
	case err == nil:
		if err := r.keysService.CreateKeyAgreement(ctx, participant.ChainId, sharedSecret, blockNumber); err != nil {
			return fmt.Errorf("persisting shared secret for chain %s: %w", participant.ChainId, err)
		}
		return nil
	case errors.Is(err, contractclient.ErrOutdatedKeyAgreement):
		slog.Warn(
			"key agreement already exists for newer block number, skipping",
			slog.String("chain_id", participant.ChainId.String()),
		)
		return nil
	default:
		return fmt.Errorf("initiating key agreement for chain %s: %w", participant.ChainId, err)
	}
}

func (r *ParticipantRegistrar) ensureAuditInfo(ctx context.Context, blockNumber uint64, storage ParticipantStorage) error {
	blockBig := new(big.Int).SetUint64(blockNumber)

	_, err := storage.GetAuditInfo(ctx, r.chainID, blockBig)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, contractclient.ErrNoAuditInfo):
		// fall through to registration
	default:
		return fmt.Errorf("getting audit info: %w", err)
	}

	venChainInfo, err := storage.GetVenOperatorChainInfo(ctx, blockBig)
	if err != nil {
		return fmt.Errorf("getting ven operator chain info: %w", err)
	}

	publicKeyHex, encrypted, mac, initialBlock, err := r.encrypt.EncryptPrivateKey(ctx, venChainInfo.ChainId.Uint64(), blockNumber)
	if err != nil {
		return fmt.Errorf("encrypting private view key for ven operator: %w", err)
	}

	if err := storage.SetAuditInfo(
		ctx,
		r.chainID,
		new(big.Int).SetUint64(initialBlock),
		publicKeyHex,
		encrypted,
		mac,
	); err != nil {
		return fmt.Errorf("setting audit info: %w", err)
	}
	return nil
}

func (r *ParticipantRegistrar) ensurePaymentSpendKey(ctx context.Context, storage ParticipantStorage) error {
	keys, err := r.keysService.GetPaymentSpendKey(ctx)
	switch {
	case err == nil:
	case errors.Is(err, service.ErrNoPaymentSpendKeySet):
		slog.Info("no payment spend key found — creating initial key")
		keys, err = r.keysService.CreatePaymentSpendKey(ctx, r.chainID)
		if err != nil {
			return fmt.Errorf("creating payment spend key: %w", err)
		}
	default:
		return fmt.Errorf("getting payment spend key: %w", err)
	}

	localPublicKeyHash, err := cryptography.GetPoseidonHashModNumber(
		[]*big.Int{keys.SecretKey, keys.SecretKey},
		cryptography.JubJubPrimeSubGroup,
	)
	if err != nil {
		return fmt.Errorf("computing payment-spend public key hash: %w", err)
	}

	onChain, err := storage.GetPaymentSpendPublicKey(ctx, r.chainID)
	switch {
	case err == nil:
		if localPublicKeyHash.Cmp(onChain) != 0 {
			return ErrPaymentSpendKeyMismatch
		}
		return nil
	case errors.Is(err, contractclient.ErrNoPaymentSpendPublicKey):
		if err := storage.SetPaymentSpendPublicKey(ctx, r.chainID, localPublicKeyHash, r.relayerHubAddresses); err != nil {
			return fmt.Errorf("setting payment-spend public key: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("getting payment-spend public key: %w", err)
	}
}

// buildHubParticipantStorage is the default BuildParticipantStorageFunc.
// Wires AuthGen + key queue + LocalExecutor + DeploymentProxyRegistry +
// ParticipantStorageClient against the Hub backend.
func buildHubParticipantStorage(
	hub *ethclient.Client,
	deploymentProxyRegistry common.Address,
	plChainID, venChainID *big.Int,
	hubKeys domain.RaylsSignKeyList,
) (ParticipantStorage, error) {
	authGen, err := contractclient.NewAuthGen(context.Background(), hub)
	if err != nil {
		return nil, fmt.Errorf("auth gen: %w", err)
	}

	executor := contractclient.NewLocalExecutor(authGen, newBootKeyQueue(hubKeys), hub)

	proxy, err := contractclient.NewDeploymentProxyRegistryClient(deploymentProxyRegistry, hub)
	if err != nil {
		return nil, fmt.Errorf("deployment proxy registry: %w", err)
	}

	psAddr, err := proxy.GetContractAddress("ParticipantStorage")
	if err != nil {
		return nil, fmt.Errorf("resolving ParticipantStorage address: %w", err)
	}

	return contractclient.NewParticipantStorageClient(psAddr, plChainID, venChainID, executor), nil
}

// getHubBlockNumber is the default GetBlockNumberFunc.
func getHubBlockNumber(ctx context.Context, hub *ethclient.Client) (uint64, error) {
	return hub.BlockNumber(ctx)
}
