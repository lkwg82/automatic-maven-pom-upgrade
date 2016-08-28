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

func TestDetectionOfMavenWrapper(t *testing.T) {
	setup()
	defer cleanup()

	content := "#!/bin/sh\necho -n x"
	ioutil.WriteFile("./mvnw", []byte(content), 0700)

	err := maven.DetermineCommand()

	assert.Nil(t, err)
	assert.Equal(t, maven.command, "./mvnw")
}

func TestMavenNotFound(t *testing.T) {
	setup()
	defer cleanup()

	os.Setenv("PATH", "")

	err := maven.DetermineCommand()

	assert.Error(t, err)
}

func TestMavenWrapperFound(t *testing.T) {
	setup()
	defer cleanup()

	os.Setenv("PATH", ".")
	content := "#!/bin/sh\necho -n x"
	ioutil.WriteFile("mvnw", []byte(content), 0700)

	// action
	err := maven.DetermineCommand()

	assert.Nil(t, err)
	assert.Equal(t, maven.command, "./mvnw")
}

func TestMavenParentPomUpdate(t *testing.T) {
	maven := setupWithTestProject(t, "simple-parent-update")
	defer cleanup()

	// action
	updateMessage, err := maven.UpdateParent()

	assert.Nil(t, err)
	assert.NotZero(t, updateMessage)
	assert.True(t, strings.HasPrefix(updateMessage, "Updating parent from 1.3.7.RELEASE to "), "but was : "+updateMessage)
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
