package users_postgres_repository

import (
	"messenger/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	ID           uuid.UUID
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
		user.PasswordHash,
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
			user.PasswordHash,
		)
	}

	return userDomains
}
