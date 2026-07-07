package users_transport_http

import (
	"messenger/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type UserDTOResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username" validate:"required" example:"qwerty"`
	FirstName string    `json:"first_name" validate:"required" example:"Ivan"`
	LastName  *string   `json:"last_name" example:"Ivanov"`
	CreatedAt time.Time `json:"created_at" example:"2026-02-26T10:30:00Z"`
	Bio       *string   `json:"bio" example:"We didn't choose this path. Circumstance chose it for us. We're simply trying to keep climbing."`
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		Bio:       user.Bio,
	}
}

func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}

	return usersDTO
}
