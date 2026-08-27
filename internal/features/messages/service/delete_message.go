package messages_service

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/core/domain"

	"github.com/google/uuid"
)

func (s *MessagesService) DeleteMessage(
	ctx context.Context,
	command DeleteMessageCommand,
) error {
	if err := command.validate(); err != nil {
		return fmt.Errorf("validate delete message command: %w", err)
	}
	if err := s.messagesRepo.CheckParticipant(
		ctx,
		command.ChatID,
		command.SenderID,
	); err != nil {
		return fmt.Errorf("check participant: %w", err)
	}
	existing, err := s.messagesRepo.GetMessage(
		ctx,
		command.MessageID,
	)
	if err != nil {
		return fmt.Errorf("get message: %w", err)
	}
	if command.ChatID != existing.ChatID || existing.SenderID != command.SenderID {
		return domain.ErrNotFound
	}

	if err := s.txmanager.WithinTransaction(
		ctx,
		func(ctx context.Context) error {
			chat, err := s.chatsRepo.GetChatForUpdate(ctx, command.ChatID)
			if err != nil {
				return fmt.Errorf("get chat for update: %w", err)
			}
			if chat.LastMessageID == nil {
				return errors.New(
					"db inconsistency: chat has messages but no last message id",
				)
			}
			if *chat.LastMessageID == command.MessageID {
				msgs, err := s.messagesRepo.GetMessages(ctx, command.ChatID, nil, 2)
				if err != nil {
					return fmt.Errorf("get messages: %w", err)
				}
				if len(msgs) == 1 {
					chat.LastMessageID = nil
				} else if len(msgs) == 2 {
					chat.LastMessageID = new(msgs[1].ID)
				} else {
					return errors.New("invalid number of message from get messages query")
				}
				if err := s.chatsRepo.UpdateChatLastMsgID(
					ctx,
					chat.ID,
					chat.LastMessageID,
				); err != nil {
					return fmt.Errorf("update chat last message id: %w", err)
				}
			}
			if err := s.messagesRepo.DeleteMessage(ctx, command.MessageID); err != nil {
				return fmt.Errorf("delete message: %w", err)
			}

			return nil
		},
	); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

type DeleteMessageCommand struct {
	ChatID    uuid.UUID
	SenderID  uuid.UUID
	MessageID uuid.UUID
}

func (c DeleteMessageCommand) validate() error {
	details := make(map[string]string)
	if c.ChatID == uuid.Nil {
		details["chat_id"] = "chat id is nil"
	}
	if c.SenderID == uuid.Nil {
		details["sender_id"] = "sender id is nil"
	}
	if c.MessageID == uuid.Nil {
		details["message_id"] = "message id is nil"
	}

	if len(details) > 0 {
		return domain.DetailedError{
			Err:     ErrInvalidInput,
			Details: details,
		}
	}

	return nil
}
