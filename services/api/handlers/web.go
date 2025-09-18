package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

// WriteJSONResponse writes a JSON response with proper error handling and logging
func WriteJSONResponse[T any](
	ctx context.Context,
	w http.ResponseWriter,
	logger *slog.Logger,
	data T,
	statusCode int,
) {
	// Set content type header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Encode and write JSON response
	if err := json.NewEncoder(w).Encode(data); err != nil {

		logger.ErrorContext(ctx, "Failed to encode JSON response", "error", err)
		// Write error response
		http.Error(w, ErrInternalServerError.Error(), http.StatusInternalServerError)
	}
}

// WriteJSONErrorResponse writes a JSON error response with proper error handling and logging
func WriteJSONErrorResponse(
	ctx context.Context,
	w http.ResponseWriter,
	logger *slog.Logger,
	err error,
	statusCode int,
) {
	logger.ErrorContext(ctx, "Failed to write JSON error response", "error", err)
	http.Error(w, err.Error(), http.StatusInternalServerError) // Never return the actual error to the client
}

// WriteJSONSuccessResponse writes a successful JSON response (200 OK)
func WriteJSONSuccessResponse[T any](
	ctx context.Context,
	w http.ResponseWriter,
	logger *slog.Logger,
	data T,
) {
	WriteJSONResponse(ctx, w, logger, data, http.StatusOK)
}
