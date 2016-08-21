package lib

import (
	"io/ioutil"
	"bufio"
	"os/exec"
	"io"
)

func readFile(pathname string) (string, error) {
	bytes, err := ioutil.ReadFile(pathname)
	n := len(bytes)
	return string(bytes[:n]), err
}

func execCommand(log *bufio.Writer, command string, arg ...string) (error) {
	execCommand := exec.Command(command, arg...)
	stdout, _ := execCommand.StdoutPipe()
	stderr, _ := execCommand.StderrPipe()
	go io.Copy(log, stdout)
	go io.Copy(log, stderr)

	defer log.Flush()
	return execCommand.Run()
}
