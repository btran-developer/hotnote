package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	exitorrors "hotnotego/internal/errors"
	"hotnotego/internal/pathutil"
	slugifypkg "hotnotego/internal/slugify"
	"hotnotego/internal/workspace"
)

var folderCreateCmd = &cobra.Command{
	Use:     "create <folder>",
	Short:   "Create a folder in the current workspace",
	Aliases: []string{"new", "cr"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		folder := slugifypkg.SlugifyPath(args[0])
		if folder == "" {
			if jsonFlag {
				outputJSONError("invalid folder name: produces empty slug")
			} else {
				fmt.Println("invalid folder name: produces empty slug")
			}
			os.Exit(exitorrors.ExitInvalidInput)
		}

		wm, err := workspace.NewManager()
		if err != nil {
			handleWorkspaceError(err)
		}

		_, wsPath, err := wm.Current()
		if err != nil {
			if jsonFlag {
				outputJSONError(exitorrors.ErrWorkspaceNotInit.Error())
			} else {
				fmt.Println(exitorrors.ErrWorkspaceNotInit.Error())
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

		folderPath := validation.FolderPath

		if _, err := os.Stat(folderPath); err == nil {
			if jsonFlag {
				outputJSONError(exitorrors.ErrFolderExists.Error())
			} else {
				fmt.Println(exitorrors.ErrFolderExists.Error())
			}
			os.Exit(exitorrors.ExitGeneral)
		} else if !os.IsNotExist(err) {
			if jsonFlag {
				outputJSONError(exitorrors.ErrFolderRead.Error())
			} else {
				fmt.Println(exitorrors.ErrFolderRead.Error())
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if err := os.MkdirAll(folderPath, 0755); err != nil {
			if jsonFlag {
				outputJSONError(exitorrors.WithContext(exitorrors.ErrFolderCreate, folderPath))
			} else {
				fmt.Println(exitorrors.WithContext(exitorrors.ErrFolderCreate, folderPath))
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if jsonFlag {
			response := map[string]string{
				"status": "created",
				"folder": folder,
				"path":   folderPath,
			}
			outputJSON(response)
		} else {
			fmt.Printf("Created folder: %s\n", folder)
		}
	},
}

func init() {
	folderCmd.AddCommand(folderCreateCmd)
}
