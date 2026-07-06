package users_transport_http

import (
	"encoding/json"
	"messenger/internal/core/domain"
	core_errors "messenger/internal/core/errors"
	core_http_response "messenger/internal/core/transport/http/response"
	core_test_utils "messenger/internal/core/utils/test"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetUsers(t *testing.T) {
	tests := []struct {
		name              string
		limit             string
		offset            string
		serviceUsers      []domain.User
		serviceErr        error
		wantUsers         []UserDTOResponse
		wantServiceLimit  *int
		wantServiceOffset *int
		wantServiceCalled bool
		wantStatus        int
		wantError         error
	}{
		{
			name:              "return all users",
			wantServiceCalled: true,
			wantStatus:        http.StatusOK,
			serviceUsers:      core_test_utils.Users,
			wantUsers: []UserDTOResponse{
				{
					ID:        1,
					Username:  "user_1",
					FirstName: "Username",
					LastName:  new("1"),
					CreatedAt: core_test_utils.CreatedAt,
					Bio:       new("I'm user 1"),
				},
				{
					ID:        2,
					Username:  "user_2",
					FirstName: "Username",
					LastName:  new("2"),
					CreatedAt: core_test_utils.CreatedAt,
					Bio:       new("I'm user 2"),
				},
				{
					ID:        3,
					Username:  "user_3",
					FirstName: "Username",
					LastName:  new("3"),
					CreatedAt: core_test_utils.CreatedAt,
					Bio:       new("I'm user 3"),
				},
			},
		},
		{
			name:              "limit users",
			limit:             "1",
			wantServiceLimit:  new(1),
			serviceUsers:      core_test_utils.Users[:1],
			wantServiceCalled: true,
			wantStatus:        http.StatusOK,
			wantUsers: []UserDTOResponse{
				{
					ID:        1,
					Username:  "user_1",
					FirstName: "Username",
					LastName:  new("1"),
					CreatedAt: core_test_utils.CreatedAt,
					Bio:       new("I'm user 1"),
				},
			},
		},
		{
			name:       "negative limit",
			limit:      "-1",
			wantError:  core_errors.ErrInvalidArgument,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:              "offset users",
			offset:            "1",
			wantServiceOffset: new(1),
			serviceUsers:      core_test_utils.Users[1:],
			wantServiceCalled: true,
			wantStatus:        http.StatusOK,
			wantUsers: []UserDTOResponse{
				{
					ID:        2,
					Username:  "user_2",
					FirstName: "Username",
					LastName:  new("2"),
					CreatedAt: core_test_utils.CreatedAt,
					Bio:       new("I'm user 2"),
				},
				{
					ID:        3,
					Username:  "user_3",
					FirstName: "Username",
					LastName:  new("3"),
					CreatedAt: core_test_utils.CreatedAt,
					Bio:       new("I'm user 3"),
				},
			},
		},
		{
			name:       "negative offset",
			offset:     "-1",
			wantError:  core_errors.ErrInvalidArgument,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:              "limit offset users",
			limit:             "1",
			offset:            "1",
			wantServiceLimit:  new(1),
			wantServiceOffset: new(1),
			serviceUsers:      core_test_utils.Users[1:2],
			wantServiceCalled: true,

			wantStatus: http.StatusOK,
			wantUsers: []UserDTOResponse{
				{
					ID:        2,
					Username:  "user_2",
					FirstName: "Username",
					LastName:  new("2"),
					CreatedAt: core_test_utils.CreatedAt,
					Bio:       new("I'm user 2"),
				},
			},
		},
		{
			name:              "empty users",
			limit:             "1",
			offset:            "2",
			wantServiceLimit:  new(1),
			wantServiceOffset: new(2),
			serviceUsers:      core_test_utils.Users[2:2],
			wantServiceCalled: true,
			wantStatus:        http.StatusOK,
			wantUsers:         []UserDTOResponse{},
		},
		{
			name:       "invalid limit",
			limit:      "asadf",
			wantError:  core_errors.ErrInvalidArgument,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid offset",
			offset:     "asadf",
			wantError:  core_errors.ErrInvalidArgument,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			var (
				serviceCalled    bool
				serviceGotLimit  *int
				serviceGotOffset *int
			)

			service := StubUsersService{
				GetUsersFn: func(limit, offset *int) ([]domain.User, error) {
					serviceCalled = true
					serviceGotLimit = limit
					serviceGotOffset = offset

					return tt.serviceUsers, tt.serviceErr
				},
			}

			handler := NewUsersHTTPHandler(&service)
			rec := httptest.NewRecorder()
			values := url.Values{}
			if tt.limit != "" {
				values.Set("limit", tt.limit)
			}
			if tt.offset != "" {
				values.Set("offset", tt.offset)
			}

			req := httptest.NewRequest(
				http.MethodGet,
				"/users?"+values.Encode(),
				nil,
			)

			ctx := core_test_utils.GetLoggerContext(req.Context())
			// action
			handler.GetUsers(rec, req.WithContext(ctx))

			// check
			if serviceCalled != tt.wantServiceCalled {
				t.Fatalf(
					"service called = %v, want %v",
					serviceCalled,
					tt.wantServiceCalled,
				)
			}

			if serviceCalled {
				if diff := cmp.Diff(tt.wantServiceLimit, serviceGotLimit); diff != "" {
					t.Fatalf("limit mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(tt.wantServiceOffset, serviceGotOffset); diff != "" {
					t.Fatalf("offset mismatch (-want +got):\n%s", diff)
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
				var gotResponse []UserDTOResponse
				if err := json.NewDecoder(rec.Body).Decode(&gotResponse); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if diff := cmp.Diff(tt.wantUsers, gotResponse); diff != "" {
					t.Fatalf("GetUsersResponse mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
