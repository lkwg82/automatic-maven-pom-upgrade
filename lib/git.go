package lib

import (
	"bufio"
	"os"
	"strings"
	"os/exec"
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
	internalName := PREFIX + "_" + branch
	args := []string{"branch", "--list", internalName}
	output, err := exec.Command(g.command, args...).Output()

	if err != nil {
		log.Panic(err)
	}

	n := len(output)
	content := string(output[:n])
	log.Println("out" + content)
	lines := strings.Split(content, "\n")
	return len(lines) == 1
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

func (g *Git) BranchCheckout(branch string) {
	if g.BranchExists(branch) {
		g.exec("checkout " + PREFIX + "_" + branch)
	}
	g.exec("checkout -b " + PREFIX + "_" + branch)
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
