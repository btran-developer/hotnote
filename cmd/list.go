package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	exitorrors "hotnotego/internal/errors"
	"hotnotego/internal/storage"
	"hotnotego/internal/workspace"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all notes",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		wm, err := workspace.NewManager()
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("create workspace manager: %v", err))
			} else {
				fmt.Printf("create workspace manager: %v\n", err)
			}
			os.Exit(exitorrors.ExitConfigError)
		}

		store := storage.NewStore(wm)

		notes, err := store.List()
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("list notes: %v", err))
			} else {
				fmt.Printf("list notes: %v\n", err)
			}
			os.Exit(exitorrors.ExitConfigError)
		}

		if jsonFlag {
			var jsonNotes []map[string]string
			for _, note := range notes {
				jsonNotes = append(jsonNotes, map[string]string{
					"slug":       note.Slug,
					"path":       note.RelPath,
					"updated_at": note.ModTime.UTC().Format(time.RFC3339),
				})
			}
			if err := outputJSON(jsonNotes); err != nil {
				outputJSONError(fmt.Sprintf("marshal JSON: %v", err))
			}
		} else {
			for _, note := range notes {
				dateStr := note.ModTime.Format("2006-01-02 15:04")
				fmt.Printf("%s\t%s\n", note.Slug, dateStr)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
