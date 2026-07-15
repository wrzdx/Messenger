package auth_postgres_repository

import (
	"messenger/internal/core/postgres"
	"time"
)

type AuthRepository struct {
	db      postgres.DBTX
	timeout time.Duration
}

func NewAuthRepository(db postgres.DBTX, timeout time.Duration) AuthRepository {
	return AuthRepository{
		db:      db,
		timeout: timeout,
	}
}

const (
	sessionsUsersFK          = "sessions_user_id_fkey"
	sessionsPK               = "sessions_pkey"
	sessionsCurrentTokenIDUK = "sessions_current_token_id_key"
)
