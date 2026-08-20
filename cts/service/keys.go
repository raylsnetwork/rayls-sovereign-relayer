package service

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cryptography"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/domain"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/keygen"
	"github.com/raylsnetwork/rayls-sovereign-relayer/faultinjector"
)

const defaultSignKeyCount = 5

type Encryptor interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

type RaylsViewKeysRepository interface {
	Create(ctx context.Context, key domain.EncryptedRaylsViewKeyPair) error
	GetForBlockNumber(ctx context.Context, blockNumber uint64) (domain.EncryptedRaylsViewKeyPair, error)
	DeleteByPublicKey(ctx context.Context, key string) error
}

type RaylsSignKeysRepository interface {
	CreatePublicRelayerRaylsSignKeys(ctx context.Context, keys domain.EncryptedPublicRelayerRaylsSignKeys) error
	CreatePrivateRelayerRaylsSignKeys(ctx context.Context, keys domain.EncryptedPrivateRelayerRaylsSignKeys) error
	CreateAtomicServiceRaylsSignKeys(ctx context.Context, keys domain.EncryptedAtomicServiceRaylsSignKeys) error
	GetPublicRelayerRaylsSignKeys(ctx context.Context) (domain.EncryptedPublicRelayerRaylsSignKeys, error)
	GetPrivateRelayerRaylsSignKeys(ctx context.Context) (domain.EncryptedPrivateRelayerRaylsSignKeys, error)
	GetAtomicServiceRaylsSignKeys(ctx context.Context) (domain.EncryptedAtomicServiceRaylsSignKeys, error)
}

type SharedSecretsRepository interface {
	Create(ctx context.Context, secret domain.EncryptedSharedSecret) error
	GetAll(ctx context.Context, blockNumber uint64) ([]domain.EncryptedSharedSecret, error)
	GetByChainIds(ctx context.Context, chainIds []*big.Int, blockNumber uint64) ([]domain.EncryptedSharedSecret, error)
	GetByChainId(ctx context.Context, chainId string, blockNumber uint64) (domain.EncryptedSharedSecret, error)
}

type EnygmaSelfSecretsRepository interface {
	Create(ctx context.Context, secret domain.EncryptedEnygmaSelfSecret) error
	GetByBlockNumberAndResource(ctx context.Context, blockNumber uint64, resourceID []byte) (domain.EncryptedEnygmaSelfSecret, error)
}

type PaymentSpendKeysRepository interface {
	Create(ctx context.Context, paymentSpendKeys domain.EncryptedPaymentSpendKeys) error
	Get(ctx context.Context) (domain.EncryptedPaymentSpendKeys, error)
}

type KeysService struct {
	encryptor Encryptor

	raylsViewKeysRepository     RaylsViewKeysRepository
	ecdsaKeysRepository         RaylsSignKeysRepository
	sharedSecretsRepository     SharedSecretsRepository
	enygmaSelfSecretsRepository EnygmaSelfSecretsRepository
	paymentSpendKeysRepository  PaymentSpendKeysRepository
}

func NewKeysService(
	encryptor Encryptor,
	viewRepo RaylsViewKeysRepository,
	ecdsaRepo RaylsSignKeysRepository,
	sharedSecretRepo SharedSecretsRepository,
	enygmaSelfSecretRepo EnygmaSelfSecretsRepository,
	paymentSpendRepo PaymentSpendKeysRepository,
) *KeysService {
	return &KeysService{
		encryptor: encryptor,

		raylsViewKeysRepository:     viewRepo,
		ecdsaKeysRepository:         ecdsaRepo,
		sharedSecretsRepository:     sharedSecretRepo,
		enygmaSelfSecretsRepository: enygmaSelfSecretRepo,
		paymentSpendKeysRepository:  paymentSpendRepo,
	}
}

