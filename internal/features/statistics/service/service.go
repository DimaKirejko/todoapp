package statistics_service

import (
	"context"
	"time"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetStatPGTasks(
		ctx context.Context,
		userId *int,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error)
}

func NewStatisticsService(StatisticsRepository StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: StatisticsRepository,
	}
}
