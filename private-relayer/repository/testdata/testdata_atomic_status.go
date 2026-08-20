// Decommissioning Teleport (vanilla, atomic).

package testdata

import (
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

var (
	AtomicSUM1 = types.AtomicStatusUpdateMessage{
		SharedID: "atomic-shared-id-1",
		Status:   types.AtomicExecutedStatus,
	}
	AtomicSUM2 = types.AtomicStatusUpdateMessage{
		SharedID: "atomic-shared-id-2",
		Status:   types.AtomicRejectedStatus,
	}
	AtomicSUM3 = types.AtomicStatusUpdateMessage{
		SharedID: "atomic-shared-id-3",
		Status:   types.AtomicPendingStatus,
	}
)

var (
	ModelAtomicSUM1 = repository.AtomicSUM{
		SharedId:    AtomicSUM1.SharedID,
		Status:      uint8(AtomicSUM1.Status),
		IsProcessed: false,
	}
	ModelAtomicSUM2 = repository.AtomicSUM{
		SharedId:    AtomicSUM2.SharedID,
		Status:      uint8(AtomicSUM2.Status),
		IsProcessed: false,
	}
	ModelAtomicSUM3Processed = repository.AtomicSUM{
		SharedId:    AtomicSUM3.SharedID,
		Status:      uint8(AtomicSUM3.Status),
		IsProcessed: true,
	}
)