func (s *KeysService) CreateRaylsViewKeyPair(ctx context.Context, initialBlock uint64) (domain.RaylsViewKeyPair, error) {
	pair, err := keygen.GenerateRaylsViewKeys()
	if err != nil {
		return domain.RaylsViewKeyPair{}, NewKeyServiceError(err.Error())
	}
	pair.InitialBlock = initialBlock

	// Fault injection: simulate KMS (external) being unavailable or slow
	// before the key material is sent for encryption.
	if faultErr := faultinjector.Check("kos.service.KeysService.CreateRaylsViewKeyPair.before_kms_encrypt"); faultErr != nil {
		return domain.RaylsViewKeyPair{}, fmt.Errorf("fault injection at before_kms_encrypt: %w", faultErr)
	}

	encrPair, err := pair.Encrypt(s.encryptor)
	if err != nil {
		return domain.RaylsViewKeyPair{}, fmt.Errorf("failed to encrypt rayls view key pair: %w", err)
	}

	// Fault injection: KMS encrypted the keys but the DB write hasn't run
	// yet — tests the compensation case where the KMS holds material the DB
	// never persists. The "kos." namespace is preserved in FI keys for
	// backward compatibility with the resilience e2e suite, which hard-codes
	// these identifiers even though the host package is now cts.
	if faultErr := faultinjector.Check("kos.service.KeysService.CreateRaylsViewKeyPair.before_db_insert"); faultErr != nil {
		return domain.RaylsViewKeyPair{}, fmt.Errorf("fault injection at before_db_insert: %w", faultErr)
	}

	err = s.raylsViewKeysRepository.Create(ctx, encrPair)
	if err != nil {
		return domain.RaylsViewKeyPair{}, NewKeyRepositoryError(err.Error())
	}

	// Fault injection: the key is durably persisted but the HTTP response
	// hasn't been sent — the client sees a request failure and may retry,
	// exercising the "already exists" / idempotency path on next attempt.
	if faultErr := faultinjector.Check("kos.service.KeysService.CreateRaylsViewKeyPair.after_db_insert"); faultErr != nil {
		return domain.RaylsViewKeyPair{}, fmt.Errorf("fault injection at after_db_insert: %w", faultErr)
	}

	return pair, nil
}

func (s *KeysService) DeleteRaylsViewKeyPair(ctx context.Context, publicKey string) error {
	err := s.raylsViewKeysRepository.DeleteByPublicKey(ctx, publicKey)
	if err != nil {
		return NewKeyRepositoryError(err.Error())
	}

	return nil
}

func (s *KeysService) GetRaylsViewKeyPair(ctx context.Context, blockNum uint64) (domain.RaylsViewKeyPair, error) {
	encrPair, err := s.raylsViewKeysRepository.GetForBlockNumber(ctx, blockNum)
	if err != nil {
		switch {
		case errors.Is(err, ErrNoRaylsViewKeysSet):
			return domain.RaylsViewKeyPair{}, fmt.Errorf("no rayls view keys set: %w", err)
		case errors.Is(err, ErrNoApplicableRaylsViewKeys):
			return domain.RaylsViewKeyPair{}, fmt.Errorf("no applicable rayls view keys for block number: %w", err)
		default:
			return domain.RaylsViewKeyPair{}, NewKeyRepositoryError(err.Error())
		}
	}

	pair, err := encrPair.Decrypt(s.encryptor)
	if err != nil {
		return domain.RaylsViewKeyPair{}, fmt.Errorf("failed to decrypt rayls view key pair: %w", err)
	}

	return pair, nil
}

// RecoverViewSalt performs ML-KEM-768 decapsulation against the view key pair
// valid for blockNumber, recovering the shared secret the sender used as a salt.
// The result is reduced mod JubJubPrimeGroup so it can be used directly as a
// Poseidon/BabyJubJub field element and as the symmetric secret for
// DecryptWithoutFPWithSS — matching cryptography.GenerateSalt on the sender side.
//
// The view secret key never leaves this process.
func (s *KeysService) RecoverViewSalt(ctx context.Context, blockNumber uint64, ctxt []byte) (*big.Int, error) {
	pair, err := s.GetRaylsViewKeyPair(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("get view key pair: %w", err)
	}

	sharedSecret, err := pair.RaylsViewPrivateKey.Decapsulate(ctxt)
	if err != nil {
		return nil, fmt.Errorf("ml-kem decapsulation: %w", err)
	}

	salt := new(big.Int).SetBytes(sharedSecret)
	salt.Mod(salt, cryptography.JubJubPrimeGroup)
	return salt, nil
}

