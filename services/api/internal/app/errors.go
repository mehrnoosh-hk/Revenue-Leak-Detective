package app

import "errors"

var (
	ErrFaildToCreateConnectionPool = errors.New("failed to create connection pool")
	ErrDatabaseNotInitialized      = errors.New("database not initialized")
	ErrServerNotInitialized        = errors.New("server not initialized")
	ErrDatabaseConnection          = errors.New("database connection failed")
	ErrServerStartup               = errors.New("server startup failed")
)
