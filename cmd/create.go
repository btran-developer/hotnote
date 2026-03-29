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

var createPath string
var createSlug string

var createCmd = &cobra.Command{
	Use:     "create [title]",
	Short:   "Create a new note",
	Aliases: []string{"new"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		wm, err := workspace.NewManager()
		if err != nil {
			handleWorkspaceError(err)
		}
		store := storage.NewStore(wm)

		var baseSlug string
		if createSlug != "" {
			baseSlug = slugifypkg.Slugify(createSlug)
		} else {
			baseSlug = slugifypkg.Slugify(title)
		}
		if baseSlug == "" {
			if jsonFlag {
				outputJSONError(exitorrors.ErrEmptySlug.Error())
			} else {
				fmt.Println(exitorrors.ErrEmptySlug.Error())
			}
			os.Exit(exitorrors.ExitInvalidInput)
		}

		var fullSlug string
		if createPath != "" {
			fullSlug = filepath.Join(createPath, baseSlug)
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
					outputJSONError(exitorrors.ErrNoteExists.Error())
				} else {
					fmt.Println(exitorrors.ErrNoteExists.Error())
				}
				os.Exit(exitorrors.ExitInvalidInput)
			} else {
				if jsonFlag {
					outputJSONError(exitorrors.ErrNoteCreate.Error())
				} else {
					fmt.Println(exitorrors.ErrNoteCreate.Error())
				}
				os.Exit(exitorrors.ExitGeneral)
			}
		}

		if jsonFlag {
			path, err := store.Path(fullSlug)
			if err != nil {
				outputJSONError(exitorrors.ErrNotePath.Error())
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
	createCmd.Flags().StringVar(&createPath, "path", "", "Subfolder path for the new note (e.g., projects/todo)")
	createCmd.Flags().StringVar(&createSlug, "slug", "", "Custom slug for the note")
	RootCmd.AddCommand(createCmd)
}
