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

// HTTPConfig holds HTTP server configuration for the API service
type HTTPConfig struct {
	Host string `yaml:"API_HOST"`
	Port string `yaml:"API_PORT"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL      string `yaml:"POSTGRES_URL"`
	Host     string `yaml:"POSTGRES_HOST"`
	Port     string `yaml:"POSTGRES_PORT"`
	User     string `yaml:"POSTGRES_USER"`
	Password string `yaml:"POSTGRES_PASSWORD"`
	DBName   string `yaml:"POSTGRES_DB"`
	SSLMode  string `yaml:"POSTGRES_SSL"`
}

// EnvironmentConfig holds environment-specific configuration
type EnvironmentConfig struct {
	LogLevel    slog.Level `yaml:"LOG_LEVEL"`
	Environment string     `yaml:"ENVIRONMENT"`
	ConfigVer   string     `yaml:"CONFIG_VERSION"`
	Debug       bool       `yaml:"DEBUG"`
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
	BuildInfo   BuildInfoConfig
	Database    DatabaseConfig
	HTTP        HTTPConfig
	Environment EnvironmentConfig
}

// LoadConfig loads the configuration with a specific env file path
func LoadConfig(envFilePath string) (*Config, error) {
	// Load environment file if specified (only in non-production environments)
	if envFilePath != "" {
		if err := loadEnvFile(envFilePath); err != nil {
			return nil, fmt.Errorf("failed to load environment file: %w", err)
		}
	}

	// After loading env file, determine the environment
	env := os.Getenv("ENVIRONMENT")
	isProduction := isProductionEnvironment(env)

	config := &Config{
		HTTP: HTTPConfig{
			Host: getEnvValue("API_HOST", isProduction, "0.0.0.0"),
			Port: getEnvValue("API_PORT", isProduction, "3030"),
		},
		Database: DatabaseConfig{
			URL:      os.Getenv("POSTGRES_URL"),
			Host:     getEnvValue("POSTGRES_HOST", isProduction, "localhost"),
			Port:     getEnvValue("POSTGRES_PORT", isProduction, "5432"),
			User:     getEnvValue("POSTGRES_USER", isProduction, "postgres"),
			Password: getEnvValue("POSTGRES_PASSWORD", isProduction, "password"),
			DBName:   getEnvValue("POSTGRES_DB", isProduction, "revenue_leak_detective_dev"),
			SSLMode:  getEnvValue("POSTGRES_SSL", isProduction, "disable"),
		},
		Environment: func() EnvironmentConfig {
			debugStr := getEnvValue("DEBUG", isProduction, "false")
			debugVal, _ := strconv.ParseBool(strings.ToLower(debugStr))
			return EnvironmentConfig{
				Environment: getEnvValue("ENVIRONMENT", isProduction, "development"),
				Debug:       debugVal,
				LogLevel:    parseLogLevel(getEnvValue("LOG_LEVEL", isProduction, "INFO")),
				ConfigVer:   getEnvValue("CONFIG_VERSION", isProduction, "unknown"),
			}
		}(),
		BuildInfo: BuildInfoConfig{
			GIT_COMMIT_HASH:       getEnvValue("GIT_COMMIT_HASH", isProduction, "unknown"),
			GIT_COMMIT_FULL:       getEnvValue("GIT_COMMIT_FULL", isProduction, "unknown"),
			GIT_COMMIT_DATE:       getEnvValue("GIT_COMMIT_DATE", isProduction, "unknown"),
			GIT_COMMIT_DATE_SHORT: getEnvValue("GIT_COMMIT_DATE_SHORT", isProduction, "unknown"),
			GIT_COMMIT_MESSAGE:    getEnvValue("GIT_COMMIT_MESSAGE", isProduction, "unknown"),
			GIT_BRANCH:            getEnvValue("GIT_BRANCH", isProduction, "unknown"),
			GIT_TAG:               getEnvValue("GIT_TAG", isProduction, "unknown"),
			GIT_DIRTY:             getEnvValue("GIT_DIRTY", isProduction, "unknown"),
			BUILD_TIMESTAMP:       getEnvValue("BUILD_TIMESTAMP", isProduction, "unknown"),
		},
	}

	// Validate required configuration
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// loadEnvFile loads environment file if path is provided
func loadEnvFile(envFilePath string) error {
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

// validate ensures all required configuration is present and valid
func (c *Config) validate() error {
	// Validate required environment variables in production first
	if err := c.validateRequiredEnvVars(); err != nil {
		return fmt.Errorf("required environment variables: %w", err)
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
			return fmt.Errorf("invalid POSTGRES_URL: %w", err)
		}
		return nil
	}

	// Otherwise, validate individual database parameters
	if c.Database.Host == "" {
		return fmt.Errorf("POSTGRES_HOST is required when POSTGRES_URL is not provided")
	}
	if c.Database.User == "" {
		return fmt.Errorf("POSTGRES_USER is required when POSTGRES_URL is not provided")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("POSTGRES_DB is required when POSTGRES_URL is not provided")
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

// validateRequiredEnvVars validates that required environment variables are set in production
func (c *Config) validateRequiredEnvVars() error {
	// Only validate in production environment
	if !c.IsProduction() {
		return nil
	}

	// In production, all required environment variables must be set
	requiredVars := []string{
		"API_HOST", "API_PORT", "POSTGRES_HOST", "POSTGRES_PORT",
		"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB", "POSTGRES_SSL",
		"ENVIRONMENT", "DEBUG", "LOG_LEVEL", "CONFIG_VERSION",
	}

	var missing []string
	for _, varName := range requiredVars {
		if value, exists := os.LookupEnv(varName); !exists || value == "" {
			missing = append(missing, varName)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables in production: %v", missing)
	}

	return nil
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

// PrintEffectiveConfig prints the effective configuration (excluding secrets and build information)
func (c *Config) PrintEffectiveConfig(logger *slog.Logger) {

	// Use slog to print the configuration
	logger.Info("Effective configuration loaded")
	logger.Info(fmt.Sprintf("config_version: %s", c.Environment.ConfigVer))
	logger.Info(fmt.Sprintf("environment: %s", c.Environment.Environment))
	logger.Info(fmt.Sprintf("debug: %v", c.Environment.Debug))
	logger.Info(fmt.Sprintf("log_level: %s", c.Environment.LogLevel.String()))
	logger.Info(fmt.Sprintf("http_port: %s", c.HTTP.Port))
	logger.Info(fmt.Sprintf("db_host: %s", c.Database.Host))
	logger.Info(fmt.Sprintf("db_port: %s", c.Database.Port))
	logger.Info(fmt.Sprintf("db_name: %s", c.Database.DBName))
	logger.Info(fmt.Sprintf("db_user: %s", c.Database.User))
	logger.Info(fmt.Sprintf("db_ssl_mode: %s", c.Database.SSLMode))
}

// PrintBuildInfo prints the build information
func (c *Config) PrintBuildInfo(logger *slog.Logger) {
	logger.Info("Build information")
	logger.Info(fmt.Sprintf("version: %s", c.BuildInfo.GIT_TAG))
	logger.Info(fmt.Sprintf("commit: %s", c.BuildInfo.GIT_COMMIT_FULL))
	logger.Info(fmt.Sprintf("build_date: %s", c.BuildInfo.BUILD_TIMESTAMP))
	logger.Info(fmt.Sprintf("git_branch: %s", c.BuildInfo.GIT_BRANCH))
	logger.Info(fmt.Sprintf("git_tag: %s", c.BuildInfo.GIT_TAG))
	logger.Info(fmt.Sprintf("git_dirty: %s", c.BuildInfo.GIT_DIRTY))
	logger.Info(fmt.Sprintf("git_commit_message: %s", c.BuildInfo.GIT_COMMIT_MESSAGE))
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
	// If POSTGRES_URL is provided, use it directly
	if c.Database.URL != "" {
		return c.Database.URL
	}

	// Otherwise, construct from individual parameters
	u := &url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(c.Database.User, c.Database.Password),
		Host:   net.JoinHostPort(c.Database.Host, c.Database.Port),
		Path:   "/" + c.Database.DBName,
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
