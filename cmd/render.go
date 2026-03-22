package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"hotnotego/internal/storage"
	"hotnotego/internal/workspace"
)

var md = goldmark.New()

var renderCmd = &cobra.Command{
	Use:   "render [title]",
	Short: "Render markdown to HTML",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
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

		store := storage.NewStore(wm)
		path, err := store.Path(title)
		if err != nil {
			if jsonFlag {
				errorResponse := map[string]string{"error": fmt.Sprintf("Error getting note path: %v", err)}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
			} else {
				fmt.Printf("Error getting note path: %v\n", err)
			}
			os.Exit(1)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			if jsonFlag {
				errorResponse := map[string]string{"error": fmt.Sprintf("Error reading note: %v", err)}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
			} else {
				fmt.Printf("Error reading note: %v\n", err)
			}
			os.Exit(1)
		}
		var buf bytes.Buffer
		err = md.Convert(content, &buf)
		if err != nil {
			if jsonFlag {
				errorResponse := map[string]string{"error": fmt.Sprintf("Error rendering markdown: %v", err)}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
			} else {
				fmt.Printf("Error rendering markdown: %v\n", err)
			}
			os.Exit(1)
		}

		if jsonFlag {
			jsonResponse := map[string]string{"content": buf.String()}
			jsonData, err := json.Marshal(jsonResponse)
			if err != nil {
				errorResponse := map[string]string{"error": fmt.Sprintf("Error marshaling JSON: %v", err)}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
				os.Exit(1)
			} else {
				fmt.Println(string(jsonData))
			}
		} else {
			fmt.Println(buf.String())
		}
	},
}

func init() {
	RootCmd.AddCommand(renderCmd)
}
