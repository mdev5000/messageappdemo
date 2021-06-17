package messages

import (
	"fmt"

	"github.com/mdev5000/messageappdemo/apperrors"
)

func validateMessage(op string, message ModifyMessage) error {
	if message.Message == "" {
		return validationFieldError(op, "message", "Message field cannot be blank.")
	}

	if len([]rune(message.Message)) > MaxMessageCharLength {
		return validationFieldError(op, "message",
			fmt.Sprintf("Message cannot be longer than %d characters.", MaxMessageCharLength))
	}

	return nil
}

func validationFieldError(op string, field, error string) error {
	re := apperrors.Error{Op: op, EType: apperrors.ETInvalid}
	re.AddResponse(apperrors.FieldErrorResponse{
		Field: field,
		Error: error,
	})
	return &re
}
