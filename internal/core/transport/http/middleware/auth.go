package core_http_middleware

import (
	"errors"
	"fmt"
	core_auth "messenger/internal/core/auth"
	core_logger "messenger/internal/core/logger"
	core_http_response "messenger/internal/core/transport/http/response"
	"net/http"
	"strings"
)

var (
	ErrMissingAuthorizationHeader = errors.New("missing authorization header")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
)

type jwtProvider interface {
	ParseToken(token string) (core_auth.Claims, error)
}

func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrMissingAuthorizationHeader
	}

	const prefix = "Bearer "

	if !strings.HasPrefix(authHeader, prefix) {
		return "", ErrInvalidAuthorizationHeader
	}

	token := strings.TrimPrefix(authHeader, prefix)
	if token == "" {
		return "", ErrInvalidAuthorizationHeader
	}

	return token, nil
}

func Auth(jwt jwtProvider) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := extractBearerToken(r)
			log := core_logger.FromContext(r.Context())
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)
			if err != nil {
				responseHandler.ErrorResponse(core_http_response.Error{
					Error:   fmt.Errorf("auth middleware: %w", err),
					Status:  http.StatusUnauthorized,
					Message: err.Error(),
				})
				return
			}
			claims, err := jwt.ParseToken(token)
			if err != nil {
				responseHandler.ErrorResponse(core_http_response.Error{
					Error:   fmt.Errorf("auth middleware: %w", err),
					Status:  http.StatusUnauthorized,
					Message: err.Error(),
				})
				return
			}
			ctx := core_auth.WithUserID(r.Context(), claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
