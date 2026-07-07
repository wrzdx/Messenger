package auth_service

import (
	"context"
	"fmt"
	core_auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"
)

func (s *AuthService) Refresh(
	ctx context.Context,
	token string,
) (core_auth.AuthTokens, error) {
	claims, err := s.jwtProvider.ParseToken(token)
	if err != nil {
		return core_auth.AuthTokens{}, domain.ErrInvalidRefreshToken
	}
	user, err := s.usersRepository.GetUser(ctx, claims.UserID)
	if err != nil {
		return core_auth.AuthTokens{}, domain.ErrUserNotFound
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
