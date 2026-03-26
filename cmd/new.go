package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	exitorrors "hotnotego/internal/errors"
	slugifypkg "hotnotego/internal/slugify"
	"hotnotego/internal/storage"
	"hotnotego/internal/workspace"
)

var newPath string

var newCmd = &cobra.Command{
	Use:   "new [title]",
	Short: "Create a new note",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
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

		baseSlug := slugifypkg.Slugify(title)
		if baseSlug == "" {
			if jsonFlag {
				outputJSONError("invalid title: produces empty slug")
			} else {
				fmt.Println("invalid title: produces empty slug")
			}
			os.Exit(exitorrors.ExitInvalidInput)
		}

		var fullSlug string
		if newPath != "" {
			fullSlug = filepath.Join(newPath, baseSlug)
		} else {
			fullSlug = baseSlug
		}

		noteID := uuid.New()
		createdAt := time.Now().UTC().Format(time.RFC3339)
		updatedAt := createdAt
		frontmatter := fmt.Sprintf("---\nid: %s\ntitle: %s\ncreated_at: %s\nupdated_at: %s\ntags: []\n---\n\n# %s\n\n", noteID, title, createdAt, updatedAt, title)

		if err := store.Ensure(fullSlug, []byte(frontmatter)); err != nil {
			if errors.Is(err, os.ErrExist) {
				if jsonFlag {
					outputJSONError(fmt.Sprintf("note already exists: %s", fullSlug))
				} else {
					fmt.Printf("note already exists: %s\n", fullSlug)
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
			path, err := store.Path(fullSlug)
			if err != nil {
				outputJSONError(fmt.Sprintf("get note path: %v", err))
				os.Exit(exitorrors.ExitGeneral)
			}
			response := map[string]string{
				"status": "created",
				"slug":   fullSlug,
				"path":   path,
			}
			outputJSON(response)
		} else {
			fmt.Printf("Created note: %s\n", fullSlug)
		}
	},
}

func init() {
	newCmd.Flags().StringVar(&newPath, "path", "", "Subfolder path for the new note (e.g., projects/todo)")
	RootCmd.AddCommand(newCmd)
}
