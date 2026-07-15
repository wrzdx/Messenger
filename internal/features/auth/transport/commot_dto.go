package auth_transport_http

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"  example:"qwerty"`
	FirstName string    `json:"first_name"  example:"Ivan"`
	LastName  *string   `json:"last_name"  example:"Ivanov"`
	CreatedAt time.Time `json:"created_at" example:"2026-02-26T10:30:00Z"`
	Bio       *string   `json:"bio"  example:"We didn't choose this path. Circumstance chose it for us. We're simply trying to keep climbing."`
}
