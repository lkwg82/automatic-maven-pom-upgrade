package main

import (
	"bytes"
	"fmt"
	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/golog"
	"github.com/droundy/goopt"
	"github.com/lkwg82/automatic-maven-pom-upgrade/lib"
	"os"
	"path"
	"runtime"
)

var optQuiet = goopt.Flag([]string{"-q", "--quiet"}, nil, "suppress any output", "")
var optVersion = goopt.Flag([]string{"--version"}, nil, "show version", "")
var optType = goopt.Alternatives([]string{"-t", "--type"}, []string{"help", "parent"}, "type of upgrade")
var optVerbose = goopt.Flag([]string{"-v", "--verbose"}, nil, "output verbosely", "")

var logger golog.Logger

type funcErr func() error

func main() {

	parseParameter()

	git := lib.NewGit(logger)
	maven := lib.NewMaven(logger)

	if *optVersion {
		fmt.Printf("version %s\n", goopt.Version)
		os.Exit(0)
	}

	if *optType == "help" {
		fmt.Print(goopt.Usage())
		os.Exit(0)
	}

	exitOnError(git.CheckIsInstalled)
	exitOnError(git.CheckIsRepo)
	exitOnError(git.OptionalCheckIsDirty)

	exitOnError(maven.DetermineCommand)
	exitOnError(maven.ParseCommandline)

	switch *optType {
	case "parent":
		updateParent(git, maven)
	default:
		panic("should never reach this point, wrong goopt config")
	}
}

func exitOnError(fun funcErr) {
	err := fun()
	if err != nil {
		logger.Emergency(err)
		os.Exit(1)
	}
}

func echo(format string, arg ...string) {
	if !*optQuiet {
		if len(arg) == 0 {
			fmt.Println(format)
		} else {
			fmt.Printf(format+"\n", arg)
		}
	}
}

func updateParent(git *lib.Git, maven *lib.Maven) {
	updateBranch := "autoupdate_" + *optType
	if git.BranchExists(updateBranch) {
		git.BranchCheckoutExisting(updateBranch)
		git.OptionalAutoMergeMaster()
	}

	updated, message, err := maven.UpdateParent()
	if err != nil {
		logger.Errorf("parent update failed: %s", err)
		os.Exit(1)
	}

	if updated {
		echo("result: " + message)
		if !git.BranchExists(updateBranch) && git.BranchCurrent() != updateBranch {
			git.BranchCheckoutNew(updateBranch)
		}
		git.OptionalCommit(message, echo)
	} else {
		echo("update not needed: %s ", message)
	}
}

func parseParameter() {
	goopt.Version = "0.1"
	goopt.Summary = "automatic upgrade maven projects, " + goopt.Version
	goopt.Parse(nil)

	if *optVerbose {
		logger = *golog.New(os.Stderr, log.Debug)
		//
		oldFormatter := logger.Formatter
		logger.Formatter = func(buf *bytes.Buffer, level log.Level, args ...interface{}) {
			_, file, line, _ := runtime.Caller(3)
			args[0] = fmt.Sprintf("%s:%d %s", path.Base(file), line, args[0])
			oldFormatter(buf, level, args)

			//_, file2, line2, _ := runtime.Caller(4)
			//args[0] = fmt.Sprintf(" %s:%d -> %s", path.Base(file2), line2, args[0])
			//oldFormatter(buf, level, args)
			//
			//_, file3, line3, _ := runtime.Caller(5)
			//if !strings.Contains(file3, "/runtime/") {
			//	args[0] = fmt.Sprintf("  %s:%d -> %s", path.Base(file3), line3, args[0])
			//	oldFormatter(buf, level, args)
			//
		}

		logger.Debugf("running in verbose mode")
		logger.Debugf(" type: %s", *optType)
	} else {
		logger = *golog.New(os.Stderr, log.Warning)
	}
}
