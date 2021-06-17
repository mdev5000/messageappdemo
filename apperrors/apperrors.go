// Package apperrors is the core error handling package for the application.
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

// Error is the primary application error object within the application. It manages useful information for logging and
// can be used to determine the message type (ex. can be forwarded to the user or should be logged at return 500).
// It is based on this article https://middlemost.com/failure-is-your-domain/.
//
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

// StatusCode returns the HTTP status code for the given error. For most errors this will be 500, usually only errors
// returning error responses (see HasResponse) will have a non-500 response codes.
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

// HasResponse indicates whether the error has a user response. A user response is an error message that can be returned
// to the user (usually via JSON response) detailing to them what went wrong. Internal errors and non-application errors
// will always return false.
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

// ToJSON encodes an error to JSON. It will return an error if the error does not support sending a message to the end
// user. An error that returns false to HasResponse cannot use ToJSON.
func ToJSON(err error) ([]byte, error) {
	switch e := err.(type) {
	case *Error:
		switch e.EType {
		case ETInvalid:
			return json.Marshal(errResponse{e.Responses})
		default:
			return nil, fmt.Errorf("error type for JSON, expected %s, but was %s", ETInvalid, e.EType)
		}
	default:
		return nil, fmt.Errorf("error is not an application error, err: %w", err)
	}
}

// IsInternal indicates if an error is intended to be shown to the end user.
func IsInternal(err error) bool {
	switch e := err.(type) {
	case *Error:
		return e.EType == ETInternal
	default:
		return true
	}
}

// FieldErrorResponse is a error response meant for the end user indicating an error with a specific field.
type FieldErrorResponse struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ErrResponse is a error response meant for the end user detailing the error that occurred. See ErrorResponse for
// example/
type ErrResponse struct {
	Error string `json:"error"`
}

// ErrorResponse creates an ErrResponse with a given message.
func ErrorResponse(errMsg string) ErrResponse {
	return ErrResponse{Error: errMsg}
}
