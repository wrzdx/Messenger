package domain

import "unicode/utf8"

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
	if l := utf8.RuneCountInString(p.Username); l < 5 || l > 32 {
		return ErrInvalidUsername
	}
	if l := utf8.RuneCountInString(p.FirstName); l < 1 || l > 64 {
		return ErrInvalidFirstName
	}

	if p.LastName != nil {
		if l := utf8.RuneCountInString(*p.LastName); l > 64 {
			return ErrInvalidLastName
		}
	}

	if p.Bio != nil {
		if l := utf8.RuneCountInString(*p.Bio); l > 70 {
			return ErrInvalidBio
		}
	}

	if l := utf8.RuneCountInString(p.Password); l < 8 || l > 32 {
		return ErrInvalidPassword
	}

	return nil
}
