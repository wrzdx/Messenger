package users_service

import (
	"context"
	"fmt"
	"messenger/internal/core/domain"
)

func (s *UsersService) CreateUser(
	ctx context.Context,
	user domain.User,
	credentials domain.UserCredentials,
) (domain.User, error) {
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}
	if err := credentials.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user credentials: %w", err)
	}
	passwordHash, err := s.hasher.Hash(credentials.Password)
	if err != nil {
		return domain.User{}, fmt.Errorf("hash passoword: %w", err)
	}

	user, err = s.userRepository.CreateUser(ctx, user, string(passwordHash))
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
