package repository

import (
	"context"

	"github.com/google/uuid"

	models "rdl-api/internal/domain/models"
)

// Repository interface
type Repository interface {
	HealthRepository
	UserRepository
	EventsRepository
	ActionsRepository
}

// HealthRepository defines the interface for health-related database operations
type HealthRepository interface {
	// Ping checks if the database is reachable
	Ping(ctx context.Context) error
}

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	CreateUser(ctx context.Context, arg models.CreateUserParams) (models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) (int64, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	// users table queries
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (models.User, error)
	UpdateUser(ctx context.Context, arg models.UpdateUserParams) (models.User, error)
}

// EventsRepository defines the interface for events-related database operations
type EventsRepository interface {
	CreateEvent(ctx context.Context, arg models.CreateEventParams) (models.Event, error)
	DeleteEvent(ctx context.Context, id uuid.UUID) (int64, error)
	GetAllEvents(ctx context.Context) ([]models.Event, error)
	GetEventById(ctx context.Context, id uuid.UUID) (models.Event, error)
	UpdateEvent(ctx context.Context, arg models.UpdateEventParams) (models.Event, error)
}

// ActionsRepository defines the interface for actions-related database operations
type ActionsRepository interface {
	CreateAction(ctx context.Context, arg models.CreateActionParams) (models.Action, error)
	DeleteAction(ctx context.Context, id uuid.UUID) (int64, error)
	GetAllActionsForTenant(ctx context.Context, tenantID uuid.UUID) ([]models.Action, error)
}

// Database abstracts the database connection
type Database interface {
	Ping(ctx context.Context) error
}

type UUID struct {
	Bytes [16]byte
	Valid bool
}
