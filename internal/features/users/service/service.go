package users_service

import (
	core_auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"
)

type UsersService struct {
	hasher         core_auth.PasswordHasher
	userRepository domain.UsersRepository
}

func NewUsersService(
	usersRepository domain.UsersRepository,
	hasher core_auth.PasswordHasher,
) *UsersService {
	return &UsersService{
		userRepository: usersRepository,
		hasher:         hasher,
	}
}
