package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// Test constants
const (
	expectedHealthyStatus = "OK"
	expectedVersion       = "1.0.0"
	testEndpoint          = "/healthz"
	testTimeout           = 100 * time.Millisecond
)

// Test error definitions
var (
	errTestService = errors.New("service unavailable")
)

// testCase represents a single test case for health check handler
type testCase struct {
	name           string
	method         string
	healthService  *testHealthService
	expectedStatus int
	expectedBody   string
	expectJSON     bool
	description    string
	timeout        time.Duration
}

// testHealthService is a mock implementation of the HealthService interface
type testHealthService struct {
	CheckReadinessFn func(ctx context.Context) error
	CheckLivenessFn  func(ctx context.Context) error
}

func (t *testHealthService) CheckReadiness(ctx context.Context) error {
	if t.CheckReadinessFn != nil {
		return t.CheckReadinessFn(ctx)
	}
	return nil
}

func (t *testHealthService) CheckLiveness(ctx context.Context) error {
	if t.CheckLivenessFn != nil {
		return t.CheckLivenessFn(ctx)
	}
	return nil
}

// Test builders for creating test instances
func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func newHealthyService() *testHealthService {
	return &testHealthService{
		CheckReadinessFn: func(ctx context.Context) error {
			return nil
		},
		CheckLivenessFn: func(ctx context.Context) error {
			return nil
		},
	}
}

func newUnhealthyService(err error) *testHealthService {
	if err == nil {
		err = errTestService
	}
	return &testHealthService{
		CheckReadinessFn: func(ctx context.Context) error {
			return err
		},
		CheckLivenessFn: func(ctx context.Context) error {
			return err
		},
	}
}

func newTimeoutService(delay time.Duration) *testHealthService {
	return &testHealthService{
		CheckReadinessFn: func(ctx context.Context) error {
			select {
			case <-time.After(delay):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
		CheckLivenessFn: func(ctx context.Context) error {
			select {
			case <-time.After(delay):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	}
}

// Test helpers
func createTestRequest(method string) *http.Request {
	req, err := http.NewRequest(method, testEndpoint, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to create test request: %v", err))
	}
	return req
}

func createTestRequestWithContext(method, url string, ctx context.Context) *http.Request {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to create test request with context: %v", err))
	}
	return req
}

func assertHealthResponse(t *testing.T, response HealthResponse, expectedStatus string) {
	t.Helper()
	if response.Status != expectedStatus {
		t.Errorf("Expected status %s, got %s", expectedStatus, response.Status)
	}
	if response.Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}
	if response.Version == "" {
		t.Error("Expected non-empty version")
	}
}

func assertHTTPResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int, expectedBody string, expectJSON bool) {
	t.Helper()

	// Verify status code
	if rr.Code != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, rr.Code)
	}

	// Verify response body
	body := rr.Body.String()
	if expectJSON {
		// For successful responses, verify JSON structure
		var response HealthResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to unmarshal JSON response: %v", err)
			return
		}
		assertHealthResponse(t, response, expectedBody)

		// Verify content type
		contentType := rr.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected content-type application/json, got %s", contentType)
		}
	} else {
		// For error responses, verify error message (http.Error adds newlines)
		expectedBodyWithNewline := expectedBody + "\n"
		if body != expectedBodyWithNewline {
			t.Errorf("Expected body %q, got %q", expectedBodyWithNewline, body)
		}
	}
}

func runHealthCheckTest(t *testing.T, tc testCase) {
	t.Helper()

	// Create request
	var req *http.Request
	if tc.timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
		defer cancel()
		req = createTestRequestWithContext(tc.method, testEndpoint, ctx)
	} else {
		req = createTestRequest(tc.method)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Create handler with test health service
	handler := HealthCheckHandler(newTestLogger(), tc.healthService)

	// Handle potential panic for nil service
	panicked := false
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			if tc.healthService == nil {
				// Expected panic for nil service
				t.Logf("Expected panic for nil service: %v", r)
			} else {
				t.Errorf("Unexpected panic: %v", r)
			}
		}
	}()

	// Execute request
	handler.ServeHTTP(rr, req)

	// Skip verification if we panicked (expected for nil service)
	if panicked && tc.healthService == nil {
		return
	}

	// Verify response
	assertHTTPResponse(t, rr, tc.expectedStatus, tc.expectedBody, tc.expectJSON)
}

