package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
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
var CONFIG_FILE string = ""

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

func userWantsToRestoreConfig() bool {
	resp := ""
	for {
		fmt.Println("Could not parse your config.")
		fmt.Println("Do you want to restore the last backup?")
		fmt.Println("NOTE: This will override your current config file [yY/nN]: ")
		fmt.Scanln(&resp)
		if resp == "y" || resp == "Y" {
			return true
		} else if resp == "n" || resp == "N" {
			return false
		} else {
			fmt.Println("Wrong input!")
		}
	}
}

func restoreConfig(config *Config, path string) error {
	jsonFile, err := os.Open(path + ".bak")
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return err
	}

	storeConfig(config, path)

	return nil
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
		if userWantsToRestoreConfig() {
			restoreConfig(config, path)
			return nil
		}
		return err
	}

	return nil
}

func storeConfig(config *Config, path string) error {
	str_file, err := json.MarshalIndent(config, "", "   ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, str_file, 0644)
	if err != nil {
		return err
	}

	// Create a config back-up as well
	path = path + ".bak"
	err = ioutil.WriteFile(path, str_file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func extractCommandsFromActions(project *Project, args []string) ([]string, error) {
	var commands []string

	for _, action := range args {
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

func openEditor(editor string, file string) error {
	runner := exec.Command(editor, file)
	runner.Stdin = os.Stdin
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stdout
	err := runner.Run()
	if err != nil {
		return err
	}
	return nil
}

func addToConfig(config *Config, workdir string, args []string) error {
	// If PROJECT is empty it's not in our config
	if reflect.DeepEqual(PROJECT, Project{}) {
		action := Action{Name: args[0], Steps: args[1:]}
		project := Project{Path: workdir, Actions: []Action{action}}

		config.Projects = append(config.Projects, project)

		err := storeConfig(config, CONFIG_FILE)
		if err != nil {
			return err
		}

		return nil
	}

	for i, project := range config.Projects {
		// Search for project
		if project.Path == workdir {
			// Search for matching action name in project
			for _, action := range project.Actions {
				if action.Name == args[0] {
					return errors.New("This action already exists, maybe you want to use 'append'")
				}
			}

			// Action did not exist, add it to project
			action := Action{Name: args[0], Steps: args[1:]}
			config.Projects[i].Actions = append(project.Actions, action)
			break
		}
	}

	// Store our changes
	err := storeConfig(config, CONFIG_FILE)
	if err != nil {
		return err
	}
	return nil
}

func appendToConfig(config *Config, workdir string, args []string) error {
	if reflect.DeepEqual(PROJECT, Project{}) {
		return errors.New("Cannot append action, no exititing project.")
	}

	for _, project := range config.Projects {
		if project.Path == workdir {
			for i, action := range project.Actions {
				if action.Name == args[0] {
					project.Actions[i].Steps = append(action.Steps, args[1:]...)
					err := storeConfig(config, CONFIG_FILE)
					if err != nil {
						return err
					}
					return nil
				}
			}

			return errors.New("This action doesn't exist, cannot append. Maybe you want 'add'.")
		}
	}
	return nil
}

func contains(where []string, what string) bool {
	for _, item := range where {
		if item == what {
			return true
		}
	}
	return false
}

func removeActionsFromConfig(config *Config, workdir string, args []string) error {
	if len(args) == 0 {
		return errors.New("Need at least one action to remove.")
	}

	// TODO optimize this, linked list unmarshalling ?
	for i, project := range config.Projects {
		if project.Path == workdir {
			updated_actions := make([]Action, 0, len(project.Actions)-len(args))
			for _, action := range project.Actions {
				if !contains(args, action.Name) {
					updated_actions = append(updated_actions, action)
				}
			}

			config.Projects[i].Actions = updated_actions
			err := storeConfig(CONFIG, CONFIG_FILE)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func Init(cmdConfig *cmd.CmdConfig) error {
	CMD_CONFIG = cmdConfig

	home_dir, _ := os.UserHomeDir()
	CONFIG_FILE = filepath.Join(home_dir, DEFAULT_HOME_CONFIG)
	envPath := os.Getenv("PROJECTILE_CONFIG")
	if envPath != "" {
		CONFIG_FILE = envPath
	}

	hasConfig, err := hasConfigFile(CONFIG_FILE)
	if err != nil {
		return err
	}

	if !hasConfig {
		return errors.New("No config file at: " + CONFIG_FILE)
	}

	err = parseConfig(CONFIG, CONFIG_FILE)
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
	case cmd.Edit:
		editor := os.Getenv("EDITOR")
		if editor == "" {
			return errors.New("Cannot open an editor. EDITOR not set.")
		}
		openEditor(editor, CONFIG_FILE)
	case cmd.Get:
		printAllActionsFromConfig(&PROJECT)
	case cmd.Do:
		commands, err = extractCommandsFromActions(&PROJECT, CMD_CONFIG.Args)
		err = commandRunner(&commands, CMD_CONFIG.Path)
		if err != nil {
			return err
		}
	case cmd.Add:
		err = addToConfig(CONFIG, CMD_CONFIG.Path, CMD_CONFIG.Args)
		if err != nil {
			return err
		}
	case cmd.Append:
		err = appendToConfig(CONFIG, CMD_CONFIG.Path, CMD_CONFIG.Args)
		if err != nil {
			return err
		}
	case cmd.Remove:
		err = removeActionsFromConfig(CONFIG, CMD_CONFIG.Path, CMD_CONFIG.Args)
		if err != nil {
			return err
		}
	}

	return nil
}
