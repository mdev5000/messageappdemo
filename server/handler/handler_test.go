package handler

import (
	"errors"
	"github.com/mdev5000/qlik_message/apperrors"
	"github.com/mdev5000/qlik_message/logging"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
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
	require.Equal(t, `{"errors":[{"error":"something happened"}]}`+"\n", rr.Body.String())
}

func TestSendErrorResponse_returns404WhenNotFound(t *testing.T) {
	log := logging.NoLog()
	rr := httptest.NewRecorder()
	err := &apperrors.Error{EType: apperrors.ETNotFound}
	SendErrorResponse(log, "op", rr, err)
	require.Equal(t, http.StatusNotFound, rr.Code)
	require.Nil(t, rr.Body.Bytes())
}
