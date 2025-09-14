// Package repository provides implementations of data access patterns for domain entities.
// It acts as an abstraction layer between the application/business logic and the underlying database,
// users.go provides CRUD operations and encapsulating SQLC-generated query usage for users persistence and retrieval.
package repository

import (
	"context"
	"log/slog"

	db "rdl-api/internal/db/sqlc"
	models "rdl-api/internal/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewUserRepository(pool *pgxpool.Pool, l *slog.Logger) UserRepository {
	return &userRepository{
		db:     pool,
		logger: l,
	}
}

// CreateUser creates a new user within tenant context
func (r *userRepository) CreateUser(ctx context.Context, arg models.CreateUserParams, tenantID uuid.UUID) (models.User, error) {
	r.logger.InfoContext(ctx, "Creating user",
		"email", arg.Email,
		"tenant_id", tenantID)

	var user models.User
	err := WithTenantContext(ctx, r.db, tenantID, func(queries *db.Queries) error {
		dbUser, err := queries.CreateUser(ctx, toCreateUserDBParams(arg))
		if err != nil {
			r.logger.ErrorContext(ctx, ErrFailedToCreateUser.Error(),
				"email", arg.Email,
				"tenant_id", tenantID,
				"error", err)
			return err
		}
		user = toUserDomain(dbUser)
		return nil
	})

	if err != nil {
		r.logger.ErrorContext(ctx, ErrFailedToCreateUser.Error(),
			"email", arg.Email,
			"tenant_id", tenantID,
			"error", err)
		return user, err
	}

	r.logger.InfoContext(ctx, "User created successfully",
		"user_id", user.ID,
		"email", user.Email,
		"tenant_id", tenantID)
	return user, err
}

func (r *userRepository) DeleteUser(ctx context.Context, userID uuid.UUID, tenantID uuid.UUID) (int64, error) {
	r.logger.InfoContext(ctx, "Deleting user",
		"user_id", userID,
		"tenant_id", tenantID)

	var rowsAffected int64
	err := WithTenantContext(ctx, r.db, tenantID, func(queries *db.Queries) error {
		rows, err := queries.DeleteUser(ctx, convertUUIDToPgtypeUUID(userID))
		if err != nil {
			r.logger.ErrorContext(ctx, ErrFailedToDeleteUser.Error(),
				"user_id", userID,
				"tenant_id", tenantID,
				"error", err)
			return err
		}
		rowsAffected = rows
		return nil
	})

	if err != nil {
		r.logger.ErrorContext(ctx, ErrFailedToDeleteUser.Error(),
			"user_id", userID,
			"tenant_id", tenantID,
			"error", err)
		return rowsAffected, err
	}

	r.logger.InfoContext(ctx, "User deleted successfully",
		"user_id", userID,
		"tenant_id", tenantID,
		"rows_affected", rowsAffected)
	return rowsAffected, err
}

func (r *userRepository) GetAllUsers(ctx context.Context, tenantID uuid.UUID) ([]models.User, error) {
	r.logger.DebugContext(ctx, "Retrieving all users",
		"tenant_id", tenantID)

	var domainUsers []models.User
	err := WithTenantContext(ctx, r.db, tenantID, func(queries *db.Queries) error {
		dbUsers, err := queries.GetAllUsers(ctx)
		if err != nil {
			r.logger.ErrorContext(ctx, ErrFailedToGetAllUsers.Error(),
				"tenant_id", tenantID,
				"error", err)
			return err
		}
		for _, dbUser := range dbUsers {
			domainUsers = append(domainUsers, toUserDomain(dbUser))
		}
		return nil
	})

	if err != nil {
		r.logger.ErrorContext(ctx, ErrFailedToGetAllUsers.Error(),
			"tenant_id", tenantID,
			"error", err)
		return domainUsers, err
	}

	r.logger.DebugContext(ctx, "Retrieved all users successfully",
		"tenant_id", tenantID,
		"user_count", len(domainUsers))
	return domainUsers, err
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string, tenantID uuid.UUID) (models.User, error) {
	r.logger.DebugContext(ctx, "Retrieving user by email",
		"email", email,
		"tenant_id", tenantID)

	var domainUser models.User
	err := WithTenantContext(ctx, r.db, tenantID, func(queries *db.Queries) error {
		dbUser, err := queries.GetUserByEmail(ctx, email)
		if err != nil {
			r.logger.ErrorContext(ctx, ErrFailedToGetUserByEmail.Error(),
				"email", email,
				"tenant_id", tenantID,
				"error", err)
			return err
		}
		domainUser = toUserDomain(dbUser)
		return nil
	})

	if err != nil {
		r.logger.ErrorContext(ctx, ErrFailedToGetUserByEmail.Error(),
			"email", email,
			"tenant_id", tenantID,
			"error", err)
		return domainUser, err
	}

	r.logger.DebugContext(ctx, "Retrieved user by email successfully",
		"user_id", domainUser.ID,
		"email", email,
		"tenant_id", tenantID)
	return domainUser, err
}

