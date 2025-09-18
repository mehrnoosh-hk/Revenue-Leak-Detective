// Package repository provides implementations of data access patterns for domain entities.
// It acts as an abstraction layer between the application/business logic and the underlying database,
// actions.go provides CRUD operations and encapsulating SQLC-generated query usage for action persistence and retrieval.
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

// ActionsRepositoryImplementation implements the ActionsRepository interface using sqlc-generated queries.
type ActionsRepositoryImplementation struct {
	Pool   *pgxpool.Pool
	Tx     pgx.Tx
	Logger *slog.Logger
}

// CreateAction persists a new action in the database.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - arg: CreateActionParams containing the action details as a domain model.
//
// Returns:
//   - models.Action: The created action as a domain model.
//   - error: Any error encountered during creation.
func (r *ActionsRepositoryImplementation) CreateAction(ctx context.Context, arg models.CreateActionParams, tenantID uuid.UUID) (models.Action, error) {
	r.Logger.InfoContext(ctx, "Creating action", "leak_id", arg.LeakID, "action_type", arg.ActionType, "tenant_id", tenantID)

	var action models.Action
	err := WithTenantContext(ctx, r.Pool, tenantID, func(queries *db.Queries) error {

		params := toCreateActionDBParams(arg)

		dbAction, err := queries.CreateAction(ctx, params)
		if err != nil {
			return r.handleDatabaseError(ctx, err, &arg.LeakID, &tenantID)
		}

		action = toActionDomain(dbAction)
		r.Logger.InfoContext(ctx, "Action created successfully", "action_id", action.ID, "leak_id", arg.LeakID, "tenant_id", tenantID)
		return nil
	})

	if err != nil {
		r.Logger.ErrorContext(ctx, "Failed to create action", "error", err, "leak_id", arg.LeakID, "tenant_id", tenantID)
		return models.Action{}, err
	}

	return action, nil
}

