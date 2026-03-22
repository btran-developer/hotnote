package cmd

import (
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
			fmt.Printf("Error creating workspace manager: %v\n", err)
			os.Exit(1)
		}
		if err := wm.Init(); err != nil {
			if errors.Is(err, workspace.ErrWorkspaceAlreadyExists) {
				fmt.Println("Error: workspace already initialized")
			} else {
				fmt.Printf("Error initializing workspace: %v\n", err)
			}
			os.Exit(1)
		}
		fmt.Println("Initialized workspace: default")
	},
}

var workspaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workspaces",
	Run: func(cmd *cobra.Command, args []string) {
		wm, err := workspace.NewManager()
		if err != nil {
			fmt.Printf("Error creating workspace manager: %v\n", err)
			os.Exit(1)
		}
		workspaces, current, err := wm.List()
		if err != nil {
			fmt.Printf("Error listing workspaces: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Found %d workspaces\n", len(workspaces)) // Debug line
		for name, path := range workspaces {
			if name == current {
				fmt.Printf("* %s\t%s\n", name, path)
			} else {
				fmt.Printf("  %s\t%s\n", name, path)
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
			fmt.Printf("Error creating workspace manager: %v\n", err)
			os.Exit(1)
		}
		if err := wm.Use(name); err != nil {
			if errors.Is(err, workspace.ErrWorkspaceDoesNotExist) {
				fmt.Printf("Error: workspace '%s' not found\n", name)
			} else {
				fmt.Printf("Error using workspace: %v\n", err)
			}
			os.Exit(1)
		}
		fmt.Printf("Switched to workspace: %s\n", name)
	},
}

var workspaceNewCmd = &cobra.Command{
	Use:   "new <name> [--path <path>]",
	Short: "Create a new workspace",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: workspace name is required")
			os.Exit(1)
		}
		name := args[0]

		path := ""
		if len(args) > 1 {
			path = args[1]
		}

		wm, err := workspace.NewManager()
		if err != nil {
			fmt.Printf("Error creating workspace manager: %v\n", err)
			os.Exit(1)
		}
		if err := wm.New(name, path); err != nil {
			if errors.Is(err, workspace.ErrWorkspaceAlreadyExists) {
				fmt.Printf("Error: workspace '%s' already exists\n", name)
			} else {
				fmt.Printf("Error creating workspace: %v\n", err)
			}
			os.Exit(1)
		}
		fmt.Printf("Created workspace: %s\n", name)
	},
}

func init() {
	workspaceCmd.AddCommand(workspaceInitCmd)
	workspaceCmd.AddCommand(workspaceListCmd)
	workspaceCmd.AddCommand(workspaceUseCmd)
	workspaceCmd.AddCommand(workspaceNewCmd)
	RootCmd.AddCommand(workspaceCmd)
}
