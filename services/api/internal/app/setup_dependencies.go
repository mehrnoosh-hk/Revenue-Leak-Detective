package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"rdl-api/config"
	"rdl-api/internal/domain/services"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
)

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

func setupServer(cfg *config.Config) *AppServer {
	mux := http.NewServeMux()
	return &AppServer{
		mux: mux,
		server: &http.Server{
			Addr:         ":" + cfg.GetPort(),
			Handler:      mux,
			ReadTimeout:  15 * time.Second, // TODO: Should be loaded from config
			WriteTimeout: 15 * time.Second, // TODO: Should be loaded from config
			IdleTimeout:  60 * time.Second, // TODO: Should be loaded from config
		},
	}
}

// setupPgxPool
func setupPgxPool(cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFaildToCreateConnectionPool, err)
	}
	return pool, nil
}

// setupDomainServices
func setupDomainServices(pool *pgxpool.Pool, logger *slog.Logger) Services {

	hService := services.NewHealthService(pool, logger)
	uService := services.NewUserService(pool, logger)
	eService := services.NewEventService(pool, logger)

	domainServices := Services{
		HealthService: hService,
		UsersService:  uService,
		EventsService: eService,
	}

	return domainServices
}
