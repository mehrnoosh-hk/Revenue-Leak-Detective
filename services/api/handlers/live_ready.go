package handlers

import (
    "context"
    "encoding/json"
    "log/slog"
    "net/http"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type probeResponse struct {
    Status    string    `json:"status"`
    Timestamp time.Time `json:"timestamp"`
}

// LiveHandler returns a simple liveness check handler.
// It should be fast and avoid any external dependencies.
func LiveHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        resp := probeResponse{Status: "OK", Timestamp: time.Now().UTC()}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
    }
}

// ReadyHandler checks readiness by pinging the database.
// If the DB is unavailable, it returns 503 to indicate not ready.
func ReadyHandler(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
        defer cancel()

        if db == nil {
            http.Error(w, "DB not initialized", http.StatusServiceUnavailable)
            return
        }

        if err := db.Ping(ctx); err != nil {
            if logger != nil {
                logger.Error("readiness DB ping failed", slog.Any("error", err))
            }
            http.Error(w, "Not Ready", http.StatusServiceUnavailable)
            return
        }

        resp := probeResponse{Status: "OK", Timestamp: time.Now().UTC()}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
    }
}