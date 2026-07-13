package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewSession(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	userID := uuid.New()
	tokenID := uuid.New()
	expiresAt := now.Add(time.Hour)

	t.Run("creates valid session", func(t *testing.T) {
		session, err := NewSession(id, userID, tokenID, now, expiresAt)
		require.NoError(t, err)
		require.Equal(t, Session{
			ID:             id,
			UserID:         userID,
			CurrentTokenID: tokenID,
			LastUsedAt:     now,
			CreatedAt:      now,
			ExpiresAt:      expiresAt,
		}, session)
	})

	t.Run("rejects invalid session", func(t *testing.T) {
		session, err := NewSession(uuid.Nil, userID, tokenID, now, expiresAt)
		require.ErrorIs(t, err, ErrInvalidSession)
		require.Zero(t, session)
	})
}

func TestSessionValidate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		session   Session
		wantError error
	}{
		{
			name: "correct",
			session: Session{
				ID:             uuid.New(),
				UserID:         uuid.New(),
				CurrentTokenID: uuid.New(),
				LastUsedAt:     now,
				CreatedAt:      now,
				ExpiresAt:      now.Add(time.Hour),
			},
			wantError: nil,
		},
		{
			name: "nil id",
			session: Session{
				UserID:         uuid.New(),
				CurrentTokenID: uuid.New(),
				LastUsedAt:     now,
				CreatedAt:      now,
				ExpiresAt:      now.Add(time.Hour),
			},
			wantError: ErrInvalidSession,
		},
		{
			name: "nil user_id",
			session: Session{
				ID:             uuid.New(),
				CurrentTokenID: uuid.New(),
				LastUsedAt:     now,
				CreatedAt:      now,
				ExpiresAt:      now.Add(time.Hour),
			},
			wantError: ErrInvalidSession,
		},
		{
			name: "nil current_token_id",
			session: Session{
				ID:         uuid.New(),
				UserID:     uuid.New(),
				LastUsedAt: now,
				CreatedAt:  now,
				ExpiresAt:  now.Add(time.Hour),
			},
			wantError: ErrInvalidSession,
		},
		{
			name: "zero last_used_at",
			session: Session{
				ID:             uuid.New(),
				UserID:         uuid.New(),
				CurrentTokenID: uuid.New(),
				CreatedAt:      now,
				ExpiresAt:      now.Add(time.Hour),
			},
			wantError: ErrInvalidSession,
		},
		{
			name: "zero created_at",
			session: Session{
				ID:             uuid.New(),
				UserID:         uuid.New(),
				CurrentTokenID: uuid.New(),
				LastUsedAt:     now,
				ExpiresAt:      now.Add(time.Hour),
			},
			wantError: ErrInvalidSession,
		},
		{
			name: "zero expires_at",
			session: Session{
				ID:             uuid.New(),
				UserID:         uuid.New(),
				CurrentTokenID: uuid.New(),
				LastUsedAt:     now,
				CreatedAt:      now,
			},
			wantError: ErrInvalidSession,
		},
		{
			name: "expires_at before created_at",
			session: Session{
				ID:             uuid.New(),
				UserID:         uuid.New(),
				CurrentTokenID: uuid.New(),
				LastUsedAt:     now,
				CreatedAt:      now,
				ExpiresAt:      now.Add(-time.Hour),
			},
			wantError: ErrInvalidSession,
		},
		{
			name: "expires_at equal to created_at",
			session: Session{
				ID:             uuid.New(),
				UserID:         uuid.New(),
				CurrentTokenID: uuid.New(),
				LastUsedAt:     now,
				CreatedAt:      now,
				ExpiresAt:      now,
			},
			wantError: ErrInvalidSession,
		},
		{
			name: "last_used_at after expires_at",
			session: Session{
				ID:             uuid.New(),
				UserID:         uuid.New(),
				CurrentTokenID: uuid.New(),
				LastUsedAt:     now.Add(2 * time.Hour),
				CreatedAt:      now,
				ExpiresAt:      now.Add(time.Hour),
			},
			wantError: ErrInvalidSession,
		},
		{
			name: "last_used_at equal to expires_at",
			session: Session{
				ID:             uuid.New(),
				UserID:         uuid.New(),
				CurrentTokenID: uuid.New(),
				LastUsedAt:     now.Add(time.Hour),
				CreatedAt:      now,
				ExpiresAt:      now.Add(time.Hour),
			},
		},
		{
			name: "last_used_at before created_at",
			session: Session{
				ID:             uuid.New(),
				UserID:         uuid.New(),
				CurrentTokenID: uuid.New(),
				LastUsedAt:     now.Add(-time.Hour),
				CreatedAt:      now,
				ExpiresAt:      now.Add(time.Hour),
			},
			wantError: ErrInvalidSession,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.ErrorIs(t, tt.session.Validate(), tt.wantError)
		})
	}
}
