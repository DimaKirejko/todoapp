package tasks_service

import (
	"context"
	"fmt"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
	core_errors "github.com/DimaKirejko/todoapp/internal/core/errors"
)

func (s *TasksService) GetTasksService(
	ctx context.Context,
	userID *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be non-negativ: %w",
			core_errors.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("limit must be non-negativ: %w",
			core_errors.ErrInvalidArgument)
	}

	tasks, err := s.tasksRepository.GetPGTasks(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get tasks from repository: %w", err)
	}

	return tasks, nil
}
