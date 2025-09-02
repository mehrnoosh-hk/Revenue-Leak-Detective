package health

import (
	"context"
	"time"

	"rdl-api/internal/db/repository"
)

// HealthService defines the interface for health check operations
type HealthService interface {
	CheckReadiness(ctx context.Context) error
	CheckLiveness(ctx context.Context) error
}

// healthService implements HealthService
type healthService struct {
	healthRepo repository.HealthRepository
}

// NewHealthService creates a new health service instance
func NewHealthService(healthRepo repository.HealthRepository) HealthService {
	return &healthService{
		healthRepo: healthRepo,
	}
}

// CheckReadiness verifies if the application is ready to serve requests
func (h *healthService) CheckReadiness(ctx context.Context) error {
	if h.healthRepo == nil {
		return ErrDatabaseNotInitialized
	}

	// Use a short timeout for readiness checks
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := h.healthRepo.Ping(ctx); err != nil {
		return ErrDatabaseUnavailable
	}

	return nil
}

// CheckLiveness verifies if the application is alive (no external dependencies)
func (h *healthService) CheckLiveness(ctx context.Context) error {
	// Liveness check should be fast and not depend on external services
	// Just return nil to indicate the application is alive
	return nil
}
