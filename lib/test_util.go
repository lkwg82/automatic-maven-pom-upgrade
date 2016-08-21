package lib

import (
	"os"
	"io/ioutil"
	"os/exec"
)

type TemporaryDirectoryForTests struct {
	Cwd           string
	testDirectory string
}

func (t *TemporaryDirectoryForTests) Init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	t.Cwd = cwd
}

func (t *TemporaryDirectoryForTests) Setup() {
	if len(t.testDirectory) > 0 {
		panic("testDirectory is already set: " + t.testDirectory)
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
	t.testDirectory = path
}

func (t *TemporaryDirectoryForTests) Cleanup() {
	if _cwd, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		if len(t.testDirectory) > 0 && _cwd != t.Cwd {
			if err := os.Chdir(t.Cwd); err != nil {
				panic(err)
			}
			if err := os.RemoveAll(t.testDirectory); err != nil {
				panic(err)
			}
			t.testDirectory = ""
		}
	}
}


func execCmd(cmd string, args []string) {
	err := exec.Command(cmd, args...).Run()
	if err != nil {
		panic(err)
	}
}