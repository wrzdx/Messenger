package auth_postgres_repository

import "messenger/internal/core/repository/postgres"

type AuthRepository struct {
	db postgres.DB
}


func NewAuthRepository(db postgres.DB) AuthRepository {
	return AuthRepository{
		db: db,
	}

}



