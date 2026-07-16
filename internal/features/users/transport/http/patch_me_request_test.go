package users_transport_http

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatchProfileRequestDecoding(t *testing.T) {
	t.Run("leaves omitted fields unset", func(t *testing.T) {
		request := decodePatchProfileRequest(t, `{}`)
		command := request.ToCommand()

		require.Nil(t, command.Username)
		require.Nil(t, command.FirstName)
		require.False(t, command.LastName.Set)
		require.False(t, command.Bio.Set)
	})

	t.Run("treats null non-nullable fields as omitted", func(t *testing.T) {
		request := decodePatchProfileRequest(t, `{
			"username": null,
			"first_name": null
		}`)
		command := request.ToCommand()

		require.Nil(t, command.Username)
		require.Nil(t, command.FirstName)
	})

	t.Run("keeps provided non-nullable field values", func(t *testing.T) {
		request := decodePatchProfileRequest(t, `{
			"username": "Updated_user",
			"first_name": "Updated name"
		}`)
		command := request.ToCommand()

		require.NotNil(t, command.Username)
		require.Equal(t, "Updated_user", *command.Username)
		require.NotNil(t, command.FirstName)
		require.Equal(t, "Updated name", *command.FirstName)
	})

	t.Run("distinguishes omitted nullable fields from null", func(t *testing.T) {
		request := decodePatchProfileRequest(t, `{
			"last_name": null,
			"bio": null
		}`)
		command := request.ToCommand()

		require.True(t, command.LastName.Set)
		require.Nil(t, command.LastName.Value)
		require.True(t, command.Bio.Set)
		require.Nil(t, command.Bio.Value)
	})

	t.Run("keeps provided nullable field values", func(t *testing.T) {
		request := decodePatchProfileRequest(t, `{
			"last_name": "Anderson",
			"bio": "Updated bio"
		}`)
		command := request.ToCommand()

		require.True(t, command.LastName.Set)
		require.NotNil(t, command.LastName.Value)
		require.Equal(t, "Anderson", *command.LastName.Value)
		require.True(t, command.Bio.Set)
		require.NotNil(t, command.Bio.Value)
		require.Equal(t, "Updated bio", *command.Bio.Value)
	})

	t.Run("rejects value of wrong type", func(t *testing.T) {
		var request PatchProfileRequest

		err := json.Unmarshal([]byte(`{"bio": 42}`), &request)

		require.Error(t, err)
	})
}

func decodePatchProfileRequest(t *testing.T, body string) PatchProfileRequest {
	t.Helper()

	var request PatchProfileRequest
	require.NoError(t, json.Unmarshal([]byte(body), &request))
	return request
}
