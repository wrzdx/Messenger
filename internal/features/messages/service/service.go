package messages_service

type MessagesService struct {
	messagesRepo MessagesRepository
	chatsRepo    ChatsRepository
	txmanager    TXManager
}

func NewMessagesService(
	messagesRepo MessagesRepository,
	chatsRepo ChatsRepository,
	txmanager TXManager,
) *MessagesService {
	return &MessagesService{
		messagesRepo: messagesRepo,
		chatsRepo:    chatsRepo,
		txmanager:    txmanager,
	}
}
