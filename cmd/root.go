package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	projectPath string

	rootCmd = &cobra.Command{
		Use:   "projectile [command]",
		Short: "Execute mutiple commands as single actions",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&projectPath, "path", "p", "project path (default is the current workind directory)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
