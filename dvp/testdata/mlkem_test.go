package testdata

import (
	"crypto/mlkem"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMlkemFixture_Cached(t *testing.T) {
	t.Run("subsequent calls return the same key bytes (cached at init)", func(t *testing.T) {
		a := MlkemEncapsulationKey()
		b := MlkemEncapsulationKey()
		assert.Equal(t, a, b, "encapsulation key must be stable across calls so encrypt/decrypt round-trips work")

		c := MlkemDecapsulationKey()
		d := MlkemDecapsulationKey()
		assert.Equal(t, c, d)
	})

	t.Run("returned slice is a defensive copy — mutating the caller's copy does not affect the cache", func(t *testing.T) {
		a := MlkemEncapsulationKey()
		original := make([]byte, len(a))
		copy(original, a)
		a[0] ^= 0xFF // mutate caller's copy
		b := MlkemEncapsulationKey()
		assert.Equal(t, original, b, "fixture must hand out independent slices")
	})
}

func TestMlkemFixture_RoundTrip(t *testing.T) {
	t.Run("the cached keypair supports a full encapsulate -> decapsulate round-trip", func(t *testing.T) {
		encPK := MlkemEncapsulationKey()
		decSK := MlkemDecapsulationKey()

		ek, err := mlkem.NewEncapsulationKey768(encPK)
		require.NoError(t, err)
		shared, ctxt := ek.Encapsulate()

		dk, err := mlkem.NewDecapsulationKey768(decSK)
		require.NoError(t, err)
		recovered, err := dk.Decapsulate(ctxt)
		require.NoError(t, err)
		assert.Equal(t, shared, recovered)
	})
}

func TestMlkemFixture_Hex(t *testing.T) {
	t.Run("hex-encoded forms decode back to the raw bytes", func(t *testing.T) {
		encPKHex, decSKHex := MlkemKeyPairHex()
		decodedEnc, err := hex.DecodeString(encPKHex)
		require.NoError(t, err)
		assert.Equal(t, MlkemEncapsulationKey(), decodedEnc)

		decodedDec, err := hex.DecodeString(decSKHex)
		require.NoError(t, err)
		assert.Equal(t, MlkemDecapsulationKey(), decodedDec)
	})
}
