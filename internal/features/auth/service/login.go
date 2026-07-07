package auth_service

import (
	"context"
	"fmt"
	"messenger/internal/core/domain"
)

func (s *AuthService) Login(
	ctx context.Context,
	username string,
	password string,
) (domain.Token, domain.Token, error) {
	user, err := s.authRepository.GetUserByUsername(
		ctx,
		username,
	)
	if err != nil {
		return domain.Token{}, domain.Token{}, domain.ErrInvalidCredentials
	}

	if err := s.hasher.Compare(user.PasswordHash, password); err != nil {
		return domain.Token{}, domain.Token{}, fmt.Errorf(
			"compare passwords: %w",
			err,
		)
	}

	refresh, err := s.jwtProvider.GenerateRefreshToken(user.ID)
	if err != nil {
		return domain.Token{}, domain.Token{}, fmt.Errorf(
			"generate refresh token: %w",
			err,
		)
	}
	access, err := s.jwtProvider.GenerateAccessToken(user.ID)
	if err != nil {
		return domain.Token{}, domain.Token{}, fmt.Errorf(
			"generate access token: %w",
			err,
		)
	}

	return refresh, access, nil
}
