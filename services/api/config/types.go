package config

import (
	"log/slog"
)

// HTTPConfig holds HTTP server configuration for the API service
type HTTPConfig struct {
	// Host is the HTTP server host address to bind to
	// Default: "0.0.0.0" (all interfaces)
	// Environment variable: API_HOST
	Host string `yaml:"API_HOST" json:"host" example:"0.0.0.0" validate:"required"`

	// Port is the HTTP server port to listen on
	// Must be between 1-65535
	// Default: "3030"
	// Environment variable: API_PORT
	Port string `yaml:"API_PORT" json:"port" example:"3030" validate:"required"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	// URL is the complete database connection URL
	// Takes precedence over individual parameters if provided
	// Format: postgresql://user:password@host:port/dbname?sslmode=mode
	// Environment variable: POSTGRES_URL
	URL string `yaml:"POSTGRES_URL" json:"url" example:"postgresql://user:pass@localhost:5432/dbname"`

	// Host is the database server hostname or IP address
	// Required if URL is not provided
	// Default: "localhost"
	// Environment variable: POSTGRES_HOST
	Host string `yaml:"POSTGRES_HOST" json:"host" example:"localhost" validate:"required_without=URL"`

	// Port is the database server port
	// Must be between 1-65535
	// Default: "5432"
	// Environment variable: POSTGRES_PORT
	Port string `yaml:"POSTGRES_PORT" json:"port" example:"5432" validate:"required_without=URL"`

	// User is the database username
	// Required if URL is not provided
	// Default: "postgres"
	// Environment variable: POSTGRES_USER
	User string `yaml:"POSTGRES_USER" json:"user" example:"postgres" validate:"required_without=URL"`

	// Password is the database password
	// Default: "password"
	// Environment variable: POSTGRES_PASSWORD
	Password string `yaml:"POSTGRES_PASSWORD" json:"password" example:"password"`

	// DBName is the database name
	// Required if URL is not provided
	// Default: "revenue_leak_detective_dev"
	// Environment variable: POSTGRES_DB
	DBName string `yaml:"POSTGRES_DB" json:"db_name" example:"revenue_leak_detective_dev" validate:"required_without=URL"`

	// SSLMode is the SSL connection mode
	// Options: disable, require, verify-ca, verify-full
	// Default: "disable"
	// Environment variable: POSTGRES_SSL
	SSLMode string `yaml:"POSTGRES_SSL" json:"ssl_mode" example:"disable" validate:"oneof=disable require verify-ca verify-full"`
}

// EnvironmentConfig holds environment-specific configuration
type EnvironmentConfig struct {
	// LogLevel is the logging level for the application
	// Options: DEBUG, INFO, WARN, ERROR
	// Default: "INFO"
	// Environment variable: LOG_LEVEL
	LogLevel slog.Level `yaml:"LOG_LEVEL" json:"log_level" example:"INFO" validate:"oneof=DEBUG INFO WARN ERROR"`

	// Environment is the application environment
	// Options: development, dev, staging, production, prod, test
	// Default: "development"
	// Environment variable: ENVIRONMENT
	Environment string `yaml:"ENVIRONMENT" json:"environment" example:"development" validate:"oneof=development dev staging production prod test"`

	// ConfigVer is the configuration schema version
	// Used for configuration migration and validation
	// Default: "unknown"
	// Environment variable: CONFIG_VERSION
	ConfigVer string `yaml:"CONFIG_VERSION" json:"config_version" example:"1.0.0"`

	// Debug enables debug mode for additional logging and features
	// Default: false
	// Environment variable: DEBUG
	Debug bool `yaml:"DEBUG" json:"debug" example:"false"`
}

// BuildInfoConfig holds build information configuration
type BuildInfoConfig struct {
	// GIT_COMMIT_HASH is the short git commit hash
	// Example: "a1b2c3d"
	GIT_COMMIT_HASH string `yaml:"GIT_COMMIT_HASH" json:"git_commit_hash" example:"a1b2c3d"`

	// GIT_COMMIT_FULL is the full git commit hash
	// Example: "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0"
	GIT_COMMIT_FULL string `yaml:"GIT_COMMIT_FULL" json:"git_commit_full" example:"a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0"`

	// GIT_COMMIT_DATE is the git commit date in full format
	// Example: "2024-01-15T10:30:00Z"
	GIT_COMMIT_DATE string `yaml:"GIT_COMMIT_DATE" json:"git_commit_date" example:"2024-01-15T10:30:00Z"`

	// GIT_COMMIT_DATE_SHORT is the git commit date in short format
	// Example: "2024-01-15"
	GIT_COMMIT_DATE_SHORT string `yaml:"GIT_COMMIT_DATE_SHORT" json:"git_commit_date_short" example:"2024-01-15"`

	// GIT_COMMIT_MESSAGE is the git commit message
	// Example: "feat: add new configuration system"
	GIT_COMMIT_MESSAGE string `yaml:"GIT_COMMIT_MESSAGE" json:"git_commit_message" example:"feat: add new configuration system"`

	// GIT_BRANCH is the git branch name
	// Example: "main"
	GIT_BRANCH string `yaml:"GIT_BRANCH" json:"git_branch" example:"main"`

	// GIT_TAG is the git tag (if any)
	// Example: "v1.0.0"
	GIT_TAG string `yaml:"GIT_TAG" json:"git_tag" example:"v1.0.0"`

	// GIT_DIRTY indicates if the working directory has uncommitted changes
	// Example: "true" or "false"
	GIT_DIRTY string `yaml:"GIT_DIRTY" json:"git_dirty" example:"false"`

	// BUILD_TIMESTAMP is the build timestamp
	// Example: "2024-01-15T10:30:00Z"
	BUILD_TIMESTAMP string `yaml:"BUILD_TIMESTAMP" json:"build_timestamp" example:"2024-01-15T10:30:00Z"`
}

// Config holds the complete application configuration
// This is the main configuration struct that combines all configuration aspects
type Config struct {
	// BuildInfo contains build-time information like git commit, version, etc.
	BuildInfo BuildInfoConfig `json:"build_info" yaml:"build_info"`

	// Database contains database connection configuration
	Database DatabaseConfig `json:"database" yaml:"database"`

	// HTTP contains HTTP server configuration
	HTTP HTTPConfig `json:"http" yaml:"http"`

	// Environment contains environment-specific configuration
	Environment EnvironmentConfig `json:"environment" yaml:"environment"`
}

// Valid environments
var ValidEnvironments = []string{"development", "dev", "staging", "production", "prod", "test"}

// Valid log levels
var ValidLogLevels = map[string]slog.Level{
	"DEBUG":   slog.LevelDebug,
	"INFO":    slog.LevelInfo,
	"WARN":    slog.LevelWarn,
	"WARNING": slog.LevelWarn,
	"ERROR":   slog.LevelError,
}

// Default configuration values
const (
	DefaultAPIPort     = "3030"
	DefaultAPIHost     = "0.0.0.0"
	DefaultDBHost      = "localhost"
	DefaultDBPort      = "5432"
	DefaultDBUser      = "postgres"
	DefaultDBPassword  = "password"
	DefaultDBName      = "revenue_leak_detective_dev"
	DefaultSSLMode     = "disable"
	DefaultEnvironment = "development"
	DefaultLogLevel    = "INFO"
	DefaultConfigVer   = "unknown"
	DefaultDebug       = "false"
)

// Environment variable names
const (
	EnvAPIHost          = "API_HOST"
	EnvAPIPort          = "API_PORT"
	EnvPostgresURL      = "POSTGRES_URL"
	EnvPostgresHost     = "POSTGRES_HOST"
	EnvPostgresPort     = "POSTGRES_PORT"
	EnvPostgresUser     = "POSTGRES_USER"
	EnvPostgresPassword = "POSTGRES_PASSWORD" //nolint:gosec // This is an environment variable name, not a hardcoded password
	EnvPostgresDB       = "POSTGRES_DB"
	EnvPostgresSSL      = "POSTGRES_SSL"
	EnvEnvironment      = "ENVIRONMENT"
	EnvLogLevel         = "LOG_LEVEL"
	EnvConfigVer        = "CONFIG_VERSION"
	EnvDebug            = "DEBUG"
)
