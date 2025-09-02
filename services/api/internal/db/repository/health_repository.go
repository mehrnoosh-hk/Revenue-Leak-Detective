package repository

import (
	"context"
)

// HealthRepository defines the interface for health-related database operations
type HealthRepository interface {
	// Ping checks if the database is reachable
	Ping(ctx context.Context) error
}

// healthRepository implements HealthRepository using sqlc
type healthRepository struct {
	db Database
}

// Database abstracts the database connection for health operations
type Database interface {
	Ping(ctx context.Context) error
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
