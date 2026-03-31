package users_srvice

import (
	"context"
	"fmt"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
)

func (s *UsersService) PatchUserService(
	ctx context.Context,
	id int,
	patch domain.UserPatch,
) (domain.User, error) {
	user, err := s.usersRepository.GetPGUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user: %w", err)
	}

	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf("apply user patch: %w", err)
	}

	patchedUser, err := s.usersRepository.PatchPGUser(ctx, id, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("patch user: %w", err)
	}

	return patchedUser, err

}
