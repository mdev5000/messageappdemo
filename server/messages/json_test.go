package messages

import (
	"encoding/json"
	"github.com/mdev5000/qlik_message/messages"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQueryMessageValueToJson_canEncodeNoValues(t *testing.T) {
	jsonValue := queryMessageToJsonValue(&messages.Message{}, fieldsMap{})
	out, err := json.Marshal(jsonValue)
	require.NoError(t, err)
	require.Equal(t, `{}`, string(out))
}

func TestQueryMessageValueToJson_canEncodeAllValues(t *testing.T) {
	loc, err := time.LoadLocation("UTC")
	require.NoError(t, err)

	jsonValue := queryMessageToJsonValue(&messages.Message{
		Id:        5,
		Version:   2,
		CreatedAt: time.Date(2020, 10, 5, 2, 3, 4, 2, loc),
		UpdatedAt: time.Date(2020, 10, 5, 3, 3, 4, 2, loc),
		Message:   "some message",
	}, fieldsMap{
		messages.FieldId:        struct{}{},
		messages.FieldVersion:   struct{}{},
		messages.FieldCreatedAt: struct{}{},
		messages.FieldUpdatedAt: struct{}{},
		messages.FieldMessage:   struct{}{},
		fieldIsPalindrome:       struct{}{},
	})
	out, err := json.Marshal(jsonValue)
	require.NoError(t, err)
	require.Equal(t, `{`+
		`"id":5,"version":2,`+
		`"created_at":"2020-10-05T02:03:04.000000002Z","updated_at":"2020-10-05T03:03:04.000000002Z"`+
		`,"message":"some message","isPalindrome":false`+
		`}`, string(out))
}
