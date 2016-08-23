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

func NewGit(logfile *os.File) (g *Git) {
	return &Git{
		log:     bufio.NewWriter(logfile),
		logFile: logfile,
		command: "git",
	}
}

func (g *Git) HasRepo() bool {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return false
	}
	return true
}

func (g *Git) IsInstalled() bool {
	err := execCommand(g.log, g.command, []string{"--version"}...)
	return err == nil
}

func (g *Git) BranchExists(branch string) bool {
	args := []string{"branch", "--list", branch}
	output, err := exec.Command(g.command, args...).Output()

	if err != nil {
		log.Panic(err)
	}

	n := len(output)
	lines := strings.Split(string(output[:n]), "\n")
	return lines[0] == "* " + branch
}

func (g *Git) IsDirty() bool {
	args := []string{"status", "--porcelain"}
	output, err := exec.Command(g.command, args...).Output()

	if err != nil {
		panic(err)
	}

	n := len(output)
	lines := strings.Split(string(output[:n]), "\n")

	// empty line has at least one
	return len(lines) == 1
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
