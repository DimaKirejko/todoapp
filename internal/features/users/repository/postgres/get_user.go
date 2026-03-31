package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
	core_errors "github.com/DimaKirejko/todoapp/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) GetPGUser(
	ctx context.Context,
	id int,
) (domain.User, error) {

	query := `
	SELECT id, version, full_name, phone_number
	FROM todoapp.users
	WHERE id=$1;
	`

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	row := r.pool.QueryRow(ctx, query, id)

	var userModel UserModel

	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id='%d': %w", id, core_errors.ErrNotFound)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomani := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)

	return userDomani, nil
}
