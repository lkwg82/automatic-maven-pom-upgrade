package lib

import (
	"testing"
	"os"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
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

	bytes, _ := ioutil.ReadFile("maven.log")
	n := len(bytes)
	logContent := string(bytes[:n])

	assert.Equal(t, maven.command, "mvn")
	assert.Equal(t, logContent, "x")
}