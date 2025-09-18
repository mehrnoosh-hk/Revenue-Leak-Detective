package middleware

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
)

// contextKey is defined in the middleware package
const tenantIDKey contextKey = "tenantID"

var (
	ErrMissingOrInvalidTenantContext = errors.New("missing or invalid tenant context")
)

// GetTenantID returns the tenant ID stored in the request context, if any.
func GetTenantID(r *http.Request) (uuid.UUID, bool) {
	if v := r.Context().Value(tenantIDKey); v != nil {
		if id, ok := v.(uuid.UUID); ok {
			return id, true
		}
	}
	return uuid.Nil, false
}

// TenantContext extracts the tenant ID from JWT token (or header in development)
// and stores it in the request context. It skips tenant validation for excluded paths.
//
// excludedPaths: List of paths that should skip tenant validation (e.g., health endpoints)
// Example: []string{"/healthz", "/health", "/live", "/ready"}
func TenantContext(l *slog.Logger, isDevelopment bool, excludedPaths []string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip tenant validation for excluded paths
			if isPathExcluded(r.URL.Path, excludedPaths) {
				next.ServeHTTP(w, r)
				return
			}

			// Extract tenant ID from Authorization header (JWT token)
			// or from X-Tenant-ID header for development/testing
			tenantID := extractTenantID(l, r, isDevelopment)

			if tenantID == uuid.Nil {
				http.Error(w, ErrMissingOrInvalidTenantContext.Error(), http.StatusUnauthorized)
				return
			}

			// Add tenant ID to request context
			ctx := context.WithValue(r.Context(), tenantIDKey, tenantID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func extractTenantID(l *slog.Logger, r *http.Request, isDevelopment bool) uuid.UUID {
	// Extract from X-Tenant-ID header (for development/testing only)
	if isDevelopment {
		tenantHeader := r.Header.Get("X-Tenant-ID")
		if tenantHeader != "" {
			if tenantID, err := uuid.Parse(tenantHeader); err == nil {
				l.Debug("Tenant ID extracted from X-Tenant-ID header", "tenantID", tenantID)
				return tenantID
			}
		}
	}
	// Extract from JWT token (recommended for production)
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if tenantID, err := extractTenantFromJWT(token); err == nil {
			l.Debug("Tenant ID extracted from JWT token", "tenantID", tenantID)
			return tenantID
		}
	}

	return uuid.Nil
}

func extractTenantFromJWT(token string) (uuid.UUID, error) {

	// TODO: Implement JWT token parsing and tenant ID extraction
	// This would typically involve:
	// 1. Verify JWT signature
	// 2. Parse claims
	// 3. Extract tenant_id from claims
	// 4. Return tenant UUID

	// For now, return a placeholder
	_ = token // Just for the linter
	return uuid.Nil, fmt.Errorf("JWT parsing not implemented")
}

// isPathExcluded checks if the given path matches any of the excluded paths.
func isPathExcluded(path string, excludedPaths []string) bool {
	return slices.Contains(excludedPaths, path)
}
