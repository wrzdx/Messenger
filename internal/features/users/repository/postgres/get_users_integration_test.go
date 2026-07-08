// //go:build integration

package users_postgres_repository

// import (
// 	"errors"
// 	"messenger/internal/core/domain"
// 	pgx_pool "messenger/internal/core/repository/postgres/pgx"
// 	test_utils "messenger/internal/core/utils/test"
// 	"testing"

// 	"github.com/google/go-cmp/cmp"
// )

// func TestGetUsers(t *testing.T) {
// 	// common setup

// 	var tests = []struct {
// 		name      string
// 		users     []domain.User
// 		limit     *int
// 		offset    *int
// 		wantError error
// 	}{
// 		{
// 			name:  "all users",
// 			users: test_utils.Users,
// 		},
// 		{
// 			name:  "limit users",
// 			users: test_utils.Users[:1],
// 			limit: new(1),
// 		},
// 		{
// 			name:      "negative limit",
// 			users:     nil,
// 			limit:     new(-1),
// 			wantError: domain.ErrNegativeLimit,
// 		},
// 		{
// 			name:   "offset users",
// 			users:  test_utils.Users[1:],
// 			offset: new(1),
// 		},
// 		{
// 			name:      "negative offset",
// 			users:     nil,
// 			offset:    new(-1),
// 			wantError: domain.ErrNegativeOffset,
// 		},
// 		{
// 			name:   "limit offset users",
// 			users:  test_utils.Users[1:2],
// 			limit:  new(1),
// 			offset: new(1),
// 		},
// 		{
// 			name:  "empty users",
// 			users: []domain.User{},
// 			limit: new(0),
// 		},
// 	}
// 	pool, err := pgx_pool.NewPool(t.Context(), pgx_pool.NewConfigMust())
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// 	defer pool.Close()

// 	// subtests
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tx, err := pool.Begin(t.Context())
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			defer tx.Rollback(t.Context())
// 			test_utils.LoadData(t, tx)
// 			repository := NewUsersRepository(tx)
// 			pagination := domain.NewPagination(tt.limit, tt.offset)
// 			// action
// 			gotUsers, gotErr := repository.GetUsers(t.Context(), pagination)

// 			// assertion
// 			if !errors.Is(gotErr, tt.wantError) {
// 				t.Fatalf("want %v, got %v", tt.wantError, gotErr)
// 			}

// 			if diff := cmp.Diff(tt.users, gotUsers); diff != "" {
// 				t.Fatalf("GetUsers mismatch got users (-want +got):\n%s", diff)
// 			}
// 		})
// 	}
// }
