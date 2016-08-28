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

func (e *Exec) Logger(logger golog.Logger) {
	e.logger = logger
}

func (e *Exec) execCommand(command string, arg ...string) error {
	execCommand := exec.Command(command, arg...)
	e.logger.Debugf("executing: %s %s", command, strings.Join(arg, " "))

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

	return execCommand.Run()
}

func (e *Exec) execCommand2(cmdline string) error {
	parts := strings.Split(cmdline, " ")
	validArgs := make([]string, 0)
	for _, p := range parts {
		if p != "" {
			validArgs = append(validArgs, p)
		}
	}
	return e.execCommand(validArgs[0], validArgs[1:]...)
}
