package core_errors

import (
	"errors"
	"messenger/internal/core/auth"
	"net/http"
)

func authError(e error) (Error, bool) {
	err := Error{
		err:     e,
		Message: e.Error(),
	}

	switch {
	case errors.Is(e, auth.ErrInvalidToken):
		err.Code = http.StatusUnauthorized
	case errors.Is(e, auth.ErrPasswordMismatch):
		err.Code = http.StatusForbidden
	default:
		return Error{}, false
	}

	return err, true
}