// Test case generators
func generateGETTests() []testCase {
	return []testCase{
		{
			name:           "GET_healthy_service",
			method:         http.MethodGet,
			healthService:  newHealthyService(),
			expectedStatus: http.StatusOK,
			expectedBody:   expectedHealthyStatus,
			expectJSON:     true,
			description:    "GET request with healthy service should return 200 OK",
		},
		{
			name:           "GET_unhealthy_service",
			method:         http.MethodGet,
			healthService:  newUnhealthyService(errTestService),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   ErrHealthCheckFailed.Error(),
			expectJSON:     false,
			description:    "GET request with unhealthy service should return 500 Internal Server Error",
		},
		{
			name:           "GET_nil_service",
			method:         http.MethodGet,
			healthService:  nil,
			expectedStatus: 0, // Will panic before setting status
			expectedBody:   "",
			expectJSON:     false,
			description:    "GET request with nil service should panic (expected behavior)",
		},
	}
}

func generateMethodNotAllowedTests() []testCase {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodHead, http.MethodOptions}
	services := []struct {
		name    string
		service *testHealthService
	}{
		{"healthy", newHealthyService()},
		{"unhealthy", newUnhealthyService(errTestService)},
		{"nil", nil},
	}

	var cases []testCase
	for _, method := range methods {
		for _, svc := range services {
			cases = append(cases, testCase{
				name:           fmt.Sprintf("%s_%s_service", method, svc.name),
				method:         method,
				healthService:  svc.service,
				expectedStatus: http.StatusMethodNotAllowed,
				expectedBody:   ErrMethodNotAllowed.Error(),
				expectJSON:     false,
				description:    fmt.Sprintf("%s request should return 405 Method Not Allowed", method),
			})
		}
	}
	return cases
}

func generateTimeoutTests() []testCase {
	return []testCase{
		{
			name:           "GET_timeout_service_success",
			method:         http.MethodGet,
			healthService:  newTimeoutService(50 * time.Millisecond),
			expectedStatus: http.StatusOK,
			expectedBody:   expectedHealthyStatus,
			expectJSON:     true,
			description:    "GET request with timeout service that completes in time should return 200 OK",
			timeout:        100 * time.Millisecond,
		},
		{
			name:           "GET_timeout_service_failure",
			method:         http.MethodGet,
			healthService:  newTimeoutService(200 * time.Millisecond),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   ErrHealthCheckFailed.Error(),
			expectJSON:     false,
			description:    "GET request with timeout service that exceeds timeout should return 500 Internal Server Error",
			timeout:        100 * time.Millisecond,
		},
	}
}

func generateErrorScenarioTests() []testCase {
	errorScenarios := []struct {
		name        string
		err         error
		expectedErr string
	}{
		{"database_error", errors.New("database unavailable"), ErrHealthCheckFailed.Error()},
		{"timeout_error", context.DeadlineExceeded, ErrHealthCheckFailed.Error()},
		{"network_error", errors.New("network unreachable"), ErrHealthCheckFailed.Error()},
	}

	var cases []testCase
	for _, scenario := range errorScenarios {
		cases = append(cases, testCase{
			name:           fmt.Sprintf("GET_%s", scenario.name),
			method:         http.MethodGet,
			healthService:  newUnhealthyService(scenario.err),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   scenario.expectedErr,
			expectJSON:     false,
			description:    fmt.Sprintf("GET request with %s should return 500 Internal Server Error", scenario.name),
		})
	}
	return cases
}

// TestHealthCheckHandler tests all combinations of HTTP methods and service health states
func TestHealthCheckHandler(t *testing.T) {
	// Generate all test cases
	var allTestCases []testCase
	allTestCases = append(allTestCases, generateGETTests()...)
	allTestCases = append(allTestCases, generateMethodNotAllowedTests()...)
	allTestCases = append(allTestCases, generateTimeoutTests()...)
	allTestCases = append(allTestCases, generateErrorScenarioTests()...)

	// Run all generated test cases
	for _, tc := range allTestCases {
		t.Run(tc.name, func(t *testing.T) {
			runHealthCheckTest(t, tc)
		})
	}
}

