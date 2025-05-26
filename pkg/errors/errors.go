// Package errors provides a single package for all error-related and extends
// the standard library's errors package.
package errors

import (
	"errors"
)

// New returns an error with the supplied message.
func New(text string) error {
	return errors.New(text)
}

// As calls the standard library's errors.As function.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Unwrap calls the standard library's errors.Unwrap function.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}
