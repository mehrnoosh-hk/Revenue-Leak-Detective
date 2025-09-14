package main

import (
	"context"
	"log/slog"
	"os"
	"rdl-api/config"
	"rdl-api/internal/app"
)

func main() {
	// Parse command line flags
	flags := parseFlags()

	// Load configuration
	cfg, err := config.LoadConfig(flags.EnvFile)
	if err != nil {
		slog.Error("Failed to load configuration", "error", err, "env_file", flags.EnvFile)
		os.Exit(1)
	}

	slog.Info("Starting Revenue Leak Detective API")

	// Create application
	application, err := app.New(cfg)
	if err != nil {
		slog.Error("Failed to create application", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// Handle flags
	flags.HandleFlags(ctx, application, cfg)

	if err := application.StartUp(ctx); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
