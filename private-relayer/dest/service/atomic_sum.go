// Decommissioning Teleport (vanilla, atomic).

package service

import (
	"context"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

type sumAtomicStatusRepository interface {
	BatchCreate(ctx context.Context, atomicStatuses []types.AtomicStatusUpdateMessage) error
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type SUMService struct {
	atomicStatusRepository sumAtomicStatusRepository
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func NewSUMService(atomicStatusRepository sumAtomicStatusRepository) *SUMService {
	return &SUMService{
		atomicStatusRepository: atomicStatusRepository,
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (s *SUMService) BatchCreate(ctx context.Context, sums []types.AtomicStatusUpdateMessage) error {
	return s.atomicStatusRepository.BatchCreate(ctx, sums)
}
