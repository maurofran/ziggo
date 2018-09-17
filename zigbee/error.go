package zigbee

// NewError will create a new zigbee error.
func NewError(message string) error {
	return &customError{message: message}
}

// NewErrorWithCause will create a new zigbee error with supplied cause
func NewErrorWithCause(message string, cause error) error {
	return &customError{message: message, cause: cause}
}

// IsError check if the supplied error is a zigbee error.
func IsError(err error) bool {
	_, ok := err.(*customError)
	return ok
}

type customError struct {
	message string
	cause   error
}

func (e *customError) Error() string {
	return e.message
}
