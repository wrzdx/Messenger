package core_errors

import (
	"errors"
)

type ErrorCode string

const (
	VALIDATION_ERROR    = "VALIDATION_ERROR"
	USER_ALREADY_EXISTS = "USER_ALREADY_EXISTS"
	USER_NOT_FOUND      = "USER_NOT_FOUND"
	INVALID_CREDENTIALS = "INVALID_CREDENTIALS"
	WRONG_PASSWORD      = "WRONG_PASSWORD"
	INVALID_TOKEN       = "INVALID_TOKEN"
	INTERNAL_ERROR      = "INTERNAL_ERROR"
)

var (
	ErrValidation = errors.New("failed to validate")
)

func ValidationError(fields map[string]string) Error {
	return Error{
		err:     ErrValidation,
		Code:    VALIDATION_ERROR,
		Message: ErrValidation.Error(),
		Fields:  fields,
	}
}

func MapError(e error) Error {
	if err, ok := e.(Error); ok {
		return err
	}
	if err, ok := authError(e); ok {
		return err
	}
	if err, ok := domainError(e); ok {
		return err
	}
	if errors.Is(e, ErrValidation) {
		return Error{
			err:     e,
			Code:    VALIDATION_ERROR,
			Message: e.Error(),
		}
	}
	return Error{
		err:     e,
		Code:    INTERNAL_ERROR,
		Message: "internal server error",
	}
}

type Error struct {
	err     error
	Code    ErrorCode
	Message string
	Fields  map[string]string
}

func (e Error) Error() string {
	return e.err.Error()
}
