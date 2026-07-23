package chats_service

import (
	"context"
	"fmt"
	"messenger/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

func (s *ChatsService) ListGroupParticipants(
	ctx context.Context,
	requesterID uuid.UUID,
	query ListGroupParticipantsQuery,
) (GroupParticipantPage, error) {
	if requesterID == uuid.Nil {
		return GroupParticipantPage{}, fmt.Errorf("nil requester id: %w", domain.ErrNotFound)
	}

	query = query.normalize()
	if err := query.validate(); err != nil {
		return GroupParticipantPage{}, fmt.Errorf("validate query: %w", err)
	}

	_, err := s.chatsRepo.GetGroupParticipant(ctx, query.ChatID, requesterID)
	if err != nil {
		return GroupParticipantPage{}, fmt.Errorf("get chat participant: %w", err)
	}

	allParticipants, err := s.chatsRepo.ListGroupParticipants(
		ctx,
		query.ChatID,
		query.Before,
		query.Limit+1,
	)

	if err != nil {
		return GroupParticipantPage{}, fmt.Errorf("list group participants: %w", err)
	}
	if allParticipants == nil {
		allParticipants = make([]ParticipantInfo, 0)
	}

	var page GroupParticipantPage
	hasMore := len(allParticipants) > query.Limit
	if hasMore {
		allParticipants = allParticipants[:query.Limit]
	}

	page.Participants = allParticipants

	if hasMore {
		last := page.Participants[len(page.Participants)-1]
		page.NextCursor = &GroupParticipantCursor{
			ParticipantID: last.ID,
			JoinedAt:      last.JoinedAt,
		}
	}
	return page, nil
}

type GroupParticipantCursor struct {
	ParticipantID uuid.UUID
	JoinedAt      time.Time
}

type ListGroupParticipantsQuery struct {
	ChatID uuid.UUID
	Before *GroupParticipantCursor
	Limit  int
}

func (q ListGroupParticipantsQuery) normalize() ListGroupParticipantsQuery {
	if q.Limit == 0 {
		q.Limit = 50
	}
	return q
}

func (q ListGroupParticipantsQuery) validate() error {
	fields := make(map[string]string)
	if q.ChatID == uuid.Nil {
		fields["chat_id"] = "chat_id is nil"
	}
	if q.Limit < 0 || q.Limit > 100 {
		fields["limit"] = "limit must be between 1 and 100"
	}

	if q.Before != nil {
		if q.Before.JoinedAt.IsZero() {
			fields["joined_at"] = "joined_at of group participant cursor cannot be zero value"
		}
		if q.Before.ParticipantID == uuid.Nil {
			fields["participant_id"] = "participant id of group participant cursor cannot be nil"
		}
	}
	if len(fields) > 0 {
		return domain.DetailedError{
			Err:     ErrInvalidListGroupParticipantsQuery,
			Details: fields,
		}
	}

	return nil
}

type ParticipantInfo struct {
	ID        uuid.UUID
	FirstName string
	LastName  *string
	Role      string
	JoinedAt  time.Time
}

func NewParticipantInfo(
	participant domain.ChatParticipant,
	role domain.GroupRole,
	profile domain.UserProfile,
) (ParticipantInfo, error) {
	gparticipant, err := domain.NewGroupParticipant(
		participant.ChatID,
		participant.UserID,
		participant.LastReadMessageID,
		participant.JoinedAt,
		role,
	)
	if err != nil {
		return ParticipantInfo{}, fmt.Errorf(
			"validate participant: %v: %w",
			err,
			ErrInvalidGroupParticipantItem,
		)
	}
	if err := profile.Validate(); err != nil {
		return ParticipantInfo{}, fmt.Errorf(
			"validate participant: %v: %w",
			err,
			ErrInvalidGroupParticipantItem,
		)
	}
	info := ParticipantInfo{
		ID:        gparticipant.UserID,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Role:      string(gparticipant.Role()),
		JoinedAt:  gparticipant.JoinedAt,
	}

	return info, nil
}

type GroupParticipantPage struct {
	Participants []ParticipantInfo
	NextCursor   *GroupParticipantCursor
}
