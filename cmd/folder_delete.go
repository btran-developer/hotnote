package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	exitorrors "hotnotego/internal/errors"
	"hotnotego/internal/pathutil"
	"hotnotego/internal/workspace"
)

var folderDeleteForce bool

var folderDeleteCmd = &cobra.Command{
	Use:     "delete <folder>",
	Short:   "Delete a folder from the current workspace",
	Aliases: []string{"del"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		folder := args[0]

		wm, err := workspace.NewManager()
		if err != nil {
			handleWorkspaceError(err)
		}

		_, wsPath, err := wm.Current()
		if err != nil {
			if jsonFlag {
				outputJSONError("workspace not initialized")
			} else {
				fmt.Println("workspace not initialized")
			}
			os.Exit(exitorrors.ExitConfigError)
		}

		validation, err := pathutil.ValidateFolderPath(wsPath, folder)
		if err != nil {
			if errors.Is(err, pathutil.ErrPathOutsideWS) || errors.Is(err, pathutil.ErrInvalidPath) {
				if jsonFlag {
					outputJSONError(err.Error())
				} else {
					fmt.Println(err.Error())
				}
				os.Exit(exitorrors.ExitInvalidInput)
			}
			if jsonFlag {
				outputJSONError(exitorrors.ErrInvalidFolderPath.Error())
			} else {
				fmt.Println(exitorrors.ErrInvalidFolderPath.Error())
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if validation.AbsFolderPath == validation.AbsWsPath {
			if jsonFlag {
				outputJSONError("cannot delete workspace root")
			} else {
				fmt.Println("cannot delete workspace root")
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		folderPath := validation.FolderPath

		if _, err := os.Stat(folderPath); err != nil {
			if os.IsNotExist(err) {
				if jsonFlag {
					outputJSONError(exitorrors.ErrFolderNotFound.Error())
				} else {
					fmt.Println(exitorrors.ErrFolderNotFound.Error())
				}
				os.Exit(exitorrors.ExitNotFound)
			}
			if jsonFlag {
				outputJSONError(exitorrors.ErrFolderRead.Error())
			} else {
				fmt.Println(exitorrors.ErrFolderRead.Error())
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		entries, err := os.ReadDir(folderPath)
		if err != nil {
			if jsonFlag {
				outputJSONError(exitorrors.ErrFolderRead.Error())
			} else {
				fmt.Println(exitorrors.ErrFolderRead.Error())
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		hasContents := len(entries) > 0

		if hasContents && !folderDeleteForce && jsonFlag {
			outputJSONError("folder not empty: use --force to delete")
			os.Exit(exitorrors.ExitGeneral)
		}

		if hasContents && !folderDeleteForce {
			fmt.Printf("Delete folder '%s' and all contents? [y/N]: ", folder)
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

		if err := os.RemoveAll(folderPath); err != nil {
			if jsonFlag {
				outputJSONError(exitorrors.ErrFolderDelete.Error())
			} else {
				fmt.Println(exitorrors.ErrFolderDelete.Error())
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if jsonFlag {
			response := map[string]string{
				"status": "deleted",
				"folder": folder,
			}
			outputJSON(response)
		} else {
			fmt.Printf("Deleted folder: %s\n", folder)
		}
	},
}

func init() {
	folderDeleteCmd.Flags().BoolVar(&folderDeleteForce, "force", false, "Skip confirmation prompt")
	folderCmd.AddCommand(folderDeleteCmd)
}
