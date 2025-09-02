package config

// Error constants for configuration validation and loading
const (
	// Validation errors
	ErrInvalidPort           = "invalid port"
	ErrPortOutOfRange        = "port out of range"
	ErrMissingDBHost         = "missing database host"
	ErrMissingDBUser         = "missing database user"
	ErrMissingDBName         = "missing database name"
	ErrInvalidDBURL          = "invalid database URL"
	ErrInvalidEnvironment    = "invalid environment"
	ErrMissingRequiredEnvVar = "missing required environment variable"

	// Loading errors
	ErrEnvFileNotFound        = "environment file not found"
	ErrEnvFileLoadFailed      = "failed to load environment file"
	ErrConfigValidationFailed = "configuration validation failed"
)
