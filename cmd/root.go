package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "hotnote",
	Short: "A terminal-first markdown note system",
	Long:  `A CLI for managing markdown notes`,
}

var dataDir string

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	RootCmd.PersistentFlags().StringVarP(&dataDir, "data-dir", "d", "notes", "Data directory for notes")
	// Add subcommands
	RootCmd.AddCommand(newCmd)
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(openCmd)
	RootCmd.AddCommand(renderCmd)
	RootCmd.AddCommand(aiCmd)
}

// Create, list, open, render, and ai commands will be defined in their respective files
