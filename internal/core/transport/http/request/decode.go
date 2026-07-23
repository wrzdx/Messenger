package http_request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New(validator.WithRequiredStructEnabled())

type validatable interface {
	Validate() map[string]string
}

func Validate(dest any) map[string]string {
	fields := make(map[string]string)
	verr := requestValidator.Struct(dest)
	if verr != nil {
		if verr, ok := errors.AsType[validator.ValidationErrors](verr); ok {
			for _, fieldErr := range verr {
				field, msg := formatField(fieldErr)
				fields[field] = msg
			}
		}
	}
	return fields
}

func toSnakeCase(str string) string {
	var matchFirstCap = true
	var buf bytes.Buffer

	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) {
			if matchFirstCap {
				buf.WriteRune('_')
			}
		}
		buf.WriteRune(unicode.ToLower(r))
		matchFirstCap = unicode.IsLower(r) || unicode.IsDigit(r)
	}
	return buf.String()
}

func formatField(err validator.FieldError) (string, string) {
	field := toSnakeCase(err.Field())
	var msg string
	switch err.Tag() {
	case "required":
		msg = field + " is required"
	case "uuid":
		msg = "invalid uuid"
	default:
		msg = err.Tag()
	}
	return field, msg
}

func DecodeAndValidateRequestBody(r *http.Request, dest any) error {
	err := json.NewDecoder(r.Body).Decode(dest)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("%w: decode json: %w", ErrInvalidRequest, err)
	}

	if fields := Validate(dest); len(fields) != 0 {
		return NewFieldError(fields)
	}
	return nil
}
