package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	exitorrors "hotnotego/internal/errors"
	slugifypkg "hotnotego/internal/slugify"
	"hotnotego/internal/storage"
	"hotnotego/internal/workspace"
)

var renameForce bool

var renameCmd = &cobra.Command{
	Use:   "rename <old> <new>",
	Short: "Rename a note in the current workspace",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		oldSlug := args[0]
		newTitle := args[1]

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

		resolvedSlug, err := store.Resolve(oldSlug)
		if errors.Is(err, storage.ErrNoteNotFound) {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("note not found: %s", oldSlug))
			} else {
				fmt.Printf("note not found: %s\n", oldSlug)
			}
			os.Exit(exitorrors.ExitNotFound)
		}
		if errors.Is(err, storage.ErrMultipleMatches) {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("multiple notes match '%s': use a more specific path", oldSlug))
			} else {
				fmt.Printf("multiple notes match '%s': use a more specific path\n", oldSlug)
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

		newSlug := computeNewSlug(resolvedSlug, newTitle)

		if newSlug == "" {
			if jsonFlag {
				outputJSONError("invalid title: produces empty slug")
			} else {
				fmt.Println("invalid title: produces empty slug")
			}
			os.Exit(exitorrors.ExitInvalidInput)
		}

		if resolvedSlug == newSlug {
			if jsonFlag {
				response := map[string]string{
					"status":  "unchanged",
					"slug":    resolvedSlug,
					"message": "slug unchanged",
				}
				outputJSON(response)
			} else {
				fmt.Printf("Slug unchanged: %s\n", resolvedSlug)
			}
			os.Exit(exitorrors.ExitSuccess)
		}

		if !renameForce && jsonFlag {
			outputJSONError("use --force to rename")
			os.Exit(exitorrors.ExitGeneral)
		}

		if !renameForce {
			fmt.Printf("Rename note '%s' to '%s'? [y/N]: ", resolvedSlug, newSlug)
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

		if err := store.Rename(resolvedSlug, newSlug); err != nil {
			if errors.Is(err, storage.ErrNoteAlreadyExists) {
				if jsonFlag {
					outputJSONError(fmt.Sprintf("note already exists: %s", newSlug))
				} else {
					fmt.Printf("note already exists: %s\n", newSlug)
				}
			} else {
				if jsonFlag {
					outputJSONError(fmt.Sprintf("rename note: %v", err))
				} else {
					fmt.Printf("rename note: %v\n", err)
				}
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if jsonFlag {
			response := map[string]string{
				"status":   "renamed",
				"old_slug": resolvedSlug,
				"new_slug": newSlug,
			}
			outputJSON(response)
		} else {
			fmt.Printf("Renamed note: %s → %s\n", resolvedSlug, newSlug)
		}
	},
}

func init() {
	renameCmd.Flags().BoolVar(&renameForce, "force", false, "Skip confirmation prompt")
	RootCmd.AddCommand(renameCmd)
}

func computeNewSlug(resolvedSlug, newTitle string) string {
	if strings.Contains(newTitle, "/") {
		parts := strings.Split(newTitle, "/")
		for i := range parts {
			parts[i] = slugifypkg.Slugify(parts[i])
		}
		return strings.Join(parts, "/")
	}

	oldBase := slugifypkg.Slugify(newTitle)
	if !strings.Contains(resolvedSlug, "/") {
		return oldBase
	}

	folder := resolvedSlug[:strings.LastIndex(resolvedSlug, "/")]
	return folder + "/" + oldBase
}
