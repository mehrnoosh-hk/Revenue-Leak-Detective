package config

import (
	"net"
	"net/url"
)

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
