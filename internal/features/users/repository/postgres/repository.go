package users_postgres_repository

import (
	"errors"
	"messenger/internal/core/domain"
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

var userConstraints = map[string]struct {
	field  string
	getVal func(u domain.User) string
}{
	"users_username_key": {
		field:  "username",
		getVal: func(u domain.User) string { return u.Username },
	},
}

func getConstraintValues(user domain.User, err error) (string, string) {
	failedField := "unknown_field"
	failedValue := ""
	var target interface{ Constraint() string }
	if errors.As(err, &target) {
		constraintName := target.Constraint()
		if cfg, ok := userConstraints[constraintName]; ok {
			failedField = cfg.field
			failedValue = cfg.getVal(user)
		}
	}
	return failedField, failedValue
}
