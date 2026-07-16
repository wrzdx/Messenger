package users_service

type UsersService struct {
	userRepository UsersRepository
	txManager      TXManager
}

func NewUsersService(
	usersRepository UsersRepository,
	txManager TXManager,
) *UsersService {
	return &UsersService{
		userRepository: usersRepository,
		txManager: txManager,
	}
}
