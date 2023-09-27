package entities

import "fmt"

type UnsupportedEventError struct {
	event fmt.Stringer
}

func NewUnsupportedEventError(s fmt.Stringer) error {
	return &UnsupportedEventError{
		event: s,
	}
}

func (e UnsupportedEventError) Error() string {
	return fmt.Sprintf("unsupported event: %s", e.event)
}
