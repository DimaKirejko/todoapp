package users_srvice

import (
	"context"
	"fmt"
)

func (s *UsersService) DeleteUserService(ctx context.Context, id int) error {
	if err := s.usersRepository.DeletePGUser(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
