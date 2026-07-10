package users_service

import (
	"context"
	"fmt"
	"messenger/internal/core/domain"
)

func (s *UsersService) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	if limit != nil {
		if err := domain.ValidateLimit(*limit); err != nil {
			return nil, err
		}
	}

	if offset != nil {
		if err := domain.ValidateOffset(*offset); err != nil {
			return nil, err
		}
	}
	users, err := s.userRepository.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	return users, nil
}
