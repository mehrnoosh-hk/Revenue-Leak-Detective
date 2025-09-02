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

func main() {
	// Parse command line flags
	versionFlag := flag.Bool("version", false, "Show version information")
	healthFlag := flag.Bool("health", false, "Run health check and exit")
	envFileFlag := flag.String("env-file", "", "Path to environment file (required)")
	flag.Parse()

	if *envFileFlag == "" {
		slog.Warn("No env file provided, using only environment variables or fallbacks to defaults(development only)")
	}

	// Load configuration
	cfg, err := config.LoadConfig(*envFileFlag)
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Handle version flag
	if *versionFlag {
		handleVersionFlag(cfg)
	}

	// Setup logger
	logger := setupLogger(cfg)

	// Log startup information
	logger.Info("Starting Revenue Leak Detective API")

	cfg.PrintBuildInfo(logger)
	cfg.PrintEffectiveConfig(logger)

	application := app.New(cfg, logger)

	// TODO: Handle health check flag (for Docker health checks)
	if *healthFlag {
		// TODO:Perform health check logic here
		fmt.Println("Health check: OK")
		os.Exit(0)
	}

	ctx := context.Background()
	if err := application.StartUp(ctx); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

func handleVersionFlag(cfg *config.Config) {
	fmt.Printf("Revenue Leak Detective API\n")
	fmt.Printf("Version: %s", cfg.BuildInfo.GIT_TAG)
	fmt.Printf("Commit: %s", cfg.BuildInfo.GIT_COMMIT_FULL)
	fmt.Printf("Built: %s", cfg.BuildInfo.BUILD_TIMESTAMP)
	os.Exit(0)
}

func setupLogger(cfg *config.Config) *slog.Logger {
	var logger *slog.Logger
	if cfg.IsDevelopment() {
		logger = slog.New(tint.NewHandler(os.Stdout, &tint.Options{
			Level:      cfg.GetLogLevel(),
			TimeFormat: time.RFC3339,
			AddSource:  true,
			NoColor:    false,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     cfg.GetLogLevel(),
			AddSource: true,
		}))
	}
	return logger
}
