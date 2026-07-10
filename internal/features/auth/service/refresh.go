package auth_service

import (
	"context"
	"fmt"
	"messenger/internal/core/auth"

	"github.com/google/uuid"
)

func (s *AuthService) Refresh(
	ctx context.Context,
	token string,
) (auth.TokenPair, error) {
	userID, tokenID, err := s.tokenProvider.ParseRefreshToken(token)
	if err != nil {
		return auth.TokenPair{}, fmt.Errorf("parse refresh: %w", err)
	}
	user, err := s.usersRepository.GetUser(ctx, userID)
	if err != nil {
		return auth.TokenPair{}, fmt.Errorf("get user: %w", err)
	}
	if user.DeletedAt != nil {
		return auth.TokenPair{}, auth.ErrInvalidToken
	}

	tokenID = uuid.New()
	tokens, err := s.tokenProvider.GenerateTokenPair(user.ID, tokenID)
	if err != nil {
		return auth.TokenPair{}, fmt.Errorf(
			"generate tokens: %w",
			err,
		)
	}

	return tokens, nil
}
