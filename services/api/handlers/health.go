package handlers

import (
	"log/slog"
	"net/http"
	"rdl-api/internal/domain/services"
	"time"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
}

// ReadyHandler returns a health check handler
func ReadyHandler(logger *slog.Logger, healthService services.HealthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
			return
		}

		if err := healthService.CheckReadiness(r.Context()); err != nil {
			WriteJSONErrorResponse(r.Context(), w, logger, ErrHealthCheckFailed, http.StatusInternalServerError)
			return
		}

		logger.DebugContext(r.Context(), "Health check endpoint accessed", "remote_addr", r.RemoteAddr)

		response := HealthResponse{
			Status:    "OK",
			Timestamp: time.Now().UTC(),
			Version:   healthService.GetVersion(),
		}

		WriteJSONSuccessResponse(r.Context(), w, logger, response)
	}
}

func LiveHandler(logger *slog.Logger, healthService services.HealthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
			return
		}

		if err := healthService.CheckLiveness(r.Context()); err != nil {
			WriteJSONErrorResponse(r.Context(),
				w,
				logger,
				err,
				http.StatusInternalServerError)
			return
		}

		response := HealthResponse{
			Status:    "OK",
			Timestamp: time.Now().UTC(),
			Version:   healthService.GetVersion(),
		}

		WriteJSONSuccessResponse(r.Context(), w, logger, response)
	}
}
