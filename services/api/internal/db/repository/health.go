// health repository implementation HealthRepository interface which is part of repository interface
package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PoolPinger defines the interface for database connection pools that can be pinged
// This is useful for testing with mock implementations
type PoolPinger interface {
	Ping(ctx context.Context) error
}

// healthRepository implements HealthRepository using sqlc
type HealthRepositoryImplementation struct {
	pool   PoolPinger
	logger *slog.Logger
}

// NewHealthRepository creates a new health repository instance with a pgxpool.Pool
func NewHealthRepositoryImplementation(pool *pgxpool.Pool, logger *slog.Logger) (HealthRepositoryImplementation, error) {
	if logger == nil {
		return HealthRepositoryImplementation{}, ErrLoggerCannotBeNil
	}
	if pool == nil {
		return HealthRepositoryImplementation{}, ErrPoolCannotBeNil
	}
	return HealthRepositoryImplementation{
		pool:   pool,
		logger: logger,
	}, nil
}

// Ping implements HealthRepository.Ping
func (h HealthRepositoryImplementation) CheckReadiness(ctx context.Context) error {
	return h.pool.Ping(ctx)
}
