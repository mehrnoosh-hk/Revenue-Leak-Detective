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

type flags struct {
	Version bool
	Health  bool
	EnvFile string
}

func parseFlags() flags {
	// Parse command line flags
	version := flag.Bool("version", false, "Show version information")
	health := flag.Bool("health", false, "Run health check and exit")
	envFile := flag.String("env-file", "", "Path to environment file (required)")
	flag.Parse()

	parsedFlags := flags{
		Version: *version,
		Health:  *health,
		EnvFile: *envFile,
	}

	if *envFile == "" {
		slog.Warn("No env file provided, using only environment variables or fallback defaults (development only)")
	}

	return parsedFlags
}

func (f flags) handleHealthFlag(ctx context.Context, application *app.Application) {
	if f.Health {
		err := application.CheckReadiness(ctx)
		if err != nil {
			slog.Error("Health check failed", "error", err)
			os.Exit(1)
		}
		fmt.Println("Health check: OK")
		os.Exit(0)
	}
}

func (f flags) handleVersionFlag(cfg *config.Config) {
	if f.Version {
		fmt.Printf("Revenue Leak Detective API\n")
		fmt.Printf("Version: %s\n", cfg.BuildInfo.GIT_TAG)
		fmt.Printf("Commit: %s\n", cfg.BuildInfo.GIT_COMMIT_FULL)
		fmt.Printf("Built: %s\n", cfg.BuildInfo.BUILD_TIMESTAMP)
		os.Exit(0)
	}
}
