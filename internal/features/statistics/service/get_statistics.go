package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
	core_errors "github.com/DimaKirejko/todoapp/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userId *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf(
				"'to' must be after from: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	tasks, err := s.statisticsRepository.GetStatPGTasks(ctx, userId, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf(
			"get tasks from repository: %w", err,
		)
	}

	statisticas := calcStatistics(tasks)

	return statisticas, nil

}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	tasksCreated := len(tasks)

	if tasksCreated == 0 {
		return domain.NewStatistics(0, 0, nil, nil)
	}

	tasksCompleted := 0
	var totalCompletionDuration time.Duration
	for _, task := range tasks {
		if task.Completed == true {
			tasksCompleted++
		}

		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totalCompletionDuration += *completionDuration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAverageCompletionTime *time.Duration
	if tasksCompleted > 0 && totalCompletionDuration != 0 {
		avg := totalCompletionDuration / time.Duration(tasksCompleted)

		tasksAverageCompletionTime = &avg
	}

	return domain.NewStatistics(
		tasksCreated,
		tasksCompleted,
		&tasksCompletedRate,
		tasksAverageCompletionTime,
	)
}
