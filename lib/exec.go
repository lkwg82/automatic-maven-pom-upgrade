package lib

import (
	"bufio"
	"github.com/alexcesaro/log/golog"
	"io"
	"os/exec"
	"strings"
)

type Exec struct {
	logger golog.Logger
}

func NewExec(logger golog.Logger) *Exec {
	return &Exec{
		logger:logger,
	}
}

func (e *Exec) Logger(logger golog.Logger) {
	e.logger = logger
}

func (e *Exec) Command(command string, arg ...string) *exec.Cmd {
	execCommand := exec.Command(command, arg...)
	e.logger.Debugf("executing: %s %s", command, strings.Join(arg, " "))
	return execCommand
}

func (e *Exec) ExecCommand(command string, arg ...string) error {
	execCommand := e.Command(command, arg...)

	if e.logger.LogDebug() {
		stdout, _ := execCommand.StdoutPipe()
		stderr, _ := execCommand.StderrPipe()

		copyToLog := func(rc io.ReadCloser) {
			in := bufio.NewScanner(rc)
			for in.Scan() {
				e.logger.Debug(in.Text())
			}
		}

		go copyToLog(stdout)
		go copyToLog(stderr)
	}

	err := execCommand.Run()
	if err != nil {
		exitError := err.(*exec.ExitError)
		e.logger.Debugf(" exit code: %s ", exitError.Error())
	}

	return err
}

func (e *Exec) execCommand2(cmdline string) error {
	parts := strings.Split(cmdline, " ")
	validArgs := make([]string, 0)
	for _, p := range parts {
		if p != "" {
			validArgs = append(validArgs, p)
		}
	}
	return e.ExecCommand(validArgs[0], validArgs[1:]...)
}
