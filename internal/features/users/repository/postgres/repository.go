package users_postgres_repository

import (
	postgres "messenger/internal/core/repository/postgres"
)

type UsersRepository struct {
	db postgres.DB
}

func NewUsersRepository(db postgres.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}
