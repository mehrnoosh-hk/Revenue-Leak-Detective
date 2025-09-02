package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PgxAdapter implements Database interface for pgxpool.Pool
type PgxAdapter struct {
	pool *pgxpool.Pool
}

// NewPgxAdapter creates a new pgx adapter
func NewPgxAdapter(pool *pgxpool.Pool) Database {
	return &PgxAdapter{
		pool: pool,
	}
}

// Ping implements Database.Ping
func (p *PgxAdapter) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}
