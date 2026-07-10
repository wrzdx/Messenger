package core_errors

import (
	"errors"
	"messenger/internal/core/domain"
)

func domainError(e error) (Error, bool) {
	err := Error{
		err:     e,
		Message: e.Error(),
	}

	switch {
	case errors.Is(e, domain.ErrInvalidUsername),
		errors.Is(e, domain.ErrInvalidFirstName),
		errors.Is(e, domain.ErrInvalidLastName),
		errors.Is(e, domain.ErrInvalidBio),
		errors.Is(e, domain.ErrInvalidPassword),
		errors.Is(e, domain.ErrNegativeLimit),
		errors.Is(e, domain.ErrNegativeOffset),
		errors.Is(e, domain.ErrInvalidChatName),
		errors.Is(e, domain.ErrEmptyMessage),
		errors.Is(e, domain.ErrLongMessage):
		err.Code = VALIDATION_ERROR

	case errors.Is(e, domain.ErrAlreadyExists):
		err.Code = ALREADY_EXISTS

	case errors.Is(e, domain.ErrNotFound):
		err.Code = NOT_FOUND

	case errors.Is(e, domain.ErrInvalidCredentials):
		err.Code = INVALID_CREDENTIALS
	case errors.Is(e, domain.ErrWrongPassword):
		err.Code = WRONG_PASSWORD
	default:
		return Error{}, false
	}
	return err, true
}
