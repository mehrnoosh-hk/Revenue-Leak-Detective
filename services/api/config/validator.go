package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// validate ensures all required configuration is present and valid
func (c *Config) validate() error {
	// Validate required environment variables in production first
	if err := c.validateRequiredEnvVars(); err != nil {
		return fmt.Errorf("%s: %w", ErrMissingRequiredEnvVar, err)
	}

	// Validate HTTP configuration
	if err := c.validateHTTP(); err != nil {
		return fmt.Errorf("HTTP config: %w", err)
	}

	// Validate database configuration
	if err := c.validateDatabase(); err != nil {
		return fmt.Errorf("database config: %w", err)
	}

	// Validate environment configuration
	if err := c.validateEnvironment(); err != nil {
		return fmt.Errorf("environment config: %w", err)
	}

	return nil
}

// validateHTTP validates HTTP server configuration
func (c *Config) validateHTTP() error {
	if err := validatePort(c.HTTP.Port); err != nil {
		return fmt.Errorf("invalid port: %w", err)
	}
	return nil
}

// validateDatabase validates database configuration
func (c *Config) validateDatabase() error {
	// If POSTGRES_URL is provided, it takes precedence
	if c.Database.URL != "" {
		if _, err := url.Parse(c.Database.URL); err != nil {
			return fmt.Errorf("%s: %w", ErrInvalidDBURL, err)
		}
		return nil
	}

	// Otherwise, validate individual database parameters
	if c.Database.Host == "" {
		return fmt.Errorf("%s: POSTGRES_HOST is required when POSTGRES_URL is not provided", ErrMissingDBHost)
	}
	if c.Database.User == "" {
		return fmt.Errorf("%s: POSTGRES_USER is required when POSTGRES_URL is not provided", ErrMissingDBUser)
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("%s: POSTGRES_DB is required when POSTGRES_URL is not provided", ErrMissingDBName)
	}

	if err := validatePort(c.Database.Port); err != nil {
		return fmt.Errorf("database port: %w", err)
	}

	return nil
}

// validateEnvironment validates environment configuration
func (c *Config) validateEnvironment() error {
	env := strings.ToLower(c.Environment.Environment)

	for _, valid := range ValidEnvironments {
		if env == valid {
			return nil
		}
	}

	return fmt.Errorf("%s: %s (valid: %v)", ErrInvalidEnvironment, env, ValidEnvironments)
}

// validateRequiredEnvVars validates that required environment variables are set in production
func (c *Config) validateRequiredEnvVars() error {
	// Only validate in production environment
	if !c.IsProduction() {
		return nil
	}

	// In production, all required environment variables must be set
	requiredVars := []string{
		EnvAPIHost, EnvAPIPort, EnvPostgresHost, EnvPostgresPort,
		EnvPostgresUser, EnvPostgresPassword, EnvPostgresDB, EnvPostgresSSL,
		EnvEnvironment, EnvDebug, EnvLogLevel, EnvConfigVer,
	}

	var missing []string
	for _, varName := range requiredVars {
		if value, exists := os.LookupEnv(varName); !exists || value == "" {
			missing = append(missing, varName)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("%s in production: %v", ErrMissingRequiredEnvVar, missing)
	}

	return nil
}

// validatePort ensures the port is valid
func validatePort(port string) error {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("%s: %s (must be a number)", ErrInvalidPort, port)
	}
	if portNum < 1 || portNum > 65535 {
		return fmt.Errorf("%s: %d (must be between 1 and 65535)", ErrPortOutOfRange, portNum)
	}
	return nil
}
