package crypto_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/crypto"
	"github.com/stretchr/testify/require"
)

type StubAWSKMSClient struct {
	plaintext string

	calledEncrypt bool
	calledDecrypt bool

	spyKeyID      string
	spyCyphertext string

	err error
}

func (s *StubAWSKMSClient) Encrypt(
	ctx context.Context,
	params *kms.EncryptInput,
	optFns ...func(*kms.Options),
) (*kms.EncryptOutput, error) {
	s.calledEncrypt = true

	s.spyKeyID = *params.KeyId
	s.plaintext = string(params.Plaintext)

	if s.err != nil {
		return nil, s.err
	}

	return &kms.EncryptOutput{
		CiphertextBlob: []byte("example cyphertext"),
	}, nil
}

func (s *StubAWSKMSClient) Decrypt(
	ctx context.Context,
	params *kms.DecryptInput,
	optFns ...func(*kms.Options),
) (*kms.DecryptOutput, error) {
	s.calledDecrypt = true
	s.spyKeyID = *params.KeyId
	s.spyCyphertext = string(params.CiphertextBlob)

	if s.err != nil {
		return nil, s.err
	}

	return &kms.DecryptOutput{
		Plaintext: []byte(s.plaintext),
	}, nil
}

func TestAWSEncryptAndDecrypt(t *testing.T) {
	t.Run("encrypts plaintext and then decrypts it", func(t *testing.T) {
		wantPlaintext := "example plaintext"

		keyID := "example key id"
		client := &StubAWSKMSClient{}
		encryptor := crypto.NewAWSEncryptor(client, keyID)

		cyphertext, err := encryptor.Encrypt([]byte(wantPlaintext))

		require.Nil(t, err)

		gotPlaintextBytes, err := encryptor.Decrypt(cyphertext)

		require.Nil(t, err)
		require.Equal(t, wantPlaintext, string(gotPlaintextBytes))
	})
}

func TestAWSEncrypt(t *testing.T) {
	t.Run("calls encrypt with key ID", func(t *testing.T) {
		wantKeyID := "example key ID"

		plaintext := "example plaintext"
		client := &StubAWSKMSClient{
			calledEncrypt: false,
			calledDecrypt: false,
		}
		encryptor := crypto.NewAWSEncryptor(client, wantKeyID)

		_, err := encryptor.Encrypt([]byte(plaintext))

		require.Nil(t, err)

		require.True(t, client.calledEncrypt)
		require.False(t, client.calledDecrypt)
		require.Equal(t, wantKeyID, client.spyKeyID)
	})

	t.Run("adds \"aws\" prefix to cyphertext", func(t *testing.T) {
		wantPrefix := crypto.AWSIdentifierPrefix
		keyID := "example key ID"

		plaintext := "example plaintext"
		client := &StubAWSKMSClient{}
		encryptor := crypto.NewAWSEncryptor(client, keyID)

		gotCyphertextBytes, _ := encryptor.Encrypt([]byte(plaintext))

		gotPrefix := string(gotCyphertextBytes)[:len(wantPrefix)]
		require.Equal(t, wantPrefix, gotPrefix)
	})

	t.Run("forwards error on fail to encrypt", func(t *testing.T) {
		wantErr := errors.New("example error")

		keyID := "example key ID"
		plaintext := "example plaintext"
		client := &StubAWSKMSClient{
			err: wantErr,
		}
		encryptor := crypto.NewAWSEncryptor(client, keyID)

		_, gotErr := encryptor.Encrypt([]byte(plaintext))

		require.ErrorIs(t, gotErr, wantErr)
	})
}

func TestAWSDecrypt(t *testing.T) {
	t.Run("returns ErrMissingAWSIdentifier on missing identifier for cyphertext", func(t *testing.T) {
		wantErr := crypto.ErrMissingAWSIdentifier

		keyID := "example key ID"
		cyphertext := "example-cyphertext"

		client := &StubAWSKMSClient{}
		encryptor := crypto.NewAWSEncryptor(client, keyID)

		_, gotErr := encryptor.Decrypt([]byte(cyphertext))

		require.ErrorIs(t, gotErr, wantErr)
	})

	t.Run("removes \"aws\" prefix from the cyphertext", func(t *testing.T) {
		wantCyphertext := "example-cyphertext"

		keyID := "example key ID"

		client := &StubAWSKMSClient{}
		encryptor := crypto.NewAWSEncryptor(client, keyID)

		cyphertext := crypto.AWSIdentifierPrefix + wantCyphertext
		_, err := encryptor.Decrypt([]byte(cyphertext))
		require.Nil(t, err)

		require.Equal(t, wantCyphertext, client.spyCyphertext)
	})

	t.Run("calls decrypt with key ID", func(t *testing.T) {
		wantKeyID := "example key ID"

		cyphertext := crypto.AWSIdentifierPrefix + "example cyphertext"
		client := &StubAWSKMSClient{
			calledEncrypt: false,
			calledDecrypt: false,
		}
		encryptor := crypto.NewAWSEncryptor(client, wantKeyID)

		_, err := encryptor.Decrypt([]byte(cyphertext))

		require.Nil(t, err)

		require.True(t, client.calledDecrypt)
		require.False(t, client.calledEncrypt)
		require.Equal(t, wantKeyID, client.spyKeyID)
	})

	t.Run("forwards error on fail to decrypt", func(t *testing.T) {
		wantErr := errors.New("example error")

		keyID := "example key ID"
		cyphertext := crypto.AWSIdentifierPrefix + "example cyphertext"
		client := &StubAWSKMSClient{
			err: wantErr,
		}
		encryptor := crypto.NewAWSEncryptor(client, keyID)

		_, gotErr := encryptor.Decrypt([]byte(cyphertext))

		require.ErrorIs(t, gotErr, wantErr)
	})
}
