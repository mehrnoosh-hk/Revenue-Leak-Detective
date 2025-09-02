package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"rdl-api/config"
	sqlc "rdl-api/internal/db/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// App represents the main application with all its dependencies.
// It follows Go best practices for dependency injection and service composition.
type App struct {
	// Core dependencies
	config *config.Config
	logger *slog.Logger

	// Database layer
	db      *pgxpool.Pool
	queries *sqlc.Queries

	// Server layer
	server *Server
}

// Application errors
// Use custom error types
var (
	ErrDatabaseNotInitialized = errors.New("database not initialized")
	ErrServerNotInitialized   = errors.New("server not initialized")
	ErrDatabaseConnection     = errors.New("database connection failed")
	ErrServerStartup          = errors.New("server startup failed")
)

// New creates a new App instance with basic dependencies properly initialized.
// This is the main entry point for dependency injection.
func New(cfg *config.Config, logger *slog.Logger) *App {

	mux := http.NewServeMux()
	return &App{
		config: cfg,
		logger: logger,
		server: &Server{
			mux: mux,
			server: &http.Server{
				Addr:         ":" + cfg.GetPort(),
				Handler:      mux,
				ReadTimeout:  15 * time.Second,
				WriteTimeout: 15 * time.Second,
				IdleTimeout:  60 * time.Second,
			},
		},
	}
}

// SetDatabase sets the database connection and queries for the application.
// This method allows for flexible database initialization and testing.
func (a *App) SetDatabase(db *pgxpool.Pool) {
	a.db = db
	a.queries = sqlc.New(db)
}

// SetServer sets the HTTP server for the application.
// This method allows for flexible server initialization and testing.
func (a *App) SetServer(server *Server) {
	a.server = server
}

func (a *App) StartServer(ctx context.Context) error {

	return a.server.Start(ctx, a.logger)
}

// Config returns the application configuration.
func (a *App) Config() *config.Config {
	return a.config
}

// Logger returns the application logger.
func (a *App) Logger() *slog.Logger {
	return a.logger
}

// DB returns the database connection.
func (a *App) DB() *pgxpool.Pool {
	return a.db
}

// Queries returns the sqlc-generated database queries.
func (a *App) Queries() *sqlc.Queries {
	return a.queries
}

// Server returns the HTTP server.
func (a *App) Server() *Server {
	return a.server
}

// Start initializes and starts the application.
// It sets up all dependencies and starts the HTTP server.
func (a *App) Start(ctx context.Context) error {
	a.logger.Info("Starting application")
	a.logger.Info(fmt.Sprintf("Environment: %s", a.config.GetEnvironment()))
	a.logger.Info(fmt.Sprintf("Port: %s", a.config.GetPort()))

	// create a *pgxpool.Pool instance from a.config.PostgresURL
	db, err := pgxpool.New(ctx, a.config.DatabaseURL())
	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}

	// Verify database connectivity with a short timeout
	a.logger.Info("Verifying database connection")
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Ping(pingCtx); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}
	a.logger.Info("Database connection verified")
	a.SetDatabase(db)

	// Validate that all required dependencies are set
	if a.db == nil {
		return ErrDatabaseNotInitialized
	}
	if a.server == nil {
		return ErrServerNotInitialized
	}

	// Start the server
	a.logger.Info("Server is ready to accept requests")
	return a.server.Start(ctx, a.logger)
}

// HealthCheck performs a comprehensive health check of the application.
// It verifies database connectivity and other critical dependencies.
func (a *App) HealthCheck(ctx context.Context) error {
	// Check if database is initialized
	if a.db == nil {
		return fmt.Errorf("database not initialized")
	}

	// Test database connectivity with a short timeout
	healthCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := a.db.Ping(healthCtx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	// Could add more health checks here (Redis, external APIs, etc.)
	return nil
}

// Shutdown gracefully shuts down the application.
// It closes database connections and stops the server.
func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down application")

	var shutdownErrors []error

	// Close database connection
	if a.db != nil {
		a.db.Close()
		a.logger.Info("Database connection closed")
	}

	// Shutdown server with timeout
	if a.server != nil {
		serverCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		if err := a.server.server.Shutdown(serverCtx); err != nil {
			shutdownErrors = append(shutdownErrors, fmt.Errorf("server shutdown failed: %w", err))
		} else {
			a.logger.Info("Server shutdown completed")
		}
	}

	if len(shutdownErrors) > 0 {
		return fmt.Errorf("shutdown errors: %v", shutdownErrors)
	}

	a.logger.Info("Graceful shutdown completed successfully")

	return nil
}
