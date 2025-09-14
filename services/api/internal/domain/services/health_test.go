package services

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"rdl-api/internal/db/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

// MockHealthRepository implements HealthRepository interface for testing
type MockHealthRepository struct {
	pingError error
	pingDelay time.Duration
}

func (m *MockHealthRepository) Ping(ctx context.Context) error {
	if m.pingDelay > 0 {
		select {
		case <-time.After(m.pingDelay):
			return m.pingError
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return m.pingError
}

// Test error definitions
var (
	ErrTestDatabaseConnectionFailed = errors.New("database connection failed")
	ErrTestDatabaseContextTimeout   = errors.New("context time out")
)

// Test helpers
func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Use Error level to reduce noise in tests
	}))
}

func newMockHealthRepository(pingError error, delay time.Duration) repository.HealthRepository {
	return &MockHealthRepository{
		pingError: pingError,
		pingDelay: delay,
	}
}

func newHealthyRepository() repository.HealthRepository {
	return newMockHealthRepository(nil, 0)
}

func newUnhealthyRepository(err error) repository.HealthRepository {
	if err == nil {
		err = ErrTestDatabaseConnectionFailed
	}
	return newMockHealthRepository(err, 0)
}

func newSlowRepository(delay time.Duration) repository.HealthRepository {
	return newMockHealthRepository(nil, delay)
}

// TestNewHealthService tests the constructor
func TestNewHealthService(t *testing.T) {
	tests := []struct {
		name     string
		pool     *pgxpool.Pool
		logger   *slog.Logger
		expected bool
	}{
		{
			name:     "valid_parameters",
			pool:     nil, // We'll use nil since we're not actually connecting to a real database
			logger:   newTestLogger(),
			expected: true,
		},
		{
			name:     "nil_logger",
			pool:     nil,
			logger:   nil,
			expected: true, // Should still work with nil logger
		},
		{
			name:     "nil_pool",
			pool:     nil,
			logger:   newTestLogger(),
			expected: true, // Constructor doesn't validate pool
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewHealthService(tt.pool, tt.logger)

			if service == nil {
				t.Error("Expected non-nil service")
				return
			}

			// Verify service implements HealthService interface
			var _ HealthService = service
		})
	}
}

