package main

import (
	"log"
	"log/slog"
	"net/http"
	"rld/services/api/config"
	"rld/services/api/handlers"
)

type APIService struct {
	Config   *config.Config
	Servemux *http.ServeMux
	Logger *slog.Logger
}

func main() {
	// Initialize configuration
	cfg, _ := config.LoadConfig()

	// Set up logging with level from config and write to file go_log and standart output
	logger := slog.New(slog.NewTextHandler(log.Writer(), &slog.HandlerOptions{
		Level: slog.Level(cfg.LogLevel),
	}))
	slog.SetDefault(logger)
	

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/health", handlers.HealthCheckHandler)

	// Create API service instance
	apiService := &APIService{
		Config:   cfg,
		Servemux: mux,
		Logger: logger,
	}

	// Start the HTTP server
	apiService.Logger.Info("Starting API service", "port", apiService.Config.Port)
	err := http.ListenAndServe(":"+apiService.Config.Port, apiService.Servemux)
	if err != nil {
		apiService.Logger.Error("Failed to start server", "error", err)
		return
	}
}
