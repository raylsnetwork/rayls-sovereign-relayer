package testdata

import (
	"crypto/mlkem"
	"encoding/hex"
	"sync"
)

// ML-KEM-768 keypair generated once per process and cached. The DvP
// production code calls cryptography.GenerateSalt(viewPK) and
// cryptography.RecoverSalt(viewSK, ctxt) — both require a real
// well-formed keypair. Each fixture call below returns the same shared
// keypair so a single call to GenerateSalt + RecoverSalt produces a
// matching round-trip across tests.
var (
	mlkemOnce       sync.Once
	mlkemEncPK      []byte
	mlkemDecSK      []byte
	mlkemEncPKHex   string
	mlkemDecSKHex   string
)

func ensureMlkemKeyPair() {
	mlkemOnce.Do(func() {
		dk, err := mlkem.GenerateKey768()
		if err != nil {
			panic("ML-KEM-768 keygen failed in test fixture: " + err.Error())
		}
		mlkemEncPK = dk.EncapsulationKey().Bytes()
		mlkemDecSK = dk.Bytes()
		mlkemEncPKHex = hex.EncodeToString(mlkemEncPK)
		mlkemDecSKHex = hex.EncodeToString(mlkemDecSK)
	})
}

// MlkemEncapsulationKey returns the cached ML-KEM-768 encapsulation key
// (raw bytes). Use as input to cryptography.GenerateSalt(viewPK).
func MlkemEncapsulationKey() []byte {
	ensureMlkemKeyPair()
	// Defensive copy so callers can't mutate the cached fixture.
	out := make([]byte, len(mlkemEncPK))
	copy(out, mlkemEncPK)
	return out
}

// MlkemDecapsulationKey returns the cached ML-KEM-768 decapsulation key
// (raw bytes). Use as input to cryptography.RecoverSalt(viewSK, ctxt).
func MlkemDecapsulationKey() []byte {
	ensureMlkemKeyPair()
	out := make([]byte, len(mlkemDecSK))
	copy(out, mlkemDecSK)
	return out
}

// MlkemKeyPair returns both halves of the cached ML-KEM-768 keypair.
func MlkemKeyPair() (encPK, decSK []byte) {
	return MlkemEncapsulationKey(), MlkemDecapsulationKey()
}

// MlkemEncapsulationKeyHex returns the hex-encoded encapsulation key —
// the form ParticipantStorage.RaylsViewPublicKey and
// RaylsViewKeyPair.RaylsViewPublicKey are stored as in the relayer flow.
func MlkemEncapsulationKeyHex() string {
	ensureMlkemKeyPair()
	return mlkemEncPKHex
}

// MlkemDecapsulationKeyHex returns the hex-encoded decapsulation key —
// the form RaylsViewKeyPair.RaylsViewSecretKey is stored as.
func MlkemDecapsulationKeyHex() string {
	ensureMlkemKeyPair()
	return mlkemDecSKHex
}

// MlkemKeyPairHex returns both halves hex-encoded.
func MlkemKeyPairHex() (encPKHex, decSKHex string) {
	return MlkemEncapsulationKeyHex(), MlkemDecapsulationKeyHex()
}
