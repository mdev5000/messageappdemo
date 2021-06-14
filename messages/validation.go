package messages

import "github.com/mdev5000/qlik_message/apperrors"

func validateMessage(message ModifyMessage, op string) error {
	if message.Message == "" {
		re := apperrors.Error{Op: op, EType: apperrors.ETInvalid}
		re.AddResponse(apperrors.FieldErrorResponse{
			Field: "message",
			Error: "Message field cannot be blank.",
		})
		return &re
	}
	return nil
}
