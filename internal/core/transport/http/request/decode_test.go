package http_request

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDecodeAndValidateRequest(t *testing.T) {
	t.Run("decodes valid request", func(t *testing.T) {
		type request struct {
			Username string `json:"username" validate:"required"`
		}
		r := httptest.NewRequest(
			"POST",
			"/",
			bytes.NewBufferString(`{"username":"Username_1"}`),
		)
		var decoded request

		err := DecodeAndValidateRequest(r, &decoded)

		require.NoError(t, err)
		require.Equal(t, "Username_1", decoded.Username)
	})

	t.Run("classifies malformed JSON as invalid request", func(t *testing.T) {
		r := httptest.NewRequest(
			"POST",
			"/",
			bytes.NewBufferString(`{"username":invalid}`),
		)
		var decoded struct {
			Username string `json:"username"`
		}

		err := DecodeAndValidateRequest(r, &decoded)

		require.ErrorIs(t, err, ErrInvalidRequest)
		require.Contains(t, err.Error(), "decode json")
		var syntaxError *json.SyntaxError
		require.ErrorAs(t, err, &syntaxError)
	})

	t.Run("returns validation fields for invalid DTO", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/", bytes.NewReader(nil))
		var decoded struct {
			Username  string `json:"username" validate:"required"`
			FirstName string `json:"first_name" validate:"required"`
		}

		err := DecodeAndValidateRequest(r, &decoded)

		require.ErrorIs(t, err, ErrInvalidRequest)
		var withFields interface {
			Fields() map[string]string
		}
		require.True(t, errors.As(err, &withFields))
		require.Equal(t, map[string]string{
			"username":   "username is required",
			"first_name": "first_name is required",
		}, withFields.Fields())
	})

	t.Run("does not expose mutable validation fields", func(t *testing.T) {
		err := newFieldError(map[string]string{"username": "required"})
		fields := err.Fields()
		fields["username"] = "changed"

		require.Equal(t, "required", err.Fields()["username"])
	})

	t.Run("rejects required zero-value UUID", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{}`))
		var decoded struct {
			ID uuid.UUID `json:"id" validate:"required"`
		}

		err := DecodeAndValidateRequest(r, &decoded)

		require.ErrorIs(t, err, ErrInvalidRequest)
		var withFields interface {
			Fields() map[string]string
		}
		require.ErrorAs(t, err, &withFields)
		require.Equal(t, "id is required", withFields.Fields()["id"])
	})
}
