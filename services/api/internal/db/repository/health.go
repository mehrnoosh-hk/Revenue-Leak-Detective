// health repository implementation HealthRepository interface which is part of repository interface
package repository

import (
	"context"
)

// healthRepository implements HealthRepository using sqlc
type healthRepository struct {
	db Database
}

// NewHealthRepository creates a new health repository instance
func NewHealthRepository(db Database) HealthRepository {
	return &healthRepository{
		db: db,
	}
}

// Ping implements HealthRepository.Ping
func (h *healthRepository) Ping(ctx context.Context) error {
	return h.db.Ping(ctx)
}
