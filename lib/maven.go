package lib

import (
	"fmt"
	"github.com/alexcesaro/log/golog"
	"os"
	"os/exec"
	"strings"
	"github.com/droundy/goopt"
)

const (
	plugin_name = "org.codehaus.mojo:versions-maven-plugin"
	plugin_version = "2.3"
)

var optMavenSettingsPath = goopt.String([]string{"--maven-settings"}, "", "path to maven settings (equivalent to -s)")

type Maven struct {
	Exec
	plugin       string
	settingsPath string
}

func NewMaven(logger golog.Logger) *Maven {
	maven := &Maven{
		Exec:Exec{
			logger:logger,
		},
		plugin: fmt.Sprintf("%s:%s", plugin_name, plugin_version),
	}
	return maven
}

func (m *Maven) DetermineCommand() error {
	m.logger.Info("determine command")
	var cmd string
	if _, err := os.Stat("mvnw"); err == nil {
		m.logger.Info("maven wrapper script found")
		cmd = "./mvnw"

		exec := NewExec(m.logger, cmd)
		err = exec.CommandRun("--version")
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
	}

	m.Exec.Cmd = cmd
	return nil
}

func (m *Maven) SettingsPath(path string) error {
	m.logger.Debugf("use maven settings path: %s", path)
	if path == "" {
		m.logger.Debug("ignoring empty settings path")
		return nil
	}

	file, err := os.Stat(path)
	if os.IsNotExist(err) {
		return NewWrapError2("path '" + path + "' is not existing")
	}

	if file.IsDir() {
		return NewWrapError2("path '" + path + "' is directory, expected a file")
	}
	m.settingsPath = path
	return nil
}

func (m *Maven) ParseCommandline() error {
	if err := m.SettingsPath(*optMavenSettingsPath); err != nil {
		return err
	}
	return nil
}

func (m *Maven) UpdateParent() (bool, string, error) {
	m.logger.Info("updating parent")
	args := []string{m.plugin + ":update-parent", "-DgenerateBackupPoms=false", "--batch-mode"}
	if m.settingsPath != "" {
		args = append(args, "-s", m.settingsPath)
	}
	output, err := m.Command(args...).CombinedOutput()
	m.DebugStdoutErr(output, err)

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

	panic("something went wrong : " + content)
}
