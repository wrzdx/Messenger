//go:build integration

package users_postgres_repository

// import (
// 	"errors"
// 	"messenger/internal/core/domain"
// 	pgx_pool "messenger/internal/core/repository/postgres/pgx"
// 	test_utils "messenger/internal/core/utils/test"
// 	"strings"
// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/require"
// )

// func TestDeleteUser(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		userID    uuid.UUID
// 		wantError error
// 	}{
// 		{
// 			name:   "existing user",
// 			userID: test_utils.MockUsers[0].ID,
// 		},
// 		{
// 			name:      "non-existing user",
// 			userID:    test_utils.MockUser.ID,
// 			wantError: domain.ErrNotFound,
// 		},
// 	}

// 	pool, err := pgx_pool.NewPool(t.Context(), pgx_pool.NewConfigMust())
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// 	defer pool.Close()

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tx, err := pool.Begin(t.Context())
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			defer tx.Rollback(t.Context())
// 			repository := NewUsersRepository(tx)
// 			test_utils.LoadData(t, tx)

// 			gotErr := repository.DeleteUser(t.Context(), tt.userID)

// 			if !errors.Is(gotErr, tt.wantError) {
// 				t.Fatalf("want %v, got %v", tt.wantError, gotErr)
// 			}

// 			if tt.wantError == nil {
// 				user, err := repository.GetUser(t.Context(), tt.userID)

// 				require.NoError(t, err)
// 				require.NotNil(t, user.DeletedAt)
// 				require.Equal(t, "Deleted Account", user.FirstName)
// 				require.Nil(t, user.LastName)
// 				require.Nil(t, user.Bio)
// 				require.Equal(t, "", user.PasswordHash)
// 				require.True(t, strings.HasPrefix(user.Username, "deleted_"))
// 			}
// 		})
// 	}
// }