func (s *KeysService) GetSharedSecret(ctx context.Context, chainID *big.Int, blockNumber uint64) (domain.SharedSecret, error) {
	encrSharedSecret, err := s.sharedSecretsRepository.GetByChainId(ctx, chainID.String(), blockNumber)
	if err != nil {
		return domain.SharedSecret{}, fmt.Errorf("failed to get shared secret: %w", err)
	}

	sharedSecret, err := encrSharedSecret.Decrypt(s.encryptor)
	if err != nil {
		return domain.SharedSecret{}, fmt.Errorf("failed to decrypt shared secret: %w", err)
	}

	return sharedSecret, nil
}

func (s *KeysService) GetAllSharedSecrets(ctx context.Context, blockNumber uint64) ([]domain.SharedSecret, error) {
	encrSharedSecrets, err := s.sharedSecretsRepository.GetAll(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get shared secrets: %w", err)
	}

	sharedSecrets := make([]domain.SharedSecret, 0)
	for _, encrSharedSecret := range encrSharedSecrets {
		sharedSecret, err := encrSharedSecret.Decrypt(s.encryptor)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt shared secret: %w", err)
		}
		sharedSecrets = append(sharedSecrets, sharedSecret)
	}

	return sharedSecrets, nil
}

func (s *KeysService) GetEnygmaSharedSelfSecret(
	ctx context.Context,
	myChainID *big.Int,
	blockNumber uint64,
	resourceID []byte,
) (domain.SharedSecret, error) {
	encrSelfSecret, err := s.enygmaSelfSecretsRepository.GetByBlockNumberAndResource(ctx, blockNumber, resourceID)
	if err != nil {
		return domain.SharedSecret{}, fmt.Errorf("failed to get self secret: %w", err)
	}
	selfSecret, err := encrSelfSecret.Decrypt(s.encryptor)
	if err != nil {
		return domain.SharedSecret{}, fmt.Errorf("failed to decrypt self secret: %w", err)
	}

	return domain.SharedSecret{
		ChainId:      myChainID,
		Secret:       selfSecret.Secret,
		InitialBlock: blockNumber,
	}, nil
}

func (s *KeysService) GenerateEnygmaSharedSecrets(
	chainIDs []*big.Int,
	blockNumber uint64,
	sharedSecrets []domain.SharedSecret,
) ([]*big.Int, []*big.Int, []*big.Int, error) {
	// Build lookup map for O(1) access by chainId
	sharedSecretsMap := make(map[string]domain.SharedSecret, len(sharedSecrets))
	for _, sharedSecret := range sharedSecrets {
		sharedSecretsMap[sharedSecret.ChainId.String()] = sharedSecret
	}

	secrets := make([]*big.Int, len(chainIDs))
	hashSecrets := make([]*big.Int, len(chainIDs))
	messageTags := make([]*big.Int, len(chainIDs))

	for i, chainID := range chainIDs {
		sharedKey, exists := sharedSecretsMap[chainID.String()]
		if !exists {
			return nil, nil, nil, fmt.Errorf("missing shared secret for participant chainId=%s at index %d", chainID.String(), i)
		}

		hashedSecret, err := cryptography.GetPoseidonHashModNumber(
			[]*big.Int{
				new(big.Int).SetBytes(sharedKey.Secret),
				new(big.Int).SetBytes(sharedKey.Secret),
			},
			cryptography.JubJubPrimeSubGroup,
		)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to compute secret hash for chainId=%s: %w", chainID.String(), err)
		}

		messageTag, err := generateEnygmaMessageTag(sharedKey.Secret, blockNumber)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to generate message tag for chainId=%s: %w", chainID.String(), err)
		}

		hashSecrets[i] = hashedSecret
		secrets[i] = new(big.Int).SetBytes(sharedKey.Secret)
		messageTags[i] = new(big.Int).SetBytes(messageTag)
	}

	return secrets, hashSecrets, messageTags, nil
}

