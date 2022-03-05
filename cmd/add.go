package cmd

import (
	"github.com/quadstew/projectile/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add [action name] [steps to execute]",
	Short: "Add a new action to the project.",
    Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &config.Config{}

		err := config.ParseConfig(cfg, ConfigPath)
		if err != nil {
			return err
		}

		prj := &config.Project{}
		err = config.GetProjectFromConfig(cfg, ProjectPath, prj)
		if err != nil {
			err = config.AddToConfig(cfg, ProjectPath, args)
			if err != nil {
				return err
			}
		} else {
			err = config.AddToProject(cfg, ProjectPath, args)
			if err != nil {
				return err
			}
		}

		err = config.StoreConfig(cfg, ConfigPath)
		if err != nil {
			return err
		}

		return nil
	},
}