// TestHealthCheckHandler_EdgeCases tests edge cases and error conditions
func TestHealthCheckHandler_EdgeCases(t *testing.T) {
	t.Run("nil_logger_uses_default", func(t *testing.T) {
		req := createTestRequest(http.MethodGet)
		rr := httptest.NewRecorder()
		handler := HealthCheckHandler(nil, newHealthyService())
		handler.ServeHTTP(rr, req)

		assertHTTPResponse(t, rr, http.StatusOK, expectedHealthyStatus, true)
	})

	t.Run("context_cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		service := &testHealthService{
			CheckReadinessFn: func(ctx context.Context) error {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					return nil
				}
			},
		}

		req := createTestRequestWithContext(http.MethodGet, testEndpoint, ctx)
		rr := httptest.NewRecorder()
		handler := HealthCheckHandler(newTestLogger(), service)
		handler.ServeHTTP(rr, req)

		// Should return 500 due to context cancellation
		assertHTTPResponse(t, rr, http.StatusInternalServerError, ErrHealthCheckFailed.Error(), false)
	})

	t.Run("json_encoding_validation", func(t *testing.T) {
		req := createTestRequest(http.MethodGet)
		rr := httptest.NewRecorder()
		handler := HealthCheckHandler(newTestLogger(), newHealthyService())
		handler.ServeHTTP(rr, req)

		// Verify the response is valid JSON
		var response HealthResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Response is not valid JSON: %v", err)
		}

		// Verify all required fields are present
		assertHealthResponse(t, response, expectedHealthyStatus)
	})
}

// TestHealthCheckHandler_ConcurrentAccess tests concurrent access to health check endpoint
func TestHealthCheckHandler_ConcurrentAccess(t *testing.T) {
	const numGoroutines = 100
	const numRequests = 10

	handler := HealthCheckHandler(newTestLogger(), newHealthyService())

	var wg sync.WaitGroup
	results := make(chan int, numGoroutines*numRequests)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numRequests; j++ {
				req := createTestRequest(http.MethodGet)
				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)
				results <- rr.Code
			}
		}()
	}

	wg.Wait()
	close(results)

	// Verify all requests returned 200 OK
	for statusCode := range results {
		if statusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", statusCode)
		}
	}
}

// TestHealthCheckHandler_ConcurrentAccessWithErrors tests concurrent access with mixed healthy/unhealthy services
func TestHealthCheckHandler_ConcurrentAccessWithErrors(t *testing.T) {
	const numGoroutines = 50

	// Create handlers with different service states
	healthyHandler := HealthCheckHandler(newTestLogger(), newHealthyService())
	unhealthyHandler := HealthCheckHandler(newTestLogger(), newUnhealthyService(errTestService))

	var wg sync.WaitGroup
	healthyResults := make(chan int, numGoroutines)
	unhealthyResults := make(chan int, numGoroutines)

	// Test healthy service concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := createTestRequest(http.MethodGet)
			rr := httptest.NewRecorder()
			healthyHandler.ServeHTTP(rr, req)
			healthyResults <- rr.Code
		}()
	}

	// Test unhealthy service concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := createTestRequest(http.MethodGet)
			rr := httptest.NewRecorder()
			unhealthyHandler.ServeHTTP(rr, req)
			unhealthyResults <- rr.Code
		}()
	}

	wg.Wait()
	close(healthyResults)
	close(unhealthyResults)

	// Verify healthy service results
	for statusCode := range healthyResults {
		if statusCode != http.StatusOK {
			t.Errorf("Expected status 200 for healthy service, got %d", statusCode)
		}
	}

	// Verify unhealthy service results
	for statusCode := range unhealthyResults {
		if statusCode != http.StatusInternalServerError {
			t.Errorf("Expected status 500 for unhealthy service, got %d", statusCode)
		}
	}
}

// BenchmarkHealthCheckHandler benchmarks the health check handler performance
func BenchmarkHealthCheckHandler(b *testing.B) {
	handler := HealthCheckHandler(newTestLogger(), newHealthyService())
	req := createTestRequest(http.MethodGet)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			b.Errorf("Expected status 200, got %d", rr.Code)
		}
	}
}

// BenchmarkHealthCheckHandler_Unhealthy benchmarks the health check handler with unhealthy service
func BenchmarkHealthCheckHandler_Unhealthy(b *testing.B) {
	handler := HealthCheckHandler(newTestLogger(), newUnhealthyService(errTestService))
	req := createTestRequest(http.MethodGet)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusInternalServerError {
			b.Errorf("Expected status 500, got %d", rr.Code)
		}
	}
}

// BenchmarkHealthCheckHandler_Concurrent benchmarks concurrent health check requests
func BenchmarkHealthCheckHandler_Concurrent(b *testing.B) {
	handler := HealthCheckHandler(newTestLogger(), newHealthyService())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		req := createTestRequest(http.MethodGet)
		for pb.Next() {
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				b.Errorf("Expected status 200, got %d", rr.Code)
			}
		}
	})
}
