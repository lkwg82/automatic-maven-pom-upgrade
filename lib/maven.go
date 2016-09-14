package lib

import (
	"fmt"
	"github.com/alexcesaro/log/golog"
	"github.com/droundy/goopt"
	"os"
	"os/exec"
	"strings"
)

const (
	pluginName    = "org.codehaus.mojo:versions-maven-plugin"
	pluginVersion = "2.3"
)

var optMavenSettingsPath = goopt.String([]string{"--maven-settings"}, "", "path to maven settings (equivalent to -s)")

// Maven wraps running the external maven command
type Maven struct {
	Exec
	plugin       string
	settingsPath string
}

// NewMaven is the construtor for Maven
func NewMaven(logger golog.Logger) *Maven {
	maven := &Maven{
		Exec: Exec{
			logger: logger,
		},
		plugin: fmt.Sprintf("%s:%s", pluginName, pluginVersion),
	}
	return maven
}

// DetermineCommand decides whether to use a repo-local maven wrapper (mvnw) or a system global maven (mvn) command
func (m *Maven) DetermineCommand() error {
	m.logger.Info("determine command")
	var cmd string
	if _, err := os.Stat("mvnw"); err == nil {
		m.logger.Info("maven wrapper script found")
		cmd = "./mvnw"

		exec := NewExec(m.logger, cmd)
		err = exec.CommandRun("--version")
		if err != nil {
			return newWrapError(err, "./mvnw --version")
		}
	} else {
		m.logger.Info("no maven wrapper script found, try mvn from PATH")
		cmd = "mvn"
		_, err := exec.LookPath("mvn")
		if err != nil {
			return newWrapError(err, "missing mvn command")
		}
	}

	m.Exec.Cmd = cmd
	return nil
}

// SettingsPath sets a custom settings.xml path to use with maven
func (m *Maven) SettingsPath(path string) error {
	m.logger.Debugf("use maven settings path: %s", path)
	if path == "" {
		m.logger.Debug("ignoring empty settings path")
		return nil
	}

	file, err := os.Stat(path)
	if os.IsNotExist(err) {
		return newWrapError2("path '" + path + "' is not existing")
	}

	if file.IsDir() {
		return newWrapError2("path '" + path + "' is directory, expected a file")
	}
	m.settingsPath = path
	return nil
}

// ParseCommandline parses the command line arguments for some maven options
func (m *Maven) ParseCommandline() error {
	if err := m.SettingsPath(*optMavenSettingsPath); err != nil {
		return err
	}
	return nil
}

// UpdateParent tries to update the parent pom of the local pom.xml
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
		os.Exit(1)
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
