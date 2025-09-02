package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// LoadConfig loads the configuration with a specific env file path
func LoadConfig(envFilePath string) (*Config, error) {
	// Load environment file if specified (only in non-production environments)
	if envFilePath != "" {
		if err := loadEnvFile(envFilePath); err != nil {
			return nil, fmt.Errorf("%s: %w", ErrEnvFileLoadFailed, err)
		}
	}

	// After loading env file, determine the environment
	env := os.Getenv(EnvEnvironment)
	isProduction := isProductionEnvironment(env)

	config := &Config{
		HTTP: HTTPConfig{
			Host: getEnvValue(EnvAPIHost, isProduction, DefaultAPIHost),
			Port: getEnvValue(EnvAPIPort, isProduction, DefaultAPIPort),
		},
		Database: DatabaseConfig{
			URL:      os.Getenv(EnvPostgresURL),
			Host:     getEnvValue(EnvPostgresHost, isProduction, DefaultDBHost),
			Port:     getEnvValue(EnvPostgresPort, isProduction, DefaultDBPort),
			User:     getEnvValue(EnvPostgresUser, isProduction, DefaultDBUser),
			Password: getEnvValue(EnvPostgresPassword, isProduction, DefaultDBPassword),
			DBName:   getEnvValue(EnvPostgresDB, isProduction, DefaultDBName),
			SSLMode:  getEnvValue(EnvPostgresSSL, isProduction, DefaultSSLMode),
		},
		Environment: func() EnvironmentConfig {
			debugStr := getEnvValue(EnvDebug, isProduction, DefaultDebug)
			debugVal, err := strconv.ParseBool(strings.ToLower(debugStr))
			if err != nil {
				return EnvironmentConfig{
					Environment: getEnvValue(EnvEnvironment, isProduction, DefaultEnvironment),
					Debug:       false,
					LogLevel:    parseLogLevel(getEnvValue(EnvLogLevel, isProduction, DefaultLogLevel)),
					ConfigVer:   getEnvValue(EnvConfigVer, isProduction, DefaultConfigVer),
				}
			}
			return EnvironmentConfig{
				Environment: getEnvValue(EnvEnvironment, isProduction, DefaultEnvironment),
				Debug:       debugVal,
				LogLevel:    parseLogLevel(getEnvValue(EnvLogLevel, isProduction, DefaultLogLevel)),
				ConfigVer:   getEnvValue(EnvConfigVer, isProduction, DefaultConfigVer),
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
		return nil, fmt.Errorf("%s: %w", ErrConfigValidationFailed, err)
	}

	return config, nil
}

// loadEnvFile loads environment file if path is provided
func loadEnvFile(envFilePath string) error {
	// Check if file exists
	if _, err := os.Stat(envFilePath); err != nil {
		return fmt.Errorf("%s: %s", ErrEnvFileNotFound, envFilePath)
	}

	// Load the env file
	if err := godotenv.Load(envFilePath); err != nil {
		return fmt.Errorf("%s: %w", ErrEnvFileLoadFailed, err)
	}

	return nil
}
