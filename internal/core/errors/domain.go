package core_errors

import (
	"errors"
	"messenger/internal/core/domain"
	"net/http"
)

func domainError(e error) (Error, bool) {
	err := Error{
		err:     e,
		Message: e.Error(),
	}
	switch {
	case errors.Is(e, domain.ErrValidation):
		err.Code = http.StatusBadRequest

	case errors.Is(e, domain.ErrAlreadyExists):
		err.Code = http.StatusConflict

	case errors.Is(e, domain.ErrNotFound):
		err.Code = http.StatusNotFound

	case errors.Is(e, domain.ErrInvalidCredentials):
		err.Code = http.StatusUnauthorized
	case errors.Is(e, domain.ErrWrongPassword):
		err.Code = http.StatusForbidden
	default:
		return Error{}, false
	}

	if de, ok := errors.AsType[domain.DetailedError](e); ok {
		err.Details = de.Details
	}
	return err, true
}
