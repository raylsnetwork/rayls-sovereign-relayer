package client_test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"testing"
	"time"

	kms "cloud.google.com/go/kms/apiv1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/client"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/client/testutil"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"
)

var (
	gcpIntegrationEnabled bool
	gcpProject            string
	gcpLocation           string
)

func init() {
	flag.BoolVar(
		&gcpIntegrationEnabled,
		"gcp-integration",
		false,
		"Run AWS client integration tests (running those tests incurs a financial cost)",
	)
	flag.StringVar(&gcpProject, "gcp-project", "", "GCP project to use during testing")
	flag.StringVar(&gcpLocation, "gcp-location", "", "GCP location to use during testing")
}

func TestGCPKMSClientWrapperIntegration(t *testing.T) {
	if !gcpIntegrationEnabled {
		return
	}

	if gcpProject == "" || gcpLocation == "" {
		fmt.Println(
			"GCP project and/or GCP location are missing. Specify them using the flags before running integration tests.",
		)
		t.FailNow()
	}

	kmsClient, err := kms.NewKeyManagementClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = kmsClient.Close() }()

	keyRingName := testutil.GenerateRandomIdentifier()
	keyName := testutil.GenerateRandomIdentifier()

	t.Run("returns false when key ring does not exist", func(t *testing.T) {
		wrapper := client.NewGCPKMSClientWrapper(kmsClient, gcpProject, gcpLocation)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		exists, err := wrapper.CheckKeyRingExists(ctx, keyRingName)
		assertGRPCError(t, err)

		require.False(t, exists)
	})

	t.Run("creates key ring", func(t *testing.T) {
		wantResourceName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", gcpProject, gcpLocation, keyRingName)

		wrapper := client.NewGCPKMSClientWrapper(kmsClient, gcpProject, gcpLocation)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		gotResourceName, err := wrapper.CreateKeyRing(ctx, keyRingName)
		assertGRPCError(t, err)

		require.Equal(t, wantResourceName, gotResourceName)
	})

	t.Run("returns true when key ring does exist", func(t *testing.T) {
		wrapper := client.NewGCPKMSClientWrapper(kmsClient, gcpProject, gcpLocation)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		exists, err := wrapper.CheckKeyRingExists(ctx, keyRingName)
		assertGRPCError(t, err)

		require.True(t, exists)
	})

	t.Run("returns false when key does not exist", func(t *testing.T) {
		wrapper := client.NewGCPKMSClientWrapper(kmsClient, gcpProject, gcpLocation)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		exists, err := wrapper.CheckKeyExists(ctx, keyRingName, keyName)
		assertGRPCError(t, err)

		require.False(t, exists)
	})

	t.Run("creates key", func(t *testing.T) {
		wantResourceName := fmt.Sprintf(
			"projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
			gcpProject,
			gcpLocation,
			keyRingName,
			keyName,
		)

		wrapper := client.NewGCPKMSClientWrapper(kmsClient, gcpProject, gcpLocation)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		gotResourceName, err := wrapper.CreateKey(ctx, keyRingName, keyName)
		assertGRPCError(t, err)

		require.Equal(t, wantResourceName, gotResourceName)
	})

	t.Run("returns true when key does exist", func(t *testing.T) {
		wrapper := client.NewGCPKMSClientWrapper(kmsClient, gcpProject, gcpLocation)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		exists, err := wrapper.CheckKeyExists(ctx, keyRingName, keyName)
		assertGRPCError(t, err)

		require.True(t, exists)
	})
}

func assertGRPCError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		if grpcStatus, ok := status.FromError(err); ok {
			t.Fatalf("didn't expect error but got one - grpc status: %v", grpcStatus)
		} else {
			t.Fatalf("didn't expect error but got one: %v", err)
		}
	}
}
