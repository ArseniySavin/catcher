package pkg

import (
	"fmt"
)

var BaseError = &Error{}

type Error struct {
	Code string
	Msg  string
}

func (e *Error) New(code, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

func (e *Error) Unwrap() error {
	return nil // why nil? Because, if we return *Error. We get circle call.
}

func (e *Error) Error() string {
	return fmt.Sprintf("Code: %s Msg: %s", e.Code, e.Msg)
}

func (e *Error) NewCode(code string) *Error {
	e.Code = code
	return e

}

func (e *Error) Throw(msg string) *Error {
	e.Msg = msg
	return e

}
