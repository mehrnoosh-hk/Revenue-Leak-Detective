// Package app provides the struct and methods for the App
// server.go handles the HTTP server and routes for configurations and methods for the App server
package app

import (
	"log/slog"
	"net/http"
	"os"
	"rdl-api/handlers"
	"rdl-api/internal/middleware"
	"time"
)

func setupAppServer(c *Container) *AppServer {
	mux := http.NewServeMux()
	handler := SetupRoutes(mux, c)

	httpConfig := c.GetConfig().HTTP

	server := &http.Server{
		Addr:              httpConfig.Host + ":" + httpConfig.Port,
		Handler:           handler,
		ReadTimeout:       httpConfig.ReadTimeout * time.Second,
		ReadHeaderTimeout: httpConfig.ReadHeaderTimeout * time.Second,
		WriteTimeout:      httpConfig.WriteTimeout * time.Second,
		IdleTimeout:       httpConfig.IdleTimeout * time.Second,
	}

	return &AppServer{
		server: server,
	}
}

func SetupRoutes(mux *http.ServeMux, c *Container) http.Handler {
	// Define paths that should be excluded from tenant context validation
	// These are typically health check endpoints that don't require authentication
	excludedPaths := []string{
		"/healthz", // Kubernetes health check
		"/health",  // Alternative health check
		"/live",    // Liveness probe
		"/ready",   // Readiness probe
	}

	logger := c.GetLogger()
	services := c.GetServices()

	// Register routes
	mux.HandleFunc("/live", handlers.LiveHandler(logger, services.HealthService))
	mux.HandleFunc("/ready", handlers.ReadyHandler(logger, services.HealthService))

	isDevelopment := c.IsDevelopment()
	// Apply middleware
	return middleware.Chain(
		mux,
		middleware.Recovery(logger), // 1. Outermost - catch all panics
		middleware.CORS(),           // 2. Handle CORS early
		middleware.RequestID(),      // 3. Generate request ID early
		middleware.TenantContext(logger, isDevelopment, excludedPaths), // 4. Extract tenant context
		middleware.Logger(logger),                                      // 5. Innermost - log everything
	)
}

func Start(logger *slog.Logger, server *http.Server) {
	// Start server in a goroutine
	go func() {

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()
}
