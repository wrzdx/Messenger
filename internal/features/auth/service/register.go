package auth_service

import (
	"context"
	"fmt"
	core_auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

func (s *AuthService) Register(
	ctx context.Context,
	payload domain.RegisterUserPayload,
) (
	domain.User,
	core_auth.AuthTokens,
	error,
) {
	if err := payload.Validate(); err != nil {
		return domain.User{}, core_auth.AuthTokens{}, fmt.Errorf("validate payload: %w", err)
	}
	passwordHash, err := s.hasher.Hash(payload.Password)
	if err != nil {
		return domain.User{}, core_auth.AuthTokens{}, fmt.Errorf("hash password: %w", err)
	}

	user := domain.NewUser(
		uuid.New(),
		payload.Username,
		payload.FirstName,
		payload.LastName,
		time.Now(),
		payload.Bio,
		passwordHash,
	)
	tokens, err := s.jwtProvider.GenerateTokens(user.ID)
	if err != nil {
		return domain.User{}, core_auth.AuthTokens{}, fmt.Errorf(
			"generate refresh token: %w",
			err,
		)
	}

	user, err = s.authRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, core_auth.AuthTokens{}, err
	}

	return user, tokens, nil
}
