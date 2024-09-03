package errors

import (
    "errors"
    "fmt"
    "runtime"
    "strings"
)

type (
    // UnwrapError is a defined type for Unwrapping a single error.
    UnwrapError interface{ Unwrap() error }
    // UnwrapErrors is a defined type for Unwrapping a multiple errors.
    UnwrapErrors interface{ Unwrap() []error }
    Castable     interface{ As(any) bool }
)

type Error struct {
    Text     string
    Errors   []error
    Line     int
    Filename string
}

func (err *Error) Error() string   { return Join(errors.New(err.Text), Join(err.Errors...)).Error() }
func (err *Error) Unwrap() []error { return err.Errors }

// New functions the same as [errors.New].
// If len(values) > 0, format will be used as a format string with values.
// The format string will have all '%w' replaced with '%v', and be used with [fmt.Sprintf].
// If any element in values is an error, it will be present in the slice returned by Unwraps.
func New(format string, values ...any) error {
    var err Error
    // Helpful information
    _, err.Filename, err.Line, _ = runtime.Caller(1)

    err.Text = format
    format = strings.ReplaceAll(format, "%w", "%v")
    if len(values) > 0 && strings.Contains(format, "%") {
        // Fix format string and format
        err.Text = fmt.Sprintf(format, values...)
    }

    // Get errors for Unwrap
    for _, e := range values {
        if e, ok := e.(error); ok {
            err.Errors = append(err.Errors, e)
        }
    }
    return &err
}

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

// OnFail return two functions. defr should be deferred immediately.
// If the caller returns before success is called, fn will be called.
func OnFail(fn func()) (success func(), defr func()) {
    fail := true
    return func() { fail = false }, func() {
        if fail {
            fn()
        }
    }
}
