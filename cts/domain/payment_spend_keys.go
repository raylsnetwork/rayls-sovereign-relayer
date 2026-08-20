package domain

import (
	"fmt"
	"math/big"
)

type PaymentSpendKeys struct {
	SecretKey *big.Int
	PublicKey *big.Int
}

func (b PaymentSpendKeys) Encrypt(enc Encryptor) (EncryptedPaymentSpendKeys, error) {
	secretKeyBytes := b.SecretKey.Bytes()

	encrSecretKey, err := enc.Encrypt(secretKeyBytes)
	if err != nil {
		return EncryptedPaymentSpendKeys{}, fmt.Errorf("failed to encrypt payment spend secret key: %w", err)
	}

	publicKeyBytes := b.PublicKey.Bytes()

	return EncryptedPaymentSpendKeys{
		EncryptedSecretKey: encrSecretKey,
		PublicKey:          publicKeyBytes,
	}, nil
}

type EncryptedPaymentSpendKeys struct {
	EncryptedSecretKey []byte
	PublicKey          []byte
}

func (e EncryptedPaymentSpendKeys) Decrypt(enc Encryptor) (PaymentSpendKeys, error) {
	secretKeyBytes, err := enc.Decrypt(e.EncryptedSecretKey)
	if err != nil {
		return PaymentSpendKeys{}, fmt.Errorf("failed to decrypt payment spend secret key: %w", err)
	}

	secretKey := new(big.Int).SetBytes(secretKeyBytes)
	publicKey := new(big.Int).SetBytes(e.PublicKey)

	return PaymentSpendKeys{
		SecretKey: secretKey,
		PublicKey: publicKey,
	}, nil
}
