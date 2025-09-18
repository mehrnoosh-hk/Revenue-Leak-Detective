package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"rdl-api/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
)

// Dependency container
type Container struct {
	config   *config.Config
	logger   *slog.Logger
	pool     *pgxpool.Pool
	services Services
}

func NewContainer(ctx context.Context, cfg *config.Config) (*Container, error) {
	logger := setupLogger(cfg)
	pool, err := setupPgxPool(ctx, cfg)
	if err != nil {
		logger.Error("failed to create database connection pool", "error", err)
		return nil, err
	}

	services := setupDomainServices(pool, logger, cfg.BuildInfo.GIT_TAG) // TODO: write a function to get the version

	return &Container{
		config:   cfg,
		logger:   logger,
		pool:     pool,
		services: services,
	}, nil
}

// setupPgxPool
func setupPgxPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFaildToCreateConnectionPool, err)
	}
	return pool, nil
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

func (c *Container) Shutdown(ctx context.Context) {
	if c.pool != nil {
		c.pool.Close()
	}
}

func (c *Container) GetEnvironment() string {
	return c.config.GetEnvironment()
}

func (c *Container) IsDevelopment() bool {
	return c.config.IsDevelopment()
}

func (c *Container) GetLogLevel() slog.Level {
	return c.config.GetLogLevel()
}

func (c *Container) GetPort() string {
	return c.config.HTTP.Port
}

func (c *Container) GetConfig() *config.Config {
	return c.config
}

func (c *Container) GetLogger() *slog.Logger {
	return c.logger
}

func (c *Container) GetPool() *pgxpool.Pool {
	return c.pool
}

func (c *Container) GetServices() Services {
	return c.services
}
