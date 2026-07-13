package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUserProfile(t *testing.T) {
	t.Run("normalizes fields and preserves username case", func(t *testing.T) {
		profile, err := NewUserProfile(
			"  Username_1  ",
			"  First Name  ",
			new("  Last Name  "),
			new("  Bio  "),
		)

		require.NoError(t, err)
		require.Equal(t, "Username_1", profile.Username())
		require.Equal(t, "First Name", profile.FirstName())
		require.Equal(t, "Last Name", *profile.LastName())
		require.Equal(t, "Bio", *profile.Bio())
	})

	t.Run("normalizes blank optional fields to nil", func(t *testing.T) {
		profile, err := NewUserProfile("Username_1", "First Name", new("  "), new("\t"))

		require.NoError(t, err)
		require.Nil(t, profile.LastName())
		require.Nil(t, profile.Bio())
	})

	validCases := []struct {
		name      string
		username  string
		firstName string
		lastName  *string
		bio       *string
	}{
		{name: "minimum username length", username: "Abc_1", firstName: "First Name"},
		{name: "maximum username length", username: strings.Repeat("A", 32), firstName: "First Name"},
		{name: "maximum first_name rune length", username: "Username_1", firstName: strings.Repeat("я", 64)},
		{name: "maximum last_name rune length", username: "Username_1", firstName: "First Name", lastName: new(strings.Repeat("я", 64))},
		{name: "maximum bio rune length", username: "Username_1", firstName: "First Name", bio: new(strings.Repeat("я", 70))},
	}

	for _, tt := range validCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUserProfile(tt.username, tt.firstName, tt.lastName, tt.bio)
			require.NoError(t, err)
		})
	}

	invalidCases := []struct {
		name      string
		username  string
		firstName string
		lastName  *string
		bio       *string
	}{
		{name: "username shorter than minimum", username: "Ab_1", firstName: "First Name"},
		{name: "username longer than maximum", username: strings.Repeat("A", 33), firstName: "First Name"},
		{name: "username with forbidden character", username: "User-name", firstName: "First Name"},
		{name: "blank first_name", username: "Username_1", firstName: "  "},
		{name: "first_name longer than maximum", username: "Username_1", firstName: strings.Repeat("я", 65)},
		{name: "last_name longer than maximum", username: "Username_1", firstName: "First Name", lastName: new(strings.Repeat("я", 65))},
		{name: "bio longer than maximum", username: "Username_1", firstName: "First Name", bio: new(strings.Repeat("я", 71))},
	}

	for _, tt := range invalidCases {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := NewUserProfile(tt.username, tt.firstName, tt.lastName, tt.bio)

			require.ErrorIs(t, err, ErrInvalidUserProfile)
			require.Zero(t, profile)
		})
	}
}

func TestUserProfileOptionalFieldsAreImmutable(t *testing.T) {
	profile, err := NewUserProfile(
		"Username_1",
		"First Name",
		new("Last Name"),
		new("Bio"),
	)
	require.NoError(t, err)

	lastName := profile.LastName()
	bio := profile.Bio()
	*lastName = "Changed Last Name"
	*bio = "Changed Bio"

	require.Equal(t, "Last Name", *profile.LastName())
	require.Equal(t, "Bio", *profile.Bio())
}
