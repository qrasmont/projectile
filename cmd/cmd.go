package cmd

// The cmd package returns a struct that has all the information to know:
//      What I should do and where to do it.

import (
	"errors"
	"os"

	"github.com/jessevdk/go-flags"
)

const (
	Do     string = "do"
	Get           = "get"
	Edit          = "edit"
	Add           = "add"
	Append        = "append"
	Remove        = "rm"
)

type CmdConfig struct {
	Path    string
	Command string
	Args    []string
}

type Options struct {
	Path string `short:"p" long:"path" description:"The project path, the current working directory by default."`
}

func getWorkDir(path string) (string, error) {
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

func setCommand(args []string, cmdConfig *CmdConfig) error {

	switch args[0] {
	case Get:
		if len(args) > 1 {
			return errors.New("'get' doesn't need aditional arguments")
		}
		cmdConfig.Command = Get
		return nil
	case Do:
		if len(args) == 1 {
			return errors.New("'do' needs at least 1 argument")
		}
		cmdConfig.Command = Do
		cmdConfig.Args = args[1:]
		return nil
	case Edit:
		if len(args) > 1 {
			return errors.New("'edit' doesn't need aditional arguments")
		}
		cmdConfig.Command = Edit
		return nil
	case Add:
		cmdConfig.Command = Add
		cmdConfig.Args = args[1:]
		return nil
	case Append:
		cmdConfig.Command = Append
		cmdConfig.Args = args[1:]
		return nil
	case Remove:
		cmdConfig.Command = Remove
		cmdConfig.Args = args[1:]
		return nil
	default:
		return errors.New("Unknown command: " + args[0])
	}
}

func New() (CmdConfig, error) {
	var cmdConf CmdConfig
	var opts Options

	parser := flags.NewParser(&opts, flags.Default)
	args, err := parser.Parse()
	if err != nil {
		return cmdConf, err
	}

	if len(args) == 0 {
		return cmdConf, errors.New("No args provided")
	}

	workdir, err := getWorkDir(opts.Path)
	if err != nil {
		return cmdConf, err
	}
	cmdConf.Path = workdir

	err = setCommand(args, &cmdConf)
	if err != nil {
		return cmdConf, err
	}

	return cmdConf, nil
}
