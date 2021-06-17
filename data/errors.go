package data

import (
	"github.com/mdev5000/messageappdemo/apperrors"
	"github.com/mdev5000/messageappdemo/messages"
	"github.com/pkg/errors"
)

func idMissingError(op string, id int64) error {
	return messages.IdMissingError{Op: op, Id: id}
}

func repoError2(op string, err error) error {
	return repoError(op, err, err)
}

func repoError(op string, err, errOrig error) error {
	return &apperrors.Error{
		EType: apperrors.ETInternal,
		Op:    op,
		Err:   err,
		Stack: errors.WithStack(errOrig),
	}
}
