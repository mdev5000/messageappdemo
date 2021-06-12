package data

import (
	"fmt"
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

func repoError(repo string, err error) error {
	return errors.WithStack(RepositoryError{RepositoryIdentifier: repo, Err: err})
}

type RepositoryError struct {
	RepositoryIdentifier string
	Err                  error
}

func (e RepositoryError) Error() string {
	return fmt.Sprintf("%s repository: %s", e.RepositoryIdentifier, e.Err.Error())
}

func (e RepositoryError) Is(target error) bool {
	switch target.(type) {
	case RepositoryError:
		return true
	default:
		return false
	}
}
