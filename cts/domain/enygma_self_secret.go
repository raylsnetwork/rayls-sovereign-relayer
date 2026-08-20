package domain

import "fmt"

// EnygmaSelfSecret represents a self secret used in the Enygma protocol.
// Self secrets are scoped per (resource, initial block): the stored value is
// Poseidon(prevRFactor, paymentSpendKey), where prevRFactor is per-resource.
type EnygmaSelfSecret struct {
	Secret       []byte
	InitialBlock uint64
	ResourceID   []byte
}

// EncryptedEnygmaSelfSecret represents an encrypted self secret for storage.
type EncryptedEnygmaSelfSecret struct {
	EncryptedSecret []byte
	InitialBlock    uint64
	ResourceID      []byte
}

// Encrypt encrypts the self secret using the provided encryptor.
func (s EnygmaSelfSecret) Encrypt(enc Encryptor) (EncryptedEnygmaSelfSecret, error) {
	encrSecret, err := enc.Encrypt(s.Secret)
	if err != nil {
		return EncryptedEnygmaSelfSecret{}, fmt.Errorf("failed to encrypt self secret: %w", err)
	}

	return EncryptedEnygmaSelfSecret{
		EncryptedSecret: encrSecret,
		InitialBlock:    s.InitialBlock,
		ResourceID:      s.ResourceID,
	}, nil
}

// Decrypt decrypts the encrypted self secret using the provided encryptor.
func (e EncryptedEnygmaSelfSecret) Decrypt(enc Encryptor) (EnygmaSelfSecret, error) {
	secretBytes, err := enc.Decrypt(e.EncryptedSecret)
	if err != nil {
		return EnygmaSelfSecret{}, fmt.Errorf("failed to decrypt self secret: %w", err)
	}

	return EnygmaSelfSecret{
		Secret:       secretBytes,
		InitialBlock: e.InitialBlock,
		ResourceID:   e.ResourceID,
	}, nil
}
