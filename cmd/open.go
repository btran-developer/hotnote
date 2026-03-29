package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	exitorrors "hotnotego/internal/errors"
	"hotnotego/internal/storage"
	"hotnotego/internal/workspace"
)

var openCmd = &cobra.Command{
	Use:     "open [title]",
	Short:   "Open a note for editing",
	Aliases: []string{"op"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		wm, err := workspace.NewManager()
		if err != nil {
			handleWorkspaceError(err)
		}

		store := storage.NewStore(wm)

		resolvedSlug, err := store.Resolve(title)
		if errors.Is(err, storage.ErrNoteNotFound) {
			if jsonFlag {
				outputJSONError(exitorrors.ErrNoteNotFound.Error())
			} else {
				fmt.Println(exitorrors.ErrNoteNotFound.Error())
			}
			os.Exit(exitorrors.ExitNotFound)
		}
		if errors.Is(err, storage.ErrMultipleMatches) {
			if jsonFlag {
				outputJSONError(exitorrors.ErrMultipleMatches.Error())
			} else {
				fmt.Println(exitorrors.ErrMultipleMatches.Error())
			}
			os.Exit(exitorrors.ExitInvalidInput)
		}
		if err != nil {
			if jsonFlag {
				outputJSONError(exitorrors.ErrNoteResolve.Error())
			} else {
				fmt.Println(exitorrors.ErrNoteResolve.Error())
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		path, err := store.Path(resolvedSlug)
		if err != nil {
			if jsonFlag {
				outputJSONError(exitorrors.ErrNotePath.Error())
			} else {
				fmt.Println(exitorrors.ErrNotePath.Error())
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if jsonFlag {
			response := map[string]string{
				"status": "opened",
				"path":   path,
			}
			outputJSON(response)
			return
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim"
		}

		editorCmd := exec.Command(editor, path)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		if err := editorCmd.Run(); err != nil {
			fmt.Println(exitorrors.WithContext(exitorrors.ErrOpenEditor, editor))
			os.Exit(exitorrors.ExitGeneral)
		}
	},
}

func init() {
	RootCmd.AddCommand(openCmd)
}
