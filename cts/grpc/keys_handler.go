package grpc

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/domain"
	keyspb "github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/keys"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type KeysService interface {
	CreateRaylsViewKeyPair(ctx context.Context, initialBlock uint64) (domain.RaylsViewKeyPair, error)
	GetRaylsViewKeyPair(ctx context.Context, blockNum uint64) (domain.RaylsViewKeyPair, error)
	DeleteRaylsViewKeyPair(ctx context.Context, publicKey string) error
	RecoverViewSalt(ctx context.Context, blockNumber uint64, ctxt []byte) (*big.Int, error)

	CreatePublicRelayerRaylsSignKeys(ctx context.Context) (domain.PublicRelayerRaylsSignKeys, error)
	CreatePrivateRelayerRaylsSignKeys(ctx context.Context) (domain.PrivateRelayerRaylsSignKeys, error)

	GetPublicRelayerRaylsSignKeys(ctx context.Context) (domain.PublicRelayerRaylsSignKeys, error)
	GetPrivateRelayerRaylsSignKeys(ctx context.Context) (domain.PrivateRelayerRaylsSignKeys, error)

	GetPublicRelayerRaylsSignPublicAddresses(ctx context.Context) (domain.PublicRelayerRaylsSignPublicAddresses, error)
	GetPrivateRelayerRaylsSignPublicAddresses(ctx context.Context) (domain.PrivateRelayerRaylsSignPublicAddresses, error)

	GenerateKeyAgreement(chainID *big.Int, publicKey []byte) (ciphertext []byte, sharedSecret []byte, digest []byte, err error)
	CreateKeyAgreement(ctx context.Context, chainID *big.Int, sharedSecret []byte, blockNum uint64) error
	CompleteKeyAgreement(ctx context.Context, chainID *big.Int, ciphertext []byte, blockNum uint64) error

	CreatePaymentSpendKey(ctx context.Context, chainId *big.Int) (domain.PaymentSpendKeys, error)
	GetPaymentSpendKey(ctx context.Context) (domain.PaymentSpendKeys, error)

	CreateSelfSecret(ctx context.Context, rFactor *big.Int, blockNumber uint64, resourceID []byte) error

	GetAllSharedSecrets(ctx context.Context, blockNumber uint64) ([]domain.SharedSecret, error)
	GetEnygmaSharedSecrets(ctx context.Context, chainIDs []*big.Int, myChainID *big.Int, blockNumber uint64, resourceID []byte) ([]domain.SharedSecret, error)
	GenerateEnygmaSharedSecrets(chainIDs []*big.Int, blockNumber uint64, sharedSecrets []domain.SharedSecret) ([]*big.Int, []*big.Int, []*big.Int, error)
}

type KeysHandler struct {
	keyspb.UnimplementedKeysServiceServer

	svc KeysService
}

func NewKeysHandler(svc KeysService) *KeysHandler {
	return &KeysHandler{svc: svc}
}

// View Keys

func (h *KeysHandler) CreateViewKeyPair(
	ctx context.Context,
	req *keyspb.CreateViewKeyPairRequest,
) (*keyspb.CreateViewKeyPairResponse, error) {
	pair, err := h.svc.CreateRaylsViewKeyPair(ctx, req.GetInitialBlock())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create view key pair: %v", err)
	}

	return &keyspb.CreateViewKeyPairResponse{
		PublicKey: hex.EncodeToString(pair.RaylsViewPublicKey.Bytes()),
	}, nil
}

