package repository

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"
)

// MockPool implements PoolPinger interface for testing
type MockPool struct {
	pingError error
}

// newTestHealthRepositoryWithPinger creates a new health repository instance with a custom pinger
// This is useful for testing with mock implementations
func newTestHealthRepositoryWithPinger(pool PoolPinger, logger *slog.Logger) HealthRepository {
	return &healthRepository{
		pool:   pool,
		logger: logger,
	}
}

// Ping implements the Ping method for testing
func (m *MockPool) Ping(ctx context.Context) error {
	// Check if context is canceled or timed out
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return m.pingError
	}
}

func TestHealthRepository_Ping(t *testing.T) {
	// Create a test logger
	logger := slog.Default()

	tests := []struct {
		name        string
		pool        *MockPool
		expectError bool
		expectedErr error
	}{
		{
			name:        "database available",
			pool:        &MockPool{pingError: nil},
			expectError: false,
			expectedErr: nil,
		},
		{
			name:        "database unavailable",
			pool:        &MockPool{pingError: ErrDatabaseUnavailable},
			expectError: true,
			expectedErr: ErrDatabaseUnavailable,
		},
		{
			name:        "database connection failed",
			pool:        &MockPool{pingError: ErrDatabaseConnection},
			expectError: true,
			expectedErr: ErrDatabaseConnection,
		},
		{
			name:        "context timeout",
			pool:        &MockPool{pingError: context.DeadlineExceeded},
			expectError: true,
			expectedErr: context.DeadlineExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create repository with mock pool
			repo := newTestHealthRepositoryWithPinger(tt.pool, logger)
			// Create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Execute ping
			err := repo.Ping(ctx)

			// Verify error expectations
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestHealthRepository_Ping_ContextCancellation(t *testing.T) {
	logger := slog.Default()
	pool := &MockPool{pingError: nil}
	repo := newTestHealthRepositoryWithPinger(pool, logger)

	// Create a context that's already canceled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Execute ping with canceled context
	err := repo.Ping(ctx)

	// Should return context canceled error
	if err == nil {
		t.Error("expected context canceled error but got none")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}

func TestHealthRepository_Ping_ContextTimeout(t *testing.T) {
	logger := slog.Default()
	pool := &MockPool{pingError: context.DeadlineExceeded}
	repo := newTestHealthRepositoryWithPinger(pool, logger)

	// Create a context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Execute ping
	err := repo.Ping(ctx)

	// Should return timeout error
	if err == nil {
		t.Error("expected timeout error but got none")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded error, got %v", err)
	}
}

func TestNewHealthRepository(t *testing.T) {
	logger := slog.Default()
	pool := &MockPool{pingError: nil}

	// Test repository creation
	repo := newTestHealthRepositoryWithPinger(pool, logger)

	// Verify repository is not nil
	if repo == nil {
		t.Error("expected repository to be created, got nil")
	}

	// Verify repository implements HealthRepository interface
	var _ = repo
}
