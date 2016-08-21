package lib

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"github.com/rafecolton/go-fileutils"
	"path"
	"strings"
	"log"
)

func TestDetectionOfMavenWrapper(t *testing.T) {
	setup()
	defer cleanup()

	content := "#!/bin/sh\necho -n x"
	ioutil.WriteFile("./mvnw", []byte(content), 0700)
	file, _ := os.Create("maven.log")

	maven, _ := NewMaven(file)

	assert.Equal(t, maven.command, "./mvnw")
}

func TestMavenNotFound(t *testing.T) {
	setup()
	defer cleanup()

	os.Setenv("PATH", "")
	file, _ := os.Create("maven.log")

	_, err := NewMaven(file)

	assert.NotNil(t, err, "should raised an error")
}

func TestMavenWrapperFound(t *testing.T) {
	setup()
	defer cleanup()

	os.Setenv("PATH", ".")

	content := "#!/bin/sh\necho -n x"
	ioutil.WriteFile("mvn", []byte(content), 0700)
	file, _ := os.Create("maven.log")

	// action
	maven, _ := NewMaven(file)

	logContent, _ := readFile("maven.log")

	assert.Equal(t, maven.command, "mvn")
	assert.Equal(t, logContent, "x")
}

func setupWithTestProject(testProjectName string) (*Maven) {
	setup()
	sourcePath := path.Dir(temporaryDirectoryForTests.Cwd + "/../test-projects/" + testProjectName)
	if err := fileutils.CpR(sourcePath, "x"); err != nil {
		panic(err)
	}
	if err := os.Chdir("x/" + testProjectName); err != nil {
		panic(err)
	}
	logFile, _ := os.Create("maven.log")

	maven, err := NewMaven(logFile)
	if err != nil {
		log.Print(readFile(logFile.Name()))
		panic(err)
	}

	return maven
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
