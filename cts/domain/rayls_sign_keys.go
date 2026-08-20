package domain

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Encryptor interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

type PublicRelayerRaylsSignKeys struct {
	PublicChainKeys  RaylsSignKeyList
	PrivateChainKeys RaylsSignKeyList
}

func (p PublicRelayerRaylsSignKeys) Encrypt(enc Encryptor) (EncryptedPublicRelayerRaylsSignKeys, error) {
	encrPublicKeys, err := p.PublicChainKeys.Encrypt(enc)
	if err != nil {
		return EncryptedPublicRelayerRaylsSignKeys{}, fmt.Errorf("failed to encrypt public chain keys: %w", err)
	}

	encrRaylsSignPrivateKeys, err := p.PrivateChainKeys.Encrypt(enc)
	if err != nil {
		return EncryptedPublicRelayerRaylsSignKeys{}, fmt.Errorf("failed to encrypt private chain keys: %w", err)
	}

	return EncryptedPublicRelayerRaylsSignKeys{
		PublicChainKeys:  encrPublicKeys,
		PrivateChainKeys: encrRaylsSignPrivateKeys,
	}, nil
}

type EncryptedPublicRelayerRaylsSignKeys struct {
	PublicChainKeys  EncryptedRaylsSignKeyList
	PrivateChainKeys EncryptedRaylsSignKeyList
}

func (e EncryptedPublicRelayerRaylsSignKeys) Decrypt(enc Encryptor) (PublicRelayerRaylsSignKeys, error) {
	publicKeys, err := e.PublicChainKeys.Decrypt(enc)
	if err != nil {
		return PublicRelayerRaylsSignKeys{}, fmt.Errorf("failed to decrypt public chain keys: %w", err)
	}

	privateKeys, err := e.PrivateChainKeys.Decrypt(enc)
	if err != nil {
		return PublicRelayerRaylsSignKeys{}, fmt.Errorf("failed to decrypt private chain keys: %w", err)
	}

	return PublicRelayerRaylsSignKeys{
		PublicChainKeys:  publicKeys,
		PrivateChainKeys: privateKeys,
	}, nil
}

type PrivateRelayerRaylsSignKeys struct {
	PrivateHubKeys            RaylsSignKeyList
	PrivateNodeKeys           RaylsSignKeyList
	PrivateHubDvpOperatorKeys RaylsSignKeyList
}

func (p PrivateRelayerRaylsSignKeys) Encrypt(enc Encryptor) (EncryptedPrivateRelayerRaylsSignKeys, error) {
	encrPrivateHubKeys, err := p.PrivateHubKeys.Encrypt(enc)
	if err != nil {
		return EncryptedPrivateRelayerRaylsSignKeys{}, fmt.Errorf("failed to encrypt private hub keys: %w", err)
	}

	encrPrivateHubDvpOperatorKeys, err := p.PrivateHubDvpOperatorKeys.Encrypt(enc)
	if err != nil {
		return EncryptedPrivateRelayerRaylsSignKeys{}, fmt.Errorf(
			"failed to encrypt private hub dvp operator keys: %w",
			err,
		)
	}

	encrRaylsSignPrivateNodeKeys, err := p.PrivateNodeKeys.Encrypt(enc)
	if err != nil {
		return EncryptedPrivateRelayerRaylsSignKeys{}, fmt.Errorf("failed to encrypt private chain keys: %w", err)
	}

	return EncryptedPrivateRelayerRaylsSignKeys{
		PrivateHubKeys:            encrPrivateHubKeys,
		PrivateNodeKeys:           encrRaylsSignPrivateNodeKeys,
		PrivateHubDvpOperatorKeys: encrPrivateHubDvpOperatorKeys,
	}, nil
}

type EncryptedPrivateRelayerRaylsSignKeys struct {
	PrivateHubKeys            EncryptedRaylsSignKeyList
	PrivateNodeKeys           EncryptedRaylsSignKeyList
	PrivateHubDvpOperatorKeys EncryptedRaylsSignKeyList
}

func (e EncryptedPrivateRelayerRaylsSignKeys) Decrypt(enc Encryptor) (PrivateRelayerRaylsSignKeys, error) {
	privateHubKeys, err := e.PrivateHubKeys.Decrypt(enc)
	if err != nil {
		return PrivateRelayerRaylsSignKeys{}, fmt.Errorf("failed to decrypt private hub keys: %w", err)
	}

	privateHubDvpOperatorKeys, err := e.PrivateHubDvpOperatorKeys.Decrypt(enc)
	if err != nil {
		return PrivateRelayerRaylsSignKeys{}, fmt.Errorf("failed to decrypt private hub dvp operator keys: %w", err)
	}

	privateNodeKeys, err := e.PrivateNodeKeys.Decrypt(enc)
	if err != nil {
		return PrivateRelayerRaylsSignKeys{}, fmt.Errorf("failed to decrypt private chain keys: %w", err)
	}

	return PrivateRelayerRaylsSignKeys{
		PrivateHubKeys:            privateHubKeys,
		PrivateNodeKeys:           privateNodeKeys,
		PrivateHubDvpOperatorKeys: privateHubDvpOperatorKeys,
	}, nil
}

