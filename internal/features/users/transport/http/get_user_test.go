package users_transport_http

import (
	"encoding/json"
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

func TestGetUser(t *testing.T) {
	user := core_test_utils.Users[0]
	tests := []struct {
		name              string
		userID            string
		serviceUser       domain.User
		serviceErr        error
		wantUser          UserDTOResponse
		wantServiceUserId uuid.UUID
		wantServiceCalled bool
		wantStatus        int
		wantError         error
	}{
		{
			name:        "existing user",
			userID:      user.ID.String(),
			serviceUser: user,
			wantUser: UserDTOResponse{
				ID:        user.ID,
				Username:  user.Username,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				CreatedAt: user.CreatedAt,
				Bio:       user.Bio,
			},
			wantServiceUserId: user.ID,
			wantServiceCalled: true,
			wantStatus:        http.StatusOK,
		},
		{
			name:              "non-existing user",
			userID:            core_test_utils.ID.String(),
			wantServiceUserId: core_test_utils.ID,
			wantServiceCalled: true,
			serviceErr:        domain.ErrUserNotFound,
			wantStatus:        http.StatusNotFound,
			wantError:         domain.ErrUserNotFound,
		},
		{
			name:       "invalid user id",
			userID:     "asdf",
			wantStatus: http.StatusBadRequest,
			wantError:  ErrInvalidArgument,
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
				"/users/"+tt.userID,
				nil,
			)
			req.SetPathValue("id", tt.userID)
			ctx := core_test_utils.GetLoggerContext(req.Context())
			handler := NewUsersHTTPHandler(&service)

			// action
			handler.GetUser(rec, req.WithContext(ctx))

			// check
			if serviceCalled != tt.wantServiceCalled {
				t.Fatalf(
					"service called = %v, want %v",
					serviceCalled,
					tt.wantServiceCalled,
				)
			}

			if serviceCalled {
				if diff := cmp.Diff(tt.wantServiceUserId, serviceGotID); diff != "" {
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

				if !strings.HasPrefix(gotError.Error, tt.wantError.Error()) {
					t.Fatalf(
						"ErrorResponse mismatch:\nwant: %s\ngot: %s",
						tt.wantError.Error(),
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
