package cmd

import (
	"github.com/quadstew/projectile/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "List project actions.",
    Args: cobra.NoArgs,
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

		config.PrintAllActionsFromConfig(prj)

		return nil
	},
}
