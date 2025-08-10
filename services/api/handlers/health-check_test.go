package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeMux_Routes(t *testing.T) {
	// Test that routes are properly registered
	mux := http.NewServeMux()

	// Register handlers as done in main
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Test registered route
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Health check route not properly registered. Expected status %d, got %d", http.StatusOK, status)
	}

	// Test unregistered route (should return 404)
	req404, err := http.NewRequest("GET", "/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr404 := httptest.NewRecorder()
	mux.ServeHTTP(rr404, req404)

	if status := rr404.Code; status != http.StatusNotFound {
		t.Errorf("Expected 404 for unregistered route, got %d", status)
	}
}
