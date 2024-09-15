package dinfra

import "fmt"

type (
	ErrorCode int

	CodeError struct {
		Msg  string
		Code ErrorCode
	}
)

func (err *CodeError) Error() string {
	return fmt.Sprintf("[%d] %s", err.Code, err.Msg)
}

func IsError(bad bool, err error) error {
	if bad {
		return err
	}
	return nil
}

func IsErrorf(bad bool, format string, a ...any) error {
	if bad {
		return fmt.Errorf(format, a...)
	}
	return nil
}

func HasErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func NewError(code ErrorCode, msg string) *CodeError {
	return &CodeError{
		Msg:  msg,
		Code: code,
	}
}

func NewErrorf(code ErrorCode, format string, a ...any) *CodeError {
	return &CodeError{
		Msg:  fmt.Sprintf(format, a...),
		Code: code,
	}
}
