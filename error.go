package validate

import (
	"errors"
	"fmt"

	"go.uber.org/multierr"
)

// Error represents validation errors, and implements Go's error type. Field
// indicates the struct field the validation error is relevant to, which is the
// full nested path relative to the top-level object being validated.
type Error struct {
	Field string
	Msg   string
	Err   error
}

func (s *Error) Error() string {
	msg := s.Msg
	if msg == "" && s.Err != nil {
		msg = s.Err.Error()
	}

	if msg == "" {
		msg = "unknown error"
	}

	if s.Field == "" {
		return msg
	}

	return fmt.Sprintf("%s: %s", s.Field, msg)
}

func (s *Error) Is(target error) bool {
	return errors.Is(s.Err, target)
}

func (s *Error) Unwrap() error {
	return s.Err
}

// Append combines two errors together into a single new error which internally
// keeps track of multiple errors via go.uber.org/multierr. If either error is a
// previously combined multierr, the returned error will be a flattened list of
// all errors.
func Append(errs error, err error) error {
	return multierr.Append(errs, err)
}

// AppendError appends a new *Error type to errs with the Msg field populated
// with the provided msg.
func AppendError(errs error, msg string) error {
	return multierr.Append(errs, &Error{Msg: msg})
}

// AppendFieldError appends a new *Error type to errs with Field and Msg
// populated with given field and msg values.
func AppendFieldError(errs error, field, msg string) error {
	return multierr.Append(errs, &Error{Field: field, Msg: msg})
}

// Errors returns a slice of all errors appended into the given error.
func Errors(err error) []error {
	return multierr.Errors(err)
}
