// Package repository provides implementations of data access patterns for domain entities.
// It acts as an abstraction layer between the application/business logic and the underlying database,
// events.go provides CRUD operations and encapsulating SQLC-generated query usage for event persistence and retrieval.
package repository

import (
	"context"
	"errors"
	"log/slog"
	db "rdl-api/internal/db/sqlc"
	"rdl-api/internal/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// eventsRepository implements the EventsRepository interface using sqlc-generated queries.
type eventsRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewEventsRepository creates a new instance of EventsRepository backed by the provided db.Queries.
//
// Parameters:
//   - queries: Pointer to db.Queries, which provides access to SQLC-generated query methods.
//
// Returns:
//   - EventsRepository: An implementation of the EventsRepository interface.
func NewEventsRepository(pool *pgxpool.Pool, logger *slog.Logger) EventsRepository {
	return &eventsRepository{pool: pool, logger: logger}
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
func (r *eventsRepository) CreateEvent(ctx context.Context, arg models.CreateEventParams, tenantID uuid.UUID) (models.Event, error) {
	r.logger.InfoContext(ctx, "Creating event", "event_id", arg.EventID, "tenant_id", tenantID, "event_type", arg.EventType)

	var event models.Event
	err := WithTenantContext(ctx, r.pool, tenantID, func(queries *db.Queries) error {
		params, err := toCreateEventDBParams(arg)
		if err != nil {
			r.logger.ErrorContext(ctx, "Failed to convert event params", "error", err, "event_id", arg.EventID, "tenant_id", tenantID)
			return WrapError("event creation", ErrConvertingDataToJSONb, arg.EventID, tenantID.String())
		}

		dbEvent, err := queries.CreateEvent(ctx, params)
		if err != nil {
			return r.handleDatabaseError(err, "create event", arg.EventID, tenantID.String())
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
func (r *eventsRepository) DeleteEvent(ctx context.Context, eventID uuid.UUID, tenantID uuid.UUID) (int64, error) {
	r.logger.InfoContext(ctx, "Deleting event", "event_id", eventID, "tenant_id", tenantID)

	var rowsAffected int64
	err := WithTenantContext(ctx, r.pool, tenantID, func(queries *db.Queries) error {
		rows, err := queries.DeleteEvent(ctx, db.ConvertUUIDToPgtypeUUID(eventID))
		if err != nil {
			return r.handleDatabaseError(err, "delete event", eventID.String(), tenantID.String())
		}

		if rows == 0 {
			r.logger.WarnContext(ctx, "Event not found for deletion", "event_id", eventID, "tenant_id", tenantID)
			return WrapError("event deletion", ErrEventNotFound, eventID.String(), tenantID.String())
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
func (r *eventsRepository) GetAllEvents(ctx context.Context, tenantID uuid.UUID) ([]models.Event, error) {
	r.logger.DebugContext(ctx, "Retrieving all events", "tenant_id", tenantID)

	var events []models.Event
	err := WithTenantContext(ctx, r.pool, tenantID, func(queries *db.Queries) error {
		dbEvents, err := queries.GetAllEvents(ctx)
		if err != nil {
			return r.handleDatabaseError(err, "get all events", "", tenantID.String())
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
func (r *eventsRepository) GetAllEventsPaginated(ctx context.Context, tenantID uuid.UUID, params models.PaginationParams) (models.PaginatedResponse[models.Event], error) {
	r.logger.DebugContext(ctx, "Retrieving events with pagination", "tenant_id", tenantID, "limit", params.Limit, "offset", params.Offset)

	var events []models.Event
	var totalCount int64

	err := WithTenantContext(ctx, r.pool, tenantID, func(queries *db.Queries) error {
		// Get total count
		count, err := queries.CountAllEvents(ctx)
		if err != nil {
			return r.handleDatabaseError(err, "count events", "", tenantID.String())
		}
		totalCount = count

		// Get paginated events
		dbEvents, err := queries.GetAllEventsPaginated(ctx, db.GetAllEventsPaginatedParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
		if err != nil {
			return r.handleDatabaseError(err, "get paginated events", "", tenantID.String())
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
func (r *eventsRepository) GetEventByID(ctx context.Context, eventID uuid.UUID, tenantID uuid.UUID) (models.Event, error) {
	r.logger.DebugContext(ctx, "Retrieving event by ID", "event_id", eventID, "tenant_id", tenantID)

	var event models.Event
	err := WithTenantContext(ctx, r.pool, tenantID, func(queries *db.Queries) error {
		dbEvent, err := queries.GetEventByID(ctx, db.ConvertUUIDToPgtypeUUID(eventID))
		if err != nil {
			// Check if it's a "no rows" error
			if errors.Is(err, pgx.ErrNoRows) {
				r.logger.WarnContext(ctx, "Event not found", "event_id", eventID, "tenant_id", tenantID)
				return WrapError("event retrieval", ErrEventNotFound, eventID.String(), tenantID.String())
			}

			return r.handleDatabaseError(err, "get event by ID", eventID.String(), tenantID.String())
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
func (r *eventsRepository) UpdateEvent(ctx context.Context, arg models.UpdateEventParams, tenantID uuid.UUID) (models.Event, error) {
	r.logger.InfoContext(ctx, "Updating event", "event_id", arg.ID, "tenant_id", tenantID)

	params, err := toUpdateEventDBParams(arg)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to convert update params", "error", err, "event_id", arg.ID, "tenant_id", tenantID)
		return models.Event{}, WrapError("event update", ErrConvertingDataToJSONb, arg.ID.String(), tenantID.String())
	}

	var domainEvent models.Event
	err = WithTenantContext(ctx, r.pool, tenantID, func(queries *db.Queries) error {
		dbEvent, err := queries.UpdateEvent(ctx, params)
		if err != nil {
			// Check if it's a "no rows" error
			if errors.Is(err, pgx.ErrNoRows) {
				r.logger.WarnContext(ctx, "Event not found for update", "event_id", arg.ID, "tenant_id", tenantID)
				return WrapError("event update", ErrEventNotFound, arg.ID.String(), tenantID.String())
			}

			return r.handleDatabaseError(err, "update event", arg.ID.String(), tenantID.String())
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
func (r *eventsRepository) CountAllEvents(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	r.logger.DebugContext(ctx, "Counting all events", "tenant_id", tenantID)

	var count int64
	err := WithTenantContext(ctx, r.pool, tenantID, func(queries *db.Queries) error {
		c, err := queries.CountAllEvents(ctx)
		if err != nil {
			return r.handleDatabaseError(err, "count events", "", tenantID.String())
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

// WithTransaction executes a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// If the function completes successfully, the transaction is committed.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - fn: Function to execute within the transaction.
//
// Returns:
//   - error: Any error encountered during transaction execution.
func (r *eventsRepository) WithTransaction(ctx context.Context, fn func(EventsRepository) error) error {
	r.logger.DebugContext(ctx, "Starting database transaction")

	// Begin transaction
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to begin transaction", "error", err)
		return WrapError("transaction", ErrDatabaseConnection, "", "")
	}

	// Ensure transaction is properly handled
	defer func() {
		if p := recover(); p != nil {
			// Rollback on panic
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				r.logger.ErrorContext(ctx, "Failed to rollback transaction after panic", "error", rollbackErr, "panic", p)
			}
			panic(p) // Re-panic after rollback
		}
	}()

	// Create a new repository instance with the transaction
	txRepo := &eventsRepository{
		pool:   r.pool,
		logger: r.logger,
	}

	// Execute the function with the transaction-aware repository
	err = fn(txRepo)
	if err != nil {
		// Rollback on error
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			r.logger.ErrorContext(ctx, "Failed to rollback transaction", "error", rollbackErr, "original_error", err)
			return WrapError("transaction rollback", rollbackErr, "", "")
		}
		r.logger.WarnContext(ctx, "Transaction rolled back due to error", "error", err)
		return err
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		r.logger.ErrorContext(ctx, "Failed to commit transaction", "error", err)
		return WrapError("transaction commit", err, "", "")
	}

	r.logger.DebugContext(ctx, "Transaction committed successfully")
	return nil
}

// CreateEventTx creates a new event within a transaction.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - tx: Database transaction.
//   - arg: CreateEventParams containing the event details as a domain model.
//   - tenantID: UUID of the tenant that owns the event.
//
// Returns:
//   - models.Event: The created event as a domain model.
//   - error: Any error encountered during creation.
func (r *eventsRepository) CreateEventTx(ctx context.Context, tx pgx.Tx, arg models.CreateEventParams, tenantID uuid.UUID) (models.Event, error) {
	r.logger.InfoContext(ctx, "Creating event in transaction", "event_id", arg.EventID, "tenant_id", tenantID, "event_type", arg.EventType)

	params, err := toCreateEventDBParams(arg)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to convert event params", "error", err, "event_id", arg.EventID, "tenant_id", tenantID)
		return models.Event{}, WrapError("event creation", ErrConvertingDataToJSONb, arg.EventID, tenantID.String())
	}

	// Set tenant context in transaction
	_, err = tx.Exec(ctx, "SET app.current_tenant_id = $1", tenantID.String())
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to set tenant context in transaction", "error", err, "tenant_id", tenantID)
		return models.Event{}, WrapError("event creation", ErrFailedToSetTenantID, arg.EventID, tenantID.String())
	}

	// Set service account flag to false
	_, err = tx.Exec(ctx, "SET app.is_service_account = false")
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to set service account flag in transaction", "error", err, "tenant_id", tenantID)
		return models.Event{}, WrapError("event creation", ErrFailedToSetServiceAccount, arg.EventID, tenantID.String())
	}

	// Create queries instance with transaction
	queries := db.New(tx)
	dbEvent, err := queries.CreateEvent(ctx, params)
	if err != nil {
		return models.Event{}, r.handleDatabaseError(err, "create event", arg.EventID, tenantID.String())
	}

	event := toEventDomain(dbEvent)
	r.logger.InfoContext(ctx, "Event created successfully in transaction", "event_id", event.ID, "tenant_id", tenantID)
	return event, nil
}

// DeleteEventTx deletes an event within a transaction.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - tx: Database transaction.
//   - eventID: UUID of the event to delete.
//   - tenantID: UUID of the tenant that owns the event.
//
// Returns:
//   - int64: Number of rows affected (should be 1 if successful).
//   - error: Any error encountered during deletion.
func (r *eventsRepository) DeleteEventTx(ctx context.Context, tx pgx.Tx, eventID uuid.UUID, tenantID uuid.UUID) (int64, error) {
	r.logger.InfoContext(ctx, "Deleting event in transaction", "event_id", eventID, "tenant_id", tenantID)

	// Set tenant context in transaction
	_, err := tx.Exec(ctx, "SET app.current_tenant_id = $1", tenantID.String())
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to set tenant context in transaction", "error", err, "tenant_id", tenantID)
		return 0, WrapError("event deletion", ErrFailedToSetTenantID, eventID.String(), tenantID.String())
	}

	// Set service account flag to false
	_, err = tx.Exec(ctx, "SET app.is_service_account = false")
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to set service account flag in transaction", "error", err, "tenant_id", tenantID)
		return 0, WrapError("event deletion", ErrFailedToSetServiceAccount, eventID.String(), tenantID.String())
	}

	// Create queries instance with transaction
	queries := db.New(tx)
	rows, err := queries.DeleteEvent(ctx, db.ConvertUUIDToPgtypeUUID(eventID))
	if err != nil {
		return 0, r.handleDatabaseError(err, "delete event", eventID.String(), tenantID.String())
	}

	if rows == 0 {
		r.logger.WarnContext(ctx, "Event not found for deletion in transaction", "event_id", eventID, "tenant_id", tenantID)
		return 0, WrapError("event deletion", ErrEventNotFound, eventID.String(), tenantID.String())
	}

	r.logger.InfoContext(ctx, "Event deleted successfully in transaction", "event_id", eventID, "tenant_id", tenantID, "rows_affected", rows)
	return rows, nil
}

// UpdateEventTx updates an existing event within a transaction.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - tx: Database transaction.
//   - arg: UpdateEventParams containing the fields to update.
//   - tenantID: UUID of the tenant that owns the event.
//
// Returns:
//   - models.Event: The updated event as a domain model.
//   - error: Any error encountered during update.
func (r *eventsRepository) UpdateEventTx(ctx context.Context, tx pgx.Tx, arg models.UpdateEventParams, tenantID uuid.UUID) (models.Event, error) {
	r.logger.InfoContext(ctx, "Updating event in transaction", "event_id", arg.ID, "tenant_id", tenantID)

	params, err := toUpdateEventDBParams(arg)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to convert update params", "error", err, "event_id", arg.ID, "tenant_id", tenantID)
		return models.Event{}, WrapError("event update", ErrConvertingDataToJSONb, arg.ID.String(), tenantID.String())
	}

	// Set tenant context in transaction
	_, err = tx.Exec(ctx, "SET app.current_tenant_id = $1", tenantID.String())
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to set tenant context in transaction", "error", err, "tenant_id", tenantID)
		return models.Event{}, WrapError("event update", ErrFailedToSetTenantID, arg.ID.String(), tenantID.String())
	}

	// Set service account flag to false
	_, err = tx.Exec(ctx, "SET app.is_service_account = false")
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to set service account flag in transaction", "error", err, "tenant_id", tenantID)
		return models.Event{}, WrapError("event update", ErrFailedToSetServiceAccount, arg.ID.String(), tenantID.String())
	}

	// Create queries instance with transaction
	queries := db.New(tx)
	dbEvent, err := queries.UpdateEvent(ctx, params)
	if err != nil {
		// Check if it's a "no rows" error
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.WarnContext(ctx, "Event not found for update in transaction", "event_id", arg.ID, "tenant_id", tenantID)
			return models.Event{}, WrapError("event update", ErrEventNotFound, arg.ID.String(), tenantID.String())
		}
		return models.Event{}, r.handleDatabaseError(err, "update event", arg.ID.String(), tenantID.String())
	}

	domainEvent := toEventDomain(dbEvent)
	r.logger.InfoContext(ctx, "Event updated successfully in transaction", "event_id", arg.ID, "tenant_id", tenantID)
	return domainEvent, nil
}

// GetEventByIDTx retrieves a single event by its UUID within a transaction.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - tx: Database transaction.
//   - eventID: UUID of the event to retrieve.
//   - tenantID: UUID of the tenant that owns the event.
//
// Returns:
//   - models.Event: The event domain model if found.
//   - error: Any error encountered during retrieval.
func (r *eventsRepository) GetEventByIDTx(ctx context.Context, tx pgx.Tx, eventID uuid.UUID, tenantID uuid.UUID) (models.Event, error) {
	r.logger.DebugContext(ctx, "Retrieving event by ID in transaction", "event_id", eventID, "tenant_id", tenantID)

	// Set tenant context in transaction
	_, err := tx.Exec(ctx, "SET app.current_tenant_id = $1", tenantID.String())
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to set tenant context in transaction", "error", err, "tenant_id", tenantID)
		return models.Event{}, WrapError("event retrieval", ErrFailedToSetTenantID, eventID.String(), tenantID.String())
	}

	// Set service account flag to false
	_, err = tx.Exec(ctx, "SET app.is_service_account = false")
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to set service account flag in transaction", "error", err, "tenant_id", tenantID)
		return models.Event{}, WrapError("event retrieval", ErrFailedToSetServiceAccount, eventID.String(), tenantID.String())
	}

	// Create queries instance with transaction
	queries := db.New(tx)
	dbEvent, err := queries.GetEventByID(ctx, db.ConvertUUIDToPgtypeUUID(eventID))
	if err != nil {
		// Check if it's a "no rows" error
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.WarnContext(ctx, "Event not found in transaction", "event_id", eventID, "tenant_id", tenantID)
			return models.Event{}, WrapError("event retrieval", ErrEventNotFound, eventID.String(), tenantID.String())
		}
		return models.Event{}, r.handleDatabaseError(err, "get event by ID", eventID.String(), tenantID.String())
	}

	event := toEventDomain(dbEvent)
	r.logger.DebugContext(ctx, "Event retrieved successfully in transaction", "event_id", eventID, "tenant_id", tenantID)
	return event, nil
}

// CreateEventsBatch creates multiple events in a single transaction.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - args: Slice of CreateEventParams containing the event details.
//   - tenantID: UUID of the tenant that owns the events.
//
// Returns:
//   - []models.Event: Slice of created events as domain models.
//   - error: Any error encountered during batch creation.
func (r *eventsRepository) CreateEventsBatch(ctx context.Context, args []models.CreateEventParams, tenantID uuid.UUID) ([]models.Event, error) {
	if len(args) == 0 {
		if r.logger != nil {
			r.logger.WarnContext(ctx, "Empty batch provided for event creation", "tenant_id", tenantID)
		}
		return []models.Event{}, nil
	}

	if r.logger != nil {
		r.logger.InfoContext(ctx, "Creating events batch", "count", len(args), "tenant_id", tenantID)
	}

	var createdEvents []models.Event
	err := r.WithTransaction(ctx, func(txRepo EventsRepository) error {
		// Pre-allocate slice with known capacity
		createdEvents = make([]models.Event, 0, len(args))

		// Get the transaction from the repository (we need to access the actual transaction)
		// For now, we'll use the regular CreateEvent method within the transaction context
		for i, arg := range args {
			event, err := txRepo.CreateEvent(ctx, arg, tenantID)
			if err != nil {
				r.logger.ErrorContext(ctx, "Failed to create event in batch", "error", err, "index", i, "event_id", arg.EventID, "tenant_id", tenantID)
				return WrapError("batch event creation", err, arg.EventID, tenantID.String())
			}
			createdEvents = append(createdEvents, event)
		}

		r.logger.InfoContext(ctx, "Events batch created successfully", "count", len(createdEvents), "tenant_id", tenantID)
		return nil
	})

	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to create events batch", "error", err, "tenant_id", tenantID)
		return nil, err
	}

	return createdEvents, nil
}

// UpdateEventsBatch updates multiple events in a single transaction.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - args: Slice of UpdateEventParams containing the update details.
//   - tenantID: UUID of the tenant that owns the events.
//
// Returns:
//   - []models.Event: Slice of updated events as domain models.
//   - error: Any error encountered during batch update.
func (r *eventsRepository) UpdateEventsBatch(ctx context.Context, args []models.UpdateEventParams, tenantID uuid.UUID) ([]models.Event, error) {
	if len(args) == 0 {
		if r.logger != nil {
			r.logger.WarnContext(ctx, "Empty batch provided for event update", "tenant_id", tenantID)
		}
		return []models.Event{}, nil
	}

	if r.logger != nil {
		r.logger.InfoContext(ctx, "Updating events batch", "count", len(args), "tenant_id", tenantID)
	}

	var updatedEvents []models.Event
	err := r.WithTransaction(ctx, func(txRepo EventsRepository) error {
		// Pre-allocate slice with known capacity
		updatedEvents = make([]models.Event, 0, len(args))

		for i, arg := range args {
			event, err := txRepo.UpdateEvent(ctx, arg, tenantID)
			if err != nil {
				r.logger.ErrorContext(ctx, "Failed to update event in batch", "error", err, "index", i, "event_id", arg.ID, "tenant_id", tenantID)
				return WrapError("batch event update", err, arg.ID.String(), tenantID.String())
			}
			updatedEvents = append(updatedEvents, event)
		}

		r.logger.InfoContext(ctx, "Events batch updated successfully", "count", len(updatedEvents), "tenant_id", tenantID)
		return nil
	})

	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to update events batch", "error", err, "tenant_id", tenantID)
		return nil, err
	}

	return updatedEvents, nil
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

// handleDatabaseError processes database-specific errors and returns appropriate wrapped errors.
// This method handles common PostgreSQL errors and converts them to domain-specific errors.
// If the error doesn't need special handling, it returns the original error wrapped with context.
//
// Parameters:
//   - err: The original database error
//   - operation: The operation being performed (for context)
//   - eventID: The event ID (for context)
//   - tenantID: The tenant ID (for context)
//
// Returns:
//   - error: A wrapped error with appropriate context, or the original error wrapped with context
func (r *eventsRepository) handleDatabaseError(err error, operation, eventID, tenantID string) error {
	if err == nil {
		return nil
	}

	// Handle PostgreSQL-specific errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			r.logger.WarnContext(context.Background(), "Unique constraint violation", "operation", operation, "event_id", eventID, "tenant_id", tenantID, "pg_error", pgErr.Message)
			return WrapError(operation, ErrEventAlreadyExists, eventID, tenantID)

		case "23503": // foreign_key_violation
			r.logger.WarnContext(context.Background(), "Foreign key constraint violation", "operation", operation, "event_id", eventID, "tenant_id", tenantID, "pg_error", pgErr.Message)
			return WrapError(operation, ErrInvalidEventData, eventID, tenantID)

		case "23502": // not_null_violation
			r.logger.WarnContext(context.Background(), "Not null constraint violation", "operation", operation, "event_id", eventID, "tenant_id", tenantID, "pg_error", pgErr.Message)
			return WrapError(operation, ErrInvalidEventData, eventID, tenantID)

		case "23514": // check_violation
			r.logger.WarnContext(context.Background(), "Check constraint violation", "operation", operation, "event_id", eventID, "tenant_id", tenantID, "pg_error", pgErr.Message)
			return WrapError(operation, ErrInvalidEventData, eventID, tenantID)

		case "42P01": // undefined_table
			r.logger.ErrorContext(context.Background(), "Database table not found", "operation", operation, "event_id", eventID, "tenant_id", tenantID, "pg_error", pgErr.Message)
			return WrapError(operation, ErrDatabaseUnavailable, eventID, tenantID)

		case "08006": // connection_failure
			r.logger.ErrorContext(context.Background(), "Database connection failure", "operation", operation, "event_id", eventID, "tenant_id", tenantID, "pg_error", pgErr.Message)
			return WrapError(operation, ErrDatabaseConnection, eventID, tenantID)

		default:
			// Log unknown PostgreSQL errors and wrap with context
			r.logger.ErrorContext(context.Background(), "Unknown PostgreSQL error", "operation", operation, "event_id", eventID, "tenant_id", tenantID, "pg_code", pgErr.Code, "pg_error", pgErr.Message)
			return WrapError(operation, err, eventID, tenantID)
		}
	}

	// Handle context cancellation
	if errors.Is(err, context.Canceled) {
		r.logger.WarnContext(context.Background(), "Operation canceled", "operation", operation, "event_id", eventID, "tenant_id", tenantID)
		return WrapError(operation, errors.New("operation canceled"), eventID, tenantID)
	}

	// Handle context timeout
	if errors.Is(err, context.DeadlineExceeded) {
		r.logger.WarnContext(context.Background(), "Operation timeout", "operation", operation, "event_id", eventID, "tenant_id", tenantID)
		return WrapError(operation, errors.New("operation timeout"), eventID, tenantID)
	}

	// For other errors, wrap with context and return
	return WrapError(operation, err, eventID, tenantID)
}
