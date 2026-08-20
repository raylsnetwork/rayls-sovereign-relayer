package keygen

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/cryptography"
)

const (
	// combinedKeyBitShift is the bit shift used when combining the random value, timestamp, and chainId
	// into a single 256-bit number.
	combinedKeyBitShift = 80

	// timestampBitShift is the bit shift applied to the timestamp before combining with the chainId.
	timestampBitShift = 32
)

func GenerateRandomModJubJubPrimeSubGroupWithChainId(chainId *big.Int) (*big.Int, error) {
	prime := cryptography.JubJubPrimeSubGroup
	// Get the current timestamp in milliseconds
	timestamp := big.NewInt(time.Now().UnixNano() / int64(time.Millisecond))

	// Create a buffer from the integer and timestamp
	buffer := make([]byte, 12)
	binary.BigEndian.PutUint32(
		buffer,
		uint32(chainId.Uint64()), //nolint:gosec // chainId fits in uint32 for buffer seeding
	)
	binary.BigEndian.PutUint64(
		buffer[4:],
		uint64(timestamp.Int64()), //nolint:gosec // timestamp in milliseconds is positive
	)

	// Use the buffer to seed a PRNG
	seed := sha256.Sum256(buffer)
	prng := hmac.New(sha256.New, seed[:])

	// Generate a 256-bit random number using the PRNG
	randomBytes := make([]byte, 32) // 32 bytes = 256 bits, needs this amount to generate from 0 to prime -1
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, fmt.Errorf("reading random bytes for key generation: %w", err)
	}
	prng.Write(randomBytes)
	randomBytes = prng.Sum(nil)[:32] // Limit to 256 bits

	// Convert random bytes to a BigInt
	randomBigInt := new(big.Int).SetBytes(randomBytes)

	// Combine the randomBigInt, timestamp, and chainId, ensuring 256 bits
	combined := new(big.Int).Lsh(randomBigInt, combinedKeyBitShift)
	combined.Or(combined, new(big.Int).Lsh(timestamp, timestampBitShift))
	combined.Add(combined, chainId)

	// Reduce combined modulo prime
	combined.Mod(combined, prime)

	// Return the result
	return combined, nil
}

func GenerateRandomModJubJubPrimeGroupWithChainId(chainId *big.Int) (*big.Int, error) {
	prime := cryptography.JubJubPrimeGroup
	// Get the current timestamp in milliseconds
	timestamp := big.NewInt(time.Now().UnixNano() / int64(time.Millisecond))

	// Create a buffer from the integer and timestamp
	buffer := make([]byte, 12)
	binary.BigEndian.PutUint32(
		buffer,
		uint32(chainId.Uint64()), //nolint:gosec // chainId fits in uint32 for buffer seeding
	)
	binary.BigEndian.PutUint64(
		buffer[4:],
		uint64(timestamp.Int64()), //nolint:gosec // timestamp in milliseconds is positive
	)

	// Use the buffer to seed a PRNG
	seed := sha256.Sum256(buffer)
	prng := hmac.New(sha256.New, seed[:])

	// Generate a 256-bit random number using the PRNG
	randomBytes := make([]byte, 32) // 32 bytes = 256 bits, needs this amount to generate from 0 to prime -1
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, fmt.Errorf("reading random bytes for key generation: %w", err)
	}
	prng.Write(randomBytes)
	randomBytes = prng.Sum(nil)[:32] // Limit to 256 bits

	// Convert random bytes to a BigInt
	randomBigInt := new(big.Int).SetBytes(randomBytes)

	// Combine the randomBigInt, timestamp, and chainId, ensuring 256 bits
	combined := new(big.Int).Lsh(randomBigInt, combinedKeyBitShift)
	combined.Or(combined, new(big.Int).Lsh(timestamp, timestampBitShift))
	combined.Add(combined, chainId)

	// Reduce combined modulo prime
	combined.Mod(combined, prime)

	// Return the result
	return combined, nil
}

// CCFindPrivateKeyByAddress

func GetPaymentSpendPublicKeyFromSpendSecretKey(sk *big.Int) (*big.Int, error) {
	return cryptography.GetPoseidonHashModNumber(
		[]*big.Int{sk, sk},
		cryptography.JubJubPrimeSubGroup,
	)
}
