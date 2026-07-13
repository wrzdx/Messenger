package auth_postgres_repository

import (
	"context"
	"messenger/internal/core/domain"
)

func (r *AuthRepository) CreateSession(
	ctx context.Context,
	session domain.Session,
) (domain.Session, error)
