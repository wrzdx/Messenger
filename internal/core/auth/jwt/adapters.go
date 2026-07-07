package core_auth_jwt

import (
	core_auth "messenger/internal/core/auth"

	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	core_auth.Claims
	jwt.RegisteredClaims
}
