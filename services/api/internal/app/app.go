// Package app provides the struct and methods for the App
// app.go handles the main application methods and dependencies for the App
package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"rdl-api/config"
	"syscall"
	"time"
)

// AppServer encapsulates the HTTP server
type AppServer struct {
	server *http.Server
}

// Application lifecycle manager
type Application struct {
	container *Container
	server    *AppServer
}

// NewApplication creates a new Application instance with minimal dependencies properly initialized.
func NewApplication(ctx context.Context, cfg *config.Config) (*Application, error) {
	container, err := NewContainer(ctx, cfg)
	if err != nil {
		return nil, err
	}

	appServer := setupAppServer(container)
	return &Application{
		container: container,
		server:    appServer,
	}, nil
}

// StartUp initializes and starts the application.
// It sets up all dependencies and starts the HTTP server.
func (a *Application) StartUp(ctx context.Context) error {
	l := a.container.GetLogger()
	c := a.container
	server := a.server.server
	l.Info("Starting application")
	l.Info(fmt.Sprintf("Environment: %s", c.GetEnvironment()))

	// Verify database connectivity with a short timeout
	if err := a.container.GetServices().HealthService.CheckReadiness(ctx); err != nil {
		return err
	}
	l.Info("Database connection verified")

	// Start the server
	Start(l, server)
	l.Info("Server is ready to accept requests on port " + server.Addr)

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal
	select {
	case sig := <-quit:
		l.Info("Shutting down Application", "signal", sig.String())
		err := a.Shutdown(ctx)
		if err != nil {
			return err
		}
	case <-ctx.Done():
		l.Info("Shutting down server", "reason", "context canceled")
		err := a.Shutdown(ctx)
		if err != nil {
			return err
		}
	}
	l.Info("Application shut down gracefully")
	return nil
}

// Shutdown gracefully shuts down the application.
// It closes database connections and stops the server.
func (a *Application) Shutdown(ctx context.Context) error {
	l := a.container.GetLogger()

	var shutdownErrors []error

	a.container.Shutdown(ctx)
	l.Info("Database connection pool closed")

	// Shutdown server with timeout
	if a.server != nil && a.server.server != nil {
		serverCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		if err := a.server.server.Shutdown(serverCtx); err != nil {
			shutdownErrors = append(shutdownErrors, fmt.Errorf("server shutdown failed: %w", err))
		}
	}

	if len(shutdownErrors) > 0 {
		return fmt.Errorf("shutdown errors: %v", shutdownErrors)
	}

	return nil
}

func (a *Application) CheckReadiness(ctx context.Context) error {
	return a.container.GetServices().HealthService.CheckReadiness(ctx)
}

func (a *Application) CheckLiveness(ctx context.Context) error {
	return a.container.GetServices().HealthService.CheckLiveness(ctx)
}
