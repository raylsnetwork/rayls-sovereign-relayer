package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

func generateAESKey(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

// EncryptRequest encrypts the plaintext using AES-GCM with the given secret key.
func EncryptData(secret string, plaintext []byte) ([]byte, error) {
	key := generateAESKey(secret)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, withstack.Wrap(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, withstack.Wrap(err)
	}

	// Generate a nonce for AES-GCM (must be unique for each encryption).
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, withstack.Wrap(err)
	}

	// Encrypt and prepend nonce to the ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// DecryptRequest decrypts the AES-GCM encrypted ciphertext using the provided secret key.
func DecryptData(secret string, encodedCiphertext []byte) ([]byte, error) {
	key := generateAESKey(secret)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, withstack.Wrap(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, withstack.Wrap(err)
	}

	// Decode the base64-encoded ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(string(encodedCiphertext))
	if err != nil {
		return nil, withstack.Wrap(err)
	}

	// Separate nonce and actual ciphertext
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, withstack.Wrap(errors.New("ciphertext too short"))
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the ciphertext using the nonce
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, withstack.Wrap(err)
	}

	return plaintext, nil
}
