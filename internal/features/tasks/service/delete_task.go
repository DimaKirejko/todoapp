package tasks_service

import (
	"context"
	"fmt"
)

func (s *TasksService) DeleteTaskService(
	ctx context.Context,
	taskID int,
) error {
	if err := s.tasksRepository.DeletePGTask(ctx, taskID); err != nil {
		return fmt.Errorf("delete task from repository: %w", err)
	}

	return nil
}
