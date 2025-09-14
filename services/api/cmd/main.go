package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"rdl-api/config"
	"rdl-api/internal/app"
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

	// Log startup information
	slog.Info("Starting Revenue Leak Detective API")

	// Move to app package
	// cfg.PrintBuildInfo(logger)
	// cfg.PrintEffectiveConfig(logger)

	application, err := app.New(cfg)
	if err != nil {
		slog.Error("Failed to create application", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	if *healthFlag {
		err := application.CheckReadiness(ctx)
		if err != nil {
			slog.Error("Health check failed", "error", err)
			os.Exit(1)
		}
		fmt.Println("Health check: OK")
		os.Exit(0)
	}

	if err := application.StartUp(ctx); err != nil {
		slog.Error("Server failed to start", "error", err)
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
