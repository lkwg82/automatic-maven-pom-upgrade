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
)

func init() {
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

func changeToTestDir() {
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
	changeToTestDir()
	defer cleanup()

	content := "#!/bin/sh\necho -n x"
	ioutil.WriteFile("./mvnw", []byte(content), 0700)
	file, _ := os.Create("maven.log")

	maven, _ := NewMaven(file)

	assert.Equal(t, maven.command, "./mvnw")
}

func TestMavenNotFound(t *testing.T) {
	changeToTestDir()
	defer cleanup()

	os.Setenv("PATH", "")
	file, _ := os.Create("maven.log")

	_, err := NewMaven(file)

	assert.NotNil(t, err, "should raised an error")
}

func TestMavenWrapperFound(t *testing.T) {
	changeToTestDir()
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

func copyTestProjectToTestDirectory(testProjectName string) {
	sourcePath := path.Dir(cwd + "/../test-projects/" + testProjectName)
	if err := fileutils.CpR(sourcePath, "x"); err != nil {
		panic(err)
	}
	if err := os.Chdir("x/" + testProjectName); err != nil {
		panic(err)
	}
}

func TestMavenParentPomUpdate(t *testing.T) {
	changeToTestDir()
	defer cleanup()

	copyTestProjectToTestDirectory("simple-parent-update")

	file, _ := os.Create("maven.log")

	// action
	maven, _ := NewMaven(file)
	maven.UpdateParent()

	var updateMessage string
	errors := make([]string, 0)
	logContent, _ := readFile("maven.log")
	lines := strings.Split(logContent, "\n")
	for _, line := range lines {
		updateToken := "[INFO] Updating parent from "
		if strings.HasPrefix(line, updateToken) {
			updateMessage = line
		} else {
			warnToken := "[WARNING]"
			errorToken := "[ERROR]"
			if strings.HasPrefix(line, warnToken) || strings.HasPrefix(line, errorToken) {
				errors=append(errors, line)
			}
		}
	}


	assert.Empty(t, errors)
	assert.NotZero(t, updateMessage)
	assert.True(t, strings.HasPrefix(updateMessage, "[INFO] Updating parent from 1.3.7.RELEASE to "))
}
