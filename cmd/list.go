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
			handleWorkspaceError(err)
		}

		store := storage.NewStore(wm)

		notes, err := store.List()
		if err != nil {
			if jsonFlag {
				outputJSONError(exitorrors.ErrNoteList.Error())
			} else {
				fmt.Println(exitorrors.ErrNoteList.Error())
			}
			os.Exit(exitorrors.ExitConfigError)
		}

		if jsonFlag {
			var jsonNotes []map[string]string
			for _, note := range notes {
				jsonNotes = append(jsonNotes, map[string]string{
					"slug":       note.Slug,
					"path":       note.RelPath,
					"created_at": note.CrTime.UTC().Format(time.RFC3339),
					"updated_at": note.ModTime.UTC().Format(time.RFC3339),
				})
			}
			if err := outputJSON(jsonNotes); err != nil {
				outputJSONError(exitorrors.ErrMarshalJSON.Error())
				os.Exit(exitorrors.ExitGeneral)
			}
		} else {
			for _, note := range notes {
				crDateStr := note.CrTime.Format("2006-01-02 15:04")
				modDateStr := note.ModTime.Format("2006-01-02 15:04")
				fmt.Printf("%s\t%s\t%s\n", note.Slug, crDateStr, modDateStr)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
