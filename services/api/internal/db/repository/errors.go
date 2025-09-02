package repository

import "errors"

// Database repository errors
var (
	ErrDatabaseNotInitialized = errors.New("database not initialized")
	ErrDatabaseUnavailable    = errors.New("database unavailable")
	ErrDatabaseConnection     = errors.New("database connection failed")
)
