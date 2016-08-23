package lib

import (
	"bufio"
	"github.com/go-errors/errors"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

func readFile(pathname string) (string, error) {
	bytes, err := ioutil.ReadFile(pathname)
	n := len(bytes)
	return string(bytes[:n]), err
}

func execCommand2(log *bufio.Writer, cmdline string) error {
	parts := strings.Split(cmdline, " ")
	validArgs := make([]string, 0)
	for _, p := range parts {
		if p != "" {
			validArgs = append(validArgs, p)
		}
	}
	return execCommand(log, validArgs[0], validArgs[1:]...)
}

func execCommand(log *bufio.Writer, command string, arg ...string) error {
	execCommand := exec.Command(command, arg...)
	stdout, _ := execCommand.StdoutPipe()
	stderr, _ := execCommand.StderrPipe()
	go io.Copy(log, stdout)
	go io.Copy(log, stderr)

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
