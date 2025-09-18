package config

import (
	"fmt"
	"log/slog"
)

// printEffectiveConfig prints the effective configuration (excluding secrets and build information)
// This provides a human-readable overview of the current configuration
func printEffectiveConfig(c *Config, logger *slog.Logger) {
	logger.Info("Effective configuration:")
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

// printBuildInfo prints the build information
func printBuildInfo(c *Config, logger *slog.Logger) {
	logger.Info("Build information:")
	logger.Info(fmt.Sprintf("version: %s", c.BuildInfo.GIT_TAG))
	logger.Info(fmt.Sprintf("commit: %s", c.BuildInfo.GIT_COMMIT_FULL))
	logger.Info(fmt.Sprintf("build_date: %s", c.BuildInfo.BUILD_TIMESTAMP))
	logger.Info(fmt.Sprintf("git_branch: %s", c.BuildInfo.GIT_BRANCH))
	logger.Info(fmt.Sprintf("git_tag: %s", c.BuildInfo.GIT_TAG))
	logger.Info(fmt.Sprintf("git_dirty: %s", c.BuildInfo.GIT_DIRTY))
	logger.Info(fmt.Sprintf("git_commit_message: %s", c.BuildInfo.GIT_COMMIT_MESSAGE))
}
