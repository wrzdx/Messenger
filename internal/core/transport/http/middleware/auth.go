package http_middleware

import (
	"fmt"
	context "messenger/internal/core/context"
	core_context "messenger/internal/core/context"
	logger "messenger/internal/core/logger"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	const prefix = "Bearer "

	if !strings.HasPrefix(authHeader, prefix) {
		return ""
	}

	token := strings.TrimPrefix(authHeader, prefix)
	if token == "" {
		return ""
	}

	return token
}

type TokenProvider interface {
	ParseAccessToken(token string) (uuid.UUID, error)
}

func Auth(jwt TokenProvider) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			log := logger.FromContext(r.Context())
			sender := http_response.NewHTTPSender(log, w)
			userID, err := jwt.ParseAccessToken(token)
			if err != nil {
				sender.Error(fmt.Errorf("failed to authenticate: %w", err))
				return
			}
			appClaims := context.ContextClaims{
				UserID: userID,
			}
			ctx := core_context.WithClaims(r.Context(), appClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
