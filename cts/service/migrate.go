package service

import (
	"context"
	"fmt"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/domain"
)

type Encryptable[T any] interface {
	Encrypt(enc domain.Encryptor) (T, error)
}

type Decryptable[T any] interface {
	Decrypt(enc domain.Encryptor) (T, error)
}

type MigrationRepository[T any] interface {
	CreateAll(context.Context, []T) error
	GetAll(context.Context) ([]T, error)
}

func MigrateRepository[T Decryptable[E], E Encryptable[T]](
	ctx context.Context,
	fromRepo, toRepo MigrationRepository[T],
	fromEnc, toEnc domain.Encryptor,
) error {
	fromRecords, err := fromRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("fetching records from source repository: %w", err)
	}

	toRecords := make([]T, 0, len(fromRecords))
	for _, record := range fromRecords {
		plainRecord, err := record.Decrypt(fromEnc)
		if err != nil {
			return fmt.Errorf("failed to decrypt record: %w", err)
		}

		encrRecord, err := plainRecord.Encrypt(toEnc)
		if err != nil {
			return fmt.Errorf("failed to encrypt record: %w", err)
		}

		toRecords = append(toRecords, encrRecord)
	}

	return toRepo.CreateAll(ctx, toRecords)
}
