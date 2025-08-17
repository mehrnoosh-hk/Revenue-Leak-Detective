package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChain(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("handled"))
		require.NoError(t, err)
	})

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

	chainedHandler := Chain(handler, middleware1, middleware2)
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	chainedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "handled", rr.Body.String())
	assert.Equal(t, "applied", rr.Header().Get("X-Middleware-1"))
	assert.Equal(t, "applied", rr.Header().Get("X-Middleware-2"))
}

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		require.NoError(t, err)
	})

	wrappedHandler := Logger(logger)(handler)
	req := httptest.NewRequest("GET", "/test?param=value", nil)
	req.Header.Set("User-Agent", "test-agent")
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())

	logOutput := buf.String()
	assert.Contains(t, logOutput, "HTTP request")
	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "/test")
	assert.Contains(t, logOutput, "test-agent")
	assert.Contains(t, logOutput, "200")
}

func TestRecovery(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	wrappedHandler := Recovery(logger)(handler)
	req := httptest.NewRequest("GET", "/panic", nil)
	rr := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		wrappedHandler.ServeHTTP(rr, req)
	})

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "Internal Server Error\n", rr.Body.String())

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Panic recovered")
	assert.Contains(t, logOutput, "test panic")
	assert.Contains(t, logOutput, "/panic")
}

func TestCORS(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		require.NoError(t, err)
	})

	wrappedHandler := CORS()(handler)

	t.Run("regular request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rr, req)

		assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, "OK", rr.Body.String())
	})

	t.Run("OPTIONS request", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/", nil)
		rr := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rr, req)

		assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, http.StatusNoContent, rr.Code)
		assert.Empty(t, rr.Body.String())
	})
}

func TestResponseWriter(t *testing.T) {
	rr := httptest.NewRecorder()
	rw := &responseWriter{
		ResponseWriter: rr,
		statusCode:     http.StatusOK,
	}

	rw.WriteHeader(http.StatusCreated)
	assert.Equal(t, http.StatusCreated, rw.statusCode)
	assert.Equal(t, http.StatusCreated, rr.Code)

	data := []byte("test data")
	n, err := rw.Write(data)
	require.NoError(t, err)
	assert.Equal(t, len(data), n)
	assert.Equal(t, "test data", rr.Body.String())
}

func TestMultipleMiddleware(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		require.NoError(t, err)
	})

	chainedHandler := Chain(handler, Logger(logger), Recovery(logger), CORS())
	req := httptest.NewRequest("POST", "/api/test", strings.NewReader(`{"test": "data"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	chainedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
	assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))

	logOutput := buf.String()
	assert.Contains(t, logOutput, "HTTP request")
	assert.Contains(t, logOutput, "POST")
}