// TestHealthService_CheckReadiness tests the CheckReadiness method
func TestHealthService_CheckReadiness(t *testing.T) {
	tests := []struct {
		name        string
		healthRepo  repository.HealthRepository
		expectError bool
		expectedErr error
		description string
	}{
		{
			name:        "database_available",
			healthRepo:  newHealthyRepository(),
			expectError: false,
			expectedErr: nil,
			description: "CheckReadiness should return nil when database is available",
		},
		{
			name:        "database_unavailable",
			healthRepo:  newUnhealthyRepository(repository.ErrDatabaseUnavailable),
			expectError: true,
			expectedErr: ErrDatabaseUnavailable,
			description: "CheckReadiness should return ErrDatabaseUnavailable when database ping fails",
		},
		{
			name:        "database_connection_failed",
			healthRepo:  newUnhealthyRepository(repository.ErrDatabaseConnection),
			expectError: true,
			expectedErr: ErrDatabaseUnavailable,
			description: "CheckReadiness should return ErrDatabaseUnavailable for any database error",
		},
		{
			name:        "generic_database_error",
			healthRepo:  newUnhealthyRepository(ErrTestDatabaseConnectionFailed),
			expectError: true,
			expectedErr: ErrDatabaseUnavailable,
			description: "CheckReadiness should return ErrDatabaseUnavailable for generic database errors",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create service with mock repository
			service := &healthService{
				healthRepo: tt.healthRepo,
				logger:     newTestLogger(),
			}

			ctx := context.Background()
			err := service.CheckReadiness(ctx)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestHealthService_CheckReadiness_Timeout tests timeout behavior
func TestHealthService_CheckReadiness_Timeout(t *testing.T) {
	tests := []struct {
		name           string
		pingDelay      time.Duration
		contextTimeout time.Duration
		expectError    bool
		description    string
	}{
		{
			name:           "ping_within_timeout",
			pingDelay:      1 * time.Second,
			contextTimeout: 2 * time.Second,
			expectError:    false,
			description:    "CheckReadiness should succeed when ping completes within timeout",
		},
		{
			name:           "ping_exceeds_timeout",
			pingDelay:      2 * time.Second,
			contextTimeout: 1 * time.Second,
			expectError:    true,
			description:    "CheckReadiness should fail when ping exceeds timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create service with slow repository
			service := &healthService{
				healthRepo: newSlowRepository(tt.pingDelay),
				logger:     newTestLogger(),
			}

			// Create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), tt.contextTimeout)
			defer cancel()

			err := service.CheckReadiness(ctx)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				// Should be either context deadline exceeded or timeout error
				if !errors.Is(err, context.DeadlineExceeded) && err != ErrDatabaseUnavailable {
					t.Errorf("Expected timeout or database unavailable error, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestHealthService_CheckReadiness_NilRepository tests nil repository handling
func TestHealthService_CheckReadiness_NilRepository(t *testing.T) {
	service := &healthService{
		healthRepo: nil,
		logger:     newTestLogger(),
	}

	ctx := context.Background()
	err := service.CheckReadiness(ctx)

	if err == nil {
		t.Error("Expected error for nil repository but got none")
		return
	}

	if !errors.Is(err, ErrDatabaseNotInitialized) {
		t.Errorf("Expected ErrDatabaseNotInitialized, got %v", err)
	}
}

// TestHealthService_CheckLiveness tests the CheckLiveness method
func TestHealthService_CheckLiveness(t *testing.T) {
	tests := []struct {
		name        string
		healthRepo  repository.HealthRepository
		description string
	}{
		{
			name:        "always_succeeds",
			healthRepo:  newHealthyRepository(),
			description: "CheckLiveness should always return nil regardless of repository state",
		},
		{
			name:        "nil_repository",
			healthRepo:  nil,
			description: "CheckLiveness should return nil even with nil repository",
		},
		{
			name:        "unhealthy_repository",
			healthRepo:  newUnhealthyRepository(ErrTestDatabaseConnectionFailed),
			description: "CheckLiveness should return nil even with unhealthy repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &healthService{
				healthRepo: tt.healthRepo,
				logger:     newTestLogger(),
			}

			ctx := context.Background()
			err := service.CheckLiveness(ctx)

			if err != nil {
				t.Errorf("CheckLiveness should always return nil, got error: %v", err)
			}
		})
	}
}

// TestHealthService_CheckLiveness_ContextCancellation tests context cancellation
func TestHealthService_CheckLiveness_ContextCancellation(t *testing.T) {
	service := &healthService{
		healthRepo: newHealthyRepository(),
		logger:     newTestLogger(),
	}

	// Test with canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := service.CheckLiveness(ctx)
	if err != nil {
		t.Errorf("CheckLiveness should return nil even with canceled context, got error: %v", err)
	}
}

// TestHealthService_EdgeCases tests edge cases and error conditions
func TestHealthService_EdgeCases(t *testing.T) {
	t.Run("context_with_deadline", func(t *testing.T) {
		service := &healthService{
			healthRepo: newHealthyRepository(),
			logger:     newTestLogger(),
		}

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour))
		defer cancel()

		err := service.CheckReadiness(ctx)
		if err != nil {
			t.Errorf("Expected no error with valid deadline, got: %v", err)
		}
	})

	t.Run("context_with_value", func(t *testing.T) {
		service := &healthService{
			healthRepo: newHealthyRepository(),
			logger:     newTestLogger(),
		}

		type testKey string
		const sampleKey testKey = "test-key"

		ctx := context.WithValue(context.Background(), sampleKey, "test-value")
		err := service.CheckReadiness(ctx)
		if err != nil {
			t.Errorf("Expected no error with context value, got: %v", err)
		}
	})

	t.Run("nil_logger_handling", func(t *testing.T) {
		service := &healthService{
			healthRepo: newHealthyRepository(),
			logger:     nil,
		}

		ctx := context.Background()
		err := service.CheckReadiness(ctx)
		if err != nil {
			t.Errorf("Expected no error with nil logger, got: %v", err)
		}
	})
}

// TestHealthService_ConcurrentAccess tests concurrent access to health service methods
func TestHealthService_ConcurrentAccess(t *testing.T) {
	const numGoroutines = 50
	const numRequests = 10

	service := &healthService{
		healthRepo: newHealthyRepository(),
		logger:     newTestLogger(),
	}

	var wg sync.WaitGroup
	results := make(chan error, numGoroutines*numRequests)

	// Test CheckReadiness concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numRequests; j++ {
				ctx := context.Background()
				err := service.CheckReadiness(ctx)
				results <- err
			}
		}()
	}

	wg.Wait()
	close(results)

	// Verify all requests succeeded
	for err := range results {
		if err != nil {
			t.Errorf("Expected no error in concurrent access, got: %v", err)
		}
	}
}

// TestHealthService_ConcurrentAccessWithErrors tests concurrent access with mixed repository states
func TestHealthService_ConcurrentAccessWithErrors(t *testing.T) {
	const numGoroutines = 25

	// Create services with different repository states
	healthyService := &healthService{
		healthRepo: newHealthyRepository(),
		logger:     newTestLogger(),
	}

	unhealthyService := &healthService{
		healthRepo: newUnhealthyRepository(ErrTestDatabaseConnectionFailed),
		logger:     newTestLogger(),
	}

	var wg sync.WaitGroup
	healthyResults := make(chan error, numGoroutines)
	unhealthyResults := make(chan error, numGoroutines)

	// Test healthy service concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()
			err := healthyService.CheckReadiness(ctx)
			healthyResults <- err
		}()
	}

	// Test unhealthy service concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()
			err := unhealthyService.CheckReadiness(ctx)
			unhealthyResults <- err
		}()
	}

	wg.Wait()
	close(healthyResults)
	close(unhealthyResults)

	// Verify healthy service results
	for err := range healthyResults {
		if err != nil {
			t.Errorf("Expected no error for healthy service, got: %v", err)
		}
	}

	// Verify unhealthy service results
	for err := range unhealthyResults {
		if err == nil {
			t.Error("Expected error for unhealthy service, got none")
		} else if !errors.Is(err, ErrDatabaseUnavailable) {
			t.Errorf("Expected ErrDatabaseUnavailable for unhealthy service, got: %v", err)
		}
	}
}

