package users_service

import (
	"context"
	"fmt"
	"messenger/internal/core/auth"

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
	if user.DeletedAt != nil {
		return auth.ErrInvalidToken
	}
	if err := s.hasher.Compare(user.PasswordHash, old_password); err != nil {
		return fmt.Errorf("compare passwords: %w", err)
	}

	if err := auth.ValidatePassword(new_password); err != nil {
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
