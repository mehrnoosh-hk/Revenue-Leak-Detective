package repository

import (
	"errors"
	"fmt"
)

// Database repository errors
var (
	ErrDatabaseNotInitialized = errors.New("database not initialized")
	ErrDatabaseUnavailable    = errors.New("database unavailable")
	ErrDatabaseConnection     = errors.New("database connection failed")
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

// ErrorWrapper wraps errors with additional context
type ErrorWrapper struct {
	Operation string
	EventID   string
	TenantID  string
	Err       error
}

func (e *ErrorWrapper) Error() string {
	if e.EventID != "" && e.TenantID != "" {
		return fmt.Sprintf("%s failed for event %s (tenant %s): %v", e.Operation, e.EventID, e.TenantID, e.Err)
	}
	if e.EventID != "" {
		return fmt.Sprintf("%s failed for event %s: %v", e.Operation, e.EventID, e.Err)
	}
	if e.TenantID != "" {
		return fmt.Sprintf("%s failed for tenant %s: %v", e.Operation, e.TenantID, e.Err)
	}
	return fmt.Sprintf("%s failed: %v", e.Operation, e.Err)
}

func (e *ErrorWrapper) Unwrap() error {
	return e.Err
}

// WrapError creates a new ErrorWrapper with context
func WrapError(operation string, err error, eventID, tenantID string) *ErrorWrapper {
	return &ErrorWrapper{
		Operation: operation,
		EventID:   eventID,
		TenantID:  tenantID,
		Err:       err,
	}
}

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
