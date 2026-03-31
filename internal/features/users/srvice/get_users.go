package users_srvice

import (
	"context"
	"fmt"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
	core_errors "github.com/DimaKirejko/todoapp/internal/core/errors"
)

func (s *UsersService) GetUsersService(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be non-negativ: %w",
			core_errors.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 { // 8:34
		return nil, fmt.Errorf("limit must be non-negativ: %w",
			core_errors.ErrInvalidArgument)
	}

	users, err := s.usersRepository.GetPGUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get users from repository: %w", err)
	}

	return users, nil

}
