package apperrors_test

import (
	"errors"
	"fmt"
	"github.com/mdev5000/qlik_message/apperrors"
	errors2 "github.com/pkg/errors"
)

func ExampleError_userErrors() {
	origErr := errors.New("original error")

	// This is an error who message should be forwarded on the user.
	// Important! Only ETInvalid errors have user responses. If the message does not have an ETInvalid type then
	// it will not be returned to the user, even if it has a response added.
	aErr := apperrors.Error{
		EType: apperrors.ETInvalid,
		Op:    "identifier to track origin of the error",
		Err:   origErr,
		Stack: errors2.WithStack(origErr),
	}
	aErr.AddResponse(apperrors.ErrorResponse("you did this wrong"))

	var err error = &aErr

	fmt.Println(apperrors.HasResponse(err))
	fmt.Println(apperrors.StatusCode(err))

	d, jsonErr := apperrors.ToJSON(err)
	if jsonErr != nil {
		panic(jsonErr)
	}
	fmt.Print(string(d))

	// Output:
	// true
	// 400
	// {"errors":[{"error":"you did this wrong"}]}
}

func ExampleError_internalErrors() {
	origErr := errors.New("original error")

	// This is an error who message should be forwarded on the user
	var err error = &apperrors.Error{
		EType: apperrors.ETInternal,
		Op:    "identifier to track origin of the error",
		Err:   origErr,
		Stack: errors2.WithStack(origErr),
	}

	fmt.Println(apperrors.HasResponse(err))
	fmt.Println(apperrors.StatusCode(err))
	fmt.Println(errors.Is(err, origErr))

	// Output:
	// false
	// 500
	// true
}

func ExampleFieldErrorResponse() {
	origErr := errors.New("original error")

	// This is an error who message should be forwarded on the user
	aErr := apperrors.Error{
		EType: apperrors.ETInvalid,
		Op:    "identifier to track origin of the error",
		Err:   origErr,
		Stack: errors2.WithStack(origErr),
	}
	aErr.AddResponse(apperrors.FieldErrorResponse{
		Field: "myField",
		Error: "this went wrong",
	})

	var err error = &aErr
	d, jsonErr := apperrors.ToJSON(err)
	if jsonErr != nil {
		panic(jsonErr)
	}
	fmt.Print(string(d))

	// Output:
	// {"errors":[{"field":"myField","error":"this went wrong"}]}
}

func ExampleErrorResponse() {
	origErr := errors.New("original error")

	// This is an error who message should be forwarded on the user
	aErr := apperrors.Error{
		EType: apperrors.ETInvalid,
		Op:    "identifier to track origin of the error",
		Err:   origErr,
		Stack: errors2.WithStack(origErr),
	}
	aErr.AddResponse(apperrors.ErrorResponse("you did this wrong"))

	var err error = &aErr
	d, jsonErr := apperrors.ToJSON(err)
	if jsonErr != nil {
		panic(jsonErr)
	}
	fmt.Print(string(d))

	// Output:
	// {"errors":[{"error":"you did this wrong"}]}
}
