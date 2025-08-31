package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"rdl-api/config"
	"rdl-api/internal/app"
	"time"

	"github.com/lmittmann/tint"
)

// Build info - TODO: Use build flags to set these values
var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func main() {
	// Parse command line flags
	versionFlag := flag.Bool("version", false, "Show version information")
	healthFlag := flag.Bool("health", false, "Run health check and exit")
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Printf("Revenue Leak Detective API\n")
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("Built: %s\n", Date)
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Setup logger
	handlerOptions := &slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.LogLevel,
	}

	var logger *slog.Logger
	if cfg.IsDevelopment() {
		logger = slog.New(tint.NewHandler(os.Stdout, &tint.Options{
			Level:      cfg.LogLevel,
			TimeFormat: time.RFC3339,
			AddSource:  true,
			NoColor:    false,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))
	}

	// Log startup information
	logger.Info("Starting Revenue Leak Detective API",
		slog.String("version", Version),
		slog.String("commit", Commit),
		slog.String("build_date", Date),
		slog.String("environment", cfg.Env))

	logger.Info("Initializing Application")

	application := app.New(cfg, logger)

	// TODO: Handle health check flag (for Docker health checks)
	if *healthFlag {
		// TODO:Perform health check logic here
		fmt.Println("Health check: OK")
		os.Exit(0)
	}

	ctx := context.Background()
	if err := application.Start(ctx); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
