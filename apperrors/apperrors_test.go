package apperrors

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToJSON_canEncodeToJSON(t *testing.T) {
	e := Error{EType: ETInvalid}
	e.AddResponse(FieldErrorResponse{
		Field: "myfield",
		Error: "what went wrong",
	})
	b := bytes.NewBuffer(nil)
	enc := json.NewEncoder(b)
	require.NoError(t, ToJSON(enc, &e))
	require.Equal(t, `{"errors":[{"field":"myfield","error":"what went wrong"}]}`+"\n", b.String())
}

func TestError_canGetInternalError(t *testing.T) {
	internalErr := errors.New("internal")
	e := &Error{
		Op:    "Some.Operation",
		EType: ETInternal,
		Err:   internalErr,
	}
	require.True(t, errors.Is(e, &Error{}))
	require.EqualError(t, e, "Error [internal] (Some.Operation): internal")
	unwrapped := errors.Unwrap(e)
	require.EqualError(t, unwrapped, "internal")
	require.Same(t, unwrapped, internalErr)
}

func TestIsInternal_indicatesWhenAnInternalError(t *testing.T) {
	require.True(t, IsInternal(&Error{EType: ETInternal}))
	require.False(t, IsInternal(&Error{EType: ETInvalid}))
}

func TestIsInternal_indicatesNonApplicationErrorsAsInternal(t *testing.T) {
	require.True(t, IsInternal(fmt.Errorf("some error")))
}
