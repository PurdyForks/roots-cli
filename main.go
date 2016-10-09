package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	version "github.com/hashicorp/go-version"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func checkMsg(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		panic(err)
	}
}

func checkTrellisRequirements() {
	checkAnsible()
}

func checkAnsible() {
	ansibleCmd, err := exec.LookPath("ansible")
	checkMsg(err, "You need ansible")
	ansibleVersionCommand := exec.Command(ansibleCmd, "--version")
	ansibleVersionOutput, err := ansibleVersionCommand.Output()
	versionText := string(ansibleVersionOutput)
	versionArray := strings.Split(versionText, "\n")
	userVersion, err := version.NewVersion(strings.Split(versionArray[0], " ")[1])
	check(err)

	acceptableVersions := ">= 2.0.2.0"
	constraint, err := version.NewConstraint(acceptableVersions)
	check(err)

	if !constraint.Check(userVersion) {
		fmt.Println("Invalid version of ansible. We require " + acceptableVersions)
		os.Exit(1)
	}
}

func cloneRepo(repo string, location string) {
	git, err := exec.LookPath("git")
	check(err)

	gitCommand := exec.Command(
		git,
		"clone",
		"https://github.com/roots/"+repo,
		location,
	)
	gitCommand.Stderr = os.Stdout
	gitCommand.Run()
}

func generateTrellis(args []string) {
	checkTrellisRequirements()
	var target string
	var err error
	if len(args) == 0 {
		target, err = filepath.Abs("trellis")
	} else {
		target, err = filepath.Abs(args[0])
	}
	check(err)
	cloneRepo("trellis", target)
	fmt.Println(target)
}

func scaffoldProject(args []string) {
	if args[0] == "trellis" {
		generateTrellis(args[1:])
	}
}

func main() {
	fmt.Println("hello")
	args := os.Args[1:]

	command := args[0]

	if command == "new" {
		scaffoldProject(args[1:])
	}
}
