package main

import (
	"fmt"
	"github.com/droundy/goopt"
	"log"
	"os"
	. "github.com/lkwg82/automatic-maven-pom-upgrade/lib"
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

	if *optType == "parent" {
		file, _ := os.Create("maven.log")

		maven, err := NewMaven(file)
		if err != nil {
			log.Fatalf("failed to initialize maven: %s",  err)
		}
		maven.UpdateParent()
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
