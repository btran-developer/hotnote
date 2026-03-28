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
			if jsonFlag {
				outputJSONError(fmt.Sprintf("create workspace manager: %v", err))
			} else {
				fmt.Printf("create workspace manager: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		store := storage.NewStore(wm)

		resolvedSlug, err := store.Resolve(title)
		if errors.Is(err, storage.ErrNoteNotFound) {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("note not found: %s", title))
			} else {
				fmt.Printf("note not found: %s\n", title)
			}
			os.Exit(exitorrors.ExitNotFound)
		}
		if errors.Is(err, storage.ErrMultipleMatches) {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("multiple notes match '%s': use a more specific path", title))
			} else {
				fmt.Printf("multiple notes match '%s': use a more specific path\n", title)
			}
			os.Exit(exitorrors.ExitInvalidInput)
		}
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("resolve note: %v", err))
			} else {
				fmt.Printf("resolve note: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		path, err := store.Path(resolvedSlug)
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("get note path: %v", err))
			} else {
				fmt.Printf("get note path: %v\n", err)
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
			fmt.Printf("open editor: %v\n", err)
			os.Exit(exitorrors.ExitGeneral)
		}
	},
}

func init() {
	RootCmd.AddCommand(openCmd)
}