func (s *KeysService) GetEnygmaSharedSecret(
	ctx context.Context,
	chainID *big.Int,
	myChainID *big.Int,
	blockNumber uint64,
	resourceID []byte,
) (domain.SharedSecret, error) {
	if chainID.Cmp(myChainID) == 0 {
		return s.GetEnygmaSharedSelfSecret(ctx, myChainID, blockNumber, resourceID)
	}

	return s.GetSharedSecret(ctx, chainID, blockNumber)
}

func (s *KeysService) GetEnygmaSharedSecrets(
	ctx context.Context,
	chainIDs []*big.Int,
	myChainID *big.Int,
	blockNumber uint64,
	resourceID []byte,
) ([]domain.SharedSecret, error) {
	otherChainIDs := make([]*big.Int, 0)
	for _, chainID := range chainIDs {
		if chainID.Cmp(myChainID) != 0 {
			otherChainIDs = append(otherChainIDs, chainID)
		}
	}

	encrSharedSecrets, err := s.sharedSecretsRepository.GetByChainIds(ctx, otherChainIDs, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get shared secrets: %w", err)
	}

	sharedSecrets := make([]domain.SharedSecret, 0)
	var sharedSecret domain.SharedSecret
	for _, encrSharedSecret := range encrSharedSecrets {
		sharedSecret, err = encrSharedSecret.Decrypt(s.encryptor)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt shared secret: %w", err)
		}
		sharedSecrets = append(sharedSecrets, sharedSecret)
	}

	sharedSelfSecret, err := s.GetEnygmaSharedSelfSecret(ctx, myChainID, blockNumber, resourceID)
	if err != nil {
		//  Self secret is available only for the "sender" side. If we're just a receiver, we don't need it.
		if errors.Is(err, ErrNoApplicableEnygmaSelfSecret) {
			return sharedSecrets, nil
		}
		return nil, fmt.Errorf("failed to get self secret: %w", err)
	}

	sharedSecrets = append(sharedSecrets, sharedSelfSecret)

	return sharedSecrets, nil
}

// GenerateKeyAgreement generates ciphertext and shared secret from the peer's public key (no DB write).
func (s *KeysService) GenerateKeyAgreement(
	chainID *big.Int,
	publicKey []byte,
) (ciphertext []byte, sharedSecret []byte, digest []byte, err error) {
	ciphertext, sharedSecret, err = keygen.GenerateSharedSecret(publicKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate shared secret: %w", err)
	}

	// There is a small probability that the shared secret is greater than the prime group,
	// this would cause troubles with the Enygma flow.
	// So, we must ensure the shared secret is in the prime group.
	for big.NewInt(0).SetBytes(sharedSecret).Cmp(cryptography.JubJubPrimeSubGroup) == 1 {
		ciphertext, sharedSecret, err = keygen.GenerateSharedSecret(publicKey)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to generate shared secret: %w", err)
		}
	}

	digest, err = keygen.GenerateKeyDigest(sharedSecret)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate shared secret: %w", err)
	}

	return ciphertext, sharedSecret, digest, nil
}

func (s *KeysService) CreateKeyAgreement(ctx context.Context, chainID *big.Int, sharedSecret []byte, blockNum uint64) error {
	sharedSecretModel := domain.SharedSecret{
		ChainId:      chainID,
		Secret:       sharedSecret,
		InitialBlock: blockNum,
	}

	// Fault injection: simulate KMS unavailable before shared-secret encrypt.
	if faultErr := faultinjector.Check("kos.service.KeysService.CreateKeyAgreement.before_kms_encrypt"); faultErr != nil {
		return fmt.Errorf("fault injection at before_kms_encrypt: %w", faultErr)
	}

	encrSharedSecret, err := sharedSecretModel.Encrypt(s.encryptor)
	if err != nil {
		return fmt.Errorf("failed to encrypt shared secret: %w", err)
	}

	// Fault injection: KMS encrypted the shared secret but DB persist hasn't
	// run — tests recovery when the initiator must re-derive the secret on
	// retry.
	if faultErr := faultinjector.Check("kos.service.KeysService.CreateKeyAgreement.before_db_insert"); faultErr != nil {
		return fmt.Errorf("fault injection at before_db_insert: %w", faultErr)
	}

	err = s.sharedSecretsRepository.Create(ctx, encrSharedSecret)
	if err != nil {
		return fmt.Errorf("failed to create shared secret: %w", err)
	}

	// Fault injection: shared secret durably persisted but HTTP response
	// lost — relayer retries CreateKeyAgreement; the repository must handle
	// duplicate inserts idempotently.
	if faultErr := faultinjector.Check("kos.service.KeysService.CreateKeyAgreement.after_db_insert"); faultErr != nil {
		return fmt.Errorf("fault injection at after_db_insert: %w", faultErr)
	}

	return nil
}

