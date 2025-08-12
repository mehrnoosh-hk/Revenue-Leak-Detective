package config

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// Store original env vars to restore later
	originalPort := os.Getenv("API_PORT")
	originalLogLevel := os.Getenv("LOG_LEVEL")
	originalEnv := os.Getenv("ENVIRONMENT")

	// Clean up after test
	defer func() {
		os.Setenv("API_PORT", originalPort)
		os.Setenv("LOG_LEVEL", originalLogLevel)
		os.Setenv("ENVIRONMENT", originalEnv)
	}()

	t.Run("default configuration", func(t *testing.T) {
		// Clear env vars
		os.Unsetenv("API_PORT")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("ENVIRONMENT")

		cfg, err := LoadConfig()
		require.NoError(t, err)

		assert.Equal(t, "8080", cfg.Port)
		assert.Equal(t, slog.LevelInfo, cfg.LogLevel)
		assert.Equal(t, "development", cfg.Env)
	})

	t.Run("custom configuration from env vars", func(t *testing.T) {
		os.Setenv("API_PORT", "3000")
		os.Setenv("LOG_LEVEL", "DEBUG")
		os.Setenv("ENVIRONMENT", "production")

		cfg, err := LoadConfig()
		require.NoError(t, err)

		assert.Equal(t, "3000", cfg.Port)
		assert.Equal(t, slog.LevelDebug, cfg.LogLevel)
		assert.Equal(t, "production", cfg.Env)
	})

	t.Run("invalid port", func(t *testing.T) {
		os.Setenv("API_PORT", "invalid")

		_, err := LoadConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid port")
	})

	t.Run("port out of range", func(t *testing.T) {
		os.Setenv("API_PORT", "70000")

		_, err := LoadConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "port must be between 1 and 65535")
	})

	t.Run("invalid log level", func(t *testing.T) {
		os.Setenv("API_PORT", "8080")
		os.Setenv("LOG_LEVEL", "INVALID")

		_, err := LoadConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid log level")
	})
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name    string
		port    string
		wantErr bool
	}{
		{"valid port", "8080", false},
		{"minimum port", "1", false},
		{"maximum port", "65535", false},
		{"invalid non-numeric", "abc", true},
		{"port too low", "0", true},
		{"port too high", "65536", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePort(tt.port)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		expected slog.Level
		wantErr  bool
	}{
		{"debug", "DEBUG", slog.LevelDebug, false},
		{"info", "INFO", slog.LevelInfo, false},
		{"warn", "WARN", slog.LevelWarn, false},
		{"warning", "WARNING", slog.LevelWarn, false},
		{"error", "ERROR", slog.LevelError, false},
		{"lowercase debug", "debug", slog.LevelDebug, false},
		{"mixed case", "Info", slog.LevelInfo, false},
		{"invalid", "INVALID", slog.LevelInfo, true},
		{"empty", "", slog.LevelInfo, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, err := parseLogLevel(tt.level)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, level)
			}
		})
	}
}

func TestConfigEnvironmentMethods(t *testing.T) {
	tests := []struct {
		name          string
		env           string
		isDevelopment bool
		isProduction  bool
	}{
		{"development", "development", true, false},
		{"dev", "dev", true, false},
		{"production", "production", false, true},
		{"prod", "prod", false, true},
		{"staging", "staging", false, false},
		{"test", "test", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Env: tt.env}
			assert.Equal(t, tt.isDevelopment, cfg.IsDevelopment())
			assert.Equal(t, tt.isProduction, cfg.IsProduction())
		})
	}
}
