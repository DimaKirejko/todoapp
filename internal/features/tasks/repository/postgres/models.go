package task_postgres_repository

import (
	"time"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func taskDomaniFromModels(TaskModel []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(TaskModel))

	for i, model := range TaskModel {
		domains[i] = taskDomaniFromModel(model)
	}

	return domains
}

func taskDomaniFromModel(taskModel TaskModel) domain.Task {
	return domain.NewDomainTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserID,
	)
}
