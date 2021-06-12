package data

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRepositoryError(t *testing.T) {
	err := repoError("myrepo", errors.New("my error"))
	require.True(t, errors.Is(err, RepositoryError{}))
	require.False(t, errors.Is(err, IdMissingError{}), "correctly indicates it is not other errors")
	require.EqualError(t, err, "myrepo repository: my error")
}

func TestIdMissingError(t *testing.T) {
	err := idMissingError("myrepo", 5)
	require.True(t, errors.Is(err, IdMissingError{}))
	require.False(t, errors.Is(err, RepositoryError{}), "correctly indicates it is not other errors")
	require.EqualError(t, err, "myrepo repository: no rows in result for get by id with id 5")
}
