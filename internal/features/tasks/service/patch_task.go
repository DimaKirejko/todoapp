package tasks_service

import (
	"context"
	"fmt"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
)

func (s *TasksService) PatchTaskService(
	ctx context.Context,
	id int,
	Pathc domain.TaskPatch,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetPGTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}

	if err := task.ApplyPatch(Pathc); err != nil {
		return domain.Task{}, fmt.Errorf("applay task patch: %w", err)
	}

	patchedTask, err := s.tasksRepository.PatchPGTask(ctx, id, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("patch task: %w", err)
	}

	return patchedTask, nil
}
