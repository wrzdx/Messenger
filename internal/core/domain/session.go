package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidSession = errors.New("invalid session")
)

type Session struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	CurrentTokenID uuid.UUID
	LastUsedAt     time.Time
	CreatedAt      time.Time
	ExpiresAt      time.Time
}

func NewSession(
	id, userID, currentTokenID uuid.UUID,
	createdAt, expiresAt time.Time,
) (Session, error) {
	session := Session{
		ID:             id,
		UserID:         userID,
		CurrentTokenID: currentTokenID,
		LastUsedAt:     createdAt,
		CreatedAt:      createdAt,
		ExpiresAt:      expiresAt,
	}
	if err := session.Validate(); err != nil {
		return Session{}, err
	}
	return session, nil
}

func (s Session) Validate() error {
	if s.ID == uuid.Nil {
		return fmt.Errorf("id is nil: %w", ErrInvalidSession)
	}
	if s.UserID == uuid.Nil {
		return fmt.Errorf("user_id s nil: %w", ErrInvalidSession)
	}

	if s.CurrentTokenID == uuid.Nil {
		return fmt.Errorf("current_token_id is nil: %w", ErrInvalidSession)
	}
	if s.ExpiresAt.IsZero() {
		return fmt.Errorf("expires_at is zero value: %w", ErrInvalidSession)
	}
	if s.LastUsedAt.IsZero() {
		return fmt.Errorf("last_user_at is zero value: %w", ErrInvalidSession)
	}
	if s.CreatedAt.IsZero() {
		return fmt.Errorf("created_at is zero value: %w", ErrInvalidSession)
	}
	if !s.ExpiresAt.After(s.CreatedAt) {
		return fmt.Errorf("expires_at must be after created_at: %w", ErrInvalidSession)
	}
	if s.LastUsedAt.After(s.ExpiresAt) {
		return fmt.Errorf("last_used_at cannot be after expires_at: %w", ErrInvalidSession)
	}
	if s.LastUsedAt.Before(s.CreatedAt) {
		return fmt.Errorf("last_user_at cannot be before created_at: %w", ErrInvalidSession)
	}
	return nil
}
