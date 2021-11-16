package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/quadstew/projectile/cmd"
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

func extractCommandsFromActions(config *Config, actions *[]string) []string {
	var commands []string

	for _, action := range *actions {
		matched := false

		for _, config_action := range config.Actions {
			if action == config_action.Name {
				matched = true
				for _, cmd := range config_action.Steps {
					commands = append(commands, cmd)
				}
			}
		}

		if !matched {
			log.Fatal(errors.New("No match for action in config."))
		}
	}

	return commands
}

func extractAllCommands(config *Config) []string {
	var commands []string

	for _, config_action := range config.Actions {
		for _, cmd := range config_action.Steps {
			commands = append(commands, cmd)
		}
	}

	return commands
}

func printAllActionsFromConfig(config *Config) {
	for _, config_action := range config.Actions {
		fmt.Println(config_action.Name)
	}
}

func commandRunner(commands *[]string, workdir *string) {
	for _, cmd := range *commands {
		args := strings.Fields(cmd)
		runner := exec.Command(args[0], args[1:]...)
		runner.Dir = *workdir
		runner.Stdout = os.Stdout
		runner.Stderr = os.Stdout
		err := runner.Run()
		if err != nil {
			panic(err)
		}
	}
}

func Bye(err error) {

	fmt.Println(err)
	os.Exit(1)
}

func main() {
	cmdConfig, err := cmd.New()
	if err != nil {
		Bye(err)
	}

	config_file := cmdConfig.Path + "/" + CONFIG_FILE
	if !hasConfigFile(&config_file) {
		log.Fatal(errors.New("No .projectile.json found!"))
	}

	var config Config
	parseConfig(&config, &config_file)

	if cmdConfig.Command == cmd.Get {
		printAllActionsFromConfig(&config)
		return
	}

	var commands []string
	if cmdConfig.Command != cmd.All {
		commands = extractCommandsFromActions(&config, &cmdConfig.Actions)
	} else {
		commands = extractAllCommands(&config)
	}

	commandRunner(&commands, &cmdConfig.Path)
}
