package users_postgres_repository

import (
	"messenger/internal/core/domain"
	"time"
)

type UserModel struct {
	ID           int
	Username     string
	FirstName    string
	LastName     *string
	CreatedAt    time.Time
	Bio          *string
	PasswordHash string
}

func UserDomainFromModel(user UserModel) domain.User {
	return domain.NewUser(
		user.ID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.CreatedAt,
		user.Bio,
	)
}

func userDomainsFromModels(users []UserModel) []domain.User {
	userDomains := make([]domain.User, len(users))
	for i, user := range users {
		userDomains[i] = domain.NewUser(
			user.ID,
			user.Username,
			user.FirstName,
			user.LastName,
			user.CreatedAt,
			user.Bio,
		)
	}

	return userDomains
}
