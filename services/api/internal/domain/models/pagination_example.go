// Package models contains domain models for the business entities.
// This file demonstrates how to use the generic pagination models with different entity types.
package models

// Example usage of pagination with different entity types:

// Example 1: Events pagination (already implemented)
// func (r *eventsRepository) GetAllEventsPaginated(ctx context.Context, tenantID uuid.UUID, params PaginationParams) (PaginatedResponse[Event], error)

// Example 2: Users pagination (future implementation)
// func (r *usersRepository) GetAllUsersPaginated(ctx context.Context, tenantID uuid.UUID, params PaginationParams) (PaginatedResponse[User], error)

// Example 3: Actions pagination (future implementation)
// func (r *actionsRepository) GetAllActionsPaginated(ctx context.Context, tenantID uuid.UUID, params PaginationParams) (PaginatedResponse[Action], error)

// Example 4: Leaks pagination (future implementation)
// func (r *leaksRepository) GetAllLeaksPaginated(ctx context.Context, tenantID uuid.UUID, params PaginationParams) (PaginatedResponse[Leak], error)

// Example 5: Generic usage in service layer
// func (s *eventService) GetEventsWithPagination(ctx context.Context, tenantID uuid.UUID, page, pageSize int) (PaginatedResponse[Event], error) {
//     params := PaginationParams{
//         Limit:  int32(pageSize),
//         Offset: int32((page - 1) * pageSize),
//     }
//     return s.repo.GetAllEventsPaginated(ctx, tenantID, params)
// }

// Example 6: HTTP handler usage
// func (h *EventHandler) GetEvents(c *gin.Context) {
//     var params PaginationParams
//     if err := c.ShouldBindQuery(&params); err != nil {
//         c.JSON(400, gin.H{"error": err.Error()})
//         return
//     }
//
//     response, err := h.service.GetEventsWithPagination(c.Request.Context(), tenantID, params)
//     if err != nil {
//         c.JSON(500, gin.H{"error": err.Error()})
//         return
//     }
//
//     c.JSON(200, response)
// }
