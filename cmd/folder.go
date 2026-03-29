package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	exitorrors "hotnotego/internal/errors"
	"hotnotego/internal/pathutil"
	slugifypkg "hotnotego/internal/slugify"
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
		newFolder := slugifypkg.SlugifyPath(args[1])
		if newFolder == "" {
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
				outputJSONError(exitorrors.ErrInvalidFolderPath.Error())
			} else {
				fmt.Println(exitorrors.ErrInvalidFolderPath.Error())
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
				outputJSONError(exitorrors.ErrInvalidFolderPath.Error())
			} else {
				fmt.Println(exitorrors.ErrInvalidFolderPath.Error())
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if _, err := os.Stat(newValidation.FolderPath); err != nil {
			if !os.IsNotExist(err) {
				if jsonFlag {
					outputJSONError(exitorrors.ErrFolderRead.Error())
				} else {
					fmt.Println(exitorrors.ErrFolderRead.Error())
				}
				os.Exit(exitorrors.ExitGeneral)
			}
		} else {
			if jsonFlag {
				outputJSONError(exitorrors.ErrFolderExists.Error())
			} else {
				fmt.Println(exitorrors.ErrFolderExists.Error())
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		parentDir := filepath.Dir(newValidation.AbsFolderPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			if jsonFlag {
				outputJSONError(exitorrors.ErrFolderCreate.Error())
			} else {
				fmt.Println(exitorrors.ErrFolderCreate.Error())
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if err := os.Rename(oldValidation.FolderPath, newValidation.FolderPath); err != nil {
			if jsonFlag {
				outputJSONError(exitorrors.ErrFolderRename.Error())
			} else {
				fmt.Println(exitorrors.ErrFolderRename.Error())
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
