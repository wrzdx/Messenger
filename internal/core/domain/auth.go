package domain

import (
	"time"

	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type Claims struct {
	UserID uuid.UUID
	Type   TokenType
}

func NewClaims(userID uuid.UUID, tokenType TokenType) Claims {
	return Claims{
		UserID: userID,
		Type:   tokenType,
	}
}

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

type RegisterUserPayload struct {
	Username  string
	FirstName string
	LastName  *string
	Bio       *string
	Password  string
}

func NewRegisterUserPayload(
	username string,
	firstName string,
	lastName *string,
	bio *string,
	password string,
) RegisterUserPayload {
	return RegisterUserPayload{
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Bio:       bio,
		Password:  password,
	}
}

func (p *RegisterUserPayload) Validate() error {
	if err := ValidateUsername(p.Username); err != nil {
		return err
	}

	if err := ValidateFirstName(p.FirstName); err != nil {
		return err
	}

	if p.LastName != nil {
		if err := ValidateLastName(*p.LastName); err != nil {
			return err
		}
	}

	if p.Bio != nil {
		if err := ValidateBio(*p.Bio); err != nil {
			return err
		}
	}

	if err := ValidatePassword(p.Password); err != nil {
		return err
	}

	return nil
}
