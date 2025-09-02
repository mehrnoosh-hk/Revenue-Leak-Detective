package health

import (
	"context"
	"testing"
	"time"
)

// MockDatabaseChecker implements DatabaseChecker for testing
type MockDatabaseChecker struct {
	pingError error
}

func (m *MockDatabaseChecker) Ping(ctx context.Context) error {
	return m.pingError
}

func TestHealthService_CheckReadiness(t *testing.T) {
	tests := []struct {
		name        string
		dbChecker   DatabaseChecker
		expectError bool
	}{
		{
			name:        "database available",
			dbChecker:   &MockDatabaseChecker{pingError: nil},
			expectError: false,
		},
		{
			name:        "database unavailable",
			dbChecker:   &MockDatabaseChecker{pingError: ErrDatabaseUnavailable},
			expectError: true,
		},
		{
			name:        "database not initialized",
			dbChecker:   nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewHealthService(tt.dbChecker)
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
