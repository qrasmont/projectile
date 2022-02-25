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

const DEFAULT_HOME_CONFIG = ".config/projectile.json"

type Action struct {
	Name  string
	Steps []string
}

type Project struct {
	Path    string
	Actions []Action
}

type Config struct {
	Projects []Project
}

var CONFIG *Config = &Config{}
var CMD_CONFIG *cmd.CmdConfig
var PROJECT Project = Project{}

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

func extractCommandsFromActions(project *Project, actions []string) ([]string, error) {
	var commands []string

	for _, action := range actions {
		matched := false

		for _, config_action := range project.Actions {
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

func printAllActionsFromConfig(project *Project) {
	for _, project_action := range project.Actions {
		fmt.Println(project_action.Name)
	}
}

func setProject(config *Config, workdir string) {
	for _, project := range config.Projects {
		if project.Path == workdir {
			PROJECT = project
			return
		}
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

	home_dir, _ := os.UserHomeDir()
	config_file := filepath.Join(home_dir, DEFAULT_HOME_CONFIG)
	envPath := os.Getenv("PROJECTILE_CONFIG")
	if envPath != "" {
		config_file = envPath
	}

	hasConfig, err := hasConfigFile(config_file)
	if err != nil {
		return err
	}

	if !hasConfig {
		return errors.New("No config file at: " + config_file)
	}

	err = parseConfig(CONFIG, config_file)
	if err != nil {
		return err
	}

	setProject(CONFIG, CMD_CONFIG.Path)

	return nil
}

func Run() error {
	var commands []string
	var err error

	switch CMD_CONFIG.Command {
	case cmd.Get:
		printAllActionsFromConfig(&PROJECT)
	case cmd.Do:
		commands, err = extractCommandsFromActions(&PROJECT, CMD_CONFIG.Actions)
		err = commandRunner(&commands, CMD_CONFIG.Path)
		if err != nil {
			return err
		}
	}

	return nil
}
