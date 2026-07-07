package users_transport_http

import (
	"encoding/json"
	"messenger/internal/core/domain"
	core_http_response "messenger/internal/core/transport/http/response"
	core_test_utils "messenger/internal/core/utils/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestDeleteMe(t *testing.T) {
	tests := []struct {
		name              string
		userID            uuid.UUID
		serviceErr        error
		wantServiceCalled bool
		wantStatus        int
		wantError         string
		withoutClaims     bool
	}{
		{
			name:              "existing user",
			wantServiceCalled: true,
			wantStatus:        http.StatusNoContent,
		},
		{
			name:              "non-existing user",
			serviceErr:        domain.ErrUserNotFound,
			wantServiceCalled: true,
			wantStatus:        http.StatusNotFound,
			wantError:         core_http_response.MapError(domain.ErrUserNotFound).Message,
		},
		{
			name:          "without claims",
			withoutClaims: true,
			wantStatus:    http.StatusUnauthorized,
			wantError:     core_http_response.MapError(core_http_response.ErrMissingClaims).Message,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				serviceCalled bool
				serviceGotID  uuid.UUID
			)

			service := StubUsersService{
				DeleteUserFn: func(id uuid.UUID) error {
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

			claims := domain.Claims{
				UserID: tt.userID,
			}

			ctx := core_test_utils.GetLoggerContext(req.Context())
			if !tt.withoutClaims {
				ctx = core_test_utils.GetClaimsContext(ctx, claims)
			}

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

			if tt.wantError != "" {
				var gotError core_http_response.ErrorResponse
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
				if rec.Body.Len() != 0 {
					t.Fatalf("expected empty response body, got %q", rec.Body.String())
				}
			}
		})
	}
}
