package auth_service

import (
	"context"
	"errors"
	"fmt"
	core_auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"
	core_errors "messenger/internal/core/errors"
)

func (s *AuthService) Login(
	ctx context.Context,
	credentials domain.UserCredentials,
) (domain.Token, domain.Token, error) {
	if err := credentials.Validate(); err != nil {
		return domain.Token{}, domain.Token{}, fmt.Errorf(
			"validate credentials: %w",
			err,
		)
	}
	userAuth, err := s.userRepository.GetUserAuth(
		ctx,
		credentials.Username,
	)
	if err != nil {
		return domain.Token{}, domain.Token{}, fmt.Errorf(
			"get user password hash: %w",
			err,
		)
	}

	if err := s.hasher.Compare(userAuth.PasswordHash, credentials.Password); err != nil {
		if errors.Is(err, core_auth.ErrInvalidCredentials) {
			return domain.Token{}, domain.Token{}, fmt.Errorf(
				"%v: %w",
				err,
				core_errors.ErrUnauthorized,
			)
		}
		return domain.Token{}, domain.Token{}, fmt.Errorf(
			"compare passwords: %w",
			core_errors.ErrUnauthorized,
		)
	}

	refresh, refreshExpires, err := s.jwtProvider.GenerateRefreshToken(int(userAuth.UserID))
	if err != nil {
		return domain.Token{}, domain.Token{}, fmt.Errorf(
			"generate refresh token: %w",
			err,
		)
	}
	access, accessExpires, err := s.jwtProvider.GenerateAccessToken(int(userAuth.UserID))
	if err != nil {
		return domain.Token{}, domain.Token{}, fmt.Errorf(
			"generate access token: %w",
			err,
		)
	}

	refreshDomain := domain.NewToken(string(refresh), refreshExpires)
	accessDomain := domain.NewToken(string(access), accessExpires)

	return refreshDomain, accessDomain, nil
}
