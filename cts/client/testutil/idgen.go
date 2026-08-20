package testutil

import (
	"crypto/rand"
	"math/big"
	"strings"
)

func GenerateRandomIdentifier() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	nameLength := 12
	dashInterval := 4

	// Create a slice to hold the random characters
	randomBytes := make([]byte, nameLength)

	// Fill the slice with random characters from the charset
	for i := range randomBytes {
		b, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic("Failed to generate random name")
		}
		randomBytes[i] = charset[b.Int64()]
	}

	// Insert dashes every dashInterval characters
	var builder strings.Builder
	for i, b := range randomBytes {
		if i > 0 && i%dashInterval == 0 {
			builder.WriteByte('-')
		}
		builder.WriteByte(b)
	}

	return builder.String()
}
