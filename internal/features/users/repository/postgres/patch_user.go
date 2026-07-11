package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/core/domain"
	"messenger/internal/core/repository/postgres"

	"github.com/google/uuid"
)

func (r *UsersRepository) PatchUser(
	ctx context.Context,
	id uuid.UUID,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OptTimeout())
	defer cancel()

	query := `
	UPDATE users
	SET 
		username=$1,
		first_name=$2,
		last_name=$3,
		bio=$4
	WHERE id=$5
	RETURNING id,
    username,
    first_name,
    last_name,
    created_at,
    deleted_at,
    bio,
    password_hash;`

	row := r.db.QueryRow(
		ctx,
		query,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Bio,
		id,
	)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Username,
		&userModel.FirstName,
		&userModel.LastName,
		&userModel.CreatedAt,
		&user.DeletedAt,
		&userModel.Bio,
		&userModel.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, postgres.ErrViolatesUnique) {
			return domain.User{}, domain.ErrAlreadyExists
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Username,
		userModel.FirstName,
		userModel.LastName,
		userModel.CreatedAt,
		userModel.DeletedAt,
		userModel.Bio,
		userModel.PasswordHash,
	)

	return userDomain, nil
}
