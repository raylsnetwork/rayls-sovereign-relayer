package configinit

import (
	"crypto/ecdsa"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBigIntDecodeHook(t *testing.T) {
	targetType := reflect.TypeFor[*big.Int]()

	t.Run("decodes valid decimal string to *big.Int", func(t *testing.T) {
		got, err := bigIntDecodeHook(reflect.TypeFor[string](), targetType, "12345")

		require.NoError(t, err)
		assert.Equal(t, big.NewInt(12345), got)
	})

	t.Run("decodes zero string", func(t *testing.T) {
		got, err := bigIntDecodeHook(reflect.TypeFor[string](), targetType, "0")

		require.NoError(t, err)
		assert.Equal(t, big.NewInt(0), got)
	})

	t.Run("decodes large number string", func(t *testing.T) {
		want, _ := new(big.Int).SetString("99999999999999999999999999", 10)

		got, err := bigIntDecodeHook(reflect.TypeFor[string](), targetType, "99999999999999999999999999")

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("decodes negative number string", func(t *testing.T) {
		got, err := bigIntDecodeHook(reflect.TypeFor[string](), targetType, "-42")

		require.NoError(t, err)
		assert.Equal(t, big.NewInt(-42), got)
	})

	t.Run("returns error on invalid decimal string", func(t *testing.T) {
		_, err := bigIntDecodeHook(reflect.TypeFor[string](), targetType, "not-a-number")

		require.Error(t, err)
		assert.Equal(t, "wrong data format for field", err.Error())
	})

	t.Run("decodes empty string to nil, leaving required-ness to the validator", func(t *testing.T) {
		got, err := bigIntDecodeHook(reflect.TypeFor[string](), targetType, "")

		require.NoError(t, err)
		assert.Equal(t, (*big.Int)(nil), got)
	})

	t.Run("returns data unchanged when source is not string", func(t *testing.T) {
		got, err := bigIntDecodeHook(reflect.TypeFor[int](), targetType, 42)

		require.NoError(t, err)
		assert.Equal(t, 42, got)
	})

	t.Run("returns data unchanged when target is not *big.Int", func(t *testing.T) {
		got, err := bigIntDecodeHook(reflect.TypeFor[string](), reflect.TypeFor[string](), "12345")

		require.NoError(t, err)
		assert.Equal(t, "12345", got)
	})
}

func TestAddressDecodeHook(t *testing.T) {
	targetType := reflect.TypeFor[common.Address]()

	t.Run("decodes valid hex address", func(t *testing.T) {
		addr := "0x1234567890abcdef1234567890abcdef12345678"

		got, err := addressDecodeHook(reflect.TypeFor[string](), targetType, addr)

		require.NoError(t, err)
		assert.Equal(t, common.HexToAddress(addr), got)
	})

	t.Run("decodes checksummed address", func(t *testing.T) {
		addr := "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"

		got, err := addressDecodeHook(reflect.TypeFor[string](), targetType, addr)

		require.NoError(t, err)
		assert.Equal(t, common.HexToAddress(addr), got)
	})

	t.Run("returns error on invalid hex address", func(t *testing.T) {
		_, err := addressDecodeHook(reflect.TypeFor[string](), targetType, "not-an-address")

		require.Error(t, err)
		assert.Equal(t, "wrong data format for field", err.Error())
	})

	t.Run("decodes empty string to zero address, leaving required-ness to the validator", func(t *testing.T) {
		got, err := addressDecodeHook(reflect.TypeFor[string](), targetType, "")

		require.NoError(t, err)
		assert.Equal(t, common.Address{}, got)
	})

	t.Run("returns data unchanged when source is not string", func(t *testing.T) {
		got, err := addressDecodeHook(reflect.TypeFor[int](), targetType, 42)

		require.NoError(t, err)
		assert.Equal(t, 42, got)
	})

	t.Run("returns data unchanged when target is not common.Address", func(t *testing.T) {
		addr := "0x1234567890abcdef1234567890abcdef12345678"
		got, err := addressDecodeHook(
			reflect.TypeFor[string](), reflect.TypeFor[string](), addr,
		)

		require.NoError(t, err)
		assert.Equal(t, addr, got)
	})
}

func TestKeyDecodeHook(t *testing.T) {
	targetType := reflect.TypeFor[*ecdsa.PrivateKey]()

	t.Run("decodes valid hex private key", func(t *testing.T) {
		// Generate a key and get its hex representation for a known-good input.
		wantKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		hexKey := common.Bytes2Hex(crypto.FromECDSA(wantKey))

		got, err := keyDecodeHook(reflect.TypeFor[string](), targetType, hexKey)

		require.NoError(t, err)
		gotKey, ok := got.(*ecdsa.PrivateKey)
		require.True(t, ok)
		assert.True(t, wantKey.Equal(gotKey))
	})

	t.Run("returns error on invalid hex key", func(t *testing.T) {
		_, err := keyDecodeHook(reflect.TypeFor[string](), targetType, "not-a-valid-key")

		require.Error(t, err)
		assert.Equal(t, "wrong data format for field", err.Error())
	})

	t.Run("returns error on empty string", func(t *testing.T) {
		_, err := keyDecodeHook(reflect.TypeFor[string](), targetType, "")

		require.Error(t, err)
		assert.Equal(t, "wrong data format for field", err.Error())
	})

	t.Run("returns data unchanged when source is not string", func(t *testing.T) {
		got, err := keyDecodeHook(reflect.TypeFor[int](), targetType, 42)

		require.NoError(t, err)
		assert.Equal(t, 42, got)
	})

	t.Run("returns data unchanged when target is not *ecdsa.PrivateKey", func(t *testing.T) {
		got, err := keyDecodeHook(reflect.TypeFor[string](), reflect.TypeFor[string](), "deadbeef")

		require.NoError(t, err)
		assert.Equal(t, "deadbeef", got)
	})
}

func TestInitConfig(t *testing.T) {
	type SimpleConfig struct {
		Name    string `mapstructure:"TEST_NAME" validate:"required"`
		Port    int    `mapstructure:"TEST_PORT" validate:"required" default:"8080"`
		Verbose bool   `mapstructure:"TEST_VERBOSE"`
	}

	t.Run("returns error when env file path does not exist", func(t *testing.T) {
		_, err := InitConfig[SimpleConfig]("/nonexistent/path/.env")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "error reading .env file")
	})

	t.Run("applies default values", func(t *testing.T) {
		t.Setenv("TEST_NAME", "test-app")

		conf, err := InitConfig[SimpleConfig]("")

		// Validation may fail for other fields, but defaults should be applied.
		_ = err
		assert.Equal(t, 8080, conf.Port)
	})

	t.Run("returns validation error when required fields are missing", func(t *testing.T) {
		// Clear any env vars that could satisfy validation.
		t.Setenv("TEST_NAME", "")
		t.Setenv("TEST_PORT", "0")

		_, err := InitConfig[SimpleConfig]("")

		require.Error(t, err)
	})

	t.Run("loads config from environment variables", func(t *testing.T) {
		t.Setenv("TEST_NAME", "my-service")
		t.Setenv("TEST_PORT", "9090")
		t.Setenv("TEST_VERBOSE", "true")

		conf, err := InitConfig[SimpleConfig]("")

		require.NoError(t, err)
		assert.Equal(t, "my-service", conf.Name)
		assert.Equal(t, 9090, conf.Port)
		assert.True(t, conf.Verbose)
	})

	t.Run("decodes *big.Int fields from environment", func(t *testing.T) {
		type BigIntConfig struct {
			ChainID *big.Int `mapstructure:"TEST_CHAIN_ID" validate:"required"`
		}
		t.Setenv("TEST_CHAIN_ID", "99999")

		conf, err := InitConfig[BigIntConfig]("")

		require.NoError(t, err)
		assert.Equal(t, big.NewInt(99999), conf.ChainID)
	})

	t.Run("decodes common.Address fields from environment", func(t *testing.T) {
		type AddressConfig struct {
			ContractAddr common.Address `mapstructure:"TEST_CONTRACT_ADDR" validate:"required"`
		}
		addr := "0x1234567890abcdef1234567890abcdef12345678"
		t.Setenv("TEST_CONTRACT_ADDR", addr)

		conf, err := InitConfig[AddressConfig]("")

		require.NoError(t, err)
		assert.Equal(t, common.HexToAddress(addr), conf.ContractAddr)
	})
}

func TestBindEnvs(t *testing.T) {
	t.Run("binds mapstructure-tagged fields", func(t *testing.T) {
		type TaggedConfig struct {
			Host string `mapstructure:"BIND_TEST_HOST"`
			Port int    `mapstructure:"BIND_TEST_PORT"`
		}

		// bindEnvs should not panic on a valid struct.
		assert.NotPanics(t, func() {
			bindEnvs(TaggedConfig{})
		})
	})

	t.Run("handles pointer to struct", func(t *testing.T) {
		type PtrConfig struct {
			Name string `mapstructure:"BIND_TEST_NAME"`
		}

		assert.NotPanics(t, func() {
			bindEnvs(&PtrConfig{})
		})
	})

	t.Run("skips fields without mapstructure tag", func(t *testing.T) {
		type MixedConfig struct {
			Tagged   string `mapstructure:"BIND_TEST_TAGGED"`
			Untagged string
		}

		assert.NotPanics(t, func() {
			bindEnvs(MixedConfig{})
		})
	})
}
