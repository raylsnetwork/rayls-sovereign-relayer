package service

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cryptography"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/domain"
)

type EncryptKeysService interface {
	GetSharedSecret(ctx context.Context, chainID *big.Int, blockNumber uint64) (domain.SharedSecret, error)
	GetAllSharedSecrets(ctx context.Context, blockNumber uint64) ([]domain.SharedSecret, error)
	GetEnygmaSharedSecrets(ctx context.Context, chainIDs []*big.Int, myChainID *big.Int, blockNumber uint64, resourceID []byte) ([]domain.SharedSecret, error)
	GetEnygmaSharedSelfSecret(ctx context.Context, myChainID *big.Int, blockNumber uint64, resourceID []byte) (domain.SharedSecret, error)
	GetRaylsViewKeyPair(ctx context.Context, blockNum uint64) (domain.RaylsViewKeyPair, error)
}

type EncryptService struct {
	myChainID   *big.Int
	keysService EncryptKeysService
	encryptor   Encryptor
}

func NewEncryptService(myChainID *big.Int, keysService EncryptKeysService, encryptor Encryptor) EncryptService {
	return EncryptService{
		myChainID:   myChainID,
		keysService: keysService,
		encryptor:   encryptor,
	}
}

func (e *EncryptService) GCMEncryptWithNormalMessageTag(
	ctx context.Context,
	blob []byte,
	chainID uint64,
	blockNumber uint64,
	prevBlockHash string,
) ([]byte, string, error) {
	sharedSecret, err := e.keysService.GetSharedSecret(ctx, new(big.Int).SetUint64(chainID), blockNumber)
	if err != nil {
		return nil, "", fmt.Errorf("get shared secret for chain %d: %w", chainID, err)
	}

	// Use receiver's chainID (sharedSecret.ChainId) for message tag generation
	msgTag := generateNormalMessageTag(prevBlockHash, sharedSecret.Secret, sharedSecret.ChainId)

	symKey, err := cryptography.DeriveSymmetricKey(sharedSecret.Secret)
	if err != nil {
		return nil, "", fmt.Errorf("failed to derive symmetric key: %w", err)
	}

	encryptedData, err := cryptography.EncryptGCM(
		symKey,
		blob,
		msgTag,
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to encrypt data: %w", err)
	}

	return encryptedData, common.Bytes2Hex(msgTag), nil
}

func (e *EncryptService) GCMDecryptWithNormalMessageTag(
	ctx context.Context,
	encryptedData []byte,
	messageTag string,
	blockNumber uint64,
	prevBlockHash string,
) ([]byte, error) {
	sharedSecrets, err := e.keysService.GetAllSharedSecrets(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("get shared secrets: %w", err)
	}

	// Generate expected message tag using my chainID (receiver's chainID)
	// The sender used receiver's chainID when creating the tag
	for _, potentialSharedSecret := range sharedSecrets {
		potentialMessageTag := generateNormalMessageTag(prevBlockHash, potentialSharedSecret.Secret, e.myChainID)

		if messageTag != common.Bytes2Hex(potentialMessageTag) {
			continue
		}

		symKey, err := cryptography.DeriveSymmetricKey(potentialSharedSecret.Secret)
		if err != nil {
			return nil, fmt.Errorf("derive symmetric key: %w", err)
		}

		_, plaintext, err := cryptography.DecryptGCM(encryptedData, symKey)
		if err != nil {
			// Tag matched one of our shared secrets but AEAD still failed.
			// This is the "tampered" case — tag collision or stale-key replay.
			if errors.Is(err, cryptography.ErrAuthFailed) {
				return nil, ErrAuthFailed
			}
			return nil, fmt.Errorf("decrypt data: %w", err)
		}

		return plaintext, nil
	}

	// Walked every shared secret we hold; none produced this tag. The message
	// is cryptographically provably not for us.
	return nil, ErrNotForRecipient
}

func (e *EncryptService) GCMEncrypt(
	ctx context.Context,
	blob []byte,
	chainID uint64,
	blockNumber uint64,
) ([]byte, error) {
	sharedSecret, err := e.keysService.GetSharedSecret(ctx, new(big.Int).SetUint64(chainID), blockNumber)
	if err != nil {
		return nil, fmt.Errorf("get shared secret for chain %d: %w", chainID, err)
	}

	symKey, err := cryptography.DeriveSymmetricKey(sharedSecret.Secret)
	if err != nil {
		return nil, fmt.Errorf("derive symmetric key: %w", err)
	}

	associatedData := make([]byte, 16)
	encryptedData, err := cryptography.EncryptGCM(symKey, blob, associatedData)
	if err != nil {
		return nil, fmt.Errorf("encrypt data: %w", err)
	}

	return encryptedData, nil
}

