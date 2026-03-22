package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	exitorrors "hotnotego/internal/errors"
)

var RootCmd = &cobra.Command{
	Use:     "hotnote",
	Short:   "A terminal-first markdown note system",
	Long:    `A CLI for managing markdown notes`,
	Version: "0.1.0",
}

var dataDir string
var jsonFlag bool

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(exitorrors.ExitGeneral)
	}
}

func init() {
	// Global flags
	RootCmd.PersistentFlags().StringVarP(&dataDir, "data-dir", "d", "notes", "Data directory for notes")
	RootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Output in JSON format")

	// Set the run function for the root command (when no subcommand is given)
	RootCmd.Run = func(cmd *cobra.Command, args []string) {
		// If no command is given, print help
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	}

	// Note: Subcommands are added via init() functions in their respective files
}

// Create, list, open, render, and ai commands will be defined in their respective files