// TODO: handle not found case
func (h *KeysHandler) GetViewPublicKey(
	ctx context.Context,
	req *keyspb.GetViewPublicKeyRequest,
) (*keyspb.GetViewPublicKeyResponse, error) {
	pair, err := h.svc.GetRaylsViewKeyPair(ctx, req.GetBlockNumber())
	if err != nil {
		if errors.Is(err, service.ErrNoRaylsViewKeysSet) {
			return nil, status.Errorf(codes.NotFound, "get view key pair: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "get view key pair: %v", err)
	}

	return &keyspb.GetViewPublicKeyResponse{
		PublicKey: hex.EncodeToString(pair.RaylsViewPublicKey.Bytes()),
	}, nil
}

func (h *KeysHandler) DeleteViewKeyPair(
	ctx context.Context,
	req *keyspb.DeleteViewKeyPairRequest,
) (*keyspb.DeleteViewKeyPairResponse, error) {
	if req.GetPublicKey() == "" {
		return nil, status.Error(codes.InvalidArgument, "public_key is required")
	}

	err := h.svc.DeleteRaylsViewKeyPair(ctx, req.GetPublicKey())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete view key pair: %v", err)
	}

	return &keyspb.DeleteViewKeyPairResponse{}, nil
}

func (h *KeysHandler) RecoverViewSalt(
	ctx context.Context,
	req *keyspb.RecoverViewSaltRequest,
) (*keyspb.RecoverViewSaltResponse, error) {
	if len(req.GetCtxt()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "ctxt is required")
	}

	salt, err := h.svc.RecoverViewSalt(ctx, req.GetBlockNumber(), req.GetCtxt())
	if err != nil {
		if errors.Is(err, service.ErrNoRaylsViewKeysSet) {
			return nil, status.Errorf(codes.NotFound, "recover view salt: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "recover view salt: %v", err)
	}

	return &keyspb.RecoverViewSaltResponse{
		Salt: salt.Bytes(),
	}, nil
}

// Sign Keys

func (h *KeysHandler) CreateSignKeys(
	ctx context.Context,
	req *keyspb.CreateSignKeysRequest,
) (*keyspb.CreateSignKeysResponse, error) {
	keySets := make(map[int32]*keyspb.StringList)
	var err error

	switch req.GetServiceType() {
	case keyspb.ServiceType_SERVICE_TYPE_PUBLIC_RELAYER:
		var k domain.PublicRelayerRaylsSignKeys
		k, err = h.svc.CreatePublicRelayerRaylsSignKeys(ctx)
		if err == nil {
			keySets[int32(keyspb.KeySetName_KEY_SET_NAME_PUBLIC_CHAIN)] = &keyspb.StringList{Values: marshalECDSAKeys(k.PublicChainKeys)}
			keySets[int32(keyspb.KeySetName_KEY_SET_NAME_PRIVATE_CHAIN)] = &keyspb.StringList{Values: marshalECDSAKeys(k.PrivateChainKeys)}
		}
	case keyspb.ServiceType_SERVICE_TYPE_PRIVATE_RELAYER:
		var k domain.PrivateRelayerRaylsSignKeys
		k, err = h.svc.CreatePrivateRelayerRaylsSignKeys(ctx)
		if err == nil {
			keySets[int32(keyspb.KeySetName_KEY_SET_NAME_PRIVATE_HUB)] = &keyspb.StringList{Values: marshalECDSAKeys(k.PrivateHubKeys)}
			keySets[int32(keyspb.KeySetName_KEY_SET_NAME_PRIVATE_HUB_DVP_OPERATOR)] = &keyspb.StringList{Values: marshalECDSAKeys(k.PrivateHubDvpOperatorKeys)}
			keySets[int32(keyspb.KeySetName_KEY_SET_NAME_PRIVACY_NODE)] = &keyspb.StringList{Values: marshalECDSAKeys(k.PrivateNodeKeys)}
		}
	default:
		return nil, status.Error(codes.InvalidArgument, "service_type is required")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "create sign keys: %v", err)
	}

	return &keyspb.CreateSignKeysResponse{KeySets: keySets}, nil
}

func (h *KeysHandler) GetSignKeys(
	ctx context.Context,
	req *keyspb.GetSignKeysRequest,
) (*keyspb.GetSignKeysResponse, error) {
	keySets := make(map[int32]*keyspb.StringList)
	var err error

	switch req.GetServiceType() {
	case keyspb.ServiceType_SERVICE_TYPE_PUBLIC_RELAYER:
		var k domain.PublicRelayerRaylsSignKeys
		k, err = h.svc.GetPublicRelayerRaylsSignKeys(ctx)
		if err == nil {
			keySets[int32(keyspb.KeySetName_KEY_SET_NAME_PUBLIC_CHAIN)] = &keyspb.StringList{Values: marshalECDSAKeys(k.PublicChainKeys)}
			keySets[int32(keyspb.KeySetName_KEY_SET_NAME_PRIVATE_CHAIN)] = &keyspb.StringList{Values: marshalECDSAKeys(k.PrivateChainKeys)}
		}
	case keyspb.ServiceType_SERVICE_TYPE_PRIVATE_RELAYER:
		var k domain.PrivateRelayerRaylsSignKeys
		k, err = h.svc.GetPrivateRelayerRaylsSignKeys(ctx)
		if err == nil {
			keySets[int32(keyspb.KeySetName_KEY_SET_NAME_PRIVATE_HUB)] = &keyspb.StringList{Values: marshalECDSAKeys(k.PrivateHubKeys)}
			keySets[int32(keyspb.KeySetName_KEY_SET_NAME_PRIVATE_HUB_DVP_OPERATOR)] = &keyspb.StringList{Values: marshalECDSAKeys(k.PrivateHubDvpOperatorKeys)}
			keySets[int32(keyspb.KeySetName_KEY_SET_NAME_PRIVACY_NODE)] = &keyspb.StringList{Values: marshalECDSAKeys(k.PrivateNodeKeys)}
		}
	default:
		return nil, status.Error(codes.InvalidArgument, "service_type is required")
	}

	if err != nil {
		if errors.Is(err, service.ErrNoRaylsSignKeys) {
			return nil, status.Errorf(codes.NotFound, "get sign keys: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "get sign keys: %v", err)
	}

	return &keyspb.GetSignKeysResponse{KeySets: keySets}, nil
}

func (h *KeysHandler) GetSignAddresses(
	ctx context.Context,
	req *keyspb.GetSignAddressesRequest,
) (*keyspb.GetSignAddressesResponse, error) {
	addressSets := make(map[int32]*keyspb.StringList)
	var err error

	switch req.GetServiceType() {
	case keyspb.ServiceType_SERVICE_TYPE_PUBLIC_RELAYER:
		var a domain.PublicRelayerRaylsSignPublicAddresses
		a, err = h.svc.GetPublicRelayerRaylsSignPublicAddresses(ctx)
		if err == nil {
			addressSets[int32(keyspb.KeySetName_KEY_SET_NAME_PUBLIC_CHAIN)] = &keyspb.StringList{Values: marshalAddresses(a.PublicChainAddresses)}
			addressSets[int32(keyspb.KeySetName_KEY_SET_NAME_PRIVATE_CHAIN)] = &keyspb.StringList{Values: marshalAddresses(a.PrivateChainAddresses)}
		}
	case keyspb.ServiceType_SERVICE_TYPE_PRIVATE_RELAYER:
		var a domain.PrivateRelayerRaylsSignPublicAddresses
		a, err = h.svc.GetPrivateRelayerRaylsSignPublicAddresses(ctx)
		if err == nil {
			addressSets[int32(keyspb.KeySetName_KEY_SET_NAME_PRIVATE_HUB)] = &keyspb.StringList{Values: marshalAddresses(a.PrivateHubAddresses)}
			addressSets[int32(keyspb.KeySetName_KEY_SET_NAME_PRIVACY_NODE)] = &keyspb.StringList{Values: marshalAddresses(a.PrivateChainAddresses)}
			addressSets[int32(keyspb.KeySetName_KEY_SET_NAME_PRIVATE_HUB_DVP_OPERATOR)] = &keyspb.StringList{Values: marshalAddresses(a.PrivateHubDvpOperatorAddresses)}
		}
	default:
		return nil, status.Error(codes.InvalidArgument, "service_type is required")
	}

	if err != nil {
		if errors.Is(err, service.ErrNoRaylsSignKeys) {
			return nil, status.Errorf(codes.NotFound, "get sign addresses: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "get sign addresses: %v", err)
	}

	return &keyspb.GetSignAddressesResponse{AddressSets: addressSets}, nil
}

// Key Agreement

func (h *KeysHandler) GenerateKeyAgreement(
	ctx context.Context,
	req *keyspb.GenerateKeyAgreementRequest,
) (*keyspb.GenerateKeyAgreementResponse, error) {
	if req.GetChainID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chain_id is required")
	}
	if len(req.GetPublicKey()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "public_key is required")
	}

	ciphertext, sharedSecret, digest, err := h.svc.GenerateKeyAgreement(
		new(big.Int).SetUint64(req.GetChainID()),
		req.GetPublicKey(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate key agreement: %v", err)
	}

	return &keyspb.GenerateKeyAgreementResponse{
		Ciphertext:   ciphertext,
		SharedSecret: sharedSecret,
		Digest:       digest,
	}, nil
}

func (h *KeysHandler) CreateKeyAgreement(
	ctx context.Context,
	req *keyspb.CreateKeyAgreementRequest,
) (*keyspb.CreateKeyAgreementResponse, error) {
	if req.GetChainID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chain_id is required")
	}
	if len(req.GetSharedSecret()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "shared_secret is required")
	}

	err := h.svc.CreateKeyAgreement(
		ctx,
		new(big.Int).SetUint64(req.GetChainID()),
		req.GetSharedSecret(),
		req.GetBlockNumber(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create key agreement: %v", err)
	}

	return &keyspb.CreateKeyAgreementResponse{}, nil
}

func (h *KeysHandler) CompleteKeyAgreement(
	ctx context.Context,
	req *keyspb.CompleteKeyAgreementRequest,
) (*keyspb.CompleteKeyAgreementResponse, error) {
	if req.GetChainID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chain_id is required")
	}
	if len(req.GetCiphertext()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "ciphertext is required")
	}

	err := h.svc.CompleteKeyAgreement(
		ctx,
		new(big.Int).SetUint64(req.GetChainID()),
		req.GetCiphertext(),
		req.GetBlockNumber(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "complete key agreement: %v", err)
	}

	return &keyspb.CompleteKeyAgreementResponse{}, nil
}

// Payment Spend Keys

func (h *KeysHandler) CreatePaymentSpendKey(
	ctx context.Context,
	req *keyspb.CreatePaymentSpendKeyRequest,
) (*keyspb.PaymentSpendKeyResponse, error) {
	if req.GetChainID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chain_id is required")
	}

	key, err := h.svc.CreatePaymentSpendKey(ctx, new(big.Int).SetUint64(req.GetChainID()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create payment spend key: %v", err)
	}

	return &keyspb.PaymentSpendKeyResponse{
		SecretKey: key.SecretKey.Bytes(),
		PublicKey: key.PublicKey.Bytes(),
	}, nil
}

func (h *KeysHandler) GetPaymentSpendKey(
	ctx context.Context,
	req *keyspb.GetPaymentSpendKeyRequest,
) (*keyspb.PaymentSpendKeyResponse, error) {
	key, err := h.svc.GetPaymentSpendKey(ctx)
	if err != nil {
		if errors.Is(err, service.ErrNoPaymentSpendKeySet) {
			return nil, status.Errorf(codes.NotFound, "get payment spend key: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "get payment spend key: %v", err)
	}

	return &keyspb.PaymentSpendKeyResponse{
		SecretKey: key.SecretKey.Bytes(),
		PublicKey: key.PublicKey.Bytes(),
	}, nil
}

// Self Secret

func (h *KeysHandler) CreateSelfSecret(
	ctx context.Context,
	req *keyspb.CreateSelfSecretRequest,
) (*keyspb.CreateSelfSecretResponse, error) {
	if len(req.GetRFactor()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "r_factor is required")
	}
	if len(req.GetResourceID()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "resource_id is required")
	}

	rFactor := new(big.Int).SetBytes(req.GetRFactor())

	err := h.svc.CreateSelfSecret(ctx, rFactor, req.GetBlockNumber(), req.GetResourceID())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create self secret: %v", err)
	}

	return &keyspb.CreateSelfSecretResponse{}, nil
}

