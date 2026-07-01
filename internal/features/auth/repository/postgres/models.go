package auth_postgres_repository

import "messenger/internal/core/domain"

type UserAuthModel struct {
	UserID       int
	PasswordHash string
}

func UserAuthDomainFromModel(m UserAuthModel) domain.UserAuth {
	return domain.NewUserAuth(m.UserID, m.PasswordHash)
}
