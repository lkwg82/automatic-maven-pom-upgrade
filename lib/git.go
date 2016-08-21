package lib

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Git struct {
	log           *bufio.Writer
	logFile       *os.File
	command       string
	CommitMessage string
}

func NewGit(logfile *os.File) (g *Git, err error) {
	g = &Git{
		log:     bufio.NewWriter(logfile),
		logFile: logfile,
		command: "git",
	}

	err = g.checkInstalled()
	if err != nil {
		return g, NewWrapError(err, "git not installed")
	}

	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return g, NewWrapError2("missing git repository")
	}

	status, err := g.IsDirty()
	if err != nil {
		return g, NewWrapError(err, "unexpected error")
	}
	if !status {
		return g, NewWrapError2("repository is dirty")
	}
	return g, err
}

func (g *Git) checkInstalled() error {
	return execCommand(g.log, g.command, []string{"--version"}...)
}

func (g *Git) BranchExists(branch string) bool {
	args := []string{"branch", "--list", branch}
	output, err := exec.Command(g.command, args...).Output()

	if err != nil {
		log.Panic(err)
	}

	n := len(output)
	lines := strings.Split(string(output[:n]), "\n")
	return lines[0] == "* "+branch
}

func (g *Git) IsDirty() (bool, error) {
	args := []string{"status", "--porcelain"}
	output, err := exec.Command(g.command, args...).Output()

	if err != nil {
		return false, err
	}

	n := len(output)
	lines := strings.Split(string(output[:n]), "\n")
	// empty line has at least one
	return len(lines) == 1, err
}

func (g *Git) BranchCurrent() string {
	args := []string{"symbolic-ref", "--short", "HEAD"}
	output, err := exec.Command(g.command, args...).Output()
	if err != nil {
		log.Panic(err)
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

func (g *Git) Commit() {
	args := []string{"commit", "-m", g.CommitMessage, "pom.xml"}
	err := execCommand(g.log, g.command, args...)
	if err != nil {
		log.Panic(err)
	}
}

func (g *Git) exec(arguments string) {
	parts := strings.Split(arguments, " ")
	validArgs := make([]string, 0)
	for _, p := range parts {
		if p != "" {
			validArgs = append(validArgs, p)
		}
	}
	err := execCommand(g.log, g.command, validArgs...)
	if err != nil {
		log.Panic(err)
	}
}
