package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(editCmd)

}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open the config in $EDITOR",
    Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			return errors.New("Cannot open an editor. EDITOR not set.")
		}

		openEditor(editor, ConfigPath)

		return nil
	},
}

func openEditor(editor string, file string) error {
	runner := exec.Command(editor, file)
	runner.Stdin = os.Stdin
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stdout
	err := runner.Run()
	if err != nil {
		return err
	}
	return nil
}
