package task_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/DimaKirejko/todoapp/internal/core/errors"
)

func (r *TasksRepository) DeletePGTask(
	ctx context.Context,
	taskID int,
) error {
	query := `
	DELETE FROM todoapp.tasks
	WHERE id=$1;
	`

	ctx, cencel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cencel()

	cmdTag, err := r.pool.Exec(ctx, query, taskID)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("task with id='%d: %w", taskID, core_errors.ErrNotFound)
	}

	return nil
}
