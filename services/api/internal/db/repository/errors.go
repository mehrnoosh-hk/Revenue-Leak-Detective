package repository

import (
	"errors"
)

// Tenant scope errors
var (
	ErrFailedToAcquireConnection = errors.New("failed to acquire connection")
	ErrSettingTenantID           = errors.New("failed to set tenant ID")
	ErrFailedToSetServiceAccount = errors.New("failed to set service account")
)

// Database repository errors
var (
	ErrDatabaseNotInitialized = errors.New("database not initialized")
	ErrDatabaseUnavailable    = errors.New("database unavailable")
	ErrDatabaseConnection     = errors.New("database connection")
	ErrForeignKeyViolation    = errors.New("foreign key violation")
	ErrNotNullViolation       = errors.New("not null violation")
	ErrCheckViolation         = errors.New("check violation")
)

// Repository construction errors
var (
	ErrPoolCannotBeNil   = errors.New("pool cannot be nil")
	ErrLoggerCannotBeNil = errors.New("logger cannot be nil")
)

// Events repository errors
var (
	ErrConvertingDataToJSONb = errors.New("error converting data to jsonb")
	ErrEventNotFound         = errors.New("event not found")
	ErrEventAlreadyExists    = errors.New("event already exists")
	ErrInvalidEventData      = errors.New("invalid event data")
	ErrEventUpdateFailed     = errors.New("event update failed")
	ErrEventDeleteFailed     = errors.New("event delete failed")
	ErrEventCreationFailed   = errors.New("event creation failed")
	ErrEventRetrievalFailed  = errors.New("event retrieval failed")
)

// Actions repository errors
var (
	ErrActionNotFound            = errors.New("action not found")
	ErrActionAlreadyExists       = errors.New("action already exists")
	ErrActionForeignKeyViolation = errors.New("action foreign key violation")
	ErrActionNotNullViolation    = errors.New("action not null violation")
	ErrDatabaseOperation         = errors.New("database operation")
)

// Users repository errors
var (
	ErrFailedToCreateUser     = errors.New("failed to create user")
	ErrFailedToDeleteUser     = errors.New("failed to delete user")
	ErrFailedToGetAllUsers    = errors.New("failed to get all users")
	ErrFailedToGetUserByEmail = errors.New("failed to get user by email")
	ErrFailedToGetUserByID    = errors.New("failed to get user by ID")
	ErrFailedToUpdateUser     = errors.New("failed to update user")
	ErrUserNotFound           = errors.New("user not found")
)
