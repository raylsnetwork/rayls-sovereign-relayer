package cloud

import "context"

type GCPKMSClient interface {
	CheckKeyRingExists(ctx context.Context, name string) (bool, error)
	CreateKeyRing(ctx context.Context, name string) (string, error)
	CheckKeyExists(ctx context.Context, keyRingName, keyName string) (bool, error)
	CreateKey(ctx context.Context, keyRingName, keyName string) (string, error)
}

type GCPKMSInfra struct {
	client GCPKMSClient

	keyRing   string
	cryptoKey string
}

func NewGCPKMSInfra(client GCPKMSClient, keyRing, cryptoKey string) *GCPKMSInfra {
	return &GCPKMSInfra{
		client:    client,
		keyRing:   keyRing,
		cryptoKey: cryptoKey,
	}
}

func (i *GCPKMSInfra) GetGCPKey(ctx context.Context) (string, error) {
	var resourceName string

	keyRingExists, err := i.client.CheckKeyRingExists(ctx, i.keyRing)
	if err != nil {
		return "", WrapInGCPKMSClientError("error when trying to check if key ring exists", err)
	}
	if !keyRingExists {
		_, err = i.client.CreateKeyRing(ctx, i.keyRing)
		if err != nil {
			return "", WrapInGCPKMSClientError("error when trying to create key ring", err)
		}
	}

	keyExists, err := i.client.CheckKeyExists(ctx, i.keyRing, i.cryptoKey)
	if err != nil {
		return "", WrapInGCPKMSClientError("error when trying to check if key exists", err)
	}
	if !keyExists {
		resourceName, err = i.client.CreateKey(ctx, i.keyRing, i.cryptoKey)
		if err != nil {
			return "", WrapInGCPKMSClientError("error when trying to create key", err)
		}
	}

	return resourceName, nil
}
