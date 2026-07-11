package auth_service

import (
	"context"
	"errors"
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
		nil,
		payload.Bio,
		passwordHash,
	)
	tokenID := uuid.New()
	tokens, err := s.tokenProvider.GenerateTokenPair(user.ID, tokenID)
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
	tmpUser := domain.NewUser(
		uuid.New(),
		p.Username,
		p.FirstName,
		p.LastName,
		time.Now(),
		nil,
		p.Bio,
		"",
	)
	userErr := tmpUser.Validate()
	passwordErr := domain.ValidatePassword(p.Password)
	if userErr == nil {
		if passwordErr != nil {
			return domain.ValidationErr(domain.UserEntity, map[string]string{
				"password": passwordErr.Error(),
			})
		}
		return nil

	} else {
		de, ok := userErr.(domain.DetailedError)
		if !ok {
			return errors.New("validate user return unexpected type")
		}
		if passwordErr != nil {
			de.Details["password"] = passwordErr.Error()
		}
		return de
	}
}
