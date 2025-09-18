// Package services provides the business logic implementations for domain entities.
// This file implements the UserService, which handles user-related operations.
package services

import (
	"context"
	"log/slog"
	"rdl-api/internal/db/repository"
	"rdl-api/internal/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersService interface {
	CreateUser(ctx context.Context, params models.CreateUserParams, tenantID uuid.UUID) (models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (int64, error)
	GetAllUsers(ctx context.Context, tenantID uuid.UUID) ([]models.User, error)
	GetUserByEmail(ctx context.Context, email string, tenantID uuid.UUID) (models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (models.User, error)
	UpdateUser(ctx context.Context, params models.UpdateUserParams, tenantID uuid.UUID) (models.User, error)
}

// userService implements the UserService interface using Tenant Aware User Repository
type userService struct {
	userRepository UserRepository
}

func NewUserService(pool *pgxpool.Pool, l *slog.Logger) UsersService {
	return &userService{
		userRepository: repository.NewUserRepository(pool, l),
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
