package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	// Mock handler that writes "handled" to response
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("handled"))
	})

	// Mock middleware that adds headers
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware-1", "applied")
			next.ServeHTTP(w, r)
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware-2", "applied")
			next.ServeHTTP(w, r)
		})
	}

	// Chain middlewares
	chainedHandler := Chain(handler, middleware1, middleware2)

	// Test request
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	chainedHandler.ServeHTTP(rr, req)

	// Verify response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "handled", rr.Body.String())
	assert.Equal(t, "applied", rr.Header().Get("X-Middleware-1"))
	assert.Equal(t, "applied", rr.Header().Get("X-Middleware-2"))
}

func TestLogger(t *testing.T) {
	// Create buffer to capture log output
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Mock handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Apply logger middleware
	loggerMiddleware := Logger(logger)
	wrappedHandler := loggerMiddleware(handler)

	// Test request
	req := httptest.NewRequest("GET", "/test?param=value", nil)
	req.Header.Set("User-Agent", "test-agent")
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	// Verify response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())

	// Verify log output contains expected fields
	logOutput := buf.String()
	assert.Contains(t, logOutput, "HTTP request")
	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "/test")
	assert.Contains(t, logOutput, "test-agent")
	assert.Contains(t, logOutput, "200")
}

func TestRecovery(t *testing.T) {
	// Create buffer to capture log output
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Mock handler that panics
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	// Apply recovery middleware
	recoveryMiddleware := Recovery(logger)
	wrappedHandler := recoveryMiddleware(handler)

	// Test request
	req := httptest.NewRequest("GET", "/panic", nil)
	rr := httptest.NewRecorder()

	// This should not panic due to recovery middleware
	assert.NotPanics(t, func() {
		wrappedHandler.ServeHTTP(rr, req)
	})

	// Verify response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "Internal Server Error\n", rr.Body.String())

	// Verify panic was logged
	logOutput := buf.String()
	assert.Contains(t, logOutput, "Panic recovered")
	assert.Contains(t, logOutput, "test panic")
	assert.Contains(t, logOutput, "/panic")
}

func TestCORS(t *testing.T) {
	// Mock handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Apply CORS middleware
	corsMiddleware := CORS()
	wrappedHandler := corsMiddleware(handler)

	t.Run("regular request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rr, req)

		// Verify CORS headers
		assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, "OK", rr.Body.String())
	})

	t.Run("OPTIONS request", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/", nil)
		rr := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rr, req)

		// Verify CORS headers
		assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, http.StatusNoContent, rr.Code)
		assert.Empty(t, rr.Body.String())
	})
}

func TestResponseWriter(t *testing.T) {
	// Create a mock ResponseWriter
	rr := httptest.NewRecorder()
	rw := &responseWriter{
		ResponseWriter: rr,
		statusCode:     http.StatusOK,
	}

	// Test WriteHeader
	rw.WriteHeader(http.StatusCreated)
	assert.Equal(t, http.StatusCreated, rw.statusCode)
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Test Write
	data := []byte("test data")
	n, err := rw.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)
	assert.Equal(t, "test data", rr.Body.String())
}

func TestMultipleMiddleware(t *testing.T) {
	// Create buffer to capture log output
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Mock handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Chain multiple middlewares
	chainedHandler := Chain(
		handler,
		Logger(logger),
		Recovery(logger),
		CORS(),
	)

	// Test request
	req := httptest.NewRequest("POST", "/api/test", strings.NewReader(`{"test": "data"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	chainedHandler.ServeHTTP(rr, req)

	// Verify response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())

	// Verify CORS headers
	assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))

	// Verify logging occurred
	logOutput := buf.String()
	assert.Contains(t, logOutput, "HTTP request")
	assert.Contains(t, logOutput, "POST")
}
