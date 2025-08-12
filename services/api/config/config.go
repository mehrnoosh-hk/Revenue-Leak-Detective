package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port string `json:"port"`
	LogLevel int `json:"golang_log_level"`
}

// LoadConfig loads the configuration from a file or environment variables.
func LoadConfig() (*Config, error) {
	config := &Config{
		Port: "8080", // Default port
		LogLevel: 0, // Default log level
	}

	if port, exists := os.LookupEnv("API_PORT"); exists {
		config.Port = port
	}

	if logLevel, exists := os.LookupEnv("Golang_LOG_LEVEL"); exists {		
		 if level, err := strconv.ParseInt(logLevel, 10, 0); err == nil {
			config.LogLevel = int(level)
		}
	}

	return config, nil
}
