package auth_postgres_repository

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	LastUsedTokenID uuid.UUID
	CreatedAt       time.Time
	ExpiresAt       time.Time
}
