package lib

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDetectionGitNotInstalled(t *testing.T) {
	setup()
	defer cleanup()

	path := os.Getenv("PATH")
	defer os.Setenv("PATH", path)
	os.Setenv("PATH", ".")

	file, _ := os.Create("git.log")
	g := NewGit(file)

	assert.False(t, g.IsInstalled())
}

func TestDetectionOfMissingGitDirectory(t *testing.T) {
	setup()
	defer cleanup()

	file, _ := os.Create("git.log")
	g := NewGit(file)

	assert.False(t, g.HasRepo())
}

func TestDetectionOfNonDirtyGitRepository(t *testing.T) {
	setup()
	defer cleanup()

	execCmd("git", []string{"init"})

	file, _ := os.Create("git.log")
	// need to removed
	os.Remove("git.log")

	g := NewGit(file)

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

	file, _ := os.Create("git.log")
	g := NewGit(file)

	assert.True(t, g.IsDirty())
}

func TestGit_BranchExists(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()

	file, _ := os.Create("git.log")
	os.Remove("git.log")
	git := NewGit(file)

	execCmd("git", []string{"checkout", "-b", "test"})

	assert.True(t, git.BranchExists("test"))
}

func TestGit_BranchCheckoutNew(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()

	file, _ := os.Create("git.log")
	os.Remove("git.log")
	git := NewGit(file)

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

	file, _ := os.Create("git.log")
	os.Remove("git.log")
	git := NewGit(file)

	git.BranchCheckoutExisting("test")

	assert.True(t, git.BranchExists("test"), "missing branch test")
	assert.Equal(t, git.BranchCurrent(), "test")
}
