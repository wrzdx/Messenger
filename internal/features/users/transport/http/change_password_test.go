package users_transport_http

import (
	"bytes"
	"encoding/json"
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

func TestChangePassword(t *testing.T) {
	tests := []struct {
		name              string
		userID            uuid.UUID
		body              ChangePasswordRequest
		serviceErr        error
		wantServiceCalled bool
		wantStatus        int
		wantError         string
	}{
		{
			name:              "valid password",
			userID:            core_test_utils.ID,
			wantServiceCalled: true,
			wantStatus:        http.StatusNoContent,
			body: ChangePasswordRequest{
				OldPassword: "old_password",
				NewPassword: "new_password",
			},
		},
		{
			name:              "invalid credentials",
			userID:            core_test_utils.ID,
			wantServiceCalled: true,
			serviceErr:        domain.ErrInvalidCredentials,
			wantStatus:        http.StatusUnauthorized,
			wantError:         core_http_response.MapError(domain.ErrInvalidCredentials).Message,
			body: ChangePasswordRequest{
				OldPassword: "wrong_password",
				NewPassword: "new_password",
			},
		},
		{
			name:              "invalid password",
			userID:            core_test_utils.ID,
			wantServiceCalled: true,
			serviceErr:        domain.ErrInvalidPassword,
			wantStatus:        http.StatusBadRequest,
			wantError:         core_http_response.MapError(domain.ErrInvalidPassword).Message,
			body: ChangePasswordRequest{
				OldPassword: "old_password",
				NewPassword: "123",
			},
		},
		{
			name:       "missing old password",
			wantStatus: http.StatusBadRequest,
			wantError:  core_http_response.ErrInvalidArgument.Error(),
			body: ChangePasswordRequest{
				NewPassword: "new_password",
			},
		},
		{
			name:       "missing new password",
			wantStatus: http.StatusBadRequest,
			wantError:  core_http_response.ErrInvalidArgument.Error(),
			body: ChangePasswordRequest{
				OldPassword: "old_password",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				serviceCalled     bool
				serviceGotID      uuid.UUID
				serviceGotOldPass string
				serviceGotNewPass string
			)

			service := StubUsersService{
				ChangePasswordFn: func(
					id uuid.UUID,
					oldPassword string,
					newPassword string,
				) error {
					serviceCalled = true
					serviceGotID = id
					serviceGotOldPass = oldPassword
					serviceGotNewPass = newPassword

					return tt.serviceErr
				},
			}

			handler := NewUsersHTTPHandler(&service)

			body, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatal(err)
			}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodPut,
				"/users/me/password",
				bytes.NewReader(body),
			)

			ctx := core_test_utils.GetLoggerContext(req.Context())
			ctx = core_auth.WithUserID(ctx, tt.userID)

			handler.ChangePassword(rec, req.WithContext(ctx))

			if serviceCalled != tt.wantServiceCalled {
				t.Fatalf(
					"service called = %v, want %v",
					serviceCalled,
					tt.wantServiceCalled,
				)
			}

			if serviceCalled {
				if diff := cmp.Diff(tt.userID, serviceGotID); diff != "" {
					t.Fatalf("userID mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.body.OldPassword, serviceGotOldPass); diff != "" {
					t.Fatalf("oldPassword mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.body.NewPassword, serviceGotNewPass); diff != "" {
					t.Fatalf("newPassword mismatch (-want +got):\n%s", diff)
				}
			}

			if rec.Code != tt.wantStatus {
				t.Fatalf("got status %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantError != "" {
				var gotError core_http_response.ErrorResponse

				if err := json.NewDecoder(rec.Body).Decode(&gotError); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if !strings.HasSuffix(gotError.Error, tt.wantError) {
					t.Fatalf(
						"ErrorResponse mismatch:\nwant: %s\ngot: %s",
						tt.wantError,
						gotError.Error,
					)
				}
			} else if rec.Body.Len() != 0 {
				t.Fatalf("expected empty response body, got %q", rec.Body.String())
			}
		})
	}
}
