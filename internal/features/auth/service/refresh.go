package auth_service

import (
	"context"
	"fmt"
	auth "messenger/internal/core/auth"

	"github.com/google/uuid"
)

func (s *AuthService) Refresh(
	ctx context.Context,
	token string,
) (auth.TokenPair, error) {
	payload, err := s.tokenService.ParseRefreshToken(token)
	if err != nil {
		return auth.TokenPair{}, fmt.Errorf("parse refresh: %w", err)
	}
	user, err := s.usersRepository.GetUser(ctx, payload.UserID)
	if err != nil {
		return auth.TokenPair{}, fmt.Errorf("get user: %w", err)
	}

	tokenID := uuid.New()
	claims := auth.AccessClaims{
		UserID: user.ID,
	}
	tokens, err := s.tokenService.GenerateTokenPair(claims, tokenID)
	if err != nil {
		return auth.TokenPair{}, fmt.Errorf(
			"generate tokens: %w",
			err,
		)
	}

	return tokens, nil
}
