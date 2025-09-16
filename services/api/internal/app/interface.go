package app

import (
	"context"
	"rdl-api/internal/domain/models"

	"github.com/google/uuid"
)

type Services struct {
	HealthService HealthService
	UsersService UsersService
	EventsService EventsService
	ActionsService ActionsService
}

type HealthService interface {
	CheckReadiness(ctx context.Context) error
	CheckLiveness(ctx context.Context) error
}

type UsersService interface {
	CreateUser(ctx context.Context, args models.CreateUserParams, tenantID uuid.UUID) (models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (int64, error)
	GetAllUsers(ctx context.Context, tenantID uuid.UUID) ([]models.User, error)
	GetUserByEmail(ctx context.Context, email string, tenantID uuid.UUID) (models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (models.User, error)
	UpdateUser(ctx context.Context, args models.UpdateUserParams, tenantID uuid.UUID) (models.User, error)
}

type EventsService interface {
	CreateEvent(ctx context.Context, args models.CreateEventParams, tenantID uuid.UUID) (models.Event, error)
	DeleteEvent(ctx context.Context, eventID uuid.UUID, tenantID uuid.UUID) (int64, error)
	GetAllEvents(ctx context.Context, tenantID uuid.UUID) ([]models.Event, error)
	GetAllEventsPaginated(ctx context.Context, tenantID uuid.UUID, params models.PaginationParams) (models.PaginatedResponse[models.Event], error)
	GetEventByID(ctx context.Context, eventID uuid.UUID, tenantID uuid.UUID) (models.Event, error)
	UpdateEvent(ctx context.Context, args models.UpdateEventParams, tenantID uuid.UUID) (models.Event, error)
	CountAllEvents(ctx context.Context, tenantID uuid.UUID) (int64, error)
}

type ActionsService interface {
	CreateAction(ctx context.Context, args models.CreateActionParams, tenantID uuid.UUID) (models.Action, error)
	DeleteAction(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (int64, error)
	GetAllActions(ctx context.Context, tenantID uuid.UUID) ([]models.Action, error)
	GetAllActionsPaginated(ctx context.Context, tenantID uuid.UUID, params models.PaginationParams) (models.PaginatedResponse[models.Action], error)
	GetActionByID(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (models.Action, error)
	CountAllActions(ctx context.Context, tenantID uuid.UUID) (int64, error)
}