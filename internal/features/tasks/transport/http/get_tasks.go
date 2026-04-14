package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/DimaKirejko/todoapp/internal/core/logger"
	core_http_request "github.com/DimaKirejko/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/DimaKirejko/todoapp/internal/core/transport/http/response"
)

type GetTasksResponseDTO []TaskDTOResponse

func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	loger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(loger, rw)

	userID, limit, offset, err := getUserIdLimitOfSetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user id query params",
		)

		return
	}

	tasksDomains, err := h.tasksService.GetTasksService(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)

		return
	}

	response := taskDTOsFromDomains(tasksDomains)

	responseHandler.JSONResponse(response, http.StatusOK)

}

// retutn UserId -> *int Limit -> *int, Ofset -> *int, Error -> error
func getUserIdLimitOfSetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		userIDQueryParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userID, err := core_http_request.GetIntQuertParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	limit, err := core_http_request.GetIntQuertParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' quert param: %w", err)
	}

	offset, err := core_http_request.GetIntQuertParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' quert param: %w", err)
	}

	return userID, limit, offset, nil
}
