package users_transport_http

import (
	"encoding/json"
	"fmt"
	"messenger/internal/core/domain"
	core_http_response "messenger/internal/core/transport/http/response"
	core_test_utils "messenger/internal/core/utils/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestGetMe(t *testing.T) {
	user := core_test_utils.Users[0]
	tests := []struct {
		name              string
		serviceUser       domain.User
		serviceErr        error
		wantUser          UserDTOResponse
		userID            uuid.UUID
		wantServiceCalled bool
		wantStatus        int
		wantError         string
		withoutClaims     bool
	}{
		{
			name:        "existing user",
			serviceUser: user,
			wantUser: UserDTOResponse{
				ID:        user.ID,
				Username:  user.Username,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				CreatedAt: user.CreatedAt,
				Bio:       user.Bio,
			},
			userID:            user.ID,
			wantServiceCalled: true,
			wantStatus:        http.StatusOK,
		},
		{
			name:              "non-existing user",
			wantServiceCalled: true,
			serviceErr:        domain.ErrUserNotFound,
			wantStatus:        http.StatusNotFound,

			wantError: core_http_response.MapError(domain.ErrUserNotFound).Message,
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
			var serviceCalled bool
			var serviceGotID uuid.UUID
			service := StubUsersService{
				GetUserFn: func(id uuid.UUID) (domain.User, error) {
					serviceCalled = true
					serviceGotID = id

					return tt.serviceUser, tt.serviceErr
				},
			}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/users/%v", tt.userID),
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
			handler.GetMe(rec, req.WithContext(ctx))

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
				var gotResponse UserDTOResponse
				if err := json.NewDecoder(rec.Body).Decode(&gotResponse); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if diff := cmp.Diff(tt.wantUser, gotResponse); diff != "" {
					t.Fatalf("GetUsersResponse mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
