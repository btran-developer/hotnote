package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"hotnotego/internal/workspace"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
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

		_, workspacePath, err := wm.Current()
		if err != nil {
			if jsonFlag {
				errorResponse := map[string]string{"error": fmt.Sprintf("Error getting current workspace: %v", err)}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
			} else {
				fmt.Printf("Error getting current workspace: %v\n", err)
			}
			os.Exit(1)
		}

		files, err := os.ReadDir(workspacePath)
		if err != nil {
			if jsonFlag {
				errorResponse := map[string]string{"error": fmt.Sprintf("Error reading notes directory: %v", err)}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
			} else {
				fmt.Printf("Error reading notes directory: %v\n", err)
			}
			os.Exit(1)
		}

		if jsonFlag {
			var notes []map[string]string
			for _, file := range files {
				if !file.IsDir() && filepath.Ext(file.Name()) == ".md" {
					slug := file.Name()[:len(file.Name())-3] // Remove .md extension
					// Get file info for timestamp
					info, err := file.Info()
					if err != nil {
						continue
					}
					notes = append(notes, map[string]string{
						"slug":       slug,
						"path":       filepath.Join(workspacePath, file.Name()),
						"updated_at": info.ModTime().UTC().Format(time.RFC3339),
					})
				}
			}
			jsonData, err := json.Marshal(notes)
			if err != nil {
				errorResponse := map[string]string{"error": fmt.Sprintf("Error marshaling JSON: %v", err)}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
			} else {
				fmt.Println(string(jsonData))
			}
		} else {
			for _, file := range files {
				if !file.IsDir() && filepath.Ext(file.Name()) == ".md" {
					fmt.Println(file.Name())
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
