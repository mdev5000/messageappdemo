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

func DecodeJsonOrError(w http.ResponseWriter, r *http.Request, v interface{}) error {
	d := json.NewDecoder(r.Body)
	if err := d.Decode(v); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// @todo probably send message indicating why this request is bad.
		return err
	}
	return nil
}

func SendErrorResponse(log *logging.Logger, op string, err error, w http.ResponseWriter) {
	if apperrors.IsInternal(err) {
		log.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		enc := json.NewEncoder(w)
		if jsonErr := apperrors.ToJSON(enc, err); jsonErr != nil {
			log.LogFailedToEncode(op, err, jsonErr, errors.WithStack(jsonErr))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
}
