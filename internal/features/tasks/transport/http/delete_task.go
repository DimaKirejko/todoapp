package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/DimaKirejko/todoapp/internal/core/logger"
	core_http_request "github.com/DimaKirejko/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/DimaKirejko/todoapp/internal/core/transport/http/response"
)

// DeleteTask   godoc
// @Summary     Видалення задачу
// @Description Видалення задачу
// @Tags        tasks
// @Param       id  path int true                             "ID задачі яку слід видалити"
// @Success     204                                           "Успішне видаленя задачі"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "Task not found
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [delete]
func (h *TasksHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
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

	if err := h.tasksService.DeleteTaskService(ctx, taskID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete task",
		)

		return
	}

	responseHandler.NoContentResponse()

}
