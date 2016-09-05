package lib

import (
	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/golog"
	"github.com/rafecolton/go-fileutils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

var maven *Maven

func init() {
	logger := *golog.New(os.Stderr, log.Warning)
	maven = NewMaven(logger)
}

func TestMaven_NotFound(t *testing.T) {
	setup()
	defer cleanup()

	os.Unsetenv("PATH")

	err := maven.DetermineCommand()

	assert.Error(t, err)
}

func TestMaven_SettingsPathIsMissing(t *testing.T) {
	defer cleanup()
	maven := setupWithTestProject(t, "simple-parent-update")

	err := maven.SettingsPath("x")

	assert.Error(t, err)
}

func TestMaven_UpdateParentPom(t *testing.T) {
	defer cleanup()
	maven := setupWithTestProject(t, "simple-parent-update")

	updated, updateMessage, err := maven.UpdateParent()

	assert.Nil(t, err)
	assert.NotZero(t, updateMessage)
	assert.True(t, updated)
	assert.True(t, strings.HasPrefix(updateMessage, "Updating parent from 1.3.7.RELEASE to "), "but was : " + updateMessage)
}

func TestMaven_UpdateParentPomTwice(t *testing.T) {
	defer cleanup()
	maven := setupWithTestProject(t, "simple-parent-update")

	updated, updateMessage, err := maven.UpdateParent()
	updated, updateMessage, err = maven.UpdateParent()

	assert.Nil(t, err)
	assert.False(t, updated)
	assert.NotEmpty(t, updateMessage)
}

func TestMaven_WrapperFound(t *testing.T) {
	setup()
	defer cleanup()

	os.Setenv("PATH", ".")
	content := "#!/bin/sh\necho -n x"
	ioutil.WriteFile("mvnw", []byte(content), 0700)

	// action
	err := maven.DetermineCommand()

	assert.Nil(t, err)
	assert.Equal(t, maven.Exec.Cmd, "./mvnw")
}

func setupWithTestProject(t *testing.T, testProjectName string) *Maven {
	setup()
	sourcePath := path.Dir(temporaryDirectoryForTests.Cwd + "/../test-projects/" + testProjectName)
	if err := fileutils.CpR(sourcePath, "x"); err != nil {
		panic(err)
	}
	if err := os.Chdir("x/" + testProjectName); err != nil {
		panic(err)
	}

	if _, exists := os.LookupEnv("JAVA_HOME"); !exists {
		t.Skip("missing JAVA_HOME env, try again this test in on the shell")
	}

	err := maven.DetermineCommand()
	if err != nil {
		panic(err)
	}

	return maven
}
