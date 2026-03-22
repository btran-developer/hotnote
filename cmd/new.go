package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"hotnotego/internal/storage"
)

var newCmd = &cobra.Command{
	Use:   "new [title]",
	Short: "Create a new note",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		// In a real implementation, we'd generate a slug from title
		id := title // placeholder
		store := storage.NewStore(dataDir)
		file, err := store.Ensure(id)
		if err != nil {
			fmt.Printf("Error creating note: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		fmt.Printf("Created note: %s\n", id)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}