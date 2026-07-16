package users_service

type UsersService struct {
	userRepository     UsersRepository
	sessionsRepository SessionsRepository
	txManager          TXManager
}

func NewUsersService(
	usersRepository UsersRepository,
	sessionsRepository SessionsRepository,
	txManager TXManager,
) *UsersService {
	return &UsersService{
		userRepository:     usersRepository,
		sessionsRepository: sessionsRepository,
		txManager:          txManager,
	}
}
