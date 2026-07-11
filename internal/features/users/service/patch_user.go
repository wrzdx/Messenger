package users_service

import (
	"context"
	"fmt"
	"messenger/internal/core/auth"
	"messenger/internal/core/domain"
	core_types "messenger/internal/core/types"

	"github.com/google/uuid"
)

func (s *UsersService) PatchUser(
	ctx context.Context,
	id uuid.UUID,
	patch UserPatch,
) (domain.User, error) {
	user, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user: %w", err)
	}
	if user.DeletedAt != nil {
		return domain.User{}, auth.ErrInvalidToken
	}
	user, err = ApplyPatch(user, patch)
	if err != nil {
		return domain.User{}, fmt.Errorf("apply user patch: %w", err)
	}

	patchedUser, err := s.userRepository.PatchUser(ctx, id, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("patch user: %w", err)
	}

	return patchedUser, nil
}

type UserPatch struct {
	Username  *string
	FirstName *string
	LastName  core_types.Nullable[string]
	Bio       core_types.Nullable[string]
}

func ApplyPatch(user domain.User, p UserPatch) (domain.User, error) {
	if p.Username != nil {
		user.Username = *p.Username
	}

	if p.FirstName != nil {
		user.FirstName = *p.FirstName
	}

	if p.LastName.Set {
		user.LastName = p.LastName.Value
	}

	if p.Bio.Set {
		user.Bio = p.Bio.Value
	}

	if err := user.Validate(); err != nil {
		return domain.User{}, err
	}

	return user, nil
}
