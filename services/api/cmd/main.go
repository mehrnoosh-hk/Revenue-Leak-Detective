package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"rdl-api/config"
	"rdl-api/internal/app"
)

var (
	Version   string
	Commit    string
	BuildDate string
)

func main() {

	fmt.Println(Version, Commit, BuildDate)
	// Parse command line flags
	flags := parseFlags()

	// Load configuration
	cfg, err := config.LoadConfig(flags.EnvFile)
	if err != nil {
		slog.Error("Failed to load configuration", "error", err, "env_file", flags.EnvFile)
		os.Exit(1)
	}

	flags.handleVersionFlag(cfg)

	slog.Info("Starting Revenue Leak Detective API")

	ctx := context.Background()

	// Create application
	application, err := app.NewApplication(ctx, cfg)
	if err != nil {
		slog.Error("Failed to create application", "error", err)
		os.Exit(1)
	}

	flags.handleHealthFlag(ctx, application)

	if err := application.StartUp(ctx); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
