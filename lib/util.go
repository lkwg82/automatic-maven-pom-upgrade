package lib

import (
	"github.com/go-errors/errors"
)

type WrapError struct {
	error
	msg string
}

func NewWrapError(err error, message string) *WrapError {
	if err == nil {
		panic("no nil error allowed")
	}
	we := &WrapError{
		error: errors.Wrap(err, 1),
		msg:   message,
	}

	if message == "" {
		we.msg = we.error.Error()
	}

	return we
}

func NewWrapError2(message string) *WrapError {
	return NewWrapError(errors.New(message), "")
}

func (w *WrapError) CausedErr() error {
	return w.error
}

func (w *WrapError) string() string {
	return w.msg
}
