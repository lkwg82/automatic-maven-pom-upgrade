package lib

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const plugin_version = "2.3"

var (
	plugin = fmt.Sprintf("org.codehaus.mojo:versions-maven-plugin:%s", plugin_version)
)

type Maven struct {
	log           *bufio.Writer
	logFile       *os.File
	command       string
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

func (m *Maven) UpdateParent() (string, error) {
	log.Print("updating parent")
	err := execCommand(m.log, m.command, []string{plugin + ":update-parent", "-DgenerateBackupPoms=false"}...)

	if err != nil {
		log.Fatalf("something failed: %s", err)
		return "", err
	}
	content, err := readFile(m.logFile.Name())
	if err != nil {
		panic(err)
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		updateToken := "[INFO] Updating parent from "
		if strings.HasPrefix(line, updateToken) {
			return line[7:], err
		}
	}
	panic("missed the line with the message : " + content)
}
