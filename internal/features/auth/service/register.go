package auth_service

import (
	"context"
	"fmt"
	auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

func (s *AuthService) Register(
	ctx context.Context,
	user domain.User,
	password string,
) (
	domain.User,
	auth.TokenPair,
	error,
) {
	if err := user.Validate(); err != nil {
		return domain.User{}, auth.TokenPair{}, fmt.Errorf("validate user: %w", err)
	}
	if err := domain.ValidatePassword(password); err != nil {
		return domain.User{}, auth.TokenPair{}, fmt.Errorf("validate password: %w", err)
	}
	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		return domain.User{}, auth.TokenPair{}, fmt.Errorf("hash password: %w", err)
	}

	user = domain.NewUser(
		uuid.New(),
		payload.Username,
		payload.FirstName,
		payload.LastName,
		time.Now(),
		payload.Bio,
		passwordHash,
	)
	tokens, err := s.tokenService.GenerateTokenPair(user.ID)
	if err != nil {
		return domain.User{}, domain.TokenPair{}, fmt.Errorf(
			"generate refresh token: %w",
			err,
		)
	}

	user, err = s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, domain.TokenPair{}, err
	}

	return user, tokens, nil
}
