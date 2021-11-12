package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const CONFIG_FILE = ".projectile.json"

type Action struct {
	Name  string
	Steps []string
}

type Config struct {
	Actions []Action
}

func parseConfig(config *Config, path *string) {
	jsonFile, err := os.Open(*path)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, config)
	if err != nil {
		fmt.Println(err)
	}
}

func hasConfigFile(file *string) bool {

	matches, err := filepath.Glob(*file)

	if err != nil {
		fmt.Println(err)
	}

	if len(matches) == 1 {
		return true
	}
	return false
}

func executeActions(config *Config, actions *[]string, workdir *string) {

	for _, a := range *actions {
		matched := false

		for _, b := range config.Actions {
			if a == b.Name {
				matched = true
				for _, cmd := range b.Steps {
					args := strings.Fields(cmd)
					runner := exec.Command(args[0], args[1:]...)
					runner.Dir = *workdir
					err := runner.Run()
					if err != nil {
						panic(err)
					}
					fmt.Printf("%s", cmd)
				}
			}
		}
		if !matched {
			log.Fatal(errors.New("No match for action in config."))
		}
	}
}

func main() {
	path := flag.String("p", "", "The project's path.")

	flag.Parse()
	actions := flag.Args()

	if len(actions) == 0 {
		log.Fatal(errors.New("Need at list one action."))
	}

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

	config_file := workdir + "/" + CONFIG_FILE
	if !hasConfigFile(&config_file) {
		log.Fatal(errors.New("No .projectile.json found!"))
	}

	var config Config
	parseConfig(&config, &config_file)

	executeActions(&config, &actions, &workdir)
}
