package config

import (
	"log/slog"
	"os"
	"strings"
)

// isProductionEnvironment checks if the given environment is production
func isProductionEnvironment(env string) bool {
	env = strings.ToLower(env)
	return env == "production" || env == "prod"
}

// getEnvValue gets an environment variable with different behavior based on environment
// In production: requires the environment variable to be set, returns error if not found
// In non-production: falls back to default value if not set
func getEnvValue(key string, isProduction bool, defaultValue string) string {
	if isProduction {
		// In production, require the environment variable to be set
		if value, exists := os.LookupEnv(key); exists && value != "" {
			return value
		}
		// If not found in production, return empty string (will be caught by validation)
		return ""
	}

	// In non-production, use default behavior
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseLogLevel converts string log level to slog.Level
func parseLogLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	env := strings.ToLower(c.Environment.Environment)
	return env == "development" || env == "dev"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	env := strings.ToLower(c.Environment.Environment)
	return env == "production" || env == "prod"
}

// IsStaging returns true if the environment is staging
func (c *Config) IsStaging() bool {
	env := strings.ToLower(c.Environment.Environment)
	return env == "staging"
}

// IsTest returns true if the environment is test
func (c *Config) IsTest() bool {
	env := strings.ToLower(c.Environment.Environment)
	return env == "test"
}

// GetEnvironment returns the configured environment
func (c *Config) GetEnvironment() string {
	return c.Environment.Environment
}

// GetLogLevel returns the configured log level
func (c *Config) GetLogLevel() slog.Level {
	return c.Environment.LogLevel
}
