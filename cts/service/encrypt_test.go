package service_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/domain"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// passthroughEncryptor is a no-op Encryptor used for tests — Encrypt and
// Decrypt both return the input unchanged. Lets us exercise the
// EncryptService flow without pulling in a real KMS client.
type passthroughEncryptor struct{}

func (passthroughEncryptor) Encrypt(data []byte) ([]byte, error) { return data, nil }
func (passthroughEncryptor) Decrypt(data []byte) ([]byte, error) { return data, nil }

// stubEncryptKeysService returns a fixed shared secret for every call. The
// real service fetches from the repository layer; we're only interested in
// verifying that GCMEncrypt → GCMDecrypt is a round-trip, so a static
// secret is enough.
type stubEncryptKeysService struct {
	secret domain.SharedSecret
}

func (s *stubEncryptKeysService) GetSharedSecret(_ context.Context, _ *big.Int, _ uint64) (domain.SharedSecret, error) {
	return s.secret, nil
}

func (s *stubEncryptKeysService) GetAllSharedSecrets(_ context.Context, _ uint64) ([]domain.SharedSecret, error) {
	return []domain.SharedSecret{s.secret}, nil
}

func (s *stubEncryptKeysService) GetEnygmaSharedSecrets(_ context.Context, _ []*big.Int, _ *big.Int, _ uint64, _ []byte) ([]domain.SharedSecret, error) {
	return []domain.SharedSecret{s.secret}, nil
}

func (s *stubEncryptKeysService) GetEnygmaSharedSelfSecret(_ context.Context, _ *big.Int, _ uint64, _ []byte) (domain.SharedSecret, error) {
	return s.secret, nil
}

func (s *stubEncryptKeysService) GetRaylsViewKeyPair(_ context.Context, _ uint64) (domain.RaylsViewKeyPair, error) {
	return domain.RaylsViewKeyPair{}, nil
}

func TestGCMEncryptDecryptRoundTrip(t *testing.T) {
	t.Parallel()

	keysSvc := &stubEncryptKeysService{
		secret: domain.SharedSecret{
			ChainId:      big.NewInt(1),
			Secret:       []byte("test-shared-secret-32-bytes-long"),
			InitialBlock: 100,
		},
	}
	svc := service.NewEncryptService(big.NewInt(1), keysSvc, passthroughEncryptor{})

	plaintext := []byte("hello world")
	const chainID = uint64(1)
	const blockNumber = uint64(100)

	ctx := context.Background()
	ciphertext, err := svc.GCMEncrypt(ctx, plaintext, chainID, blockNumber)
	require.NoError(t, err)
	assert.NotEmpty(t, ciphertext)
	assert.NotEqual(t, plaintext, ciphertext, "ciphertext must differ from plaintext")

	decrypted, err := svc.GCMDecrypt(ctx, ciphertext, chainID, blockNumber)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted, "round-trip must preserve plaintext")
}

// TestGCMDecrypt_TamperedCiphertext_ReturnsErrAuthFailed verifies that the
// caller-addressed decrypt path treats AEAD-fail as suspicious (tampered),
// not as "not for me" — see Option A in cts/service/encrypt_errors.go.
func TestGCMDecrypt_TamperedCiphertext_ReturnsErrAuthFailed(t *testing.T) {
	t.Parallel()

	keysSvc := &stubEncryptKeysService{
		secret: domain.SharedSecret{
			ChainId:      big.NewInt(1),
			Secret:       []byte("test-shared-secret-32-bytes-long"),
			InitialBlock: 100,
		},
	}
	svc := service.NewEncryptService(big.NewInt(1), keysSvc, passthroughEncryptor{})

	ctx := context.Background()
	ciphertext, err := svc.GCMEncrypt(ctx, []byte("payload"), 1, 100)
	require.NoError(t, err)

	// Flip one byte deep in the AEAD body so the tag verification fails.
	tampered := make([]byte, len(ciphertext))
	copy(tampered, ciphertext)
	tampered[len(tampered)-1] ^= 0xff

	_, err = svc.GCMDecrypt(ctx, tampered, 1, 100)
	require.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrAuthFailed), "GCMDecrypt must return ErrAuthFailed on AEAD failure; got %v", err)
	assert.False(t, errors.Is(err, service.ErrNotForRecipient), "GCMDecrypt must NOT return ErrNotForRecipient — caller addressed a specific chain")
}

// TestGCMDecryptWithProvidedSS_WrongSalt_ReturnsErrNotForRecipient verifies
// that the salt-based DVP path treats AEAD-fail as "not for me" rather than
// tampered — per Option A, since on this path AEAD-fail is the ONLY signal
// that the salt didn't match.
func TestGCMDecryptWithProvidedSS_WrongSalt_ReturnsErrNotForRecipient(t *testing.T) {
	t.Parallel()

	keysSvc := &stubEncryptKeysService{
		secret: domain.SharedSecret{ChainId: big.NewInt(1), Secret: []byte("placeholder")},
	}
	svc := service.NewEncryptService(big.NewInt(1), keysSvc, passthroughEncryptor{})

	correctSalt := []byte("correct-salt-32-bytes-long-XXXXX")
	wrongSalt := []byte("wrong-salt-32-bytes-long-YYYYYYY")

	ciphertext, err := svc.GCMEncryptWithProvidedSS([]byte("payload"), correctSalt)
	require.NoError(t, err)

	_, err = svc.GCMDecryptWithProvidedSS(ciphertext, wrongSalt)
	require.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrNotForRecipient), "wrong-salt decrypt must surface as NotForRecipient; got %v", err)
	assert.False(t, errors.Is(err, service.ErrAuthFailed), "wrong-salt decrypt must NOT surface as AuthFailed — on the salt path that's the normal not-for-me signal")
}
