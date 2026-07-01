package auth_service

import (
	"context"
	"errors"
	"fmt"
	core_auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"
	core_errors "messenger/internal/core/errors.go"
)

func (s *AuthService) Login(
	ctx context.Context,
	credentials domain.UserCredentials,
) (domain.RefreshToken, domain.AccessToken, error) {
	if err := credentials.Validate(); err != nil {
		return domain.RefreshToken{}, domain.AccessToken{}, fmt.Errorf(
			"validate credentials: %w",
			err,
		)
	}
	userAuth, err := s.userRepository.GetUserAuth(
		ctx,
		credentials.Username,
	)
	if err != nil {
		return domain.RefreshToken{}, domain.AccessToken{}, fmt.Errorf(
			"get user password hash: %w",
			err,
		)
	}

	if err := s.hasher.Compare(userAuth.PasswordHash, credentials.Password); err != nil {
		if errors.Is(err, core_auth.ErrInvalidCredentials) {
			return domain.RefreshToken{}, domain.AccessToken{}, fmt.Errorf(
				"%v: %w",
				err,
				core_errors.ErrUnauthorized,
			)
		}
		return domain.RefreshToken{}, domain.AccessToken{}, fmt.Errorf(
			"compare passwords: %w",
			core_errors.ErrUnauthorized,
		)
	}

	refresh, refreshExpires, err := s.jwtProvider.GenerateRefreshToken(int(userAuth.UserID))
	if err != nil {
		return domain.RefreshToken{}, domain.AccessToken{}, fmt.Errorf(
			"generate refresh token: %w",
			err,
		)
	}
	access, accessExpires, err := s.jwtProvider.GenerateAccessToken(int(userAuth.UserID))
	if err != nil {
		return domain.RefreshToken{}, domain.AccessToken{}, fmt.Errorf(
			"generate access token: %w",
			err,
		)
	}

	refreshDomain := domain.RefreshToken(domain.NewToken(string(refresh), refreshExpires))
	accessDomain := domain.AccessToken(domain.NewToken(string(access), accessExpires))

	return refreshDomain, accessDomain, nil
}
