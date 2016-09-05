package lib

import (
	"github.com/alexcesaro/log/golog"
	"os"
	"os/exec"
	"strings"
	"github.com/droundy/goopt"
)

var optGitNoUpdate = goopt.Flag([]string{"--git-no-update"}, nil, "skip automerge updates from master", "")
var optGitNoDirtyCheck = goopt.Flag([]string{"--git-no-dirty-check"}, nil, "skip dirty check", "")
var optNoCommit = goopt.Flag([]string{"--git-no-commit"}, nil, "skip commit", "")
var optHookAfterCommit = goopt.String([]string{"--hook-after"}, "/bin/echo", "command to call after commit (commit message is 1st arg)")

type Git struct {
	Exec
}

func NewGit(logger golog.Logger) *Git {
	git := &Git{
		Exec : Exec{
			logger:logger,
			Cmd:"git",
		},
	}
	return git
}

func (g *Git) BranchCheckoutExisting(branch string) {
	g.CommandRunExitOnErr("checkout", branch)
}

func (g *Git) BranchCheckoutNew(branch string) {
	g.CommandRunExitOnErr("checkout", "-b", branch)
}

func (g *Git) BranchCurrent() string {
	output, err := g.Command("symbolic-ref", "--short", "HEAD").Output()
	g.DebugStdoutErr(output, err)
	if err != nil {
		g.logger.Emergency(err)
		os.Exit(1)
	}

	n := len(output)
	lines := strings.Split(string(output[:n]), "\n")
	return strings.Replace(lines[0], "refs/heads/", "", 0)
}

func (g *Git) BranchExists(branch string) bool {
	output, err := g.Command("branch", "--list", branch).Output()
	g.DebugStdoutErr(output, err)

	n := len(output)
	content := strings.TrimSpace(string(output[:n]))
	lines := strings.Split(content, "\n")

	if err != nil {
		g.logger.Emergencyf("checking of branch '%s' exists: %s \n %s", string(output[:n]), err)
		os.Exit(1)
	}

	return lines[0] == "* " + branch || lines[0] == branch
}

func (g *Git) Commit(message string) {
	g.CommandRunExitOnErr("add", "pom.xml")
	args := []string{"commit", "-m", "'" + message + "'", "pom.xml"}
	g.CommandRunExitOnErr(args...)
}

func (g *Git) DominantMergeFrom(branch, message string) {
	output, err := g.Command("merge", "--commit", "--strategy-option=theirs", branch, "--message", "\"" + message + "\"").CombinedOutput()
	g.DebugStdoutErr(output, err)
	if err != nil {
		panic(err)
	}
}

func (g *Git) Fetch() {
	g.CommandRunExitOnErr("fetch")
}

func (g *Git) HasMergeConflict(branch string) (result bool) {
	output, err := g.Command("merge", "--no-commit", branch).CombinedOutput()
	g.DebugStdoutErr(output, err)

	n := len(output)
	content := strings.TrimSpace(string(output[:n]))
	if content == "Already up-to-date." {
		result = false
	} else {
		result = err != nil
		output, err = g.Command("merge", "--abort").CombinedOutput()
		g.DebugStdoutErr(output, err)
	}
	return result
}

func (g *Git) IsDirty() bool {
	output, err := g.Command("status", "--porcelain").Output()
	g.DebugStdoutErr(output, err)

	if err != nil {
		panic(err)
	}

	n := len(output)
	content := strings.TrimSpace(string(output[:n]))
	if content == "" {
		return false
	}

	lines := strings.Split(content, "\n")

	// empty line has at least one line
	return len(lines) > 0
}

func (g *Git) IsInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func (g *Git) IsInSyncWith(branch string) (result bool) {
	output, err := g.Command("merge", "--no-commit", branch).CombinedOutput()
	g.DebugStdoutErr(output, err)

	n := len(output)
	content := strings.TrimSpace(string(output[:n]))
	if content == "Already up-to-date." {
		result = true
	} else {
		result = false
		output, err = g.Command("merge", "--abort").CombinedOutput()
		g.DebugStdoutErr(output, err)
	}
	return result
}

func (g *Git) IsRepo() bool {
	err := g.Command("status").Run()
	g.DebugErr(err)
	return err == nil
}

func (g *Git) OptionalAutoMergeMaster() {
	if !*optGitNoUpdate &&  !g.IsInSyncWith("master") {
		g.DominantMergeFrom("master", "merge updates from master")
	}
}

func (g *Git) CheckIsInstalled() error {
	if !g.IsInstalled() {
		return NewWrapError2("need git to be installed or in the PATH")
	}
	return nil
}

func (g *Git) CheckIsRepo() error {
	if !g.IsRepo() {
		return NewWrapError2("need called from a directory, which has a repository")
	}
	return nil
}

func (g *Git) OptionalCheckIsDirty() error {
	if !*optGitNoDirtyCheck &&                !g.IsDirty() {
		return NewWrapError2("repository is dirty, plz commit or reset")
	}
	return nil
}

func (g *Git) OptionalCommit(message string, echo func(string, ...string)) {
	if !*optNoCommit {
		echo("committing '%s'", message)
		g.Commit(message)
		echo("executing afterCommitHook")
		g.execAfterCommitHook(message)
	} else {
		echo("skipping commit")
	}
}

func (g *Git)execAfterCommitHook(message string) {
	cmd := NewExec(g.logger, *optHookAfterCommit)
	cmd.CommandRunExitOnErr(message)
}
