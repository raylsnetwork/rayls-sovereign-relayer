package cryptography_test

import (
	"testing"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/cryptography"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// generate32ByteKey returns a deterministic 32-byte key for test use.
func generate32ByteKey() []byte {
	return cryptography.HashIt([]byte("test-key-seed"))
}

func TestGenerateSecureRandomValue(t *testing.T) {
	t.Run("returns slice of requested size", func(t *testing.T) {
		val, err := cryptography.GenerateSecureRandomValue(32)
		require.NoError(t, err)
		assert.Len(t, val, 32)
	})

	t.Run("returns different values on successive calls", func(t *testing.T) {
		a, err := cryptography.GenerateSecureRandomValue(16)
		require.NoError(t, err)

		b, err := cryptography.GenerateSecureRandomValue(16)
		require.NoError(t, err)

		assert.NotEqual(t, a, b)
	})

	t.Run("returns empty slice for size zero", func(t *testing.T) {
		val, err := cryptography.GenerateSecureRandomValue(0)
		require.NoError(t, err)
		assert.Len(t, val, 0)
	})
}

func TestHashIt(t *testing.T) {
	t.Run("returns 32 bytes for SHA3-256", func(t *testing.T) {
		result := cryptography.HashIt([]byte("data"))
		assert.Len(t, result, 32)
	})

	t.Run("is deterministic for the same input", func(t *testing.T) {
		h1 := cryptography.HashIt([]byte("same"))
		h2 := cryptography.HashIt([]byte("same"))
		assert.Equal(t, h1, h2)
	})

	t.Run("produces different hashes for different inputs", func(t *testing.T) {
		h1 := cryptography.HashIt([]byte("a"))
		h2 := cryptography.HashIt([]byte("b"))
		assert.NotEqual(t, h1, h2)
	})

	t.Run("hashes empty input without error", func(t *testing.T) {
		result := cryptography.HashIt([]byte{})
		assert.Len(t, result, 32)
	})
}

func TestGetSharedMessageTag(t *testing.T) {
	t.Run("returns 32-byte tag", func(t *testing.T) {
		result := cryptography.GetSharedMessageTag([]byte("secret"), []byte("random"))
		assert.Len(t, result, 32)
	})

	t.Run("is deterministic for same inputs", func(t *testing.T) {
		f1 := cryptography.GetSharedMessageTag([]byte("secret"), []byte("random"))
		f2 := cryptography.GetSharedMessageTag([]byte("secret"), []byte("random"))
		assert.Equal(t, f1, f2)
	})

	t.Run("changes when secret changes", func(t *testing.T) {
		f1 := cryptography.GetSharedMessageTag([]byte("secret-A"), []byte("random"))
		f2 := cryptography.GetSharedMessageTag([]byte("secret-B"), []byte("random"))
		assert.NotEqual(t, f1, f2)
	})

	t.Run("changes when randomness changes", func(t *testing.T) {
		f1 := cryptography.GetSharedMessageTag([]byte("secret"), []byte("random-A"))
		f2 := cryptography.GetSharedMessageTag([]byte("secret"), []byte("random-B"))
		assert.NotEqual(t, f1, f2)
	})
}

func TestEncryptGCM(t *testing.T) {
	t.Run("produces ciphertext with correct structure", func(t *testing.T) {
		key := generate32ByteKey()
		plaintext := []byte("hello")
		ad := make([]byte, 16)

		result, err := cryptography.EncryptGCM(key, plaintext, ad)
		require.NoError(t, err)

		// Structure: AD(16) + Nonce(12) + Ciphertext(len(plaintext)) + GCM tag(16)
		expectedMinLen := 16 + 12 + len(plaintext) + 16
		assert.GreaterOrEqual(t, len(result), expectedMinLen)
	})

	t.Run("succeeds with key shorter than 32 bytes", func(t *testing.T) {
		shortKey := []byte("16-byte-key!!!!!")
		plaintext := []byte("test data")
		ad := make([]byte, 16)

		_, err := cryptography.EncryptGCM(shortKey, plaintext, ad)
		require.NoError(t, err)
	})

	t.Run("produces different ciphertexts for same input due to random nonce", func(t *testing.T) {
		key := generate32ByteKey()
		plaintext := []byte("same-data")
		ad := make([]byte, 16)

		ct1, err := cryptography.EncryptGCM(key, plaintext, ad)
		require.NoError(t, err)
		ct2, err := cryptography.EncryptGCM(key, plaintext, ad)
		require.NoError(t, err)
		assert.NotEqual(t, ct1, ct2)
	})

	t.Run("preserves associated data in output prefix", func(t *testing.T) {
		key := generate32ByteKey()
		plaintext := []byte("test")
		ad := []byte("1234567890123456") // 16 bytes

		result, err := cryptography.EncryptGCM(key, plaintext, ad)
		require.NoError(t, err)
		assert.Equal(t, ad, result[:16])
	})

	t.Run("roundtrip with DecryptGCM recovers plaintext and associated data", func(t *testing.T) {
		key := generate32ByteKey()
		plaintext := []byte("hello world")
		ad := []byte("associated-data!")

		ct, err := cryptography.EncryptGCM(key, plaintext, ad)
		require.NoError(t, err)

		recoveredAD, recoveredPlaintext, err := cryptography.DecryptGCM(ct, key)
		require.NoError(t, err)
		assert.Equal(t, plaintext, recoveredPlaintext)
		assert.Equal(t, ad, recoveredAD)
	})
}

func TestDecryptGCM(t *testing.T) {
	t.Run("fails with wrong key", func(t *testing.T) {
		keyA := generate32ByteKey()
		keyB := cryptography.HashIt([]byte("different-seed"))
		plaintext := []byte("secret")
		ad := make([]byte, 16)

		ct, err := cryptography.EncryptGCM(keyA, plaintext, ad)
		require.NoError(t, err)

		_, _, err = cryptography.DecryptGCM(ct, keyB)
		require.Error(t, err)
	})

	t.Run("fails with tampered ciphertext", func(t *testing.T) {
		key := generate32ByteKey()
		plaintext := []byte("secret data")
		ad := make([]byte, 16)

		ct, err := cryptography.EncryptGCM(key, plaintext, ad)
		require.NoError(t, err)

		// Tamper with the encrypted portion (after AD + nonce = 28 bytes)
		if len(ct) > 29 {
			ct[29] ^= 0xFF
		}

		_, _, err = cryptography.DecryptGCM(ct, key)
		require.Error(t, err)
	})

	t.Run("panics with ciphertext shorter than expected header", func(t *testing.T) {
		key := generate32ByteKey()
		shortCT := make([]byte, 10)

		assert.Panics(t, func() {
			_, _, _ = cryptography.DecryptGCM(shortCT, key)
		}, "DecryptGCM does not bounds-check before slicing AD and nonce")
	})
}
