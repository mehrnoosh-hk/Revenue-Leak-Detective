package services

import (
	"context"
	"log/slog"
	"rdl-api/internal/db/repository"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthService interface {
	CheckReadiness(ctx context.Context) error
	CheckLiveness(ctx context.Context) error
	GetVersion() string
}

type healthService struct {
	healthRepo HealthRepository
	logger     *slog.Logger
	version    string
}

func NewHealthService(p *pgxpool.Pool, logger *slog.Logger, version string) (HealthService, error) {
	if logger == nil {
		return nil, ErrLoggerCannotBeNil
	}
	if p == nil {
		return nil, ErrPoolCannotBeNil
	}
	healthRepo, err := repository.NewHealthRepositoryImplementation(p, logger)
	if err != nil {
		return nil, err
	}
	return healthService{
		healthRepo: &healthRepo,
		logger:     logger,
		version:    version,
	}, nil
}

func (h healthService) CheckReadiness(ctx context.Context) error {

	// Use a short timeout for readiness checks
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := h.healthRepo.CheckReadiness(ctx); err != nil {
		return ErrDatabaseUnavailable
	}

	return nil
}

func (h healthService) CheckLiveness(ctx context.Context) error {
	return nil
}

func (h healthService) GetVersion() string {
	return h.version
}
