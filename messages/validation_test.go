package messages

import (
	"github.com/mdev5000/qlik_message/apperrors"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestValidation_errorWhenEmpty(t *testing.T) {
	requireHasResponseErrors(t, validateMessage("", ModifyMessage{Message: ""}), apperrors.FieldErrorResponse{
		Field: "message",
		Error: "Message field cannot be blank.",
	})
}

func TestValidation_errorWhenLargeThanMaxCharLimit(t *testing.T) {
	msg := strings.Repeat("m", MaxMessageCharLength+1)
	requireHasResponseErrors(t, validateMessage("", ModifyMessage{Message: msg}), apperrors.FieldErrorResponse{
		Field: "message",
		Error: "Message cannot be longer than 512 characters.",
	})
}

func TestValidation_validWhenAtCharacterLimit(t *testing.T) {
	msg := strings.Repeat("m", MaxMessageCharLength)
	require.Nil(t, validateMessage("", ModifyMessage{Message: msg}))
}
