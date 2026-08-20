package crypto

import (
	"context"

	kms "cloud.google.com/go/kms/apiv1"
)

func NewGCPKMSClient() (*kms.KeyManagementClient, error) {
	return kms.NewKeyManagementClient(context.Background())
}
