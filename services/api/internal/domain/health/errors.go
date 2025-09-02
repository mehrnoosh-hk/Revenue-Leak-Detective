package health

import "errors"

// Health domain errors
var (
	ErrDatabaseNotInitialized = errors.New("database not initialized")
	ErrDatabaseUnavailable    = errors.New("database unavailable")
)
