package lib

import (
	"bufio"
	"github.com/alexcesaro/log/golog"
	"io"
	"os/exec"
	"strings"
	"path"
	"runtime"
	"bytes"
	"os"
)

type Exec struct {
	logger golog.Logger
	Cmd    string
}

func NewExec(logger golog.Logger, cmd string) *Exec {
	return &Exec{
		logger:logger,
		Cmd:cmd,
	}
}

func (e *Exec) Command(arg ...string) *exec.Cmd {
	execCommand := exec.Command(e.Cmd, arg...)
	if e.logger.LogDebug() {
		name := e.getFirstPublicLibCaller()
		e.logger.Debugf("call from: %s ", name)
		e.logger.Debugf("executing: %s %s", e.Cmd, strings.Join(arg, " "))
	}
	return execCommand
}

func (e *Exec) CommandRun(arg ...string) error {
	execCommand := e.Command(arg...)

	if e.logger.LogDebug() {
		stdout, _ := execCommand.StdoutPipe()
		stderr, _ := execCommand.StderrPipe()

		go e.copyToLog(stdout, "stdout")
		go e.copyToLog(stderr, "stderr")
	}

	err := execCommand.Run()
	e.DebugErr(err)
	return err
}

func (e *Exec) CommandRunExitOnErr(arg ...string) {
	err := e.CommandRun(arg...)
	if err != nil {
		e.logger.Emergency(err)
		os.Exit(1)
	}
}

func (e *Exec) DebugErr(err error) {
	if err != nil {
		exitError := err.(*exec.ExitError)
		e.logger.Debugf(" exit code: %s ", exitError.Error())
	}
}

func (e *Exec) DebugStdout(output []byte) {
	if e.logger.LogDebug() {
		e.copyToLog(bytes.NewReader(output), "stdout")
	}
}

func (e *Exec) DebugStdoutErr(output []byte, err error) {
	if e.logger.LogDebug() {
		e.copyToLog(bytes.NewReader(output), "stdout")
		e.DebugErr(err)
	}
}

func (e *Exec) DebugStderr(output []byte) {
	if e.logger.LogDebug() {
		e.copyToLog(bytes.NewReader(output), "stderr")
	}
}

func (e *Exec) copyToLog(rc io.Reader, channel string) {
	in := bufio.NewScanner(rc)
	for in.Scan() {
		e.logger.Debugf("   %s: %s", channel, in.Text())
	}
}

func (e *Exec) execCommand2(cmdline string) error {
	parts := strings.Split(cmdline, " ")
	validArgs := make([]string, 0)
	for _, p := range parts {
		if p != "" {
			validArgs = append(validArgs, p)
		}
	}
	return e.CommandRun(validArgs[0:]...)
}

func (e *Exec) getFirstPublicLibCaller() string {
	var chain string;
	var last *runtime.Func
	for i := 0; i < 10; i++ {
		pc, _, _, _ := runtime.Caller(i)
		curr := runtime.FuncForPC(pc)

		baseCurr := strings.Replace(path.Base(curr.Name()), "lib.", "", 1)
		baseLast := strings.Replace(path.Base(last.Name()), "lib.", "", 1)
		if strings.Index(baseCurr, "(") == 0 {
			if strings.Index(baseCurr, "(*Exec") != 0 {
				if chain == "" {
					chain = baseCurr
				} else {
					chain = baseCurr + " -> " + chain
				}
			}
		}
		if strings.Index(baseCurr, "(") != 0 && strings.Index(baseLast, "(") == 0 {
			if strings.Index(baseLast, "(*Exec") != 0 {
				//fmt.Println(chain)
				return chain
			}
		}
		if curr.Name() == "testing.tRunner" {
			return baseLast
		}

		last = curr
	}
	//os.Exit(1)
	return path.Base(last.Name())
}
