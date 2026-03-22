package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	exitorrors "hotnotego/internal/errors"
	"hotnotego/internal/storage"
	"hotnotego/internal/workspace"
)

var openCmd = &cobra.Command{
	Use:   "open [title]",
	Short: "Open a note for editing",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		wm, err := workspace.NewManager()
		if err != nil {
			fmt.Printf("create workspace manager: %v\n", err)
			os.Exit(exitorrors.ExitGeneral)
		}

		store := storage.NewStore(wm)
		path, err := store.Path(title)
		if err != nil {
			fmt.Printf("get note path: %v\n", err)
			os.Exit(exitorrors.ExitGeneral)
		}

		// Check if file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Println("note not found")
			os.Exit(exitorrors.ExitNotFound)
		}

		// Determine editor to use
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim" // default fallback
		}

		// Open the file in the editor
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
