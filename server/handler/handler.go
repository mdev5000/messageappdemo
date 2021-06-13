// Package handler
// Utility functions for implementing http.Handlers.
package handler

import (
	"encoding/json"
	"github.com/mdev5000/qlik_message/apperrors"
	"github.com/mdev5000/qlik_message/logging"
	"github.com/pkg/errors"
	"net/http"
)

func DecodeJsonOrError(log *logging.Logger, op string, w http.ResponseWriter, r *http.Request, v interface{}) bool {
	d := json.NewDecoder(r.Body)
	if err := d.Decode(v); err != nil {
		appErr := apperrors.Error{Op: op, EType: apperrors.ETInvalid, Err: err, Stack: errors.WithStack(err)}
		appErr.AddResponse(apperrors.ErrorResponse("invalid json"))
		SendErrorResponse(log, op, w, &appErr)
		return false
	}
	return true
}

func SendErrorResponse(log *logging.Logger, op string, w http.ResponseWriter, err error) {
	if apperrors.IsInternal(err) {
		log.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	enc := json.NewEncoder(w)
	if jsonErr := apperrors.ToJSON(enc, err); jsonErr != nil {
		log.LogFailedToEncode(op, err, jsonErr, errors.WithStack(jsonErr))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func EncodeJsonOrError(op string, log *logging.Logger, w http.ResponseWriter, v interface{}) bool {
	enc := json.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		log.LogFailedToEncode(op, err, err, errors.WithStack(err))
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	return true
}
