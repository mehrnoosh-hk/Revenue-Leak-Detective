package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Write a test for the HealthCheckHandler
func TestHealthCheckHandler(t *testing.T) {
	// Create a test logger that discards output
	testLogger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Test GET request
	t.Run("GET request returns OK", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/healthz", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := HealthCheckHandler(testLogger)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		contentType := rr.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("handler returned wrong content type: got %v want %v",
				contentType, "application/json")
		}

		// Check that response contains status field
		var gotBody struct {
			Status string `json:"status"`
			Timestamp time.Time `json:"timestamp"`
			Version string `json:"version"`
		}
		if err := json.Unmarshal(rr.Body.Bytes(), &gotBody); err != nil {
			t.Errorf("failed to unmarshal response body, invalid json: %v", err)
		}
		if gotBody.Status != "OK" {
			t.Errorf("handler returned wrong status: got %v want %v",
				gotBody.Status, "OK")
		}

		// check the timestamp and version is not empty
		if gotBody.Timestamp.IsZero() {
			t.Errorf("handler returned empty timestamp")
		}
		if gotBody.Version == "" {
			t.Errorf("handler returned empty version")
		}
	})

	// Test non-GET request returns Method Not Allowed
	t.Run("POST request returns Method Not Allowed", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/healthz", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := HealthCheckHandler(testLogger)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusMethodNotAllowed)
		}
	})
}
