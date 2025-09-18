package services

import "errors"

var (
	ErrDatabaseNotInitialized = errors.New("database not initialized")
	ErrDatabaseUnavailable    = errors.New("database unavailable")

	// Service construction errors
	ErrLoggerCannotBeNil = errors.New("logger cannot be nil")
	ErrPoolCannotBeNil   = errors.New("pool cannot be nil")
)
