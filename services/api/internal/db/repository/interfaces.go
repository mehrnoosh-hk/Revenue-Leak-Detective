package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	models "rdl-api/internal/domain/models"
)

// TenantAwareRepository interface that contains all the repositories that are tenant aware
type TenantAwareRepository interface {
	// HealthRepository
	UserRepository
	EventsRepository
	// ActionsRepository
}

// HealthRepository defines the interface for health-related database operations
type HealthRepository interface {
	// Ping checks if the database is reachable
	Ping(ctx context.Context) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, arg models.CreateUserParams, tenantID uuid.UUID) (models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (int64, error)
	GetAllUsers(ctx context.Context, tenantID uuid.UUID) ([]models.User, error)
	GetUserByEmail(ctx context.Context, email string, tenantID uuid.UUID) (models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (models.User, error)
	UpdateUser(ctx context.Context, arg models.UpdateUserParams, tenantID uuid.UUID) (models.User, error)
}

// EventsRepository defines the interface for events-related database operations
type EventsRepository interface {
	CreateEvent(ctx context.Context, arg models.CreateEventParams, tenantID uuid.UUID) (models.Event, error)
	DeleteEvent(ctx context.Context, eventID uuid.UUID, tenantID uuid.UUID) (int64, error)
	GetAllEvents(ctx context.Context, tenantID uuid.UUID) ([]models.Event, error)
	GetAllEventsPaginated(ctx context.Context, tenantID uuid.UUID, params models.PaginationParams) (models.PaginatedResponse[models.Event], error)
	GetEventByID(ctx context.Context, eventID uuid.UUID, tenantID uuid.UUID) (models.Event, error)
	UpdateEvent(ctx context.Context, arg models.UpdateEventParams, tenantID uuid.UUID) (models.Event, error)
	CountAllEvents(ctx context.Context, tenantID uuid.UUID) (int64, error)

	// Transaction-aware methods
	WithTransaction(ctx context.Context, fn func(EventsRepository) error) error
	CreateEventTx(ctx context.Context, tx pgx.Tx, arg models.CreateEventParams, tenantID uuid.UUID) (models.Event, error)
	DeleteEventTx(ctx context.Context, tx pgx.Tx, eventID uuid.UUID, tenantID uuid.UUID) (int64, error)
	UpdateEventTx(ctx context.Context, tx pgx.Tx, arg models.UpdateEventParams, tenantID uuid.UUID) (models.Event, error)
	GetEventByIDTx(ctx context.Context, tx pgx.Tx, eventID uuid.UUID, tenantID uuid.UUID) (models.Event, error)

	// Batch operations
	CreateEventsBatch(ctx context.Context, args []models.CreateEventParams, tenantID uuid.UUID) ([]models.Event, error)
	UpdateEventsBatch(ctx context.Context, args []models.UpdateEventParams, tenantID uuid.UUID) ([]models.Event, error)
}

// ActionsRepository defines the interface for actions-related database operations
type ActionsRepository interface {
	CreateAction(ctx context.Context, arg models.CreateActionParams) (models.Action, error)
	DeleteAction(ctx context.Context, id uuid.UUID) (int64, error)
	GetAllActionsForTenant(ctx context.Context, tenantID uuid.UUID) ([]models.Action, error)
	GetActionByID(ctx context.Context, id uuid.UUID) (models.Action, error)
}

// Database abstracts the database connection pool
type Database interface {
	Ping(ctx context.Context) error
}
