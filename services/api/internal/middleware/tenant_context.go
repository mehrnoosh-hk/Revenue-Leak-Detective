package middleware

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

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
// and stores it in the request context.
func TenantContext(l *slog.Logger, isDevelopment bool) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
				l.Info("Tenant ID extracted from X-Tenant-ID header", "tenantID", tenantID)
				return tenantID
			}
		}
	}
	// Extract from JWT token (recommended for production)
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if tenantID, err := extractTenantFromJWT(token); err == nil {
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
	return uuid.Nil, fmt.Errorf("JWT parsing not implemented, %s", token)
}
