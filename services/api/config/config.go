package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type ServerConfig struct {
	Port string `json:"server_port"`
}

type DatabaseConfig struct {
	DBPort     string `json:"db_port"`
	DBHost     string `json:"db_host"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
}

// Config holds the application configuration
type Config struct {
	ServerConfig   ServerConfig   `json:"server_config"`
	DatabaseConfig DatabaseConfig `json:"database_config"`
	Env            string         `json:"environment"`
	LogLevel       slog.Level     `json:"log_level"`
}

// LoadConfig loads the configuration from environment variables with validation
func LoadConfig() (*Config, error) {
	config := &Config{
		ServerConfig: ServerConfig{
			Port: "3030",
		},
		DatabaseConfig: DatabaseConfig{
			DBPort:     "5432",
			DBHost:     "localhost",
			DBUser:     "postgres",
			DBPassword: "password",
			DBName:     "revenue_leak_detective_dev",
		},
		LogLevel: slog.LevelInfo, // Default log level
		Env:      "development",  // Default environment
	}

	// Load port
	if port, exist := os.LookupEnv("API_PORT"); exist {
		if err := validatePort(port); err != nil {
			return nil, fmt.Errorf("invalid port: %w", err)
		}
		config.ServerConfig.Port = port
	}

	// Load log level
	if logLevel, exist := os.LookupEnv("LOG_LEVEL"); exist {
		level, err := parseLogLevel(logLevel)
		if err != nil {
			return nil, fmt.Errorf("invalid log level: %w", err)
		}
		config.LogLevel = level
	}

	// Load environment
	if env, exist := os.LookupEnv("ENVIRONMENT"); exist {
		config.Env = strings.ToLower(env)
	}

	return config, nil
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
func parseLogLevel(level string) (slog.Level, error) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "WARN", "WARNING":
		return slog.LevelWarn, nil
	case "ERROR":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unknown log level: %s", level)
	}
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Env == "development" || c.Env == "dev"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Env == "production" || c.Env == "prod"
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", c.DatabaseConfig.DBUser, c.DatabaseConfig.DBPassword, c.DatabaseConfig.DBHost, c.DatabaseConfig.DBPort, c.DatabaseConfig.DBName)
}
