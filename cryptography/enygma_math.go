package cryptography

import (
	"math/big"

	"github.com/iden3/go-iden3-crypto/babyjub"
)

var (
	// Constant curve values
	gx, _ = new(
		big.Int,
	).SetString("16540640123574156134436876038791482806971768689494387082833631921987005038935", 10)
	gy, _ = new(
		big.Int,
	).SetString("20819045374670962167435360035096875258406992893633759881276124905556507972311", 10)
	G     = &babyjub.Point{X: gx, Y: gy}
	hx, _ = new(
		big.Int,
	).SetString("10100005861917718053548237064487763771145251762383025193119768015180892676690", 10)
	hy, _ = new(
		big.Int,
	).SetString("7512830269827713629724023825249861327768672768516116945507944076335453576011", 10)
	H                   = &babyjub.Point{X: hx, Y: hy}
	JubJubPrimeGroup, _ = new(
		big.Int,
	).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
	// group prime in babyjubjub is 21888242871839275222246405745257275088548364400416034343698204186575808495617, a 254 bit number
	JubJubPrimeSubGroup, _ = new(
		big.Int,
	).SetString("2736030358979909402780800718157159386076813972158567259200215660948447373041", 10)
	// sub group prime in babyjubjub is 2736030358979909402780800718157159386076813972158567259200215660948447373041, a 251 bit number
	outOfFieldOffset            int64 = 14
	OutOfJubJubPrimeFieldNumber       = new(
		big.Int,
	).Add(JubJubPrimeGroup, big.NewInt(outOfFieldOffset))
	// This value is impossible since it is out of the prime field, should be used in null/zero/empty structs
)

// Points X and Y belong to Group, r and v values belong to subGroup

func GetPK(v *big.Int) *babyjub.Point {
	rG := babyjub.NewPoint().Mul(v, G)
	return rG
}

func GetH(v *big.Int) *babyjub.Point {
	rG := babyjub.NewPoint().Mul(v, H)
	return rG
}

func PedersenCommitmentEnygma(v *big.Int, r *big.Int) *babyjub.Point {
	vG := GetPK(v)
	vH := GetH(r)
	return addPks(vG, vH)
}

func GetNegative(x *big.Int) *big.Int {
	// If x is 0, return 0
	if x.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0)
	}
	// Otherwise, calculate the inverse as p - x
	inverse := big.NewInt(0).Sub(JubJubPrimeSubGroup, x)
	return inverse
}

func addPks(pk1 *babyjub.Point, pk2 *babyjub.Point) *babyjub.Point {
	return babyjub.NewPoint().Projective().Add(pk1.Projective(), pk2.Projective()).Affine()
}

// AddMod performs modular addition: (a + b) % p
func AddMod(a, b, p *big.Int) *big.Int {
	result := new(big.Int).Add(a, b)
	result.Mod(result, p)
	return result
}

// SubMod performs modular subtraction: (a - b) % p
func SubMod(a, b, p *big.Int) *big.Int {
	result := new(big.Int).Sub(a, b)
	result.Mod(result, p)
	if result.Sign() < 0 { // Ensure the result is non-negative
		result.Add(result, p)
	}
	return result
}
