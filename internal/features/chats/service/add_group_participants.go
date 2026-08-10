package chats_service

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/core/domain"
	"slices"
	"time"

	"github.com/google/uuid"
)

type AddGroupParticipantStatus string

const (
	Added         AddGroupParticipantStatus = "added"
	Unavailable   AddGroupParticipantStatus = "unavailable"
	AlreadyMember AddGroupParticipantStatus = "already_member"
)

func (s *ChatsService) AddGroupParticipants(
	ctx context.Context,
	command AddGroupParticipantsCommand,
) ([]AddGroupParticipantResult, error) {
	if err := command.validate(); err != nil {
		return nil, err
	}
	requester, err := s.chatsRepo.GetGroupParticipant(
		ctx,
		command.GroupID,
		command.RequesterID,
	)
	if err != nil {
		return nil, fmt.Errorf("get requester: %w", err)
	}
	if requester.Role() != domain.OwnerRole && requester.Role() != domain.AdminRole {
		return nil, fmt.Errorf(
			"requester has no enough rights: %w",
			ErrNotEnoughRights,
		)
	}
	toAdd := make([]domain.GroupParticipant, 0)
	positions := make([]int, 0)
	result := make([]AddGroupParticipantResult, len(command.ParticipantIDs))

	statuses, err := s.chatsRepo.GetParticipantsStatus(
		ctx,
		command.ParticipantIDs,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get participant statuses: %w",
			err,
		)
	}

	now := time.Now()
	foundByUserID := make(map[uuid.UUID]struct{})
	for _, status := range statuses {
		if status.Found {
			foundByUserID[status.UserID] = struct{}{}
		}
	}
	for i, id := range command.ParticipantIDs {
		result[i] = AddGroupParticipantResult{
			UserID: id,
			Status: Added,
		}
		if _, ok := foundByUserID[id]; !ok {
			result[i].Status = Unavailable
		} else {
			gp, err := domain.NewGroupParticipant(
				command.GroupID,
				id,
				nil,
				now,
				domain.MemberRole,
			)
			if err != nil {
				return nil, fmt.Errorf(
					"new group participant: %w",
					err,
				)
			} else {
				toAdd = append(toAdd, gp)
				positions = append(positions, i)
			}
		}
	}
	if err := s.txmanager.WithinTransaction(ctx, func(ctx context.Context) error {
		addedParticipants, err := s.chatsRepo.AddGroupParticipants(
			ctx,
			command.GroupID,
			toAdd,
		)
		if err != nil {
			return fmt.Errorf("add group participants: %w", err)
		}
		if len(addedParticipants) != len(toAdd) {
			return errors.New("len of add group result and toAdd not match")
		}

		for i, p := range toAdd {
			if !addedParticipants[i] {
				result[positions[i]] = AddGroupParticipantResult{
					UserID: p.UserID,
					Status: AlreadyMember,
				}
			}
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %w", err)
	}

	return result, nil
}

type AddGroupParticipantsCommand struct {
	GroupID        uuid.UUID
	RequesterID    uuid.UUID
	ParticipantIDs []uuid.UUID
}

type AddGroupParticipantResult struct {
	UserID uuid.UUID
	Status AddGroupParticipantStatus
}

func (c AddGroupParticipantsCommand) validate() error {
	if c.GroupID == uuid.Nil {
		return fmt.Errorf("groupID is nil: %w", ErrInvalidAddGroupParticipantsQuery)
	}
	if c.RequesterID == uuid.Nil {
		return fmt.Errorf("requesterID is nil: %w", ErrInvalidAddGroupParticipantsQuery)
	}
	if slices.Contains(c.ParticipantIDs, uuid.Nil) {
		return fmt.Errorf(
			"nil id in participants list: %w",
			ErrInvalidAddGroupParticipantsQuery,
		)
	}
	if len(c.ParticipantIDs) > 100 {
		return fmt.Errorf(
			"size of participants list is at most 100 per request: %w",
			ErrInvalidAddGroupParticipantsQuery,
		)
	}
	return nil
}
