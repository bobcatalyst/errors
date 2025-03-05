package errors

import (
    "errors"
    "fmt"
    "github.com/bobcatalyst/debug"
    "runtime"
    "strconv"
    "strings"
)

type (
    // UnwrapError is a defined type for Unwrapping a single error.
    UnwrapError interface{ Unwrap() error }
    // UnwrapErrors is a defined type for Unwrapping a multiple errors.
    UnwrapErrors interface{ Unwrap() []error }
    Castable     interface{ As(any) bool }
    Comparable   interface{ Is(error) bool }
)

type Error struct {
    Text     string
    Values   []any
    Errors   []error
    Line     int
    Filename string
}

func (err *Error) Error() string {
    var sb strings.Builder
    if debug.Debug {
        sb.WriteRune('{')
        sb.WriteString(err.Filename)
        sb.WriteRune(':')
        sb.WriteString(strconv.Itoa(err.Line))
        sb.WriteString("} ")
    }

    _, _ = fmt.Fprintf(&sb, err.Text, err.Values...)

    for _, e := range err.Errors {
        _, _ = fmt.Fprintf(&sb, " \\\n\t%v", e)
    }

    return sb.String()
}

func (err *Error) Unwrap() []error { return err.Errors }

// New functions the same as [errors.New].
// If len(values) > 0, format will be used as a format string with values.
// The format string will have all '%w' replaced with '%v', and be used with [fmt.Sprintf].
// If any element in values is an error, it will be present in the slice returned by Unwraps.
func New(format string, values ...any) error {
    err := Error{
        Text:   strings.ReplaceAll(format, "%w", "%v"),
        Values: values,
    }

    if debug.Debug {
        _, err.Filename, err.Line, _ = runtime.Caller(1)
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

// Must panics if err == nil, returning t otherwise.
func Must[T any](t T, err error) T {
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

// OnFail runs the functions supplied to onFail until success is called.
// Usage:
//   ok, onFail := errors.OnFail()
//   f := errors.Must(os.Create('file.txt'))
//   defer onFail(func() {
//     if err := f.Close(); err != nil {
//       slog.Error("failed to write file", "error", err)
//     }
//   })
//   errors.Must(f.Write([]byte("Hello World"))
//   ok()
func OnFail() (success func(), onFail func(func())) {
    fail := true
    return func() { fail = false }, func(fn func()) {
        if fail {
            fn()
        }
    }
}
