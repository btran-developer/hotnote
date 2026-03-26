package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	exitorrors "hotnotego/internal/errors"
	"hotnotego/internal/pathutil"
	"hotnotego/internal/workspace"
)

var mkdirCmd = &cobra.Command{
	Use:   "mkdir <folder>",
	Short: "Create a folder in the current workspace",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		folder := args[0]

		wm, err := workspace.NewManager()
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("create workspace manager: %v", err))
			} else {
				fmt.Printf("create workspace manager: %v\n", err)
			}
			os.Exit(exitorrors.ExitConfigError)
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
				outputJSONError(fmt.Sprintf("validate path: %v", err))
			} else {
				fmt.Printf("validate path: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		folderPath := validation.FolderPath

		if _, err := os.Stat(folderPath); err == nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("folder already exists: %s", folder))
			} else {
				fmt.Printf("folder already exists: %s\n", folder)
			}
			os.Exit(exitorrors.ExitGeneral)
		} else if !os.IsNotExist(err) {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("check folder: %v", err))
			} else {
				fmt.Printf("check folder: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if err := os.MkdirAll(folderPath, 0755); err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("create folder: %v", err))
			} else {
				fmt.Printf("create folder: %v\n", err)
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
	RootCmd.AddCommand(mkdirCmd)
}
