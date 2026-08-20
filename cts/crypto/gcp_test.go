package crypto_test

import (
	"context"
	"flag"
	"fmt"
	"hash/crc32"
	"testing"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/googleapis/gax-go/v2"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/crypto"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var (
	gcpIntegrationEnabled bool
	gcpResourceName       string
)

type StubGCPKMSClient struct {
	ciphertext string
	plaintext  string

	spyCiphertext string
}

func (s *StubGCPKMSClient) Encrypt(
	ctx context.Context,
	req *kmspb.EncryptRequest,
	opts ...gax.CallOption,
) (*kmspb.EncryptResponse, error) {
	return &kmspb.EncryptResponse{
		Ciphertext: []byte(s.ciphertext),
		CiphertextCrc32C: &wrapperspb.Int64Value{
			Value: int64(s.crc32c([]byte(s.ciphertext))),
		},
		VerifiedPlaintextCrc32C: true,
	}, nil
}

func (s *StubGCPKMSClient) Decrypt(
	ctx context.Context,
	req *kmspb.DecryptRequest,
	opts ...gax.CallOption,
) (*kmspb.DecryptResponse, error) {
	s.spyCiphertext = string(req.Ciphertext)

	return &kmspb.DecryptResponse{
		Plaintext: []byte(s.plaintext),
		PlaintextCrc32C: &wrapperspb.Int64Value{
			Value: int64(s.crc32c([]byte(s.plaintext))),
		},
	}, nil
}

func (s *StubGCPKMSClient) crc32c(data []byte) uint32 {
	t := crc32.MakeTable(crc32.Castagnoli)
	return crc32.Checksum(data, t)
}

func init() {
	flag.BoolVar(
		&gcpIntegrationEnabled,
		"gcp-integration",
		false,
		"Run AWS client integration tests (running those tests incurs a financial cost)",
	)
	flag.StringVar(&gcpResourceName, "gcp-resource-name", "", "The GCP KMS resource name to be used in tests")
}

func TestGCPEncryptor(t *testing.T) {
	if !gcpIntegrationEnabled {
		return
	}

	if gcpResourceName == "" {
		fmt.Println("no resource name provided - cannot use GCP KMS without a resource name")
		t.FailNow()
	}

	client, err := kms.NewKeyManagementClient(context.Background())
	if err != nil {
		t.Fatalf("couldn't initialize kms client, %v", err)
	}
	defer func() { _ = client.Close() }()

	t.Run("encrypts plaintext and decrypts it's ciphertext", func(t *testing.T) {
		wantText := "this is a sample text"

		gcpEncryptor := crypto.NewGCPEncryptor(client, gcpResourceName)

		ciphertext, err := gcpEncryptor.Encrypt([]byte(wantText))
		require.Nil(t, err, "expected no error but got one during encrypt")

		gotTextBytes, err := gcpEncryptor.Decrypt(ciphertext)
		require.Nil(t, err, "expected no error but got one during decrypt")

		require.Equal(t, wantText, string(gotTextBytes))
	})
}

func TestGCPEncrypt(t *testing.T) {
	t.Run("adds \"gcp\" prefix to ciphertext", func(t *testing.T) {
		wantPrefix := crypto.GCPIdentifierPrefix
		wantCiphertext := "example-ciphertext"

		gcpResourceName := "example resource name"
		plaintext := "example plaintext"

		stubClient := &StubGCPKMSClient{
			ciphertext: wantCiphertext,
		}
		gcpEncryptor := crypto.NewGCPEncryptor(stubClient, gcpResourceName)

		gotCiphertextWithPrefix, err := gcpEncryptor.Encrypt([]byte(plaintext))
		require.Nil(t, err)

		gotPrefix := string(gotCiphertextWithPrefix[:len(crypto.GCPIdentifierPrefix)])
		gotCiphertext := string(gotCiphertextWithPrefix[len(crypto.GCPIdentifierPrefix):])

		require.Equal(t, wantPrefix, gotPrefix)
		require.Equal(t, wantCiphertext, gotCiphertext)
	})
}

func TestGCPDecrypt(t *testing.T) {
	t.Run("returns ErrMissingGCPIdentifier on missing identifier for cyphertext", func(t *testing.T) {
		wantErr := crypto.ErrMissingGCPIdentifier

		gcpResourceName := "example resource name"
		cyphertext := "example-cyphertext"

		client := &StubGCPKMSClient{}
		encryptor := crypto.NewGCPEncryptor(client, gcpResourceName)

		_, gotErr := encryptor.Decrypt([]byte(cyphertext))

		require.Equal(t, wantErr, gotErr)
	})

	t.Run("removes \"gcp\" prefix from ciphertext", func(t *testing.T) {
		wantCiphertext := "example ciphertext"

		gcpResourceName := "example resource name"
		plaintext := "example plain text"

		stubClient := &StubGCPKMSClient{
			plaintext: plaintext,
		}
		gcpEncryptor := crypto.NewGCPEncryptor(stubClient, gcpResourceName)

		ciphertext := crypto.GCPIdentifierPrefix + wantCiphertext
		_, err := gcpEncryptor.Decrypt([]byte(ciphertext))
		require.Nil(t, err)

		require.Equal(t, wantCiphertext, stubClient.spyCiphertext)
	})
}
