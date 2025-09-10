package services

import (
	"context"
	"rdl-api/internal/domain/models"
)

type EventServiceInterface interface {
	CreateEvent(ctx context.Context, params models.CreateEventParams) (models.Event, error)
}
