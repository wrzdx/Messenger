package auth_service

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/core/auth"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

func (s *AuthService) Login(
	ctx context.Context,
	username string,
	password string,
) (auth.TokenPair, error) {
	user, err := s.usersRepository.GetUserByUsername(
		ctx,
		username,
	)
	if err != nil {
		return auth.TokenPair{}, domain.ErrInvalidCredentials
	}

	if err := s.hasher.Compare(user.PasswordHash, password); err != nil {
		if errors.Is(err, auth.ErrPasswordMismatch) {
			return auth.TokenPair{}, domain.ErrInvalidCredentials
		}
		return auth.TokenPair{}, fmt.Errorf(
			"compare passwords: %w",
			err,
		)
	}
	tokenID := uuid.New()
	tokens, err := s.tokenProvider.GenerateTokenPair(user.ID, tokenID)
	if err != nil {
		return auth.TokenPair{}, fmt.Errorf(
			"generate tokens: %w",
			err,
		)
	}

	return tokens, nil
}