func (s *KeysService) CompleteKeyAgreement(ctx context.Context, chainID *big.Int, ciphertext []byte, blockNum uint64) error {
	viewKeyPair, err := s.GetRaylsViewKeyPair(ctx, blockNum)
	if err != nil {
		return fmt.Errorf("failed to get view key pair: %w", err)
	}

	sharedSecret, err := keygen.RecoverSharedSecret(viewKeyPair.RaylsViewPrivateKey, ciphertext)
	if err != nil {
		return fmt.Errorf("failed to recover shared secret: %w", err)
	}

	sharedSecretModel := domain.SharedSecret{
		ChainId:      chainID,
		Secret:       sharedSecret,
		InitialBlock: blockNum,
	}

	encrSharedSecret, err := sharedSecretModel.Encrypt(s.encryptor)
	if err != nil {
		return fmt.Errorf("failed to encrypt shared secret: %w", err)
	}

	err = s.sharedSecretsRepository.Create(ctx, encrSharedSecret)
	if err != nil {
		return fmt.Errorf("failed to create shared secret: %w", err)
	}

	return nil
}

func (s *KeysService) CreateSelfSecret(ctx context.Context, rFactor *big.Int, blockNumber uint64, resourceID []byte) error {
	spendKey, err := s.GetPaymentSpendKey(ctx)
	if err != nil {
		return fmt.Errorf("failed to get view key pair: %w", err)
	}

	selfSecret, err := cryptography.GetPoseidonHashModNumber(
		[]*big.Int{rFactor, spendKey.SecretKey},
		cryptography.JubJubPrimeSubGroup,
	)
	if err != nil {
		return fmt.Errorf("failed to compute self secret: %w", err)
	}

	selfSecretModel := domain.EnygmaSelfSecret{
		Secret:       selfSecret.Bytes(),
		InitialBlock: blockNumber,
		ResourceID:   resourceID,
	}

	encrypted, err := selfSecretModel.Encrypt(s.encryptor)
	if err != nil {
		return fmt.Errorf("failed to encrypt self secret: %w", err)
	}

	err = s.enygmaSelfSecretsRepository.Create(ctx, encrypted)
	if err != nil {
		return fmt.Errorf("failed to store self secret: %w", err)
	}

	return nil
}

func (s *KeysService) GetSelfSecret(ctx context.Context, blockNumber uint64, resourceID []byte) (domain.EnygmaSelfSecret, error) {
	encrypted, err := s.enygmaSelfSecretsRepository.GetByBlockNumberAndResource(ctx, blockNumber, resourceID)
	if err != nil {
		return domain.EnygmaSelfSecret{}, fmt.Errorf("failed to get self secret: %w", err)
	}

	secret, err := encrypted.Decrypt(s.encryptor)
	if err != nil {
		return domain.EnygmaSelfSecret{}, fmt.Errorf("failed to decrypt self secret: %w", err)
	}

	return secret, nil
}

