package tasks_transport_http

import (
	"time"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
)

type TaskDTOResponse struct {
	ID           int        `json:"id"             example:"5"`
	Version      int        `json:"version"        example:"2"`
	Title        string     `json:"title"          example:"do something"`
	Description  *string    `json:"description"    example:"do something in detalse"`
	Completed    bool       `json:"completed"      example:"false"`
	CreatedAt    time.Time  `json:"created_at"     example:"2026-01-01T10:30:00Z"`
	CompletedAt  *time.Time `json:"completed_at"   example:"null"`
	AuthorUserID int        `json:"author_user_id" example:"12"`
}

func taskDTOfromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func taskDTOsFromDomains(tasks []domain.Task) []TaskDTOResponse {
	dtos := make([]TaskDTOResponse, len(tasks))

	for i, task := range tasks {
		dtos[i] = taskDTOfromDomain(task)
	}

	return dtos
}
