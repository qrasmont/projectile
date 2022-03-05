package cmd

import (
	"fmt"
	"os"

	"github.com/quadstew/projectile/utils"
	"github.com/spf13/cobra"
)

var (
	ProjectPath string
	ConfigPath  string
	pathFlag    string

	rootCmd = &cobra.Command{
		Use:   "projectile [command]",
		Short: "Execute mutiple commands as single actions",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if pathFlag != "" {
				ProjectPath = pathFlag
			} else {
				wd, err := utils.GetWorkDir()
				if err != nil {
					return err
				}
				ProjectPath = wd
			}

			ConfigPath = utils.GetConfigPath()

			return nil
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&pathFlag, "path", "", "project path (default is the current working directory)")
}
