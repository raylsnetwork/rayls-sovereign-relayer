// Decommissioning Teleport (vanilla, atomic).

package app

import (
	sharedrepository "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/public-relayer/repository"
)

type Repositories struct {
	lastProcessedBlockRepository *sharedrepository.LastProcessedBlockRepository

	publicRevertSignature  *repository.RevertSignatureRepository
	privateRevertSignature *repository.RevertSignatureRepository

	publicMessageRecord  *repository.MessageRecordRepository
	privateMessageRecord *repository.MessageRecordRepository
}

func (p *PublicRelayer) initializeRepositories() {
	p.repositories = &Repositories{
		lastProcessedBlockRepository: sharedrepository.NewLastProcessedBlockRepository(p.pool),

		publicRevertSignature:  repository.NewRevertSignatureRepository("public_revert_signature", p.pool),
		privateRevertSignature: repository.NewRevertSignatureRepository("private_revert_signature", p.pool),

		publicMessageRecord:  repository.NewMessageRecordRepository("public_message_record", p.pool),
		privateMessageRecord: repository.NewMessageRecordRepository("private_message_record", p.pool),
	}
}
