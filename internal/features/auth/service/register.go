package auth_service

import (
	"context"
	"fmt"
	"messenger/internal/core/auth"
	"messenger/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

func (s *AuthService) Register(
	ctx context.Context,
	payload RegisterPayload,
) (
	domain.User,
	auth.TokenPair,
	error,
) {
	profile, err := domain.NewUserProfile(
		payload.Username,
		payload.FirstName,
		payload.LastName,
		payload.Bio,
	)
	if err != nil {
		return domain.User{}, auth.TokenPair{}, fmt.Errorf("new user profile: %w", err)
	}
	if err := auth.ValidatePassword(payload.Password); err != nil {
		return domain.User{}, auth.TokenPair{}, fmt.Errorf("validate password: %w", err)
	}
	passwordHash, err := s.hasher.Hash(payload.Password)
	if err != nil {
		return domain.User{}, auth.TokenPair{}, fmt.Errorf("hash password: %w", err)
	}
	now := time.Now()
	user, err := domain.NewUser(
		uuid.New(),
		profile,
		now,
		nil,
		passwordHash,
	)
	if err != nil {
		return domain.User{}, auth.TokenPair{}, fmt.Errorf("new user: %w", err)
	}

	session, err := domain.NewSession(
		uuid.New(),
		user.ID,
		uuid.New(),
		now,
		now.Add(s.config.SessionTTL),
	)
	if err != nil {
		return domain.User{}, auth.TokenPair{}, fmt.Errorf("new session: %w", err)
	}
	tokens, err := s.generateTokenPair(session, now)
	if err != nil {
		return domain.User{}, auth.TokenPair{}, fmt.Errorf(
			"generate token pair: %w",
			err,
		)
	}
	err = s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := s.usersRepository.CreateUser(ctx, user); err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		if err := s.sessionsRepository.CreateSession(ctx, session); err != nil {
			return fmt.Errorf("create session: %w", err)
		}
		return nil
	})

	if err != nil {
		return domain.User{}, auth.TokenPair{}, fmt.Errorf("transaction: %w", err)
	}

	return user, tokens, nil
}

type RegisterPayload struct {
	Username  string
	FirstName string
	LastName  *string
	Bio       *string
	Password  string
}
