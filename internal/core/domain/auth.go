package domain

import (
	"fmt"
	core_errors "messenger/internal/core/errors"
	"time"
)

type Token struct {
	Token   string
	Expires time.Time
}

func NewToken(token string, expires time.Time) Token {
	return Token{
		Token:   token,
		Expires: expires,
	}
}

type UserCredentials struct {
	Username string
	Password string
}

func NewCredentials(username string, password string) UserCredentials {
	return UserCredentials{
		Username: username,
		Password: password,
	}
}

func (c *UserCredentials) Validate() error {
	passwordLen := len([]rune(c.Password))
	if passwordLen < 8 || passwordLen > 32 {
		return fmt.Errorf(
			"invalid `Password` len: %d: %w",
			passwordLen,
			core_errors.ErrInvalidArgument,
		)
	}

	usernameLen := len([]rune(c.Username))
	if usernameLen < 5 || usernameLen > 32 {
		return fmt.Errorf(
			"invalid `Username` len: %d: %w",
			usernameLen,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

type UserAuth struct {
	UserID       int
	PasswordHash string
}

func NewUserAuth(id int, passwordHash string) UserAuth {
	return UserAuth{
		UserID:       id,
		PasswordHash: passwordHash,
	}
}
