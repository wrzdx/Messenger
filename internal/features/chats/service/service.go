package chats_service

type ChatsService struct {
	chatsRepo ChatsRepository
	usersRepo UsersRepository
	txmanager TXManager
}

func NewChatsService(
	chatsRepo ChatsRepository,
	usersRepo UsersRepository,
	txmanager TXManager,
) *ChatsService {
	return &ChatsService{
		chatsRepo: chatsRepo,
		usersRepo: usersRepo,
		txmanager: txmanager,
	}
}
