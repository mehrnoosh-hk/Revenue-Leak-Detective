package services

import "errors"

var (
	ErrDatabaseNotInitialized = errors.New("database not initialized")
	ErrDatabaseUnavailable    = errors.New("database unavailable")
)
