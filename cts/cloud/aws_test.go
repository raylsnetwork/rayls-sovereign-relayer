package cloud_test

import (
	"context"
	"errors"
	"testing"

	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/cloud"
	"github.com/stretchr/testify/require"
)

type StubAWSKMSClient struct {
	key string

	spyAlias string
	spyKeyID string

	getKeyIDFromAliasErr error
	createKeyErr         error
}

func (s *StubAWSKMSClient) CreateKey(ctx context.Context) (string, error) {
	return s.key, s.createKeyErr
}

func (s *StubAWSKMSClient) GetKeyIDFromAlias(ctx context.Context, alias string) (string, error) {
	s.spyAlias = alias
	return s.key, s.getKeyIDFromAliasErr
}

func (s *StubAWSKMSClient) CreateAliasForKeyID(ctx context.Context, alias, keyID string) error {
	s.spyAlias = alias
	s.spyKeyID = keyID
	return nil
}

func TestKMS(t *testing.T) {
	t.Run("checks whether key alias exists and returns key ID if it does", func(t *testing.T) {
		wantAlias := "example-alias"
		wantKey := "this-is-a-sample-key"

		client := &StubAWSKMSClient{
			key: wantKey,
		}
		infra := cloud.NewAWSKMSInfra(client, wantAlias)

		gotKey, err := infra.GetAWSKey(context.Background())
		require.Nil(t, err)

		require.Equal(t, wantAlias, client.spyAlias)
		require.Equal(t, wantKey, gotKey)
	})

	t.Run("creates new key on ErrNoKey from client and sets it's alias", func(t *testing.T) {
		wantAlias := "example-alias"
		wantKey := "this-is-a-sample-key"

		client := &StubAWSKMSClient{
			key:                  wantKey,
			getKeyIDFromAliasErr: cloud.ErrNoKey,
		}
		infra := cloud.NewAWSKMSInfra(client, wantAlias)

		gotKey, err := infra.GetAWSKey(context.Background())
		require.Nil(t, err)

		require.Equal(t, wantAlias, client.spyAlias)
		require.Equal(t, wantKey, client.spyKeyID)

		require.Equal(t, wantKey, gotKey)
	})

	t.Run("wraps unknown error from AWSKMSClient in AWSKMSClientError", func(t *testing.T) {
		wantError := errors.New("example error")

		client := &StubAWSKMSClient{
			getKeyIDFromAliasErr: wantError,
		}
		infra := cloud.NewAWSKMSInfra(client, "example alias")

		_, gotError := infra.GetAWSKey(context.Background())

		require.IsType(t, &cloud.AWSKMSClientError{}, gotError)
		require.Equal(t, wantError, errors.Unwrap(gotError))
	})
}
