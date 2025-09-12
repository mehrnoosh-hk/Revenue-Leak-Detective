package repository

import (
	"context"
	"errors"
	db "rdl-api/internal/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrFailedToAcquireConnection = errors.New("failed to acquire connection")
	ErrFailedToSetTenantID       = errors.New("failed to set tenant ID")
	ErrFailedToSetServiceAccount = errors.New("failed to set service account")
)

// WithTenantContext executes a function with tenant context set
func WithTenantContext(ctx context.Context, pool *pgxpool.Pool, tenantID uuid.UUID, fn func(*db.Queries) error) error {
	// Get a connection from the pool
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return ErrFailedToAcquireConnection
	}
	defer conn.Release()

	// Set the tenant ID in the session
	_, err = conn.Exec(ctx, "SET app.current_tenant_id = $1", tenantID.String())
	if err != nil {
		return ErrFailedToSetTenantID
	}

	// Set service account flag to false for regular operations
	_, err = conn.Exec(ctx, "SET app.is_service_account = false")
	if err != nil {
		return ErrFailedToSetServiceAccount
	}

	// Create a new Queries instance with the connection that has the session context
	tenantQueries := db.New(conn.Conn())
	return fn(tenantQueries)
}

// WithServiceAccountContext returns a Queries instance with service account privileges
// It uses a service account privileged connection to create a new instance of queries
func WithServiceAccountContext(ctx context.Context, pool *pgxpool.Pool) (serviceQueries *db.Queries, err error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, ErrFailedToAcquireConnection
	}
	defer conn.Release()

	// Set service account flag to true
	_, err = conn.Exec(ctx, "SET app.is_service_account = true")
	if err != nil {
		return nil, ErrFailedToSetServiceAccount
	}

	// Create a new Queries instance with the connection that has the session context
	serviceQueries = db.New(conn.Conn())

	return serviceQueries, nil
}
