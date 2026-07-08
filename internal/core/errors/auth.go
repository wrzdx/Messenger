package core_errors

import (
	"errors"
	"messenger/internal/core/auth"
)

func authError(e error) (Error, bool) {
	err := Error{
		err:     e,
		Message: e.Error(),
	}

	switch {
	case errors.Is(e, auth.ErrInvalidToken):
		err.Code = INVALID_TOKEN
	case errors.Is(e, auth.ErrPasswordMismatch):
		err.Code = WRONG_PASSWORD
	default:
		return Error{}, false
	}

	return err, true
}
