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
	_, err := NewGit(file)

	assert.Equal(t, err.(*WrapError).string(), "git not installed")
}

func TestDetectionOfMissingGitDirectory(t *testing.T) {
	setup()
	defer cleanup()

	file, _ := os.Create("git.log")
	_, err := NewGit(file)

	assert.Equal(t, err.(*WrapError).string(), "missing git repository")
}

func TestDetectionOfNonDirtyGitRepository(t *testing.T) {
	setup()
	defer cleanup()

	execCmd("git", []string{"init"})

	file, _ := os.Create("git.log")
	// need to removed
	os.Remove("git.log")

	_, err := NewGit(file)

	assert.Nil(t, err)
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
	_, err = NewGit(file)

	assert.Equal(t, err.(*WrapError).string(), "repository is dirty")
}

func TestGit_BranchExists(t *testing.T) {
	setup()
	defer cleanup()

	execCmd("git", []string{"init"})

	file, _ := os.Create("git.log")
	os.Remove("git.log")
	git, err := NewGit(file)

	execCmd("git", []string{"checkout", "-b", "test"})

	assert.Nil(t, err)
	assert.True(t, git.BranchExists("test"))
}

//func TestGit_BranchCheckoutNew(t *testing.T) {
//	setup()
//	defer cleanup()
//
//	execCmd("git", []string{"init"})
//	os.Create("test")
//	execCmd("git", []string{"add", "test"})
//	execCmd("git", []string{"commit", "-m", "'test'", "test"})
//
//	file, _ := os.Create("git.log")
//	os.Remove("git.log")
//	git, err := NewGit(file)
//
//	git.BranchCheckout("test")
//
//	assert.Nil(t, err)
//	assert.True(t, git.BranchExists("test"))
//	assert.Equal(t, git.BranchCurrent(), "test")
//}

//func TestGit_BranchCheckoutExisting(t *testing.T) {
//	setup()
//	defer cleanup()
//
//	execCmd("git", []string{"init"})
//
//	file, _ := os.Create("git.log")
//	os.Remove("git.log")
//	git, err := NewGit(file)
//
//	execCmd("git", []string{"checkout", "-b", "test"})
//
//	assert.Nil(t, err)
//	assert.True(t, git.BranchExists("test"))
//}