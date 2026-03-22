package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/spf13/cobra"
	"hotnotego/internal/storage"
)

var openCmd = &cobra.Command{
	Use:   "open [title]",
	Short: "Open a note for editing",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		store := storage.NewStore(dataDir)
		path := store.Path(title)
		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("Error opening note: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(content))
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}