// GCMEncryptEnygma seals an Enygma transfer batch for chainID, mirroring the secret selection of
// GCMDecryptWithEnygmaMessageTag so that seal-key == decrypt-key for every destination:
//   - the sender's own-chain (self/change) batch (chainID == myChainID) is sealed with the rotating
//     per-resource self-secret — the SAME secret used for the on-chain message tag and the ZK proof;
//   - every other destination is sealed with the per-chain pairwise shared secret.
//
// The plain GCMEncrypt path used GetSharedSecret(chainID) for ALL destinations, which for the own-chain
// batch returns the static own-chain pairwise row instead of the self-secret. The decrypt expects the
// self-secret (GetEnygmaSharedSecrets appends the self-secret as the own-chain entry), so the own-chain
// batch failed AEAD (OUTCOME_TAMPERED) and was dropped/un-resyncable.
func (e *EncryptService) GCMEncryptEnygma(
	ctx context.Context,
	blob []byte,
	chainID uint64,
	blockNumber uint64,
	resourceID []byte,
) ([]byte, error) {
	var sharedSecret domain.SharedSecret
	var err error
	if new(big.Int).SetUint64(chainID).Cmp(e.myChainID) == 0 {
		// Self / change batch — use the self-secret (matches the message tag, the ZK proof, and the decrypt).
		sharedSecret, err = e.keysService.GetEnygmaSharedSelfSecret(ctx, e.myChainID, blockNumber, resourceID)
		if err != nil {
			return nil, fmt.Errorf("get enygma self secret for own chain %d: %w", chainID, err)
		}
	} else {
		sharedSecret, err = e.keysService.GetSharedSecret(ctx, new(big.Int).SetUint64(chainID), blockNumber)
		if err != nil {
			return nil, fmt.Errorf("get shared secret for chain %d: %w", chainID, err)
		}
	}

	return e.GCMEncryptWithProvidedSS(blob, sharedSecret.Secret)
}

func (e *EncryptService) GCMDecrypt(
	ctx context.Context,
	encryptedData []byte,
	chainID uint64,
	blockNumber uint64,
) ([]byte, error) {
	sharedSecret, err := e.keysService.GetSharedSecret(ctx, new(big.Int).SetUint64(chainID), blockNumber)
	if err != nil {
		return nil, fmt.Errorf("get shared secret for chain %d: %w", chainID, err)
	}

	symKey, err := cryptography.DeriveSymmetricKey(sharedSecret.Secret)
	if err != nil {
		return nil, fmt.Errorf("derive symmetric key: %w", err)
	}

	_, plaintext, err := cryptography.DecryptGCM(encryptedData, symKey)
	if err != nil {
		// Caller explicitly addressed chainID and we used that chain's key;
		// AEAD failure here means a tag-matched-but-failed-AEAD analogue —
		// stale key, replay, or tampering. Treat as suspicious, not "not for me".
		if errors.Is(err, cryptography.ErrAuthFailed) {
			return nil, ErrAuthFailed
		}
		return nil, fmt.Errorf("decrypt data: %w", err)
	}

	return plaintext, nil
}

func (e *EncryptService) GCMEncryptWithProvidedSS(blob []byte, ss []byte) ([]byte, error) {
	symKey, err := cryptography.DeriveSymmetricKey(ss)
	if err != nil {
		return nil, fmt.Errorf("derive symmetric key: %w", err)
	}
	associatedData := make([]byte, 16)
	encryptedData, err := cryptography.EncryptGCM(symKey, blob, associatedData)
	if err != nil {
		return nil, fmt.Errorf("encrypt data: %w", err)
	}
	return encryptedData, nil
}

func (e *EncryptService) GCMDecryptWithProvidedSS(encryptedData []byte, ss []byte) ([]byte, error) {
	symKey, err := cryptography.DeriveSymmetricKey(ss)
	if err != nil {
		return nil, fmt.Errorf("derive symmetric key: %w", err)
	}
	_, plaintext, err := cryptography.DecryptGCM(encryptedData, symKey)
	if err != nil {
		// On the salt-based DVP path the only "not for me" signal we have is
		// AEAD-fail — ML-KEM decapsulation succeeds for everyone, but produces
		// the right salt only for the intended recipient. Treat AEAD-fail as
		// the normal "not for me" outcome rather than tampering.
		if errors.Is(err, cryptography.ErrAuthFailed) {
			return nil, ErrNotForRecipient
		}
		return nil, fmt.Errorf("decrypt data: %w", err)
	}
	return plaintext, nil
}

