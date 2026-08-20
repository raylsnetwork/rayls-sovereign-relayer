package cryptography

import (
	"fmt"
	"math/big"

	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

// GetPoseidonHash returns the raw Poseidon hash without modular reduction
// Use this when the hash will be used as input to another Poseidon call (chained hashing)
// The circuit computes intermediate hashes without reduction, only reducing the final result
func GetPoseidonHash(inputs []*big.Int) (*big.Int, error) {
	return poseidon.Hash(inputs)
}

// GetPoseidonHashModNumber calculates the Poseidon hash of input(s) and returns the result modulo the specified number
func GetPoseidonHashModNumber(input interface{}, modulus *big.Int) (*big.Int, error) {
	var inputs []*big.Int

	switch v := input.(type) {
	case *big.Int:
		// Handle single *big.Int value
		inputs = []*big.Int{v}
	case []*big.Int:
		// Handle slice of *big.Int values
		inputs = v
	default:
		return nil, fmt.Errorf("unsupported input type: %T, expected *big.Int or []*big.Int", input)
	}

	// Calculate Poseidon hash
	poseidonHash, err := poseidon.Hash(inputs)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("computing poseidon hash: %w", err))
	}

	// Apply modulo operation
	if modulus != nil && modulus.Sign() > 0 {
		poseidonHash.Mod(poseidonHash, modulus)
	}

	return poseidonHash, nil
}
