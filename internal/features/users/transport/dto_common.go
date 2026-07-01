package users_transport_http

import (
	"messenger/internal/core/domain"
	"time"
)

type UserDTOResponse struct {
	ID        int       `json:"id" example:"1"`
	Username  string    `json:"username" validate:"required,min=5,max=32" example:"qwerty"`
	FirstName string    `json:"first_name" validate:"required,min=1,max=64" example:"Ivan"`
	LastName  *string   `json:"last_name" validate:"max=64" example:"Ivanov"`
	CreatedAt time.Time `json:"created_at" example:"2026-02-26T10:30:00Z"`
	Bio       *string   `json:"bio" validate:"max=70" example:"We didn't choose this path. Circumstance chose it for us. We're simply trying to keep climbing."`
}

func userDTOFromDomain(user domain.User) domain.User {
	return domain.NewUser(
		user.ID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.CreatedAt,
		user.Bio,
	)
}
