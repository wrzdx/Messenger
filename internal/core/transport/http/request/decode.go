package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "messenger/internal/core/errors"

	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

func validationError(err error) error {
	if err == nil {
		return nil
	}
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}
	e := errs[0]

	switch e.Tag() {
	case "required":
		return fmt.Errorf("%s is required: %w",
			e.Field(),
			core_errors.ErrInvalidArgument,
		)

	case "min":
		return fmt.Errorf("%s must contain at least %s characters: %w",
			e.Field(),
			e.Param(),
			core_errors.ErrInvalidArgument,
		)

	case "max":
		return fmt.Errorf("%s must contain at most %s characters: %w",
			e.Field(),
			e.Param(),
			core_errors.ErrInvalidArgument,
		)

	default:
		return core_errors.ErrInvalidArgument
	}
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"decode json: %v: %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	v, ok := dest.(validatable)

	var err error
	if ok {
		err = v.Validate()
	} else {
		err = validationError(requestValidator.Struct(dest))
	}

	if err != nil {
		return fmt.Errorf("request validation: %w", err)
	}

	return nil
}