func (s *KeysService) CreatePublicRelayerRaylsSignKeys(ctx context.Context) (domain.PublicRelayerRaylsSignKeys, error) {
	const raylsSignKeyCount = 5
	publicChainKeys, err := keygen.GenerateRaylsSignKeys(raylsSignKeyCount)
	if err != nil {
		return domain.PublicRelayerRaylsSignKeys{}, fmt.Errorf("failed to generate public chain keys: %w", err)
	}

	privateChainKeys, err := keygen.GenerateRaylsSignKeys(raylsSignKeyCount)
	if err != nil {
		return domain.PublicRelayerRaylsSignKeys{}, fmt.Errorf("failed to generate private chain keys: %w", err)
	}

	keysPair := domain.PublicRelayerRaylsSignKeys{
		PublicChainKeys:  publicChainKeys,
		PrivateChainKeys: privateChainKeys,
	}

	// Fault injection: simulate KMS unavailable before sign-key material
	// is encrypted.
	if faultErr := faultinjector.Check("kos.service.KeysService.CreatePublicRelayerRaylsSignKeys.before_kms_encrypt"); faultErr != nil {
		return domain.PublicRelayerRaylsSignKeys{}, fmt.Errorf("fault injection at before_kms_encrypt: %w", faultErr)
	}

	encrKeyPair, err := keysPair.Encrypt(s.encryptor)
	if err != nil {
		return domain.PublicRelayerRaylsSignKeys{}, fmt.Errorf("failed to encrypt keys pair: %w", err)
	}

	// Fault injection: KMS encrypted the keys but the DB write hasn't run.
	if faultErr := faultinjector.Check("kos.service.KeysService.CreatePublicRelayerRaylsSignKeys.before_db_insert"); faultErr != nil {
		return domain.PublicRelayerRaylsSignKeys{}, fmt.Errorf("fault injection at before_db_insert: %w", faultErr)
	}

	err = s.ecdsaKeysRepository.CreatePublicRelayerRaylsSignKeys(ctx, encrKeyPair)
	if err != nil {
		return domain.PublicRelayerRaylsSignKeys{}, NewKeyRepositoryError(err.Error())
	}

	// Fault injection: keys durably persisted but the HTTP response hasn't
	// been sent — exercises the idempotency path when public-relayer retries
	// its startup key bootstrap.
	if faultErr := faultinjector.Check("kos.service.KeysService.CreatePublicRelayerRaylsSignKeys.after_db_insert"); faultErr != nil {
		return domain.PublicRelayerRaylsSignKeys{}, fmt.Errorf("fault injection at after_db_insert: %w", faultErr)
	}

	return keysPair, nil
}

func (s *KeysService) CreatePrivateRelayerRaylsSignKeys(ctx context.Context) (domain.PrivateRelayerRaylsSignKeys, error) {
	privateHubKeys, err := keygen.GenerateRaylsSignKeys(defaultSignKeyCount)
	if err != nil {
		return domain.PrivateRelayerRaylsSignKeys{}, fmt.Errorf("failed to generate private hub keys: %w", err)
	}

	privateHubDvpOperatorKeys, err := keygen.GenerateRaylsSignKeys(1)
	if err != nil {
		return domain.PrivateRelayerRaylsSignKeys{}, fmt.Errorf(
			"failed to generate private hub dvp operator keys: %w",
			err,
		)
	}

	privateNodeKeys, err := keygen.GenerateRaylsSignKeys(defaultSignKeyCount)
	if err != nil {
		return domain.PrivateRelayerRaylsSignKeys{}, fmt.Errorf("failed to generate private chain keys: %w", err)
	}

	keysPair := domain.PrivateRelayerRaylsSignKeys{
		PrivateHubKeys:            privateHubKeys,
		PrivateHubDvpOperatorKeys: privateHubDvpOperatorKeys,
		PrivateNodeKeys:           privateNodeKeys,
	}

	encrKeyPair, err := keysPair.Encrypt(s.encryptor)
	if err != nil {
		return domain.PrivateRelayerRaylsSignKeys{}, fmt.Errorf("failed to encrypt keys pair: %w", err)
	}

	err = s.ecdsaKeysRepository.CreatePrivateRelayerRaylsSignKeys(ctx, encrKeyPair)
	if err != nil {
		return domain.PrivateRelayerRaylsSignKeys{}, NewKeyRepositoryError(err.Error())
	}

	return keysPair, nil
}

func (s *KeysService) GetPublicRelayerRaylsSignKeys(ctx context.Context) (domain.PublicRelayerRaylsSignKeys, error) {
	encrKeys, err := s.ecdsaKeysRepository.GetPublicRelayerRaylsSignKeys(ctx)
	if err != nil {
		if errors.Is(err, ErrNoRaylsSignKeys) {
			return domain.PublicRelayerRaylsSignKeys{}, fmt.Errorf("failed to get public relayer rayls sign keys: %w", err)
		} else {
			return domain.PublicRelayerRaylsSignKeys{}, NewKeyRepositoryError(err.Error())
		}
	}

	keys, err := encrKeys.Decrypt(s.encryptor)
	if err != nil {
		return domain.PublicRelayerRaylsSignKeys{}, fmt.Errorf("failed to decrypt keys pair: %w", err)
	}

	return keys, nil
}

