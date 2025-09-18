// Package repository provides implementations of data access patterns for domain entities.
// It acts as an abstraction layer between the application/business logic and the underlying database,
// events.go provides CRUD operations and encapsulating SQLC-generated query usage for event persistence and retrieval.
package repository

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	db "rdl-api/internal/db/sqlc"
	"rdl-api/internal/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EventsRepositoryImplementation implements the EventsRepository interface using sqlc-generated queries.
type EventsRepositoryImplementation struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewEventsRepository creates a new instance of EventsRepository backed by the provided pgx.Tx and pgxpool.Pool.
//
// Parameters:
//   - pool: Pointer to pgxpool.Pool, which provides access to the database.
//   - logger: Pointer to slog.Logger, which provides access to the logger.
//
// Returns:
//   - EventsRepository: An implementation of the EventsRepository interface.
//   - error: Any error encountered during initialization.
func NewEventsRepository(pool *pgxpool.Pool, l *slog.Logger) (EventsRepositoryImplementation, error) {
	if pool == nil {
		return EventsRepositoryImplementation{}, ErrPoolCannotBeNil
	}
	if l == nil {
		return EventsRepositoryImplementation{}, ErrLoggerCannotBeNil
	}
	return EventsRepositoryImplementation{pool: pool, logger: l}, nil
}

// CreateEvent persists a new event in the database.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - arg: CreateEventParams containing the event details as a domain model.
//   - tenantID: UUID of the tenant that owns the event.
//
// Returns:
//   - models.Event: The created event as a domain model.
//   - error: Any error encountered during creation.
func (r EventsRepositoryImplementation) CreateEvent(ctx context.Context, arg models.CreateEventParams, tenantID uuid.UUID) (models.Event, error) {
	r.logger.InfoContext(ctx, "Creating event", "event_id", arg.EventID, "tenant_id", tenantID, "event_type", arg.EventType)

	var event models.Event
	err := WithTenantContext(ctx, r.pool, tenantID, func(queries *db.Queries) error {
		params, err := toCreateEventDBParams(arg)
		if err != nil {
			r.logger.ErrorContext(ctx, "Failed to convert event params", "error", err, "event_id", arg.EventID, "tenant_id", tenantID)
			return ErrConvertingDataToJSONb
		}

		dbEvent, err := queries.CreateEvent(ctx, params)
		if err != nil {
			return r.handleDatabaseError(ctx, err, "create event", arg.EventID, tenantID.String())
		}

		event = toEventDomain(dbEvent)
		r.logger.InfoContext(ctx, "Event created successfully", "event_id", event.ID, "tenant_id", tenantID)
		return nil
	})

	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to create event", "error", err, "event_id", arg.EventID, "tenant_id", tenantID)
		return models.Event{}, err
	}

	return event, nil
}

// DeleteEvent removes an event from the database by its UUID.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - eventID: UUID of the event to delete.
//   - tenantID: UUID of the tenant that owns the event.
//
// Returns:
//   - int64: Number of rows affected (should be 1 if successful).
//   - error: Any error encountered during deletion.
func (r EventsRepositoryImplementation) DeleteEvent(ctx context.Context, eventID uuid.UUID, tenantID uuid.UUID) (int64, error) {
	r.logger.InfoContext(ctx, "Deleting event", "event_id", eventID, "tenant_id", tenantID)

	var rowsAffected int64
	err := WithTenantContext(ctx, r.pool, tenantID, func(queries *db.Queries) error {
		rows, err := queries.DeleteEvent(ctx, convertUUIDToPgtypeUUID(eventID))
		if err != nil {
			return r.handleDatabaseError(ctx, err, "delete event", eventID.String(), tenantID.String())
		}

		if rows == 0 {
			r.logger.WarnContext(ctx, "Event not found for deletion", "event_id", eventID, "tenant_id", tenantID)
			return ErrEventNotFound
		}

		rowsAffected = rows
		r.logger.InfoContext(ctx, "Event deleted successfully", "event_id", eventID, "tenant_id", tenantID, "rows_affected", rowsAffected)
		return nil
	})

	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to delete event", "error", err, "event_id", eventID, "tenant_id", tenantID)
		return 0, err
	}

	return rowsAffected, nil
}

// GetAllEvents retrieves all events from the database.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - tenantID: UUID of the tenant that owns the events.
//
// Returns:
//   - []models.Event: Slice of all event domain models.
//   - error: Any error encountered during retrieval.
func (r EventsRepositoryImplementation) GetAllEvents(ctx context.Context, tenantID uuid.UUID) ([]models.Event, error) {
	r.logger.DebugContext(ctx, "Retrieving all events", "tenant_id", tenantID)

	var events []models.Event
	err := WithTenantContext(ctx, r.pool, tenantID, func(queries *db.Queries) error {
		// Use sensible defaults for pagination (limit 1000, offset 0)
		dbEvents, err := queries.GetAllEvents(ctx, db.GetAllEventsParams{
			Limit:  1000,
			Offset: 0,
		})
		if err != nil {
			return r.handleDatabaseError(ctx, err, "get all events", "", tenantID.String())
		}

		// Pre-allocate slice with known capacity for better performance
		events = make([]models.Event, 0, len(dbEvents))
		for _, dbEvent := range dbEvents {
			events = append(events, toEventDomain(dbEvent))
		}

		r.logger.DebugContext(ctx, "Retrieved events successfully", "tenant_id", tenantID, "count", len(events))
		return nil
	})

	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to retrieve events", "error", err, "tenant_id", tenantID)
		return nil, err
	}

	return events, nil
}

