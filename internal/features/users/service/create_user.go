package users_service

import (
	"context"
	"fmt"
	"messenger/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

func (s *UsersService) CreateUser(
	ctx context.Context,
	payload domain.RegisterUserPayload,
) (domain.User, error) {
	if err := payload.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("%w: validate payload", err)
	}
	passwordHash, err := s.hasher.Hash(payload.Password)
	if err != nil {
		return domain.User{}, fmt.Errorf("%w: hash password", err)
	}

	user := domain.NewUser(
		uuid.New(),
		payload.Username,
		payload.FirstName,
		payload.LastName,
		time.Now(),
		payload.Bio,
		string(passwordHash),
	)

	user, err = s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
