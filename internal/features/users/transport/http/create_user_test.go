package users_transport_http

import (
	"bytes"
	"encoding/json"
	"messenger/internal/core/domain"
	core_errors "messenger/internal/core/errors"
	core_logger "messenger/internal/core/logger"
	core_http_response "messenger/internal/core/transport/http/response"
	core_test_utils "messenger/internal/core/utils/test"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var tests = []struct {
	name       string
	serviceErr error
	wantStatus int
	wantError  error
	body       CreateUserRequest
}{
	{
		name:       "valid user",
		wantStatus: http.StatusCreated,
		body: CreateUserRequest{
			Username:  "i.ivanov",
			FirstName: "Ivan",
			LastName:  new("Ivanov"),
			Bio:       new("I like pizza!"),
			Password:  "password",
		},
	},
	{
		name:       "username already exists",
		serviceErr: core_errors.ErrConflict,
		wantStatus: http.StatusConflict,
		wantError:  core_errors.ErrConflict,
		body: CreateUserRequest{
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
		wantError:  core_errors.ErrInvalidArgument,
		body: CreateUserRequest{
			FirstName: "Ivan",
			Password:  "password",
		},
	},
	{
		name:       "username too short",
		wantStatus: http.StatusBadRequest,
		wantError:  core_errors.ErrInvalidArgument,
		body: CreateUserRequest{
			Username:  "ivan",
			FirstName: "Ivan",
			Password:  "password",
		},
	},
	{
		name:       "username too long",
		wantStatus: http.StatusBadRequest,
		wantError:  core_errors.ErrInvalidArgument,
		body: CreateUserRequest{
			Username:  strings.Repeat("a", 33),
			FirstName: "Ivan",
			Password:  "password",
		},
	},
	{
		name:       "missing first name",
		wantStatus: http.StatusBadRequest,
		wantError:  core_errors.ErrInvalidArgument,
		body: CreateUserRequest{
			Username: "i.ivanov",
			Password: "password",
		},
	},
	{
		name:       "first name too long",
		wantStatus: http.StatusBadRequest,
		wantError:  core_errors.ErrInvalidArgument,
		body: CreateUserRequest{
			Username:  "i.ivanov",
			FirstName: strings.Repeat("a", 65),
			Password:  "password",
		},
	},
	{
		name:       "last name too long",
		wantStatus: http.StatusBadRequest,
		wantError:  core_errors.ErrInvalidArgument,
		body: CreateUserRequest{
			Username:  "i.ivanov",
			FirstName: "Ivan",
			LastName:  new(strings.Repeat("a", 65)),
			Password:  "password",
		},
	},
	{
		name:       "bio too long",
		wantStatus: http.StatusBadRequest,
		wantError:  core_errors.ErrInvalidArgument,
		body: CreateUserRequest{
			Username:  "i.ivanov",
			FirstName: "Ivan",
			Bio:       new(strings.Repeat("a", 71)),
			Password:  "password",
		},
	},
	{
		name:       "password too short",
		wantStatus: http.StatusBadRequest,
		wantError:  core_errors.ErrInvalidArgument,
		body: CreateUserRequest{
			Username:  "i.ivanov",
			FirstName: "Ivan",
			Password:  "pass",
		},
	},
	{
		name:       "password too long",
		wantStatus: http.StatusBadRequest,
		wantError:  core_errors.ErrInvalidArgument,
		body: CreateUserRequest{
			Username:  "i.ivanov",
			FirstName: "Ivan",
			Password:  strings.Repeat("a", 33),
		},
	},
}

func TestCreateUser(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			want := CreateUserResponse{
				ID:        core_test_utils.ID,
				Username:  tt.body.Username,
				FirstName: tt.body.FirstName,
				LastName:  tt.body.LastName,
				CreatedAt: core_test_utils.CreatedAt,
				Bio:       tt.body.Bio,
			}

			WantServiceGotUser := domain.NewUserUninitialized(
				tt.body.Username,
				tt.body.FirstName,
				tt.body.LastName,
				tt.body.Bio,
			)
			WantServiceGotUser.ID = core_test_utils.ID
			WantServiceGotUser.CreatedAt = core_test_utils.CreatedAt
			serviceGotCreds := domain.NewCredentials(
				tt.body.Username,
				tt.body.Password,
			)
			service := StubUsersService{
				ReturnUser: domain.NewUser(
					core_test_utils.ID,
					tt.body.Username,
					tt.body.FirstName,
					tt.body.LastName,
					core_test_utils.CreatedAt,
					tt.body.Bio,
				),
				ReturnError: tt.serviceErr,
			}

			handler := NewUsersHTTPHandler(&service)
			body, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatal(err)
			}
			log := core_logger.NewTestLogger()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("", "/", bytes.NewReader(body))
			ctx := core_logger.ToContext(req.Context(), log)

			// Action
			handler.CreateUser(rec, req.WithContext(ctx))

			// Check
			if service.Called {
				service.GotUser.ID = core_test_utils.ID
				service.GotUser.CreatedAt = core_test_utils.CreatedAt
				if diff := cmp.Diff(WantServiceGotUser, service.GotUser); diff != "" {
					t.Fatalf("ServiceGotUser mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(serviceGotCreds, service.GotCreds); diff != "" {
					t.Fatalf("ServiceGotCreds mismatch (-want +got):\n%s", diff)
				}
			}

			if tt.wantError != nil {
				var gotError core_http_response.ErrorResponse
				if err := json.NewDecoder(rec.Body).Decode(&gotError); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if !strings.HasSuffix(gotError.Error, tt.wantError.Error()) {
					t.Fatalf(
						"ErrorResponse mismatch:\nwant: %s\ngot: %s",
						tt.wantError.Error(),
						gotError.Error,
					)
				}
			} else {
				if rec.Code != http.StatusCreated {
					t.Fatalf("got status %d, want %d", rec.Code, http.StatusCreated)
				}
				var gotResponse CreateUserResponse
				if err := json.NewDecoder(rec.Body).Decode(&gotResponse); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				gotResponse.CreatedAt = core_test_utils.CreatedAt

				if diff := cmp.Diff(want, gotResponse); diff != "" {
					t.Fatalf("CreateUserResponse mismatch (-want +got):\n%s", diff)
				}
			}

		})
	}
}

func Ptr(s string) {
	panic("unimplemented")
}
