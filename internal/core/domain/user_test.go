package domain

import (
	"errors"
	core_errors "messenger/internal/core/errors.go"
	"strings"
	"testing"
)

func TestUserValidate(t *testing.T) {
	lastName := "Smith"
	longLastName := strings.Repeat("a", 65)

	bio := "Hello!"
	longBio := strings.Repeat("a", 71)

	tests := []struct {
		name string
		user User
	}{
		{
			name: "valid",
			user: User{
				Username:  "username",
				FirstName: "Andrew",
				LastName:  &lastName,
				Bio:       &bio,
			},
		},
		{
			name: "short username",
			user: User{
				Username:  "usr",
				FirstName: "Andrew",
			},
		},
		{
			name: "short first name",
			user: User{
				Username:  "username",
				FirstName: "Tom",
			},
		},
		{
			name: "last name too long",
			user: User{
				Username:  "username",
				FirstName: "Andrew",
				LastName:  &longLastName,
			},
		},
		{
			name: "bio too long",
			user: User{
				Username:  "username",
				FirstName: "Andrew",
				Bio:       &longBio,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if tt.name == "valid" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatal("expected validation error")
			}

			if !errors.Is(err, core_errors.ErrInvalidArgument) {
				t.Fatalf("expected ErrInvalidArgument, got %v", err)
			}
		})
	}
}

func TestUserPatchValidate(t *testing.T) {
	username := "new_username"

	tests := []struct {
		name  string
		patch UserPatch
		valid bool
	}{
		{
			name: "valid",
			patch: UserPatch{
				Username: Nullable[string]{
					Set:   true,
					Value: &username,
				},
			},
			valid: true,
		},
		{
			name: "username to null",
			patch: UserPatch{
				Username: Nullable[string]{
					Set: true,
				},
			},
		},
		{
			name: "first name to null",
			patch: UserPatch{
				FirstName: Nullable[string]{
					Set: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.patch.Validate()

			if tt.valid {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatal("expected validation error")
			}

			if !errors.Is(err, core_errors.ErrInvalidArgument) {
				t.Fatalf("expected ErrInvalidArgument, got %v", err)
			}
		})
	}
}

func TestUserApplyPatch(t *testing.T) {
	user := User{
		Username:  "username",
		FirstName: "Andrew",
	}

	newUsername := "new_username"
	newBio := "Hello"

	t.Run("success", func(t *testing.T) {
		u := user

		err := u.ApplyPatch(UserPatch{
			Username: Nullable[string]{
				Set:   true,
				Value: &newUsername,
			},
			Bio: Nullable[string]{
				Set:   true,
				Value: &newBio,
			},
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if u.Username != newUsername {
			t.Fatalf("username wasn't updated")
		}

		if u.Bio == nil || *u.Bio != newBio {
			t.Fatalf("bio wasn't updated")
		}
	})

	t.Run("invalid patch", func(t *testing.T) {
		u := user

		err := u.ApplyPatch(UserPatch{
			Username: Nullable[string]{
				Set: true,
			},
		})

		if err == nil {
			t.Fatal("expected error")
		}

		if u.Username != user.Username {
			t.Fatal("user should not be modified")
		}
	})

	t.Run("patched user invalid", func(t *testing.T) {
		u := user

		short := "abc"

		err := u.ApplyPatch(UserPatch{
			Username: Nullable[string]{
				Set:   true,
				Value: &short,
			},
		})

		if err == nil {
			t.Fatal("expected error")
		}

		if u.Username != user.Username {
			t.Fatal("user should not be modified")
		}
	})
}
