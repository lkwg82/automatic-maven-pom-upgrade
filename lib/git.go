package lib

import (
	"github.com/alexcesaro/log/golog"
	"github.com/droundy/goopt"
	"os"
	"os/exec"
	"strings"
)

var optGitNoUpdate = goopt.Flag([]string{"--git-no-update"}, nil, "skip automerge updates from master", "")
var optGitNoDirtyCheck = goopt.Flag([]string{"--git-no-dirty-check"}, nil, "skip dirty check", "")
var optNoCommit = goopt.Flag([]string{"--git-no-commit"}, nil, "skip commit", "")
var optHookAfterCommit = goopt.String([]string{"--hook-after"}, "/bin/echo", "command to call after commit (commit message is 1st arg)")

// Git wraps git command execution
type Git struct {
	Exec
}

// NewGit constructs new Git
func NewGit(logger golog.Logger) *Git {
	git := &Git{
		Exec: Exec{
			logger: logger,
			Cmd:    "git",
		},
	}
	return git
}

// BranchCheckoutExisting checks out already existing branch
func (g *Git) BranchCheckoutExisting(branch string) {
	g.CommandRunExitOnErr("checkout", branch)
}

// BranchCheckoutNew creates and checks out new branch
func (g *Git) BranchCheckoutNew(branch string) {
	g.CommandRunExitOnErr("checkout", "-b", branch)
}

// BranchCurrent returns current branch
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

// BranchExists detects if given branch exists
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

	return lines[0] == "* "+branch || lines[0] == branch
}

// Commit perform a git commit with a given message
func (g *Git) Commit(message string) {
	g.CommandRunExitOnErr("add", "pom.xml")
	args := []string{"commit", "-m", "'" + message + "'", "pom.xml"}
	g.CommandRunExitOnErr(args...)
}

// DominantMergeFrom merges given branch into current, overwrites any local conflicting changes
func (g *Git) DominantMergeFrom(branch, message string) {
	output, err := g.Command("merge", "--commit", "--strategy-option=theirs", branch, "--message", "\""+message+"\"").CombinedOutput()
	g.DebugStdoutErr(output, err)
	if err != nil {
		panic(err)
	}
}

// Fetch perform a git fetch
func (g *Git) Fetch() {
	g.CommandRunExitOnErr("fetch")
}

// HasMergeConflict detects if merging a given branch into current would result merge conflicts
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

// IsDirty detects if there are any uncommitted changes
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

// IsInstalled checks if git is installed
func (g *Git) IsInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// IsInSyncWith checks if the current branch does not need updates from given branch
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

// IsRepo checks if current directory is considered a git repository
func (g *Git) IsRepo() bool {
	err := g.Command("status").Run()
	g.DebugErr(err)
	return err == nil
}

// OptionalAutoMergeMaster does a optional merge from the master, it dependends on the cmd args
func (g *Git) OptionalAutoMergeMaster() {
	if !*optGitNoUpdate && !g.IsInSyncWith("master") {
		g.DominantMergeFrom("master", "merge updates from master")
	}
}

// CheckIsInstalled returns an error if git is not installed
func (g *Git) CheckIsInstalled() error {
	if !g.IsInstalled() {
		return newWrapError2("need git to be installed or in the PATH")
	}
	return nil
}

// CheckIsRepo returns an error if current directory is not a git repository
func (g *Git) CheckIsRepo() error {
	if !g.IsRepo() {
		return newWrapError2("need called from a directory, which has a repository")
	}
	return nil
}

// OptionalCheckIsDirty returns on optional error if git repository is dirty, it dependends on the cmd args
func (g *Git) OptionalCheckIsDirty() error {
	if !*optGitNoDirtyCheck && g.IsDirty() {
		return newWrapError2("repository is dirty, plz commit or reset")
	}
	return nil
}

// OptionalCommit does a optional commit, it dependends on the cmd args
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

func (g *Git) execAfterCommitHook(message string) {
	cmd := NewExec(g.logger, *optHookAfterCommit)
	cmd.CommandRunExitOnErr(message)
}
