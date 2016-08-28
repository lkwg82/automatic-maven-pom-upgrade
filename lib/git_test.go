package lib

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"github.com/alexcesaro/log/golog"
	"github.com/alexcesaro/log"
)

var (
	logger golog.Logger
	execCmd func(string,...string) error
)

func init() {
	logger = *golog.New(os.Stderr, log.Debug)

	exec := &Exec{
		logger :logger,
	}
	execCmd = exec.execCommand
}

func TestDetectionGitNotInstalled(t *testing.T) {
	setup()
	defer cleanup()

	path := os.Getenv("PATH")
	defer os.Setenv("PATH", path)
	os.Setenv("PATH", ".")

	g := initGit()

	assert.False(t, g.IsInstalled())
}

func TestDetectionOfMissingGitDirectory(t *testing.T) {
	setup()
	defer cleanup()

	g := initGit()

	assert.False(t, g.HasRepo())
}

func TestDetectionOfNonDirtyGitRepository(t *testing.T) {
	setup()
	defer cleanup()

	execCmd("git", "init")

	g := initGit()

	assert.True(t, g.IsDirty())
}

func TestDetectionOfDirtyGitRepository(t *testing.T) {
	setup()
	defer cleanup()

	execCmd("git", "init")

	_, err := os.Create("test")
	if err != nil {
		assert.Error(t, err)
	}

	g := initGit()

	assert.True(t, g.IsDirty())
}

func TestGit_BranchExists(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()

	git := initGit()

	execCmd("git", "checkout", "-b", "test")

	assert.True(t, git.BranchExists("test"))
}

func TestGit_BranchCheckoutNew(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()

	git := initGit()

	git.BranchCheckoutNew("test")

	assert.True(t, git.BranchExists("test"), "missing branch test")
	assert.Equal(t, git.BranchCurrent(), "test")
}

func initGit() *Git {
	return NewGit(logger)
}

func createRepoWithSingleCommit() {
	execCmd("git", "init")

	execCmd("git", "config", "user.email", "test@ci.com")
	execCmd("git", "config", "user.name", "test")

	os.Create("test")
	execCmd("git", "add", "test")
	execCmd("git", "commit", "-m", "'test'", "test")
}

func TestGit_BranchCheckoutExisting(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()
	execCmd("git", "checkout", "-b", "test")
	execCmd("git", "checkout", "master")

	git := initGit()

	git.BranchCheckoutExisting("test")

	assert.True(t, git.BranchExists("test"), "missing branch test")
	assert.Equal(t, git.BranchCurrent(), "test")
}
