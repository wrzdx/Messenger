package domain

import (
	"errors"
	core_errors "messenger/internal/core/errors.go"
	"testing"
)

func TestUserCredentialsValidate(t *testing.T) {
	tests := []struct {
		name    string
		creds   UserCredentials
		wantErr bool
	}{
		{
			name: "valid credentials",
			creds: UserCredentials{
				Username: "username",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "password too short",
			creds: UserCredentials{
				Username: "username",
				Password: "short",
			},
			wantErr: true,
		},
		{
			name: "password too long",
			creds: UserCredentials{
				Username: "username",
				Password: "123456789012345678901234567890123",
			},
			wantErr: true,
		},
		{
			name: "username too short",
			creds: UserCredentials{
				Username: "usr",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "username too long",
			creds: UserCredentials{
				Username: "123456789012345678901234567890123",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "both invalid",
			creds: UserCredentials{
				Username: "usr",
				Password: "123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.creds.Validate()

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if !errors.Is(err, core_errors.ErrInvalidArgument) {
					t.Fatalf("expected ErrInvalidArgument, got %v", err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
