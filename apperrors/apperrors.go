package apperrors

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type errResponse struct {
	Errors []interface{} `json:"errors"`
}

// Error https://middlemost.com/failure-is-your-domain/
type Error struct {
	// The type of error. Determines how the error is handled, ex. ETInternal errors would result in internal logging
	// of the error and a 500 response code being returned to the user.
	EType string

	// Descriptor of the operation that caused the error. Ex. MessagesService.Create.
	Op string

	Err error

	// An error containing the full stack trace to the source of the error.
	Stack error

	// A list of error responses to return to the user.
	Responses []interface{}
}

func (e *Error) AddResponse(r interface{}) {
	e.Responses = append(e.Responses, r)
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error [%s] (%s): %s", e.EType, e.Op, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Is(err error) bool {
	switch err.(type) {
	case *Error:
		return true
	default:
		return errors.Is(e.Err, err)
	}
}

func StatusCode(err error) int {
	switch e := err.(type) {
	case *Error:
		switch e.EType {
		case ETInvalid:
			return http.StatusBadRequest
		case ETNotFound:
			return http.StatusNotFound
		default:
			return http.StatusInternalServerError
		}
	default:
		return http.StatusInternalServerError
	}
}

func HasResponse(err error) bool {
	switch e := err.(type) {
	case *Error:
		switch e.EType {
		case ETInvalid:
			return len(e.Responses) > 0
		}
	}
	return false
}

func ToJSON(encoder *json.Encoder, err error) error {
	switch e := err.(type) {
	case *Error:
		switch e.EType {
		case ETInvalid:
			return encoder.Encode(errResponse{e.Responses})
		default:
			return fmt.Errorf("error type for JSON, expected %s, but was %s", ETInvalid, e.EType)
		}
	default:
		return fmt.Errorf("error is not an application error, err: %w", err)
	}
}

func IsInternal(err error) bool {
	switch e := err.(type) {
	case *Error:
		return e.EType == ETInternal
	default:
		return true
	}
}

type FieldErrorResponse struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

func ErrorResponse(errMsg string) ErrResponse {
	return ErrResponse{Error: errMsg}
}
