// Package repository provides implementations of data access patterns for domain entities.
// It acts as an abstraction layer between the application/business logic and the underlying database,
// enabling CRUD operations and encapsulating SQLC-generated query usage for event persistence and retrieval.
package repository

import (
	"context"
	db "rdl-api/internal/db/sqlc"
	"rdl-api/internal/domain/models"

	"github.com/google/uuid"
)

// eventsRepository implements the EventsRepository interface using sqlc-generated queries.
// It provides methods for CRUD operations on Event entities in the database.
type eventsRepository struct {
	queries *db.Queries
}

// NewEventsRepository creates a new instance of EventsRepository backed by the provided db.Queries.
//
// Parameters:
//   - queries: Pointer to db.Queries, which provides access to SQLC-generated query methods.
//
// Returns:
//   - EventsRepository: An implementation of the EventsRepository interface.
func NewEventsRepository(queries *db.Queries) EventsRepository {
	return &eventsRepository{queries: queries}
}

// CreateEvent persists a new event in the database.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - arg: CreateEventParams containing the event details as a domain model.
//
// Returns:
//   - models.Event: The created event as a domain model.
//   - error: Any error encountered during creation.
func (r *eventsRepository) CreateEvent(ctx context.Context, arg models.CreateEventParams) (models.Event, error) {
	params, err := toCreateEventDBParams(arg)
	if err != nil {
		return models.Event{}, err
	}

	event, err := r.queries.CreateEvent(ctx, params)
	if err != nil {
		return models.Event{}, err
	}
	return toEventDomain(event), nil
}

// DeleteEvent removes an event from the database by its UUID.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - id: UUID of the event to delete.
//
// Returns:
//   - int64: Number of rows affected (should be 1 if successful).
//   - error: Any error encountered during deletion.
func (r *eventsRepository) DeleteEvent(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.queries.DeleteEvent(ctx, db.ConvertUUIDToPgtypeUUID(id))
}

// GetAllEvents retrieves all events from the database.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//
// Returns:
//   - []models.Event: Slice of all event domain models.
//   - error: Any error encountered during retrieval.
func (r *eventsRepository) GetAllEvents(ctx context.Context) ([]models.Event, error) {
	events, err := r.queries.GetAllEvents(ctx)
	if err != nil {
		return nil, err
	}

	var result []models.Event
	for _, e := range events {
		result = append(result, toEventDomain(e))
	}
	return result, nil
}

// GetEventById retrieves a single event by its UUID.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - id: UUID of the event to retrieve.
//
// Returns:
//   - models.Event: The event domain model if found.
//   - error: Any error encountered during retrieval.
func (r *eventsRepository) GetEventById(ctx context.Context, id uuid.UUID) (models.Event, error) {
	event, err := r.queries.GetEventByID(ctx, db.ConvertUUIDToPgtypeUUID(id))
	if err != nil {
		return models.Event{}, err
	}
	return toEventDomain(event), nil
}

// UpdateEvent updates an existing event in the database.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - arg: UpdateEventParams containing the fields to update.
//
// Returns:
//   - models.Event: The updated event domain model.
//   - error: Any error encountered during update.
func (r *eventsRepository) UpdateEvent(ctx context.Context, arg models.UpdateEventParams) (models.Event, error) {
	params, err := toUpdateEventDBParams(arg)
	if err != nil {
		return models.Event{}, err
	}

	event, err := r.queries.UpdateEvent(ctx, params)
	if err != nil {
		return models.Event{}, err
	}
	return toEventDomain(event), nil
}

// toEventDomain converts a db.Event (database model) to a models.Event (domain model).
//
// Parameters:
//   - e: db.Event struct as returned by SQLC queries.
//
// Returns:
//   - models.Event: The corresponding domain model.
func toEventDomain(e db.Event) models.Event {
	return models.Event{
		ID:         uuid.UUID(e.ID.Bytes),
		TenantID:   uuid.UUID(e.TenantID.Bytes),
		ProviderID: uuid.UUID(e.ProviderID.Bytes),
		EventType:  models.EventTypeEnum(e.EventType),
		EventID:    e.EventID,
		Status:     models.EventStatusEnum(e.Status),
		Data:       any(e.Data),
		CreatedAt:  e.CreatedAt.Time,
		UpdatedAt:  e.UpdatedAt.Time,
	}
}

// toCreateEventDBParams converts a domain CreateEventParams to a db.CreateEventParams for persistence.
//
// Parameters:
//   - arg: models.CreateEventParams containing the event creation details.
//
// Returns:
//   - db.CreateEventParams: The database model for event creation.
//   - error: Any error encountered during conversion (e.g., data serialization).
func toCreateEventDBParams(arg models.CreateEventParams) (db.CreateEventParams, error) {
	data, err := db.ConvertInterfaceToBytes(arg.Data)
	if err != nil {
		return db.CreateEventParams{}, err
	}
	return db.CreateEventParams{
		TenantID:   db.ConvertUUIDToPgtypeUUID(arg.TenantID),
		ProviderID: db.ConvertUUIDToPgtypeUUID(arg.ProviderID),
		EventType:  db.EventTypeEnum(arg.EventType),
		EventID:    arg.EventID,
		Status:     db.EventStatusEnum(arg.Status),
		Data:       data,
	}, nil
}

// toUpdateEventDBParams converts a domain UpdateEventParams to a db.UpdateEventParams for persistence.
//
// Parameters:
//   - arg: models.UpdateEventParams containing the event update details.
//
// Returns:
//   - db.UpdateEventParams: The database model for event update.
//   - error: Any error encountered during conversion (e.g., data serialization).
func toUpdateEventDBParams(arg models.UpdateEventParams) (db.UpdateEventParams, error) {
	var data []byte
	var err error

	if arg.Data != nil {
		data, err = db.ConvertInterfaceToBytes(arg.Data)
		if err != nil {
			return db.UpdateEventParams{}, err
		}
	}

	resultEventType := db.ConvertEnumsToNullableEnum[*db.EventTypeEnum, db.NullEventTypeEnum]((*db.EventTypeEnum)(arg.EventType))

	resultEventStatus := db.ConvertEnumsToNullableEnum[*db.EventStatusEnum, db.NullEventStatusEnum]((*db.EventStatusEnum)(arg.Status))

	return db.UpdateEventParams{
		ID:         db.ConvertUUIDToPgtypeUUID(arg.ID),
		TenantID:   db.ConvertNullableUUIDToPgtypeUUID(arg.TenantID),
		ProviderID: db.ConvertNullableUUIDToPgtypeUUID(arg.ProviderID),
		EventType:  resultEventType,
		EventID:    arg.EventID,
		Status:     resultEventStatus,
		Data:       data,
	}, nil
}
