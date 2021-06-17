package apperrors

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestToJSON_canEncodeToJSON(t *testing.T) {
	e := Error{EType: ETInvalid}
	e.AddResponse(FieldErrorResponse{
		Field: "myfield",
		Error: "what went wrong",
	})
	d, err := ToJSON(&e)
	require.NoError(t, err)
	require.Equal(t, `{"errors":[{"field":"myfield","error":"what went wrong"}]}`, string(d))
}

func TestHasResponse_trueWhenInvalidAndContainsAResponse(t *testing.T) {
	e := Error{EType: ETInvalid}
	e.AddResponse(FieldErrorResponse{
		Field: "myfield",
		Error: "what went wrong",
	})
	require.True(t, HasResponse(&e))
}

func TestHasResponse_falseWhenNotInvalid(t *testing.T) {
	cases := []struct {
		name  string
		etype string
	}{
		{"internal", ETInternal},
		{"not found", ETNotFound},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			e := Error{EType: c.etype}
			e.AddResponse(FieldErrorResponse{
				Field: "myfield",
				Error: "what went wrong",
			})
			require.False(t, HasResponse(&e))
		})
	}
}

func TestError_canGetInternalError(t *testing.T) {
	internalErr := errors.New("internal")
	e := &Error{
		Op:    "Some.Operation",
		EType: ETInternal,
		Err:   internalErr,
	}
	require.True(t, errors.Is(e, &Error{}))
	require.True(t, errors.Is(e, internalErr))
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

func TestStatusCode_returnsCorrectCode(t *testing.T) {
	cases := []struct {
		name string
		code int
		err  error
	}{
		{name: "internal", code: http.StatusInternalServerError, err: &Error{EType: ETInternal}},
		{name: "not found", code: http.StatusNotFound, err: &Error{EType: ETNotFound}},
		{name: "invalid", code: http.StatusBadRequest, err: &Error{EType: ETInvalid}},
		{name: "non app error", code: http.StatusInternalServerError, err: fmt.Errorf("some error")},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.code, StatusCode(c.err))
		})
	}
}
