package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func hasConfigFile(workdir *string) bool {

	matches, err := filepath.Glob(*workdir + "/.projectile.json")

	if err != nil {
		fmt.Println(err)
	}

	if len(matches) == 1 {
		return true
	}
	return false
}

func main() {
	path := flag.String("p", "", "The project's path.")

	flag.Parse()
	actions := flag.Args()

	fmt.Printf("path: %s, action: %s\n", *path, actions)

	workdir := ""
	if *path != "" {
		workdir = *path
	} else {
		// Get the cwd
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(errors.New("Could not get cwd."))
		}
		workdir = cwd
	}

	if !hasConfigFile(&workdir) {
		log.Fatal(errors.New("No .projectile.json found!"))
	}

	fmt.Printf("workdir set to: %s\n", workdir)
}
