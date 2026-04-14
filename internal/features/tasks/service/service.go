package tasks_service

import (
	"context"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
)

type TasksService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	CreatePGTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)

	GetPGTasks(
		ctx context.Context,
		userID *int,
		limit *int,
		offset *int,
	) ([]domain.Task, error)

	GetPGTask(
		ctx context.Context,
		taskID int,
	) (domain.Task, error)

	DeletePGTask(
		ctx context.Context,
		taskID int,
	) error

	PatchPGTask(
		ctx context.Context,
		taskID int,
		task domain.Task,
	) (domain.Task, error)
}

func NewTasksService(
	tasksRepository TasksRepository,
) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
	}
}
