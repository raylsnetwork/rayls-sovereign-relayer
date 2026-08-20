package domain_test

import (
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/domain"
	"github.com/stretchr/testify/require"
)

// stubEncryptor is a minimal Encryptor that returns the canned `out` from
// Decrypt and `errOut` from Encrypt; used to drive specific Decrypt outcomes
// without depending on the service-layer mocks.
type stubEncryptor struct {
	encryptFn func([]byte) ([]byte, error)
	decryptFn func([]byte) ([]byte, error)
}

func (s *stubEncryptor) Encrypt(b []byte) ([]byte, error) { return s.encryptFn(b) }
func (s *stubEncryptor) Decrypt(b []byte) ([]byte, error) { return s.decryptFn(b) }

func TestEncryptedRaylsSignKeyList_Decrypt(t *testing.T) {
	t.Run("returns parse error when decrypted bytes are not a valid ECDSA key", func(t *testing.T) {
		// Encryptor.Decrypt yields garbage bytes that crypto.ToECDSA rejects.
		// Pre-fix this would silently produce a nil *ecdsa.PrivateKey and panic
		// later at crypto.PubkeyToAddress; post-fix the error must surface here.
		enc := &stubEncryptor{
			decryptFn: func([]byte) ([]byte, error) { return []byte{0x00}, nil },
		}
		encrypted := domain.EncryptedRaylsSignKeyList{[]byte("ciphertext")}

		got, err := encrypted.Decrypt(enc)

		require.Error(t, err)
		require.Contains(t, err.Error(), "parsing ECDSA key from decrypted bytes")
		// Outer wrapper must say "decrypt", not "encrypt" (typo in earlier version).
		require.Contains(t, err.Error(), "failed to decrypt")
		require.NotContains(t, err.Error(), "failed to encrypt")
		require.Empty(t, got)
	})

	t.Run("propagates encryptor decrypt errors", func(t *testing.T) {
		wantErr := errors.New("decryption backend offline")
		enc := &stubEncryptor{
			decryptFn: func([]byte) ([]byte, error) { return nil, wantErr },
		}
		encrypted := domain.EncryptedRaylsSignKeyList{[]byte("ciphertext")}

		got, err := encrypted.Decrypt(enc)

		require.ErrorIs(t, err, wantErr)
		require.Empty(t, got)
	})

	t.Run("round-trips a valid key via Encrypt/Decrypt", func(t *testing.T) {
		// Identity encryptor lets the round-trip be a no-op so we exercise the
		// happy path of Encrypt → Decrypt purely through marshal/unmarshal.
		enc := &stubEncryptor{
			encryptFn: func(b []byte) ([]byte, error) { return b, nil },
			decryptFn: func(b []byte) ([]byte, error) { return b, nil },
		}
		key, err := crypto.GenerateKey()
		require.NoError(t, err)

		encrypted, err := domain.RaylsSignKeyList{key}.Encrypt(enc)
		require.NoError(t, err)

		decrypted, err := encrypted.Decrypt(enc)
		require.NoError(t, err)
		require.Len(t, decrypted, 1)
		require.Equal(t,
			crypto.PubkeyToAddress(key.PublicKey),
			crypto.PubkeyToAddress(decrypted[0].PublicKey),
		)
	})
}
