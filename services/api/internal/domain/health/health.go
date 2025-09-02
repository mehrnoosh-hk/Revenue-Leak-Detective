package health

import (
	"context"
	"time"
)

// HealthService defines the interface for health check operations
type HealthService interface {
	CheckReadiness(ctx context.Context) error
	CheckLiveness(ctx context.Context) error
}

// healthService implements HealthService
type healthService struct {
	dbChecker DatabaseChecker
}

// DatabaseChecker abstracts database operations for health checks
type DatabaseChecker interface {
	Ping(ctx context.Context) error
}

// NewHealthService creates a new health service instance
func NewHealthService(dbChecker DatabaseChecker) HealthService {
	return &healthService{
		dbChecker: dbChecker,
	}
}

// CheckReadiness verifies if the application is ready to serve requests
func (h *healthService) CheckReadiness(ctx context.Context) error {
	if h.dbChecker == nil {
		return ErrDatabaseNotInitialized
	}

	// Use a short timeout for readiness checks
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return h.dbChecker.Ping(ctx)
}

// CheckLiveness verifies if the application is alive (no external dependencies)
func (h *healthService) CheckLiveness(ctx context.Context) error {
	// Liveness check should be fast and not depend on external services
	// Just return nil to indicate the application is alive
	return nil
}
