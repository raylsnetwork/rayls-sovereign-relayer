package crypto

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/service"
)

const (
	AWSIdentifierPrefix = "aws:"
	// awsPrefixLength is the minimum length of the AWS identifier prefix in ciphertext.
	awsPrefixLength = 4
)

var ErrMissingAWSIdentifier = errors.New("missing aws identifier for cyphertext")

type AWSKMSClient interface {
	Encrypt(ctx context.Context, params *kms.EncryptInput, optFns ...func(*kms.Options)) (*kms.EncryptOutput, error)
	Decrypt(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error)
}

type AWSEncryptor struct {
	keyID string

	kms AWSKMSClient
}

var _ service.Encryptor = &AWSEncryptor{}

func NewAWSEncryptor(kms AWSKMSClient, keyID string) *AWSEncryptor {
	return &AWSEncryptor{
		keyID: keyID,
		kms:   kms,
	}
}

func NewDefaultAWSEncryptor(ctx context.Context, profile, keyID string) (*AWSEncryptor, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, fmt.Errorf("loading AWS default config: %w", err)
	}

	kmsClient := kms.NewFromConfig(cfg)

	return &AWSEncryptor{
		keyID: keyID,
		kms:   kmsClient,
	}, nil
}

func (e *AWSEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	encryptInput := &kms.EncryptInput{
		KeyId:     aws.String(e.keyID),
		Plaintext: plaintext,
	}

	result, err := e.kms.Encrypt(context.TODO(), encryptInput)
	if err != nil {
		return nil, fmt.Errorf("encrypting with AWS KMS: %w", err)
	}

	cyphertext := append([]byte(AWSIdentifierPrefix), result.CiphertextBlob...)

	return cyphertext, nil
}

func (e *AWSEncryptor) Decrypt(cyphertext []byte) ([]byte, error) {
	if len(cyphertext[:len(AWSIdentifierPrefix)]) < awsPrefixLength {
		return nil, ErrMissingAWSIdentifier
	}
	if string(cyphertext[:len(AWSIdentifierPrefix)]) != AWSIdentifierPrefix {
		return nil, ErrMissingAWSIdentifier
	}

	decryptInput := &kms.DecryptInput{
		KeyId:          aws.String(e.keyID),
		CiphertextBlob: cyphertext[len(AWSIdentifierPrefix):],
	}

	result, err := e.kms.Decrypt(context.TODO(), decryptInput)
	if err != nil {
		return nil, fmt.Errorf("decrypting with AWS KMS: %w", err)
	}

	return result.Plaintext, nil
}
