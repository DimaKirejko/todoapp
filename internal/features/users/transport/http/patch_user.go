package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/DimaKirejko/todoapp/internal/core/domain"
	core_logger "github.com/DimaKirejko/todoapp/internal/core/logger"
	core_http_request "github.com/DimaKirejko/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/DimaKirejko/todoapp/internal/core/transport/http/response"
	core_http_types "github.com/DimaKirejko/todoapp/internal/core/transport/http/types"
	core_http_utils "github.com/DimaKirejko/todoapp/internal/core/transport/http/utils"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("FullName can't be NULL")
		}

		fullNameLen := len([]rune(*r.FullName.Value)) // речек
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("fullName must be between 3 and 100 symbols")
		}
	} //10:17

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("Phone number must be between 10 and 15")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("PhoneNumber must startswith '+' symbol")
			}
		}
	}

	return nil
}

type PatchUserResponse UsersDTOREsponse

func (h *UsersHttpHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_utils.GetIntPathValues(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	domainedUserPatch := UserPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUserService(ctx, userID, domainedUserPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed  to patch_user",
		)

		return
	}

	log.Debug(
		fmt.Sprintf(
			"PatchUserRequest fields:\nFullName: '%s'\nPhoneNumber: '%s'",
			request.FullName,
			request.PhoneNumber,
		),
	)

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func UserPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.UserPatch{
		FullName:    request.FullName.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}
}
