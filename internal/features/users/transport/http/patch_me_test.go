package users_transport_http

import (
	"bytes"
	"encoding/json"
	"errors"
	core_auth "messenger/internal/core/auth"
	"messenger/internal/core/domain"
	core_http_response "messenger/internal/core/transport/http/response"
	core_test_utils "messenger/internal/core/utils/test"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestPatchMe(t *testing.T) {
	user := core_test_utils.Users[0]

	tests := []struct {
		name string
		body map[string]any

		serviceErr        error
		wantServiceCalled bool
		wantServicePatch  domain.UserPatch
		wantStatus        int
		wantError         error
		wantResponse      PatchUserResponse
	}{
		{
			name:              "patch username",
			wantServiceCalled: true,
			wantStatus:        http.StatusOK,
			body: map[string]any{
				"username": "new_username",
			},
			wantServicePatch: domain.UserPatch{
				Username: domain.Nullable[string]{
					Value: new("new_username"),
					Set:   true,
				},
			},
			wantResponse: PatchUserResponse{
				ID:        user.ID,
				Username:  "new_username",
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Bio:       user.Bio,
				CreatedAt: user.CreatedAt,
			},
		},
		{
			name:              "patch all fields",
			wantServiceCalled: true,
			wantStatus:        http.StatusOK,
			body: map[string]any{
				"username":   "new_username",
				"first_name": "Ivan",
				"last_name":  "Ivanov",
				"Bio":        "I like pizza!",
			},
			wantServicePatch: domain.UserPatch{
				Username: domain.Nullable[string]{
					Value: new("new_username"),
					Set:   true,
				},
				FirstName: domain.Nullable[string]{
					Value: new("Ivan"),
					Set:   true,
				},
				LastName: domain.Nullable[string]{
					Value: new("Ivanov"),
					Set:   true,
				},
				Bio: domain.Nullable[string]{
					Value: new("I like pizza!"),
					Set:   true,
				},
			},
			wantResponse: PatchUserResponse{
				ID:        user.ID,
				Username:  "new_username",
				FirstName: "Ivan",
				LastName:  new("Ivanov"),
				Bio:       new("I like pizza!"),
				CreatedAt: user.CreatedAt,
			},
		},
		{
			name:       "username too short",
			wantStatus: http.StatusBadRequest,
			wantError:  ErrInvalidArgument,
			body: map[string]any{
				"username": "abc",
			},
		},
		{
			name:       "username too long",
			wantStatus: http.StatusBadRequest,
			wantError:  ErrInvalidArgument,
			body: map[string]any{
				"username": strings.Repeat("a", 33),
			},
		},
		{
			name:       "username is null",
			wantStatus: http.StatusBadRequest,
			wantError:  ErrInvalidArgument,
			body: map[string]any{
				"username": nil,
			},
		},
		{
			name:       "first name is null",
			wantStatus: http.StatusBadRequest,
			wantError:  ErrInvalidArgument,
			body: map[string]any{
				"first_name": nil,
			},
		},
		{
			name:              "service error",
			wantServiceCalled: true,
			wantStatus:        http.StatusInternalServerError,
			wantError:         errors.New("service error"),
			serviceErr:        errors.New("service error"),
			body: map[string]any{
				"username": "new_username",
			},
			wantServicePatch: domain.UserPatch{
				Username: domain.Nullable[string]{
					Value: new("new_username"),
					Set:   true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				serviceCalled bool
				serviceGotID  uuid.UUID
				servicePatch  domain.UserPatch
			)

			service := StubUsersService{
				PatchUserFn: func(
					id uuid.UUID,
					patch domain.UserPatch,
				) (domain.User, error) {

					serviceCalled = true
					serviceGotID = id
					servicePatch = patch

					return domain.NewUser(
						user.ID,
						tt.wantResponse.Username,
						tt.wantResponse.FirstName,
						tt.wantResponse.LastName,
						user.CreatedAt,
						tt.wantResponse.Bio,
						core_test_utils.PasswordHash,
					), tt.serviceErr
				},
			}

			handler := NewUsersHTTPHandler(&service)
			data, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(
				http.MethodPatch,
				"/users/me",
				bytes.NewReader(data),
			)

			claims := core_auth.Claims{
				UserID: user.ID,
			}

			ctx := core_test_utils.GetLoggerContext(req.Context())
			ctx = core_test_utils.GetClaimsContext(ctx, claims)

			rec := httptest.NewRecorder()

			// action
			handler.PatchMe(rec, req.WithContext(ctx))

			// check

			if serviceCalled != tt.wantServiceCalled {
				t.Fatalf(
					"service called = %v, want %v",
					serviceCalled,
					tt.wantServiceCalled,
				)
			}

			if serviceCalled {
				if diff := cmp.Diff(user.ID, serviceGotID); diff != "" {
					t.Fatalf("userID mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantServicePatch, servicePatch); diff != "" {
					t.Fatalf("patch mismatch (-want +got):\n%s", diff)
				}
			}

			if rec.Code != tt.wantStatus {
				t.Fatalf("got status %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantError != nil {
				var gotError core_http_response.ErrorResponse

				if err := json.NewDecoder(rec.Body).Decode(&gotError); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if !strings.HasPrefix(gotError.Error, tt.wantError.Error()) {
					t.Fatalf(
						"ErrorResponse mismatch:\nwant: %s\ngot: %s",
						tt.wantError.Error(),
						gotError.Error,
					)
				}
			} else {
				var gotResponse PatchUserResponse

				if err := json.NewDecoder(rec.Body).Decode(&gotResponse); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				gotResponse.CreatedAt = core_test_utils.CreatedAt

				if diff := cmp.Diff(tt.wantResponse, gotResponse); diff != "" {
					t.Fatalf("PatchMeResponse mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
