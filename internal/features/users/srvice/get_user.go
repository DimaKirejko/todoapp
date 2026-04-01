package users_srvice

import (
	"context"
	"fmt"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
)

func (s *UsersService) GetUserService(
	ctx context.Context,
	id int,
) (domain.User, error) {
	user, err := s.usersRepository.GetPGUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	return user, nil
}
