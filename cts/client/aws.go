package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/cloud"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

type AWSKMSClient interface {
	CreateKey(context.Context, *kms.CreateKeyInput, ...func(*kms.Options)) (*kms.CreateKeyOutput, error)
	DescribeKey(
		ctx context.Context,
		params *kms.DescribeKeyInput,
		optFns ...func(*kms.Options),
	) (*kms.DescribeKeyOutput, error)
	CreateAlias(
		ctx context.Context,
		params *kms.CreateAliasInput,
		optFns ...func(*kms.Options),
	) (*kms.CreateAliasOutput, error)
}

type AWSKMSClientWrapper struct {
	client AWSKMSClient
}

var _ cloud.AWSKMSClient = &AWSKMSClientWrapper{}

func NewAWSKMSClientWrapper(client AWSKMSClient) *AWSKMSClientWrapper {
	return &AWSKMSClientWrapper{
		client: client,
	}
}

func NewDefaultAWSKMSClientWrapper(ctx context.Context, profile string) (*AWSKMSClientWrapper, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("loading AWS config with profile %s: %w", profile, err))
	}

	client := kms.NewFromConfig(cfg)

	return &AWSKMSClientWrapper{
		client: client,
	}, nil
}

func (w *AWSKMSClientWrapper) CreateKey(ctx context.Context) (string, error) {
	input := &kms.CreateKeyInput{
		Description: aws.String("Generated key for encryption"),
	}

	result, err := w.client.CreateKey(ctx, input)
	if err != nil {
		return "", withstack.Wrap(fmt.Errorf("creating AWS KMS key: %w", err))
	}

	return *result.KeyMetadata.KeyId, nil
}

func (w *AWSKMSClientWrapper) GetKeyIDFromAlias(ctx context.Context, alias string) (string, error) {
	input := &kms.DescribeKeyInput{
		KeyId: &alias,
	}

	result, err := w.client.DescribeKey(ctx, input)
	if err != nil {
		notFoundException := &types.NotFoundException{}
		if errors.As(err, &notFoundException) {
			return "", cloud.ErrNoKey
		}
		return "", withstack.Wrap(fmt.Errorf("describing AWS KMS key %s: %w", alias, err))
	}

	return *result.KeyMetadata.KeyId, nil
}

func (w *AWSKMSClientWrapper) CreateAliasForKeyID(ctx context.Context, aliasName, keyID string) error {
	input := &kms.CreateAliasInput{
		AliasName:   &aliasName,
		TargetKeyId: &keyID,
	}

	_, err := w.client.CreateAlias(ctx, input)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("creating AWS KMS alias %s: %w", aliasName, err))
	}
	return nil
}