// Enygma

func (h *KeysHandler) GenerateEnygmaSharedSecrets(
	ctx context.Context,
	req *keyspb.GenerateEnygmaSharedSecretsRequest,
) (*keyspb.GenerateEnygmaSharedSecretsResponse, error) {
	if len(req.GetChainIDs()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "chain_ids is required")
	}
	if len(req.GetResourceID()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "resource_id is required")
	}

	chainIDs := make([]*big.Int, len(req.GetChainIDs()))
	for i, id := range req.GetChainIDs() {
		chainIDs[i] = new(big.Int).SetUint64(id)
	}

	senderChainID := new(big.Int).SetUint64(req.GetSenderChainID())
	prevRFactor := new(big.Int).SetBytes(req.GetPrevRFactor())
	resourceID := req.GetResourceID()

	if err := h.svc.CreateSelfSecret(ctx, prevRFactor, req.GetBlockNumber(), resourceID); err != nil {
		return nil, status.Errorf(codes.Internal, "create self secret: %v", err)
	}

	sharedSecrets, err := h.svc.GetEnygmaSharedSecrets(ctx, chainIDs, senderChainID, req.GetBlockNumber(), resourceID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get shared secrets: %v", err)
	}

	secrets, hashSecrets, messageTags, err := h.svc.GenerateEnygmaSharedSecrets(chainIDs, req.GetBlockNumber(), sharedSecrets)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate enygma shared secrets: %v", err)
	}

	secretsBytes := make([][]byte, len(secrets))
	hashSecretsBytes := make([][]byte, len(hashSecrets))
	messageTagsBytes := make([][]byte, len(messageTags))
	for i := range secrets {
		secretsBytes[i] = secrets[i].Bytes()
		hashSecretsBytes[i] = hashSecrets[i].Bytes()
		messageTagsBytes[i] = messageTags[i].Bytes()
	}

	return &keyspb.GenerateEnygmaSharedSecretsResponse{
		Secrets:     secretsBytes,
		HashSecrets: hashSecretsBytes,
		MessageTags: messageTagsBytes,
	}, nil
}

// Helpers

func marshalECDSAKeys(keys []*ecdsa.PrivateKey) []string {
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		result = append(result, hexutil.Encode(crypto.FromECDSA(key))[2:])
	}
	return result
}

func marshalAddresses(addresses domain.AddressList) []string {
	result := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		result = append(result, addr.Hex())
	}
	return result
}
