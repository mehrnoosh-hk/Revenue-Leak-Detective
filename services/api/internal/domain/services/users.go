// Package services provides the business logic implementations for domain entities.
// This file implements the UserService, which handles user-related operations.
package services

import (
	"context"
	"rdl-api/internal/db/repository"
	"rdl-api/internal/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// userService implements the UserService interface using Tenant Aware User Repository
type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(db *pgxpool.Pool) UserService {
	return &userService{
		userRepository: repository.NewUserRepository(db),
	}
}

func (u *userService) CreateUser(ctx context.Context, params models.CreateUserParams, tenantID uuid.UUID) (models.User, error) {
	return u.userRepository.CreateUser(ctx, params, tenantID)
}

func (u *userService) DeleteUser(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (int64, error) {
	return u.userRepository.DeleteUser(ctx, id, tenantID)
}

func (u *userService) GetAllUsers(ctx context.Context, tenantID uuid.UUID) ([]models.User, error) {
	return u.userRepository.GetAllUsers(ctx, tenantID)
}

func (u *userService) GetUserByEmail(ctx context.Context, email string, tenantID uuid.UUID) (models.User, error) {
	return u.userRepository.GetUserByEmail(ctx, email, tenantID)
}

func (u *userService) GetUserByID(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (models.User, error) {
	return u.userRepository.GetUserByID(ctx, id, tenantID)
}

func (u *userService) UpdateUser(ctx context.Context, params models.UpdateUserParams, tenantID uuid.UUID) (models.User, error) {
	return u.userRepository.UpdateUser(ctx, params, tenantID)
}
