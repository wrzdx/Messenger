package auth_service

import (
	"context"
	"fmt"
	core_auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"
)

func (s *AuthService) Login(
	ctx context.Context,
	username string,
	password string,
) (core_auth.AuthTokens, error) {
	user, err := s.authRepository.GetUserByUsername(
		ctx,
		username,
	)
	if err != nil {
		return core_auth.AuthTokens{}, domain.ErrInvalidCredentials
	}

	if err := s.hasher.Compare(user.PasswordHash, password); err != nil {
		return core_auth.AuthTokens{}, fmt.Errorf(
			"compare passwords: %w",
			err,
		)
	}

	tokens, err := s.jwtProvider.GenerateTokens(user.ID)
	if err != nil {
		return core_auth.AuthTokens{}, fmt.Errorf(
			"generate refresh token: %w",
			err,
		)
	}

	return tokens, nil
}
