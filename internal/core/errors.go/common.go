package core_errors

import "errors"

var (
	ErrorNotFound      = errors.New("not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrConflict        = errors.New("conflict")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrInternalServer          = errors.New("internal server error")
)
