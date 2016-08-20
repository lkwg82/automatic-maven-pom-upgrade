package lib

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

var (
	maven_command string
)

type Maven struct {
	log     *bufio.Writer
	command string
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile | log.Ltime)

}

func NewMaven(logfile *os.File) (m *Maven, err error) {
	m = &Maven{
		log: bufio.NewWriter(logfile),
	}
	m.command, err = m.determineCommand()
	return m, err
}

func (m *Maven) determineCommand() (cmd string, err error) {
	if _, err := os.Stat("mvnw"); err == nil {
		log.Print("maven wrapper script found")
		cmd = "./mvnw"
	} else {
		log.Print("no maven wrapper script found, try mvn from PATH")
		cmd = "mvn"
	}

	command := exec.Command(cmd, []string{"--version"}...)
	stdout, _ := command.StdoutPipe()
	stderr, _ := command.StderrPipe()
	go io.Copy(m.log, stdout)
	go io.Copy(m.log, stderr)

	defer m.log.Flush()
	err = command.Run()

	if err != nil {
		return "", err
	}
	return cmd, err
}

func Version() {
	cmd := maven_command
	args := []string{"--version"}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
