package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/DimaKirejko/todoapp/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"decode json: %v :%w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	var (
		err error
	)

	v, ok := dest.(validatable) // тайп асершн (треба перевірити)
	if ok {
		err = v.Validate() // ти не можеш написати: dest.Validate(). тип змінної dest — any, а в any немає методу Validate(). Компiлятор НЕ знає: що там User що там є Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf(
			"request validation: %v: %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
