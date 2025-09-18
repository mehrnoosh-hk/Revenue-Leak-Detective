// Package models contains domain models for the business entities.
// This package defines the core data structures and types used throughout
// the application.
package models

// PaginationParams represents parameters for paginated queries.
// This struct is used to control pagination behavior for list operations.
//
// Fields:
//   - Limit: Maximum number of items to return (required, must be > 0)
//   - Offset: Number of items to skip (required, must be >= 0)
//
// Example usage:
//   - First page (10 items): Limit=10, Offset=0
//   - Second page (10 items): Limit=10, Offset=10
//   - Third page (10 items): Limit=10, Offset=20
type PaginationParams struct {
	Limit  int32 `json:"limit" validate:"required,min=1,max=1000"`
	Offset int32 `json:"offset" validate:"required,min=0"`
}

// PaginatedResponse represents a paginated response containing items and metadata.
// This struct provides both the requested data and pagination information.
//
// Fields:
//   - Items: The actual data items for the current page
//   - TotalCount: Total number of items across all pages
//   - Limit: Number of items per page
//   - Offset: Number of items skipped
//   - HasNext: Whether there are more items after the current page
//   - HasPrevious: Whether there are items before the current page
type PaginatedResponse[T any] struct {
	Items       []T   `json:"items"`
	TotalCount  int64 `json:"total_count"`
	Limit       int32 `json:"limit"`
	Offset      int32 `json:"offset"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}

// NewPaginatedResponse creates a new PaginatedResponse with calculated metadata.
//
// Parameters:
//   - items: The items for the current page
//   - totalCount: Total number of items across all pages
//   - limit: Number of items per page
//   - offset: Number of items skipped
//
// Returns:
//   - PaginatedResponse[T]: A properly configured paginated response
func NewPaginatedResponse[T any](items []T, totalCount int64, limit, offset int32) PaginatedResponse[T] {
	return PaginatedResponse[T]{
		Items:       items,
		TotalCount:  totalCount,
		Limit:       limit,
		Offset:      offset,
		HasNext:     int64(offset+limit) < totalCount,
		HasPrevious: offset > 0,
	}
}
