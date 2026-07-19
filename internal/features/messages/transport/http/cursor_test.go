package messages_transport_http

import (
	"encoding/base64"
	"errors"
	http_request "messenger/internal/core/transport/http/request"
	messages_service "messenger/internal/features/messages/service"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMessageCursor(t *testing.T) {
	t.Run("encodes nil cursor as nil", func(t *testing.T) {
		encoded, err := encodeMessageCursor(nil)

		require.NoError(t, err)
		require.Nil(t, encoded)
	})

	t.Run("round trips cursor", func(t *testing.T) {
		expected := &messages_service.MessageCursor{
			MessageID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			CreatedAt: time.Date(2026, time.July, 19, 12, 30, 0, 0, time.UTC),
		}

		encoded, err := encodeMessageCursor(expected)
		require.NoError(t, err)
		require.NotNil(t, encoded)

		actual, err := decodeMessageCursor(*encoded)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("decodes empty cursor as nil", func(t *testing.T) {
		cursor, err := decodeMessageCursor("")

		require.NoError(t, err)
		require.Nil(t, cursor)
	})

	testCases := []struct {
		name   string
		cursor string
	}{
		{
			name:   "rejects malformed base64",
			cursor: "%%%",
		},
		{
			name:   "rejects malformed json",
			cursor: base64.URLEncoding.EncodeToString([]byte(`{"message_id":`)),
		},
		{
			name: "rejects nil message id",
			cursor: base64.URLEncoding.EncodeToString([]byte(
				`{"created_at":"2026-07-19T12:30:00Z"}`,
			)),
		},
		{
			name: "rejects zero created at",
			cursor: base64.URLEncoding.EncodeToString([]byte(
				`{"message_id":"00000000-0000-0000-0000-000000000001"}`,
			)),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cursor, err := decodeMessageCursor(testCase.cursor)

			require.ErrorIs(t, err, http_request.ErrInvalidRequest)
			require.Nil(t, cursor)
		})
	}

	t.Run("does not classify encoding failure as invalid request", func(t *testing.T) {
		_, err := encodeMessageCursor(&messages_service.MessageCursor{
			MessageID: uuid.New(),
			CreatedAt: time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC),
		})

		require.Error(t, err)
		require.False(t, errors.Is(err, http_request.ErrInvalidRequest))
	})
}
