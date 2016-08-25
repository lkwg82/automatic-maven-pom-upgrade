package main

import (
	"fmt"
	"github.com/droundy/goopt"
	. "github.com/lkwg82/automatic-maven-pom-upgrade/lib"
	"log"
	"os"
)

var optVerbose = goopt.Flag([]string{
	"-v", "--verbose"},
	[]string{"--quiet"}, "output verbosely",
	"be quiet, instead")

var optType = goopt.Alternatives([]string{"--type"}, []string{"help", "parent"}, "type of upgrade")

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
}

func main() {
	parseParameter()

	gitLog, _ := os.Create("git.log")
	git := NewGit(gitLog)

	assert(git.IsInstalled(), "need git to be installed or in the PATH")
	assert(git.HasRepo(), "need called from a directory, which has a repository")
	assert(!git.IsDirty(), "repository is dirty, plz commit or reset")

	mavenLog, _ := os.Create("maven.log")
	maven := NewMaven(mavenLog)
	err := maven.DetermineCommand()
	if err != nil {
		log.Fatal(err)
	}

	if *optType == "parent" {
		if message, err := maven.UpdateParent(); err != nil {
			git.CommitMessage = message
		}
	}

	git.Commit()
}

func assert(status bool, hint string) {
	if (!status) {
		log.Fatal("ERROR: " + hint)
	}
}

func parseParameter() {
	goopt.Summary = "automatic upgrade maven projects"
	goopt.Version = "0.1"
	goopt.Parse(nil)

	if *optType == "help" {
		fmt.Print(goopt.Usage())
		os.Exit(0)
	}
}
