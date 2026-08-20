package client_test

import (
	"context"
	"errors"
	"flag"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/client"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/client/testutil"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/cloud"
	"github.com/stretchr/testify/require"
)

var (
	awsProfile            string
	awsIntegrationEnabled bool
)

type StubAWSKMSClient struct {
	err error
}

// CreateAlias implements client.AWSKMSClient.
func (s *StubAWSKMSClient) CreateAlias(
	ctx context.Context,
	params *kms.CreateAliasInput,
	optFns ...func(*kms.Options),
) (*kms.CreateAliasOutput, error) {
	return &kms.CreateAliasOutput{}, s.err
}

// DescribeKey implements client.AWSKMSClient.
func (s *StubAWSKMSClient) DescribeKey(
	ctx context.Context,
	params *kms.DescribeKeyInput,
	optFns ...func(*kms.Options),
) (*kms.DescribeKeyOutput, error) {
	return &kms.DescribeKeyOutput{}, s.err
}

func (s *StubAWSKMSClient) CreateKey(
	context.Context,
	*kms.CreateKeyInput,
	...func(*kms.Options),
) (*kms.CreateKeyOutput, error) {
	return &kms.CreateKeyOutput{
		KeyMetadata: &types.KeyMetadata{
			KeyId: new(string),
		},
	}, s.err
}

type SlowAWSKMSClient struct {
	sleepDuration time.Duration
}

func (s *SlowAWSKMSClient) CreateKey(
	ctx context.Context,
	input *kms.CreateKeyInput,
	args ...func(*kms.Options),
) (*kms.CreateKeyOutput, error) {
	select {
	case <-ctx.Done():
	case <-time.NewTimer(s.sleepDuration).C:
	}
	return nil, nil //nolint:nilnil // returning nil,nil is intentional for test stub
}

func init() {
	flag.StringVar(&awsProfile, "aws-profile", "default", "Specify AWS profile to be used during integration tests.")
	flag.BoolVar(
		&awsIntegrationEnabled,
		"aws-integration",
		false,
		"Run AWS client integration tests (running those tests incurs a financial cost)",
	)
}

func TestAWSKMSClientWrapper(t *testing.T) {
	t.Run("forwards error from client", func(t *testing.T) {
		wantErr := errors.New("this is an example error")

		stubClient := StubAWSKMSClient{
			err: wantErr,
		}
		wrapper := client.NewAWSKMSClientWrapper(&stubClient)

		_, gotErr := wrapper.CreateKey(context.Background())

		require.ErrorIs(t, gotErr, wantErr)
	})

	t.Run("exits gracefully on canceled context", func(t *testing.T) {
		stubClient := StubAWSKMSClient{}
		wrapper := client.NewAWSKMSClientWrapper(&stubClient)

		// The grace period is a safety bound, not a target: under normal
		// load the goroutine returns within microseconds of ctx cancel,
		// and the fixture observes that via the shutdownSignal channel.
		// Using a sub-millisecond grace caused intermittent flakes under
		// concurrent CPU load (goroutine scheduling can exceed 1 ms).
		hasGracefulShutdown := gracefulShutdownFixture(wrapper.CreateKey, 2*time.Second)

		require.True(t, hasGracefulShutdown)
	})
}

func TestAWSKMSClientWrapperIntegration(t *testing.T) {
	var (
		keyID string
		alias string
	)

	if !awsIntegrationEnabled {
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(awsProfile))
	if err != nil {
		t.Fatal(err)
	}
	kmsClient := kms.NewFromConfig(cfg)

	alias = "alias/" + testutil.GenerateRandomIdentifier()

	t.Run("returns ErrNoKey on key alias that does not exist", func(t *testing.T) {
		wantErr := cloud.ErrNoKey

		wrapper := client.NewAWSKMSClientWrapper(kmsClient)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, gotErr := wrapper.GetKeyIDFromAlias(ctx, alias)

		require.ErrorIs(t, gotErr, wantErr)
	})

	t.Run("creates new key and returns it's id", func(t *testing.T) {
		wrapper := client.NewAWSKMSClientWrapper(kmsClient)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		keyID, err = wrapper.CreateKey(ctx)

		require.Nil(t, err)
		require.NotEmpty(t, keyID)
	})

	t.Run("creates alias for the key", func(t *testing.T) {
		wrapper := client.NewAWSKMSClientWrapper(kmsClient)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		err = wrapper.CreateAliasForKeyID(ctx, alias, keyID)

		require.Nil(t, err)
	})

	t.Run("returns true on key alias that exists", func(t *testing.T) {
		wrapper := client.NewAWSKMSClientWrapper(kmsClient)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		gotKeyID, err := wrapper.GetKeyIDFromAlias(ctx, alias)
		require.Nil(t, err)

		require.Equal(t, keyID, gotKeyID)
	})
}

func gracefulShutdownFixture(fn func(ctx context.Context) (string, error), gracePeriod time.Duration) bool {
	shutdownSignal := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_, _ = fn(ctx)
		shutdownSignal <- struct{}{}
	}()
	cancel()

	select {
	case <-time.NewTimer(gracePeriod).C:
		return false
	case <-shutdownSignal:
		return true
	}
}
