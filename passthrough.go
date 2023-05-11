// Package errors passes through the existing [errors] library.
// The Following methods are added:
package errors

import "errors"

// New functions the same as [errors.New].
func New(text string) error { return errors.New(text) }

// Is functions the same as [errors.Is].
func Is(err, target error) bool { return errors.Is(err, target) }

// As functions the same as [errors.As].
func As(err error, target any) bool { return errors.As(err, &target) }

// Join functions the same as [errors.Join].
func Join(errs ...error) error { return errors.Join(errs...) }

// Unwrap functions the same as [errors.Unwrap].
func Unwrap(err error) error { return errors.Unwrap(err) }
