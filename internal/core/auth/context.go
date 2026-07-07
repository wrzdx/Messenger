package core_auth

import (
	"context"
	"messenger/internal/core/domain"
)

type claimsContextKey struct{}

var (
	claimsKey = claimsContextKey{}
)

func WithClaims(ctx context.Context, claims domain.Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

func ClaimsFromContext(ctx context.Context) (domain.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(domain.Claims)
	return claims, ok
}
