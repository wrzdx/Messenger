package core_auth

import (
	"context"

	"github.com/google/uuid"
)

type userIDContextKey struct{}

var (
	key = userIDContextKey{}
)

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(key).(uuid.UUID)

	return userID, ok
}

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, key, userID)
}
