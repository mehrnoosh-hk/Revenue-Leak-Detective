package handlers

import (
	"time"

	"github.com/google/uuid"
)

type CreateTenantRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type CreateTenantResponse struct {
	ID uuid.UUID `json:"id"`
}

type GetTenantRequest struct {
	ID uuid.UUID `json:"id"`
}

type GetTenantResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
