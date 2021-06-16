package apperrors

const (
	// ETInvalid is returned when user provided input is invalid.
	ETInvalid = "invalid"

	// ETInternal is returned for internal application errors that should not be forwarded to the client.
	ETInternal = "internal"

	// ETNotFound is returned when a resource could not be found with the given identifier.
	ETNotFound = "not found"
)
