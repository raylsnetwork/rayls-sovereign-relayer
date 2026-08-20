package conv_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/raylsnetwork/rayls-sovereign-relayer/conv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringToBigInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    *big.Int
		wantErr bool
	}{
		{
			name:  "positive integer",
			input: "123",
			want:  big.NewInt(123),
		},
		{
			name:  "three digit number",
			input: "999",
			want:  big.NewInt(999),
		},
		{
			name:  "zero",
			input: "0",
			want:  big.NewInt(0),
		},
		{
			name:  "negative integer",
			input: "-42",
			want:  big.NewInt(-42),
		},
		{
			name:  "max uint64",
			input: "18446744073709551615",
			want:  new(big.Int).SetUint64(18446744073709551615),
		},
		{
			name:    "non-numeric string",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "float",
			input:   "1.5",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := conv.StringToBigInt(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, 0, got.Cmp(tt.want))
		})
	}
}

func TestStringToBytes32(t *testing.T) {
	t.Parallel()

	validHex := "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"
	validBytes, _ := hex.DecodeString(validHex)
	var expectedValid [32]byte
	copy(expectedValid[:], validBytes)

	tests := []struct {
		name    string
		input   string
		want    [32]byte
		wantErr bool
	}{
		{
			name:  "valid 32-byte hex",
			input: validHex,
			want:  expectedValid,
		},
		{
			name:  "all zeros",
			input: "0000000000000000000000000000000000000000000000000000000000000000",
			want:  [32]byte{},
		},
		{
			name:  "all ff",
			input: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			want: [32]byte{
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			},
		},
		{
			name:    "too short",
			input:   "abcdef",
			wantErr: true,
		},
		{
			name:    "too long",
			input:   "abcdef0123456789abcdef0123456789abcdef0123456789abcdef01234567890000",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid hex characters",
			input:   "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
			wantErr: true,
		},
		{
			name:    "63 characters (odd length)",
			input:   "abcdef0123456789abcdef0123456789abcdef0123456789abcdef012345678",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := conv.StringToBytes32(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
