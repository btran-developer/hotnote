package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	exitorrors "hotnotego/internal/errors"
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
				errorResponse := map[string]string{"error": fmt.Sprintf("create workspace manager: %v", err)}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
			} else {
				fmt.Printf("create workspace manager: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		store := storage.NewStore(wm)
		path, err := store.Path(title)
		if err != nil {
			if jsonFlag {
				errorResponse := map[string]string{"error": fmt.Sprintf("get note path: %v", err)}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
			} else {
				fmt.Printf("get note path: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			if jsonFlag {
				errorResponse := map[string]string{"error": "note not found"}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
			} else {
				fmt.Println("note not found")
			}
			os.Exit(exitorrors.ExitNotFound)
		}
		var buf bytes.Buffer
		err = md.Convert(content, &buf)
		if err != nil {
			if jsonFlag {
				errorResponse := map[string]string{"error": fmt.Sprintf("render markdown: %v", err)}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
			} else {
				fmt.Printf("render markdown: %v\n", err)
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if jsonFlag {
			jsonResponse := map[string]string{"content": buf.String()}
			jsonData, err := json.Marshal(jsonResponse)
			if err != nil {
				errorResponse := map[string]string{"error": fmt.Sprintf("marshal JSON: %v", err)}
				jsonError, _ := json.Marshal(errorResponse)
				fmt.Println(string(jsonError))
				os.Exit(exitorrors.ExitGeneral)
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
