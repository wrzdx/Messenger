package users_postgres_repository

import core_postgres "messenger/internal/core/repository/postgres"

type UsersRepository struct {
	db core_postgres.DB
}

func NewUsersRepository(db core_postgres.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}
