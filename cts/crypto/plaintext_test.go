package crypto_test

import (
	"testing"

	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/crypto"
	"github.com/stretchr/testify/require"
)

func TestPlaintextEncryptor(t *testing.T) {
	t.Run("returns the same text on encrypt", func(t *testing.T) {
		wantText := "sample plain text"
		ptEncryptor := crypto.PlaintextEncryptor{}

		gotTextBytes, err := ptEncryptor.Encrypt([]byte(wantText))

		require.Nil(t, err, "expected no error but got one")
		require.Equal(t, wantText, string(gotTextBytes), "didn't get the same text")
	})

	t.Run("returns the same text on decrypt", func(t *testing.T) {
		wantText := "sample plain text"
		ptEncryptor := crypto.PlaintextEncryptor{}

		gotTextBytes, err := ptEncryptor.Decrypt([]byte(wantText))

		require.Nil(t, err, "expected no error but got one")
		require.Equal(t, wantText, string(gotTextBytes), "didn't get the same text")
	})
}
