package cmd

import (
	"errors"
	"fmt"
	"os"
	"sort"

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
	Use:   "list",
	Short: "List all workspaces",
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
	Use:   "new <name> [--path <path>]",
	Short: "Create a new workspace",
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

		path := ""
		if len(args) > 1 {
			path = args[1]
		}

		wm, err := workspace.NewManager()
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("create workspace manager: %v", err))
			} else {
				fmt.Printf("create workspace manager: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}
		if err := wm.New(name, path); err != nil {
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

func init() {
	workspaceCmd.AddCommand(workspaceInitCmd)
	workspaceCmd.AddCommand(workspaceListCmd)
	workspaceCmd.AddCommand(workspaceUseCmd)
	workspaceCmd.AddCommand(workspaceNewCmd)
	RootCmd.AddCommand(workspaceCmd)
}
