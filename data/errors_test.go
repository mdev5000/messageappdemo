package data

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIdMissingError(t *testing.T) {
	err := idMissingError("myrepo", 5)
	require.True(t, errors.Is(err, IdMissingError{}))
	require.False(t, errors.Is(err, errors.New("another error")), "correctly indicates it is not other errors")
	require.EqualError(t, err, "myrepo repository: no rows in result for get by id with id 5")
}
