package lib

import (
	"fmt"
	"github.com/alexcesaro/log/golog"
	"os"
	"os/exec"
	"strings"
)

const (
	plugin_name = "org.codehaus.mojo:versions-maven-plugin"
	plugin_version = "2.3"
)

type Maven struct {
	Exec
	command string
	plugin  string
}

func NewMaven(logger golog.Logger) *Maven {
	maven := &Maven{
		plugin: fmt.Sprintf("%s:%s", plugin_name, plugin_version),
	}
	maven.Logger(logger)
	return maven
}

func (m *Maven) DetermineCommand() (err error) {
	m.logger.Info("determine command")
	var cmd string
	if _, err := os.Stat("mvnw"); err == nil {
		m.logger.Info("maven wrapper script found")
		cmd = "./mvnw"

		err = m.execCommand(cmd, []string{"--version"}...)
		if err != nil {
			return NewWrapError(err, "./mvnw --version")
		}
	} else {
		m.logger.Info("no maven wrapper script found, try mvn from PATH")
		cmd = "mvn"
		_, err = exec.LookPath("mvn")
		if err != nil {
			return NewWrapError(err, "missing mvn command")
		}
		m.command = "mvn"
	}

	m.command = cmd
	return err
}

func (m *Maven) UpdateParent() (string, error) {
	m.logger.Info("updating parent")
	args := []string{m.plugin + ":update-parent", "-DgenerateBackupPoms=false", "--batch-mode"}
	m.logger.Debugf("executing: %s %s", m.command, strings.Join(args, " "))
	command := exec.Command(m.command, args...)

	output, err := command.CombinedOutput()

	if err != nil {
		n := len(output)
		m.logger.Error("something failed: %s\n %s", err, string(output[:n]))
		os.Exit(1)
		return "", err
	}

	n := len(output)
	content := string(output[:n])
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		updateToken := "[INFO] Updating parent from "
		if strings.HasPrefix(line, updateToken) {
			return line[7:], err
		}
	}
	panic("missed the line with the message : " + content)
}
