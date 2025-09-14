// Package app provides the struct and methods for the App
// app.go handles the main application methods and dependencies for the App
package app

import (
	"context"
	"fmt"
	"log/slog"
	"rdl-api/config"
	"rdl-api/internal/domain/services"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// App represents the main application with all its dependencies.
type App struct {
	// Core dependencies
	config *config.Config
	logger *slog.Logger

	// Repository layer
	pool *pgxpool.Pool

	// Domain services
	Services Services

	// Server layer
	server *AppServer
}

type Services struct {
	HealthService services.HealthService
	UsersService  services.UsersService
	EventsService services.EventsService
}

// New creates a new App instance with minimal dependencies properly initialized.
// This is the main entry point for dependency injection.
func New(cfg *config.Config) (*App, error) {
	logger := setupLogger(cfg)
	server := setupServer(cfg)
	pool, err := setupPgxPool(cfg)
	if err != nil {
		logger.Error("failed to create database connection pool", "error", err)
		return nil, err
	}
	Services := setupDomainServices(pool, logger)
	return &App{
		config:   cfg,
		logger:   logger,
		pool:     pool,
		Services: Services,
		server:   server,
	}, nil
}

// StartUp initializes and starts the application.
// It sets up all dependencies and starts the HTTP server.
func (a *App) StartUp(ctx context.Context) error {
	a.logger.Info("Starting application")
	a.logger.Info(fmt.Sprintf("Environment: %s", a.config.GetEnvironment()))
	a.logger.Info(fmt.Sprintf("Port: %s", a.config.GetPort()))

	// Verify database connectivity with a short timeout
	a.logger.Info("Verifying database connection at app Startup")
	if err := a.Services.HealthService.CheckReadiness(ctx); err != nil {
		return err
	}
	a.logger.Info("Database connection verified")

	// Start the server
	a.logger.Info("Server is ready to accept requests")
	return a.server.Start(ctx, a.logger, a.Services, a.config.IsDevelopment())
}

// Shutdown gracefully shuts down the application.
// It closes database connections and stops the server.
func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down application")

	var shutdownErrors []error

	// Close database connection pool
	if a.pool != nil {
		a.pool.Close()
		a.logger.Info("Database connection pool closed")
	}

	// Shutdown server with timeout
	if a.server != nil {
		serverCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		if err := a.server.server.Shutdown(serverCtx); err != nil {
			shutdownErrors = append(shutdownErrors, fmt.Errorf("server shutdown failed: %w", err))
		}
	}

	if len(shutdownErrors) > 0 {
		return fmt.Errorf("shutdown errors: %v", shutdownErrors)
	}

	a.logger.Info("Graceful shutdown completed successfully")

	return nil
}
