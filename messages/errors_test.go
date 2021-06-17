package messages

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIdMissingError(t *testing.T) {
	err := IdMissingError{"myrepo", 5}
	require.True(t, errors.Is(err, IdMissingError{}))
	require.False(t, errors.Is(err, errors.New("another error")), "correctly indicates it is not other errors")
	require.EqualError(t, err, "myrepo: no row in result with id 5")
}
