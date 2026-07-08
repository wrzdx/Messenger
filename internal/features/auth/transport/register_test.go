package auth_transport_http

import (
	"bytes"
	"encoding/json"
	"messenger/internal/core/domain"
	http_response "messenger/internal/core/transport/http/response"
	test_utils "messenger/internal/core/utils/test"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name              string
		serviceErr        error
		wantServiceCalled bool
		wantStatus        int
		wantError         string
		body              RegisterRequest
	}{
		{
			name:              "valid user",
			wantServiceCalled: true,
			wantStatus:        http.StatusCreated,
			body: RegisterRequest{
				Username:  "i.ivanov",
				FirstName: "Ivan",
				LastName:  new("Ivanov"),
				Bio:       new("I like pizza!"),
				Password:  "password",
			},
		},
		{
			name:              "username already exists",
			wantServiceCalled: true,
			serviceErr:        domain.ErrUserAlreadyExists,
			wantStatus:        http.StatusConflict,
			wantError:         http_response.MapError(domain.ErrUserAlreadyExists).Message,
			body: RegisterRequest{
				Username:  "i.ivanov",
				FirstName: "Ivan",
				LastName:  new("Ivanov"),
				Bio:       new("I like pizza!"),
				Password:  "password",
			},
		},
		{
			name:       "missing username",
			wantStatus: http.StatusBadRequest,
			wantError:  "username is required",
			body: RegisterRequest{
				FirstName: "Ivan",
				Password:  "password",
			},
		},
		{
			name:              "username too short",
			wantStatus:        http.StatusBadRequest,
			wantServiceCalled: true,
			serviceErr:        domain.ErrInvalidUsername,
			wantError:         http_response.MapError(domain.ErrInvalidUsername).Message,
			body: RegisterRequest{
				Username:  "ivan",
				FirstName: "Ivan",
				Password:  "password",
			},
		},
		{
			name:              "username too long",
			wantStatus:        http.StatusBadRequest,
			wantServiceCalled: true,
			serviceErr:        domain.ErrInvalidUsername,
			wantError:         http_response.MapError(domain.ErrInvalidUsername).Message,
			body: RegisterRequest{
				Username:  strings.Repeat("a", 33),
				FirstName: "Ivan",
				Password:  "password",
			},
		},
		{
			name:       "missing first name",
			wantStatus: http.StatusBadRequest,
			wantError:  "first_name is required",
			body: RegisterRequest{
				Username: "i.ivanov",
				Password: "password",
			},
		},
		{
			name:              "first name too long",
			wantStatus:        http.StatusBadRequest,
			wantServiceCalled: true,
			serviceErr:        domain.ErrInvalidFirstName,
			wantError:         http_response.MapError(domain.ErrInvalidFirstName).Message,
			body: RegisterRequest{
				Username:  "i.ivanov",
				FirstName: strings.Repeat("a", 65),
				Password:  "password",
			},
		},
		{
			name:              "last name too long",
			wantStatus:        http.StatusBadRequest,
			wantServiceCalled: true,
			serviceErr:        domain.ErrInvalidLastName,
			wantError:         http_response.MapError(domain.ErrInvalidLastName).Message,
			body: RegisterRequest{
				Username:  "i.ivanov",
				FirstName: "Ivan",
				LastName:  new(strings.Repeat("a", 65)),
				Password:  "password",
			},
		},
		{
			name:              "bio too long",
			wantStatus:        http.StatusBadRequest,
			wantServiceCalled: true,
			serviceErr:        domain.ErrInvalidBio,
			wantError:         http_response.MapError(domain.ErrInvalidBio).Message,
			body: RegisterRequest{
				Username:  "i.ivanov",
				FirstName: "Ivan",
				Bio:       new(strings.Repeat("a", 71)),
				Password:  "password",
			},
		},
		{
			name:              "password too short",
			wantStatus:        http.StatusBadRequest,
			wantServiceCalled: true,
			serviceErr:        domain.ErrInvalidPassword,
			wantError:         http_response.MapError(domain.ErrInvalidPassword).Message,
			body: RegisterRequest{
				Username:  "i.ivanov",
				FirstName: "Ivan",
				Password:  "pass",
			},
		},
		{
			name:              "password too long",
			wantStatus:        http.StatusBadRequest,
			wantServiceCalled: true,
			serviceErr:        domain.ErrInvalidPassword,
			wantError:         http_response.MapError(domain.ErrInvalidPassword).Message,
			body: RegisterRequest{
				Username:  "i.ivanov",
				FirstName: "Ivan",
				Password:  strings.Repeat("a", 33),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			userResponse := UserResponse{
				ID:        test_utils.ID,
				Username:  tt.body.Username,
				FirstName: tt.body.FirstName,
				LastName:  tt.body.LastName,
				CreatedAt: test_utils.CreatedAt,
				Bio:       tt.body.Bio,
			}

			want := RegisterResponse{
				User: userResponse,
			}

			wantServiceGotPayload := domain.NewRegisterUserPayload(
				tt.body.Username,
				tt.body.FirstName,
				tt.body.LastName,
				tt.body.Bio,
				tt.body.Password,
			)
			var servicePayload domain.RegisterUserPayload
			var serviceCalled bool
			service := StubAuthService{
				CreateUserFn: func(
					payload domain.RegisterUserPayload,
				) (domain.User, domain.TokenPair, error) {
					serviceCalled = true
					servicePayload = domain.NewRegisterUserPayload(
						payload.Username,
						payload.FirstName,
						payload.LastName,
						payload.Bio,
						payload.Password,
					)
					return domain.NewUser(
						test_utils.ID,
						tt.body.Username,
						tt.body.FirstName,
						tt.body.LastName,
						test_utils.CreatedAt,
						tt.body.Bio,
						test_utils.PasswordHash,
					), domain.TokenPair{}, tt.serviceErr
				},
			}

			handler := NewAuthHTTPHandler(&service, &StubCookieManager{})
			body, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatal(err)
			}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("", "/users", bytes.NewReader(body))
			ctx := test_utils.GetLoggerContext(req.Context())

			// Action
			handler.Register(rec, req.WithContext(ctx))

			// Check
			if serviceCalled != tt.wantServiceCalled {
				t.Fatalf(
					"service called = %v, want %v",
					serviceCalled,
					tt.wantServiceCalled,
				)
			}
			if serviceCalled {
				if diff := cmp.Diff(wantServiceGotPayload, servicePayload); diff != "" {
					t.Fatalf("ServiceGotUser mismatch (-want +got):\n%s", diff)
				}
			}

			if rec.Code != tt.wantStatus {
				t.Fatalf("got status %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantError != "" {
				var gotError http_response.ErrorResponse
				if err := json.NewDecoder(rec.Body).Decode(&gotError); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if gotError.Error != tt.wantError {
					t.Fatalf(
						"ErrorResponse mismatch:\nwant: %s\ngot: %s",
						tt.wantError,
						gotError.Error,
					)
				}
			} else {
				var gotResponse RegisterResponse
				if err := json.NewDecoder(rec.Body).Decode(&gotResponse); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				gotResponse.User.CreatedAt = test_utils.CreatedAt

				if diff := cmp.Diff(want, gotResponse); diff != "" {
					t.Fatalf("CreateUserResponse mismatch (-want +got):\n%s", diff)
				}
			}

		})
	}
}