// GetAllEventsPaginated retrieves events from the database with pagination support.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - tenantID: UUID of the tenant that owns the events.
//   - params: Pagination parameters (limit and offset).
//
// Returns:
//   - models.PaginatedResponse[models.Event]: Paginated response containing events and metadata.
//   - error: Any error encountered during retrieval.
func (r EventsRepositoryImplementation) GetAllEventsPaginated(ctx context.Context, tenantID uuid.UUID, params models.PaginationParams) (models.PaginatedResponse[models.Event], error) {
	r.logger.DebugContext(ctx, "Retrieving events with pagination", "tenant_id", tenantID, "limit", params.Limit, "offset", params.Offset)

	var events []models.Event
	var totalCount int64

	err := WithTenantContext(ctx, r.pool, tenantID, func(queries *db.Queries) error {
		// Get total count
		count, err := queries.CountAllEvents(ctx)
		if err != nil {
			return r.handleDatabaseError(ctx, err, "count events", "", tenantID.String())
		}
		totalCount = count

		// Get paginated events
		dbEvents, err := queries.GetAllEventsPaginated(ctx, db.GetAllEventsPaginatedParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
		if err != nil {
			return r.handleDatabaseError(ctx, err, "get paginated events", "", tenantID.String())
		}

		// Convert to domain models
		events = make([]models.Event, 0, len(dbEvents))
		for _, dbEvent := range dbEvents {
			events = append(events, toEventDomain(dbEvent))
		}

		r.logger.DebugContext(ctx, "Retrieved paginated events successfully", "tenant_id", tenantID, "count", len(events), "total_count", totalCount)
		return nil
	})

	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to retrieve paginated events", "error", err, "tenant_id", tenantID)
		return models.PaginatedResponse[models.Event]{}, err
	}

	// Create paginated response
	response := models.NewPaginatedResponse(events, totalCount, params.Limit, params.Offset)
	return response, nil
}

// GetEventByID retrieves a single event by its UUID.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - eventID: UUID of the event to retrieve.
//   - tenantID: UUID of the tenant that owns the event.
//
// Returns:
//   - models.Event: The event domain model if found.
//   - error: Any error encountered during retrieval.
func (r EventsRepositoryImplementation) GetEventByID(ctx context.Context, eventID uuid.UUID, tenantID uuid.UUID) (models.Event, error) {
	r.logger.DebugContext(ctx, "Retrieving event by ID", "event_id", eventID, "tenant_id", tenantID)

	var event models.Event
	err := WithTenantContext(ctx, r.pool, tenantID, func(queries *db.Queries) error {
		dbEvent, err := queries.GetEventByID(ctx, convertUUIDToPgtypeUUID(eventID))
		if err != nil {
			// Check if it's a "no rows" error
			if errors.Is(err, pgx.ErrNoRows) {
				r.logger.WarnContext(ctx, "Event not found", "event_id", eventID, "tenant_id", tenantID)
				return ErrEventNotFound
			}

			return r.handleDatabaseError(ctx, err, "get event by ID", eventID.String(), tenantID.String())
		}

		event = toEventDomain(dbEvent)
		r.logger.DebugContext(ctx, "Event retrieved successfully", "event_id", eventID, "tenant_id", tenantID)
		return nil
	})

	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to retrieve event", "error", err, "event_id", eventID, "tenant_id", tenantID)
		return models.Event{}, err
	}

	return event, nil
}

// UpdateEvent updates an existing event in the database.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - arg: UpdateEventParams containing the fields to update.
//   - tenantID: UUID of the tenant that owns the event.
//
// Returns:
//   - models.Event: The updated event domain model.
//   - error: Any error encountered during update.
func (r EventsRepositoryImplementation) UpdateEvent(ctx context.Context, arg models.UpdateEventParams, tenantID uuid.UUID) (models.Event, error) {
	r.logger.InfoContext(ctx, "Updating event", "event_id", arg.ID, "tenant_id", tenantID)

	params, err := toUpdateEventDBParams(arg)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to convert update params", "error", err, "event_id", arg.ID, "tenant_id", tenantID)
		return models.Event{}, ErrConvertingDataToJSONb
	}

	var domainEvent models.Event
	err = WithTenantContext(ctx, r.pool, tenantID, func(queries *db.Queries) error {
		dbEvent, dbErr := queries.UpdateEvent(ctx, params)
		if dbErr != nil {
			// Check if it's a "no rows" error
			if errors.Is(dbErr, pgx.ErrNoRows) {
				r.logger.WarnContext(ctx, "Event not found for update", "event_id", arg.ID, "tenant_id", tenantID)
				return ErrEventNotFound
			}

			return r.handleDatabaseError(ctx, dbErr, "update event", arg.ID.String(), tenantID.String())
		}

		domainEvent = toEventDomain(dbEvent)
		r.logger.InfoContext(ctx, "Event updated successfully", "event_id", arg.ID, "tenant_id", tenantID)
		return nil
	})

	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to update event", "error", err, "event_id", arg.ID, "tenant_id", tenantID)
		return models.Event{}, err
	}

	return domainEvent, nil
}

