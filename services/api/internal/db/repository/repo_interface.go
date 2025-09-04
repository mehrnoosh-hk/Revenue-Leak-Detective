package repository

import (
	"context"
	sqlc "rdl-api/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

// Repository interface
type Repository interface {
	HealthRepository
	UserRepository
}

// HealthRepository defines the interface for health-related database operations
type HealthRepository interface {
	// Ping checks if the database is reachable
	Ping(ctx context.Context) error
}

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	DeleteUser(ctx context.Context, id pgtype.UUID) (int64, error)
	GetAllUsers(ctx context.Context) ([]sqlc.User, error)
	// users table queries
	GetUserByEmail(ctx context.Context, email string) (sqlc.User, error)
	GetUserById(ctx context.Context, id pgtype.UUID) (sqlc.User, error)
	UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error)
}

// Database abstracts the database connection
type Database interface {
	Ping(ctx context.Context) error
}
