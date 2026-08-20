package keygen

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateRaylsSignKeys(count int) ([]*ecdsa.PrivateKey, error) {
	keys := []*ecdsa.PrivateKey{}
	for range count {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			return nil, fmt.Errorf("error generating private key: %w", err)
		}
		keys = append(keys, privateKey)
	}
	return keys, nil
}
