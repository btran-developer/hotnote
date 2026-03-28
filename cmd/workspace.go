package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	exitorrors "hotnotego/internal/errors"
	"hotnotego/internal/workspace"
)

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Workspace management",
	Long:  `Manage hotnote workspaces`,
}

var workspaceInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the default workspace",
	Run: func(cmd *cobra.Command, args []string) {
		wm, err := workspace.NewManager()
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("create workspace manager: %v", err))
			} else {
				fmt.Printf("create workspace manager: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}
		if err := wm.Init(); err != nil {
			if errors.Is(err, workspace.ErrWorkspaceAlreadyExists) {
				if jsonFlag {
					outputJSONError("workspace already initialized")
				} else {
					fmt.Println("workspace already initialized")
				}
			} else {
				if jsonFlag {
					outputJSONError(fmt.Sprintf("init workspace: %v", err))
				} else {
					fmt.Printf("init workspace: %v\n", err)
				}
			}
			os.Exit(exitorrors.ExitGeneral)
		}
		if jsonFlag {
			response := map[string]string{"message": "Initialized workspace: default"}
			outputJSON(response)
		} else {
			fmt.Println("Initialized workspace: default")
		}
	},
}

var workspaceListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all workspaces",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		wm, err := workspace.NewManager()
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("create workspace manager: %v", err))
			} else {
				fmt.Printf("create workspace manager: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}
		workspaces, current, err := wm.List()
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("list workspaces: %v", err))
			} else {
				fmt.Printf("list workspaces: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		// Get sorted workspace names for deterministic output
		workspaceNames := make([]string, 0, len(workspaces))
		for name := range workspaces {
			workspaceNames = append(workspaceNames, name)
		}
		sort.Strings(workspaceNames)

		if jsonFlag {
			var wsList []map[string]interface{}
			for _, name := range workspaceNames {
				path := workspaces[name]
				ws := map[string]interface{}{
					"name": name,
					"path": path,
				}
				if name == current {
					ws["current"] = true
				} else {
					ws["current"] = false
				}
				wsList = append(wsList, ws)
			}
			if err := outputJSON(wsList); err != nil {
				outputJSONError(fmt.Sprintf("marshal JSON: %v", err))
			}
		} else {
			fmt.Printf("Found %d workspaces\n", len(workspaces))
			for _, name := range workspaceNames {
				path := workspaces[name]
				if name == current {
					fmt.Printf("* %s\t%s\n", name, path)
				} else {
					fmt.Printf("  %s\t%s\n", name, path)
				}
			}
		}
	},
}

var workspaceUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Switch to a workspace",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		wm, err := workspace.NewManager()
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("create workspace manager: %v", err))
			} else {
				fmt.Printf("create workspace manager: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}
		if err := wm.Use(name); err != nil {
			if errors.Is(err, workspace.ErrWorkspaceDoesNotExist) {
				if jsonFlag {
					outputJSONError("workspace not found")
				} else {
					fmt.Println("workspace not found")
				}
				os.Exit(exitorrors.ExitNotFound)
			} else {
				if jsonFlag {
					outputJSONError(fmt.Sprintf("use workspace: %v", err))
				} else {
					fmt.Printf("use workspace: %v\n", err)
				}
				os.Exit(exitorrors.ExitGeneral)
			}
		}
		if jsonFlag {
			response := map[string]string{"message": fmt.Sprintf("Switched to workspace: %s", name)}
			outputJSON(response)
		} else {
			fmt.Printf("Switched to workspace: %s\n", name)
		}
	},
}

var workspaceNewCmd = &cobra.Command{
	Use:     "create <name> [--path <path>]",
	Short:   "Create a new workspace",
	Aliases: []string{"new"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			if jsonFlag {
				outputJSONError("workspace name required")
			} else {
				fmt.Println("workspace name required")
			}
			os.Exit(exitorrors.ExitGeneral)
		}
		name := args[0]

		wm, err := workspace.NewManager()
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("create workspace manager: %v", err))
			} else {
				fmt.Printf("create workspace manager: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}
		if err := wm.New(name, workspaceNewPath); err != nil {
			if errors.Is(err, workspace.ErrWorkspaceAlreadyExists) {
				if jsonFlag {
					outputJSONError("workspace already exists")
				} else {
					fmt.Println("workspace already exists")
				}
			} else {
				if jsonFlag {
					outputJSONError(fmt.Sprintf("create workspace: %v", err))
				} else {
					fmt.Printf("create workspace: %v\n", err)
				}
			}
			os.Exit(exitorrors.ExitGeneral)
		}
		if jsonFlag {
			response := map[string]string{"message": fmt.Sprintf("Created workspace: %s", name)}
			outputJSON(response)
		} else {
			fmt.Printf("Created workspace: %s\n", name)
		}
	},
}

var workspaceDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a workspace or clear default workspace contents",
	Long: `Delete a workspace entirely, or clear contents of the default workspace.

For non-default workspaces:
  - Deletes the workspace directory and all contents
  - Removes the workspace from configuration
  - Cannot delete the currently active workspace (switch first)

For the default workspace:
  - Clears all contents (files and folders) within the workspace
  - Preserves the workspace structure and configuration
  - Can be used even when default is the current workspace

Examples:
  hotnote workspace delete work              # Delete 'work' workspace
  hotnote workspace delete default           # Clear all contents in default
  hotnote workspace delete work --force      # Skip confirmation prompt`,
	Aliases: []string{"del"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		wm, err := workspace.NewManager()
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("create workspace manager: %v", err))
			} else {
				fmt.Printf("create workspace manager: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		isDefault := name == "default"

		if isDefault {
			if !workspaceDeleteForce && jsonFlag {
				outputJSONError("use --force to delete")
				os.Exit(exitorrors.ExitGeneral)
			}

			if !workspaceDeleteForce {
				fmt.Printf("Clear all contents in default workspace? [y/N]: ")
				reader := bufio.NewReader(os.Stdin)
				input, err := reader.ReadString('\n')
				if err != nil {
					os.Exit(exitorrors.ExitGeneral)
				}
				input = strings.TrimSpace(input)
				if input != "y" && input != "Y" {
					os.Exit(exitorrors.ExitGeneral)
				}
				fmt.Println("Note: The default workspace structure will remain")
			}

			if err := wm.ClearDefaultWorkspace(); err != nil {
				if jsonFlag {
					outputJSONError(fmt.Sprintf("clear default workspace: %v", err))
				} else {
					fmt.Printf("clear default workspace: %v\n", err)
				}
				os.Exit(exitorrors.ExitGeneral)
			}

			if jsonFlag {
				response := map[string]string{
					"status":    "cleared",
					"workspace": "default",
				}
				outputJSON(response)
			} else {
				fmt.Println("Cleared default workspace")
			}
		} else {
			exists, err := wm.Exists(name)
			if err != nil {
				if jsonFlag {
					outputJSONError(fmt.Sprintf("check workspace: %v", err))
				} else {
					fmt.Printf("check workspace: %v\n", err)
				}
				os.Exit(exitorrors.ExitGeneral)
			}
			if !exists {
				if jsonFlag {
					outputJSONError("workspace not found")
				} else {
					fmt.Println("workspace not found")
				}
				os.Exit(exitorrors.ExitNotFound)
			}

			if !workspaceDeleteForce && jsonFlag {
				outputJSONError("use --force to delete")
				os.Exit(exitorrors.ExitGeneral)
			}

			if !workspaceDeleteForce {
				fmt.Printf("Delete workspace '%s'? [y/N]: ", name)
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

			if err := wm.Delete(name); err != nil {
				if errors.Is(err, workspace.ErrCannotDeleteCurrent) {
					if jsonFlag {
						outputJSONError("cannot delete current workspace: switch to another workspace first")
					} else {
						fmt.Println("cannot delete current workspace: switch to another workspace first")
					}
					os.Exit(exitorrors.ExitInvalidInput)
				}
				if jsonFlag {
					outputJSONError(fmt.Sprintf("delete workspace: %v", err))
				} else {
					fmt.Printf("delete workspace: %v\n", err)
				}
				os.Exit(exitorrors.ExitGeneral)
			}

			if jsonFlag {
				response := map[string]string{
					"status":    "deleted",
					"workspace": name,
				}
				outputJSON(response)
			} else {
				fmt.Printf("Deleted workspace: %s\n", name)
			}
		}
	},
}

var workspaceRenameCmd = &cobra.Command{
	Use:     "rename <old> <new>",
	Short:   "Rename a workspace",
	Aliases: []string{"rn"},
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		oldName := args[0]
		newName := args[1]

		wm, err := workspace.NewManager()
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("create workspace manager: %v", err))
			} else {
				fmt.Printf("create workspace manager: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if err := wm.Rename(oldName, newName); err != nil {
			if errors.Is(err, workspace.ErrWorkspaceDoesNotExist) {
				if jsonFlag {
					outputJSONError("workspace not found")
				} else {
					fmt.Println("workspace not found")
				}
				os.Exit(exitorrors.ExitNotFound)
			}
			if errors.Is(err, workspace.ErrWorkspaceAlreadyExists) {
				if jsonFlag {
					outputJSONError("workspace already exists")
				} else {
					fmt.Println("workspace already exists")
				}
				os.Exit(exitorrors.ExitInvalidInput)
			}
			if errors.Is(err, workspace.ErrCannotRenameDefault) {
				if jsonFlag {
					outputJSONError("cannot rename default workspace")
				} else {
					fmt.Println("cannot rename default workspace")
				}
				os.Exit(exitorrors.ExitInvalidInput)
			}
			if jsonFlag {
				outputJSONError(fmt.Sprintf("rename workspace: %v", err))
			} else {
				fmt.Printf("rename workspace: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if jsonFlag {
			response := map[string]string{
				"status": "renamed",
				"old":    oldName,
				"new":    newName,
			}
			outputJSON(response)
		} else {
			fmt.Printf("Renamed workspace: %s → %s\n", oldName, newName)
		}
	},
}

var workspaceDeleteForce bool
var workspaceNewPath string

func init() {
	workspaceCmd.AddCommand(workspaceInitCmd)
	workspaceCmd.AddCommand(workspaceListCmd)
	workspaceCmd.AddCommand(workspaceUseCmd)
	workspaceCmd.AddCommand(workspaceNewCmd)
	workspaceCmd.AddCommand(workspaceDeleteCmd)
	workspaceCmd.AddCommand(workspaceRenameCmd)
	workspaceNewCmd.Flags().StringVar(&workspaceNewPath, "path", "", "Custom path for workspace")
	workspaceDeleteCmd.Flags().BoolVar(&workspaceDeleteForce, "force", false, "Skip confirmation prompt")
	RootCmd.AddCommand(workspaceCmd)
}
