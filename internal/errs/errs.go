package errs

import "fmt"

type CustomError struct {
	Code    ErrorCodeEnum
	Message string
}

func (c *CustomError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", c.Code, c.Message)
}

func New(code ErrorCodeEnum, msg string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: msg,
	}
}

type ErrorCodeEnum int

const (
	ErrUnknownCode ErrorCodeEnum = iota
)

var (
	ErrUnknown = New(ErrUnknownCode, "unknown error")
)
