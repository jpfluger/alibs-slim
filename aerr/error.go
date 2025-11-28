package aerr

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Error wraps the built-in error interface to allow JSON marshaling.
type Error struct {
	error
}

// New creates a new Error instance from a non-nil error.
// Returns nil if the input error is nil.
func New(format string) *Error {
	return NewError(errors.New(format))
}

func Newf(format string, a ...interface{}) *Error {
	return NewError(fmt.Errorf(format, a...))
}

// NewError creates a new Error instance from a non-nil error.
// Returns nil if the input error is nil.
func NewError(err error) *Error {
	if err == nil {
		return nil
	}
	return &Error{error: err}
}

// NewErrorFromString creates a new Error instance from a string.
func NewErrorFromString(err string) *Error {
	return &Error{error: errors.New(err)}
}

// IsNil checks if the Error instance or the embedded error is nil.
func (err *Error) IsNil() bool {
	return err == nil || err.error == nil
}

// MarshalJSON customizes the JSON marshaling for Error.
func (err Error) MarshalJSON() ([]byte, error) {
	if err.error == nil {
		return []byte(`null`), nil
	}
	return json.Marshal(err.Error())
}

// UnmarshalJSON customizes the JSON unmarshaling for Error.
func (err *Error) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		err.error = nil
		return nil
	}

	var errMsg string
	if err := json.Unmarshal(b, &errMsg); err != nil {
		return err
	}

	err.error = errors.New(errMsg)
	return nil
}

// Error returns the string representation of the embedded error.
func (err *Error) Error() string {
	if err == nil || err.error == nil {
		return ""
	}
	return err.error.Error()
}

// ToError returns the embedded error.
func (err *Error) ToError() error {
	return err.error
}

// String returns the string representation of the embedded error.
func (err *Error) String() string {
	return err.Error()
}

// Unwrap returns the embedded error, allowing compatibility with errors.Unwrap.
func (err *Error) Unwrap() error {
	return err.error
}

// IsEqual compares the embedded error with another error.
// Returns true if both errors are the same or both are nil.
func (e *Error) IsEqual(err error) bool {
	if e == nil {
		return err == nil
	}
	return e.error == err || (err != nil && e.error.Error() == err.Error())
}
