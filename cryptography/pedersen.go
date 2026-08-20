package cryptography

import (
	"crypto/elliptic"
	"math/big"
)

// Point a point on the elliptic curve.
type Point struct {
	X, Y *big.Int
}

// PedersenCommitment computes the Pedersen commitment for a given value and blinding factor.
func PedersenCommitment(value, blindingFactor *big.Int) Point {
	curve := elliptic.P256() // You can use a different elliptic curve if needed

	// Generator point G
	Gx, Gy := curve.ScalarBaseMult([]byte("generator"))
	G := Point{Gx, Gy}

	// Compute the commitment C = G * value + H * blindingFactor
	commitment := G.ScalarMult(value).Add(G.ScalarMult(blindingFactor))

	return commitment
}

// Add adds two points on the elliptic curve.
func (p1 Point) Add(p2 Point) Point {
	curve := elliptic.P256()
	x, y := curve.Add(p1.X, p1.Y, p2.X, p2.Y)
	return Point{x, y}
}

// ScalarMult multiplies a point by a scalar.
func (p Point) ScalarMult(scalar *big.Int) Point {
	curve := elliptic.P256()
	x, y := curve.ScalarMult(p.X, p.Y, scalar.Bytes())
	return Point{x, y}
}
