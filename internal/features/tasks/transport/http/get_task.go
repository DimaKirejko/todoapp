package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/DimaKirejko/todoapp/internal/core/logger"
	core_http_request "github.com/DimaKirejko/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/DimaKirejko/todoapp/internal/core/transport/http/response"
)

type GetTaskRespone TaskDTOResponse

// GetTask      godoc
// @Summary     Отримати параметри задачі
// @Description Отримати параметри задачі
// @Tags        tasks
// @Produce     json
// @Param       id  path int true                             "ID задачі"
// @Success     200 {object} GetTaskRespone                  "Успішно віднайдений користувач"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "Task not found
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [GET]
func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValues(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get TaskID path value",
		)

		return
	}

	domainTask, err := h.tasksService.GetTaskService(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get Task",
		)

		return
	}

	response := GetTaskRespone(taskDTOfromDomain(domainTask))

	responseHandler.JSONResponse(response, http.StatusOK)
}
