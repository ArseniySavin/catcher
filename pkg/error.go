package pkg

import (
	"fmt"
	"github.com/ArseniySavin/catcher/pkg/internal"
)

var BaseError = &Error{}

type Error struct {
	Code string
	Msg  string
	Stk  string
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

func (e *Error) Stack() string {
	return fmt.Sprintf("Stack: %s", e.Stk)
}

func (e *Error) NewCode(code string) *Error {
	e.Code = code
	return e

}

func (e *Error) Throw(msg string) *Error {
	e.Msg = msg
	e.Stk = internal.MarshalStruct(internal.CallInfo(2))
	return e

}
