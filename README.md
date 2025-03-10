Drop in replacement for the standard errors library.

Extra funcs:

- Check
  - panics if the given error is not nil
- Get
  - panics if the error in a given tuple is not nil, returns the value otherwise
- Do
  - runs a func returning an error, panicing if the returned error is not nil
- DoSet
  - runs a function setting the given error pointer if the returned err is not nil
- Unwraps
  - unwraps multiple errors
- To
  - converts an error to a concrete type
- OnFail
  - recovers from a panic or early return

Extra interfaces:

- UnwrapError
  - an error implementing `Unwrap() error`
- UnwrapErrors
  - an error implementing `Unwrap() []error`
- Castable
  - an error implementing `As(any) bool`
- Comparable
  - an error implementing `Is(error) bool`