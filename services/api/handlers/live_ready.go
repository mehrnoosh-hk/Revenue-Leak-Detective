package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"rdl-api/internal/domain/services"
	"time"
)

type probeResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// LiveHandler returns a simple liveness check handler.
// It should be fast and avoid any external dependencies.
func LiveHandler(logger *slog.Logger, healthService services.HealthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
			return
		}

		// Liveness check should be fast and not depend on external services
		if err := healthService.CheckLiveness(r.Context()); err != nil {
			if logger != nil {
				logger.Error("liveness check failed", slog.Any("error", err))
			}
			http.Error(w, ErrNotAlive.Error(), http.StatusInternalServerError)
			return
		}

		resp := probeResponse{Status: "OK", Timestamp: time.Now().UTC()}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			if logger != nil {
				logger.Error("failed to encode response", slog.Any("error", err))
			}
			http.Error(w, ErrInternalServerError.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// ReadyHandler checks readiness by using the health service.
// If the service is unavailable, it returns 503 to indicate not ready.
func ReadyHandler(logger *slog.Logger, healthService services.HealthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := healthService.CheckReadiness(r.Context()); err != nil {
			if logger != nil {
				logger.ErrorContext(r.Context(), "readiness check failed", "error", err)
			}
			http.Error(w, ErrNotReady.Error(), http.StatusServiceUnavailable)
			return
		}

		resp := probeResponse{Status: "OK", Timestamp: time.Now().UTC()}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			if logger != nil {
				logger.ErrorContext(r.Context(), "failed to encode response", "error", err)
			}
			http.Error(w, ErrInternalServerError.Error(), http.StatusInternalServerError)
			return
		}
	}
}
