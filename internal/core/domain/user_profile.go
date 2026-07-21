package domain

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	UsernamePattern = regexp.MustCompile("^[a-zA-Z0-9_]{5,32}$")
)

var ErrInvalidUserProfile = errors.New("invalid user profile")

type UserProfile struct {
	Username  string
	FirstName string
	LastName  *string
	Bio       *string
}

func NewUserProfile(
	username string,
	firstName string,
	lastName *string,
	bio *string,
) (UserProfile, error) {
	profile := UserProfile{
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Bio:       bio,
	}
	profile = profile.normalize()
	if err := profile.Validate(); err != nil {
		return UserProfile{}, err
	}
	return profile, nil
}

func (p UserProfile) normalize() UserProfile {
	p.Username = strings.TrimSpace(p.Username)
	p.FirstName = strings.TrimSpace(p.FirstName)
	if p.LastName != nil {
		if trimmed := strings.TrimSpace(*p.LastName); trimmed != "" {
			p.LastName = &trimmed
		} else {
			p.LastName = nil
		}
	}
	if p.Bio != nil {
		if trimmed := strings.TrimSpace(*p.Bio); trimmed != "" {
			p.Bio = &trimmed
		} else {
			p.Bio = nil
		}
	}

	return p
}
func (p UserProfile) Validate() error {
	fields := make(map[string]string)
	if !UsernamePattern.MatchString(p.Username) {
		fields["username"] = "must contain between 5 and 32 characters, and only ASCII letters, digits, and underscores"
	}
	if l := utf8.RuneCountInString(p.FirstName); l < 1 || l > 64 {
		fields["first_name"] = "fists_name contain between 1 and 64 characters"
	}
	if p.LastName != nil {
		if l := utf8.RuneCountInString(*p.LastName); l > 64 {
			fields["last_name"] = "last_name must contain at most 64 characters"
		}
	}

	if p.Bio != nil {
		if l := utf8.RuneCountInString(*p.Bio); l > 70 {
			fields["bio"] = "bio must contain at most 70 characters"
		}
	}

	if len(fields) > 0 {
		return DetailedError{
			Err:     ErrInvalidUserProfile,
			Details: fields,
		}
	}

	return nil
}
