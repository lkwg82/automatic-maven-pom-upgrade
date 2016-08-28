package main

import (
	"fmt"
	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/golog"
	"github.com/droundy/goopt"
	. "github.com/lkwg82/automatic-maven-pom-upgrade/lib"
	"os"
)

var optVerbose = goopt.Flag([]string{
	"-v", "--verbose"},
	[]string{"--quiet"}, "output verbosely",
	"be quiet, instead")

var optType = goopt.Alternatives([]string{"--type"}, []string{"help", "parent"}, "type of upgrade")
var logger golog.Logger

func main() {
	parseParameter()

	if *optType == "help" {
		fmt.Print(goopt.Usage())
		os.Exit(0)
	}

	git := NewGit(logger)

	assert(git.IsInstalled(), "need git to be installed or in the PATH")
	assert(git.HasRepo(), "need called from a directory, which has a repository")
	assert(!git.IsDirty(), "repository is dirty, plz commit or reset")

	maven := NewMaven(logger)
	err := maven.DetermineCommand()
	if err != nil {
		logger.Emergency(err)
		os.Exit(1)
	}

	switch *optType {
	case "parent": updateParent(git, maven)
	default:
		panic("should never reach this point, wrong goopt config")
	}
}

func updateParent(git *Git, maven *Maven) {
	changeBranch(git)
	updated, message, err := maven.UpdateParent()
	if err != nil {
		logger.Errorf("parent update failed: %s", err)
		os.Exit(1)
	}

	if updated {
		git.Commit(message)
	} else {
		fmt.Printf("update not needed: %s \n",message)
	}
}
func changeBranch(git *Git) {
	if branch := "autoupdate_" + *optType; git.BranchExists(branch) {
		git.BranchCheckoutExisting(branch)
	} else {
		git.BranchCheckoutNew(branch)
	}
}

func assert(status bool, hint string) {
	if !status {
		logger.Error(hint)
		os.Exit(1)
	}
}

func parseParameter() {
	goopt.Summary = "automatic upgrade maven projects"
	goopt.Version = "0.1"
	goopt.Parse(nil)

	if *optVerbose {
		logger = *golog.New(os.Stderr, log.Debug)
	} else {
		logger = *golog.New(os.Stderr, log.Warning)
	}
}
