package users_service

import (
	"context"
	"fmt"
	"messenger/internal/core/auth"
	"messenger/internal/core/domain"
	http_types "messenger/internal/core/transport/http/types"

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
	Username  http_types.Nullable[string]
	FirstName http_types.Nullable[string]
	LastName  http_types.Nullable[string]
	Bio       http_types.Nullable[string]
}

func ApplyPatch(user domain.User, p UserPatch) (domain.User, error) {
	if p.Username.Set {
		if p.Username.Value == nil {
			return domain.User{}, domain.ErrInvalidUsername
		}
		user.Username = *p.Username.Value
	}

	if p.FirstName.Set {
		if p.FirstName.Value == nil {
			return domain.User{}, domain.ErrInvalidFirstName
		}
		user.FirstName = *p.FirstName.Value
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
