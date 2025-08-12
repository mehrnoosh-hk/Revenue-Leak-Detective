package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"rld/services/api/config"
	"strings"
	"testing"
)

// Write a test for the HealthCheckHandler
func TestHealthCheckHandler(t *testing.T) {
	// Create a test logger that discards output
	testLogger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	hd := &HandlerDependencies{
		Config: &config.Config{Port: "8080", LogLevel: slog.LevelInfo, Env: "test"},
		Logger: testLogger,
	}

	// Test GET request
	t.Run("GET request returns OK", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/healthz", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := HealthCheckHandler(hd)
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
		if !strings.Contains(rr.Body.String(), `"status":"OK"`) {
			t.Errorf("handler returned unexpected body: %v", rr.Body.String())
		}
	})

	// Test non-GET request returns Method Not Allowed
	t.Run("POST request returns Method Not Allowed", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/healthz", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := HealthCheckHandler(hd)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusMethodNotAllowed)
		}
	})
}
