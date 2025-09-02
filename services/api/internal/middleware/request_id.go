package middleware

import (
    "context"
    "net/http"

    "github.com/google/uuid"
)

// contextKey is a private type to avoid key collisions in context.
// Using an unexported type ensures uniqueness across packages.
type contextKey string

const requestIDKey contextKey = "requestID"

// GetRequestID returns the request ID stored in the request context, if any.
func GetRequestID(r *http.Request) string {
    if v := r.Context().Value(requestIDKey); v != nil {
        if id, ok := v.(string); ok {
            return id
        }
    }
    return ""
}

// RequestID is middleware that ensures every request has a unique ID.
// The ID is generated at the edge of the server, stored in the request
// context for downstream handlers, and set as the `X-Request-ID` response
// header so that clients and logs can correlate activity.
func RequestID() Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Check if a request ID is already provided by the client via header.
            reqID := r.Header.Get("X-Request-ID")
            if reqID == "" {
                // Generate a new UUID if none is provided.
                reqID = uuid.NewString()
            }

            // Store the request ID in the context so that it can be retrieved downstream.
            ctx := context.WithValue(r.Context(), requestIDKey, reqID)
            r = r.WithContext(ctx)

            // Set the request ID on the response so that clients see the value.
            w.Header().Set("X-Request-ID", reqID)

            next.ServeHTTP(w, r)
        })
    }
}