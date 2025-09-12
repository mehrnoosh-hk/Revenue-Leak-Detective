package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"rdl-api/handlers"
	"rdl-api/internal/middleware"
	"syscall"
	"time"
)

type Server struct {
	mux    *http.ServeMux
	server *http.Server
}

func (s *Server) Start(ctx context.Context, logger *slog.Logger, services *DomainServices, isDevelopment bool) error {

	// Setup routes with services reference
	s.SetupRoutes(logger, services, isDevelopment)

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", slog.Any("error", err))
			quit <- syscall.SIGTERM
		}
	}()

	// Wait for interrupt signal
	select {
	case sig := <-quit:
		logger.Info("Shutting down server", slog.String("signal", sig.String()))
	case <-ctx.Done():
		logger.Info("Shutting down server", slog.String("reason", "context canceled"))
	}

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	logger.Info("Server stopped gracefully")
	return nil
}

func (s *Server) SetupRoutes(logger *slog.Logger, services *DomainServices, isDevelopment bool) {
	// Apply middleware
	handler := middleware.Chain(
		s.mux,
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Recovery(logger),
		middleware.CORS(),
		middleware.TenantContext(logger, isDevelopment),
	)

	s.server.Handler = handler

	// Register routes
	s.mux.HandleFunc("/healthz", handlers.HealthCheckHandler(logger))
	s.mux.HandleFunc("/health", handlers.HealthCheckHandler(logger)) // Alternative endpoint
	s.mux.HandleFunc("/live", handlers.LiveHandler(services.healthService, logger))
	s.mux.HandleFunc("/ready", handlers.ReadyHandler(services.healthService, logger))
}
