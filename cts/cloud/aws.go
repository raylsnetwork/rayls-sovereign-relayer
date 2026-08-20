package cloud

import (
	"context"
	"errors"
)

type AWSKMSClient interface {
	CreateKey(context.Context) (string, error)
	GetKeyIDFromAlias(ctx context.Context, alias string) (string, error)
	CreateAliasForKeyID(ctx context.Context, aliasName, keyID string) error
}

type AWSKMSInfra struct {
	client AWSKMSClient

	keyAlias string
}

func NewAWSKMSInfra(client AWSKMSClient, keyAlias string) *AWSKMSInfra {
	return &AWSKMSInfra{
		client: client,

		keyAlias: keyAlias,
	}
}

func (a *AWSKMSInfra) GetAWSKey(ctx context.Context) (string, error) {
	keyID, err := a.client.GetKeyIDFromAlias(ctx, a.keyAlias)
	if err == nil {
		return keyID, nil
	} else if !errors.Is(err, ErrNoKey) {
		return "", WrapInAWSKMSClientError("error when getting key from alias", err)
	}

	keyID, err = a.client.CreateKey(ctx)
	if err != nil {
		return "", WrapInAWSKMSClientError("error when trying to create key", err)
	}

	err = a.client.CreateAliasForKeyID(ctx, a.keyAlias, keyID)
	if err != nil {
		return "", WrapInAWSKMSClientError("error when trying to create alias for key", err)
	}

	return keyID, nil
}
