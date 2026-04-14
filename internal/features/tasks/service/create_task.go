package tasks_service

import (
	"context"
	"fmt"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
)

func (s *TasksService) CreateTaskService(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf(
			"validate task domain: %w",
			err,
		)
	}

	task, err := s.tasksRepository.CreatePGTask(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}
