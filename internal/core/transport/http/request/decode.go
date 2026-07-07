package core_http_request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

func formatValidationErrors(errs validator.ValidationErrors) error {
	var messages []string

	for _, err := range errs {
		var msg string
		switch err.Tag() {
		case "required":
			msg = "is required"
		case "min":
			msg = fmt.Sprintf("must be at least %s characters", err.Param())
		case "max":
			msg = fmt.Sprintf("must be at most %s characters", err.Param())
		default:
			msg = err.Tag()
		}

		messages = append(messages, fmt.Sprintf("%s %s", err.Field(), msg))
	}

	return errors.New(strings.Join(messages, "; "))
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"decode json: %w",
			err,
		)
	}

	v, ok := dest.(validatable)

	var err error
	if ok {
		err = v.Validate()
	} else {
		err = requestValidator.Struct(dest)
		if err != nil {
			if verr, ok := errors.AsType[validator.ValidationErrors](err); ok {
				err = formatValidationErrors(verr)
			}
		}
	}
	return err
}
