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
					outputJSONError(fmt.Sprintf("validate path: %v", err))
				} else {
					fmt.Printf("validate path: %v\n", err)
				}
				os.Exit(exitorrors.ExitGeneral)
			}
			targetPath = validation.FolderPath
		}

		info, err := os.Stat(targetPath)
		if os.IsNotExist(err) {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("path not found: %s", userPathArg))
			} else {
				fmt.Printf("path not found: %s\n", userPathArg)
			}
			os.Exit(exitorrors.ExitNotFound)
		}
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("check path: %v", err))
			} else {
				fmt.Printf("check path: %v\n", err)
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
				outputJSONError(fmt.Sprintf("read directory: %v", err))
			} else {
				fmt.Printf("read directory: %v\n", err)
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
				outputJSONError(fmt.Sprintf("marshal JSON: %v", err))
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
