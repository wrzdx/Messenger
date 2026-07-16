package users_transport_http

import (
	"context"
	"messenger/internal/core/domain"
	users_service "messenger/internal/features/users/service"

	"github.com/google/uuid"
)

type UsersService interface {
	GetUser(
		ctx context.Context,
		id uuid.UUID,
	) (domain.User, error)
	UpdateProfile(
		ctx context.Context,
		userID uuid.UUID,
		command users_service.UpdateProfileCommand,
	) (domain.User, error)
}
