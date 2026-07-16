package users_transport_http

import (
	"context"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

type UsersService interface {
	GetUser(
		ctx context.Context,
		id uuid.UUID,
	) (domain.User, error)
}
