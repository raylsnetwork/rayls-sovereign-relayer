package crypto

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"

	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/googleapis/gax-go/v2"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/service"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	GCPIdentifierPrefix = "gcp:"
)

var (
	ErrMissingGCPIdentifier   = errors.New("missing gcp indentifier for ciphertext")
	ErrReqCorruptedInTransit  = errors.New("request corrupted in transit")
	ErrRespCorruptedInTransit = errors.New("response corrupted in transit")
)

type GCPKMSClient interface {
	Encrypt(context.Context, *kmspb.EncryptRequest, ...gax.CallOption) (*kmspb.EncryptResponse, error)
	Decrypt(context.Context, *kmspb.DecryptRequest, ...gax.CallOption) (*kmspb.DecryptResponse, error)
}

type GCPEncryptor struct {
	client       GCPKMSClient
	resourceName string
}

var _ service.Encryptor = (*GCPEncryptor)(nil)

func NewGCPEncryptor(client GCPKMSClient, resourceName string) *GCPEncryptor {
	return &GCPEncryptor{
		client:       client,
		resourceName: resourceName,
	}
}

func (e *GCPEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	plaintextCRC32C := e.crc32c(plaintext)

	// Build the request.
	req := &kmspb.EncryptRequest{
		Name:            e.resourceName,
		Plaintext:       plaintext,
		PlaintextCrc32C: wrapperspb.Int64(int64(plaintextCRC32C)),
	}
	result, err := e.client.Encrypt(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt: %w", err)
	}

	if !result.VerifiedPlaintextCrc32C {
		return nil, ErrReqCorruptedInTransit
	}
	if int64(e.crc32c(result.Ciphertext)) != result.CiphertextCrc32C.Value {
		return nil, ErrRespCorruptedInTransit
	}

	gcpCiphertext := append([]byte(GCPIdentifierPrefix), result.Ciphertext...)

	return gcpCiphertext, nil
}

func (e *GCPEncryptor) Decrypt(gcpCiphertext []byte) ([]byte, error) {
	if len(gcpCiphertext) < len(GCPIdentifierPrefix) {
		return nil, ErrMissingGCPIdentifier
	}
	if string(gcpCiphertext[:len(GCPIdentifierPrefix)]) != GCPIdentifierPrefix {
		return nil, ErrMissingGCPIdentifier
	}

	ciphertext := gcpCiphertext[len(GCPIdentifierPrefix):]

	ciphertextCRC32C := e.crc32c(ciphertext)
	req := &kmspb.DecryptRequest{
		Name:             e.resourceName,
		Ciphertext:       ciphertext,
		CiphertextCrc32C: wrapperspb.Int64(int64(ciphertextCRC32C)),
	}

	result, err := e.client.Decrypt(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt ciphertext: %w", err)
	}

	if int64(e.crc32c(result.Plaintext)) != result.PlaintextCrc32C.Value {
		return nil, ErrRespCorruptedInTransit
	}

	return result.Plaintext, nil
}

func (e *GCPEncryptor) crc32c(data []byte) uint32 {
	t := crc32.MakeTable(crc32.Castagnoli)
	return crc32.Checksum(data, t)
}
