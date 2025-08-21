package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
}

// HealthCheckHandler returns a health check handler
func HealthCheckHandler(hd *HandlerDependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		hd.Logger.Debug("Health check endpoint accessed",
			slog.String("remote_addr", r.RemoteAddr))

		response := HealthResponse{
			Status:    "OK",
			Timestamp: time.Now().UTC(),
			Version:   "1.0.0", // Could be loaded from build info
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			hd.Logger.Error("Failed to encode health response",
				slog.Any("error", err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
