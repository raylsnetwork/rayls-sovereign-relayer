package cryptography_test

import (
	"math/big"
	"testing"

	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cryptography"
	"github.com/stretchr/testify/assert"
)

// addBabyJubjubPoints adds two BabyJubjub points using projective coordinates.
func addBabyJubjubPoints(p1, p2 *babyjub.Point) *babyjub.Point {
	return babyjub.NewPoint().Projective().Add(p1.Projective(), p2.Projective()).Affine()
}

func TestGetPK(t *testing.T) {
	t.Run("returns generator G for scalar 1", func(t *testing.T) {
		result := cryptography.GetPK(big.NewInt(1))
		assert.Equal(t, 0, result.X.Cmp(cryptography.G.X))
		assert.Equal(t, 0, result.Y.Cmp(cryptography.G.Y))
	})

	t.Run("returns identity point for scalar 0", func(t *testing.T) {
		result := cryptography.GetPK(big.NewInt(0))
		assert.Equal(t, 0, result.X.Cmp(big.NewInt(0)))
		assert.Equal(t, 0, result.Y.Cmp(big.NewInt(1)))
	})

	t.Run("returns different points for different scalars", func(t *testing.T) {
		p1 := cryptography.GetPK(big.NewInt(1))
		p2 := cryptography.GetPK(big.NewInt(2))
		assert.False(t, p1.X.Cmp(p2.X) == 0 && p1.Y.Cmp(p2.Y) == 0)
	})
}

func TestGetH(t *testing.T) {
	t.Run("returns H base point for scalar 1", func(t *testing.T) {
		result := cryptography.GetH(big.NewInt(1))
		assert.Equal(t, 0, result.X.Cmp(cryptography.H.X))
		assert.Equal(t, 0, result.Y.Cmp(cryptography.H.Y))
	})

	t.Run("returns different point than GetPK for same scalar", func(t *testing.T) {
		pk := cryptography.GetPK(big.NewInt(5))
		h := cryptography.GetH(big.NewInt(5))
		assert.False(t, pk.X.Cmp(h.X) == 0 && pk.Y.Cmp(h.Y) == 0)
	})
}

func TestPedersenCommitmentEnygma(t *testing.T) {
	t.Run("is deterministic for same value and randomness", func(t *testing.T) {
		c1 := cryptography.PedersenCommitmentEnygma(big.NewInt(10), big.NewInt(20))
		c2 := cryptography.PedersenCommitmentEnygma(big.NewInt(10), big.NewInt(20))
		assert.Equal(t, 0, c1.X.Cmp(c2.X))
		assert.Equal(t, 0, c1.Y.Cmp(c2.Y))
	})

	t.Run("changes when value changes", func(t *testing.T) {
		c1 := cryptography.PedersenCommitmentEnygma(big.NewInt(10), big.NewInt(20))
		c2 := cryptography.PedersenCommitmentEnygma(big.NewInt(11), big.NewInt(20))
		assert.False(t, c1.X.Cmp(c2.X) == 0 && c1.Y.Cmp(c2.Y) == 0)
	})

	t.Run("changes when randomness changes", func(t *testing.T) {
		c1 := cryptography.PedersenCommitmentEnygma(big.NewInt(10), big.NewInt(20))
		c2 := cryptography.PedersenCommitmentEnygma(big.NewInt(10), big.NewInt(21))
		assert.False(t, c1.X.Cmp(c2.X) == 0 && c1.Y.Cmp(c2.Y) == 0)
	})

	t.Run("zero value equals r times H", func(t *testing.T) {
		r := big.NewInt(42)
		c := cryptography.PedersenCommitmentEnygma(big.NewInt(0), r)
		h := cryptography.GetH(r)
		assert.Equal(t, 0, c.X.Cmp(h.X))
		assert.Equal(t, 0, c.Y.Cmp(h.Y))
	})

	t.Run("zero randomness equals v times G", func(t *testing.T) {
		v := big.NewInt(42)
		c := cryptography.PedersenCommitmentEnygma(v, big.NewInt(0))
		g := cryptography.GetPK(v)
		assert.Equal(t, 0, c.X.Cmp(g.X))
		assert.Equal(t, 0, c.Y.Cmp(g.Y))
	})

	t.Run("homomorphic addition property holds", func(t *testing.T) {
		v1 := big.NewInt(100)
		r1 := big.NewInt(200)
		v2 := big.NewInt(300)
		r2 := big.NewInt(400)

		c1 := cryptography.PedersenCommitmentEnygma(v1, r1)
		c2 := cryptography.PedersenCommitmentEnygma(v2, r2)
		cSum := addBabyJubjubPoints(c1, c2)

		vSum := cryptography.AddMod(v1, v2, cryptography.JubJubPrimeSubGroup)
		rSum := cryptography.AddMod(r1, r2, cryptography.JubJubPrimeSubGroup)
		cExpected := cryptography.PedersenCommitmentEnygma(vSum, rSum)

		assert.Equal(t, 0, cSum.X.Cmp(cExpected.X))
		assert.Equal(t, 0, cSum.Y.Cmp(cExpected.Y))
	})
}

