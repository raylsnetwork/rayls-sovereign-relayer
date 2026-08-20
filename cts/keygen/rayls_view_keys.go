package keygen

import (
	"crypto/mlkem"
	"fmt"
	"math/big"

	"github.com/iden3/go-iden3-crypto/constants"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/domain"
)

func GenerateRaylsViewKeys() (domain.RaylsViewKeyPair, error) {
	dk, err := mlkem.GenerateKey768()
	if err != nil {
		return domain.RaylsViewKeyPair{}, fmt.Errorf("generating ML-KEM 768 key pair: %w", err)
	}

	ek := dk.EncapsulationKey()

	return domain.RaylsViewKeyPair{
		RaylsViewPrivateKey: dk,
		RaylsViewPublicKey:  ek,
	}, nil
}

// Uses encapsulation key to generate a new fresh shared secret and a ciphertext
func GenerateSharedSecret(encapsulationKey []byte) ([]byte, []byte, error) {
	ek, err := mlkem.NewEncapsulationKey768(encapsulationKey)
	if err != nil {
		return nil, nil, fmt.Errorf("creating ML-KEM 768 encapsulation key: %w", err)
	}

	sharedSecret, ciphertext := ek.Encapsulate()

	return ciphertext, sharedSecret, nil
}

// GenerateKeyDigest takes a shared secret and generates a key digest to be posted as an identifier
// receives an initial shared secret ([]byte)
// returns the Poseidon hash of the secret and an error
func GenerateKeyDigest(secret []byte) ([]byte, error) {
	// create a new big.Int
	z := new(big.Int)

	// set the value of z (big.Int) to the shared secret ([]byte)
	z.SetBytes(secret)

	// reduce it to fit in the field
	z.Mod(z, constants.Q)

	inputs := []*big.Int{z}

	digest, err := poseidon.Hash(inputs)
	if err != nil {
		return nil, fmt.Errorf("computing Poseidon hash for key digest: %w", err)
	}

	return digest.Bytes(), nil
}

// IMPORTANT NOTE:
// ML-KEM has non-zero probability of failure, meaning two honest parties may derive different shared secrets. This causes handshake failure.
// ML-KEM has a cryptographically small failure rate less than 2^-138; Clients should retry if a failure is encountered.

// Decapsulate takes Alice's DecapsulationKey and Bob's ciphertext
// receives a decapsulation key and ciphertext
// returns a (decapsulated) shared secret and an error
func RecoverSharedSecret(dk *mlkem.DecapsulationKey768, ciphertext []byte) ([]byte, error) {
	secret, err := dk.Decapsulate(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("decapsulating shared secret: %w", err)
	}

	return secret, nil
}
