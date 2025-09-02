package repository

import (
	"context"
	"testing"
	"time"
)

// MockDatabase implements Database interface for testing
type MockDatabase struct {
	pingError error
}

func (m *MockDatabase) Ping(ctx context.Context) error {
	return m.pingError
}

func TestHealthRepository_Ping(t *testing.T) {
	tests := []struct {
		name        string
		db          Database
		expectError bool
	}{
		{
			name:        "database available",
			db:          &MockDatabase{pingError: nil},
			expectError: false,
		},
		{
			name:        "database unavailable",
			db:          &MockDatabase{pingError: ErrDatabaseUnavailable},
			expectError: true,
		},
		{
			name:        "database connection failed",
			db:          &MockDatabase{pingError: ErrDatabaseConnection},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewHealthRepository(tt.db)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := repo.Ping(ctx)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
