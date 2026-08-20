package client

import (
	"context"
	"fmt"

	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

type GCPKMSClient interface {
	GetKeyRing(context.Context, *kmspb.GetKeyRingRequest, ...gax.CallOption) (*kmspb.KeyRing, error)
	CreateKeyRing(context.Context, *kmspb.CreateKeyRingRequest, ...gax.CallOption) (*kmspb.KeyRing, error)
	GetCryptoKey(context.Context, *kmspb.GetCryptoKeyRequest, ...gax.CallOption) (*kmspb.CryptoKey, error)
	CreateCryptoKey(
		ctx context.Context,
		req *kmspb.CreateCryptoKeyRequest,
		opts ...gax.CallOption,
	) (*kmspb.CryptoKey, error)
}

type GCPKMSClientWrapper struct {
	client GCPKMSClient
	parent string
}

func NewGCPKMSClientWrapper(client GCPKMSClient, project, location string) *GCPKMSClientWrapper {
	return &GCPKMSClientWrapper{
		client: client,
		parent: fmt.Sprintf("projects/%s/locations/%s", project, location),
	}
}

func (w *GCPKMSClientWrapper) CheckKeyRingExists(ctx context.Context, name string) (bool, error) {
	_, err := w.client.GetKeyRing(ctx, &kmspb.GetKeyRingRequest{
		Name: fmt.Sprintf("%s/keyRings/%s", w.parent, name),
	})
	if err != nil {
		if grpcStatus, ok := status.FromError(err); ok {
			if grpcStatus.Code() == codes.NotFound {
				return false, nil
			} else {
				return false, withstack.Wrap(fmt.Errorf("checking GCP key ring %s: %w", name, err))
			}
		} else {
			return false, withstack.Wrap(fmt.Errorf("checking GCP key ring %s: %w", name, err))
		}
	}
	return true, nil
}

func (w *GCPKMSClientWrapper) CreateKeyRing(ctx context.Context, name string) (string, error) {
	keyRing, err := w.client.CreateKeyRing(ctx, &kmspb.CreateKeyRingRequest{
		Parent:    w.parent,
		KeyRingId: name,
		KeyRing:   &kmspb.KeyRing{},
	})

	if err != nil {
		return "", withstack.Wrap(fmt.Errorf("creating GCP key ring %s: %w", name, err))
	}
	return keyRing.Name, nil
}

func (w *GCPKMSClientWrapper) CheckKeyExists(ctx context.Context, keyRingName, keyName string) (bool, error) {
	_, err := w.client.GetCryptoKey(ctx, &kmspb.GetCryptoKeyRequest{
		Name: fmt.Sprintf("%s/keyRings/%s/cryptoKeys/%s", w.parent, keyRingName, keyName),
	})
	if err != nil {
		if grpcStatus, ok := status.FromError(err); ok {
			if grpcStatus.Code() == codes.NotFound {
				return false, nil
			} else {
				return false, withstack.Wrap(fmt.Errorf("checking GCP crypto key %s/%s: %w", keyRingName, keyName, err))
			}
		} else {
			return false, withstack.Wrap(fmt.Errorf("checking GCP crypto key %s/%s: %w", keyRingName, keyName, err))
		}
	}
	return true, nil
}

func (w *GCPKMSClientWrapper) CreateKey(ctx context.Context, keyRingName, keyName string) (string, error) {
	key, err := w.client.CreateCryptoKey(ctx, &kmspb.CreateCryptoKeyRequest{
		Parent:      fmt.Sprintf("%s/keyRings/%s", w.parent, keyRingName),
		CryptoKeyId: keyName,
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ENCRYPT_DECRYPT,
		},
	})

	if err != nil {
		return "", withstack.Wrap(fmt.Errorf("creating GCP crypto key %s/%s: %w", keyRingName, keyName, err))
	}
	return key.Name, nil
}
