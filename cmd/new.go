package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"hotnotego/internal/storage"
	"hotnotego/internal/workspace"
)

var newCmd = &cobra.Command{
	Use:   "new [title]",
	Short: "Create a new note",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		// Generate a slug from title (lowercase, hyphen-separated)
		id := slugify(title)
		wm, err := workspace.NewManager()
		if err != nil {
			fmt.Printf("Error creating workspace manager: %v\n", err)
			os.Exit(1)
		}
		store := storage.NewStore(wm)
		file, err := store.Ensure(id)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				fmt.Printf("Error: note '%s' already exists\n", id)
			} else {
				fmt.Printf("Error creating note: %v\n", err)
			}
			os.Exit(1)
		}

		// Create frontmatter with UUID and timestamps
		noteID := uuid.New()
		createdAt := time.Now().UTC().Format(time.RFC3339)
		updatedAt := createdAt
		frontmatter := fmt.Sprintf("---\nid: %s\ntitle: %s\ncreated_at: %s\nupdated_at: %s\ntags: []\n---\n\n# %s\n\n", noteID, title, createdAt, updatedAt, title)

		// Write frontmatter to the file
		if _, err := file.Write([]byte(frontmatter)); err != nil {
			fmt.Printf("Error writing frontmatter: %v\n", err)
			os.Exit(1)
		}
		// Ensure data is written to disk
		if err := file.Sync(); err != nil {
			fmt.Printf("Error syncing file: %v\n", err)
			os.Exit(1)
		}

		defer file.Close()
		fmt.Printf("Created note: %s\n", id)
	},
}

func init() {
	RootCmd.AddCommand(newCmd)
}

// slugify converts a string to a slug (lowercase, hyphen-separated)
func slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)
	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	// Remove any non-alphanumeric characters except hyphens
	// This is a simple implementation - in a real app you might want to use a proper slug library
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}
