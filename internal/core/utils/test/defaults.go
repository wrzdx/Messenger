package test_utils

import (
	"messenger/internal/core/domain"
	logger "messenger/internal/core/logger"
	"time"

	"github.com/google/uuid"
)

var (
	MockUser = domain.NewUser(
		uuid.New(),
		"ecorp",
		"Tyrell",
		new("Wellick"),
		time.Now(),
		new("Dead"),
		"fsociety_hash",
	)
	Log = logger.NewTestLogger()
)

var MockUsers = []domain.User{
	{
		ID:           uuid.New(),
		Username:     "user_1",
		FirstName:    "Username",
		LastName:     new("1"),
		CreatedAt:    time.Now(),
		Bio:          new("I'm user 1"),
		PasswordHash: "password hash",
	},
	{
		ID:           uuid.New(),
		Username:     "user_2",
		FirstName:    "Username",
		LastName:     new("2"),
		CreatedAt:    time.Now(),
		Bio:          new("I'm user 2"),
		PasswordHash: "password hash",
	},
	{
		ID:           uuid.New(),
		Username:     "user_3",
		FirstName:    "Username",
		LastName:     new("3"),
		CreatedAt:    time.Now(),
		Bio:          new("I'm user 3"),
		PasswordHash: "password hash",
	},
}