type AtomicServiceRaylsSignKeys struct {
	PrivateHubKeys   RaylsSignKeyList
	PrivateChainKeys RaylsSignKeyList
}

func (a AtomicServiceRaylsSignKeys) Encrypt(enc Encryptor) (EncryptedAtomicServiceRaylsSignKeys, error) {
	encrPrivateHubKeys, err := a.PrivateHubKeys.Encrypt(enc)
	if err != nil {
		return EncryptedAtomicServiceRaylsSignKeys{}, fmt.Errorf("failed to encrypt private hub keys: %w", err)
	}

	encrRaylsSignPrivateKeys, err := a.PrivateChainKeys.Encrypt(enc)
	if err != nil {
		return EncryptedAtomicServiceRaylsSignKeys{}, fmt.Errorf("failed to encrypt private chain keys: %w", err)
	}

	return EncryptedAtomicServiceRaylsSignKeys{
		PrivateHubKeys:   encrPrivateHubKeys,
		PrivateChainKeys: encrRaylsSignPrivateKeys,
	}, nil
}

type EncryptedAtomicServiceRaylsSignKeys struct {
	PrivateHubKeys   EncryptedRaylsSignKeyList
	PrivateChainKeys EncryptedRaylsSignKeyList
}

func (e EncryptedAtomicServiceRaylsSignKeys) Decrypt(enc Encryptor) (AtomicServiceRaylsSignKeys, error) {
	privateHubKeys, err := e.PrivateHubKeys.Decrypt(enc)
	if err != nil {
		return AtomicServiceRaylsSignKeys{}, fmt.Errorf("failed to decrypt private hub keys: %w", err)
	}

	privateKeys, err := e.PrivateChainKeys.Decrypt(enc)
	if err != nil {
		return AtomicServiceRaylsSignKeys{}, fmt.Errorf("failed to decrypt private chain keys: %w", err)
	}

	return AtomicServiceRaylsSignKeys{
		PrivateHubKeys:   privateHubKeys,
		PrivateChainKeys: privateKeys,
	}, nil
}

type RaylsSignKeyList []*ecdsa.PrivateKey

func (k RaylsSignKeyList) Encrypt(enc Encryptor) (EncryptedRaylsSignKeyList, error) {
	encrKeys, err := encryptKeySlice(k, enc)
	if err != nil {
		return EncryptedRaylsSignKeyList{}, fmt.Errorf("failed to encrypt list of keys: %w", err)
	}

	return EncryptedRaylsSignKeyList(encrKeys), nil
}

func encryptKeySlice(keys []*ecdsa.PrivateKey, encr Encryptor) ([][]byte, error) {
	encrKeys := [][]byte{}

	for _, key := range keys {
		encrKey, err := encr.Encrypt(crypto.FromECDSA(key))
		if err != nil {
			return nil, fmt.Errorf("encrypting ECDSA key: %w", err)
		}
		encrKeys = append(encrKeys, encrKey)
	}

	return encrKeys, nil
}

type EncryptedRaylsSignKeyList [][]byte

func (e EncryptedRaylsSignKeyList) Decrypt(enc Encryptor) (RaylsSignKeyList, error) {
	keys, err := decryptKeySlice(e, enc)
	if err != nil {
		return RaylsSignKeyList{}, fmt.Errorf("failed to decrypt key list: %w", err)
	}

	return RaylsSignKeyList(keys), nil
}

func decryptKeySlice(encrKeys [][]byte, encr Encryptor) ([]*ecdsa.PrivateKey, error) {
	keys := []*ecdsa.PrivateKey{}

	for _, encrKey := range encrKeys {
		keyBytes, err := encr.Decrypt(encrKey)
		if err != nil {
			return nil, fmt.Errorf("decrypting ECDSA key: %w", err)
		}
		key, err := crypto.ToECDSA(keyBytes)
		if err != nil {
			return nil, fmt.Errorf("parsing ECDSA key from decrypted bytes: %w", err)
		}
		keys = append(keys, key)
	}

	return keys, nil
}

type AddressList []common.Address

type PublicRelayerRaylsSignPublicAddresses struct {
	PublicChainAddresses  AddressList
	PrivateChainAddresses AddressList
}

type PrivateRelayerRaylsSignPublicAddresses struct {
	PrivateHubAddresses            AddressList
	PrivateHubDvpOperatorAddresses AddressList
	PrivateChainAddresses          AddressList
}

type AtomicServiceRaylsSignPublicAddresses struct {
	PrivateHubAddresses   AddressList
	PrivateChainAddresses AddressList
}

type AddressesPair struct {
	FirstChainAddresses  []common.Address
	SecondChainAddresses []common.Address
}