func (s *KeysService) GetPrivateRelayerRaylsSignKeys(ctx context.Context) (domain.PrivateRelayerRaylsSignKeys, error) {
	encrKeys, err := s.ecdsaKeysRepository.GetPrivateRelayerRaylsSignKeys(ctx)
	if err != nil {
		if errors.Is(err, ErrNoRaylsSignKeys) {
			return domain.PrivateRelayerRaylsSignKeys{}, fmt.Errorf("failed to get private relayer rayls sign keys: %w", err)
		} else {
			return domain.PrivateRelayerRaylsSignKeys{}, NewKeyRepositoryError(err.Error())
		}
	}

	keys, err := encrKeys.Decrypt(s.encryptor)
	if err != nil {
		return domain.PrivateRelayerRaylsSignKeys{}, fmt.Errorf("failed to decrypt keys pair: %w", err)
	}

	return keys, nil
}

func (s *KeysService) GetPublicRelayerRaylsSignPublicAddresses(ctx context.Context) (domain.PublicRelayerRaylsSignPublicAddresses, error) {
	encrKeys, err := s.ecdsaKeysRepository.GetPublicRelayerRaylsSignKeys(ctx)
	if err != nil {
		if errors.Is(err, ErrNoRaylsSignKeys) {
			return domain.PublicRelayerRaylsSignPublicAddresses{}, fmt.Errorf("failed to get public relayer sign keys for addresses: %w", err)
		}
		return domain.PublicRelayerRaylsSignPublicAddresses{}, NewKeyRepositoryError(err.Error())
	}

	keys, err := encrKeys.Decrypt(s.encryptor)
	if err != nil {
		return domain.PublicRelayerRaylsSignPublicAddresses{}, fmt.Errorf("failed to decrypt keys pair: %w", err)
	}

	return domain.PublicRelayerRaylsSignPublicAddresses{
		PublicChainAddresses:  privateKeyListToAddressList(keys.PublicChainKeys),
		PrivateChainAddresses: privateKeyListToAddressList(keys.PrivateChainKeys),
	}, nil
}

func (s *KeysService) GetPrivateRelayerRaylsSignPublicAddresses(ctx context.Context) (domain.PrivateRelayerRaylsSignPublicAddresses, error) {
	encrKeys, err := s.ecdsaKeysRepository.GetPrivateRelayerRaylsSignKeys(ctx)
	if err != nil {
		if errors.Is(err, ErrNoRaylsSignKeys) {
			return domain.PrivateRelayerRaylsSignPublicAddresses{}, fmt.Errorf("failed to get private relayer sign keys for addresses: %w", err)
		}
		return domain.PrivateRelayerRaylsSignPublicAddresses{}, NewKeyRepositoryError(err.Error())
	}

	keys, err := encrKeys.Decrypt(s.encryptor)
	if err != nil {
		return domain.PrivateRelayerRaylsSignPublicAddresses{}, fmt.Errorf("failed to decrypt keys pair: %w", err)
	}

	return domain.PrivateRelayerRaylsSignPublicAddresses{
		PrivateHubAddresses:            privateKeyListToAddressList(keys.PrivateHubKeys),
		PrivateHubDvpOperatorAddresses: privateKeyListToAddressList(keys.PrivateHubDvpOperatorKeys),
		PrivateChainAddresses:          privateKeyListToAddressList(keys.PrivateNodeKeys),
	}, nil
}

