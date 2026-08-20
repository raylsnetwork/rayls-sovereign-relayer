package keygen_test

import (
	"testing"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/keygen"
	"github.com/stretchr/testify/require"
)

func TestGenerateRaylsViewKeys(t *testing.T) {
	t.Run("returns valid ML-KEM key pair", func(t *testing.T) {
		pair, err := keygen.GenerateRaylsViewKeys()

		require.Nil(t, err)
		require.NotNil(t, pair.RaylsViewPrivateKey)
		require.NotNil(t, pair.RaylsViewPublicKey)
	})

	t.Run("keys can be serialized to bytes", func(t *testing.T) {
		pair, err := keygen.GenerateRaylsViewKeys()
		require.Nil(t, err)

		privateKeyBytes := pair.RaylsViewPrivateKey.Bytes()
		publicKeyBytes := pair.RaylsViewPublicKey.Bytes()

		require.NotEmpty(t, privateKeyBytes)
		require.NotEmpty(t, publicKeyBytes)
	})
}

func TestGenerateSharedSecret(t *testing.T) {
	t.Run("returns non-empty ciphertext and shared secret", func(t *testing.T) {
		pair, err := keygen.GenerateRaylsViewKeys()
		require.Nil(t, err)

		publicKeyBytes := pair.RaylsViewPublicKey.Bytes()
		ciphertext, sharedSecret, err := keygen.GenerateSharedSecret(publicKeyBytes)

		require.Nil(t, err)
		require.NotEmpty(t, ciphertext)
		require.NotEmpty(t, sharedSecret)
	})

	t.Run("returns error on invalid encapsulation key", func(t *testing.T) {
		invalidKey := []byte("invalid key bytes")

		_, _, err := keygen.GenerateSharedSecret(invalidKey)

		require.NotNil(t, err)
	})
}

func TestGenerateKeyDigest(t *testing.T) {
	t.Run("returns non-empty digest", func(t *testing.T) {
		secret := []byte("test shared secret for hashing")

		digest, err := keygen.GenerateKeyDigest(secret)

		require.Nil(t, err)
		require.NotEmpty(t, digest)
	})

	t.Run("digest is deterministic for same input", func(t *testing.T) {
		secret := []byte("deterministic test secret")

		digest1, err1 := keygen.GenerateKeyDigest(secret)
		digest2, err2 := keygen.GenerateKeyDigest(secret)

		require.Nil(t, err1)
		require.Nil(t, err2)
		require.Equal(t, digest1, digest2)
	})

	t.Run("different inputs produce different digests", func(t *testing.T) {
		secret1 := []byte("first secret")
		secret2 := []byte("second secret")

		digest1, err1 := keygen.GenerateKeyDigest(secret1)
		digest2, err2 := keygen.GenerateKeyDigest(secret2)

		require.Nil(t, err1)
		require.Nil(t, err2)
		require.NotEqual(t, digest1, digest2)
	})
}

func TestRecoverSharedSecret(t *testing.T) {
	t.Run("recovers the same shared secret from ciphertext", func(t *testing.T) {
		pair, err := keygen.GenerateRaylsViewKeys()
		require.Nil(t, err)

		publicKeyBytes := pair.RaylsViewPublicKey.Bytes()
		ciphertext, originalSecret, err := keygen.GenerateSharedSecret(publicKeyBytes)
		require.Nil(t, err)

		recoveredSecret, err := keygen.RecoverSharedSecret(pair.RaylsViewPrivateKey, ciphertext)

		require.Nil(t, err)
		require.Equal(t, originalSecret, recoveredSecret)
	})

	t.Run("returns error on invalid ciphertext", func(t *testing.T) {
		pair, err := keygen.GenerateRaylsViewKeys()
		require.Nil(t, err)

		invalidCiphertext := []byte("invalid ciphertext")

		_, err = keygen.RecoverSharedSecret(pair.RaylsViewPrivateKey, invalidCiphertext)

		require.NotNil(t, err)
	})
}

func TestEncapsulateDecapsulateRoundtrip(t *testing.T) {
	t.Run("full key agreement roundtrip works correctly", func(t *testing.T) {
		// Party A generates their key pair
		partyA, err := keygen.GenerateRaylsViewKeys()
		require.Nil(t, err)

		// Party B generates their key pair
		partyB, err := keygen.GenerateRaylsViewKeys()
		require.Nil(t, err)

		// Party B encapsulates using Party A's public key
		ciphertext, partyBSecret, err := keygen.GenerateSharedSecret(partyA.RaylsViewPublicKey.Bytes())
		require.Nil(t, err)

		// Party A decapsulates using their own private key
		partyASecret, err := keygen.RecoverSharedSecret(partyA.RaylsViewPrivateKey, ciphertext)
		require.Nil(t, err)

		// Party A generate digest using shared secret
		partyADigest, err := keygen.GenerateKeyDigest(partyASecret)
		require.Nil(t, err)

		// Party B generate digest using shared secret
		partyBDigest, err := keygen.GenerateKeyDigest(partyBSecret)
		require.Nil(t, err)

		// Both parties should have the same shared secret and digest
		require.Equal(t, partyBSecret, partyASecret)
		require.Equal(t, partyADigest, partyBDigest)

		// Verify Party B's keys are different from Party A's (they are independent parties)
		require.NotEqual(t, partyA.RaylsViewPublicKey.Bytes(), partyB.RaylsViewPublicKey.Bytes())
	})

	t.Run("different key pairs produce different shared secrets", func(t *testing.T) {
		pair1, err := keygen.GenerateRaylsViewKeys()
		require.Nil(t, err)

		pair2, err := keygen.GenerateRaylsViewKeys()
		require.Nil(t, err)

		_, secret1, err := keygen.GenerateSharedSecret(pair1.RaylsViewPublicKey.Bytes())
		require.Nil(t, err)

		_, secret2, err := keygen.GenerateSharedSecret(pair2.RaylsViewPublicKey.Bytes())
		require.Nil(t, err)

		require.NotEqual(t, secret1, secret2)
	})
}
