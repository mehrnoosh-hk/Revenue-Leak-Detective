package services

import (
	"context"
	"log/slog"
	"rdl-api/internal/db/repository"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type healthService struct {
	healthRepo repository.HealthRepository
	logger     *slog.Logger
}

func NewHealthService(p *pgxpool.Pool, logger *slog.Logger) HealthService {
	return &healthService{
		healthRepo: repository.NewHealthRepository(p, logger),
		logger:     logger,
	}
}

func (h *healthService) CheckReadiness(ctx context.Context) error {
	if h.healthRepo == nil {
		return ErrDatabaseNotInitialized
	}

	// Use a short timeout for readiness checks
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := h.healthRepo.Ping(ctx); err != nil {
		return ErrDatabaseUnavailable
	}

	return nil
}

func (h *healthService) CheckLiveness(ctx context.Context) error {
	return nil
}
