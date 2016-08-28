package lib

import (
	"log"
	"os/exec"
	"strings"
	"github.com/alexcesaro/log/golog"
)

type Git struct {
	Exec
	command       string
	CommitMessage string
}

func NewGit(logger golog.Logger) *Git {
	git := &Git{
		command: "git",
	}
	git.Logger(logger)
	return git
}

func (g *Git) HasRepo() bool {
	err := g.execCommand(g.command, "status")
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
	g.execCommand2(g.command + " add pom.xml")
	args := []string{"commit", "-m", g.CommitMessage, "pom.xml"}
	err := g.execCommand(g.command, args...)
	if err != nil {
		log.Panic(err)
	}
}

func (g *Git) exec(arguments string) {
	err := g.execCommand2(g.command + " " + arguments)
	if err != nil {
		log.Panic(err)
	}
}
