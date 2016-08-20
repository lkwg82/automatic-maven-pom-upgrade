package lib

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"github.com/rafecolton/go-fileutils"
	"path"
	"strings"
)

var (
	cwd string
	testDirectory string
	logFile *os.File
)

func init() {
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

func setup() {
	if len(testDirectory) > 0 {
		panic("testDirectory is already set: " + testDirectory)
	}

	path, err := ioutil.TempDir("", "x")
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(path, 0655); err != nil {
		panic(err)
	}

	if err := os.Chdir(path); err != nil {
		panic(err)
	}
	testDirectory = path
}

func cleanup() {
	if _cwd, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		if len(testDirectory) > 0 && _cwd != cwd {
			if err := os.Chdir(cwd); err != nil {
				panic(err)
			}
			if err := os.RemoveAll(testDirectory); err != nil {
				panic(err)
			}
			testDirectory = ""
		}
	}
}

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
	sourcePath := path.Dir(cwd + "/../test-projects/" + testProjectName)
	if err := fileutils.CpR(sourcePath, "x"); err != nil {
		panic(err)
	}
	if err := os.Chdir("x/" + testProjectName); err != nil {
		panic(err)
	}
	logFile, _ = os.Create("maven.log")

	maven, err := NewMaven(logFile)
	if err != nil {
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
