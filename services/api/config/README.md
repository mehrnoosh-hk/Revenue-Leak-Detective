# Configuration System

The Revenue Leak Detective API uses a comprehensive configuration system that provides structured, validated configuration management with environment-specific loading and secure logging.

## Overview

This package provides a comprehensive configuration management system organized as follows:

### Files:
- `types.go`:      Configuration structs and constants
- `errors.go`:     Error constants and messages
- `loader.go`:     Configuration loading and environment file handling
- `validator.go`:  Configuration validation logic
- `environment.go`: Environment detection and utilities
- `database.go`:   Database-specific configuration utilities
- `http.go`:       HTTP-specific configuration utilities
- `build_info.go`: Build information utilities
- `docs.go`:       Documentation generator and examples

### Structure
```
config/
├── README.md        # 📚 Complete documentation & entry point
├── types.go         # 📋 Configuration structs & constants
├── errors.go        # ⚠️ Error constants & messages
├── loader.go        # 🔄 Configuration loading logic
├── validator.go     # ✅ Validation logic
├── environment.go   # 🌍 Environment utilities
├── database.go      # 🗄️ Database utilities
├── http.go          # �� HTTP utilities
├── build_info.go    # 📦 Build info & printing utilities
├── docs.go          # 📚 Documentation generator
└── config_test.go   # 🧪 Tests
```

### Usage:
```go
import "rdl-api/config"

// Load configuration
cfg, err := config.LoadConfig("/path/to/.env")
if err != nil {
    log.Fatal(err)
}

// Print configuration overview
cfg.PrintEffectiveConfig(logger)
```

### Supported Features:
- Environment-specific loading (.env files)
- Production vs development behavior
- Comprehensive validation
- Database URL construction
- Build information tracking
- Auto-generated documentation
- Enhanced IDE support with detailed comments and tags

## Features

- **Structured Configuration**: Separate structs for HTTP, Database, and Environment configuration
- **Environment-Specific Loading**: Automatically loads `.env.dev` files in development environments
- **Validation**: Comprehensive validation of required configuration values
- **Security**: Masks sensitive information in logs and configuration output
- **Versioning**: Includes `CONFIG_VERSION` for configuration schema versioning
- **Flexible Database Configuration**: Supports both `DATABASE_URL` and individual parameter configuration
- **Auto-Generated Documentation**: Documentation is generated from code comments and struct tags
- **Enhanced IDE Support**: Comprehensive comments and validation tags for better IDE integration

## Documentation Generation

The configuration system includes auto-generated documentation that is extracted from the code comments and struct tags. This ensures that documentation is always up-to-date with the actual code.

### Generating Documentation

```go
import "rdl-api/config"

// Generate configuration documentation
docs := config.ConfigDocs()
fmt.Println(docs)

// Generate example configuration
example := config.GenerateExampleConfig()
fmt.Println(example)
```

### IDE Support

The configuration structs include comprehensive comments and validation tags that provide:

- **Autocomplete**: Detailed field descriptions and examples
- **Validation**: IDE validation for required fields and data types
- **Documentation**: Inline documentation for each configuration option
- **Examples**: Example values for each configuration field

## Configuration Structure
```go
type HTTPConfig struct {
    Port string `json:"port" yaml:"port"`
}
```

### Database Configuration
```go
type DatabaseConfig struct {
    URL      string `json:"url" yaml:"url"`           // Full DATABASE_URL (takes precedence)
    Host     string `json:"host" yaml:"host"`         // Database host
    Port     string `json:"port" yaml:"port"`         // Database port
    User     string `json:"user" yaml:"user"`         // Database user
    Password string `json:"password" yaml:"password"` // Database password
    Name     string `json:"name" yaml:"name"`         // Database name
    SSLMode  string `json:"ssl_mode" yaml:"ssl_mode"` // SSL mode
}
```

### Environment Configuration
```go
type EnvironmentConfig struct {
    Environment string     `json:"environment" yaml:"environment"` // Environment name
    LogLevel    slog.Level `json:"log_level" yaml:"log_level"`    // Logging level
    ConfigVer   string     `json:"config_version" yaml:"config_version"` // Config schema version
}
```

## Environment Variables

### HTTP Server
- `API_PORT`: HTTP server port (default: "3030")

