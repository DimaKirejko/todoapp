package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
	core_http_server "github.com/DimaKirejko/todoapp/internal/core/transport/http/server"
)

type TasksHTTPHandler struct {
	tasksService TasksService
}

type TasksService interface {
	CreateTaskService(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)

	GetTasksService(
		ctx context.Context,
		userID *int,
		limit *int,
		offset *int,
	) ([]domain.Task, error)

	GetTaskService(
		ctx context.Context,
		taskID int,
	) (domain.Task, error)

	DeleteTaskService(
		ctx context.Context,
		taskID int,
	) error

	PatchTaskService(
		ctx context.Context,
		id int,
		Pathc domain.TaskPatch,
	) (domain.Task, error)
}

func NewTasksHTTPHandler(
	tasksService TasksService,
) *TasksHTTPHandler {
	return &TasksHTTPHandler{
		tasksService: tasksService,
	}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: h.CreateTask,
		},

		{
			Method:  http.MethodGet,
			Path:    "/tasks",
			Handler: h.GetTasks,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks/{id}",
			Handler: h.GetTask,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/tasks/{id}",
			Handler: h.DeleteTask,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/tasks/{id}",
			Handler: h.PathTask,
		},
	}
}
