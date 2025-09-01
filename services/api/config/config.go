package config

import (
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	CONFIG_VERSION = "1.0.0"
)

// HTTPConfig holds HTTP server configuration for the API service
type HTTPConfig struct {
	Host string `yaml:"api_host"`
	Port string `yaml:"api_port"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL      string `yaml:"POSTGRESQL_URL"`
	Host     string `yaml:"POSTGRESQL_HOST"`
	Port     string `yaml:"POSTGRESQL_PORT"`
	User     string `yaml:"POSTGRESQL_USER"`
	Password string `yaml:"POSTGRESQL_PASSWORD"`
	Name     string `yaml:"POSTGRESQL_NAME"`
	SSLMode  string `yaml:"POSTGRES_SSL"`
}

// EnvironmentConfig holds environment-specific configuration
type EnvironmentConfig struct {
	Environment string     `yaml:"ENVIRONMENT"`
	Debug       bool       `yaml:"DEBUG"`
	LogLevel    slog.Level `yaml:"LOG_LEVEL"`
	ConfigVer   string     `yaml:"CONFIG_VERSION"`
}

// BuildInfoConfig holds build information configuration
type BuildInfoConfig struct {
	GIT_COMMIT_HASH       string `yaml:"GIT_COMMIT_HASH"`
	GIT_COMMIT_FULL       string `yaml:"GIT_COMMIT_FULL"`
	GIT_COMMIT_DATE       string `yaml:"GIT_COMMIT_DATE"`
	GIT_COMMIT_DATE_SHORT string `yaml:"GIT_COMMIT_DATE_SHORT"`
	GIT_COMMIT_MESSAGE    string `yaml:"GIT_COMMIT_MESSAGE"`
	GIT_BRANCH            string `yaml:"GIT_BRANCH"`
	GIT_TAG               string `yaml:"GIT_TAG"`
	GIT_DIRTY             string `yaml:"GIT_DIRTY"`
	BUILD_TIMESTAMP       string `yaml:"BUILD_TIMESTAMP"`
}

// Config holds the complete application configuration
type Config struct {
	HTTP        HTTPConfig
	Database    DatabaseConfig
	Environment EnvironmentConfig
	BuildInfo   BuildInfoConfig
}


// LoadConfigWithEnvFile loads the configuration with a specific env file path
func LoadConfigWithEnvFile(envFilePath string) (*Config, error) {
	// Load environment file if specified
	if err := loadEnvFile(envFilePath); err != nil {
		return nil, fmt.Errorf("failed to load environment file: %w", err)
	}

	config := &Config{
		HTTP: HTTPConfig{
			Host: getEnvOrDefault("API_HOST", "0.0.0.0"),
			Port: getEnvOrDefault("API_PORT", "3030"),
		},
		Database: DatabaseConfig{
			URL:      os.Getenv("POSTGRES_URL"),
			Host:     getEnvOrDefaultStrict("POSTGRES_HOST", "localhost"),
			Port:     getEnvOrDefaultStrict("POSTGRES_PORT", "5432"),
			User:     getEnvOrDefaultStrict("POSTGRES_USER", "postgres"),
			Password: getEnvOrDefaultStrict("POSTGRES_PASSWORD", "password"),
			Name:     getEnvOrDefaultStrict("POSTGRES_NAME", "revenue_leak_detective_dev"),
			SSLMode:  getEnvOrDefaultStrict("POSTGRES_SSL", "disable"),
		},
		Environment: EnvironmentConfig{
			Environment: getEnvOrDefault("ENVIRONMENT", "development"),
			Debug:       getEnvOrDefault("DEBUG", "false") == "true",
			LogLevel:    parseLogLevel(getEnvOrDefault("LOG_LEVEL", "INFO")),
			ConfigVer:   CONFIG_VERSION,
		},
	}

	// Validate required configuration
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Print effective configuration (excluding secrets)
	config.printEffectiveConfig()

	return config, nil
}

// loadEnvFile loads environment file if path is provided
func loadEnvFile(envFilePath string) error {
	// Only load env file if path is explicitly provided
	if envFilePath == "" {
		return nil
	}

	// Check if file exists
	if _, err := os.Stat(envFilePath); err != nil {
		return fmt.Errorf("env file not found: %s", envFilePath)
	}

	// Load the env file
	if err := godotenv.Load(envFilePath); err != nil {
		return fmt.Errorf("failed to load %s: %w", envFilePath, err)
	}

	return nil
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvOrDefaultStrict gets an environment variable or returns a default value
// This version treats empty strings as "not set" and returns the default
func getEnvOrDefaultStrict(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

// validate ensures all required configuration is present and valid
func (c *Config) validate() error {
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
	// If DATABASE_URL is provided, it takes precedence
	if c.Database.URL != "" {
		if _, err := url.Parse(c.Database.URL); err != nil {
			return fmt.Errorf("invalid DATABASE_URL: %w", err)
		}
		return nil
	}

	// Otherwise, validate individual database parameters
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required when DATABASE_URL is not provided")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required when DATABASE_URL is not provided")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required when DATABASE_URL is not provided")
	}

	if err := validatePort(c.Database.Port); err != nil {
		return fmt.Errorf("invalid database port: %w", err)
	}

	return nil
}

// validateEnvironment validates environment configuration
func (c *Config) validateEnvironment() error {
	validEnvs := []string{"development", "dev", "staging", "production", "prod", "test"}
	env := strings.ToLower(c.Environment.Environment)

	for _, valid := range validEnvs {
		if env == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid environment: %s (valid: %v)", env, validEnvs)
}

// validatePort ensures the port is valid
func validatePort(port string) error {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("port must be a number: %s", port)
	}
	if portNum < 1 || portNum > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got: %d", portNum)
	}
	return nil
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

// printEffectiveConfig prints the effective configuration (excluding secrets)
func (c *Config) printEffectiveConfig() {

	// Use slog to print the configuration
	logger := slog.Default()
	logger.Info("Effective configuration loaded",
		slog.String("config_version", CONFIG_VERSION),
		slog.String("environment", c.Environment.Environment),
		slog.Bool("debug", c.Environment.Debug),
		slog.String("log_level", c.Environment.LogLevel.String()),
		slog.String("http_port", c.HTTP.Port),
		slog.String("db_host", c.Database.Host),
		slog.String("db_port", c.Database.Port),
		slog.String("db_name", c.Database.Name),
		slog.String("db_user", c.Database.User),
		slog.String("db_ssl_mode", c.Database.SSLMode),
	)
}

// maskSecret masks sensitive parts of a URL
func maskSecret(urlStr string) string {
	if urlStr == "" {
		return ""
	}

	parsed, err := url.Parse(urlStr)
	if err != nil {
		return "***"
	}

	if parsed.User != nil {
		if _, ok := parsed.User.Password(); ok {
			parsed.User = url.UserPassword(parsed.User.Username(), "***")
		}
	}

	return parsed.String()
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

// DatabaseURL returns the database connection URL
func (c *Config) DatabaseURL() string {
	// If DATABASE_URL is provided, use it directly
	if c.Database.URL != "" {
		return c.Database.URL
	}

	// Otherwise, construct from individual parameters
	u := &url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(c.Database.User, c.Database.Password),
		Host:   net.JoinHostPort(c.Database.Host, c.Database.Port),
		Path:   "/" + c.Database.Name,
	}

	// Add SSL mode as query parameter if specified
	if c.Database.SSLMode != "" {
		q := u.Query()
		q.Set("sslmode", c.Database.SSLMode)
		u.RawQuery = q.Encode()
	}

	return u.String()
}

// GetLogLevel returns the configured log level
func (c *Config) GetLogLevel() slog.Level {
	return c.Environment.LogLevel
}

// GetPort returns the configured HTTP port
func (c *Config) GetPort() string {
	return c.HTTP.Port
}

// GetEnvironment returns the configured environment
func (c *Config) GetEnvironment() string {
	return c.Environment.Environment
}
