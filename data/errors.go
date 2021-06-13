package data

import (
	"fmt"
	"github.com/mdev5000/qlik_message/apperrors"
	"github.com/pkg/errors"
)

func idMissingError(repo string, id int64) error {
	return errors.WithStack(IdMissingError{RepositoryIdentifier: repo, Id: id})
}

type IdMissingError struct {
	RepositoryIdentifier string
	Id                   int64
}

func (e IdMissingError) Error() string {
	return fmt.Sprintf("%s repository: no rows in result for get by id with id %d", e.RepositoryIdentifier, e.Id)
}

func (e IdMissingError) Is(target error) bool {
	switch target.(type) {
	case IdMissingError:
		return true
	default:
		return false
	}
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
