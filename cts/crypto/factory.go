package crypto

import (
	"context"
	"errors"
	"fmt"

	kms "cloud.google.com/go/kms/apiv1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/client"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/cloud"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/service"
)

type CleanupFunc func() error

type EncryptorFactory interface {
	New(ctx context.Context) (service.Encryptor, error)
	Cleanup() error
}

func GetEncryptor(
	ctx context.Context,
	name string,
	plaintextFactory, awsFactory, gcpFactory EncryptorFactory,
) (service.Encryptor, CleanupFunc, error) {
	var (
		encryptor service.Encryptor
		cleanup   CleanupFunc
		err       error
	)
	switch name {
	case "plaintext":
		encryptor, err = plaintextFactory.New(ctx)
		cleanup = plaintextFactory.Cleanup
	case "aws":
		encryptor, err = awsFactory.New(ctx)
		cleanup = awsFactory.Cleanup
	case "gcp":
		encryptor, err = gcpFactory.New(ctx)
		cleanup = gcpFactory.Cleanup
	default:
		return nil, nil, errors.New("unsupported encryption service")
	}

	if err != nil {
		return nil, nil, fmt.Errorf("creating %s encryptor: %w", name, err)
	}

	return encryptor, cleanup, nil
}

type PlaintextEncryptorFactory struct{}

func NewPlaintextEncryptorFactory() *PlaintextEncryptorFactory {
	return &PlaintextEncryptorFactory{}
}

func (f *PlaintextEncryptorFactory) New(context.Context) (service.Encryptor, error) {
	return &PlaintextEncryptor{}, nil
}

func (f *PlaintextEncryptorFactory) Cleanup() error {
	return nil
}

type AWSEncryptorFactory struct {
	profile  string
	keyAlias string
}

func NewAWSEncryptorFactory(profile, keyAlias string) *AWSEncryptorFactory {
	return &AWSEncryptorFactory{
		profile:  profile,
		keyAlias: keyAlias,
	}
}

func (f *AWSEncryptorFactory) New(ctx context.Context) (service.Encryptor, error) {
	kmsClientWrapper, err := client.NewDefaultAWSKMSClientWrapper(ctx, f.profile)
	if err != nil {
		return nil, fmt.Errorf("creating AWS KMS client wrapper: %w", err)
	}

	kmsInfra := cloud.NewAWSKMSInfra(kmsClientWrapper, f.keyAlias)

	keyID, err := kmsInfra.GetAWSKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting AWS KMS key: %w", err)
	}

	encryptor, err := NewDefaultAWSEncryptor(ctx, f.profile, keyID)
	if err != nil {
		return nil, fmt.Errorf("creating AWS encryptor: %w", err)
	}

	return encryptor, nil
}

func (f *AWSEncryptorFactory) Cleanup() error {
	return nil
}

type GCPEncryptorFactory struct {
	project   string
	location  string
	keyRing   string
	cryptoKey string

	kmsClient *kms.KeyManagementClient
}

func NewGCPEncryptorFactory(project, location, keyRing, cryptoKey string) *GCPEncryptorFactory {
	return &GCPEncryptorFactory{
		project:   project,
		location:  location,
		keyRing:   keyRing,
		cryptoKey: cryptoKey,
	}
}

func (f *GCPEncryptorFactory) New(ctx context.Context) (service.Encryptor, error) {
	kmsClient, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating GCP KMS client: %w", err)
	}
	f.kmsClient = kmsClient

	kmsClientWrapper := client.NewGCPKMSClientWrapper(kmsClient, f.project, f.location)
	kmsInfra := cloud.NewGCPKMSInfra(kmsClientWrapper, f.keyRing, f.cryptoKey)

	resourceName, err := kmsInfra.GetGCPKey(ctx)
	if err != nil {
		_ = kmsClient.Close()
		return nil, fmt.Errorf("getting GCP KMS key: %w", err)
	}

	encryptor := NewGCPEncryptor(kmsClient, resourceName)
	if err != nil {
		_ = kmsClient.Close()
		return nil, fmt.Errorf("creating GCP encryptor: %w", err)
	}

	return encryptor, nil
}

func (f *GCPEncryptorFactory) Cleanup() error {
	return f.kmsClient.Close()
}
