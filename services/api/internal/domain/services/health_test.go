package services

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

// mockHealthyRepository implements HealthRepository interface which is always healthy
type mockHealthyRepository struct{}

func (m *mockHealthyRepository) CheckReadiness(ctx context.Context) error {
	return nil
}

// mockUnhealthyRepository implements HealthRepository interface which is always unhealthy
type mockUnhealthyRepository struct{}

func (m *mockUnhealthyRepository) CheckReadiness(ctx context.Context) error {
	return ErrDatabaseUnavailable
}

// Test helpers
func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Use Error level to reduce noise in tests
	}))
}

// TestNewHealthService tests the constructor
func TestNewHealthService(t *testing.T) {
	tests := []struct {
		name        string
		pool        *pgxpool.Pool
		logger      *slog.Logger
		expectedErr error
	}{
		{
			name:        "nil_logger",
			pool:        nil,
			logger:      nil,
			expectedErr: ErrLoggerCannotBeNil, // Should still work with nil logger
		},
		{
			name:        "nil_pool",
			pool:        nil,
			logger:      newTestLogger(),
			expectedErr: ErrPoolCannotBeNil, // Constructor doesn't validate pool
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewHealthService(tt.pool, tt.logger, "test")

			assert.Equal(t, healthService{}, service)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

// TestHealthService_CheckReadiness tests the CheckReadiness method
func TestHealthService_CheckReadiness(t *testing.T) {
	tests := []struct {
		name        string
		healthRepo  HealthRepository
		expectError bool
		expectedErr error
		description string
	}{
		{
			name:        "Readiness succeeds",
			healthRepo:  &mockHealthyRepository{},
			expectError: false,
			expectedErr: nil,
			description: "CheckReadiness should return nil when database is available",
		},
		{
			name:        "Readiness fails",
			healthRepo:  &mockUnhealthyRepository{},
			expectError: true,
			expectedErr: ErrDatabaseUnavailable,
			description: "CheckReadiness should return ErrDatabaseUnavailable when database is not available",
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

// TestHealthService_CheckLiveness_ContextCancellation tests context cancellation
func TestHealthService_CheckLiveness_ContextCancellation(t *testing.T) {
	service := &healthService{
		healthRepo: &mockHealthyRepository{},
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
			healthRepo: &mockHealthyRepository{},
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
			healthRepo: &mockHealthyRepository{},
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
}

// TestHealthService_ConcurrentAccess tests concurrent access to health service methods
func TestHealthService_ConcurrentAccess(t *testing.T) {
	const numGoroutines = 50
	const numRequests = 10

	service := &healthService{
		healthRepo: &mockHealthyRepository{},
		logger:     newTestLogger(),
	}

	var wg sync.WaitGroup
	results := make(chan error, numGoroutines*numRequests)

	// Test CheckReadiness concurrently
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
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
		healthRepo: &mockHealthyRepository{},
		logger:     newTestLogger(),
	}

	unhealthyService := &healthService{
		healthRepo: &mockUnhealthyRepository{},
		logger:     newTestLogger(),
	}

	var wg sync.WaitGroup
	healthyResults := make(chan error, numGoroutines)
	unhealthyResults := make(chan error, numGoroutines)

	wg.Add(numGoroutines)
	// Test healthy service concurrently
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			ctx := context.Background()
			err := healthyService.CheckReadiness(ctx)
			healthyResults <- err
		}()
	}

	// Test unhealthy service concurrently
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
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
		healthRepo: &mockHealthyRepository{},
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
		healthRepo: &mockHealthyRepository{},
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
		healthRepo: &mockUnhealthyRepository{},
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
		healthRepo: &mockHealthyRepository{},
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
