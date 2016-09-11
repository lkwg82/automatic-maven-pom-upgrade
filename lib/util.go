package lib

import (
	"github.com/go-errors/errors"
)

type wrapError struct {
	error
	msg string
}

func newWrapError(err error, message string) *wrapError {
	if err == nil {
		panic("no nil error allowed")
	}
	we := &wrapError{
		error: errors.Wrap(err, 1),
		msg:   message,
	}

	if message == "" {
		we.msg = we.error.Error()
	}

	return we
}

func newWrapError2(message string) *wrapError {
	return newWrapError(errors.New(message), "")
}

func (w *wrapError) CausedErr() error {
	return w.error
}

func (w *wrapError) string() string {
	return w.msg
}