func (r *userRepository) GetUserByID(ctx context.Context, userID uuid.UUID, tenantID uuid.UUID) (models.User, error) {
	r.logger.DebugContext(ctx, "Retrieving user by ID",
		"user_id", userID,
		"tenant_id", tenantID)

	var user models.User
	err := WithTenantContext(ctx, r.db, tenantID, func(queries *db.Queries) error {
		dbUser, err := queries.GetUserByID(ctx, convertUUIDToPgtypeUUID(userID))
		if err != nil {
			r.logger.ErrorContext(ctx, ErrFailedToGetUserByID.Error(),
				"user_id", userID,
				"tenant_id", tenantID,
				"error", err)
			return err
		}

		user = toUserDomain(dbUser)
		return nil
	})

	if err != nil {
		r.logger.ErrorContext(ctx, ErrFailedToGetUserByID.Error(),
			"user_id", userID,
			"tenant_id", tenantID,
			"error", err)
		return user, err
	}

	r.logger.DebugContext(ctx, "Retrieved user by ID successfully",
		"user_id", userID,
		"tenant_id", tenantID)
	return user, err
}

func (r *userRepository) UpdateUser(ctx context.Context, arg models.UpdateUserParams, tenantID uuid.UUID) (models.User, error) {
	r.logger.InfoContext(ctx, "Updating user",
		"user_id", arg.ID,
		"tenant_id", tenantID)

	var user models.User
	err := WithTenantContext(ctx, r.db, tenantID, func(queries *db.Queries) error {
		dbUser, err := queries.UpdateUser(ctx, toUpdateUserDBParams(arg))
		if err != nil {
			r.logger.ErrorContext(ctx, ErrFailedToUpdateUser.Error(),
				"user_id", arg.ID,
				"tenant_id", tenantID,
				"error", err)
			return err
		}
		user = toUserDomain(dbUser)
		return nil
	})

	if err != nil {
		r.logger.ErrorContext(ctx, ErrFailedToUpdateUser.Error(),
			"user_id", arg.ID,
			"tenant_id", tenantID,
			"error", err)
		return user, err
	}

	r.logger.InfoContext(ctx, "User updated successfully",
		"user_id", user.ID,
		"tenant_id", tenantID)
	return user, err
}

// toUserDomain converts a db.User (database model) to a models.User (domain model)
func toUserDomain(u db.User) models.User {
	return models.User{
		ID:         convertPgtypeUUIDToUUID(u.ID),
		Email:      u.Email,
		Name:       u.Name,
		ExternalID: u.ExternalID,
		CreatedAt:  u.CreatedAt.Time,
		UpdatedAt:  u.UpdatedAt.Time,
		TenantID:   convertPgtypeUUIDToUUID(u.ID),
	}
}

// toCreateUserDBParams converts domain CreateUserParams to database CreateUserParams
func toCreateUserDBParams(arg models.CreateUserParams) db.CreateUserParams {
	return db.CreateUserParams{
		Email:      arg.Email,
		Name:       arg.Name,
		ExternalID: arg.ExternalID,
	}
}

// toUpdateUserDBParams converts domain UpdateUserParams to database UpdateUserParams
func toUpdateUserDBParams(arg models.UpdateUserParams) db.UpdateUserParams {
	return db.UpdateUserParams{
		ID:         convertUUIDToPgtypeUUID(arg.ID),
		Email:      arg.Email,
		Name:       arg.Name,
		ExternalID: arg.ExternalID,
	}
}
