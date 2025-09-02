package config

// GetPort returns the configured HTTP port
func (c *Config) GetPort() string {
	return c.HTTP.Port
}
