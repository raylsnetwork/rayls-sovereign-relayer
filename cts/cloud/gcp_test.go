package cloud_test

import (
	"context"
	"errors"
	"testing"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/cloud"
	"github.com/stretchr/testify/require"
)

type StubGCPKMSClient struct {
	keyRingExists bool
	keyExists     bool

	spyCheckKeyRingExists bool
	spyCheckKeyExists     bool

	spyKeyRing   string
	spyCryptoKey string

	parentName   string //  parent up to the crypto key
	resourceName string //  parent with crypto key included

	err error
}

func (c *StubGCPKMSClient) CheckKeyRingExists(ctx context.Context, name string) (bool, error) {
	c.spyCheckKeyRingExists = true
	c.spyKeyRing = name

	return c.keyRingExists, c.err
}

func (c *StubGCPKMSClient) CreateKeyRing(ctx context.Context, name string) (string, error) {
	c.spyKeyRing = name
	return c.parentName, c.err
}

func (c *StubGCPKMSClient) CheckKeyExists(ctx context.Context, keyRing, cryptoKey string) (bool, error) {
	c.spyCheckKeyExists = true

	c.spyKeyRing = keyRing
	c.spyCryptoKey = cryptoKey

	return c.keyRingExists, c.err
}

func (c *StubGCPKMSClient) CreateKey(ctx context.Context, keyRing, cryptoKey string) (string, error) {
	c.spyKeyRing = keyRing
	c.spyCryptoKey = cryptoKey
	return c.resourceName, c.err
}

func TestGCPKMSInfra(t *testing.T) {
	t.Run("checks whether key ring and key exists", func(t *testing.T) {
		keyRing := "example key ring"
		cryptoKey := "example crypto key"

		client := &StubGCPKMSClient{}
		infra := cloud.NewGCPKMSInfra(client, keyRing, cryptoKey)

		_, err := infra.GetGCPKey(context.Background())
		require.Nil(t, err)

		require.True(t, client.spyCheckKeyRingExists)
		require.True(t, client.spyCheckKeyRingExists)
	})

	t.Run("creates key ring and key if it does not exist", func(t *testing.T) {
		wantKeyRing := "example key ring"
		wantCryptoKey := "example crypto key"

		client := &StubGCPKMSClient{
			keyRingExists: false,
			keyExists:     false,
		}
		infra := cloud.NewGCPKMSInfra(client, wantKeyRing, wantCryptoKey)

		_, err := infra.GetGCPKey(context.Background())
		require.Nil(t, err)

		require.Equal(t, wantKeyRing, client.spyKeyRing)
		require.Equal(t, wantCryptoKey, client.spyCryptoKey)
	})

	t.Run("returns resource name", func(t *testing.T) {
		keyRing := "example key ring"
		cryptoKey := "example crypto key"

		wantResourceName := "example-resource-name"

		client := &StubGCPKMSClient{
			resourceName: wantResourceName,
		}
		infra := cloud.NewGCPKMSInfra(client, keyRing, cryptoKey)

		gotResourceName, err := infra.GetGCPKey(context.Background())
		require.Nil(t, err)

		require.Equal(t, wantResourceName, gotResourceName)
	})

	t.Run("wraps errors from gcp client in GCPKMSClientError", func(t *testing.T) {
		keyRing := "example key ring"
		cryptoKey := "example crypto key"

		wantError := errors.New("example error")

		client := &StubGCPKMSClient{
			err: wantError,
		}
		infra := cloud.NewGCPKMSInfra(client, keyRing, cryptoKey)

		_, gotError := infra.GetGCPKey(context.Background())

		require.IsType(t, &cloud.GCPKMSClientError{}, gotError)
		require.Equal(t, wantError, errors.Unwrap(gotError))
	})
}
