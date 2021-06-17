// Package handler
// Utility functions for implementing http.Handlers.
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mdev5000/messageappdemo/apperrors"
	"github.com/mdev5000/messageappdemo/logging"
	"github.com/pkg/errors"
)

const ContentTypeJson = "application/json; charset=UTF-8"

func DecodeJsonOrError(log *logging.Logger, op string, w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if r.Body == nil {
		appErr := apperrors.Error{Op: op, EType: apperrors.ETInvalid}
		appErr.AddResponse(apperrors.ErrorResponse("invalid json"))
		SendErrorResponse(log, op, w, &appErr)
		return false
	}
	d := json.NewDecoder(r.Body)
	if err := d.Decode(v); err != nil {
		if err.Error() == "http: request body too large" {
			appErr := apperrors.Error{Op: op, EType: apperrors.ETInvalid, Err: err, Stack: errors.WithStack(err)}
			appErr.AddResponse(apperrors.ErrorResponse("request body too large"))
			SendErrorResponse(log, op, w, &appErr)
			return false
		}
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

	if !apperrors.HasResponse(err) {
		w.WriteHeader(code)
		return
	}

	out, jsonErr := apperrors.ToJSON(err)
	if jsonErr != nil {
		log.LogFailedToEncode(op, err, jsonErr, errors.WithStack(jsonErr))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	contentTypeJson(w)
	writeData(op, log, w, out)
}

func writeData(op string, log *logging.Logger, w http.ResponseWriter, data []byte) bool {
	if _, errWrite := w.Write(data); errWrite != nil {
		log.LogError(&apperrors.Error{
			EType:     apperrors.ETInternal,
			Op:        op,
			Err:       errWrite,
			Stack:     errors.WithStack(errWrite),
			Responses: nil,
		})
		return false
	}
	return true
}

func EncodeJsonOrError(op string, log *logging.Logger, w http.ResponseWriter, r *http.Request, v interface{}) bool {
	contentTypeJson(w)
	// Don't return content if a HEAD request.
	if r.Method == "HEAD" {
		return true
	}
	d, jsonErr := json.Marshal(v)
	if jsonErr != nil {
		log.LogFailedToEncode(op, jsonErr, jsonErr, errors.WithStack(jsonErr))
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	return writeData(op, log, w, d)
}

func contentTypeJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", ContentTypeJson)
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