func (e *EncryptService) GCMDecryptWithEnygmaMessageTag(
	ctx context.Context,
	encryptedData []byte,
	messageTag []byte,
	blockNumber uint64,
	anonymitySet []*big.Int,
	chainID *big.Int,
	resourceID []byte,
) ([]byte, error) {
	// Fetch only secrets for the anonymity set participants (+ self secret if sender)
	sharedSecrets, err := e.keysService.GetEnygmaSharedSecrets(ctx, anonymitySet, chainID, blockNumber, resourceID)
	if err != nil {
		return nil, fmt.Errorf("get enygma shared secrets: %w", err)
	}

	// Self secret is used for message tag matching when we are the sender
	selfSecret, err := e.keysService.GetEnygmaSharedSelfSecret(ctx, chainID, blockNumber, resourceID)
	var selfSecretPtr *domain.SharedSecret
	if err == nil {
		selfSecretPtr = &selfSecret
	} else if !errors.Is(err, ErrNoApplicableEnygmaSelfSecret) {
		return nil, fmt.Errorf("get enygma self secret: %w", err)
	}

	actualMessageTag := common.Bytes2Hex(messageTag)

	for _, potentialSharedSecret := range sharedSecrets {
		messageTagSecret := potentialSharedSecret.Secret
		if selfSecretPtr != nil && potentialSharedSecret.ChainId.Cmp(selfSecretPtr.ChainId) == 0 {
			messageTagSecret = selfSecretPtr.Secret
		}

		potentialMessageTag, err := generateEnygmaMessageTag(messageTagSecret, blockNumber)
		if err != nil {
			return nil, fmt.Errorf("generate enygma message tag: %w", err)
		}

		if actualMessageTag != common.Bytes2Hex(potentialMessageTag) {
			continue
		}

		symKey, err := cryptography.DeriveSymmetricKey(potentialSharedSecret.Secret)
		if err != nil {
			return nil, fmt.Errorf("derive symmetric key: %w", err)
		}

		_, plaintext, err := cryptography.DecryptGCM(encryptedData, symKey)
		if err != nil {
			// Tag matched but AEAD failed — tampered / replay (see GCMDecryptWithNormalMessageTag).
			if errors.Is(err, cryptography.ErrAuthFailed) {
				return nil, ErrAuthFailed
			}
			return nil, fmt.Errorf("decrypt data: %w", err)
		}

		return plaintext, nil
	}

	// No anonymity-set participant's tag matched — message is not for us.
	return nil, ErrNotForRecipient
}

func (e *EncryptService) EncryptPrivateKey(
	ctx context.Context,
	chainID uint64,
	blockNumber uint64,
) (string, []byte, []byte, uint64, error) {
	pair, err := e.keysService.GetRaylsViewKeyPair(ctx, blockNumber)
	if err != nil {
		return "", nil, nil, 0, fmt.Errorf("get rayls view key pair: %w", err)
	}

	sharedSecret, err := e.keysService.GetSharedSecret(ctx, new(big.Int).SetUint64(chainID), blockNumber)
	if err != nil {
		return "", nil, nil, 0, fmt.Errorf("get shared secret for chain %d: %w", chainID, err)
	}

	symKey, err := cryptography.DeriveSymmetricKey(sharedSecret.Secret)
	if err != nil {
		return "", nil, nil, 0, fmt.Errorf("derive symmetric key: %w", err)
	}

	// JSON-marshal the private key bytes before encrypting.
	// The governance service expects the decrypted payload to be a JSON-encoded byte slice
	// (base64 string), matching the old KOS HTTP behavior which used json.Marshal([]byte).
	privateKeyJSON, err := json.Marshal(pair.RaylsViewPrivateKey.Bytes())
	if err != nil {
		return "", nil, nil, 0, fmt.Errorf("marshal private key: %w", err)
	}

	associatedData := make([]byte, 16)
	encryptedPrivateKey, err := cryptography.EncryptGCM(symKey, privateKeyJSON, associatedData)
	if err != nil {
		return "", nil, nil, 0, fmt.Errorf("encrypt private key: %w", err)
	}

	mac := cryptography.KMAC(sharedSecret.Secret, encryptedPrivateKey)

	return hex.EncodeToString(pair.RaylsViewPublicKey.Bytes()), encryptedPrivateKey, mac, pair.InitialBlock, nil
}

// generateNormalMessageTag creates a message tag for encrypted message identification.
// The shared secret is symmetric - both participants share the same secret. This means
// both participants can regenerate the same message tag. Without additional context,
// there is no way for the receiver to know if a message is FOR them or FROM them.
// That's why we include the receiver's chainID in the tag.
func generateNormalMessageTag(
	prevBlockHash string,
	sharedSecret []byte,
	chainID *big.Int,
) []byte {
	var messageTagData string
	if len(prevBlockHash) != 0 {
		messageTagData = prevBlockHash[2:] + chainID.String() + hex.EncodeToString(sharedSecret)
	} else {
		messageTagData = prevBlockHash + chainID.String() + hex.EncodeToString(sharedSecret)
	}

	return cryptography.HashIt([]byte(messageTagData))[0:16]
}

// poseidonHashTagSeed is the seed value used to compute the hash tag for enygma message tags.
const poseidonHashTagSeed = 12

func generateEnygmaMessageTag(
	sharedSecret []byte,
	blockNumber uint64,
) ([]byte, error) {
	hashTag, err := cryptography.GetPoseidonHash([]*big.Int{big.NewInt(poseidonHashTagSeed)})
	if err != nil {
		return nil, fmt.Errorf("failed to generate hash tag: %w", err)
	}
	msgTag, err := cryptography.GetPoseidonHashModNumber(
		[]*big.Int{
			hashTag,
			big.NewInt(0).SetBytes(sharedSecret),
			new(big.Int).SetUint64(blockNumber),
		},
		cryptography.JubJubPrimeSubGroup,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate message tag: %w", err)
	}

	return msgTag.Bytes(), nil
}
