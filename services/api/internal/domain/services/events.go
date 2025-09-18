// Package services provides business logic and orchestration for domain entities.
// This file implements the EventService, which handles event-related operations.
package services

import (
	"context"
	"log/slog"
	"rdl-api/internal/db/repository"
	"rdl-api/internal/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventsService interface {
	CreateEvent(ctx context.Context, args models.CreateEventParams, tenantID uuid.UUID) (models.Event, error)
	DeleteEvent(ctx context.Context, eventID uuid.UUID, tenantID uuid.UUID) (int64, error)
	GetAllEvents(ctx context.Context, tenantID uuid.UUID) ([]models.Event, error)
	GetAllEventsPaginated(ctx context.Context, tenantID uuid.UUID, params models.PaginationParams) (models.PaginatedResponse[models.Event], error)
	GetEventByID(ctx context.Context, eventID uuid.UUID, tenantID uuid.UUID) (models.Event, error)
	UpdateEvent(ctx context.Context, args models.UpdateEventParams, tenantID uuid.UUID) (models.Event, error)
	CountAllEvents(ctx context.Context, tenantID uuid.UUID) (int64, error)
}

type eventsService struct {
	eventsRepository EventsRepository
	logger           *slog.Logger
}

// - Pointer to an initialized EventService.
func NewEventService(pool *pgxpool.Pool, l *slog.Logger) (EventsService, error) {
	// It needs to initialze an EventsRepository with the dependencies injected from the app
	eR, err := repository.NewEventsRepository(pool, l)
	if err != nil {
		return nil, err
	}
	return &eventsService{eventsRepository: eR, logger: l}, nil
}

// CreateEvent creates a new event in the system.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - args: CreateEventParams containing the details of the event to be created.
//   - tenantID: UUID of the tenant that owns the event.
//
// Returns:
//   - The created Event domain model.
//   - An error if the creation fails.
func (s *eventsService) CreateEvent(ctx context.Context, args models.CreateEventParams, tenantID uuid.UUID) (models.Event, error) {
	return s.eventsRepository.CreateEvent(ctx, args, tenantID)
}

func (s *eventsService) DeleteEvent(ctx context.Context, eventID uuid.UUID, tenantID uuid.UUID) (int64, error) {
	return s.eventsRepository.DeleteEvent(ctx, eventID, tenantID)
}

func (s *eventsService) GetAllEvents(ctx context.Context, tenantID uuid.UUID) ([]models.Event, error) {
	return s.eventsRepository.GetAllEvents(ctx, tenantID)
}

func (s *eventsService) GetAllEventsPaginated(ctx context.Context, tenantID uuid.UUID, params models.PaginationParams) (models.PaginatedResponse[models.Event], error) {
	return s.eventsRepository.GetAllEventsPaginated(ctx, tenantID, params)
}

func (s *eventsService) GetEventByID(ctx context.Context, eventID uuid.UUID, tenantID uuid.UUID) (models.Event, error) {
	return s.eventsRepository.GetEventByID(ctx, eventID, tenantID)
}

func (s *eventsService) UpdateEvent(ctx context.Context, args models.UpdateEventParams, tenantID uuid.UUID) (models.Event, error) {
	return s.eventsRepository.UpdateEvent(ctx, args, tenantID)
}

func (s *eventsService) CountAllEvents(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	return s.eventsRepository.CountAllEvents(ctx, tenantID)
}
