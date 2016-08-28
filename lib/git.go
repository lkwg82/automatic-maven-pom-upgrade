package lib

import (
	"github.com/alexcesaro/log/golog"
	"os"
	"os/exec"
	"strings"
)

type Git struct {
	Exec
	command string
}

func NewGit(logger golog.Logger) *Git {
	git := &Git{
		command: "git",
	}
	git.Logger(logger)
	return git
}

func (g *Git) HasRepo() bool {
	err := g.Command(g.command, "status").Run()
	return err == nil
}

func (g *Git) IsInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func (g *Git) BranchExists(branch string) bool {
	output, err := g.Command(g.command, "branch", "--list", branch).Output()

	n := len(output)
	content := strings.TrimSpace(string(output[:n]))
	lines := strings.Split(content, "\n")
	if g.logger.LogDebug() {
		for _, line := range lines {
			g.logger.Debug("output " + line)
		}
	}

	if err != nil {
		g.logger.Emergencyf("checking of branch '%s' exists: %s \n %s", string(output[:n]), err)
		os.Exit(1)
	}

	return lines[0] == "* " + branch || lines[0] == branch
}

func (g *Git) IsDirty() bool {
	output, err := g.Command(g.command, "status", "--porcelain").Output()

	if err != nil {
		panic(err)
	}

	n := len(output)
	content := strings.TrimSpace(string(output[:n]))
	if content == "" {
		return false
	}

	lines := strings.Split(content, "\n")

	g.logger.Debugf("output %d", len(lines))
	g.logger.Debugf("output '%s'", string(output[:n]))

	// empty line has at least one line
	return len(lines) > 0
}

func (g *Git) BranchCurrent() string {
	output, err := g.Command(g.command, "symbolic-ref", "--short", "HEAD").Output()
	if err != nil {
		g.logger.Emergency(err)
		os.Exit(1)
	}

	n := len(output)
	lines := strings.Split(string(output[:n]), "\n")
	return strings.Replace(lines[0], "refs/heads/", "", 0)
}

func (g *Git) BranchCheckoutExisting(branch string) {
	g.exec("checkout " + branch)
}

func (g *Git) BranchCheckoutNew(branch string) {
	g.exec("checkout -b " + branch)
}

func (g *Git) Commit(message string) {
	g.execCommand2(g.command + " add pom.xml")
	args := []string{"commit", "-m", "'" + message + "'", "pom.xml"}
	err := g.ExecCommand(g.command, args...)
	if err != nil {
		g.logger.Emergency(err)
		os.Exit(1)
	}
}

func (g *Git) exec(arguments string) {
	err := g.execCommand2(g.command + " " + arguments)
	if err != nil {
		g.logger.Emergency(err)
		os.Exit(1)
	}
}