func (s *KeysService) GetAtomicServiceRaylsSignPublicAddresses(ctx context.Context) (domain.AtomicServiceRaylsSignPublicAddresses, error) {
	encrKeys, err := s.ecdsaKeysRepository.GetAtomicServiceRaylsSignKeys(ctx)
	if err != nil {
		if errors.Is(err, ErrNoRaylsSignKeys) {
			return domain.AtomicServiceRaylsSignPublicAddresses{}, fmt.Errorf("failed to get atomic service sign keys for addresses: %w", err)
		}
		return domain.AtomicServiceRaylsSignPublicAddresses{}, NewKeyRepositoryError(err.Error())
	}

	keys, err := encrKeys.Decrypt(s.encryptor)
	if err != nil {
		return domain.AtomicServiceRaylsSignPublicAddresses{}, fmt.Errorf("failed to decrypt keys pair: %w", err)
	}

	return domain.AtomicServiceRaylsSignPublicAddresses{
		PrivateHubAddresses:   privateKeyListToAddressList(keys.PrivateHubKeys),
		PrivateChainAddresses: privateKeyListToAddressList(keys.PrivateChainKeys),
	}, nil
}

func privateKeyListToAddressList(keyList []*ecdsa.PrivateKey) domain.AddressList {
	addresses := make(domain.AddressList, 0, len(keyList))
	for _, key := range keyList {
		addresses = append(addresses, crypto.PubkeyToAddress(key.PublicKey))
	}
	return addresses
}

func (s *KeysService) CreatePaymentSpendKey(ctx context.Context, chainId *big.Int) (domain.PaymentSpendKeys, error) {
	secretKey, err := keygen.GenerateRandomModJubJubPrimeSubGroupWithChainId(chainId)
	if err != nil {
		return domain.PaymentSpendKeys{}, fmt.Errorf("failed to generate payment spend secret key: %w", err)
	}

	publicKey, err := keygen.GetPaymentSpendPublicKeyFromSpendSecretKey(secretKey)
	if err != nil {
		return domain.PaymentSpendKeys{}, fmt.Errorf("failed to compute payment spend public key: %w", err)
	}

	keys := domain.PaymentSpendKeys{
		SecretKey: secretKey,
		PublicKey: publicKey,
	}

	// Fault injection: simulate KMS unavailable before payment-spend keys
	// are encrypted.
	if faultErr := faultinjector.Check("kos.service.KeysService.CreatePaymentSpendKey.before_kms_encrypt"); faultErr != nil {
		return domain.PaymentSpendKeys{}, fmt.Errorf("fault injection at before_kms_encrypt: %w", faultErr)
	}

	encrKeys, err := keys.Encrypt(s.encryptor)
	if err != nil {
		return domain.PaymentSpendKeys{}, fmt.Errorf("failed to encrypt keys: %w", err)
	}

	// Fault injection: KMS encrypted the keys but DB persist hasn't run.
	if faultErr := faultinjector.Check("kos.service.KeysService.CreatePaymentSpendKey.before_db_insert"); faultErr != nil {
		return domain.PaymentSpendKeys{}, fmt.Errorf("fault injection at before_db_insert: %w", faultErr)
	}

	err = s.paymentSpendKeysRepository.Create(ctx, encrKeys)
	if err != nil {
		return domain.PaymentSpendKeys{}, fmt.Errorf("failed to store payment spend keys: %w", err)
	}

	// Fault injection: keys durably persisted but HTTP response lost.
	if faultErr := faultinjector.Check("kos.service.KeysService.CreatePaymentSpendKey.after_db_insert"); faultErr != nil {
		return domain.PaymentSpendKeys{}, fmt.Errorf("fault injection at after_db_insert: %w", faultErr)
	}

	return keys, nil
}

func (s *KeysService) GetPaymentSpendKey(ctx context.Context) (domain.PaymentSpendKeys, error) {
	encrKeys, err := s.paymentSpendKeysRepository.Get(ctx)
	if err != nil {
		if errors.Is(err, ErrNoPaymentSpendKeySet) {
			return domain.PaymentSpendKeys{}, fmt.Errorf("failed to get payment spend keys: %w", err)
		} else {
			return domain.PaymentSpendKeys{}, NewKeyRepositoryError(err.Error())
		}
	}

	keys, err := encrKeys.Decrypt(s.encryptor)
	if err != nil {
		return domain.PaymentSpendKeys{}, fmt.Errorf("failed to decrypt keys: %w", err)
	}

	return keys, nil
}
