package handlers

import (
	"log/slog"
	"rdl-api/config"
)

// HandlerDependencies provides dependencies for handlers
type HandlerDependencies struct {
	Config *config.Config
	Logger *slog.Logger
}
