package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	exitorrors "hotnotego/internal/errors"
	"hotnotego/internal/storage"
	"hotnotego/internal/workspace"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:     "delete <slug>",
	Short:   "Delete a note from the current workspace",
	Aliases: []string{"del"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slug := args[0]

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

		resolvedSlug, err := store.Resolve(slug)
		if errors.Is(err, storage.ErrNoteNotFound) {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("note not found: %s", slug))
			} else {
				fmt.Printf("note not found: %s\n", slug)
			}
			os.Exit(exitorrors.ExitNotFound)
		}
		if errors.Is(err, storage.ErrMultipleMatches) {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("multiple notes match '%s': use a more specific path", slug))
			} else {
				fmt.Printf("multiple notes match '%s': use a more specific path\n", slug)
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

		if !deleteForce && jsonFlag {
			outputJSONError("use --force to delete")
			os.Exit(exitorrors.ExitGeneral)
		}

		if !deleteForce {
			fmt.Printf("Delete note '%s'? [y/N]: ", slug)
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				os.Exit(exitorrors.ExitGeneral)
			}
			input = strings.TrimSpace(input)
			if input != "y" && input != "Y" {
				os.Exit(exitorrors.ExitGeneral)
			}
		}

		if err := store.Delete(resolvedSlug); err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("delete note: %v", err))
			} else {
				fmt.Printf("delete note: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if jsonFlag {
			response := map[string]string{
				"status": "deleted",
				"slug":   resolvedSlug,
			}
			outputJSON(response)
		} else {
			fmt.Printf("Deleted note: %s\n", resolvedSlug)
		}
	},
}

func init() {
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Skip confirmation prompt")
	RootCmd.AddCommand(deleteCmd)
}
