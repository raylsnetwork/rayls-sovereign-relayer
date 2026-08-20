package crypto

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

func NewAWSKMSClient(ctx context.Context) (*kms.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("loading AWS default config: %w", err))
	}

	return kms.NewFromConfig(cfg), nil
}
