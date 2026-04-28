package tasks_transport_http

import (
	"net/http"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
	core_logger "github.com/DimaKirejko/todoapp/internal/core/logger"
	core_http_request "github.com/DimaKirejko/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/DimaKirejko/todoapp/internal/core/transport/http/response"
)

type CreateTaskReqDTO struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"         example:"do something"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000" example:"do something in detalse"` //omitempty!!
	AuthorUserID int     `json:"author_user_id" validate:"required"              example:"12"`
}

type CreateTaskRespDTO TaskDTOResponse

// CreateTask godoc
// @Summary Створити задачу
// @Description Створити нову задачу в системі
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body CreateTaskReqDTO true "CreateTaskReq"
// @Success 201 {object} CreateTaskRespDTO "Успішно створена задача"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks [post]
func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateTaskReqDTO
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	taskDomain := domain.NewDomainTaskUninitialized(
		request.Title,
		request.Description,
		request.AuthorUserID,
	)

	taskDomani, err := h.tasksService.CreateTaskService(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)

		return
	}

	response := CreateTaskRespDTO(taskDTOfromDomain(taskDomani)) // необхідно CreateTaskRespDTO(taskDTOfromDomain у випадку змін в CreateTaskRespDTO (type conversion)

	responseHandler.JSONResponse(response, http.StatusCreated)
}
