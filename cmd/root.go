package cmd

import (
	"encoding/json"
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
var prettyFlag bool

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(exitorrors.ExitGeneral)
	}
}

// outputJSON outputs data as JSON, using pretty formatting if prettyFlag is true
func outputJSON(v interface{}) error {
	var data []byte
	var err error

	if prettyFlag {
		data, err = json.MarshalIndent(v, "", "  ")
	} else {
		data, err = json.Marshal(v)
	}

	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}

// outputJSONError outputs an error as JSON, using pretty formatting if prettyFlag is true
func outputJSONError(errMsg string) {
	errorResponse := map[string]string{"error": errMsg}
	outputJSON(errorResponse)
}

func init() {
	// Global flags
	RootCmd.PersistentFlags().StringVarP(&dataDir, "data-dir", "d", "notes", "Data directory for notes")
	RootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Output in JSON format")
	RootCmd.PersistentFlags().BoolVar(&prettyFlag, "pretty", false, "Pretty-print JSON output (only with --json)")

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
