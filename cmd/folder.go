package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	exitorrors "hotnotego/internal/errors"
	"hotnotego/internal/pathutil"
	"hotnotego/internal/workspace"
)

var folderCmd = &cobra.Command{
	Use:   "folder",
	Short: "Folder management",
}

var folderRenameCmd = &cobra.Command{
	Use:     "rename <old> <new>",
	Short:   "Rename a folder in the current workspace",
	Aliases: []string{"rn"},
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		oldFolder := args[0]
		newFolder := args[1]

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

		oldValidation, err := pathutil.ValidateFolderPath(wsPath, oldFolder)
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

		if oldValidation.AbsFolderPath == oldValidation.AbsWsPath {
			if jsonFlag {
				outputJSONError("cannot rename workspace root")
			} else {
				fmt.Println("cannot rename workspace root")
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if _, err := os.Stat(oldValidation.FolderPath); err != nil {
			if os.IsNotExist(err) {
				if jsonFlag {
					outputJSONError(fmt.Sprintf("folder not found: %s", oldFolder))
				} else {
					fmt.Printf("folder not found: %s\n", oldFolder)
				}
				os.Exit(exitorrors.ExitNotFound)
			}
			if jsonFlag {
				outputJSONError(fmt.Sprintf("check folder: %v", err))
			} else {
				fmt.Printf("check folder: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		newValidation, err := pathutil.ValidateFolderPath(wsPath, newFolder)
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

		if _, err := os.Stat(newValidation.FolderPath); err != nil {
			if !os.IsNotExist(err) {
				if jsonFlag {
					outputJSONError(fmt.Sprintf("check folder: %v", err))
				} else {
					fmt.Printf("check folder: %v\n", err)
				}
				os.Exit(exitorrors.ExitGeneral)
			}
		} else {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("folder already exists: %s", newFolder))
			} else {
				fmt.Printf("folder already exists: %s\n", newFolder)
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		parentDir := filepath.Dir(newValidation.AbsFolderPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("create parent directory: %v", err))
			} else {
				fmt.Printf("create parent directory: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if err := os.Rename(oldValidation.FolderPath, newValidation.FolderPath); err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("rename folder: %v", err))
			} else {
				fmt.Printf("rename folder: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if jsonFlag {
			response := map[string]string{
				"status": "renamed",
				"old":    oldFolder,
				"new":    newFolder,
			}
			outputJSON(response)
		} else {
			fmt.Printf("Renamed folder: %s → %s\n", oldFolder, newFolder)
		}
	},
}

func init() {
	folderCmd.AddCommand(folderRenameCmd)
	RootCmd.AddCommand(folderCmd)
}
