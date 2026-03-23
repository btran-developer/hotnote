package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	exitorrors "hotnotego/internal/errors"
	"hotnotego/internal/storage"
	"hotnotego/internal/workspace"
)

var newCmd = &cobra.Command{
	Use:   "new [title]",
	Short: "Create a new note",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		id := slugify(title)
		wm, err := workspace.NewManager()
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("create workspace manager: %v", err))
			} else {
				fmt.Printf("create workspace manager: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}
		store := storage.NewStore(wm)

		noteID := uuid.New()
		createdAt := time.Now().UTC().Format(time.RFC3339)
		updatedAt := createdAt
		frontmatter := fmt.Sprintf("---\nid: %s\ntitle: %s\ncreated_at: %s\nupdated_at: %s\ntags: []\n---\n\n# %s\n\n", noteID, title, createdAt, updatedAt, title)

		if err := store.Ensure(id, []byte(frontmatter)); err != nil {
			if errors.Is(err, os.ErrExist) {
				if jsonFlag {
					outputJSONError("note already exists")
				} else {
					fmt.Println("note already exists")
				}
			} else {
				if jsonFlag {
					outputJSONError(fmt.Sprintf("create note: %v", err))
				} else {
					fmt.Printf("create note: %v\n", err)
				}
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if jsonFlag {
			path, err := store.Path(id)
			if err != nil {
				outputJSONError(fmt.Sprintf("get note path: %v", err))
				os.Exit(exitorrors.ExitGeneral)
			}
			response := map[string]string{
				"status": "created",
				"slug":   id,
				"path":   path,
			}
			outputJSON(response)
		} else {
			fmt.Printf("Created note: %s\n", id)
		}
	},
}

func init() {
	RootCmd.AddCommand(newCmd)
}

// slugify converts a string to a slug (lowercase, hyphen-separated)
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	s = result.String()
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	s = strings.Trim(s, "-")
	return s
}
