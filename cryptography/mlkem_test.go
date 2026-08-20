package cryptography_test

import (
	"crypto/mlkem"
	"testing"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/cryptography"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSalt_RecoverSalt_RoundTrip(t *testing.T) {
	t.Parallel()

	dk, err := mlkem.GenerateKey768()
	require.NoError(t, err)

	viewPK := dk.EncapsulationKey().Bytes()
	viewSK := dk.Bytes()

	salt, ctxt, err := cryptography.GenerateSalt(viewPK)
	require.NoError(t, err)
	require.NotNil(t, salt)
	require.NotEmpty(t, ctxt)

	recovered, err := cryptography.RecoverSalt(viewSK, ctxt)
	require.NoError(t, err)

	assert.Equal(t, salt, recovered, "recovered salt must match original")
}

func TestGenerateSalt_DifferentCallsProduceDifferentSalts(t *testing.T) {
	t.Parallel()

	dk, err := mlkem.GenerateKey768()
	require.NoError(t, err)

	viewPK := dk.EncapsulationKey().Bytes()

	salt1, _, err := cryptography.GenerateSalt(viewPK)
	require.NoError(t, err)

	salt2, _, err := cryptography.GenerateSalt(viewPK)
	require.NoError(t, err)

	assert.NotEqual(t, salt1, salt2, "each encapsulation should produce a unique salt")
}

func TestGenerateSalt_InvalidPublicKey(t *testing.T) {
	t.Parallel()

	_, _, err := cryptography.GenerateSalt([]byte("invalid"))
	require.Error(t, err)
}

func TestRecoverSalt_InvalidSecretKey(t *testing.T) {
	t.Parallel()

	_, err := cryptography.RecoverSalt([]byte("invalid"), []byte("ciphertext"))
	require.Error(t, err)
}

func TestRecoverSalt_InvalidCiphertext(t *testing.T) {
	t.Parallel()

	dk, err := mlkem.GenerateKey768()
	require.NoError(t, err)

	// ML-KEM-768 decapsulation doesn't return an error for invalid ciphertexts
	// (it returns an implicit rejection value), but the ciphertext must be the right length.
	// A too-short ciphertext will cause an error.
	_, err = cryptography.RecoverSalt(dk.Bytes(), []byte("short"))
	require.Error(t, err)
}
