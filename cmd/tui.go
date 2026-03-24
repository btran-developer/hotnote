package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"hotnotego/internal/tui"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the TUI interface",
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := tui.New()
		if err != nil {
			return fmt.Errorf("create TUI app: %w", err)
		}
		if err := app.Run(); err != nil {
			return fmt.Errorf("run TUI: %w", err)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(tuiCmd)
}
