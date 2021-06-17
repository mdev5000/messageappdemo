package messages

import (
	"testing"

	"github.com/mdev5000/messageappdemo/apperrors"
	"github.com/mdev5000/messageappdemo/logging"
	"github.com/stretchr/testify/require"
)

// Used to setup a service for testing. Note the repo is nil since this is intended for test the do not use the
// repository.
func tServiceNoRepo() *Service {
	return NewService(logging.NoLog(), nil)
}

func requireHasResponseErrors(t *testing.T, err error, expected interface{}) {
	aErr, ok := err.(*apperrors.Error)
	if !ok {
		t.Fatalf("expected err to be type *apperrors.Error but was %+v", err)
	}
	if len(aErr.Responses) == 0 {
		t.Fatalf("response is empty for %+v", aErr)
	}
	require.Equal(t, apperrors.ETInvalid, aErr.EType)
	require.Equal(t, expected, aErr.Responses[0])
}

func TestService_Create_runsValidation(t *testing.T) {
	_, err := tServiceNoRepo().Create(ModifyMessage{Message: ""})
	requireHasResponseErrors(t, err, apperrors.FieldErrorResponse{
		Field: "message",
		Error: "Message field cannot be blank.",
	})
}

func TestService_Update_runsValidation(t *testing.T) {
	_, err := tServiceNoRepo().Update(5, ModifyMessage{Message: ""})
	requireHasResponseErrors(t, err, apperrors.FieldErrorResponse{
		Field: "message",
		Error: "Message field cannot be blank.",
	})
}
