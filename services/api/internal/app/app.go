package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"rdl-api/config"
	"rdl-api/internal/db/repository"
	sqlc "rdl-api/internal/db/sqlc"
	"rdl-api/internal/domain/health"
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

	// Repository layer
	repoLayer *RepoLayer

	// Domain services
	domainServices *DomainServices

	// Server layer
	server *Server
}

type RepoLayer struct {
	healthRepo repository.HealthRepository
}

type DomainServices struct {
	healthService health.HealthService
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
		config:         cfg,
		logger:         logger,
		repoLayer:      &RepoLayer{},      // Initialize the struct
		domainServices: &DomainServices{}, // Initialize the struct
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

// setupDependencies creates all application dependencies
func (a *App) setupDependencies() error {
	if a.db == nil {
		return ErrDatabaseNotInitialized
	}

	// Create repository layer
	pgxAdapter := repository.NewPgxAdapter(a.db)
	a.repoLayer.healthRepo = repository.NewHealthRepository(pgxAdapter)

	// Create domain services
	a.domainServices.healthService = health.NewHealthService(a.repoLayer.healthRepo)

	return nil
}

func (a *App) StartServer(ctx context.Context) error {
	return a.server.Start(ctx, a.logger, a.domainServices)
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

// RepoLayer returns the repository layer.
func (a *App) RepoLayer() *RepoLayer {
	return a.repoLayer
}

// DomainServices returns the domain services.
func (a *App) DomainServices() *DomainServices {
	return a.domainServices
}

// Server returns the HTTP server.
func (a *App) Server() *Server {
	return a.server
}

// StartUp initializes and starts the application.
// It sets up all dependencies and starts the HTTP server.
func (a *App) StartUp(ctx context.Context) error {
	a.logger.Info("Starting application")
	a.logger.Info(fmt.Sprintf("Environment: %s", a.config.GetEnvironment()))
	a.logger.Info(fmt.Sprintf("Port: %s", a.config.GetPort()))

	// create a *pgxpool.Pool instance from a.config.PostgresURL
	db, err := pgxpool.New(ctx, a.config.DatabaseURL())
	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}

	// Verify database connectivity with a short timeout
	a.logger.Info("Verifying database connection at app Startup")
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Ping(pingCtx); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}
	a.logger.Info("Database connection verified")

	// Set the database connection
	a.SetDatabase(db)

	// Setup all dependencies
	if err := a.setupDependencies(); err != nil {
		return fmt.Errorf("failed to setup dependencies: %w", err)
	}

	// Validate that all required dependencies are set
	if a.server == nil {
		return ErrServerNotInitialized
	}

	// Start the server
	a.logger.Info("Server is ready to accept requests")
	return a.server.Start(ctx, a.logger, a.domainServices)
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
