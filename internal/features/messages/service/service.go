package messages_service

type MessagesService struct {
	messagesRepo MessagesRepository
	chatsRepo    ChatsRepository
	txmanager    TXManager
	notifier     Notifier
}

func NewMessagesService(
	messagesRepo MessagesRepository,
	chatsRepo ChatsRepository,
	txmanager TXManager,
	notifier Notifier,
) *MessagesService {
	return &MessagesService{
		messagesRepo: messagesRepo,
		chatsRepo:    chatsRepo,
		txmanager:    txmanager,
		notifier:     notifier,
	}
}
