package main

import (
	"fmt"
	"os"

	"github.com/quadstew/projectile/cmd"
	"github.com/quadstew/projectile/project"
)

func Bye(err error) {

	fmt.Println(err)
	os.Exit(1)
}

func main() {
	// Parse cmd args & flags
	cmdConfig, err := cmd.New()
	if err != nil {
		Bye(err)
	}

    // Parse the project config & set the actions to perform
	err = project.Init(&cmdConfig)
	if err != nil {
		Bye(err)
	}

    // Run actions
	err = project.Run()
	if err != nil {
		Bye(err)
	}
}
