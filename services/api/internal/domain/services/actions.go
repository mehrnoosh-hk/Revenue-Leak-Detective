// Package services provides business logic implementations for domain entities.
// It acts as an abstraction layer between the application handlers and the repository layer,
// actions.go provides business logic for action-related operations.
package services

import (
	"context"
	"log/slog"
	"rdl-api/internal/db/repository"
	"rdl-api/internal/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ActionsService interface {
	CreateAction(ctx context.Context, args models.CreateActionParams, tenantID uuid.UUID) (models.Action, error)
	DeleteAction(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (int64, error)
	GetAllActions(ctx context.Context, tenantID uuid.UUID) ([]models.Action, error)
	GetAllActionsPaginated(ctx context.Context, tenantID uuid.UUID, params models.PaginationParams) (models.PaginatedResponse[models.Action], error)
	GetActionByID(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (models.Action, error)
	CountAllActions(ctx context.Context, tenantID uuid.UUID) (int64, error)
}

// NewActionsRepository creates a new instance of ActionsRepositoryImplementation backed by the provided pgx.Tx and pgxpool.Pool.
//
// Parameters:
//   - pool: Pointer to pgxpool.Pool, which provides access to the database.
//   - tx: Pointer to pgx.Tx, which provides access to the database transaction.
//   - logger: Pointer to slog.Logger, which provides access to the logger.
//
// Returns:
//   - ActionsRepository: An implementation of the ActionsRepository interface.
func NewActionsRepository(pool *pgxpool.Pool, tx pgx.Tx, logger *slog.Logger) repository.ActionsRepositoryImplementation {
	return repository.ActionsRepositoryImplementation{Pool: pool, Tx: tx, Logger: logger}
}

// actionsService implements the ActionsService interface using repository pattern.
type actionsService struct {
	actionsRepo repository.ActionsRepositoryImplementation
	logger      *slog.Logger
}

// NewActionsService creates a new instance of ActionsService backed by the provided pool.
//
// Parameters:
//   - pool: Database connection pool.
//   - logger: Logger for structured logging.
//
// Returns:
//   - ActionsService: An implementation of the ActionsService interface.
func NewActionsService(pool *pgxpool.Pool, logger *slog.Logger) ActionsService {
	// It needs to initialize an ActionsRepository with the dependencies injected from the app
	aR := NewActionsRepository(pool, nil, logger)
	return &actionsService{
		actionsRepo: aR,
		logger:      logger,
	}
}

// CreateAction creates a new action in the system.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - args: CreateActionParams containing the action details.
//   - tenantID: UUID of the tenant that owns the action.
//
// Returns:
//   - models.Action: The created action.
//   - error: Any error encountered during creation.
func (s *actionsService) CreateAction(ctx context.Context, args models.CreateActionParams, tenantID uuid.UUID) (models.Action, error) {
	s.logger.InfoContext(ctx, "Creating action", "leak_id", args.LeakID, "action_type", args.ActionType, "tenant_id", tenantID)

	action, err := s.actionsRepo.CreateAction(ctx, args, tenantID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create action", "error", err, "leak_id", args.LeakID, "tenant_id", tenantID)
		return models.Action{}, err
	}

	s.logger.InfoContext(ctx, "Action created successfully", "action_id", action.ID, "leak_id", args.LeakID, "tenant_id", tenantID)
	return action, nil
}

// DeleteAction removes an action from the system by its UUID.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - id: UUID of the action to delete.
//   - tenantID: UUID of the tenant that owns the action.
//
// Returns:
//   - int64: Number of rows affected (should be 1 if successful).
//   - error: Any error encountered during deletion.
func (s *actionsService) DeleteAction(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (int64, error) {
	s.logger.InfoContext(ctx, "Deleting action", "action_id", id, "tenant_id", tenantID)

	rowsAffected, err := s.actionsRepo.DeleteAction(ctx, id, tenantID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete action", "error", err, "action_id", id, "tenant_id", tenantID)
		return 0, err
	}

	s.logger.InfoContext(ctx, "Action deleted successfully", "action_id", id, "rows_affected", rowsAffected, "tenant_id", tenantID)
	return rowsAffected, nil
}

// GetAllActions retrieves all actions for a specific tenant.
// This method uses sensible defaults for pagination (limit 1000, offset 0) to maintain backward compatibility.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - tenantID: UUID of the tenant that owns the actions.
//
// Returns:
//   - []models.Action: Slice of actions for the tenant.
//   - error: Any error encountered during retrieval.
func (s *actionsService) GetAllActions(ctx context.Context, tenantID uuid.UUID) ([]models.Action, error) {
	s.logger.DebugContext(ctx, "Retrieving all actions for tenant", "tenant_id", tenantID)

	actions, err := s.actionsRepo.GetAllActions(ctx, tenantID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to retrieve actions for tenant", "error", err, "tenant_id", tenantID)
		return nil, err
	}

	s.logger.DebugContext(ctx, "Retrieved actions successfully", "tenant_id", tenantID, "count", len(actions))
	return actions, nil
}

// GetAllActionsPaginated retrieves actions for a specific tenant with pagination support.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - tenantID: UUID of the tenant that owns the actions.
//   - params: Pagination parameters (limit and offset).
//
// Returns:
//   - models.PaginatedResponse[models.Action]: Paginated response containing actions and metadata.
//   - error: Any error encountered during retrieval.
func (s *actionsService) GetAllActionsPaginated(ctx context.Context, tenantID uuid.UUID, params models.PaginationParams) (models.PaginatedResponse[models.Action], error) {
	s.logger.DebugContext(ctx, "Retrieving actions with pagination", "tenant_id", tenantID, "limit", params.Limit, "offset", params.Offset)

	response, err := s.actionsRepo.GetAllActionsPaginated(ctx, tenantID, params)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to retrieve paginated actions for tenant", "error", err, "tenant_id", tenantID)
		return models.PaginatedResponse[models.Action]{}, err
	}

	s.logger.DebugContext(ctx, "Retrieved paginated actions successfully", "tenant_id", tenantID, "count", len(response.Items), "total_count", response.TotalCount)
	return response, nil
}

// GetActionByID retrieves an action by its UUID.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - id: UUID of the action to retrieve.
//   - tenantID: UUID of the tenant that owns the action.
//
// Returns:
//   - models.Action: The requested action.
//   - error: Any error encountered during retrieval.
func (s *actionsService) GetActionByID(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (models.Action, error) {
	s.logger.DebugContext(ctx, "Retrieving action by ID", "action_id", id, "tenant_id", tenantID)

	action, err := s.actionsRepo.GetActionByID(ctx, id, tenantID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to retrieve action by ID", "error", err, "action_id", id, "tenant_id", tenantID)
		return models.Action{}, err
	}

	s.logger.DebugContext(ctx, "Retrieved action successfully", "action_id", id, "tenant_id", tenantID)
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
func (s *actionsService) CountAllActions(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	s.logger.DebugContext(ctx, "Counting actions for tenant", "tenant_id", tenantID)

	count, err := s.actionsRepo.CountAllActions(ctx, tenantID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to count actions for tenant", "error", err, "tenant_id", tenantID)
		return 0, err
	}

	s.logger.DebugContext(ctx, "Counted actions successfully", "tenant_id", tenantID, "count", count)
	return count, nil
}
