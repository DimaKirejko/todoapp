package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/DimaKirejko/todoapp/internal/core/logger"
	core_http_request "github.com/DimaKirejko/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/DimaKirejko/todoapp/internal/core/transport/http/response"
)

type GetUsersResponse []UsersDTOREsponse

// GetUsers     godoc
// @Summary     Отримати список користувачів
// @Description Отримати список користувачів
// @Tags        users
// @Produce     json
// @Param       limit query int false                         "розмір сторінки з користувачами"
// @Param       offset query int false                        "зсув сторінки з користувачами"
// @Success     200 {object} GetUsersResponse                 "Успішне отрмання користувачів"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users [get]
func (h *UsersHttpHandler) GetUsersTransport(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := getLimitOfSetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'limit'/'offset' query param",
		)

		return
	}

	usersDomains, err := h.usersService.GetUsersService(ctx, limit, offset) ////
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get users",
		)

		return
	}

	response := GetUsersResponse(usersDTOFromDomains(usersDomains))

	responseHandler.JSONResponse(response, http.StatusOK)
}

// retutn Limit -> *int, Ofset -> *int, Error -> error
func getLimitOfSetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	limit, err := core_http_request.GetIntQuertParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' quert param: %w", err)
	}

	offset, err := core_http_request.GetIntQuertParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' quert param: %w", err)
	}

	return limit, offset, nil
}
