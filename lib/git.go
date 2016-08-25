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

func NewGit(logfile *os.File) *Git {
	return &Git{
		log:     bufio.NewWriter(logfile),
		logFile: logfile,
		command: "git",
	}
}

func (g *Git) HasRepo() bool {
	err := execCommand(g.log, g.command, "status")
	return err == nil
}

func (g *Git) IsInstalled() bool {
	_, err := exec.LookPath("git")
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

	// empty line has at least one line
	return len(lines) > 0
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
	err := execCommand2(g.log, g.command + " " + arguments)
	if err != nil {
		log.Panic(err)
	}
}
