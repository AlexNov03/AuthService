package internal

import "errors"

var (
	ErrAlreadyExists  error = errors.New("internal error: already exists")
	ErrNotFound       error = errors.New("internal error: not found")
	ErrInternalServer error = errors.New("internal error: internal server error")
)

type InternalError struct {
	Message string
	Type    error
}

func (ie *InternalError) Error() string {
	return ie.Message
}

func NewInternalError(message string, errorType error) error {
	return &InternalError{
		Message: message,
		Type:    errorType,
	}
}
