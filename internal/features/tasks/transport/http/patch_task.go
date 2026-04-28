package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
	core_logger "github.com/DimaKirejko/todoapp/internal/core/logger"
	core_http_request "github.com/DimaKirejko/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/DimaKirejko/todoapp/internal/core/transport/http/response"
	core_http_types "github.com/DimaKirejko/todoapp/internal/core/transport/http/types"
)

type PatchTaskReqDTO struct {
	Title       core_http_types.Nullable[string] `json:"title"       swaggertype:"string" example:"Name Name Name"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"null"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"   swaggertype:"boolean"`
}

type PatchTaskResponse TaskDTOResponse

func (r *PatchTaskReqDTO) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("Title can't be NULL")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("Title must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("Description must be between 1 and 1000 symbols")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("Completed can't be NULL")
		}
	}

	return nil
}

// PathTask godoc
// @Summary Змінти параметри задачі
// @Description Змінти параметри задачі
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "user id"
// @Param request body PatchTaskReqDTO true "PathTask тіло запиту"
// @Success 200 {object} PatchTaskResponse "Успішно змінена задача"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "User not found"
// @Failure 409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [patch]
func (h *TasksHTTPHandler) PathTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValues(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task id path value")
		return
	}

	var request PatchTaskReqDTO
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate http request",
		)
	}

	taskPatch := taskPatchFromReq(request)

	taskDomain, err := h.tasksService.PatchTaskService(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"faild to patch task",
		)

		return
	}

	response := PatchTaskResponse(taskDTOfromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromReq(request PatchTaskReqDTO) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
