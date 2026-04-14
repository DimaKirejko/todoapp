package task_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
	core_errors "github.com/DimaKirejko/todoapp/internal/core/errors"
	core_postgres_pool "github.com/DimaKirejko/todoapp/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) GetPGTask(
	ctx context.Context,
	taskID int,
) (domain.Task, error) {
	query := `
	SELECT  id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM todoapp.tasks
	WHERE id=$1;`

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	row := r.pool.QueryRow(ctx, query, taskID)

	var taskModel TaskModel

	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf("task with id='%d': %w", taskID, core_errors.ErrNotFound)
		}

		return domain.Task{}, fmt.Errorf("Scan error: %w", err)
	}

	taskDomain := taskDomaniFromModel(taskModel)

	return taskDomain, nil
}
