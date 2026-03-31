package users_transport_http

import (
	"net/http"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
	core_logger "github.com/DimaKirejko/todoapp/internal/core/logger"
	core_http_request "github.com/DimaKirejko/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/DimaKirejko/todoapp/internal/core/transport/http/response"
)

type CrateUserReqDTO struct {
	FullNmae    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"`
}

type CreateUserResponse UsersDTOREsponse

func (h *UsersHttpHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CrateUserReqDTO
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	userDomain := domainFromDTO(request)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to crate user")

		return
	}

	response := CreateUserResponse(UsersDTOREsponse(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)

}

func domainFromDTO(dto CrateUserReqDTO) domain.User {
	return domain.NewUserUninitialized(dto.FullNmae, dto.PhoneNumber)
}
