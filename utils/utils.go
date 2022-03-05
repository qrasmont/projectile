package utils

import (
	"errors"
	"os"
	"path/filepath"
)

const DEFAULT_HOME_CONFIG = ".config/projectile.json"

func GetWorkDir(path string) (string, error) {
	var workdir string

	if path != "" {
		workdir = path
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return "", errors.New("Could not get working directory.")
		}
		workdir = wd
	}

	return workdir, nil
}
func GetConfigPath() string {
	home_dir, _ := os.UserHomeDir()
	config_path := filepath.Join(home_dir, DEFAULT_HOME_CONFIG)
	envPath := os.Getenv("PROJECTILE_CONFIG")
	if envPath != "" {
		config_path = envPath
	}

	return config_path
}

func Exist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