// CountAllEvents counts all events in the database.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - tenantID: UUID of the tenant that owns the events.
//
// Returns:
//   - int64: Number of events.
//   - error: Any error encountered during counting.
func (r EventsRepositoryImplementation) CountAllEvents(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	r.logger.DebugContext(ctx, "Counting all events", "tenant_id", tenantID)

	var count int64
	err := WithTenantContext(ctx, r.pool, tenantID, func(queries *db.Queries) error {
		c, err := queries.CountAllEvents(ctx)
		if err != nil {
			return r.handleDatabaseError(ctx, err, "count events", "", tenantID.String())
		}
		count = c
		return nil
	})

	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to count all events", "error", err, "tenant_id", tenantID)
		return 0, err
	}

	return count, nil
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
		Data:       (*json.RawMessage)(&e.Data),
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
	data, err := convertInterfaceToBytes(arg.Data)
	if err != nil {
		return db.CreateEventParams{}, err
	}
	return db.CreateEventParams{
		TenantID:   convertUUIDToPgtypeUUID(arg.TenantID),
		ProviderID: convertUUIDToPgtypeUUID(arg.ProviderID),
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
		data, err = convertInterfaceToBytes(arg.Data)
		if err != nil {
			return db.UpdateEventParams{}, err
		}
	}

	resultEventType, err := convertEnumsToNullableEnum[*db.EventTypeEnum, db.NullEventTypeEnum]((*db.EventTypeEnum)(arg.EventType))
	if err != nil {
		return db.UpdateEventParams{}, err
	}

	resultEventStatus, err := convertEnumsToNullableEnum[*db.EventStatusEnum, db.NullEventStatusEnum]((*db.EventStatusEnum)(arg.Status))
	if err != nil {
		return db.UpdateEventParams{}, err
	}

	return db.UpdateEventParams{
		ID:        convertUUIDToPgtypeUUID(arg.ID),
		EventType: resultEventType,
		Status:    resultEventStatus,
		Data:      data,
	}, nil
}

// handleDatabaseError processes database-specific errors and returns appropriate wrapped errors.
// This method handles common PostgreSQL errors and converts them to domain-specific errors.
// If the error doesn't need special handling, it returns the original error wrapped with context.
//
// Parameters:
//   - ctx: The context for request/tracing metadata
//   - err: The original database error
//   - operation: The operation being performed (for context)
//   - eventID: The event ID (for context)
//   - tenantID: The tenant ID (for context)
//
// Returns:
//   - error: A wrapped error with appropriate context, or the original error wrapped with context
func (r EventsRepositoryImplementation) handleDatabaseError(ctx context.Context, err error, operation, eventID, tenantID string) error {
	if err == nil {
		return nil
	}

	// Handle PostgreSQL-specific errors
	var pgErr *pgconn.PgError
	var message string
	var errToReturn error

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			message = "Unique constraint violation"
			errToReturn = ErrEventAlreadyExists
		case "23503": // foreign_key_violation
			message = "Foreign key constraint violation"
			errToReturn = ErrForeignKeyViolation

		case "23502": // not_null_violation
			message = "Not null constraint violation"
			errToReturn = ErrNotNullViolation

		case "23514": // check_violation
			message = "Check constraint violation"
			errToReturn = ErrCheckViolation

		case "42P01": // undefined_table
			message = "Database table not found"
			errToReturn = ErrDatabaseUnavailable

		case "08006": // connection_failure
			message = "Database connection failure"
			errToReturn = ErrDatabaseConnection

		default:
			// Log unknown PostgreSQL errors and wrap with context
			message = "Unknown PostgreSQL error"
			errToReturn = err
		}

		r.logger.ErrorContext(ctx, message, "operation", operation, "event_id", eventID, "tenant_id", tenantID, "pg_code", pgErr.Code, "pg_error", pgErr.Message)
		return errToReturn
	}

	// Handle context cancellation
	if errors.Is(err, context.Canceled) {
		r.logger.WarnContext(ctx, "Operation canceled", "operation", operation, "event_id", eventID, "tenant_id", tenantID)
		return errors.New("operation canceled")
	}

	// Handle context timeout
	if errors.Is(err, context.DeadlineExceeded) {
		r.logger.WarnContext(ctx, "Operation timeout", "operation", operation, "event_id", eventID, "tenant_id", tenantID)
		return errors.New("operation timeout")
	}

	// For other errors, return the original error
	return err
}

func handleDatabaseErrorLogHelper(ctx context.Context, logger *slog.Logger, err error, operation, eventID, tenantID string) error {
	if err == nil {
		return nil
	}

	logger.ErrorContext(ctx, "Database error", "operation", operation, "event_id", eventID, "tenant_id", tenantID, "error", err)
	return err
}
