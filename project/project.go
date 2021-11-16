package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

var CONFIG *Config = &Config{}
var CMD_CONFIG *cmd.CmdConfig

func hasConfigFile(file string) (bool, error) {

	matches, err := filepath.Glob(file)

	if err != nil {
		return false, err
	}

	if len(matches) == 1 {
		return true, nil
	}
	return false, nil
}

func parseConfig(config *Config, path string) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return err
	}

	return nil
}

func extractCommandsFromActions(config *Config, actions []string) ([]string, error) {
	var commands []string

	for _, action := range actions {
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
			return commands, errors.New("No match for action in config.")
		}
	}

	return commands, nil
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

func commandRunner(commands *[]string, workdir string) error {
	for _, cmd := range *commands {
		args := strings.Fields(cmd)
		runner := exec.Command(args[0], args[1:]...)
		runner.Dir = workdir
		runner.Stdout = os.Stdout
		runner.Stderr = os.Stdout
		err := runner.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func Init(cmdConfig *cmd.CmdConfig) error {
	CMD_CONFIG = cmdConfig

	config_file := CMD_CONFIG.Path + "/" + CONFIG_FILE

	hasConfig, err := hasConfigFile(config_file)
	if err != nil {
		return err
	}

	if !hasConfig {
		return errors.New("No config file found in: " + CMD_CONFIG.Path)
	}

	err = parseConfig(CONFIG, config_file)
	if err != nil {
		return err
	}

	return nil
}

func Run() error {
	var commands []string
	var err error

	switch CMD_CONFIG.Command {
	case cmd.Get:
		printAllActionsFromConfig(CONFIG)
	case cmd.All:
		commands = extractAllCommands(CONFIG)
		err = commandRunner(&commands, CMD_CONFIG.Path)
		if err != nil {
			return err
		}
	case cmd.Do:
		commands, err = extractCommandsFromActions(CONFIG, CMD_CONFIG.Actions)
		err = commandRunner(&commands, CMD_CONFIG.Path)
		if err != nil {
			return err
		}
	}

	return nil
}
