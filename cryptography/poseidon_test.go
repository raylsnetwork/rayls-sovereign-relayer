package cryptography_test

import (
	"math/big"
	"testing"

	"github.com/raylsnetwork/rayls-sovereign-relayer/cryptography"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPoseidonHashModNumber(t *testing.T) {
	t.Run("hashes a single big.Int input", func(t *testing.T) {
		result, err := cryptography.GetPoseidonHashModNumber(big.NewInt(42), nil)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Sign() >= 0)
	})

	t.Run("is deterministic for same input", func(t *testing.T) {
		h1, err := cryptography.GetPoseidonHashModNumber(big.NewInt(42), nil)
		require.NoError(t, err)

		h2, err := cryptography.GetPoseidonHashModNumber(big.NewInt(42), nil)
		require.NoError(t, err)

		assert.Equal(t, 0, h1.Cmp(h2))
	})

	t.Run("produces different hashes for different inputs", func(t *testing.T) {
		h1, err := cryptography.GetPoseidonHashModNumber(big.NewInt(1), nil)
		require.NoError(t, err)

		h2, err := cryptography.GetPoseidonHashModNumber(big.NewInt(2), nil)
		require.NoError(t, err)

		assert.NotEqual(t, 0, h1.Cmp(h2))
	})

	t.Run("hashes a slice of big.Int inputs", func(t *testing.T) {
		inputs := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}
		result, err := cryptography.GetPoseidonHashModNumber(inputs, nil)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("applies modulus when provided", func(t *testing.T) {
		modulus := big.NewInt(100)
		result, err := cryptography.GetPoseidonHashModNumber(big.NewInt(42), modulus)
		require.NoError(t, err)
		assert.True(t, result.Cmp(modulus) < 0, "result should be less than modulus")
		assert.True(t, result.Sign() >= 0, "result should be non-negative")
	})

	t.Run("skips modulus when nil", func(t *testing.T) {
		result, err := cryptography.GetPoseidonHashModNumber(big.NewInt(42), nil)
		require.NoError(t, err)
		// Poseidon hash output is a large number; without modulus it stays large
		assert.NotNil(t, result)
	})

	t.Run("skips modulus when zero", func(t *testing.T) {
		withoutMod, err := cryptography.GetPoseidonHashModNumber(big.NewInt(42), nil)
		require.NoError(t, err)

		withZeroMod, err := cryptography.GetPoseidonHashModNumber(big.NewInt(42), big.NewInt(0))
		require.NoError(t, err)

		assert.Equal(t, 0, withoutMod.Cmp(withZeroMod), "zero modulus should behave like nil modulus")
	})

	t.Run("returns error for unsupported input type string", func(t *testing.T) {
		_, err := cryptography.GetPoseidonHashModNumber("invalid", nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported input type")
	})

	t.Run("returns error for unsupported input type int", func(t *testing.T) {
		_, err := cryptography.GetPoseidonHashModNumber(42, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported input type")
	})
}
