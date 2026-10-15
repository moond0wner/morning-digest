package core_errors

import (
	"net/http"
)

type Type int

const (
	ErrUserInput Type = iota
	ErrExternalAPI
	ErrInternal
	ErrNotFound
)

type Error struct {
	Type Type
	Err  error
	Msg  string
}

func (t Type) String() string {
	switch t {
	case ErrUserInput:
		return "user input error"
	case ErrNotFound:
		return "not found"
	case ErrInternal:
		return "internal error"
	case ErrExternalAPI:
		return "external api error"
	}
	return "unknown error"
}

func (e *Error) Error() string {
	if e.Err == nil {
		if e.Msg != "" {
			return e.Msg
		}
		return e.Type.String()
	}
	if e.Msg != "" {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Type.String() + ": " + e.Err.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Type == t.Type
}

func New(t Type, msg string) *Error {
	return &Error{Type: t, Msg: msg}
}

func Wrap(t Type, err error, msg string) *Error {
	return &Error{Type: t, Err: err, Msg: msg}
}

func CheckStatusCode(statusCode int) Type {
	switch statusCode {
	case http.StatusBadRequest:
		return ErrUserInput
	case http.StatusNotFound:
		return ErrNotFound
	default:
		return ErrExternalAPI
	}
}
