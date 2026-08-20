package domain

import (
	"crypto/mlkem"
	"fmt"
)

type RaylsViewKeyPair struct {
	InitialBlock        uint64
	RaylsViewPrivateKey *mlkem.DecapsulationKey768
	RaylsViewPublicKey  *mlkem.EncapsulationKey768
}

func (d RaylsViewKeyPair) Encrypt(enc Encryptor) (EncryptedRaylsViewKeyPair, error) {
	privateKeyBytes := d.RaylsViewPrivateKey.Bytes()
	publicKeyBytes := d.RaylsViewPublicKey.Bytes()

	encrPrivateKey, err := enc.Encrypt(privateKeyBytes)
	if err != nil {
		return EncryptedRaylsViewKeyPair{}, fmt.Errorf("failed to encrypt rayls view private key: %w", err)
	}

	return EncryptedRaylsViewKeyPair{
		InitialBlock:                 d.InitialBlock,
		EncryptedRaylsViewPrivateKey: encrPrivateKey,
		RaylsViewPublicKey:           publicKeyBytes,
	}, nil
}

type EncryptedRaylsViewKeyPair struct {
	InitialBlock                 uint64
	EncryptedRaylsViewPrivateKey []byte
	RaylsViewPublicKey           []byte
}

func (e EncryptedRaylsViewKeyPair) Decrypt(enc Encryptor) (RaylsViewKeyPair, error) {
	privateKeyBytes, err := enc.Decrypt(e.EncryptedRaylsViewPrivateKey)
	if err != nil {
		return RaylsViewKeyPair{}, fmt.Errorf("failed to decrypt rayls view private key: %w", err)
	}

	privateKey, err := mlkem.NewDecapsulationKey768(privateKeyBytes)
	if err != nil {
		return RaylsViewKeyPair{}, fmt.Errorf("failed to create private key: %w", err)
	}

	publicKey, err := mlkem.NewEncapsulationKey768(e.RaylsViewPublicKey)
	if err != nil {
		return RaylsViewKeyPair{}, fmt.Errorf("failed to create public key: %w", err)
	}

	return RaylsViewKeyPair{
		InitialBlock:        e.InitialBlock,
		RaylsViewPrivateKey: privateKey,
		RaylsViewPublicKey:  publicKey,
	}, nil
}
