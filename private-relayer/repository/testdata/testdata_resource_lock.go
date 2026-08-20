package testdata

import (
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

var (
	// Active ResourceLock (not expired)
	ResourceLock1 = types.ResourceLock{
		ResourceId: "resource-123",
		ExpiresAt:  time.Now().UTC().Add(1 * time.Hour).Truncate(time.Millisecond),
	}

	// Expired ResourceLock
	ResourceLock2 = types.ResourceLock{
		ResourceId: "resource-456",
		ExpiresAt:  time.Now().UTC().Add(-1 * time.Hour).Truncate(time.Millisecond),
	}

	// Soon to expire ResourceLock
	ResourceLock3 = types.ResourceLock{
		ResourceId: "resource-789",
		ExpiresAt:  time.Now().UTC().Add(5 * time.Minute).Truncate(time.Millisecond),
	}
)
