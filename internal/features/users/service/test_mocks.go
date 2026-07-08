package users_service

// import (
// 	"context"
// 	"messenger/internal/core/domain"

// 	"github.com/google/uuid"
// )

// type StubUsersRepository struct {
// 	CreateUserFn func(
// 		user domain.User,
// 	) (domain.User, error)

// 	GetUsersFn func(
// 		limit *int,
// 		offset *int,
// 	) ([]domain.User, error)

// 	GetUserFn func(
// 		id uuid.UUID,
// 	) (domain.User, error)

// 	DeleteUserFn func(
// 		id uuid.UUID,
// 	) error

// 	PatchUserFn func(
// 		id uuid.UUID,
// 		user domain.User,
// 	) (domain.User, error)

// 	ChangePasswordFn func(id uuid.UUID, hash string) error
// }

// func (s *StubUsersRepository) CreateUser(
// 	ctx context.Context,
// 	user domain.User,
// ) (domain.User, error) {
// 	return s.CreateUserFn(user)
// }

// func (s *StubUsersRepository) GetUsers(
// 	ctx context.Context,
// 	pagination domain.Pagination,
// ) ([]domain.User, error) {
// 	return s.GetUsersFn(pagination.Limit, pagination.Offset)
// }

// func (s *StubUsersRepository) GetUser(
// 	ctx context.Context,
// 	id uuid.UUID,
// ) (domain.User, error) {
// 	return s.GetUserFn(id)
// }

// func (s *StubUsersRepository) DeleteUser(
// 	ctx context.Context,
// 	id uuid.UUID,
// ) error {
// 	return s.DeleteUserFn(id)
// }

// func (s *StubUsersRepository) PatchUser(
// 	ctx context.Context,
// 	id uuid.UUID,
// 	user domain.User,
// ) (domain.User, error) {
// 	return s.PatchUserFn(id, user)
// }

// func (s *StubUsersRepository) ChangePassword(
// 	ctx context.Context,
// 	id uuid.UUID,
// 	newPasswordHash string,
// ) error {
// 	return s.ChangePasswordFn(id, newPasswordHash)
// }

// type StubHasher struct {
// 	HashFn    func(password string) (string, error)
// 	CompareFn func(hash, password string) error
// }

// func (h *StubHasher) Hash(password string) (string, error) {
// 	return h.HashFn(password)
// }
// func (h *StubHasher) Compare(hash, password string) error {
// 	return h.CompareFn(hash, password)
// }
