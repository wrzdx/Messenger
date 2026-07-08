package core_errors

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

func MapError(e error) Error {
	if err, ok := domainError(e); ok {
		return err
	}
	return Error{
		Err:     e,
		Code:    INTERNAL_ERROR,
		Message: "internal server error",
	}
}

type Error struct {
	Err     error
	Code    ErrorCode
	Message string
	Fields  map[string]string
}
