package tasks_service

import (
	"context"
	"fmt"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
)

func (s *TasksService) GetTaskService(
	ctx context.Context,
	taskID int,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetPGTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("Get Task from repository: %w", err)
	}

	return task, nil
}
