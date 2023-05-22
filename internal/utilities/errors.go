package utilities

import (
	"fmt"

	"github.com/pkg/errors"
)

// This error is used to denote something risky and unexpected happening in the system.
// Ideally this should never throw, but if it does it means something really weird is going on.
// In most cases, BadErrors indicates code that does not have a test.
// Initially, I left code untested when I felt the error can never happen.
// If this was thrown in production, it means something changed or my above feeling was incorrect.
type BadError struct {
	message string
}

func NewBadError(message string) *BadError {
	return &BadError{
		message: message,
	}
}

func WrapBadError(err error, message string) error {
	badError := NewBadError(message)
	return errors.Wrap(err, badError.Error())
}

func (e *BadError) Error() string {
	return fmt.Sprintf("THIS IS BAD: %s", e.message)
}
