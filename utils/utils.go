package utils

import (
	"errors"
	"os"
	"path/filepath"
)

const DEFAULT_HOME_CONFIG = ".config/projectile.json"

func GetWorkDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.New("Could not get working directory.")
	}

	return wd, nil
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

func FormatSubDir(path *string) {
	if *path == "" {
		return
	}

	if (*path)[:2] == "./" {
		*path = (*path)[2:]
	}

	for (*path)[:1] == "/" {
		*path = (*path)[1:]
	}

	for (*path)[len(*path) - 1:] == "/" {
		*path = (*path)[:len(*path) - 1]
	}
}
