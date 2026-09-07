package realtime_transport_ws

import (
	"context"
	"go.uber.org/zap"
	"messenger/internal/core/domain"
	"messenger/internal/core/logger"
)

type Notifier struct {
	hub          *Hub
	participants ParticipantsRepository
	log          *logger.Logger
}

func NewNotifier(hub *Hub, participants ParticipantsRepository, log *logger.Logger) *Notifier {
	return &Notifier{
		hub:          hub,
		participants: participants,
		log:          log,
	}
}

// MessageCreated makes one best-effort notification attempt after commit.
// Failures are logged here and do not change the result of saving the message.
func (p *Notifier) MessageCreated(
	ctx context.Context,
	message domain.Message,
) {
	log := p.log.With(zap.String("chat_id", message.ChatID.String()), zap.String("message_id", message.ID.String()))
	recipientIDs, err := p.participants.GetParticipants(ctx, message.ChatID)
	if err != nil {
		log.Error("get message notification recipients", zap.Error(err))
		return
	}
	event := Event{
		Type: msgCreated,
		Data: Message{
			ID:        message.ID,
			ChatID:    message.ChatID,
			SenderID:  message.SenderID,
			Content:   message.Content,
			CreatedAt: message.CreatedAt,
			UpdatedAt: message.UpdatedAt,
		},
	}
	for _, id := range recipientIDs {
		if err := p.hub.Publish(id, event); err != nil {
			log.Error("publish message created", zap.Error(err), zap.String("recipient_id", id.String()))
		}
	}
}
