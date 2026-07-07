package core_test_utils

import (
	"errors"
	"messenger/internal/core/domain"
	core_logger "messenger/internal/core/logger"
	"time"

	"github.com/google/uuid"
)

var (
	ID           = uuid.New()
	CreatedAt    = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	PasswordHash = "password_hash"
	HasherError  = errors.New("hash failed")
	RepoError    = errors.New("db error")
	Log          = core_logger.NewTestLogger()
)

var Users = []domain.User{
	{
		ID:           uuid.New(),
		Username:     "user_1",
		FirstName:    "Username",
		LastName:     new("1"),
		CreatedAt:    CreatedAt,
		Bio:          new("I'm user 1"),
		PasswordHash: PasswordHash,
	},
	{
		ID:           uuid.New(),
		Username:     "user_2",
		FirstName:    "Username",
		LastName:     new("2"),
		CreatedAt:    CreatedAt,
		Bio:          new("I'm user 2"),
		PasswordHash: PasswordHash,
	},
	{
		ID:           uuid.New(),
		Username:     "user_3",
		FirstName:    "Username",
		LastName:     new("3"),
		CreatedAt:    CreatedAt,
		Bio:          new("I'm user 3"),
		PasswordHash: PasswordHash,
	},
}
