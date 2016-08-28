package lib

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"github.com/alexcesaro/log/golog"
	"github.com/alexcesaro/log"
)

func TestDetectionGitNotInstalled(t *testing.T) {
	setup()
	defer cleanup()

	path := os.Getenv("PATH")
	defer os.Setenv("PATH", path)
	os.Setenv("PATH", ".")

	g := initGit()

	assert.False(t, g.IsInstalled())
}

func initGit() *Git {
	logger :=   *golog.New(os.Stderr, log.Debug)
	g := NewGit(logger)

	return g
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

	execCmd("git", []string{"init"})

	g := initGit()

	assert.True(t, g.IsDirty())
}

func TestDetectionOfDirtyGitRepository(t *testing.T) {
	setup()
	defer cleanup()

	execCmd("git", []string{"init"})

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
	os.Remove("git.log")

	execCmd("git", []string{"checkout", "-b", "test"})

	assert.True(t, git.BranchExists("test"))
}

func TestGit_BranchCheckoutNew(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()

	git := initGit()
	os.Remove("git.log")

	git.BranchCheckoutNew("test")

	assert.True(t, git.BranchExists("test"), "missing branch test")
	assert.Equal(t, git.BranchCurrent(), "test")
}

func createRepoWithSingleCommit() {
	execCmd("git", []string{"init"})

	execCmd("git", []string{"config", "user.email", "test@ci.com"})
	execCmd("git", []string{"config", "user.name", "test"})

	os.Create("test")
	execCmd("git", []string{"add", "test"})
	execCmd("git", []string{"commit", "-m", "'test'", "test"})
}

func TestGit_BranchCheckoutExisting(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()
	execCmd("git", []string{"checkout", "-b", "test"})
	execCmd("git", []string{"checkout", "master"})

	git := initGit()
	os.Remove("git.log")

	git.BranchCheckoutExisting("test")

	assert.True(t, git.BranchExists("test"), "missing branch test")
	assert.Equal(t, git.BranchCurrent(), "test")
}
