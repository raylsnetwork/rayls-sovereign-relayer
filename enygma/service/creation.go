package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"

	telemetry "github.com/raylsnetwork/rayls-privacy-relayer-api/otel"
	repository "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"go.opentelemetry.io/otel/trace"
)

type creationServiceTracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

type creationServiceRepositoryManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type creationServiceEnygmaRepository interface {
	CreateEnygma(ctx context.Context, enygma types.Enygma) error
}

type creationServiceEnygmaHistoryRepository interface {
	InsertEnygmaHistory(ctx context.Context, history types.EnygmaHistory) error
}

type EnygmaCreationService struct {
	tracer                  creationServiceTracer
	repositoryManager       creationServiceRepositoryManager
	enygmaRepository        creationServiceEnygmaRepository
	enygmaHistoryRepository creationServiceEnygmaHistoryRepository
}

func NewEnygmaCreationService(
	tracer creationServiceTracer,
	repositoryManager creationServiceRepositoryManager,
	enygmaRepository creationServiceEnygmaRepository,
	enygmaHistoryRepository creationServiceEnygmaHistoryRepository,
) *EnygmaCreationService {
	return &EnygmaCreationService{
		tracer:                  tracer,
		repositoryManager:       repositoryManager,
		enygmaRepository:        enygmaRepository,
		enygmaHistoryRepository: enygmaHistoryRepository,
	}
}

func (s *EnygmaCreationService) CreateEnygma(
	ctx context.Context,
	resourceId string,
	fromChainId *big.Int,
	blockNumberPrivateHub *big.Int,
) error {
	ctx, span := s.tracer.Start(ctx, telemetry.SPAN_CREATE_ENYGMA)
	defer span.End()

	err := s.repositoryManager.WithTransaction(ctx, func(txCtx context.Context) error {
		rFactor := big.NewInt(0)
		balance := big.NewInt(0)

		err := s.enygmaHistoryRepository.InsertEnygmaHistory(
			txCtx,
			types.EnygmaHistory{
				ResourceId:            resourceId,
				FromChainId:           fromChainId,
				BalanceChange:         balance,
				RFactor:               rFactor,
				BlockNumberPrivateHub: blockNumberPrivateHub,
				EventType:             types.EnygmaCreation,
			},
		)
		if err != nil {
			if errors.Is(err, repository.ErrAlreadyProcessed) {
				slog.Info("Skipping already-processed enygma creation", slog.String("resourceId", resourceId))
				return nil
			}
			return fmt.Errorf("insert enygma history: %w", err)
		}

		err = s.enygmaRepository.CreateEnygma(
			txCtx,
			types.Enygma{
				ResourceId:           resourceId,
				FinalizedR:           rFactor,
				FinalizedBalance:     balance,
				FinalizedBlockNumber: blockNumberPrivateHub,
				PendingBlockNumber:   blockNumberPrivateHub,
			},
		)
		if err != nil {
			return fmt.Errorf("create enygma: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("executing enygma creation transaction for resource %s: %w", resourceId, err)
	}

	return nil
}
