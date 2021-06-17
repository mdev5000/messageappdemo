package handler

import (
	"github.com/mdev5000/messageappdemo/apperrors"
	"github.com/pkg/errors"
)

func ResponseError(op string) apperrors.Error {
	err := errors.New("invalid data")
	return apperrors.Error{
		EType: apperrors.ETInvalid,
		Op:    op,
		Err:   err,
		Stack: errors.WithStack(err),
	}
}
