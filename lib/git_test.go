package lib

import (
	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/golog"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	logger golog.Logger
	git     *Git
	execGit func(...string) error
)

func init() {
	logger = *golog.New(os.Stderr, log.Debug)

	exec := &Exec{
		logger: logger,
	}
	execGit = func(args ...string) error {
		err := exec.execCommand("git", args...)
		if err == nil {
			return nil
		}
		panic(err)
	}
	git = NewGit(logger)
}

func TestDetectionGitNotInstalled(t *testing.T) {
	setup()
	defer cleanup()

	path := os.Getenv("PATH")
	defer os.Setenv("PATH", path)
	os.Setenv("PATH", ".")

	assert.False(t, git.IsInstalled())
}

func TestDetectionOfMissingGitRepository(t *testing.T) {
	setup()
	defer cleanup()

	assert.False(t, git.HasRepo())
}

func TestDetectionOfGitRepositoryFromSubdirectory(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()

	if err := os.MkdirAll("x/s/a", 0744); err != nil {
		panic(err)
	}
	if err := os.Chdir("x/s/a"); err != nil {
		panic(err)
	}

	assert.True(t, git.HasRepo())
}

func TestDetectionOfNonDirtyGitRepository(t *testing.T) {
	setup()
	defer cleanup()

	execGit("init")

	assert.False(t, git.IsDirty())
}

func TestDetectionOfDirtyGitRepository(t *testing.T) {
	setup()
	defer cleanup()

	execGit("init")

	_, err := os.Create("test")
	if err != nil {
		assert.Error(t, err)
	}

	assert.True(t, git.IsDirty())
}

func TestGit_BranchExists(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()

	execGit("checkout", "-b", "test")

	assert.True(t, git.BranchExists("test"))
}

func TestGit_BranchExistsNotOnBranch(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()

	execGit("checkout", "-b", "test")
	execGit("checkout", "master")

	assert.True(t, git.BranchExists("test"))
}

func TestGit_BranchCheckoutNew(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()

	git.BranchCheckoutNew("test")

	assert.True(t, git.BranchExists("test"), "missing branch test")
	assert.Equal(t, git.BranchCurrent(), "test")
}

func TestGit_BranchCheckoutExisting(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()
	execGit("checkout", "-b", "test")
	execGit("checkout", "master")

	git.BranchCheckoutExisting("test")

	assert.True(t, git.BranchExists("test"), "missing branch test")
	assert.Equal(t, git.BranchCurrent(), "test")
}

func TestGit_Commit(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()

	_, err := os.Create("pom.xml")
	if err != nil {
		panic(err)
	}

	git.Commit("initial update")

	assert.False(t, git.IsDirty())
}

func createRepoWithSingleCommit() {
	execGit("init")

	execGit("config", "user.email", "test@ci.com")
	execGit("config", "user.name", "test")

	os.Create("test")
	execGit("add", "test")
	execGit("commit", "-m", "'test'", "test")
}
