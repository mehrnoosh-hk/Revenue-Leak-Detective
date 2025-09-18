package repository

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

// MockPool implements PoolPinger interface for testing
type mockHealthyPinger struct{}

func (m *mockHealthyPinger) Ping(ctx context.Context) error {
	return nil
}

type mockUnhealthyPinger struct{}

func (m *mockUnhealthyPinger) Ping(ctx context.Context) error {
	return ErrDatabaseUnavailable
}

func TestHealthRepository_CheckReadiness(t *testing.T) {
	// Create a test logger
	tests := []struct {
		name        string
		poolPinger  PoolPinger
		expectError bool
		expectedErr error
	}{
		{
			name:        "database available",
			poolPinger:  &mockHealthyPinger{},
			expectError: false,
			expectedErr: nil,
		},
		{
			name:        "database unavailable",
			poolPinger:  &mockUnhealthyPinger{},
			expectError: true,
			expectedErr: ErrDatabaseUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Execute ping
			err := tt.poolPinger.Ping(ctx)

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

func TestNewHealthRepository(t *testing.T) {
	tests := []struct {
		name        string
		pool        *pgxpool.Pool
		logger      *slog.Logger
		expectedErr error
	}{
		{
			name:        "nil_pool",
			pool:        nil,
			logger:      slog.Default(),
			expectedErr: ErrPoolCannotBeNil,
		},
		{
			name:        "nil_logger",
			pool:        nil,
			logger:      nil,
			expectedErr: ErrLoggerCannotBeNil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := NewHealthRepositoryImplementation(tt.pool, tt.logger)
			assert.Equal(t, HealthRepositoryImplementation{}, repo)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
