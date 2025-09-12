package services

import (
	"context"
	"rdl-api/internal/domain/models"

	"github.com/google/uuid"
)

type DomainServices interface {
	UsersService
	EventsService
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
