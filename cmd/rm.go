package cmd

import (
	"github.com/quadstew/projectile/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm [actions to remove]",
	Short: "Remove the actions listed.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &config.Config{}

		err := config.ParseConfig(cfg, ConfigPath)
		if err != nil {
			return err
		}

		err = config.RemoveActionsFromConfig(cfg, ProjectPath, args)
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
