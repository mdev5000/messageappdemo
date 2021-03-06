package messages

import "fmt"

type IdMissingError struct {
	Op string
	Id int64
}

func (e IdMissingError) Error() string {
	return fmt.Sprintf("%s: no row in result with id %d", e.Op, e.Id)
}

func (e IdMissingError) Is(target error) bool {
	switch target.(type) {
	case IdMissingError:
		return true
	default:
		return false
	}
}
