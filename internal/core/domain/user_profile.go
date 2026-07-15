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
	normalizedUsername string
	username           string
	firstName          string
	lastName           *string
	bio                *string
}

func NewUserProfile(
	username string,
	firstName string,
	lastName *string,
	bio *string,
) (UserProfile, error) {
	profile := UserProfile{
		username:  username,
		firstName: firstName,
		lastName:  lastName,
		bio:       bio,
	}
	profile = profile.normalize()
	if err := profile.validate(); err != nil {
		return UserProfile{}, err
	}
	return profile, nil
}

func (p UserProfile) Username() string  { return p.username }
func (p UserProfile) FirstName() string { return p.firstName }
func (p UserProfile) LastName() *string {
	if p.lastName == nil {
		return nil
	}

	lastName := *p.lastName
	return &lastName
}
func (p UserProfile) Bio() *string {
	if p.bio == nil {
		return nil
	}
	bio := *p.bio
	return &bio
}

func (p UserProfile) normalize() UserProfile {
	p.username = strings.TrimSpace(p.username)
	p.firstName = strings.TrimSpace(p.firstName)
	if p.lastName != nil {
		if trimmed := strings.TrimSpace(*p.lastName); trimmed != "" {
			p.lastName = &trimmed
		} else {
			p.lastName = nil
		}
	}
	if p.bio != nil {
		if trimmed := strings.TrimSpace(*p.bio); trimmed != "" {
			p.bio = &trimmed
		} else {
			p.bio = nil
		}
	}

	return p
}
func (p UserProfile) validate() error {
	fields := make(map[string]string)
	if !UsernamePattern.MatchString(p.username) {
		fields["username"] = "must contain between 5 and 32 characters, and only ASCII letters, digits, and underscores"
	}
	if l := utf8.RuneCountInString(p.firstName); l < 1 || l > 64 {
		fields["first_name"] = "fists_name contain between 1 and 64 characters"
	}
	if p.lastName != nil {
		if l := utf8.RuneCountInString(*p.lastName); l > 64 {
			fields["last_name"] = "last_name must contain at most 64 characters"
		}
	}

	if p.bio != nil {
		if l := utf8.RuneCountInString(*p.bio); l > 70 {
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
