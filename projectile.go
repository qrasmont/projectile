package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	path := flag.String("p", "", "The project's path.")

	flag.Parse()
	action := flag.Args()

	fmt.Printf("path: %s, action: %s\n", *path, action)

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

	fmt.Printf("workdir set to: %s\n", workdir)
}
