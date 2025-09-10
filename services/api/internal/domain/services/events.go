// Package services provides business logic and orchestration for domain entities.
// This file implements the EventService, which handles event-related operations.
package services

import (
	"context"
	"rdl-api/internal/db/repository"
	"rdl-api/internal/domain/models"
)

// EventService provides methods for managing events in the system.
// It acts as a bridge between the application/business logic and the data repository layer.
type EventService struct {
	// EventRepository is the data access layer for event persistence.
	EventRepository repository.EventsRepository
}

// NewEventService constructs a new EventService with the given EventsRepository.
//
// Parameters:
//   - eventRepository: An implementation of EventsRepository for event data access.
//
// Returns:
//   - Pointer to an initialized EventService.
func NewEventService(eventRepository repository.EventsRepository) *EventService {
	return &EventService{EventRepository: eventRepository}
}

// CreateEvent creates a new event in the system.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - params: CreateEventParams containing the details of the event to be created.
//
// Returns:
//   - The created Event domain model.
//   - An error if the creation fails.
func (s *EventService) CreateEvent(ctx context.Context, params models.CreateEventParams) (models.Event, error) {
	return s.EventRepository.CreateEvent(ctx, params)
}