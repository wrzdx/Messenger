package core_errors

import (
	"net/http"
)

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
	return Error{
		err:     e,
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
	}
}

type Error struct {
	err     error
	Code    int
	Message string
	Details map[string]string
}

func (e Error) Error() string {
	return e.err.Error()
}