func TestGetNegative(t *testing.T) {
	t.Run("returns zero for zero input", func(t *testing.T) {
		result := cryptography.GetNegative(big.NewInt(0))
		assert.Equal(t, 0, result.Cmp(big.NewInt(0)))
	})

	t.Run("returns subgroup prime minus x for positive x", func(t *testing.T) {
		x := big.NewInt(5)
		expected := new(big.Int).Sub(cryptography.JubJubPrimeSubGroup, x)
		result := cryptography.GetNegative(x)
		assert.Equal(t, 0, result.Cmp(expected))
	})

	t.Run("x plus GetNegative of x equals subgroup prime", func(t *testing.T) {
		x := big.NewInt(12345)
		neg := cryptography.GetNegative(x)
		sum := new(big.Int).Add(x, neg)
		assert.Equal(t, 0, sum.Cmp(cryptography.JubJubPrimeSubGroup))
	})
}

func TestAddMod(t *testing.T) {
	t.Run("adds two values modulo p", func(t *testing.T) {
		result := cryptography.AddMod(big.NewInt(3), big.NewInt(4), big.NewInt(10))
		assert.Equal(t, 0, result.Cmp(big.NewInt(7)))
	})

	t.Run("wraps around when sum exceeds modulus", func(t *testing.T) {
		result := cryptography.AddMod(big.NewInt(7), big.NewInt(8), big.NewInt(10))
		assert.Equal(t, 0, result.Cmp(big.NewInt(5)))
	})

	t.Run("identity with zero operand", func(t *testing.T) {
		result := cryptography.AddMod(big.NewInt(5), big.NewInt(0), big.NewInt(10))
		assert.Equal(t, 0, result.Cmp(big.NewInt(5)))
	})

	t.Run("works with large BabyJubjub values", func(t *testing.T) {
		pMinus1 := new(big.Int).Sub(cryptography.JubJubPrimeSubGroup, big.NewInt(1))
		result := cryptography.AddMod(pMinus1, big.NewInt(1), cryptography.JubJubPrimeSubGroup)
		assert.Equal(t, 0, result.Cmp(big.NewInt(0)))
	})
}

func TestSubMod(t *testing.T) {
	t.Run("subtracts two values modulo p", func(t *testing.T) {
		result := cryptography.SubMod(big.NewInt(7), big.NewInt(3), big.NewInt(10))
		assert.Equal(t, 0, result.Cmp(big.NewInt(4)))
	})

	t.Run("wraps to non-negative when a is less than b", func(t *testing.T) {
		result := cryptography.SubMod(big.NewInt(3), big.NewInt(7), big.NewInt(10))
		assert.Equal(t, 0, result.Cmp(big.NewInt(6)))
	})

	t.Run("returns zero when a equals b", func(t *testing.T) {
		result := cryptography.SubMod(big.NewInt(5), big.NewInt(5), big.NewInt(10))
		assert.Equal(t, 0, result.Cmp(big.NewInt(0)))
	})

	t.Run("SubMod and AddMod are inverse operations", func(t *testing.T) {
		a := big.NewInt(42)
		b := big.NewInt(17)
		p := big.NewInt(100)

		sum := cryptography.AddMod(a, b, p)
		result := cryptography.SubMod(sum, b, p)
		assert.Equal(t, 0, result.Cmp(a))
	})
}
