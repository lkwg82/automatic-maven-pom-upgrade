package lib

import (
	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/golog"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	logger  golog.Logger
	git     *Git
	execGit func(...string) error
)

func init() {
	logger = *golog.New(os.Stderr, log.Warning)

	exec := &Exec{
		logger: logger,
		Cmd:    "git",
	}
	execGit = func(args ...string) error {
		err := exec.CommandRun(args...)
		if err == nil {
			return nil
		}
		panic(err)
	}
	git = NewGit(logger)
}

func TestGit_IsInstalled(t *testing.T) {
	setup()
	defer cleanup()

	path := os.Getenv("PATH")
	defer os.Setenv("PATH", path)
	os.Setenv("PATH", ".")

	assert.False(t, git.IsInstalled())
}

func TestGit_IsRepositoryWhenMissing(t *testing.T) {
	setup()
	defer cleanup()

	assert.False(t, git.IsRepo())
}

func TestGit_IsRepositoryFromSubdirectory(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()

	if err := os.MkdirAll("x/s/a", 0744); err != nil {
		panic(err)
	}
	if err := os.Chdir("x/s/a"); err != nil {
		panic(err)
	}

	assert.True(t, git.IsRepo())
}

func TestGit_IsNotDirtyGitRepository(t *testing.T) {
	setup()
	defer cleanup()

	execGit("init")

	assert.False(t, git.IsDirty())
}

func TestGit_IsDirty(t *testing.T) {
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

func TestGit_BranchExistsRmote(t *testing.T) {
	setup()
	defer cleanup()

	// remote repository
	os.Mkdir("remote", 0755)
	os.Chdir("remote")

	execGit("init")
	execGit("config", "user.email", "test@ci.com")
	execGit("config", "user.name", "test")
	os.Create("test")
	execGit("add", "test")
	execGit("commit", "-m", "'test'", "test")
	execGit("checkout", "-b", "test")
	parent, _ := os.Getwd()

	// local
	os.Chdir("..")
	execGit("clone", parent, "local")
	os.Chdir("local")
	git.Fetch()

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

func TestGit_IsInSyncWithMasterSame(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()

	execGit("checkout", "-b", "test")

	assert.True(t, git.IsInSyncWith("master"))
}

func TestGit_IsInSyncWithMasterAhead(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()

	execGit("checkout", "-b", "test")
	os.Create("test2")
	execGit("add", "test2")
	execGit("commit", "-m", "'test'", "test2")

	assert.True(t, git.IsInSyncWith("master"))
}

func TestGit_IsNotInSyncWithMaster(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()

	os.Create("test2")
	execGit("add", "test2")
	execGit("commit", "-m", "'test'", "test2")

	execGit("checkout", "-b", "test")
	execGit("reset", "--hard", "HEAD~1")

	assert.False(t, git.IsInSyncWith("master"))
	assert.False(t, git.IsDirty())
}

func TestGit_MergeMasterIntoBranch(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()

	os.Create("test2")
	execGit("add", "test2")
	execGit("commit", "-m", "'test'", "test2")

	execGit("checkout", "-b", "test")
	execGit("reset", "--hard", "HEAD~1")

	os.Create("test2")
	ioutil.WriteFile("test2", []byte("test"), 0700)
	execGit("add", "test2")
	execGit("commit", "-m", "'test in branch test2'", "test2")

	assert.False(t, git.IsInSyncWith("master"))
	assert.False(t, git.HasMergeConflict("master"))

	// action
	git.DominantMergeFrom("master", "updates from master")

	assert.True(t, git.IsInSyncWith("master"))
}

func TestGit_MergeMasterIntoBranchWithConflict(t *testing.T) {
	setup()
	defer cleanup()

	createRepoWithSingleCommit()

	// write to test in branch test
	execGit("checkout", "-b", "test")
	ioutil.WriteFile("test", []byte("test"), 0700)
	execGit("add", "test")
	execGit("commit", "-m", "'update test'", "test")

	// write to test in master different
	execGit("checkout", "master")
	ioutil.WriteFile("test", []byte("afafd"), 0700)
	execGit("add", "test")
	execGit("commit", "-m", "'update test'", "test")

	execGit("checkout", "test")

	assert.True(t, git.HasMergeConflict("master"))

	// action
	git.DominantMergeFrom("master", "updates from master")

	time.Sleep(10 * time.Millisecond)

	assert.True(t, git.IsInSyncWith("master"))
	output, _ := ioutil.ReadFile("test")
	n := len(output)
	content := strings.TrimSpace(string(output[:n]))
	assert.Equal(t, "afafd", content)
}

func createRepoWithSingleCommit() {
	execGit("init")

	execGit("config", "user.email", "test@ci.com")
	execGit("config", "user.name", "test")

	os.Create("test")
	execGit("add", "test")
	execGit("commit", "-m", "'test'", "test")
}
