package errors

import (
	"errors"
)

type (
	// UnwrapError is a defined type for Unwrapping a single error.
	UnwrapError interface{ Unwrap() error }
	// UnwrapErrors is a defined type for Unwrapping a multiple errors.
	UnwrapErrors interface{ Unwrap() []error }
	Castable     interface{ As(any) bool }
)

// Check panics if err == nil.
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// Get panics if err == nil, returning t otherwise.
func Get[T any](t T, err error) T {
	Check(err)
	return t
}

// Do attempts to do the func, panicking if it errs.
// Do will panic if fn == nil.
func Do(fn func() error) { Check(fn()) }

func DoSet(fn func() error, err *error) {
	if e := fn(); e != nil {
		*err = Join(*err, e)
	}
}

// Unwraps functions the same as [errors.Unwrap], except that it unwraps [UnwrapErrors].
func Unwraps(err error) []error {
	if u, ok := err.(UnwrapErrors); ok {
		return u.Unwrap()
	}
	return nil
}

// To is a convenience function for [errors.As], allocating the target for you.
func To[E error](err error) (e E, ok bool) {
	ok = errors.As(err, &e)
	return
}

// Catch attempts to recover from a panic.
// If it succeeds, and the panic was and error, and the panic contains the specified error fn is run.
// If fn returns an error it will be thrown.
// If fn returns nil no more panics will be thrown.
func Catch(target error, fn func(original error) error) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			if Is(err, target) {
				if err = fn(err); err != nil {
					panic(err)
				}
				return
			}
		}
		panic(r)
	}
}
