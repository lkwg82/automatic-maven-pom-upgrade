package lib

import "io/ioutil"

func readFile(pathname string) (string, error) {
	bytes, err := ioutil.ReadFile("maven.log")
	n := len(bytes)
	return string(bytes[:n]), err
}