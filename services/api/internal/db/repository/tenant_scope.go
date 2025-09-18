package repository

import (
	"context"
	db "rdl-api/internal/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WithTenantContext executes a function with tenant context set
func WithTenantContext(ctx context.Context, pool *pgxpool.Pool, tenantID uuid.UUID, fn func(*db.Queries) error) error {
	// Get a connection from the pool
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return ErrFailedToAcquireConnection
	}
	defer conn.Release()

	// Begin a transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Will be no-op if committed

	// Set the tenant ID in the session
	_, err = tx.Exec(ctx, "SET LOCAL app.current_tenant_id = $1", tenantID.String())
	if err != nil {
		return ErrSettingTenantID
	}

	// Set service account flag to false for regular operations
	_, err = tx.Exec(ctx, "SET LOCAL app.is_service_account = false")
	if err != nil {
		return ErrFailedToSetServiceAccount
	}

	// Create a new Queries instance with the connection that has the session context
	tenantQueries := db.New(tx)
	err = fn(tenantQueries)
	if err != nil {
		return err
	}

	return nil
}
