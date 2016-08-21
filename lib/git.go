package lib

import (
	"bufio"
	"os"
	"strings"
	"os/exec"
	"errors"
	"log"
)

const PREFIX = "autoupdate"

type Git struct {
	log           *bufio.Writer
	logFile       *os.File
	command       string
	CommitMessage string
}

func NewGit(logfile *os.File) (g *Git, err error) {
	g = &Git{
		log: bufio.NewWriter(logfile),
		logFile:logfile,
		command: "git",
	}

	err = g.checkInstalled()
	if err != nil {
		return
	}

	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return g, errors.New("missing git repository")
	}

	if status, err := g.IsDirty(); err != nil || !status {
		if ( !status) {
			return g, errors.New("repository is dirty")
		}
		return g, err
	}

	return g, err
}

func (g *Git) checkInstalled() error {
	return execCommand(g.log, g.command, []string{"--version"}...)
}

func (g *Git) BranchExists(branch string) bool {
	args := []string{"branch", "--list", PREFIX + "_" + branch}
	output, err := exec.Command(g.command, args...).Output()

	if err != nil {
		log.Panic(err)
	}

	n := len(output)
	lines := strings.Split(string(output[:n]), "\n")
	return len(lines) == 1
}

func (g *Git) IsDirty() (bool, error) {
	args := []string{"status", "--procelain"}
	output, err := exec.Command(g.command, args...).Output()

	if err != nil {
		return false, err
	}

	n := len(output)
	lines := strings.Split(string(output[:n]), "\n")
	return len(lines) == 0, err
}

func (g *Git) BranchCheckout(branch string) {
	if g.BranchExists(branch) {
		g.exec("checkout " + branch)
	}
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
