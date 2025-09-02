package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"rdl-api/internal/domain/health"
)

type probeResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// LiveHandler returns a simple liveness check handler.
// It should be fast and avoid any external dependencies.
func LiveHandler(healthService health.HealthService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Liveness check should be fast and not depend on external services
		if err := healthService.CheckLiveness(r.Context()); err != nil {
			if logger != nil {
				logger.Error("liveness check failed", slog.Any("error", err))
			}
			http.Error(w, "Not Alive", http.StatusInternalServerError)
			return
		}

		resp := probeResponse{Status: "OK", Timestamp: time.Now().UTC()}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			if logger != nil {
				logger.Error("failed to encode response", slog.Any("error", err))
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

// ReadyHandler checks readiness by using the health service.
// If the service is unavailable, it returns 503 to indicate not ready.
func ReadyHandler(healthService health.HealthService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := healthService.CheckReadiness(r.Context()); err != nil {
			if logger != nil {
				logger.Error("readiness check failed", slog.Any("error", err))
			}
			http.Error(w, "Not Ready", http.StatusServiceUnavailable)
			return
		}

		resp := probeResponse{Status: "OK", Timestamp: time.Now().UTC()}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			if logger != nil {
				logger.Error("failed to encode response", slog.Any("error", err))
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
