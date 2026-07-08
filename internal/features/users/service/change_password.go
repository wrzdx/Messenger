package users_service

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/core/auth"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

func (s *UsersService) ChangePassword(
	ctx context.Context,
	id uuid.UUID,
	old_password string,
	new_password string,
) error {
	user, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if err := s.hasher.Compare(user.PasswordHash, old_password); err != nil {
		if errors.Is(err, auth.ErrPasswordMismatch) {
			return domain.ErrWrongPassword
		}
		return fmt.Errorf("compare passwords: %w", err)
	}

	if err := domain.ValidatePassword(new_password); err != nil {
		return fmt.Errorf("validate new password: %w", err)
	}

	passwordHash, err := s.hasher.Hash(new_password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := s.userRepository.ChangePassword(ctx, id, passwordHash); err != nil {
		return fmt.Errorf("change user password: %w", err)
	}

	return nil
}
