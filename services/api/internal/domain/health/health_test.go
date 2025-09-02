package health

import (
	"context"
	"testing"
	"time"

	"rdl-api/internal/db/repository"
)

// MockHealthRepository implements repository.HealthRepository for testing
type MockHealthRepository struct {
	pingError error
}

func (m *MockHealthRepository) Ping(ctx context.Context) error {
	return m.pingError
}

func TestHealthService_CheckReadiness(t *testing.T) {
	tests := []struct {
		name        string
		healthRepo  repository.HealthRepository
		expectError bool
	}{
		{
			name:        "database available",
			healthRepo:  &MockHealthRepository{pingError: nil},
			expectError: false,
		},
		{
			name:        "database unavailable",
			healthRepo:  &MockHealthRepository{pingError: repository.ErrDatabaseUnavailable},
			expectError: true,
		},
		{
			name:        "database not initialized",
			healthRepo:  nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewHealthService(tt.healthRepo)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := service.CheckReadiness(ctx)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestHealthService_CheckLiveness(t *testing.T) {
	service := NewHealthService(nil)
	ctx := context.Background()

	err := service.CheckLiveness(ctx)

	if err != nil {
		t.Errorf("liveness check should never fail: %v", err)
	}
}
