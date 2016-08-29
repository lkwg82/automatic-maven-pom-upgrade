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

func (m *Maven) DetermineCommand() error {
	m.logger.Info("determine command")
	var cmd string
	if _, err := os.Stat("mvnw"); err == nil {
		m.logger.Info("maven wrapper script found")
		cmd = "./mvnw"

		err = m.ExecCommand(cmd, "--version")
		if err != nil {
			return NewWrapError(err, "./mvnw --version")
		}
	} else {
		m.logger.Info("no maven wrapper script found, try mvn from PATH")
		cmd = "mvn"
		_, err := exec.LookPath("mvn")
		if err != nil {
			return NewWrapError(err, "missing mvn command")
		}
		m.command = "mvn"
	}

	m.command = cmd
	return nil
}

func (m *Maven) UpdateParent() (bool, string, error) {
	m.logger.Info("updating parent")
	args := []string{m.plugin + ":update-parent", "-DgenerateBackupPoms=false", "--batch-mode"}
	command := m.Command(m.command, args...)

	output, err := command.CombinedOutput()

	if err != nil {
		n := len(output)
		m.logger.Error("something failed: %s\n %s", err, string(output[:n]))
		panic("something went wrong")
	}

	n := len(output)
	content := string(output[:n])
	lines := strings.Split(content, "\n")

	updateToken := "[INFO] Updating parent from "
	noupdateToken := "[INFO] Current version of "

	if m.logger.LogDebug() {
		for _, line := range strings.Split(content, "\n") {
			m.logger.Debug(line)
		}
	}

	for _, line := range lines {
		if strings.HasPrefix(line, updateToken) {
			message := line[7:]
			m.logger.Infof("updated: %s", message)
			return true, message, err
		} else if strings.HasPrefix(line, noupdateToken) {
			message := line[7:]
			m.logger.Infof("no update: %s", message)
			return false, message, err
		}
	}

	panic("something went wrgon : " + content)
}
