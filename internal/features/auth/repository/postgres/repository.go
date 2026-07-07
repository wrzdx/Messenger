package auth_postgres_repository

import core_postgres "messenger/internal/core/repository/postgres"

type AuthRepository struct {
	db core_postgres.DB
}

func NewUsersRepository(db core_postgres.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}