// BenchmarkHealthService_CheckReadiness benchmarks the CheckReadiness method
func BenchmarkHealthService_CheckReadiness(b *testing.B) {
	service := &healthService{
		healthRepo: newHealthyRepository(),
		logger:     newTestLogger(),
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := service.CheckReadiness(ctx)
		if err != nil {
			b.Errorf("Expected no error, got: %v", err)
		}
	}
}

// BenchmarkHealthService_CheckLiveness benchmarks the CheckLiveness method
func BenchmarkHealthService_CheckLiveness(b *testing.B) {
	service := &healthService{
		healthRepo: newHealthyRepository(),
		logger:     newTestLogger(),
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := service.CheckLiveness(ctx)
		if err != nil {
			b.Errorf("Expected no error, got: %v", err)
		}
	}
}

// BenchmarkHealthService_CheckReadiness_Unhealthy benchmarks CheckReadiness with unhealthy repository
func BenchmarkHealthService_CheckReadiness_Unhealthy(b *testing.B) {
	service := &healthService{
		healthRepo: newUnhealthyRepository(ErrTestDatabaseConnectionFailed),
		logger:     newTestLogger(),
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := service.CheckReadiness(ctx)
		if err == nil {
			b.Error("Expected error for unhealthy repository, got none")
		}
	}
}

// BenchmarkHealthService_Concurrent benchmarks concurrent access to health service
func BenchmarkHealthService_Concurrent(b *testing.B) {
	service := &healthService{
		healthRepo: newHealthyRepository(),
		logger:     newTestLogger(),
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		ctx := context.Background()
		for pb.Next() {
			err := service.CheckReadiness(ctx)
			if err != nil {
				b.Errorf("Expected no error, got: %v", err)
			}
		}
	})
}
