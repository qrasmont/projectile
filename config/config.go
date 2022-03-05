package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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

	StoreConfig(config, path)

	return nil
}

func ParseConfig(config *Config, path string) error {
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

func StoreConfig(config *Config, path string) error {
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

func ExtractCommandsFromActions(project *Project, args []string) ([]string, error) {
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

func PrintAllActionsFromConfig(project *Project, showSteps bool) {
	for _, project_action := range project.Actions {
		fmt.Println(project_action.Name)
		if showSteps {
			for _, step := range project_action.Steps {
				fmt.Println("\t", step)
			}
		}
	}
}

func GetProjectFromConfig(config *Config, project_path string, project *Project) error {
	for _, config_project := range config.Projects {
		if config_project.Path == project_path {
			*project = config_project
			return nil
		}
	}

	return errors.New("This project does not exists")
}

func CommandRunner(commands *[]string, workdir string) error {
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

func AddToProject(config *Config, workdir string, args []string) error {
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

	return nil
}

func AddToConfig(config *Config, workdir string, args []string) error {
	action := Action{Name: args[0], Steps: args[1:]}
	project := Project{Path: workdir, Actions: []Action{action}}

	config.Projects = append(config.Projects, project)

	return nil
}

func AppendToConfig(config *Config, workdir string, args []string) error {
	for _, project := range config.Projects {
		if project.Path == workdir {
			for i, action := range project.Actions {
				if action.Name == args[0] {
					project.Actions[i].Steps = append(action.Steps, args[1:]...)
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

func RemoveActionsFromConfig(config *Config, workdir string, args []string) error {
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
			return nil
		}
	}

	return nil
}
