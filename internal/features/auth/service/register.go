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
	if err := payload.Validate(); err != nil {
		return domain.User{}, auth.TokenPair{}, fmt.Errorf("validate register payload: %w", err)
	}
	passwordHash, err := s.hasher.Hash(payload.Password)
	if err != nil {
		return domain.User{}, auth.TokenPair{}, fmt.Errorf("hash password: %w", err)
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
	claims := auth.AccessClaims{
		UserID: user.ID,
	}
	tokenID := uuid.New()
	tokens, err := s.tokenService.GenerateTokenPair(claims, tokenID)
	if err != nil {
		return domain.User{}, auth.TokenPair{}, fmt.Errorf(
			"generate refresh token: %w",
			err,
		)
	}

	user, err = s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, auth.TokenPair{}, err
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

func NewRegisterPayload(
	username string,
	firstName string,
	lastName *string,
	bio *string,
	password string,
) RegisterPayload {
	return RegisterPayload{
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Bio:       bio,
		Password:  password,
	}
}

func (p *RegisterPayload) Validate() error {
	if err := domain.ValidateUsername(p.Username); err != nil {
		return err
	}

	if err := domain.ValidateFirstName(p.FirstName); err != nil {
		return err
	}

	if p.LastName != nil {
		if err := domain.ValidateLastName(*p.LastName); err != nil {
			return err
		}
	}

	if p.Bio != nil {
		if err := domain.ValidateBio(*p.Bio); err != nil {
			return err
		}
	}

	if err := domain.ValidatePassword(p.Password); err != nil {
		return err
	}

	return nil
}
