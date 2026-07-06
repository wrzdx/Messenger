package users_transport_http

import (
	"encoding/json"
	core_auth "messenger/internal/core/auth"
	core_errors "messenger/internal/core/errors"
	core_http_response "messenger/internal/core/transport/http/response"
	core_test_utils "messenger/internal/core/utils/test"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDeleteMe(t *testing.T) {
	tests := []struct {
		name              string
		userID            int
		serviceErr        error
		wantServiceCalled bool
		wantStatus        int
		wantError         error
	}{
		{
			name:              "existing user",
			userID:            1,
			wantServiceCalled: true,
			wantStatus:        http.StatusNoContent,
		},
		{
			name:              "non-existing user",
			userID:            -1,
			serviceErr:        core_errors.ErrorNotFound,
			wantServiceCalled: true,
			wantStatus:        http.StatusNotFound,
			wantError:         core_errors.ErrorNotFound,
		},
		{
			name:              "service error",
			userID:            1,
			serviceErr:        core_errors.ErrInternalServer,
			wantServiceCalled: true,
			wantStatus:        http.StatusInternalServerError,
			wantError:         core_errors.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				serviceCalled bool
				serviceGotID  int
			)

			service := StubUsersService{
				DeleteUserFn: func(id int) error {
					serviceCalled = true
					serviceGotID = id

					return tt.serviceErr
				},
			}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodDelete,
				"/users/me",
				nil,
			)

			claims := core_auth.Claims{
				UserID: tt.userID,
			}

			ctx := core_test_utils.GetLoggerContext(req.Context())
			ctx = core_test_utils.GetClaimsContext(ctx, claims)

			handler := NewUsersHTTPHandler(&service)

			// action
			handler.DeleteMe(rec, req.WithContext(ctx))

			// check
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
			}

			if rec.Code != tt.wantStatus {
				t.Fatalf("got status %d, want %d", rec.Code, tt.wantStatus)
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
				if rec.Body.Len() != 0 {
					t.Fatalf("expected empty response body, got %q", rec.Body.String())
				}
			}
		})
	}
}
