package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"hotnotego/internal/storage"
	"os"
	"path/filepath"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	Run: func(cmd *cobra.Command, args []string) {
		store := storage.NewStore(dataDir)
		files, err := os.ReadDir(store.Root)
		if err != nil {
			fmt.Printf("Error reading notes directory: %v\n", err)
			os.Exit(1)
		}

		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".md" {
				fmt.Println(file.Name())
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
