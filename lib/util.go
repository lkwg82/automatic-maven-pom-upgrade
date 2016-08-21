package lib

import (
	"bufio"
	"github.com/go-errors/errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func readFile(pathname string) (string, error) {
	bytes, err := ioutil.ReadFile(pathname)
	n := len(bytes)
	return string(bytes[:n]), err
}

func execCommand(log *bufio.Writer, command string, arg ...string) error {
	execCommand := exec.Command(command, arg...)
	stdout, _ := execCommand.StdoutPipe()
	stderr, _ := execCommand.StderrPipe()
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	defer log.Flush()
	return execCommand.Run()
}

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
