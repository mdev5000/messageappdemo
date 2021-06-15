// Package handler
// Utility functions for implementing http.Handlers.
package handler

import (
	"encoding/json"
	"fmt"
	"github.com/mdev5000/qlik_message/apperrors"
	"github.com/mdev5000/qlik_message/logging"
	"github.com/pkg/errors"
	"net/http"
	"time"
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

	code := apperrors.StatusCode(err)
	w.WriteHeader(code)
	if !apperrors.HasResponse(err) {
		return
	}

	contentTypeJson(w)
	enc := json.NewEncoder(w)
	if jsonErr := apperrors.ToJSON(enc, err); jsonErr != nil {
		log.LogFailedToEncode(op, err, jsonErr, errors.WithStack(jsonErr))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func EncodeJsonOrError(op string, log *logging.Logger, w http.ResponseWriter, v interface{}) bool {
	contentTypeJson(w)
	enc := json.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		log.LogFailedToEncode(op, err, err, errors.WithStack(err))
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	return true
}

func contentTypeJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func SetETagInt(w http.ResponseWriter, v int) {
	w.Header().Set("ETag", fmt.Sprintf(`"%d"`, v))
}

func SetLastModified(w http.ResponseWriter, dt time.Time) {
	w.Header().Set("Last-Modified", LastModifiedFormat(dt))
}

func LastModifiedFormat(dt time.Time) string {
	return dt.Format("Mon, 02 Jan 2006 15:04:05 GMT")
}
