package users_service

type UsersService struct {
	userRepository UsersRepository
	hasher         Hasher
}

func NewUsersService(
	usersRepository UsersRepository,
	hasher Hasher,
) *UsersService {
	return &UsersService{
		userRepository: usersRepository,
		hasher:         hasher,
	}
}
