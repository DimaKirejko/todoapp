package users_srvice

import (
	"context"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
)

type UsersService struct {
	usersRepository UsersRepository
}

type UsersRepository interface {
	CreatePGUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetPGUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetPGUser(
		ctx context.Context,
		id int,
	) (domain.User, error)

	DeletePGUser(
		ctx context.Context,
		id int,
	) error

	PatchPGUser(
		ctx context.Context,
		id int,
		user domain.User,
	) (domain.User, error)
}

func NewUsersService(usersRepository UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
