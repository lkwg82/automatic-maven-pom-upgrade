package lib

import (
	"log"
	"os/exec"
	"strings"
	"github.com/alexcesaro/log/golog"
	"os"
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
	output, err := exec.Command(g.command, "branch", "--list", branch).Output()

	if err != nil {
		log.Panic(err)
	}

	n := len(output)
	lines := strings.Split(string(output[:n]), "\n")
	return lines[0] == "* " + branch
}

func (g *Git) IsDirty() bool {
	output, err := exec.Command(g.command, "status", "--porcelain").Output()

	if err != nil {
		panic(err)
	}

	n := len(output)
	lines := strings.Split(string(output[:n]), "\n")

	// empty line has at least one line
	return len(lines) > 0
}

func (g *Git) BranchCurrent() string {
	output, err := exec.Command(g.command, "symbolic-ref", "--short", "HEAD").Output()
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
		g.logger.Emergency(err)
		os.Exit(1)
	}
}
