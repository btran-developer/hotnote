package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
	exitorrors "hotnotego/internal/errors"
	"hotnotego/internal/pathutil"
	"hotnotego/internal/workspace"
)

var folderListCmd = &cobra.Command{
	Use:     "list [path]",
	Short:   "List files and folders in a directory",
	Aliases: []string{"ls"},
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		wm, err := workspace.NewManager()
		if err != nil {
			if jsonFlag {
				outputJSONError(exitorrors.ErrWorkspaceNotInit.Error())
			} else {
				fmt.Println(exitorrors.ErrWorkspaceNotInit.Error())
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

		targetPath := wsPath
		userPathArg := "."
		if len(args) > 0 {
			userPathArg = args[0]
			validation, err := pathutil.ValidateFolderPath(wsPath, userPathArg)
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
			targetPath = validation.FolderPath
		}

		info, err := os.Stat(targetPath)
		if os.IsNotExist(err) {
			if jsonFlag {
				outputJSONError(exitorrors.ErrFolderNotFound.Error())
			} else {
				fmt.Println(exitorrors.ErrFolderNotFound.Error())
			}
			os.Exit(exitorrors.ExitNotFound)
		}
		if err != nil {
			if jsonFlag {
				outputJSONError(exitorrors.ErrFolderRead.Error())
			} else {
				fmt.Println(exitorrors.ErrFolderRead.Error())
			}
			os.Exit(exitorrors.ExitGeneral)
		}
		if !info.IsDir() {
			if jsonFlag {
				outputJSONError("not a directory")
			} else {
				fmt.Println("not a directory")
			}
			os.Exit(exitorrors.ExitInvalidInput)
		}

		entries, err := os.ReadDir(targetPath)
		if err != nil {
			if jsonFlag {
				outputJSONError(exitorrors.ErrFolderRead.Error())
			} else {
				fmt.Println(exitorrors.ErrFolderRead.Error())
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		type entry struct {
			Name string `json:"name"`
			Path string `json:"path"`
			Type string `json:"type"`
		}

		var items []entry
		var folders []entry
		var files []entry

		for _, e := range entries {
			relPath, err := filepath.Rel(wsPath, filepath.Join(targetPath, e.Name()))
			if err != nil {
				continue
			}

			item := entry{
				Name: e.Name(),
				Path: relPath,
			}
			if e.IsDir() {
				item.Type = "folder"
				folders = append(folders, item)
			} else {
				item.Type = "file"
				files = append(files, item)
			}
		}

		sort.Slice(folders, func(i, j int) bool { return folders[i].Name < folders[j].Name })
		sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })

		items = append(folders, files...)

		if jsonFlag {
			if err := outputJSON(items); err != nil {
				outputJSONError(exitorrors.ErrMarshalJSON.Error())
				os.Exit(exitorrors.ExitGeneral)
			}
		} else {
			if len(items) == 0 {
				fmt.Printf("Empty directory: %s\n", targetPath)
				return
			}
			for _, item := range items {
				if item.Type == "folder" {
					fmt.Printf("%s/\n", item.Name)
				} else {
					fmt.Printf("%s\n", item.Name)
				}
			}
		}
	},
}

func init() {
	folderCmd.AddCommand(folderListCmd)
}