### Database
- `DATABASE_URL`: Full database connection URL (recommended for production)
- `DB_HOST`: Database host (default: "localhost")
- `DB_PORT`: Database port (default: "5432")
- `DB_USER`: Database user (default: "postgres")
- `DB_PASSWORD`: Database password (default: "password")
- `DB_NAME`: Database name (default: "revenue_leak_detective_dev")
- `DB_SSL_MODE`: SSL mode (default: "disable")

### Environment
- `ENVIRONMENT`: Environment name (default: "development")
- `LOG_LEVEL`: Log level (default: "INFO")

## Environment File Loading

The system supports loading configuration from environment files using the `godotenv` library. The env file path is specified via command line flag:

```bash
./api -env-file=/path/to/your/.env.dev
```

**Note**: 
- `.env.*` files are gitignored for security
- If no `-env-file` flag is provided, only environment variables are used
- The env file path must be explicitly provided - no fallback mechanisms

## Configuration Validation

The system validates:

1. **HTTP Configuration**: Port number validity (1-65535)
2. **Database Configuration**: 
   - If `DATABASE_URL` is provided, validates URL format
   - Otherwise, requires `DB_HOST`, `DB_USER`, and `DB_NAME`
   - Validates port number if specified
3. **Environment Configuration**: Validates environment name against allowed values

## Usage Examples

### Basic Configuration Loading
```go
import "rdl-api/config"

// Load config without env file (uses only environment variables)
cfg, err := config.LoadConfig("")
if err != nil {
    log.Fatalf("Failed to load configuration: %v", err)
}

// Load config with specific env file
cfg, err := config.LoadConfig("/path/to/.env.dev")
if err != nil {
    log.Fatalf("Failed to load configuration: %v", err)
}
```

### Accessing Configuration Values
```go
// HTTP configuration
port := cfg.GetPort()

// Database configuration
dbURL := cfg.DatabaseURL()

// Environment information
env := cfg.GetEnvironment()
logLevel := cfg.GetLogLevel()

// Environment checks
if cfg.IsDevelopment() {
    // Development-specific logic
}

if cfg.IsProduction() {
    // Production-specific logic
}
```

### Database URL Construction

The system automatically constructs database URLs:

1. **With DATABASE_URL**: Uses the provided URL directly
2. **Without DATABASE_URL**: Constructs from individual parameters
   - Format: `postgresql://user:password@host:port/dbname?sslmode=mode`

## Security Features

- **Secret Masking**: Passwords and sensitive data are masked in logs
- **Environment Isolation**: Development and production configurations are properly separated
- **Validation**: Prevents invalid configuration from being loaded

## Configuration Versioning

The `CONFIG_VERSION` constant tracks the configuration schema version. This allows for:

- Schema evolution tracking
- Migration scripts
- Configuration validation across versions

## Testing

The configuration system includes comprehensive tests covering:

- Default configuration loading
- Environment variable overrides
- Validation logic
- Database URL construction
- Secret masking
- Environment detection methods

Run tests with:
```bash
go test ./config/... -v
```

## Migration from Old System

The new configuration system replaces the old ad-hoc environment variable reading with:

- **Before**: Direct `os.Getenv()` calls scattered throughout code
- **After**: Centralized, validated configuration with structured access

Update your code to use the new getter methods:
- `cfg.ServerConfig.Port` → `cfg.GetPort()`
- `cfg.Env` → `cfg.GetEnvironment()`
- `cfg.LogLevel` → `cfg.GetLogLevel()`

## Best Practices

1. **Use Getter Methods**: Always use the provided getter methods for configuration access
2. **Environment Files**: Use `-env-file` flag for development, environment variables for production
3. **Validation**: The system validates configuration at startup - fix any validation errors
4. **Secrets**: Never commit `.env` files to version control
5. **Defaults**: Provide sensible defaults for all configuration values
6. **Explicit Configuration**: Always specify the env file path explicitly - no fallback mechanisms

## Troubleshooting

### Common Issues

1. **Validation Errors**: Check that all required environment variables are set
2. **Port Conflicts**: Ensure the configured port is available and valid
3. **Database Connection**: Verify database credentials and connectivity
4. **Environment Detection**: Ensure `ENVIRONMENT` variable is set correctly

### Debug Mode

Set `LOG_LEVEL=DEBUG` to see detailed configuration loading information.

### Configuration Output

The system automatically logs the effective configuration at startup (with secrets masked) to help with debugging.

