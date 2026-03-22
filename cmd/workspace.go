package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
				errorResponse := map[string]string{"error": fmt.Sprintf("Error creating workspace manager: %v", err)}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
			} else {
				fmt.Printf("Error creating workspace manager: %v\n", err)
			}
			os.Exit(1)
		}
		if err := wm.Init(); err != nil {
			if errors.Is(err, workspace.ErrWorkspaceAlreadyExists) {
				if jsonFlag {
					errorResponse := map[string]string{"error": "workspace already initialized"}
					jsonError, _ := json.Marshal(errorResponse)
					fmt.Println(string(jsonError))
				} else {
					fmt.Println("Error: workspace already initialized")
				}
			} else {
				if jsonFlag {
					errorResponse := map[string]string{"error": fmt.Sprintf("Error initializing workspace: %v", err)}
					jsonError, _ := json.Marshal(errorResponse)
					fmt.Println(string(jsonError))
				} else {
					fmt.Printf("Error initializing workspace: %v\n", err)
				}
			}
			os.Exit(1)
		}
		if jsonFlag {
			response := map[string]string{"message": "Initialized workspace: default"}
			jsonOutput, _ := json.Marshal(response)
			fmt.Println(string(jsonOutput))
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
				errorResponse := map[string]string{"error": fmt.Sprintf("Error creating workspace manager: %v", err)}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
			} else {
				fmt.Printf("Error creating workspace manager: %v\n", err)
			}
			os.Exit(1)
		}
		workspaces, current, err := wm.List()
		if err != nil {
			if jsonFlag {
				errorResponse := map[string]string{"error": fmt.Sprintf("Error listing workspaces: %v", err)}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
			} else {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		if jsonFlag {
			var wsList []map[string]interface{}
			for name, path := range workspaces {
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
			jsonOutput, _ := json.Marshal(wsList)
			fmt.Println(string(jsonOutput))
		} else {
			fmt.Printf("Found %d workspaces\n", len(workspaces))
			for name, path := range workspaces {
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
				errorResponse := map[string]string{"error": fmt.Sprintf("Error creating workspace manager: %v", err)}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
			} else {
				fmt.Printf("Error creating workspace manager: %v\n", err)
			}
			os.Exit(1)
		}
		if err := wm.Use(name); err != nil {
			if errors.Is(err, workspace.ErrWorkspaceDoesNotExist) {
				if jsonFlag {
					errorResponse := map[string]string{"error": fmt.Sprintf("workspace '%s' not found", name)}
					jsonError, _ := json.Marshal(errorResponse)
					fmt.Println(string(jsonError))
				} else {
					fmt.Printf("Error: workspace '%s' not found\n", name)
				}
			} else {
				if jsonFlag {
					errorResponse := map[string]string{"error": fmt.Sprintf("Error using workspace: %v", err)}
					jsonError, _ := json.Marshal(errorResponse)
					fmt.Println(string(jsonError))
				} else {
					fmt.Printf("Error using workspace: %v\n", err)
				}
			}
			os.Exit(1)
		}
		if jsonFlag {
			response := map[string]string{"message": fmt.Sprintf("Switched to workspace: %s", name)}
			jsonOutput, _ := json.Marshal(response)
			fmt.Println(string(jsonOutput))
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
				errorResponse := map[string]string{"error": "workspace name is required"}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
			} else {
				fmt.Println("Error: workspace name is required")
			}
			os.Exit(1)
		}
		name := args[0]

		path := ""
		if len(args) > 1 {
			path = args[1]
		}

		wm, err := workspace.NewManager()
		if err != nil {
			if jsonFlag {
				errorResponse := map[string]string{"error": fmt.Sprintf("Error creating workspace manager: %v", err)}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
			} else {
				fmt.Printf("Error creating workspace manager: %v\n", err)
			}
			os.Exit(1)
		}
		if err := wm.New(name, path); err != nil {
			if errors.Is(err, workspace.ErrWorkspaceAlreadyExists) {
				if jsonFlag {
					errorResponse := map[string]string{"error": fmt.Sprintf("workspace '%s' already exists", name)}
					jsonError, _ := json.Marshal(errorResponse)
					fmt.Println(string(jsonError))
				} else {
					fmt.Printf("Error: workspace '%s' already exists\n", name)
				}
			} else {
				if jsonFlag {
					errorResponse := map[string]string{"error": fmt.Sprintf("Error creating workspace: %v", err)}
					jsonError, _ := json.Marshal(errorResponse)
					fmt.Println(string(jsonError))
				} else {
					fmt.Printf("Error creating workspace: %v\n", err)
				}
			}
			os.Exit(1)
		}
		if jsonFlag {
			response := map[string]string{"message": fmt.Sprintf("Created workspace: %s", name)}
			jsonOutput, _ := json.Marshal(response)
			fmt.Println(string(jsonOutput))
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
