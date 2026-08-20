package conv

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

const (
	hexStringLen    = 64
	decodedBytesLen = 32
)

func StringToBigInt(str string) (*big.Int, error) {
	result, success := new(big.Int).SetString(str, 10)
	if !success {
		return nil, fmt.Errorf("parse string to big.Int %q: failed to parse", str)
	}
	return result, nil
}

func BigIntsToUint64s(vals []*big.Int) []uint64 {
	out := make([]uint64, len(vals))
	for i, v := range vals {
		out[i] = v.Uint64()
	}
	return out
}

func StringToBytes32(s string) ([32]byte, error) {
	if len(s) != hexStringLen {
		return [32]byte{}, errors.New("string must be exactly 64 characters long")
	}
	var b [32]byte
	decoded, err := hex.DecodeString(s)
	if err != nil {
		return b, fmt.Errorf("decoding hex string: %w", err)
	}
	if len(decoded) != decodedBytesLen {
		return b, errors.New("decoded byte slice must be exactly 32 bytes long")
	}
	copy(b[:], decoded)
	return b, nil
}
