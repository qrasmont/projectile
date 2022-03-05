package cmd

import (
	"github.com/quadstew/projectile/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(appendCmd)
}

var appendCmd = &cobra.Command{
	Use:   "append [action] [steps to append]",
	Short: "Append steps to an existing action.",
    Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &config.Config{}

		err := config.ParseConfig(cfg, ConfigPath)
		if err != nil {
			return err
		}

		err = config.AppendToConfig(cfg, ProjectPath, args)
		if err != nil {
			return err
		}

		err = config.StoreConfig(cfg, ConfigPath)
		if err != nil {
			return err
		}

		return nil
	},
}
