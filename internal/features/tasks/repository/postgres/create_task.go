package task_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
	core_errors "github.com/DimaKirejko/todoapp/internal/core/errors"
	core_postgres_pool "github.com/DimaKirejko/todoapp/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) CreatePGTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	query := `
	INSERT INTO todoapp.tasks (title, description, completed, created_at, completed_at, author_user_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id;
	`

	ctx, cencel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cencel()

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserID,
	)

	var TaskModel TaskModel

	err := row.Scan(
		&TaskModel.ID,
		&TaskModel.Version,
		&TaskModel.Title,
		&TaskModel.Description,
		&TaskModel.Completed,
		&TaskModel.CreatedAt,
		&TaskModel.CompletedAt,
		&TaskModel.AuthorUserID,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Task{}, fmt.Errorf(
				"%v user with id='%d': %w",
				err,
				task.AuthorUserID,
				core_errors.ErrNotFound,
			)
		}

		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := domain.NewDomainTask(
		TaskModel.ID,
		TaskModel.Version,
		TaskModel.Title,
		TaskModel.Description,
		TaskModel.Completed,
		TaskModel.CreatedAt,
		TaskModel.CompletedAt,
		TaskModel.AuthorUserID,
	)

	return taskDomain, nil
}
