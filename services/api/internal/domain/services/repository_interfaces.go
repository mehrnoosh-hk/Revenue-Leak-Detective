package services

import (
	"context"

	"github.com/google/uuid"

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
	CheckReadiness(ctx context.Context) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, arg models.CreateUserParams, tenantID uuid.UUID) (models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (int64, error)
	GetAllUsers(ctx context.Context, tenantID uuid.UUID) ([]models.User, error)
	GetUserByEmail(ctx context.Context, email string, tenantID uuid.UUID) (models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (models.User, error)
	UpdateUser(ctx context.Context, arg models.UpdateUserParams, tenantID uuid.UUID) (models.User, error)
}

// EventsRepository defines the interface for events CRUD operations
type EventsRepository interface {
	// Create operations
	CreateEvent(ctx context.Context, arg models.CreateEventParams, tenantID uuid.UUID) (models.Event, error)

	// Read operations
	GetAllEvents(ctx context.Context, tenantID uuid.UUID) ([]models.Event, error)
	GetAllEventsPaginated(ctx context.Context, tenantID uuid.UUID, params models.PaginationParams) (models.PaginatedResponse[models.Event], error)
	GetEventByID(ctx context.Context, eventID uuid.UUID, tenantID uuid.UUID) (models.Event, error)
	CountAllEvents(ctx context.Context, tenantID uuid.UUID) (int64, error)

	// Update operations
	UpdateEvent(ctx context.Context, arg models.UpdateEventParams, tenantID uuid.UUID) (models.Event, error)

	// Delete operations
	DeleteEvent(ctx context.Context, eventID uuid.UUID, tenantID uuid.UUID) (int64, error)
}

// ActionsRepository defines the interface for actions-related database operations
type ActionsRepository interface {
	CreateAction(ctx context.Context, arg models.CreateActionParams, tenantID uuid.UUID) (models.Action, error)
	DeleteAction(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (int64, error)
	GetAllActions(ctx context.Context, tenantID uuid.UUID) ([]models.Action, error)
	GetAllActionsPaginated(ctx context.Context, tenantID uuid.UUID, params models.PaginationParams) (models.PaginatedResponse[models.Action], error)
	GetActionByID(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (models.Action, error)
	CountAllActions(ctx context.Context, tenantID uuid.UUID) (int64, error)
}

// Database abstracts the database connection pool
type Database interface {
	Ping(ctx context.Context) error
}
