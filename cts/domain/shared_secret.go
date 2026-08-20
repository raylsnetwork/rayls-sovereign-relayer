package domain

import (
	"fmt"
	"math/big"
)

type SharedSecret struct {
	ChainId      *big.Int // Chain ID this secret is associated with
	Secret       []byte   // The actual shared secret value
	InitialBlock uint64      // Block number when the secret is valid
}
type EncryptedSharedSecret struct {
	ChainId         string
	EncryptedSecret []byte
	InitialBlock    uint64
}

func (s SharedSecret) Encrypt(enc Encryptor) (EncryptedSharedSecret, error) {
	encrSecret, err := enc.Encrypt(s.Secret)
	if err != nil {
		return EncryptedSharedSecret{}, fmt.Errorf("failed to encrypt shared secret: %w", err)
	}

	return EncryptedSharedSecret{
		ChainId:         s.ChainId.String(),
		EncryptedSecret: encrSecret,
		InitialBlock:    s.InitialBlock,
	}, nil
}

func (e EncryptedSharedSecret) Decrypt(enc Encryptor) (SharedSecret, error) {
	secretBytes, err := enc.Decrypt(e.EncryptedSecret)
	if err != nil {
		return SharedSecret{}, fmt.Errorf("failed to decrypt shared secret: %w", err)
	}

	chainId := new(big.Int)
	chainId.SetString(e.ChainId, 10)

	return SharedSecret{
		ChainId:      chainId,
		Secret:       secretBytes,
		InitialBlock: e.InitialBlock,
	}, nil
}
