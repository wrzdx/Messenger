package users_service

import (
	"context"
	"fmt"
	"messenger/internal/core/domain"
)

func (s *UsersService) GetUsers(
	ctx context.Context,
	pagination domain.Pagination,
) ([]domain.User, error) {
	if err := pagination.Validate(); err != nil {
		return nil, err
	}
	users, err := s.userRepository.GetUsers(ctx, pagination)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	return users, nil
}
