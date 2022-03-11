package cmd

import (
	"github.com/quadstew/projectile/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doCmd)

}

var doCmd = &cobra.Command{
	Use:   "do [actions to run]",
	Short: "Run a project action.",
    Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &config.Config{}

		err := config.ParseConfig(cfg, ConfigPath)
		if err != nil {
			return err
		}

		prj := &config.Project{}
		err = config.GetProjectFromConfig(cfg, ProjectPath, prj)
		if err != nil {
			return err
		}

		actions, err := config.ExtractActionsFromProject(prj, args)
		if err != nil {
			return err
		}

		err = config.CommandRunner(&actions, ProjectPath)
		if err != nil {
			return err
		}

		return nil
	},
}
