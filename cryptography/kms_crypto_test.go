package cryptography_test

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/raylsnetwork/rayls-sovereign-relayer/cryptography"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptData(t *testing.T) {
	t.Run("encrypts plaintext and returns base64-encoded ciphertext", func(t *testing.T) {
		ct, err := cryptography.EncryptData("my-secret", []byte("hello world"))
		require.NoError(t, err)
		assert.NotEmpty(t, ct)
		assert.NotEqual(t, []byte("hello world"), ct)

		_, decodeErr := base64.StdEncoding.DecodeString(string(ct))
		assert.NoError(t, decodeErr, "output should be valid base64")
	})

	t.Run("returns different ciphertexts for same input due to random nonce", func(t *testing.T) {
		ct1, err := cryptography.EncryptData("secret", []byte("same-data"))
		require.NoError(t, err)

		ct2, err := cryptography.EncryptData("secret", []byte("same-data"))
		require.NoError(t, err)

		assert.NotEqual(t, ct1, ct2)
	})

	t.Run("encrypts empty plaintext without error", func(t *testing.T) {
		ct, err := cryptography.EncryptData("secret", []byte{})
		require.NoError(t, err)
		assert.NotEmpty(t, ct, "even empty plaintext produces ciphertext (nonce + auth tag)")
	})
}

func TestDecryptData(t *testing.T) {
	t.Run("roundtrip encrypt then decrypt recovers original plaintext", func(t *testing.T) {
		original := []byte("test message")
		ct, err := cryptography.EncryptData("my-secret", original)
		require.NoError(t, err)

		plaintext, err := cryptography.DecryptData("my-secret", ct)
		require.NoError(t, err)
		assert.Equal(t, original, plaintext)
	})

	t.Run("roundtrip with empty plaintext", func(t *testing.T) {
		ct, err := cryptography.EncryptData("secret", []byte{})
		require.NoError(t, err)

		plaintext, err := cryptography.DecryptData("secret", ct)
		require.NoError(t, err)
		assert.Empty(t, plaintext)
	})

	t.Run("roundtrip with large plaintext preserves all bytes", func(t *testing.T) {
		original := bytes.Repeat([]byte("A"), 10240)
		ct, err := cryptography.EncryptData("secret", original)
		require.NoError(t, err)

		plaintext, err := cryptography.DecryptData("secret", ct)
		require.NoError(t, err)
		assert.Equal(t, original, plaintext)
	})

	t.Run("fails with wrong secret", func(t *testing.T) {
		ct, err := cryptography.EncryptData("secret-A", []byte("data"))
		require.NoError(t, err)

		_, err = cryptography.DecryptData("secret-B", ct)
		require.Error(t, err)
	})

	t.Run("fails with invalid base64 input", func(t *testing.T) {
		_, err := cryptography.DecryptData("secret", []byte("not-valid-base64!@#$"))
		require.Error(t, err)
	})

	t.Run("fails with ciphertext shorter than nonce size", func(t *testing.T) {
		shortData := base64.StdEncoding.EncodeToString([]byte{0x01})
		_, err := cryptography.DecryptData("secret", []byte(shortData))
		require.Error(t, err)
	})

	t.Run("fails with tampered ciphertext", func(t *testing.T) {
		ct, err := cryptography.EncryptData("secret", []byte("sensitive data"))
		require.NoError(t, err)

		raw, err := base64.StdEncoding.DecodeString(string(ct))
		require.NoError(t, err)

		// Flip a byte in the ciphertext portion (after the nonce)
		if len(raw) > 13 {
			raw[13] ^= 0xFF
		}
		tampered := []byte(base64.StdEncoding.EncodeToString(raw))

		_, err = cryptography.DecryptData("secret", tampered)
		require.Error(t, err)
	})
}
