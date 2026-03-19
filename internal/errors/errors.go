package errors

import (
	"fmt"
)

type InternalError struct {
	code    int
	message string
}

func NewInternalError() *InternalError {
	return &InternalError{}
}

func (i *InternalError) WithCode(code int) *InternalError {
	i.code = code
	return i
}

func (i *InternalError) WithMessage(msg string) *InternalError {
	i.message = msg
	return i
}

func (i *InternalError) Error() string {
	return i.message
}

func (i *InternalError) Code() int {
	return i.code
}

var (
	ErrBadRequest     = fmt.Errorf("bad request")
	ErrInternalServer = fmt.Errorf("internal error")
	ErrEventNotFound  = fmt.Errorf("event is not found")
)
