package http_middleware

import (
	"fmt"
	auth "messenger/internal/core/auth"
	context "messenger/internal/core/context"
	core_errors "messenger/internal/core/errors"
	logger "messenger/internal/core/logger"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"
	"strings"
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

func Auth(jwt auth.TokenService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			log := logger.FromContext(r.Context())
			responseHandler := http_response.NewHTTPSender(log, w)
			payload, err := jwt.ParseAccessToken(token)
			if err != nil {
				respError := core_errors.Error{
					Err:     fmt.Errorf("auth middleware: %w", err),
					Code:    core_errors.INVALID_TOKEN,
					Message: "failed parse token",
				}
				responseHandler.Error(respError)
				return
			}
			appClaims := context.ContextClaims{
				UserID: payload.UserID,
			}
			ctx := context.WithClaims(r.Context(), appClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
