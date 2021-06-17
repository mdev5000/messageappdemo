package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"unsafe"

	"github.com/mdev5000/messageappdemo/apperrors"
	"github.com/mdev5000/messageappdemo/logging"
	"github.com/stretchr/testify/require"
)

func TestSendErrorResponse_internalErrorReturns500(t *testing.T) {
	log := logging.NoLog()
	rr := httptest.NewRecorder()
	err := &apperrors.Error{
		EType: apperrors.ETInternal,
		Err:   errors.New("some error"),
	}
	SendErrorResponse(log, "op", rr, err)
	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Nil(t, rr.Body.Bytes())
}

func TestSendErrorResponse_nonAppErrorReturns500(t *testing.T) {
	log := logging.NoLog()
	rr := httptest.NewRecorder()
	SendErrorResponse(log, "op", rr, errors.New("my error"))
	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Nil(t, rr.Body.Bytes())
}

func TestSendErrorResponse_invalidErrorReturnsErrorResponseWhenResponse(t *testing.T) {
	log := logging.NoLog()
	rr := httptest.NewRecorder()
	err := &apperrors.Error{
		EType: apperrors.ETInvalid,
		Err:   errors.New("some error"),
	}
	err.AddResponse(apperrors.ErrorResponse("something happened"))
	SendErrorResponse(log, "op", rr, err)
	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Equal(t, `{"errors":[{"error":"something happened"}]}`, rr.Body.String())
}

func TestSendErrorResponse_returns404WhenNotFound(t *testing.T) {
	log := logging.NoLog()
	rr := httptest.NewRecorder()
	err := &apperrors.Error{EType: apperrors.ETNotFound}
	SendErrorResponse(log, "op", rr, err)
	require.Equal(t, http.StatusNotFound, rr.Code)
	require.Nil(t, rr.Body.Bytes())
}

func TestSendErrorResponse_returns500WhenCannotEncodeErrorMessage(t *testing.T) {
	log := logging.NoLog()
	rr := httptest.NewRecorder()
	err := &apperrors.Error{
		EType: apperrors.ETInvalid,
		Err:   errors.New("some error"),
	}
	err.AddResponse(unsafe.Pointer(nil))
	SendErrorResponse(log, "op", rr, err)
	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestEncodeJsonOrError_canEncode(t *testing.T) {
	log := logging.NoLog()
	r, err := http.NewRequest("GET", "/", bytes.NewBuffer(nil))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	require.True(t, EncodeJsonOrError("op", log, rr, r, "encode this"))
	require.Equal(t, ContentTypeJson, rr.Header().Get("Content-Type"))
	require.Equal(t, `"encode this"`, rr.Body.String())
}

func TestEncodeJsonOrError_returns500whenEncodingFails(t *testing.T) {
	log := logging.NoLog()
	r, err := http.NewRequest("GET", "/", bytes.NewBuffer(nil))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	require.False(t, EncodeJsonOrError("op", log, rr, r, unsafe.Pointer(nil)))
	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
