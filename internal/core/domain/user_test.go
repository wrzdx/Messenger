package domain

import (
	"errors"
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
		err  error
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
			err: ErrInvalidUsername,
		},
		{
			name: "missing first name",
			user: User{
				Username:  "username",
				FirstName: "",
			},
			err: ErrInvalidFirstName,
		},
		{
			name: "last name too long",
			user: User{
				Username:  "username",
				FirstName: "Andrew",
				LastName:  &longLastName,
			},
			err: ErrInvalidLastName,
		},
		{
			name: "bio too long",
			user: User{
				Username:  "username",
				FirstName: "Andrew",
				Bio:       &longBio,
			},
			err: ErrInvalidBio,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.Join(tt.user.Validate()...)

			if tt.name == "valid" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatal("expected validation error")
			}

			if !errors.Is(err, tt.err) {
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
		err   error
	}{
		{
			name: "valid",
			patch: UserPatch{
				Username: Nullable[string]{
					Set:   true,
					Value: &username,
				},
			},
		},
		{
			name: "username to null",
			patch: UserPatch{
				Username: Nullable[string]{
					Set: true,
				},
			},
			err: ErrNullUsername,
		},
		{
			name: "first name to null",
			patch: UserPatch{
				FirstName: Nullable[string]{
					Set: true,
				},
			},
			err: ErrNullFirstname,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.Join(tt.patch.Validate()...)

			if !errors.Is(err, tt.err) {
				t.Fatalf("want: %v, got %v", tt.err, err)
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
