package cryptography_test

import (
	"crypto/elliptic"
	"math/big"
	"testing"

	"github.com/raylsnetwork/rayls-sovereign-relayer/cryptography"
	"github.com/stretchr/testify/assert"
)

func TestPedersenCommitment(t *testing.T) {
	t.Run("returns a point on the P-256 curve", func(t *testing.T) {
		c := cryptography.PedersenCommitment(big.NewInt(42), big.NewInt(7))
		assert.True(t, elliptic.P256().IsOnCurve(c.X, c.Y))
	})

	t.Run("is deterministic for same inputs", func(t *testing.T) {
		c1 := cryptography.PedersenCommitment(big.NewInt(42), big.NewInt(7))
		c2 := cryptography.PedersenCommitment(big.NewInt(42), big.NewInt(7))
		assert.Equal(t, 0, c1.X.Cmp(c2.X))
		assert.Equal(t, 0, c1.Y.Cmp(c2.Y))
	})

	t.Run("changes when value changes", func(t *testing.T) {
		c1 := cryptography.PedersenCommitment(big.NewInt(42), big.NewInt(7))
		c2 := cryptography.PedersenCommitment(big.NewInt(43), big.NewInt(7))
		assert.False(t, c1.X.Cmp(c2.X) == 0 && c1.Y.Cmp(c2.Y) == 0)
	})

	t.Run("changes when blinding factor changes", func(t *testing.T) {
		c1 := cryptography.PedersenCommitment(big.NewInt(42), big.NewInt(7))
		c2 := cryptography.PedersenCommitment(big.NewInt(42), big.NewInt(8))
		assert.False(t, c1.X.Cmp(c2.X) == 0 && c1.Y.Cmp(c2.Y) == 0)
	})
}

func TestPointAdd(t *testing.T) {
	t.Run("result is on the P-256 curve", func(t *testing.T) {
		p1 := cryptography.PedersenCommitment(big.NewInt(10), big.NewInt(20))
		p2 := cryptography.PedersenCommitment(big.NewInt(30), big.NewInt(40))
		sum := p1.Add(p2)
		assert.True(t, elliptic.P256().IsOnCurve(sum.X, sum.Y))
	})

	t.Run("addition is commutative", func(t *testing.T) {
		p1 := cryptography.PedersenCommitment(big.NewInt(10), big.NewInt(20))
		p2 := cryptography.PedersenCommitment(big.NewInt(30), big.NewInt(40))
		sum1 := p1.Add(p2)
		sum2 := p2.Add(p1)
		assert.Equal(t, 0, sum1.X.Cmp(sum2.X))
		assert.Equal(t, 0, sum1.Y.Cmp(sum2.Y))
	})
}

func TestPointScalarMult(t *testing.T) {
	t.Run("scalar multiplication by 1 returns the same point", func(t *testing.T) {
		p := cryptography.PedersenCommitment(big.NewInt(42), big.NewInt(7))
		result := p.ScalarMult(big.NewInt(1))
		assert.Equal(t, 0, result.X.Cmp(p.X))
		assert.Equal(t, 0, result.Y.Cmp(p.Y))
	})

	t.Run("scalar multiplication by 2 equals adding point to itself", func(t *testing.T) {
		p := cryptography.PedersenCommitment(big.NewInt(42), big.NewInt(7))
		doubled := p.ScalarMult(big.NewInt(2))
		added := p.Add(p)
		assert.Equal(t, 0, doubled.X.Cmp(added.X))
		assert.Equal(t, 0, doubled.Y.Cmp(added.Y))
	})

	t.Run("result is on the P-256 curve", func(t *testing.T) {
		p := cryptography.PedersenCommitment(big.NewInt(42), big.NewInt(7))
		result := p.ScalarMult(big.NewInt(5))
		assert.True(t, elliptic.P256().IsOnCurve(result.X, result.Y))
	})
}
