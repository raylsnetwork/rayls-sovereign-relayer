package cryptography

import (
	"crypto/mlkem"
	"fmt"
	"math/big"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

// GenerateSalt encapsulates a shared secret using the recipient's ML-KEM-768 view public key.
// Returns the salt (shared secret as big.Int, reduced mod JubJubPrimeGroup so it is a valid
// Poseidon/BN254 field element) and the encapsulated ciphertext (CTXT).
func GenerateSalt(recipientViewPK []byte) (salt *big.Int, ctxt []byte, err error) {
	ek, err := mlkem.NewEncapsulationKey768(recipientViewPK)
	if err != nil {
		return nil, nil, withstack.Wrap(fmt.Errorf("parsing ML-KEM encapsulation key: %w", err))
	}

	sharedSecret, ciphertext := ek.Encapsulate()

	salt = new(big.Int).SetBytes(sharedSecret)
	salt.Mod(salt, JubJubPrimeGroup)

	return salt, ciphertext, nil
}

// RecoverSalt decapsulates the ciphertext using the local ML-KEM-768 view secret key
// to recover the shared secret (salt), reduced mod JubJubPrimeGroup to match GenerateSalt.
func RecoverSalt(viewSK []byte, ctxt []byte) (*big.Int, error) {
	dk, err := mlkem.NewDecapsulationKey768(viewSK)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("parsing ML-KEM decapsulation key: %w", err))
	}

	sharedSecret, err := dk.Decapsulate(ctxt)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("ML-KEM decapsulation: %w", err))
	}

	salt := new(big.Int).SetBytes(sharedSecret)
	salt.Mod(salt, JubJubPrimeGroup)

	return salt, nil
}
