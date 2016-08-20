package lib

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

const plugin_version = "2.3"

var (
	maven_command string
	plugin = fmt.Sprintf("org.codehaus.mojo:versions-maven-plugin:%s", plugin_version)
)

type Maven struct {
	log     *bufio.Writer
	logFile *os.File
	command string
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile | log.Ltime)

}

func NewMaven(logfile *os.File) (m *Maven, err error) {
	m = &Maven{
		log: bufio.NewWriter(logfile),
		logFile:logfile,
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

	err = execCommand(m.log, cmd, []string{"--version"}...)
	if err != nil {
		return "", err
	}
	return cmd, err
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

func (m *Maven) UpdateParent() {
	log.Print("updating parent")
	err := execCommand(m.log, m.command, []string{plugin +":update-parent", "-DgenerateBackupPoms=false"}...)
	if err != nil {
		content, _ := readFile(m.logFile.Name())
		log.Print(content)
	}

}
