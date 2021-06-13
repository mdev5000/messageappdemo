package apperrors

const (
	// ETInvalid is returned user provided input is invalid.
	ETInvalid = "invalid"

	// ETInternal is returned for internal application errors that should not be forwarded to the client.
	ETInternal = "internal"
)