// DeleteAction removes an action from the database by its UUID.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - id: UUID of the action to delete.
//
// Returns:
//   - int64: Number of rows affected (should be 1 if successful).
//   - error: Any error encountered during deletion.
func (r *ActionsRepositoryImplementation) DeleteAction(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (int64, error) {
	r.Logger.InfoContext(ctx, "Deleting action", "action_id", id, "tenant_id", tenantID)

	var rowsAffected int64
	err := WithTenantContext(ctx, r.Pool, tenantID, func(queries *db.Queries) error {
		rows, err := queries.DeleteAction(ctx, convertUUIDToPgtypeUUID(id))
		if err != nil {
			return r.handleDatabaseError(ctx, err, &id, &tenantID)
		}

		if rows == 0 {
			return ErrActionNotFound
		}

		rowsAffected = rows
		r.Logger.InfoContext(ctx, "Action deleted successfully", "action_id", id, "rows_affected", rows, "tenant_id", tenantID)
		return nil
	})

	if err != nil {
		r.Logger.ErrorContext(ctx, "Failed to delete action", "error", err, "action_id", id, "tenant_id", tenantID)
		return 0, err
	}

	return rowsAffected, nil
}

// GetAllActionsForTenant retrieves all actions for a specific tenant from the database.
// This method uses sensible defaults for pagination (limit 1000, offset 0) to maintain backward compatibility.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - tenantID: UUID of the tenant that owns the actions.
//
// Returns:
//   - []models.Action: Slice of actions as domain models.
//   - error: Any error encountered during retrieval.
func (r *ActionsRepositoryImplementation) GetAllActions(ctx context.Context, tenantID uuid.UUID) ([]models.Action, error) {
	r.Logger.DebugContext(ctx, "Retrieving all actions for tenant", "tenant_id", tenantID)

	var actions []models.Action
	err := WithTenantContext(ctx, r.Pool, tenantID, func(queries *db.Queries) error {
		// Use sensible defaults for pagination (limit 1000, offset 0)
		dbActions, err := queries.GetAllActions(ctx)

		if err != nil {
			return r.handleDatabaseError(ctx, err, nil, &tenantID)
		}

		// Convert to domain models
		actions = make([]models.Action, 0, len(dbActions))
		for _, dbAction := range dbActions {
			actions = append(actions, toActionDomain(dbAction))
		}

		r.Logger.DebugContext(ctx, "Retrieved actions successfully", "tenant_id", tenantID, "count", len(actions))
		return nil
	})

	if err != nil {
		r.Logger.ErrorContext(ctx, "Failed to retrieve actions for tenant", "error", err, "tenant_id", tenantID)
		return nil, err
	}

	return actions, nil
}

// GetAllActionsForTenantPaginated retrieves actions from the database with pagination support.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - tenantID: UUID of the tenant that owns the actions.
//   - params: Pagination parameters (limit and offset).
//
// Returns:
//   - models.PaginatedResponse[models.Action]: Paginated response containing actions and metadata.
//   - error: Any error encountered during retrieval.
func (r *ActionsRepositoryImplementation) GetAllActionsPaginated(ctx context.Context, tenantID uuid.UUID, params models.PaginationParams) (models.PaginatedResponse[models.Action], error) {
	r.Logger.DebugContext(ctx, "Retrieving actions with pagination", "tenant_id", tenantID, "limit", params.Limit, "offset", params.Offset)

	var actions []models.Action
	var totalCount int64

	err := WithTenantContext(ctx, r.Pool, tenantID, func(queries *db.Queries) error {
		// Get total count
		count, err := queries.CountAllActions(ctx)
		if err != nil {
			return r.handleDatabaseError(ctx, err, nil, &tenantID)
		}
		totalCount = count

		// Get paginated actions
		dbActions, err := queries.GetAllActionsPaginated(ctx, db.GetAllActionsPaginatedParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
		if err != nil {
			return r.handleDatabaseError(ctx, err, nil, &tenantID)
		}

		// Convert to domain models
		actions = make([]models.Action, 0, len(dbActions))
		for _, dbAction := range dbActions {
			actions = append(actions, toActionDomain(dbAction))
		}

		r.Logger.DebugContext(ctx, "Retrieved paginated actions successfully", "tenant_id", tenantID, "count", len(actions), "total_count", totalCount)
		return nil
	})

	if err != nil {
		r.Logger.ErrorContext(ctx, "Failed to retrieve paginated actions for tenant", "error", err, "tenant_id", tenantID)
		return models.PaginatedResponse[models.Action]{}, err
	}

	return models.NewPaginatedResponse(actions, totalCount, params.Limit, params.Offset), nil
}

// GetActionByID retrieves an action from the database by its UUID.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - id: UUID of the action to retrieve.
//
// Returns:
//   - models.Action: The action as a domain model.
//   - error: Any error encountered during retrieval.
func (r *ActionsRepositoryImplementation) GetActionByID(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (models.Action, error) {
	r.Logger.DebugContext(ctx, "Retrieving action by ID", "action_id", id, "tenant_id", tenantID)

	var action models.Action
	err := WithTenantContext(ctx, r.Pool, tenantID, func(queries *db.Queries) error {
		dbAction, err := queries.GetActionByID(ctx, convertUUIDToPgtypeUUID(id))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				errWraper := ErrActionNotFound
				r.Logger.WarnContext(ctx, "Action not found", "action_id", id, "tenant_id", tenantID)
				return errWraper
			}
			return r.handleDatabaseError(ctx, err, &id, &tenantID)
		}

		action = toActionDomain(dbAction)
		r.Logger.DebugContext(ctx, "Retrieved action successfully", "action_id", id, "tenant_id", tenantID)
		return nil
	})

	if err != nil {
		r.Logger.ErrorContext(ctx, "Failed to retrieve action by ID", "error", err, "action_id", id, "tenant_id", tenantID)
		return models.Action{}, err
	}

	return action, nil
}

// CountAllActions counts the total number of actions for a specific tenant.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - tenantID: UUID of the tenant that owns the actions.
//
// Returns:
//   - int64: Total count of actions for the tenant.
//   - error: Any error encountered during counting.
func (r *ActionsRepositoryImplementation) CountAllActions(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	r.Logger.DebugContext(ctx, "Counting actions for tenant", "tenant_id", tenantID)

	var count int64
	err := WithTenantContext(ctx, r.Pool, tenantID, func(queries *db.Queries) error {
		totalCount, err := queries.CountAllActions(ctx)
		if err != nil {
			return r.handleDatabaseError(ctx, err, nil, &tenantID)
		}

		count = totalCount
		r.Logger.DebugContext(ctx, "Counted actions successfully", "tenant_id", tenantID, "count", count)
		return nil
	})

	if err != nil {
		r.Logger.ErrorContext(ctx, "Failed to count actions for tenant", "error", err, "tenant_id", tenantID)
		return 0, err
	}

	return count, nil
}

// handleDatabaseError handles database-specific errors and converts them to domain errors.
func (r *ActionsRepositoryImplementation) handleDatabaseError(ctx context.Context, err error, resourceID *uuid.UUID, tenantID *uuid.UUID) error{
	if err == nil {
		return nil
	}

	r.Logger.ErrorContext(ctx, "Database error", "resource_id", resourceID, "tenant_id", tenantID, "error", err)

	// Handle specific PostgreSQL errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return ErrActionAlreadyExists
		case "23503": // foreign_key_violation
			return ErrActionForeignKeyViolation
		case "23502": // not_null_violation
			return ErrActionNotNullViolation
		}
	}

	// Handle connection errors
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrActionNotFound
	}

	// Generic database error
	return ErrDatabaseOperation
}

// toCreateActionDBParams converts domain CreateActionParams to SQLC CreateActionParams.
func toCreateActionDBParams(arg models.CreateActionParams) db.CreateActionParams {
	return db.CreateActionParams{
		LeakID:     convertUUIDToPgtypeUUID(arg.LeakID),
		ActionType: convertActionTypeEnumToDB(arg.ActionType),
		Status:     convertActionStatusEnumToDB(arg.Status),
		Result:     convertActionResultEnumToDB(arg.Result),
	}
}

// toActionDomain converts SQLC Action to domain Action.
func toActionDomain(dbAction db.Action) models.Action {
	return models.Action{
		ID:         convertPgtypeUUIDToUUID(dbAction.ID),
		LeakID:     convertPgtypeUUIDToUUID(dbAction.LeakID),
		ActionType: convertActionTypeEnumFromDB(dbAction.ActionType),
		Status:     convertActionStatusEnumFromDB(dbAction.Status),
		Result:     convertActionResultEnumFromDB(dbAction.Result),
		CreatedAt:  dbAction.CreatedAt.Time,
		UpdatedAt:  dbAction.UpdatedAt.Time,
	}
}
