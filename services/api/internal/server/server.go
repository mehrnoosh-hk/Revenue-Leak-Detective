package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rld/services/api/config"
	"rld/services/api/handlers"
	"rld/services/api/internal/middleware"
)

// Server represents the HTTP server
type Server struct {
	config     *config.Config
	logger     *slog.Logger
	httpServer *http.Server
	mux        *http.ServeMux
}

// New creates a new server instance
func New(cfg *config.Config, logger *slog.Logger) *Server {
	mux := http.NewServeMux()

	server := &Server{
		config: cfg,
		logger: logger,
		mux:    mux,
		httpServer: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      mux,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	return server
}

// SetupRoutes configures all the routes for the server
func (s *Server) SetupRoutes() {
	// Apply middleware
	handler := middleware.Chain(
		s.mux,
		middleware.Logger(s.logger),
		middleware.Recovery(s.logger),
		middleware.CORS(),
	)

	s.httpServer.Handler = handler

	// API service for handlers
	apiService := &handlers.HandlerDependencies{
		Config: s.config,
		Logger: s.logger,
	}

	// Register routes
	s.mux.HandleFunc("/healthz", handlers.HealthCheckHandler(apiService))
	s.mux.HandleFunc("/health", handlers.HealthCheckHandler(apiService)) // Alternative endpoint
}

// Start starts the HTTP server with graceful shutdown
func (s *Server) Start(ctx context.Context) error {
	// Setup routes
	s.SetupRoutes()

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		s.logger.Info("Starting HTTP server",
			slog.String("addr", s.httpServer.Addr),
			slog.String("env", s.config.Env))

		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Failed to start server", slog.Any("error", err))
			quit <- syscall.SIGTERM
		}
	}()

	// Wait for interrupt signal
	select {
	case sig := <-quit:
		s.logger.Info("Shutting down server", slog.String("signal", sig.String()))
	case <-ctx.Done():
		s.logger.Info("Shutting down server", slog.String("reason", "context cancelled"))
	}

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	s.logger.Info("Server stopped gracefully")
	return nil
}

// Stop stops the server
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
