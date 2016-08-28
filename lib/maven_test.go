package lib

import (
	"github.com/rafecolton/go-fileutils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"github.com/alexcesaro/log/golog"
	"github.com/alexcesaro/log"
)

func TestDetectionOfMavenWrapper(t *testing.T) {
	setup()
	defer cleanup()

	content := "#!/bin/sh\necho -n x"
	ioutil.WriteFile("./mvnw", []byte(content), 0700)

	maven := initMaven()
	err := maven.DetermineCommand()

	assert.Nil(t, err)
	assert.Equal(t, maven.command, "./mvnw")
}

func TestMavenNotFound(t *testing.T) {
	setup()
	defer cleanup()

	os.Setenv("PATH", "")

	maven := initMaven()
	err := maven.DetermineCommand()

	assert.Error(t, err)
}

func TestMavenWrapperFound(t *testing.T) {
	setup()
	defer cleanup()

	os.Setenv("PATH", ".")
	content := "#!/bin/sh\necho -n x"
	ioutil.WriteFile("mvnw", []byte(content), 0700)

	maven := initMaven()

	// action
	err := maven.DetermineCommand()

	assert.Nil(t, err)
	assert.Equal(t, maven.command, "./mvnw")
}

func TestMavenParentPomUpdate(t *testing.T) {
	maven := setupWithTestProject("simple-parent-update")
	defer cleanup()

	// action
	updateMessage, err := maven.UpdateParent()

	assert.Nil(t, err)
	assert.NotZero(t, updateMessage)
	assert.True(t, strings.HasPrefix(updateMessage, "Updating parent from 1.3.7.RELEASE to "), "but was : " + updateMessage)
}

func initMaven() *Maven {
	logger := *golog.New(os.Stderr, log.Debug)
	maven := NewMaven(logger)
	return maven
}

func setupWithTestProject(testProjectName string) *Maven {
	setup()
	sourcePath := path.Dir(temporaryDirectoryForTests.Cwd + "/../test-projects/" + testProjectName)
	if err := fileutils.CpR(sourcePath, "x"); err != nil {
		panic(err)
	}
	if err := os.Chdir("x/" + testProjectName); err != nil {
		panic(err)
	}

	maven := initMaven()
	err := maven.DetermineCommand()
	if err != nil {
		panic(err)
	}

	return maven
}
